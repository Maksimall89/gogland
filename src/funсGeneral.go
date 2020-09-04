package src

import (
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"math"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func LogInit() {
	path := "log" // name folder for logs
	// check what folder log is exist
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(path, os.ModePerm)
	}
	dateTime := time.Now()
	path = fmt.Sprintf("%s/%s-logFile.log", path, dateTime.Format("02.01.2006-15.04.05.000"))
	// configurator for logger
	// open a file
	fileLog, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("Error opening file: %v", err)
	}
	// assign it to the standard logger
	log.SetOutput(fileLog)
	log.SetPrefix("Gogland ")
}
func SendMessageTelegram(chatId int64, message string, replyToMessageID int, bot *tgbotapi.BotAPI) error {
	var pointerStr int
	var msg tgbotapi.MessageConfig
	var newMsg tgbotapi.Message
	var err error
	isEnd := false

	if len(message) == 0 {
		return nil
	}
	if replyToMessageID != 0 {
		msg.ReplyToMessageID = replyToMessageID
	}
	msg.ChatID = chatId
	msg.ParseMode = "HTML"

	for !isEnd {
		if len(message) > 4090 { // ограничение на длину одного сообщения 4096
			pointerStr = strings.LastIndex(message[0:4090], "\n")
			msg.Text = message[0:pointerStr]
			message = message[pointerStr:]
		} else {
			msg.Text = message
			isEnd = true
		}

		newMsg, err = bot.Send(msg)
		if err != nil {
			msg.ParseMode = "Markdown"
			newMsg, err = bot.Send(msg)
			if err != nil {
				log.Println(err)
				log.Println(msg.Text)
			}
			msg.ParseMode = "HTML"
		}
		if strings.Contains(msg.Text, "&#9889;Выдан новый уровень!") {
			_, err := bot.PinChatMessage(tgbotapi.PinChatMessageConfig{ChatID: chatId, MessageID: newMsg.MessageID})
			if err != nil {
				log.Println(err)
				return err
			}
		}
	}
	return nil
}
func SearchLocation(text string) map[string]float64 {
	maps := make(map[string]float64)
	re := regexp.MustCompile(`(\d{2}[.|°,]+\d{1,8}(['|′.]\d{1,})*([.|″]\d)*["|″]*)[,|\w\s\n<brBR/>]+(\d{2}[.|°,]+\d{1,8}(['|′.]\d{1,})*([.|″]\d)*["|″]*)`)

	if re.MatchString(text) {
		var degree [2]float64
		var min [2]float64
		var sec [2]float64
		var counter int
		var arrSeparator [2]string

		strCoordinates := re.FindAllString(text, -1)
		regPattern1, _ := regexp.Compile(`(\d{2}[.|,]+\d{2,8})[,|\s\n\w<brBR/>]+(\d{2}[.|,]+\d{2,8})`)
		regPattern2, _ := regexp.Compile(`(\d{2}°+\d+\.+\d{1,8})[,\s\n\w<brBR/>]+(\d{2}°+\d+\.+\d{1,8})`)
		regPattern3, _ := regexp.Compile(`(\d{2}°+\d{1,8}'+\d[.\d{1,8}]*)["]*[,|\s\n\w<brBR/>]+(\d{2}°+\d{1,8}'+\d[.\d{1,8}]*)["]*`)

		for _, strCoordinate := range strCoordinates {
			// 40.167841 58.410761 or 40,167841 58,410761
			if regPattern1.MatchString(strCoordinate) {
				sNew := regPattern1.FindStringSubmatch(strCoordinate)
				sNew[1] = strings.ReplaceAll(sNew[1], ",", ".")
				sNew[2] = strings.ReplaceAll(sNew[2], ",", ".")

				maps["Latitude"+strconv.FormatInt(int64(counter), 10)], _ = strconv.ParseFloat(sNew[1], 64)
				maps["Longitude"+strconv.FormatInt(int64(counter), 10)], _ = strconv.ParseFloat(sNew[2], 64)

				counter++
				continue
			}

			//  56°50.683, 53°11.776
			if regPattern2.MatchString(strCoordinate) {
				sNew := regPattern2.FindStringSubmatch(strCoordinate)
				buf1 := strings.Split(sNew[1], `°`)
				buf2 := strings.Split(sNew[2], `°`)

				degree[0], _ = strconv.ParseFloat(buf1[0], 64)
				degree[1], _ = strconv.ParseFloat(buf2[0], 64)

				min[0], _ = strconv.ParseFloat(buf1[1], 64)
				min[1], _ = strconv.ParseFloat(buf2[1], 64)

				maps["Latitude"+strconv.FormatInt(int64(counter), 10)] = degree[0] + ToFixed(min[0]/60, 8)
				maps["Longitude"+strconv.FormatInt(int64(counter), 10)] = degree[1] + ToFixed(min[1]/60, 8)

				counter++
				continue
			}

			// 59°24'48.6756" 58°24'38.7396"
			// 59°24'48.6756 58°24'38.7396
			// 55°10'11"N 52°12'01"E
			// 40°15′08″ 58°26′23″
			// 40°15′08 58°26′23
			if regPattern3.MatchString(strCoordinate) {
				sNew := regPattern3.FindStringSubmatch(strCoordinate)
				buf1 := strings.Split(sNew[1], `°`)
				buf2 := strings.Split(sNew[2], `°`)

				degree[0], _ = strconv.ParseFloat(buf1[0], 64)
				degree[1], _ = strconv.ParseFloat(buf2[0], 64)

				if strings.Contains(strCoordinate, `′`) {
					arrSeparator[0] = `′`
					arrSeparator[1] = `″`
				} else {
					arrSeparator[0] = `'`
					arrSeparator[1] = `"`
				}

				buf1 = strings.Split(buf1[1], arrSeparator[0])
				buf2 = strings.Split(buf2[1], arrSeparator[0])

				min[0], _ = strconv.ParseFloat(buf1[0], 64)
				min[1], _ = strconv.ParseFloat(buf2[0], 64)

				buf1 = strings.Split(buf1[1], arrSeparator[1])
				buf2 = strings.Split(buf2[1], arrSeparator[1])

				sec[0], _ = strconv.ParseFloat(buf1[0], 64)
				sec[1], _ = strconv.ParseFloat(buf2[0], 64)

				maps["Latitude"+strconv.FormatInt(int64(counter), 10)] = degree[0] + ToFixed(min[0]/60, 7) + ToFixed(sec[0]/3600, 7)
				maps["Longitude"+strconv.FormatInt(int64(counter), 10)] = degree[1] + ToFixed(min[1]/60, 7) + ToFixed(sec[1]/3600, 7)

				counter++
				continue
			}
		}
	}
	return maps
}
func SendLocation(locationMap map[string]float64, webToBot chan MessengerStyle) {
	if len(locationMap) < 1 {
		return
	}

	msgBot := MessengerStyle{}
	msgBot.ChannelMessage = ""

	for i := 0; i < len(locationMap)/2; i++ {
		msgBot.Latitude = locationMap["Latitude"+strconv.FormatInt(int64(i), 10)]
		msgBot.Longitude = locationMap["Longitude"+strconv.FormatInt(int64(i), 10)]
		msgBot.Type = "location"
		webToBot <- msgBot
	}
}
func SearchPhoto(text string) map[int]string {
	maps := make(map[int]string)
	regPattern, _ := regexp.Compile(`<img.+?src="(.+?)".*?>`)
	if regPattern.MatchString(text) {
		s := regPattern.FindAllString(text, -1)
		for i := 0; i < len(s); i++ {
			sNew := regPattern.FindStringSubmatch(s[i])
			maps[i] = sNew[1]
		}
	}
	return maps
}
func SendPhoto(photoMap map[int]string, webToBot chan MessengerStyle) {
	if len(photoMap) < 1 {
		return
	}

	msgBot := MessengerStyle{}
	msgBot.MsgId = 0
	for i := 0; i < len(photoMap); i++ {
		msgBot.ChannelMessage = photoMap[i]
		msgBot.Type = "photo"
		webToBot <- msgBot
	}
}
func ReplaceTag(str string, subUrl string) string {
	/*
		*полужирный*
		_курсив_
		[ссылка](http://www.example.com/)
		`строчный моноширинный`

		<b>полужирный</b>, <strong>полужирный</strong>
		<i>курсив</i>
		<a href="http://www.example.com/">ссылка</a>
		<code>строчный моноширинный</code>
		<pre>блочный моноширинный (можно писать код)</pre>
	*/
	if len(str) < 1 {
		return ""
	}

	assumeHtmlTag := []string{"<b>", "</b>", "<code>", "</code>", "<pre>", "</pre>", "</a>"}

	type pairStr struct {
		original string
		replaced string
	}

	var pairs = []pairStr{
		{`&nbsp;`, " "},
		{`&hellip;`, "..."},
		{`&mdash;`, "-"},
		{`&laquo;`, "«"},
		{`&raquo;`, "»"},
		{`(\n|\t)*`, ""},
		{`<style.*?</style>`, ""},
		{`<script.*?</script>`, "<b>[В тексте найден кусок скрипта, который был вырезан]</b>"},
		{`(\/\/<!\[CDATA\[(\D|\d)*\/\/\]\]>)`, ""},
		{`<\s*(b|B) .*?>`, "<b>"},
		{`<\s*\/(b|B) .*?>`, "</b>"},
		{`<\s*\/\s*(p|P)\s*>`, "\n"},
		{`<\s*(br|BR)\s*\/\s*>`, "\n"},
		{`(<\s*(h\d+|H\d+).*?>)`, "<b>"},
		{`(<\s*\/\s*(h\d+|H\d+)\s*>)`, "</b>\n"},
		{`<\s*(hr|HR)\s*\/*\s*>`, "\n------------\n"},
		{`<\s*(font|FONT).*?>`, `[цвет:]`},
		{`<\s*\/\s*(font|FONT).*?>`, `[:цвет]`},
		{`<\s*\/\s*(li|LI)\s*>`, "\n"},
		{`<\s*(sup|SUP)\s*>`, "^"},
		{`<\s*(sub|SUB)\s*>`, "_"},
		{`<\s*(small|SMALL)\s*>`, "^"},
		{`<b><b>`, "<b>"},
		{`</b></b>`, "</b>"},
		{`<\s*(iframe|IFRAME).+?src="(?P<link>.+?)".*?<\/(iframe|IFRAME)>`, `[Вставлен iframe с другого сайта: <a href="${link}">${link}</a>]`},
		{`<\s*(audio|AUDIO).+?src="(?P<link>.+?)".*?<\/(audio|AUDIO)>`, `[Вставлено audio с другого сайта: <a href="${link}">${link}</a>]`},
		{`<\s*a\s*(?i)href\s*=\s*(?P<link>\"[^"]*\"|'[^']*'|[^'">\s]+).*?>(?P<titleLink>.+?)<\s*/\s*[a|A]\s*>`, `<a href=${link}>${titleLink}</a>`},
		{`<\s*(?i)img.+?src="(?P<link>.+?)".*?>`, `[Картинка: <a href="${link}">${link}</a>]`},
		{`<a href="(?P<link>.*?)".*?>\n*\[Картинка: (?P<titleLink>.+?)]\n*[<\/a>]*`, `[Спрятанная ссылка: <a href="${link}">${link}</a> под картинкой: ${titleLink}]`},
		{`</a></a>`, `</a>`},
		{`</a>]</a>`, `</a>]`},
		{`(?i)<\s*img.+?src=(?P<link>.+?)>`, `[Картинка: <a href="${link}">${link}</a>]`},
	}

	for _, pair := range pairs {
		reg := regexp.MustCompile(pair.original)
		str = reg.ReplaceAllString(str, pair.replaced)
	}

	// delete unknown a html teg \<(\/?[^\>]+)\>
	reg := regexp.MustCompile(`(</?[^>]+>)`)
	tegs := reg.FindAllString(str, -1)

	eval := false
	for _, teg := range tegs {
		for _, repair := range assumeHtmlTag {
			if teg == repair {
				eval = true
				break
			}
		}
		if !eval && !strings.Contains(teg, `<a href="`) {
			str = strings.Replace(str, teg, "", -1)
		}
		eval = false
	}

	var strArr []string
	// add subUrl
	if strings.Contains(str, "<a href=\"/") {
		urlString := "\n<a href=\"http://" + subUrl + "/"
		str = strings.Replace(str, "<a href=\"/", urlString, -1)
	}

	// Coordinates
	reg = regexp.MustCompile(`[:|\s](\d{2}[.|°,]+\d{1,8}(['|′.]\d)*([.|″]\d)*["|″]*)[,|\w\s\n<brBR/>]+(\d{2}[.|°,]+\d{1,8}(['|′.]\d)*([.|″]\d)*["|″]*)[\s|,.]`)
	strArr = reg.FindAllString(str, -1)
	for _, coordinate := range strArr {
		finalCoordinates := SearchLocation(coordinate)
		for i := 0; i < len(finalCoordinates)/2; i++ {
			str = strings.ReplaceAll(str, coordinate, fmt.Sprintf(` <code>%f,%f</code> <a href="https://maps.google.com/?q=%f,%f">[G]</a> <a href="https://yandex.ru/maps/?source=serp_navig&text=%f,%f">[Y]</a>, `, finalCoordinates["Latitude"+strconv.FormatInt(int64(i), 10)], finalCoordinates["Latitude"+strconv.FormatInt(int64(i), 10)], finalCoordinates["Latitude"+strconv.FormatInt(int64(i), 10)], finalCoordinates["Latitude"+strconv.FormatInt(int64(i), 10)], finalCoordinates["Latitude"+strconv.FormatInt(int64(i), 10)], finalCoordinates["Latitude"+strconv.FormatInt(int64(i), 10)]))
		}
	}
	// Convert Links
	reg = regexp.MustCompile(`<a href="(.+?)"*>(.+?)</a>`)
	strArr = reg.FindAllString(str, -1)
	for _, link := range strArr {
		newLink, err := ConvertUTF(link)
		if err != nil {
			log.Printf("Error URL conver =%s, Text=%s", err, str)
			continue
		}
		str = strings.Replace(str, link, newLink, -1)
	}
	return str
}
func DeleteMapFloat(maps map[string]float64) {
	for key := range maps {
		delete(maps, key)
	}
}
func DeleteMap(maps map[int]string) {
	for key := range maps {
		delete(maps, key)
	}
}
func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(Round(num*output)) / output
}
func ConvertUTF(str string) (string, error) {
	u, err := url.QueryUnescape(str)
	if err != nil {
		log.Printf("\nURL %s, bad request %s", str, err)
		return str, err
	}
	return u, nil
}
func Round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}
func ConvertTimeSec(times int) string {
	if times == 0 {
		return "0 секунд"
	}
	// из секунд в минуту
	str := ""
	timeSec := times % 60
	timeMin := times / 60
	timeHour := times / 3600
	if timeHour > 0 {
		timeMin = timeMin % 60
	}
	timeDay := times / 86400
	if timeDay > 0 {
		timeMin = times % 60
		timeHour = timeHour % 24
	}

	// Дни
	switch timeDay {
	case 0:
		str += ""
	case 1:
		str += fmt.Sprintf("%d день ", timeDay)
	case 2, 3, 4:
		str += fmt.Sprintf("%d дня ", timeDay)
	default:
		str += fmt.Sprintf("%d дней ", timeDay)
	}
	// Часы
	switch timeHour {
	case 0:
		str += ""
	case 1:
		str += fmt.Sprintf("%d час ", timeHour)
	case 2, 3, 4:
		str += fmt.Sprintf("%d часа ", timeHour)
	default:
		str += fmt.Sprintf("%d часов ", timeHour)
	}
	// Минуты
	switch timeMin {
	case 0:
		str += ""
	case 1:
		str += fmt.Sprintf("%d минута ", timeMin)
	case 2, 3, 4:
		str += fmt.Sprintf("%d минуты ", timeMin)
	default:
		str += fmt.Sprintf("%d минут ", timeMin)
	}
	// Секунды
	switch timeSec {
	case 0:
		str += ""
	case 1:
		str += fmt.Sprintf("%d секунда", timeSec)
	case 2, 3, 4:
		str += fmt.Sprintf("%d секунды", timeSec)
	default:
		str += fmt.Sprintf("%d секунд", timeSec)
	}
	return strings.TrimSpace(str)
}

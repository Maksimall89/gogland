package main

import (
	"fmt"
	"log"
	"math"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

func searchLocation(text string) map[string]float64 {
	maps := make(map[string]float64)
	if reg, _ := regexp.MatchString(`(\d{2}[.|°,]+\d{1,8}(['|′.]\d{1,})*([.|″]\d)*["|″]*)[,|\w\s\n<brBR/>]+(\d{2}[.|°,]+\d{1,8}(['|′.]\d{1,})*([.|″]\d)*["|″]*)`, text); reg {
		re := regexp.MustCompile(`(\d{2}[.|°,]+\d{1,8}(['|′.]\d{1,})*([.|″]\d)*["|″]*)[,|\w\s\n<brBR/>]+(\d{2}[.|°,]+\d{1,8}(['|′.]\d{1,})*([.|″]\d)*["|″]*)`)
		s := re.FindAllString(text, -1)

		var degree [2]float64
		var min [2]float64
		var sec [2]float64
		var counter int

		for i := 0; i < len(s); i++ {
			// 40.167841 58.410761 or 40,167841 58,410761
			if reg, _ := regexp.MatchString(`(\d{2}[.|,]+\d{2,8})[,|\s\n\w<brBR/>]+(\d{2}[.|,]+\d{2,8})`, s[i]); reg {
				sNew := regexp.MustCompile(`(\d{2}[.|,]+\d{2,8})[,|\s\n\w<brBR/>]+(\d{2}[.|,]+\d{2,8})`).FindStringSubmatch(s[i])
				maps["Latitude"+strconv.FormatInt(int64(counter), 10)], _ = strconv.ParseFloat(sNew[1], 64)
				maps["Longitude"+strconv.FormatInt(int64(counter), 10)], _ = strconv.ParseFloat(sNew[2], 64)

				counter++
				continue
			}

			//  56°50.683, 53°11.776
			if reg, _ := regexp.MatchString(`(\d{2}°+\d+\.+\d{1,8})[,\s\n\w<brBR/>]+(\d{2}°+\d+\.+\d{1,8})`, s[i]); reg {
				sNew := regexp.MustCompile(`(\d{2}°+\d+\.+\d{1,8})[,\s\n\w<brBR/>]+(\d{2}°+\d+\.+\d{1,8})`).FindStringSubmatch(s[i])
				buf1 := strings.Split(sNew[1], `°`)
				buf2 := strings.Split(sNew[2], `°`)

				degree[0], _ = strconv.ParseFloat(buf1[0], 64)
				degree[1], _ = strconv.ParseFloat(buf2[0], 64)

				min[0], _ = strconv.ParseFloat(buf1[1], 64)
				min[1], _ = strconv.ParseFloat(buf2[1], 64)

				maps["Latitude"+strconv.FormatInt(int64(counter), 10)] = degree[0] + toFixed(min[0]/60, 8)
				maps["Longitude"+strconv.FormatInt(int64(counter), 10)] = degree[1] + toFixed(min[1]/60, 8)

				counter++
				continue
			}

			// 40°15′08″ 58°26′23″
			// 40°15′08 58°26′23
			if reg, _ := regexp.MatchString(`(\d{2}°+\d{1,8}′+\d{1,8}[″]*)[,|\s\n\w<brBR/>]+(\d{2}°+\d{1,8}′+\d{1,8}[″]*)`, s[i]); reg {
				sNew := regexp.MustCompile(`(\d{2}°+\d{1,8}′+\d{1,8}[″]*)[,|\s\n\w<brBR/>]+(\d{2}°+\d{1,8}′+\d{1,8}[″]*)`).FindStringSubmatch(s[i])
				buf1 := strings.Split(sNew[1], `°`)
				buf2 := strings.Split(sNew[2], `°`)

				degree[0], _ = strconv.ParseFloat(buf1[0], 64)
				degree[1], _ = strconv.ParseFloat(buf2[0], 64)

				buf1 = strings.Split(buf1[1], `′`)
				buf2 = strings.Split(buf2[1], `′`)

				min[0], _ = strconv.ParseFloat(buf1[0], 64)
				min[1], _ = strconv.ParseFloat(buf2[0], 64)

				buf1 = strings.Split(buf1[1], `″`)
				buf2 = strings.Split(buf2[1], `″`)

				sec[0], _ = strconv.ParseFloat(buf1[0], 64)
				sec[1], _ = strconv.ParseFloat(buf2[0], 64)

				maps["Latitude"+strconv.FormatInt(int64(counter), 10)] = degree[0] + toFixed(min[0]/60, 7) + toFixed(sec[0]/3600, 7)
				maps["Longitude"+strconv.FormatInt(int64(counter), 10)] = degree[1] + toFixed(min[1]/60, 7) + toFixed(sec[1]/3600, 7)

				counter++
				continue
			}

			// 59°24'48.6756" 58°24'38.7396"
			// 59°24'48.6756 58°24'38.7396
			if reg, _ := regexp.MatchString(`(\d{2}°+\d{2,8}'+\d{2}\.\d{1,8})["]*[,|\s\n\w<brBR/>]+(\d{2}°+\d{2}'+\d{2}\.\d{1,8})["]*`, s[i]); reg {
				sNew := regexp.MustCompile(`(\d{2}°+\d{2,8}'+\d{2}\.\d{1,8})["]*[,|\s\n\w<brBR/>]+(\d{2}°+\d{2}'+\d{2}\.\d{1,8})["]*`).FindStringSubmatch(s[i])
				buf1 := strings.Split(sNew[1], `°`)
				buf2 := strings.Split(sNew[2], `°`)

				degree[0], _ = strconv.ParseFloat(buf1[0], 64)
				degree[1], _ = strconv.ParseFloat(buf2[0], 64)

				buf1 = strings.Split(buf1[1], `'`)
				buf2 = strings.Split(buf2[1], `'`)

				min[0], _ = strconv.ParseFloat(buf1[0], 64)
				min[1], _ = strconv.ParseFloat(buf2[0], 64)

				buf1 = strings.Split(buf1[1], `"`)
				buf2 = strings.Split(buf2[1], `"`)

				sec[0], _ = strconv.ParseFloat(buf1[0], 64)
				sec[1], _ = strconv.ParseFloat(buf2[0], 64)

				maps["Latitude"+strconv.FormatInt(int64(counter), 10)] = degree[0] + toFixed(min[0]/60, 7) + toFixed(sec[0]/3600, 7)
				maps["Longitude"+strconv.FormatInt(int64(counter), 10)] = degree[1] + toFixed(min[1]/60, 7) + toFixed(sec[1]/3600, 7)

				counter++
				continue
			}

			// 55°10'11"N 52°12'01"E
			if reg, _ := regexp.MatchString(`(\d{2}°+\d{2,8}'+\d{2,}")[,|\s\n\w<brBR/>]+(\d{2}°+\d{2}'+\d{2,}")`, s[i]); reg {

				sNew := regexp.MustCompile(`(\d{2}°+\d{2,8}'+\d{2,}")[,|\s\n\w<brBR/>]+(\d{2}°+\d{2}'+\d{2,}")`).FindStringSubmatch(s[i])
				buf1 := strings.Split(sNew[1], `°`)
				buf2 := strings.Split(sNew[2], `°`)

				degree[0], _ = strconv.ParseFloat(buf1[0], 64)
				degree[1], _ = strconv.ParseFloat(buf2[0], 64)

				buf1 = strings.Split(buf1[1], `'`)
				buf2 = strings.Split(buf2[1], `'`)

				min[0], _ = strconv.ParseFloat(buf1[0], 64)
				min[1], _ = strconv.ParseFloat(buf2[0], 64)

				buf1 = strings.Split(buf1[1], `"`)
				buf2 = strings.Split(buf2[1], `"`)

				sec[0], _ = strconv.ParseFloat(buf1[0], 64)
				sec[1], _ = strconv.ParseFloat(buf2[0], 64)

				maps["Latitude"+strconv.FormatInt(int64(counter), 10)] = degree[0] + toFixed(min[0]/60, 7) + toFixed(sec[0]/3600, 7)
				maps["Longitude"+strconv.FormatInt(int64(counter), 10)] = degree[1] + toFixed(min[1]/60, 7) + toFixed(sec[1]/3600, 7)

				counter++
				continue
			}
		}
	}
	return maps
}
func sentLocation(locationMap map[string]float64, webToBot chan MessengerStyle) {

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
	return
}
func searchPhoto(text string) map[int]string {
	maps := make(map[int]string)
	if reg, _ := regexp.MatchString(`<img.+?src="(.+?)".*?>`, text); reg {
		s := regexp.MustCompile(`<img.+?src="(.+?)".*?>`).FindAllString(text, -1)
		for i := 0; i < len(s); i++ {
			sNew := regexp.MustCompile(`<img.+?src="(.+?)".*?>`).FindStringSubmatch(s[i])
			maps[i] = sNew[1]
		}
	}
	return maps
}
func sentPhoto(photoMap map[int]string, webToBot chan MessengerStyle) {

	if len(photoMap) < 1 {
		return
	}

	msgBot := MessengerStyle{}
	msgBot.ChannelMessage = ""
	msgBot.MsgId = 0
	for i := 0; i < len(photoMap); i++ {
		msgBot.ChannelMessage = photoMap[i]
		msgBot.Type = "photo"
		webToBot <- msgBot
	}
	return
}
func replaceTag(str string, subUrl string) string {
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
		if eval == false && !strings.Contains(teg, `<a href="`) {
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
	reg = regexp.MustCompile(`(\d{2}[.|°,]+\d{1,8}(['|′.]\d{1,})*([.|″]\d)*["|″]*)[,|\w\s\n<brBR/>]+(\d{2}[.|°,]+\d{1,8}(['|′.]\d{1,})*([.|″]\d)*["|″]*)`)
	strArr = reg.FindAllString(str, -1)
	for _, coordinate := range strArr {
		finalCoordinates := searchLocation(coordinate)
		str = strings.Replace(str, coordinate, fmt.Sprintf(`[<a href=\"https://maps.google.com/?q=%f,%f\">[G]</a>] [<a href=\"https://yandex.ru/maps/?source=serp_navig&text=%f,%f\">[Y]</a>],`, finalCoordinates["Latitude0"], finalCoordinates["Longitude0"], finalCoordinates["Latitude0"], finalCoordinates["Longitude0"]), -1)
	}

	// Convert Links
	reg = regexp.MustCompile(`<a href="(.+?)"*>(.+?)</a>`)
	strArr = reg.FindAllString(str, -1)
	for _, link := range strArr {
		newLink, err := convertUTF(link)
		if err != nil {
			log.Printf("Error URL conver =%s, Text=%s", err, str)
			continue
		}
		str = strings.Replace(str, link, newLink, -1)
	}
	return str
}
func deleteMapFloat(maps map[string]float64) {
	for key := range maps {
		delete(maps, key)
	}
}
func deleteMap(maps map[int]string) {
	for key := range maps {
		delete(maps, key)
	}
}
func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
func convertUTF(str string) (string, error) {
	u, err := url.QueryUnescape(str)
	if err != nil {
		log.Printf("\nURL %s, bad request %s", str, err)
		return str, err
	}
	return u, nil
}
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}
func convertTimeSec(times int) string {

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
		break
	case 2:
		str += fmt.Sprintf("%d дня ", timeDay)
		break
	case 3:
		str += fmt.Sprintf("%d дня ", timeDay)
		break
	case 4:
		str += fmt.Sprintf("%d дня ", timeDay)
		break
	default:
		str += fmt.Sprintf("%d дней ", timeDay)
		break
	}
	// Часы
	switch timeHour {
	case 0:
		str += ""
	case 1:
		str += fmt.Sprintf("%d час ", timeHour)
		break
	case 2:
		str += fmt.Sprintf("%d часа ", timeHour)
		break
	case 3:
		str += fmt.Sprintf("%d часа ", timeHour)
		break
	case 4:
		str += fmt.Sprintf("%d часа ", timeHour)
		break
	default:
		str += fmt.Sprintf("%d часов ", timeHour)
		break
	}

	// Минуты
	switch timeMin {
	case 0:
		str += ""
	case 1:
		str += fmt.Sprintf("%d минута ", timeMin)
		break
	case 2:
		str += fmt.Sprintf("%d минуты ", timeMin)
		break
	case 3:
		str += fmt.Sprintf("%d минуты ", timeMin)
		break
	case 4:
		str += fmt.Sprintf("%d минуты ", timeMin)
		break
	default:
		str += fmt.Sprintf("%d минут ", timeMin)
		break
	}
	// Секунды
	switch timeSec {
	case 0:
		str += ""
	case 1:
		str += fmt.Sprintf("%d секунда", timeSec)
		break
	case 2:
		str += fmt.Sprintf("%d секунды", timeSec)
		break
	case 3:
		str += fmt.Sprintf("%d секунды", timeSec)
		break
	case 4:
		str += fmt.Sprintf("%d секунды", timeSec)
		break
	default:
		str += fmt.Sprintf("%d секунд", timeSec)
		break
	}

	// если не смогли распарсить
	if str == "" {
		return fmt.Sprintf("%d секунд", times)
	}
	return strings.TrimSpace(str)
}

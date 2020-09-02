package main

import (
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"os"
	"reflect"
	"strings"
	"time"
)

var (
	bufModel      Model
	photoMap      map[int]string
	locationMap   map[string]float64
	cookieJar     *cookiejar.Jar
	client        http.Client
	isAnswerBlock bool
	isWork        bool
)

func MainHandler(resp http.ResponseWriter, _ *http.Request) {
	_, _ = resp.Write([]byte("Hi there! I'm telegram bot @gogland_bot. My owner @maksimall89"))
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano()) // real random
	var err error
	var confJSON ConfigGameJSON

	if os.Getenv("Gogland_logs") == "1" {
		logInit()
		defer log.Println(recover())
	}

	// web server for heroku
	http.HandleFunc("/", MainHandler)
	go func() {
		err = http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	}()
	// create chanel
	webToBot := make(chan MessengerStyle, 100000) // chanel for send command in game and back information
	botToWeb := make(chan MessengerStyle, 100)    // chanel for send command in game and back information

	commandArr := []string{
		"/help - информация по всем доступным командам;",

		"/start /stop /restart - запуск в первый раз /остановка/перезагрузка бота после остановки;",
		"/pause /resume  - приостановить/возобновить прием кодов;",

		"/faq - краткий FAQ по моей работе;",
		"/task - задание с уровня;",
		"/msg - все сообщения с уровня;",
		"/timer - время до автоперехода;",
		"/codes - оставшиеся коды;",
		"/codesall - оставшиеся + снятые коды;",
		"/hints - все подсказки;",
		"/penalty - все штрафные подсказки;",
		"/bonuses - информация о бонусах;",
		"/joke - шутки про команду;",
		"/b - отправить бонусный код при ограничении;",
		"/getPenalty - взять штрафную подсказку;",
		"/add - добавить игрока;",

		"\n<b>Помогалки:</b>",
		"/ana - анаграммы;",
		"/smask - поиск по маске (* - мн.символов, ? - один символ);",
		"/ass - ассоциации к слову;",
		"/n2w - перевод из числа в букву;",
		"/w2n - перевод из буквы в число;",
		"/mt - из чила в номер элемента таблицы Менделеева;",
		"/mz - азбука морзе(из .- в букву англ/рус);",
		"/b2d - из 2 в 10;",
		"/d2b - из 10 в 2;",
		"/ac <code>code</code> - коды регионов или если указан код, то будет только один регион;",
		"/bra - шрифт Браиля выводит на экран (/bra 010000);",
		"/qw - перевод из одной раскладки в другую;",
	}

	// read config file
	var configuration ConfigBot
	configuration.init("config.json")

	// configuration bot
	bot, err := tgbotapi.NewBotAPI(configuration.TelegramBotToken)
	if err != nil {
		log.Println(err)
	}
	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)
	defer log.Println("Bot off!")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 0 // TIME OUT!!

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Println(err)
		log.Println("Failed to get updates")
	}

	// set default config
	var msgChannel MessengerStyle
	var chatId int64
	var buffer struct {
		Coordinate Coordinate
		TextToChat string
	}
	var update tgbotapi.Update // chanel Update from telegram

	isAnswerBlock = false // pass for enter code
	isWork = false        // state work bot
	isBonus := new(bool)  // code bonus

	// main cycle
	for {
		// В канал будут приходить все новые сообщения from web
		select {
		case msgChannel = <-webToBot:
			switch msgChannel.Type {
			case "photo":
				_, _ = bot.Send(tgbotapi.NewPhotoShare(chatId, msgChannel.ChannelMessage))
			case "location":
				if (buffer.Coordinate.Latitude == msgChannel.Latitude) && (buffer.Coordinate.Longitude == msgChannel.Longitude) {
					break
				}
				_, _ = bot.Send(tgbotapi.NewVenue(chatId, fmt.Sprintf(`%g %g`, msgChannel.Latitude, msgChannel.Longitude), "", msgChannel.Latitude, msgChannel.Longitude))
				buffer.Coordinate.Latitude = msgChannel.Latitude
				buffer.Coordinate.Longitude = msgChannel.Longitude
			case "text":
				// проверка на повтор сообщения
				if msgChannel.ChannelMessage == buffer.TextToChat || len(msgChannel.ChannelMessage) == 0 {
					continue
				} else {
					buffer.TextToChat = msgChannel.ChannelMessage
				}
				_ = sendMessageTelegram(chatId, msgChannel.ChannelMessage, msgChannel.MsgId, bot)
			}
		default:
			break
		}
		// В канал updates будут приходить все новые сообщения from telegram
		select {
		case update = <-updates:
			if update.Message == nil || update.Message.Text == "" || reflect.TypeOf(update.Message.Text).Kind() != reflect.String {
				continue
			}
		default:
			continue
		}

		if !isWork {
			chatId = update.Message.Chat.ID
			switch strings.ToLower(update.Message.Command()) {
			case "b", "pause", "resume", "restart", "hints", "penalty", "codesall", "codes", "task", "msg", "timer", "bonuses":
				_ = sendMessageTelegram(chatId, "Игра ещё не началась.", 0, bot)
				continue
			default:
				break
			}
		} else {
			// Если пишут в другой чат ему, то игнор
			if chatId != update.Message.Chat.ID {
				continue
			}
		}

		switch strings.ToLower(update.Message.Command()) {
		case "help":
			go func() {
				var str string
				for _, item := range commandArr {
					str += item + "\n"
				}
				_ = sendMessageTelegram(chatId, str, 0, bot)
			}()
		case "faq":
			_ = sendMessageTelegram(chatId, `<code>Бот отправляет в движок все сообщения содержащие английские буквы(word), слово-буквы (слово1) не разделенные пробелом. Чтобы принудительно отправить код необходимо использовать восклицательный знак — "!". Если на уровне есть ограничение на ввод, то бот остановит прием кодов и продолжит их принимать на следующем уровне автоматически. Если на уровне есть координаты, то бот их преобразует в GPS-координаты и отправит как локацию в чат, также бот преобразует все сообщения в чате написанные в одну строку 52.4456 52.4563 в координаты, при этом координаты в тексте задания, который придёт вам в один клик копируются. Все скрытие ссылки под картинками будут также отмечены в чате, а сами картинки всегда скидываются отдельным сообщением для удобства. И самое главное помните, что бот не волшебник, он только учутся и за ним нужно следить.</code>`, 0, bot)
		case "stop":
			msgChannel.ChannelMessage = "stop"
			botToWeb <- msgChannel
			isWork = false
		case "postfix":
			confJSON.Postfix = update.Message.CommandArguments()
			_ = sendMessageTelegram(chatId, "Постфикс принят", 0, bot)
		case "prefix":
			confJSON.Prefix = update.Message.CommandArguments()
			_ = sendMessageTelegram(chatId, "Префикс принят", 0, bot)
		case "b":
			*isBonus = true
			go sendCode(&client, &confJSON, update.Message.CommandArguments(), isBonus, webToBot, update.Message.MessageID)
		case "pause":
			isAnswerBlock = true
			_ = sendMessageTelegram(chatId, "Приём кодов <b>приостановлен</b>.\nДля возобновления наберите /resume Для ввода бонусных кодов /b", 0, bot)
		case "resume":
			_ = sendMessageTelegram(chatId, "Приём кодов <b>возобновлён</b>.\nДля приостановки наберите /pause Для ввода бонусных кодов /b", 0, bot)
		case "getPenalty":
			if len(update.Message.CommandArguments()) > 0 {
				getPenalty(&client, &confJSON, update.Message.Text, webToBot)
			} else {
				_ = sendMessageTelegram(chatId, "Недостаточно символов. Необходимо отправить: <code>/getPenalty 1111</code>", 0, bot)
			}
		case "restart":
			log.Printf("RESTART JSON %s change  config.", update.Message.From.UserName)
			// create cookie
			cookieJar, _ = cookiejar.New(nil)
			client = http.Client{
				Jar: cookieJar,
			}
			// start go worker!
			go workerJSON(&client, &confJSON, botToWeb, webToBot, &isWork, &isAnswerBlock)
			defer close(botToWeb)
			defer close(webToBot)
			isWork = true
			chatId = update.Message.Chat.ID
		case "start":
			// set config can only owner
			if (update.Message.From.UserName == configuration.OwnName) && (update.Message.CommandArguments() != "") {
				_ = sendMessageTelegram(chatId, confJSON.init(update.Message.CommandArguments()), 0, bot)
				log.Printf("%s change config JSON.", update.Message.From.UserName)

				// create cookie
				cookieJar, _ = cookiejar.New(nil)
				client = http.Client{
					Jar: cookieJar,
				}
				// start go worker!
				go workerJSON(&client, &confJSON, botToWeb, webToBot, &isWork, &isAnswerBlock)
				defer close(botToWeb)
				defer close(webToBot)

				_ = sendMessageTelegram(chatId, "All change is apply JSON", 0, bot)
				isWork = true
				chatId = update.Message.Chat.ID
			} else {
				if isWork {
					_ = sendMessageTelegram(chatId, "I work!", 0, bot)
				} else {
					_ = sendMessageTelegram(chatId, "You is't my own!", 0, bot)
				}
				log.Printf("%s try to change config!", update.Message.From.UserName)
			}
		case "hints":
			go func() {
				if len(bufModel.Level.Helps) > 0 {
					_ = sendMessageTelegram(chatId, getFirstHelps(bufModel.Level.Helps, confJSON), 0, bot)
				} else {
					_ = sendMessageTelegram(chatId, "Подсказок нет&#128552;!\n", 0, bot)
				}
			}()
		case "penalty":
			go func() {
				if len(bufModel.Level.PenaltyHelps) > 0 {
					_ = sendMessageTelegram(chatId, getFirstHelps(bufModel.Level.PenaltyHelps, confJSON), 0, bot)
				} else {
					_ = sendMessageTelegram(chatId, "Штрафных подсказок нет&#128532;!\n", 0, bot)
				}
			}()
		case "codesall":
			//codesall - оставшиеся + снятые коды.
			go func() {
				_ = sendMessageTelegram(chatId, getLeftCodes(bufModel.Level.Sectors, false), 0, bot)
			}()
		case "codes":
			///codes - оставшиеся коды.
			go func() {
				_ = sendMessageTelegram(chatId, getLeftCodes(bufModel.Level.Sectors, true), 0, bot)
			}()
		case "task":
			go func() {
				_ = sendMessageTelegram(chatId, getFirstTask(bufModel.Level.Tasks, confJSON), 0, bot)
			}()
		case "msg":
			go func() {
				msgLevel := getFirstMessages(bufModel.Level.Messages, confJSON)
				if len(msgLevel) == 0 {
					_ = sendMessageTelegram(chatId, "&#128495;Сообщений на уровне нет.", 0, bot)
				} else {
					_ = sendMessageTelegram(chatId, msgLevel, 0, bot)
				}
			}()
		case "timer":
			go func() {
				_ = sendMessageTelegram(chatId, getFirstTimer(bufModel.Level), 0, bot)
			}()
		case "add":
			go func() {
				_ = sendMessageTelegram(chatId, addUser(&client, &confJSON, update.Message.CommandArguments()), update.Message.MessageID, bot)
			}()
		case "bonuses":
			go func() {
				if len(bufModel.Level.Bonuses) > 0 {
					_ = sendMessageTelegram(chatId, getFirstBonuses(bufModel.Level.Bonuses, confJSON), 0, bot)
				} else {
					_ = sendMessageTelegram(chatId, "Бонусов в игре нет!\n", 0, bot)
				}
			}()
		case "joke":
			if len(configuration.Jokes) > 0 {
				_ = sendMessageTelegram(chatId, configuration.Jokes[rand.Intn(len(configuration.Jokes))], 0, bot)
			} else {
				_ = sendMessageTelegram(chatId, "Шуток нет.", 0, bot)
			}
		case "n2w":
			go func() {
				if len(update.Message.CommandArguments()) > 0 {
					_ = sendMessageTelegram(chatId, transferToAlphabet(update.Message.CommandArguments(), true), update.Message.MessageID, bot)
				} else {
					_ = sendMessageTelegram(chatId, "Недостаточно символов. Необходимо отправить: <code>/n2w 22 5</code>", update.Message.MessageID, bot)
				}
			}()
		case "w2n":
			go func() {
				if len(update.Message.CommandArguments()) > 0 {
					_ = sendMessageTelegram(chatId, transferToAlphabet(update.Message.CommandArguments(), false), update.Message.MessageID, bot)
				} else {
					_ = sendMessageTelegram(chatId, "Недостаточно символов. Необходимо отправить: <code>/n2w А D</code>", update.Message.MessageID, bot)
				}
			}()
		case "ac":
			go func() {
				if len(update.Message.CommandArguments()) > 0 {
					_ = sendMessageTelegram(chatId, autoCode(update.Message.CommandArguments()), update.Message.MessageID, bot)
				} else {
					_ = sendMessageTelegram(chatId, "Недостаточно символов. Необходимо отправить: <code>/ac 18</code>", update.Message.MessageID, bot)
				}
			}()
		case "ana":
			go func() {
				if len(update.Message.CommandArguments()) > 0 {
					_ = sendMessageTelegram(chatId, searchAnagramAndMaskWord(update.Message.CommandArguments(), true), update.Message.MessageID, bot)
				} else {
					_ = sendMessageTelegram(chatId, "Недостаточно символов. Необходимо отправить: <code>/ana сел</code>", update.Message.MessageID, bot)
				}
			}()
		case "bra":
			go func() {
				if len(update.Message.CommandArguments()) > 5 {
					_ = sendMessageTelegram(chatId, braille(update.Message.CommandArguments()), update.Message.MessageID, bot)
				} else {
					_ = sendMessageTelegram(chatId, "Недостаточно символов. Необходимо отправить: <code>/bra 101000</code>", update.Message.MessageID, bot)
				}
			}()
		case "smask":
			go func() {
				if len(update.Message.CommandArguments()) > 1 {
					_ = sendMessageTelegram(chatId, searchAnagramAndMaskWord(update.Message.CommandArguments(), false), update.Message.MessageID, bot)
				} else {
					_ = sendMessageTelegram(chatId, "Недостаточно символов. Необходимо отправить: <code>/smask а?в*</code>", update.Message.MessageID, bot)
				}
			}()
		case "mt":
			go func() {
				if len(update.Message.CommandArguments()) > 0 {
					_ = sendMessageTelegram(chatId, tableMendeleev(update.Message.CommandArguments()), update.Message.MessageID, bot)
				} else {
					_ = sendMessageTelegram(chatId, "Недостаточно символов. Необходимо отправить: <code>/mt 11</code>", update.Message.MessageID, bot)
				}
			}()
		case "mz":
			go func() {
				if len(update.Message.CommandArguments()) > 0 {
					_ = sendMessageTelegram(chatId, morse(update.Message.CommandArguments()), update.Message.MessageID, bot)
				} else {
					_ = sendMessageTelegram(chatId, "Недостаточно символов. Необходимо отправить: <code>/mz .-</code>", update.Message.MessageID, bot)
				}
			}()
		case "ass":
			go func() {
				if len(update.Message.CommandArguments()) > 0 {
					_ = sendMessageTelegram(chatId, associations(update.Message.CommandArguments()), update.Message.MessageID, bot)
				} else {
					_ = sendMessageTelegram(chatId, "Недостаточно символов. Необходимо отправить: <code>/ass лето</code>", update.Message.MessageID, bot)
				}
			}()
		case "b2d":
			go func() {
				if len(update.Message.CommandArguments()) > 0 {
					_ = sendMessageTelegram(chatId, bin(update.Message.CommandArguments(), false), update.Message.MessageID, bot)
				} else {
					_ = sendMessageTelegram(chatId, "Недостаточно символов. Необходимо отправить: <code>/b2d 101</code>", update.Message.MessageID, bot)
				}
			}()
		case "d2b":
			go func() {
				if len(update.Message.CommandArguments()) > 0 {
					_ = sendMessageTelegram(chatId, bin(update.Message.CommandArguments(), true), update.Message.MessageID, bot)
				} else {
					_ = sendMessageTelegram(chatId, "Недостаточно символов. Необходимо отправить: <code>/d2b 99</code>", update.Message.MessageID, bot)
				}
			}()
		case "qw":
			go func() {
				if len(update.Message.CommandArguments()) > 0 {
					_ = sendMessageTelegram(chatId, translateQwerty(update.Message.CommandArguments()), update.Message.MessageID, bot)
				} else {
					_ = sendMessageTelegram(chatId, "Недостаточно символов. Необходимо отправить: <code>/qw ц z</code>", update.Message.MessageID, bot)
				}
			}()
		default:
			go func() {
				sendLocation(searchLocation(update.Message.Text), webToBot)
			}()

			if isWork {
				// WTF symbol what I need ignore
				if strings.ContainsAny(update.Message.Text, ":;/, '*+@#$%^&(){}[]|") {
					break
				}
				// check  codes
				if strings.ContainsAny(strings.ToLower(update.Message.Text), "abcdefghijklmnopqrstuvwxyz0123456789!?") {
					*isBonus = false
					if !isAnswerBlock && (!strings.HasPrefix(update.Message.Text, "!") || !strings.HasPrefix(update.Message.Text, "?")) {
						arrCodes := strings.Split(update.Message.Text, "\n")
						for _, code := range arrCodes {
							go sendCode(&client, &confJSON, code, isBonus, webToBot, update.Message.MessageID)
						}
					} else {
						_ = sendMessageTelegram(chatId, "Приём кодов <b>приостановлен</b>.\nДля возобновления наберите /resume", 0, bot)
					}
				}
			}
		}
	}
}

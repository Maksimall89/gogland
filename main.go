package main

import (
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"
)

var bufModel Model
var photoMap map[int]string
var locationMap map[string]float64

func MainHandler(resp http.ResponseWriter, _ *http.Request) {
	_, _ = resp.Write([]byte("Hi there! I'm telegram bot @gogland_bot. My owner @maksimall89"))
}

func main() {
	// web server for heroku
	http.HandleFunc("/", MainHandler)
	go http.ListenAndServe(":"+os.Getenv("PORT"), nil)

	var buffTextToChat string
	var str string
	var pointerStr int
	var err error
	var confJSON ConfigGameJSON

	if os.Getenv("Gogland_logs") == "1" {
		str = "log" // name folder for logs
		// check what folder log is exist
		_, err = os.Stat(str)
		if os.IsNotExist(err) {
			_ = os.MkdirAll(str, os.ModePerm)
		}

		str = fmt.Sprintf("%s/%d-%02d-%02d-%02d-%02d-%02d-logFile.log", str, time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second())

		// configurator for logger
		// open a file
		fileLog, err := os.OpenFile(str, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			fmt.Printf("Error opening file: %v", err)
		}
		defer fileLog.Close()
		defer log.Println(recover())

		// assign it to the standard logger
		log.SetOutput(fileLog)
		log.SetPrefix("Gogland ")
	}

	// create chanel
	webToBot := make(chan MessengerStyle, 100000) // chanel for send command in game and back information
	botToWeb := make(chan MessengerStyle, 100000) // chanel for send command in game and back information

	commandArr := []string{
		"/help - информация по всем доступным командам;",

		"/start /stop /restart - запуск в первый раз /остановка/перезагрузка бота после остановки;",
		"/getstatus - статус бота;",
		"/pause /resume  - приостановить/возобновить прием кодов;",

		"/faq - краткий FAQ по моей работе;",
		"/task - задание с уровня;",
		"/msg - все сообщения с уровня;",
		"/timer - время до автоперехода;",
		"/codes - оставшиеся коды;",
		"/codestaken - снятые коды;",
		"/codesall - оставшиеся + снятые коды;",
		"/hints - все подсказки;",
		"/penalty - все штрафные подсказки;",
		"/bonuses - информация о бонусах;",
		"/joke - шутки про команду;",
		"/b - отправить бонусный код при ограничении;",
		"/getPenalty - взять штрафную подсказку;",

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

	defer log.Println("Bot off!.")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 0 // TIME OUT!!

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Println(err)
		log.Println("Failed to get updates")
	}

	// set default config
	var msgBot MessengerStyle // style for messenge chanel
	var msgChanel MessengerStyle
	var chatId int64

	var update tgbotapi.Update  // chanel Update from telegram
	var newMsg tgbotapi.Message // style message from telegram
	var msg tgbotapi.MessageConfig

	isWork := false                        // state work bot
	rand.Seed(time.Now().UTC().UnixNano()) // real random

	// main cycle
	for {
		select {
		// В канал msgChanel будут приходить все новые сообщения from web
		case msgChanel = <-webToBot:

			if isWork {
				// шлем только в чат игры
				chatId = msgBot.ChatId
			} else {
				// шлём в любой чат
				chatId = update.Message.Chat.ID
			}

			switch msgChanel.Type {
			case "photo":
				_, _ = bot.Send(tgbotapi.NewPhotoShare(chatId, msgChanel.ChannelMessage))
			case "location":
				_, _ = bot.Send(tgbotapi.NewVenue(chatId, fmt.Sprintf(`%g %g`, msgChanel.Latitude, msgChanel.Longitude), "", msgChanel.Latitude, msgChanel.Longitude))
			case "text":
				msg = tgbotapi.NewMessage(chatId, msgChanel.ChannelMessage)
				if msgChanel.MsgId != 0 {
					msg.ReplyToMessageID = msgChanel.MsgId
				} else {
					// проверка на повтор сообщения
					if msgChanel.ChannelMessage == buffTextToChat {
						break
					} else {
						buffTextToChat = msgChanel.ChannelMessage
					}
				}

				if len(msgChanel.ChannelMessage) == 0 {
					continue
				}

				msg.ParseMode = "HTML"
				msg.ChatID = msgBot.ChatId

				// если тест уровня большой, то дробим его на множество мелких сообщений
				if len(msg.Text) > 4090 { // ограничение на длину одного сообщения 4096
					log.Printf("Big massage. Len = %d", len(msg.Text))
					buffTextToChat = msg.Text
					for {
						pointerStr = strings.LastIndex(buffTextToChat[0:4090], "\n")
						msg.Text = buffTextToChat[0:pointerStr]
						newMsg, err = bot.Send(msg)
						if err != nil {
							msg.ParseMode = "Markdown"
							newMsg, _ = bot.Send(msg)
							log.Println(err)
							log.Println(msg.Text)
							log.Println(msgChanel.ChannelMessage)
						}

						if strings.Contains(msg.Text, "<b>Задание") {
							_, err := bot.PinChatMessage(tgbotapi.PinChatMessageConfig{ChatID: msg.ChatID, MessageID: newMsg.MessageID})
							if err != nil {
								log.Println(err)
							}
						}

						buffTextToChat = buffTextToChat[pointerStr:]
						if len(buffTextToChat) < 4091 {
							msg.Text = buffTextToChat
							_, err := bot.Send(msg)
							if err != nil {
								msg.ParseMode = "Markdown"
								_, _ = bot.Send(msg)
								log.Println(err)
								log.Println(msg.Text)
								log.Println(msgChanel.ChannelMessage)
							}
							break
						}
					}
				} else {
					newMsg, err = bot.Send(msg)
					if err != nil {
						log.Println(msg.Text)
						msg.ParseMode = "Markdown"
						msg.ChatID = msgBot.ChatId
						newMsg, _ = bot.Send(msg)
						log.Println(err)
						log.Println(msgChanel.ChannelMessage)
					}

					if strings.Contains(msg.Text, "<b>Задание") {
						_, err := bot.PinChatMessage(tgbotapi.PinChatMessageConfig{ChatID: msg.ChatID, MessageID: newMsg.MessageID})
						if err != nil {
							log.Println(err)
						}
					}

				}
			}
		default:
			break
		}

		select {
		// В канал updates будут приходить все новые сообщения from telegram
		case update = <-updates:
			if update.Message == nil {
				continue
			}
			if reflect.TypeOf(update.Message.Text).Kind() != reflect.String {
				continue
			}
			if update.Message.Text == "" {
				continue
			}
		default:
			continue
		}

		str = ""         // clear old state
		msgBot.MsgId = 0 // save new id

		switch strings.ToLower(update.Message.Command()) {
		case "help":
			for _, item := range commandArr {
				str += item + "\n"
			}
		case "faq":
			str = `<code>Бот отправляет в движок все сообщения содержащие английские буквы(word), слово-буквы (слово1) не разделенные пробелом. Чтобы принудительно отправить код необходимо использовать восклицательный знак — "!". Если на уровне есть ограничение на ввод, то бот остановит прием кодов и продолжит их принимать на следующем уровне автоматически. Если на уровне есть координаты, то бот их преобразует в GPS-координаты и отправит как локацию в чат, также бот преобразует все сообщения в чате написанные в одну строку 52.4456 52.4563 в координаты, при этом координаты в тексте задания, который придёт вам в один клик копируются. Все скрытие ссылки под картинками будут также отмечены в чате, а сами картинки всегда скидываются отдельным сообщением для удобства. И самое главное помните, что бот не волшебник, он только учутся и за ним нужно следить.</code>`
		case "stop":
			msgBot.ChannelMessage = "stop"
			botToWeb <- msgBot
			isWork = false
			break
		case "postfix":
			confJSON.Postfix = update.Message.CommandArguments()
			str = "Постфикс принят"
		case "prefix":
			confJSON.Prefix = update.Message.CommandArguments()
			str = "Префикс принят"
		case "b":
			msgBot.ChannelMessage = "bonuscodes " + update.Message.CommandArguments()
			msgBot.MsgId = update.Message.MessageID
			botToWeb <- msgBot
			break
		case "getPenalty":
			text := "Недостаточно символов. Необходимо отправить: <code>/getPenalty 1111</code>"
			if len(update.Message.CommandArguments()) > 0 {
				// TODO доделать взятие штрафной подсказки по аналогии с отправкой кода
				text = "getPenaltyJSON()"
			}
			msgBot.ChannelMessage = text
			msgBot.MsgId = update.Message.MessageID
			botToWeb <- msgBot
			break
		case "restart":
			if !isWork {
				log.Printf("RESTART JSON %s change  config.", update.Message.From.UserName)
				go workerJSON(confJSON, botToWeb, webToBot, &isWork)
				defer close(botToWeb)
				defer close(webToBot)

				str = "RESTART JSON! All change is applay"
				msgBot.ChatId = update.Message.Chat.ID
				isWork = true
			} else {
				str = "I work!"
			}
		case "start":
			// set config can only owner
			if (update.Message.From.UserName == configuration.OwnName) && (update.Message.CommandArguments() != "") && !isWork {
				// split arg
				args := strings.Split(update.Message.CommandArguments(), " ")

				if len(args) < 3 {
					str = "Need more arguments! [start login password http://DEMO.en.cx/GameDetails.aspx?gid=1]"
					log.Printf("%s try to change config! Need more arguments!", update.Message.From.UserName)
					break
				}
				if len(args) > 3 {
					str = "Слишком много аргументов!"
					log.Printf("%s try to change config!", update.Message.From.UserName)
					break
				}
				// configuration game
				confJSON.NickName = args[0]
				confJSON.Password = args[1]
				confJSON.URLGame = args[2]
				confJSON.separateURL()
				log.Printf("%s change config JSON.", update.Message.From.UserName)

				// start go worker!
				go workerJSON(confJSON, botToWeb, webToBot, &isWork)
				defer close(botToWeb)
				defer close(webToBot)

				str = "All change is apply JSON"
				msgBot.ChatId = update.Message.Chat.ID
				isWork = true
			} else {
				if isWork {
					str = "I work!"
				} else {
					str = "You is't my own!"
				}
				log.Printf("%s try to change config!", update.Message.From.UserName)
			}
		case "hints":
			if isWork && update.Message.Chat.ID == msgBot.ChatId {
				go func() {
					if len(bufModel.Level.Helps) > 0 {
						str = getFirstHelps(bufModel.Level.Helps, confJSON)
					} else {
						str = "Подсказок нет&#128552;!\n"
					}
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, str)
					msg.ParseMode = "HTML"
					_, _ = bot.Send(msg)
				}()
			} else {
				str = "Игра ещё не началась."
			}
		case "penalty":
			if isWork && update.Message.Chat.ID == msgBot.ChatId {
				go func() {
					if len(bufModel.Level.PenaltyHelps) > 0 {
						str = getFirstHelps(bufModel.Level.PenaltyHelps, confJSON)
					} else {
						str = "Штрафных подсказок нет&#128532;!\n"
					}
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, str)
					msg.ParseMode = "HTML"
					_, _ = bot.Send(msg)
				}()
			} else {
				str = "Игра ещё не началась."
			}
		case "codesall":
			//codesall - оставшиеся + снятые коды.
			if isWork && update.Message.Chat.ID == msgBot.ChatId {
				go func() {
					str += getLeftCodes(bufModel.Level.Sectors, false)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, str)
					msg.ParseMode = "HTML"
					_, _ = bot.Send(msg)
				}()
			} else {
				str = "Игра ещё не началась."
			}
		case "codes":
			///codes - оставшиеся коды.
			if isWork && update.Message.Chat.ID == msgBot.ChatId {
				go func() {
					str += getLeftCodes(bufModel.Level.Sectors, true)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, str)
					msg.ParseMode = "HTML"
					_, _ = bot.Send(msg)
				}()
			} else {
				str = "Игра ещё не началась."
			}
		case "task":
			if isWork && update.Message.Chat.ID == msgBot.ChatId {
				go func() {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, getFirstTask(bufModel.Level.Tasks, confJSON))
					msg.ParseMode = "HTML"
					_, _ = bot.Send(msg)
				}()
			} else {
				str = "Игра ещё не началась."
			}
		case "msg":
			if isWork && update.Message.Chat.ID == msgBot.ChatId {
				go func() {
					str = getFirstMessages(bufModel.Level.Messages, confJSON)
					if str == "" {
						str = "&#128495;Сообщений на уровне нет."
					}
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, str)
					msg.ParseMode = "HTML"
					_, _ = bot.Send(msg)
				}()
			} else {
				str = "Игра ещё не началась."
			}
		case "timer":
			if isWork {
				go func() {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, getFirstTimer(bufModel.Level))
					msg.ParseMode = "HTML"
					msg.ReplyToMessageID = update.Message.MessageID
					_, _ = bot.Send(msg)
				}()
			} else {
				str = "Игра ещё не началась."
			}
		case "bonuses":
			if isWork && update.Message.Chat.ID == msgBot.ChatId {
				go func() {
					if len(bufModel.Level.Bonuses) > 0 {
						str = getFirstBonuses(bufModel.Level.Bonuses, confJSON)
					} else {
						str = "Бонусов в игре нет!\n"
					}
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, str)
					msg.ParseMode = "HTML"
					_, _ = bot.Send(msg)
				}()
			} else {
				str = "Игра ещё не началась."
			}
		case "joke":
			if len(configuration.Jokes) > 0 {
				str = configuration.Jokes[rand.Intn(len(configuration.Jokes))]
			} else {
				str = "Шуток нет."
			}
		case "n2w":
			go func() {
				text := "Недостаточно символов. Необходимо отправить: <code>/n2w 22 5</code>"
				if len(update.Message.CommandArguments()) > 0 {
					text = transferToAlphabet(update.Message.CommandArguments(), true)
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
				msg.ParseMode = "HTML"
				msg.ReplyToMessageID = update.Message.MessageID
				_, _ = bot.Send(msg)
			}()
		case "w2n":
			go func() {
				text := "Недостаточно символов. Необходимо отправить: <code>/n2w А D</code>"
				if len(update.Message.CommandArguments()) > 0 {
					text = transferToAlphabet(update.Message.CommandArguments(), false)
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
				msg.ParseMode = "HTML"
				msg.ReplyToMessageID = update.Message.MessageID
				_, _ = bot.Send(msg)
			}()
		case "ac":
			go func() {
				text := "Недостаточно символов. Необходимо отправить: <code>/ac 18</code>"
				if len(update.Message.CommandArguments()) > 0 {
					text = autoCode(update.Message.CommandArguments())
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
				msg.ParseMode = "HTML"
				msg.ReplyToMessageID = update.Message.MessageID
				_, _ = bot.Send(msg)
			}()
		case "ana":
			go func() {
				text := "Недостаточно символов. Необходимо отправить: <code>/ana сел</code>"
				if len(update.Message.CommandArguments()) > 0 {
					text, _ = anagram(update.Message.CommandArguments())
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
				msg.ParseMode = "HTML"
				msg.ReplyToMessageID = update.Message.MessageID
				_, _ = bot.Send(msg)
			}()
		case "bra":
			go func() {
				text := "Недостаточно символов. Необходимо отправить: <code>/bra 101000</code>"
				if len(update.Message.CommandArguments()) > 6 {
					text = braille(update.Message.CommandArguments())
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
				msg.ParseMode = "HTML"
				msg.ReplyToMessageID = update.Message.MessageID
				_, _ = bot.Send(msg)
			}()
		case "smask":
			go func() {
				text := "Недостаточно символов. Необходимо отправить: <code>/smask а?в*</code>"
				if len(update.Message.CommandArguments()) > 2 {
					text, _ = searchForMask(update.Message.CommandArguments())
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
				msg.ParseMode = "HTML"
				msg.ReplyToMessageID = update.Message.MessageID
				_, _ = bot.Send(msg)
			}()
		case "mt":
			go func() {
				text := "Недостаточно символов. Необходимо отправить: <code>/mt 11</code>"
				if len(update.Message.CommandArguments()) > 0 {
					text = tableMendeleev(update.Message.CommandArguments())
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
				msg.ParseMode = "HTML"
				msg.ReplyToMessageID = update.Message.MessageID
				_, _ = bot.Send(msg)
			}()
		case "mz":
			go func() {
				text := "Недостаточно символов. Необходимо отправить: <code>/mz .-</code>"
				if len(update.Message.CommandArguments()) > 0 {
					text = morse(update.Message.CommandArguments())
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
				msg.ParseMode = "HTML"
				msg.ReplyToMessageID = update.Message.MessageID
				_, _ = bot.Send(msg)
			}()
		case "ass":
			go func() {
				text := "Недостаточно символов. Необходимо отправить: <code>/ass лето</code>"
				if len(update.Message.CommandArguments()) > 0 {
					text, _ = associations(update.Message.CommandArguments())
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
				msg.ParseMode = "HTML"
				msg.ReplyToMessageID = update.Message.MessageID
				_, _ = bot.Send(msg)
			}()
		case "b2d":
			go func() {
				text := "Недостаточно символов. Необходимо отправить: <code>/b2d 101</code>"
				if len(update.Message.CommandArguments()) > 0 {
					text = bin(update.Message.CommandArguments(), false)
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
				msg.ParseMode = "HTML"
				msg.ReplyToMessageID = update.Message.MessageID
				_, _ = bot.Send(msg)
			}()
		case "d2b":
			go func() {
				text := "Недостаточно символов. Необходимо отправить: <code>/d2b 99</code>"
				if len(update.Message.CommandArguments()) > 0 {
					text = bin(update.Message.CommandArguments(), true)
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
				msg.ParseMode = "HTML"
				msg.ReplyToMessageID = update.Message.MessageID
				_, _ = bot.Send(msg)
			}()
		case "qw":
			go func() {
				text := "Недостаточно символов. Необходимо отправить: <code>/qw ц z</code>"
				if len(update.Message.CommandArguments()) > 0 {
					text = translateQwerty(update.Message.CommandArguments())
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
				msg.ParseMode = "HTML"
				msg.ReplyToMessageID = update.Message.MessageID
				_, _ = bot.Send(msg)
			}()

		default:
			go func() {
				sentLocation(searchLocation(update.Message.Text), webToBot)
			}()

			if isWork && (update.Message.Chat.ID == msgBot.ChatId) {
				// check  codes
				if strings.HasPrefix(update.Message.Text, "!") || strings.HasPrefix(update.Message.Text, "?") || strings.HasPrefix(update.Message.Text, "/") {
					msgBot.MsgId = update.Message.MessageID
					msgBot.ChannelMessage = update.Message.Text
					botToWeb <- msgBot
					break
				}

				// WTF symbol what i need ignore
				if strings.IndexAny(strings.ToLower(update.Message.Text), ":;/, '*+@#$%^&(){}[]|") != -1 {
					break
				}

				// check english lang or number
				if strings.IndexAny(strings.ToLower(update.Message.Text), "abcdefghijklmnopqrstuvwxyz0123456789") != -1 {
					msgBot.MsgId = update.Message.MessageID
					msgBot.ChannelMessage = update.Message.Text
					botToWeb <- msgBot
					break
				}
			} else {
				break
			}
		}

		// проверка длины сообщения
		if str == "" {
			continue
		}

		// структура для отправки сообщения в чат
		// если игра запущена, то отсылаем только в нужный нам чат
		if isWork {
			msg = tgbotapi.NewMessage(msgBot.ChatId, str)
		} else {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, str)
		}

		msg.ParseMode = "HTML"
		_, err = bot.Send(msg)
		if err != nil {
			log.Println(err)
			log.Println(msg.Text)
		}
	}
}

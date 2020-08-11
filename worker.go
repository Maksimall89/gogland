package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func workerJSON(client *http.Client, game *ConfigGameJSON, botToWeb chan MessengerStyle, webToBot chan MessengerStyle, isWork *bool, isAnswerBlock *bool) string {

	photoMap = make(map[int]string)
	locationMap = make(map[string]float64)

	isGameStart := false
	isPromblemStart := false

	var str string
	var bufStr string
	var modelGame Model
	var msgBot MessengerStyle

	msgBot.Type = "text"
	msgBot.MsgId = 0

	// Заходим в движок
	msgBot.ChannelMessage = enterGameJSON(client, *game)
	webToBot <- msgBot

	// Получаем актуальное состояние игры
	for {
		// если игра уже идёт
		if isGameStart {
			break
		}
		select {
		case msg := <-botToWeb:
			if msg.ChannelMessage == "stop" {
				msgBot.ChannelMessage = "Бот выключен. Мы даже не играли &#128546; \nДля перезапуска используйте /restart"
				webToBot <- msgBot
				*isWork = false
				log.Printf("Bot %s stop.\n", game.Gid)
				return "Bot stop"
			}
		default:
			modelGame = gameEngineModel(client, *game)
			// что-то отсылаем если состояние игры изменилось лишь иначе идём на следующий круг
			switch modelGame.Event {
			case 0:
				// если у нас не сбой и реально какое-то число уровней есть в игре
				if modelGame.Level.LevelId != 0 {
					str = "Игра уже идёт!"
					isGameStart = true
					break
				} else {
					if !isPromblemStart {
						str = "Не смогу получить состояние игры..."
						isPromblemStart = true
					}
					enterGameJSON(client, *game)
					break
				}
			case 1:
				msgBot.ChannelMessage = "&#9940;Игра не существует!"
				webToBot <- msgBot
				*isWork = false
				log.Printf("Bot %s stop  - case 0.\n", game.Gid)
				return "ERROR"
			case 5:
				str = "Игра ещё не началась!"
			case 6:
				msgBot.ChannelMessage = "Игра закончилась."
				webToBot <- msgBot
				*isWork = false
				log.Printf("Bot %s stop - case 6.\n", game.Gid)
				return "FINISHED"
			case 7:
				str = "&#9940;Не подана заявка игроком, который запустил бота: " + game.NickName
			case 8:
				str = "&#9940;Не подана заявка командой!"
			case 9:
				str = "&#9940;Команда игрока <b>" + game.NickName + "</b> еще не принята в игру:" + game.URLGame
			case 10:
				str = "&#9940;У игрока нет команды, который запустил бота: " + game.NickName
			case 11:
				str = "&#9940;Игрок не активен в команд, который запустил бота: " + game.NickName
			case 12:
				str = "&#9940;В игре нет уровней!"
			case 13:
				str = "&#9940;Превышено количество участников в команде!"
			case 16:
				str = "&#9940;Уровень снят"
			case 17:
				msgBot.ChannelMessage = "&#9940;Игра закончена!"
				webToBot <- msgBot
				*isWork = false
				log.Printf("Bot %s stop - case 17.\n", game.Gid)
				return "FINISHED"
			case 18:
				str = "&#9940;Уровень снят!"
			case 19:
				str = "&#9940;Уровень пройден автопереходом!"
			case 20:
				str = "&#9940;Все сектора отгаданы!"
			case 21:
				str = "&#9940;Уровень снят!"
			case 22:
				str = "&#9940;Таймаут уровня!"
			default:
				str = "&#9940;Проблемы с игрой...."
				enterGameJSON(client, *game)
				break
			}

			if str == bufStr {
				continue
			}

			msgBot.ChannelMessage = str
			webToBot <- msgBot

			bufStr = str
		}
	}

	// Цикл самой игры
	for {
		str = ""
		select {
		// В канал msg будут приходить все новые сообщения from telegram
		case msg := <-botToWeb:
			if msg.ChannelMessage == "stop" {
				msgBot.ChannelMessage = "<b>Бот выключен.</b> \nДля перезапуска используйте /restart"
				msgBot.MsgId = 0 // clear replay
				webToBot <- msgBot
				*isWork = false
				log.Printf("Bot %s stop.\n", game.Gid)
				return "Bot stop"
			}
		default:
			modelGame = gameEngineModel(client, *game)

			// Проверка на конец игры
			if modelGame.Event == 6 || modelGame.Event == 17 {
				msgBot.ChannelMessage = "&#128293;Игра завершена!\n<b>Вы молодцы, штаб ОГОНЬ!</b>"
				webToBot <- msgBot
				*isWork = false
				log.Printf("Game finished %s", game.URLGame)
				return "FINISH GAME"
			}

			// если мы не на уровне или что-то пошло не так
			if modelGame.Event != 0 && modelGame.GameId == 0 || modelGame.Level.Number == 0 {
				enterGameJSON(client, *game)
				continue
			}

			// Если это новый уровень, то
			if bufModel.Level.Number != modelGame.Level.Number {
				deleteMap(photoMap)
				deleteMapFloat(locationMap)

				msgBot.ChannelMessage = "&#9889;Выдан новый уровень!"
				webToBot <- msgBot

				// Название
				str = fmt.Sprintf("&#128681; #Уровень <b>%d</b> %s\n", modelGame.Level.Number, modelGame.Level.Name)
				// Автопереход
				str += getFirstTimer(modelGame.Level)
				// Сектора
				str += getFirstSector(modelGame.Level)
				// Ограничение на ввод
				if modelGame.Level.HasAnswerBlockRule {
					str += fmt.Sprintf("\n&#10071;<b>Ограничение на ввод!</b>\nПриём кодов <b>приостановлен</b>.\nУ вас %d попыток на ", modelGame.Level.AttemtsNumber)
					// блокировка установлена для: 0,1 – для игрока; 2 – для команды
					if modelGame.Level.BlockTargetId == 2 {
						str += "команду"
					} else {
						str += "игрока"
					}
					str += fmt.Sprintf("за %s\nДля возобновления наберите /resume\nЧтобы отправить бонусные коды введите: /b <code>код</code>\n\n", convertTimeSec(modelGame.Level.AttemtsPeriod))
					*isAnswerBlock = true
				}
				// Сообщения
				str += getFirstMessages(modelGame.Level.Messages, *game)
				//  Задание
				str += getFirstTask(modelGame.Level.Tasks, *game)
				// Подсказки
				str += getFirstHelps(modelGame.Level.Helps, *game)
				// Штрафные подсказки
				str += getFirstHelps(modelGame.Level.PenaltyHelps, *game)
				//  Бонусы
				str += getFirstBonuses(modelGame.Level.Bonuses, *game)
				//  sent message
				msgBot.ChannelMessage = str
				webToBot <- msgBot

				// Собираемы карты с кордами и фотками
				photoMap = searchPhoto(str)
				locationMap = searchLocation(str)
				//sent maps
				sentPhoto(photoMap, webToBot)
				sentLocation(locationMap, webToBot)

			} else {

				str = ""
				// Автопереход
				if bufModel.Level.TimeoutSecondsRemain != modelGame.Level.TimeoutSecondsRemain {
					switch modelGame.Level.TimeoutSecondsRemain {
					case 60:
						str += fmt.Sprintf("&#9200;До автоперехода 1&#8419;<b> минута!</b>\n")
					case 300:
						str += fmt.Sprintf("&#9200;До автоперехода 5&#8419;<b> минут!</b>\n")
					}
				}

				//  Сектора +
				if bufModel.Level.SectorsLeftToClose != modelGame.Level.SectorsLeftToClose {
					switch modelGame.Level.SectorsLeftToClose {
					case 3:
						str += fmt.Sprintf("&#128269;Осталось найти 3&#8419; сектора.\nСобрались тряпки!\n")
					case 1:
						str += fmt.Sprintf("&#128269;Осталось найти 1&#8419; сектор.\nТушёнки, пора уже найти его!\n")
					}
				}

				msgBot.ChannelMessage = str
				webToBot <- msgBot

				// Сообщения
				compareMessages(modelGame.Level.Messages, bufModel.Level.Messages, *game, webToBot)
				//  Задание
				compareTasks(modelGame.Level.Tasks, bufModel.Level.Tasks, *game, webToBot)
				//  Подсказки
				compareHelps(modelGame.Level.Helps, bufModel.Level.Helps, *game, webToBot)
				//  Штрафные подсказки
				compareHelps(modelGame.Level.PenaltyHelps, bufModel.Level.PenaltyHelps, *game, webToBot)
				//  Бонусы
				compareBonuses(modelGame.Level.Bonuses, bufModel.Level.Bonuses, *game, webToBot)
			}

			// копируем всё в буфер
			bufModel = modelGame
			game.LevelNumber = modelGame.Level.Number

			time.Sleep(time.Duration(70000 + rand.Intn(1000)))
		}
	}
}

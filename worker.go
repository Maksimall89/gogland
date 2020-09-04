package main

import (
	"fmt"
	"log"
	"net/http"
)

func workerJSON(client *http.Client, game *ConfigGameJSON, botToWeb chan MessengerStyle, webToBot chan MessengerStyle, isWork *bool, isAnswerBlock *bool) string {
	photoMap := make(map[int]string)
	locationMap := make(map[string]float64)

	var str string
	var modelGame Model
	var msgBot MessengerStyle

	msgBot.Type = "text"
	msgBot.MsgId = 0

	// Заходим в движок
	msgBot.ChannelMessage = enterGame(client, *game)
	webToBot <- msgBot

	// Получаем актуальное состояние игры
	isWorkMu.Lock()
	err := startGame(client, game, isWork, botToWeb, webToBot)
	if err != nil {
		*isWork = false
		isWorkMu.Unlock()
		return "Bot off"
	}
	isWorkMu.Unlock()

	// Цикл самой игры
	for {
		select {
		// В канал msg будут приходить все новые сообщения from telegram
		case msg := <-botToWeb:
			if msg.ChannelMessage == "stop" {
				msgBot.ChannelMessage = "<b>Бот выключен.</b> \nДля перезапуска используйте /restart"
				msgBot.MsgId = 0 // clear replay
				webToBot <- msgBot
				isWorkMu.Lock()
				*isWork = false
				isWorkMu.Unlock()
				log.Printf("Bot %s stop.\n", game.Gid)
				return "Bot stop"
			}
		default:
			modelGame = gameEngineModel(client, *game)

			// Проверка на конец игры
			if modelGame.Event == 6 || modelGame.Event == 17 {
				msgBot.ChannelMessage = "&#128293;Игра завершена!\n<b>Вы молодцы, штаб ОГОНЬ!</b>"
				webToBot <- msgBot
				isWorkMu.Lock()
				*isWork = false
				isWorkMu.Unlock()
				log.Printf("Game finished %s", game.URLGame)
				return "FINISH GAME"
			}

			// если мы не на уровне или что-то пошло не так
			if modelGame.Event != 0 && modelGame.GameId == 0 || modelGame.Level.Number == 0 {
				enterGame(client, *game)
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
					isAnswerBlockMu.Lock()
					*isAnswerBlock = true
					isAnswerBlockMu.Unlock()
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
				sendPhoto(photoMap, webToBot)
				sendLocation(locationMap, webToBot)
			} else {
				str = ""
				// Автопереход
				if bufModel.Level.TimeoutSecondsRemain != modelGame.Level.TimeoutSecondsRemain {
					switch modelGame.Level.TimeoutSecondsRemain {
					case 60:
						str += "&#9200;До автоперехода 1&#8419;<b> минута!</b>\n"
					case 300:
						str += "&#9200;До автоперехода 5&#8419;<b> минут!</b>\n"
					}
				}

				//  Сектора +
				if bufModel.Level.SectorsLeftToClose != modelGame.Level.SectorsLeftToClose {
					switch modelGame.Level.SectorsLeftToClose {
					case 3:
						str += "&#128269;Осталось найти 3&#8419; сектора.\nСобрались тряпки!\n"
					case 1:
						str += "&#128269;Осталось найти 1&#8419; сектор.\nТушёнки, пора уже найти его!\n"
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
		}
	}
}

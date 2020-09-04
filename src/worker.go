package src

import (
	"fmt"
	"log"
	"net/http"
)

func WorkerJSON(client *http.Client, game *ConfigGameJSON, botToWeb chan MessengerStyle, webToBot chan MessengerStyle, isWork *bool, isAnswerBlock *bool) string {
	photoMap := make(map[int]string)
	locationMap := make(map[string]float64)

	var str string
	var modelGame Model
	var msgBot MessengerStyle

	msgBot.Type = "text"
	msgBot.MsgId = 0

	// Заходим в движок
	msgBot.ChannelMessage = EnterGame(client, *game)
	webToBot <- msgBot

	// Получаем актуальное состояние игры
	IsWorkMu.Lock()
	err := StartGame(client, game, isWork, botToWeb, webToBot)
	if err != nil {
		*isWork = false
		IsWorkMu.Unlock()
		return "Bot off"
	}
	IsWorkMu.Unlock()

	// Цикл самой игры
	for {
		select {
		// В канал msg будут приходить все новые сообщения from telegram
		case msg := <-botToWeb:
			if msg.ChannelMessage == "stop" {
				msgBot.ChannelMessage = "<b>Бот выключен.</b> \nДля перезапуска используйте /restart"
				msgBot.MsgId = 0 // clear replay
				webToBot <- msgBot
				IsWorkMu.Lock()
				*isWork = false
				IsWorkMu.Unlock()
				log.Printf("Bot %s stop.\n", game.Gid)
				return "Bot stop"
			}
		default:
			modelGame = GameEngineModel(client, *game)

			// Проверка на конец игры
			if modelGame.Event == 6 || modelGame.Event == 17 {
				msgBot.ChannelMessage = "&#128293;Игра завершена!\n<b>Вы молодцы, штаб ОГОНЬ!</b>"
				webToBot <- msgBot
				IsWorkMu.Lock()
				*isWork = false
				IsWorkMu.Unlock()
				log.Printf("Game finished %s", game.URLGame)
				return "FINISH GAME"
			}

			// если мы не на уровне или что-то пошло не так
			if modelGame.Event != 0 && modelGame.GameId == 0 || modelGame.Level.Number == 0 {
				EnterGame(client, *game)
				continue
			}

			// Если это новый уровень, то
			if BufModel.Level.Number != modelGame.Level.Number {
				DeleteMap(photoMap)
				DeleteMapFloat(locationMap)
				msgBot.ChannelMessage = "&#9889;Выдан новый уровень!"
				webToBot <- msgBot

				// Название
				str = fmt.Sprintf("&#128681; #Уровень <b>%d</b> %s\n", modelGame.Level.Number, modelGame.Level.Name)
				// Автопереход
				str += GetFirstTimer(modelGame.Level)
				// Сектора
				str += GetFirstSector(modelGame.Level)
				// Ограничение на ввод
				if modelGame.Level.HasAnswerBlockRule {
					str += fmt.Sprintf("\n&#10071;<b>Ограничение на ввод!</b>\nПриём кодов <b>приостановлен</b>.\nУ вас %d попыток на ", modelGame.Level.AttemtsNumber)
					// блокировка установлена для: 0,1 – для игрока; 2 – для команды
					if modelGame.Level.BlockTargetId == 2 {
						str += "команду"
					} else {
						str += "игрока"
					}
					str += fmt.Sprintf("за %s\nДля возобновления наберите /resume\nЧтобы отправить бонусные коды введите: /b <code>код</code>\n\n", ConvertTimeSec(modelGame.Level.AttemtsPeriod))
					IsAnswerBlockMu.Lock()
					*isAnswerBlock = true
					IsAnswerBlockMu.Unlock()
				}
				// Сообщения
				str += GetFirstMessages(modelGame.Level.Messages, *game)
				//  Задание
				str += GetFirstTask(modelGame.Level.Tasks, *game)
				// Подсказки
				str += GetFirstHelps(modelGame.Level.Helps, *game)
				// Штрафные подсказки
				str += GetFirstHelps(modelGame.Level.PenaltyHelps, *game)
				//  Бонусы
				str += GetFirstBonuses(modelGame.Level.Bonuses, *game)
				//  sent message
				msgBot.ChannelMessage = str
				webToBot <- msgBot
				// Собираемы карты с кордами и фотками
				photoMap = SearchPhoto(str)
				locationMap = SearchLocation(str)
				//sent maps
				SendPhoto(photoMap, webToBot)
				SendLocation(locationMap, webToBot)
			} else {
				str = ""
				// Автопереход
				if BufModel.Level.TimeoutSecondsRemain != modelGame.Level.TimeoutSecondsRemain {
					switch modelGame.Level.TimeoutSecondsRemain {
					case 60:
						str += "&#9200;До автоперехода 1&#8419;<b> минута!</b>\n"
					case 300:
						str += "&#9200;До автоперехода 5&#8419;<b> минут!</b>\n"
					}
				}

				//  Сектора +
				if BufModel.Level.SectorsLeftToClose != modelGame.Level.SectorsLeftToClose {
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
				CompareMessages(modelGame.Level.Messages, BufModel.Level.Messages, *game, webToBot)
				//  Задание
				CompareTasks(modelGame.Level.Tasks, BufModel.Level.Tasks, *game, webToBot)
				//  Подсказки
				CompareHelps(modelGame.Level.Helps, BufModel.Level.Helps, *game, webToBot)
				//  Штрафные подсказки
				CompareHelps(modelGame.Level.PenaltyHelps, BufModel.Level.PenaltyHelps, *game, webToBot)
				//  Бонусы
				CompareBonuses(modelGame.Level.Bonuses, BufModel.Level.Bonuses, *game, webToBot)
			}

			// копируем всё в буфер
			BufModel = modelGame
			game.LevelNumber = modelGame.Level.Number
		}
	}
}

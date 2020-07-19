package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func enterGameJSON(client *http.Client, game ConfigGameJSON) string {

	type JSONenter struct {
		Error                int         `json:"Error"`
		Message              string      `json:"Message"`
		IPUnblockURL         interface{} `json:"IpUnblockUrl"`
		BruteForceUnblockURL interface{} `json:"BruteForceUnblockUrl"`
		ConfirmEmailURL      interface{} `json:"ConfirmEmailUrl"`
		CaptchaURL           interface{} `json:"CaptchaUrl"`
		AdminWhoCanActivate  interface{} `json:"AdminWhoCanActivate"`
	}

	msgBot := MessengerStyle{}
	msgBot.Type = "text"

	var counter int8
	bodyJSON := &JSONenter{}

	for counter = 0; counter < 3; counter++ {

		resp, err := client.PostForm(fmt.Sprintf("http://%s/login/signin?json=1", game.SubUrl), url.Values{"Login": {game.NickName}, "Password": {game.Password}, "ddlNetwork": {"1"}})
		if err != nil || resp == nil {
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			client.PostForm(fmt.Sprintf("http://%s/Login.aspx?return=/GameDetails.aspx?gid=%s", game.SubUrl, game.Gid), url.Values{"socialAssign": {"0"}, "Login": {game.NickName}, "Password": {game.Password}, "EnButton1": {"Вход"}, "ddlNetwork": {"1"}})
			continue
		}

		err = json.Unmarshal(body, bodyJSON)
		if err != nil {
			client.PostForm(fmt.Sprintf("http://%s/Login.aspx?return=/GameDetails.aspx?gid=%s", game.SubUrl, game.Gid), url.Values{"socialAssign": {"0"}, "Login": {game.NickName}, "Password": {game.Password}, "EnButton1": {"Вход"}, "ddlNetwork": {"1"}})
			continue
		}

		if bodyJSON.Error == 0 {
			return fmt.Sprintf("&#10004;<b>Авторизация прошла успешно</b> на игру: %s", game.URLGame)
		} else {
			return fmt.Sprintf("&#9940;Авторизация прошла НЕ успешно.\n%s\n%s", bodyJSON.Message, game.SubUrl)
		}
	}
	return fmt.Sprintf("&#9940;Превышено число попыток авторизации на игру %s!", game.URLGame)
}
func gameEngineModel(client *http.Client, game ConfigGameJSON) Model {
	msgBot := MessengerStyle{}
	msgBot.Type = "text"
	var counter int8

	bodyJSON := &Model{}

	enterGameJSON(client, game)

	// 3 Попытки
	for counter = 0; counter < 3; counter++ {
		resp, err := client.Get(fmt.Sprintf("http://%s/GameEngines/Encounter/Play/%s?json=1", game.SubUrl, game.Gid))
		if err != nil || resp == nil {
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			enterGameJSON(client, game)
			continue
		}

		if strings.Contains(string(body), `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">`) {
			enterGameJSON(client, game)
			continue
		}

		err = json.Unmarshal(body, bodyJSON)
		if err != nil {
			log.Println(err)
			log.Println(string(body))
			enterGameJSON(client, game)
			continue
		}
		return *bodyJSON
	}

	return Model{}
}
func sentCodeJSON(client *http.Client, game ConfigGameJSON, code string, isBonus bool, webToBot chan MessengerStyle, MsgId int) {
	var msgBot MessengerStyle
	msgBot.MsgId = MsgId
	msgBot.Type = "text"

	// убираем знаки
	if code[0:1] == "!" || (code[0:1] == "?") || (code[0:1] == "/") {
		if len(code) > 1 {
			code = code[1:]
		} else {
			msgBot.ChannelMessage = "&#9940;Код слишком маленький!"
			webToBot <- msgBot
			return
		}
	}
	// Убираем пробелы
	code = strings.TrimSpace(code)

	// добавление префиксов и простфиксов
	if game.Prefix != "" && game.Postfix != "" {
		sectorCode := strings.Split(code, ".")
		if len(sectorCode) > 1 {
			code = game.Prefix + sectorCode[0] + game.Postfix + sectorCode[1]
		} else {
			sectorCode = strings.Split(code, "-")
			if len(sectorCode) > 1 {
				code = game.Prefix + sectorCode[0] + game.Postfix + sectorCode[1]
			}
		}
	}
	// Получаем текущие состояние игры
	ModelState := gameEngineModel(client, game)

	// Проверяем на дубль
	for _, LevelInfo := range ModelState.Level.MixedActions {
		if LevelInfo.Answer == code {
			msgBot.ChannelMessage = "Код введён &#9745;<b>ПОВТОРНО</b>!"
			webToBot <- msgBot
			return
		}
	}
	// если получили ошибку... то снова пытаемся получить статус игры
	if ModelState.Level.Number == 0 {
		for {
			ModelState = gameEngineModel(client, game)
			if ModelState.Level.Number != 0 {
				break
			}
		}
	}
	if game.LevelNumber == ModelState.Level.Number && ModelState.Level.IsPassed == false && ModelState.Level.Dismissed == false {
		var resp *http.Response
		var err error
		var errCounter int8
		for errCounter = 0; errCounter < 5; errCounter++ {
			if isBonus {
				resp, err = client.PostForm(fmt.Sprintf("http://%s/GameEngines/Encounter/Play/%s?json=1/", game.SubUrl, game.Gid), url.Values{"LevelId": {fmt.Sprintf("%d", ModelState.Level.LevelId)}, "LevelNumber": {fmt.Sprintf("%d", ModelState.Level.Number)}, "BonusAction.Answer": {code}})
			} else {
				if ModelState.Level.HasAnswerBlockRule == false || ModelState.Level.BlockDuration <= 0 {
					resp, err = client.PostForm(fmt.Sprintf("http://%s/GameEngines/Encounter/Play/%s?json=1/", game.SubUrl, game.Gid), url.Values{"LevelId": {fmt.Sprintf("%d", ModelState.Level.LevelId)}, "LevelNumber": {fmt.Sprintf("%d", ModelState.Level.Number)}, "LevelAction.Answer": {code}})
				} else {
					msgBot.ChannelMessage = fmt.Sprintf("&#128219;<b>Ограничение на ввод.</b>\nЯ не смог отправить код&#128546;\nВы сможете ввести код через %d", ModelState.Level.BlockDuration)
					webToBot <- msgBot
					return
				}
			}

			if err != nil || resp == nil {
				log.Println("Ошибка при отправке кода 1.")
				log.Println(err)
				enterGameJSON(client, game)
				continue
			}

			defer resp.Body.Close()

			// читаем всё из body
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println("Ошибка при отправке кода 2.")
				log.Println(string(body))
				log.Println(err)
				enterGameJSON(client, game)
				continue
			}

			if strings.Contains(string(body), `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">`) {
				enterGameJSON(client, game)
				continue
			}

			bodyJSON := &Model{}
			err = json.Unmarshal(body, bodyJSON)
			if err != nil {
				log.Println("Ошибка генерации JSON.")
				log.Println(err)
				enterGameJSON(client, game)
				continue
			}

			// Если получили ошибку
			if ModelState.Event != 0 {
				log.Println("Слетела авторизация...")
				enterGameJSON(client, game)
				continue
			}

			if isBonus {
				if bodyJSON.EngineAction.BonusAction.IsCorrectAnswer == true {
					msgBot.ChannelMessage = "Бонусный код &#9989;<b>ВЕРНЫЙ</b>"
					webToBot <- msgBot
					return
				} else {
					msgBot.ChannelMessage = "Бонусный код &#10060;<b>НЕВЕРНЫЙ</b>"
					webToBot <- msgBot
					return
				}
			}

			if bodyJSON.EngineAction.LevelAction.IsCorrectAnswer == true {
				msgBot.ChannelMessage = "Код &#9989;<b>ВЕРНЫЙ</b>"
				webToBot <- msgBot
				return
			} else {
				msgBot.ChannelMessage = "Код &#10060;<b>НЕВЕРНЫЙ</b>"
				webToBot <- msgBot
				return
			}
		}
	} else {
		msgBot.ChannelMessage = fmt.Sprintf("Код от предыдущего уровня не отправлен. Вы уже на следующем уровне №%d", ModelState.Level.Number)
		webToBot <- msgBot
		return
	}

	msgBot.ChannelMessage = "&#9940;Превышено число попыток отправить код!\nПовторите ещё раз."
	webToBot <- msgBot
	return
}

func getPenaltyJSON(client *http.Client, game ConfigGameJSON, penaltyID int, webToBot chan MessengerStyle) {
	var msgBot MessengerStyle
	msgBot.Type = "text"
	var errCounter int8
	var str string

	for errCounter = 0; errCounter < 5; errCounter++ {
		resp, err := client.Get(fmt.Sprintf("http://%s/GameEngines/Encounter/Play/%s?json=1&pid=%d&pact=1", game.SubUrl, game.Gid, penaltyID))
		if err != nil || resp == nil {
			log.Println("Ошибка при взятии штрафной подсказки 1.")
			log.Println(err)
			enterGameJSON(client, game)
			continue
		}
		defer resp.Body.Close()

		// читаем всё из body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Ошибка при взятии штрафной подсказки 2.")
			log.Println(string(body))
			log.Println(err)
			enterGameJSON(client, game)
			continue
		}

		if strings.Contains(string(body), `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">`) {
			enterGameJSON(client, game)
			continue
		}

		bodyJSON := &HelpsStruct{}
		err = json.Unmarshal(body, bodyJSON)
		if err != nil {
			log.Println("Ошибка генерации JSON.")
			log.Println(err)
			enterGameJSON(client, game)
			continue
		}

		if bodyJSON.Number == 0 {
			str = "&#9940;Мы не смогли взять штрафную подсказку"
		} else {
			str = fmt.Sprintf("&#9889;Штрафная подсказка №%d:\n%s", bodyJSON.Number, bodyJSON.HelpText)
		}
		msgBot.ChannelMessage = str
		webToBot <- msgBot
		return
	}

	msgBot.ChannelMessage = "&#9940;Превышено число попыток взять штрафную подсказку!\nПовторите ещё раз."
	webToBot <- msgBot
	return
}
func getFirstBonuses(bonuses []BonusesStruct, gameConfig ConfigGameJSON) (str string) {
	for _, bonus := range bonuses {
		// Если ещё недоступен
		if bonus.SecondsToStart > 0 {
			str += fmt.Sprintf("&#128488;<b>Бонус №%d</b> %s будет доступен через %s.\n", bonus.Number, bonus.Name, convertTimeSec(bonus.SecondsToStart))
		}
		// Если доступен и отгадан
		if bonus.SecondsToStart == 0 && bonus.IsAnswered == true {
			str += fmt.Sprintf("&#10004;<b>Бонус №%d</b> %s (<b>выполнен</b>, награда: %s)\n", bonus.Number, bonus.Name, convertTimeSec(bonus.AwardTime))
		}
		// Если доступен и не отгадан
		if bonus.SecondsToStart == 0 && bonus.IsAnswered == false {
			str += fmt.Sprintf("&#128488;<b>Бонус №%d</b> %s\n%s\n", bonus.Number, bonus.Name, replaceTag(bonus.Task, gameConfig.SubUrl))
		}
		// Если есть подсказка/награда
		if bonus.Help != "" {
			str += fmt.Sprintf("%s\n", replaceTag(bonus.Help, gameConfig.SubUrl))
		}
	}
	return str
}
func getFirstSector(level LevelStruct) (str string) {
	if level.SectorsLeftToClose > 0 {
		str = fmt.Sprintf("&#128269;Вам нужно найти <b>%d из %d</b> секторов.\n", level.SectorsLeftToClose, len(level.Sectors))
	} else {
		str = "&#128269;На уровне 1 код.\n"
	}
	return str
}
func getFirstTimer(level LevelStruct) (str string) {
	if level.Timeout == 0 {
		str = "&#9203;<b>Автоперехода нет</b>\n"
	} else {
		str = fmt.Sprintf("&#9200;До автоперехода %s\n", convertTimeSec(level.TimeoutSecondsRemain))
		if level.TimeoutAward != 0 {
			str += fmt.Sprintf("&#128276;Штраф за автопереход %s\n", convertTimeSec(level.TimeoutAward*(-1)))
		}
	}
	return str
}
func getFirstTask(tasks []TaskStruct, gameConfig ConfigGameJSON) (str string) {
	if len(tasks) > 0 {
		str = fmt.Sprintf("\n&#9889;<b>Задание</b>:\n")
		for _, text := range tasks {
			if text.TaskText == "" {
				str += "\n&#10060;Текст задания <b>отсуствует</b>!\n"
			} else {
				str += fmt.Sprintf("%s\n", replaceTag(text.TaskTextFormatted, gameConfig.SubUrl))
			}
		}
	} else {
		str = "&#10060;<b>Задания отсуствуют!</b>\n"
	}
	return str
}
func getLeftCodes(sectors []SectorsStruct, isNeed bool) (str string) {

	if isNeed {
		str = "Вам осталось снять сектора:\n\n"
	}

	if len(sectors) > 0 {
		for _, sector := range sectors {
			if isNeed {
				if !sector.IsAnswered {
					str += fmt.Sprintf("&#10060;Сектор <b>%s №%d</b> не отгадан.\n", sector.Name, sector.Order)
				}
			} else {
				if sector.IsAnswered {
					str += fmt.Sprintf("&#10004;Сектор <b>%s №%d</b> отгадан, ответ: %s\n", sector.Name, sector.Order, sector.Answer.Answer)
				} else {
					str += fmt.Sprintf("&#10060;Сектор <b>%s №%d</b> не отгадан.\n", sector.Name, sector.Order)
				}
			}
		}
	} else {
		str = "&#128269;На уровне 1 код.\n"
	}
	return str
}
func getFirstHelps(helps []HelpsStruct, gameConfig ConfigGameJSON) (str string) {
	if len(helps) <= 0 {
		return ""
	}
	if helps[0].IsPenalty {
		if len(helps) > 0 {
			for _, penaltyHelp := range helps {
				// проверяем через сколько будет доступна подсказка
				if penaltyHelp.RemainSeconds > 0 {
					str += fmt.Sprintf("&#10004;<b>Штрафная подсказка</b> №%d будет через %s.\n", penaltyHelp.Number, convertTimeSec(penaltyHelp.RemainSeconds))
				}
				// Если подсказка уже доступна и взята
				if penaltyHelp.HelpText != "" {
					str += fmt.Sprintf("&#10004;<b>Штрафная подсказка</b> №%d:\n%s\n\n", penaltyHelp.Number, replaceTag(penaltyHelp.HelpText, gameConfig.SubUrl))
					continue
				}
				// Если подсказка уже доступна, но не взята
				if penaltyHelp.HelpText == "" && penaltyHelp.RemainSeconds == 0 {
					str += fmt.Sprintf("&#10004;<b>Штрафная подсказка</b> №%d доступна.\n", penaltyHelp.Number)
				}
				// Проверяем, что нужно подтеверждение и мы не взяли ещё подсказку
				if penaltyHelp.RequestConfirm && penaltyHelp.HelpText == "" {
					str += fmt.Sprintf("&#9888;Треубется подтверждение взятия штрафной подсказки: %d\n", penaltyHelp.HelpId)
				}
				// Штраф за взятие если ещё ёё не взяли
				if penaltyHelp.Penalty != 0 && penaltyHelp.HelpText == "" {
					str += fmt.Sprintf("&#9888;Штраф за взятие: %s\n", convertTimeSec(penaltyHelp.Penalty))
				}
				// Описание подсказки если ещё её не взяли
				if penaltyHelp.PenaltyComment != "" && penaltyHelp.HelpText == "" {
					str += fmt.Sprintf("<b>Описание:</b> %s\n", replaceTag(penaltyHelp.PenaltyComment, gameConfig.SubUrl))
				}
				str += "\n"
			}
		} else {
			str = ""
		}
	} else {
		if len(helps) > 0 {
			for _, help := range helps {
				str += fmt.Sprintf("&#10004;<b>Подсказка</b> №%d ", help.Number)
				if help.RemainSeconds > 0 {
					str += fmt.Sprintf("будет через %s.\n", convertTimeSec(help.RemainSeconds))
				} else {
					str += fmt.Sprintf("\n%s\n", replaceTag(help.HelpText, gameConfig.SubUrl))
				}
				str += "\n"
			}
		} else {
			str = "&#10060;Подсказок нет!\n"
		}
	}

	return str
}
func getFirstMessages(msg []MessagesStruct, gameConfig ConfigGameJSON) (str string) {
	if len(msg) > 0 {
		str += "&#128495;<b>Сообщения на уровне:</b>\n"
		for _, message := range msg {
			str += fmt.Sprintf("&#128172;%s\n", replaceTag(message.MessageText, gameConfig.SubUrl))
		}
	} else {
		str = ""
	}
	return str
}
func compareHelps(newHelps []HelpsStruct, oldHelps []HelpsStruct, gameConf ConfigGameJSON, webToBot chan MessengerStyle) {

	// если у нас всё нулевой длины, то нафиг нам идти дальше...
	if (len(newHelps) == 0) && (len(oldHelps) == 0) {
		return
	}

	var StartString string
	var str string
	msgBot := MessengerStyle{}
	msgBot.Type = "text"
	msgBot.MsgId = 0

	if newHelps[0].IsPenalty {
		StartString = "Штрафная подсказка"
	} else {
		StartString = "Подсказка"
	}

	if len(newHelps) > len(oldHelps) {
		msgBot.ChannelMessage = fmt.Sprintf("&#9889;<b>%s</b> появилась в движке.\n", StartString)
	}
	if len(newHelps) < len(oldHelps) {
		msgBot.ChannelMessage = fmt.Sprintf("&#9889;<b>%s</b> исчезла в движке&#128465;.\n", StartString)
	}
	webToBot <- msgBot
	msgBot.ChannelMessage = ""

	for numberNew, helpNew := range newHelps {
		// проверяем на длину, чтобы не выйти за пределы если добавили новую подсказку
		if numberNew < len(oldHelps) {
			// старая подсказка
			if newHelps[0].IsPenalty {
				if oldHelps[numberNew].PenaltyComment != helpNew.PenaltyComment {
					msgBot.ChannelMessage = fmt.Sprintf("&#11088;<b>Изменение в описании</b> штрафной подсказки №%d:\n%s\n", helpNew.Number, replaceTag(helpNew.PenaltyComment, gameConf.SubUrl))
					webToBot <- msgBot

					//if text have location
					sentLocation(searchLocation(helpNew.PenaltyComment), webToBot)
					//if text have img
					sentPhoto(searchPhoto(helpNew.PenaltyComment), webToBot)
				}
			}
			// проверка на изменение текста
			if oldHelps[numberNew].HelpText != helpNew.HelpText {
				msgBot.ChannelMessage = fmt.Sprintf("&#11088;<b>Изменение</b> в %s №%d:\n%s\n", StartString, helpNew.Number, replaceTag(helpNew.HelpText, gameConf.SubUrl))
				webToBot <- msgBot

				//if text have location
				sentLocation(searchLocation(helpNew.HelpText), webToBot)
				//if text have img
				sentPhoto(searchPhoto(helpNew.HelpText), webToBot)
			}
			// проверка, что уже отправляли
			if helpNew.RemainSeconds != oldHelps[numberNew].RemainSeconds {
				switch helpNew.RemainSeconds {
				case 60:
					msgBot.ChannelMessage = fmt.Sprintf("&#10004;<b>%s</b> №%d через 1&#8419; минуту.\n", StartString, helpNew.Number)
				case 300:
					msgBot.ChannelMessage = fmt.Sprintf("&#10004;<b>%s</b> №%d через 5&#8419; минут.\n", StartString, helpNew.Number)
				}
				webToBot <- msgBot
				msgBot.ChannelMessage = ""
			}
		} else {
			// новая подсказка!
			if newHelps[0].IsPenalty {
				str = fmt.Sprintf("&#11088;<b>Новая штрафная подсказка</b> №%d:\n%s\n", helpNew.Number, replaceTag(helpNew.PenaltyComment, gameConf.SubUrl))
				if helpNew.RemainSeconds > 0 {
					str += fmt.Sprintf("<b>Будет открыта через</b> %s.\n", convertTimeSec(helpNew.RemainSeconds))
				}
				if helpNew.HelpText != "" {
					str += fmt.Sprintf("<b>Штрафная подсказка:</b>\n%s\n", replaceTag(helpNew.HelpText, gameConf.SubUrl))
				}
				msgBot.ChannelMessage = str
				webToBot <- msgBot

				// ОПИСАНИЕ
				//if text have location
				sentLocation(searchLocation(helpNew.PenaltyComment), webToBot)
				//if text have img
				sentPhoto(searchPhoto(helpNew.PenaltyComment), webToBot)
				//if text have location
				sentLocation(searchLocation(helpNew.HelpText), webToBot)
				//if text have img
				sentPhoto(searchPhoto(helpNew.HelpText), webToBot)

			} else {
				str = fmt.Sprintf("&#11088;<b>Новая подсказка</b> №%d", helpNew.Number)
				if helpNew.RemainSeconds > 0 {
					str += fmt.Sprintf("будет открыта через %s.\n", convertTimeSec(helpNew.RemainSeconds))
				} else {
					str += fmt.Sprintf("\n%s\n", replaceTag(helpNew.HelpText, gameConf.SubUrl))
				}
				msgBot.ChannelMessage = str
				webToBot <- msgBot

				//if text have location
				sentLocation(searchLocation(helpNew.HelpText), webToBot)
				//if text have img
				sentPhoto(searchPhoto(helpNew.HelpText), webToBot)
			}

			switch helpNew.RemainSeconds {
			case 60:
				msgBot.ChannelMessage = fmt.Sprintf("&#10004;<b>%s</b> №%d через 1&#8419; минуту.\n", StartString, helpNew.Number)
			case 300:
				msgBot.ChannelMessage = fmt.Sprintf("&#10004;<b>%s</b> №%d через 5&#8419; минут.\n", StartString, helpNew.Number)
			}
			webToBot <- msgBot
		}
	}
}
func compareBonuses(new []BonusesStruct, old []BonusesStruct, gameConf ConfigGameJSON, webToBot chan MessengerStyle) {
	// если у нас всё нулевой длины, то нафиг нам идти дальше...
	if (len(new) == 0) && (len(old) == 0) {
		return
	}

	msgBot := MessengerStyle{}
	msgBot.Type = "text"
	msgBot.MsgId = 0

	var str string

	if len(new) > len(old) {
		msgBot.ChannelMessage = fmt.Sprintf("&#9889;<b>Бонус появился</b> в движке.\n")
	}
	if len(new) < len(old) {
		msgBot.ChannelMessage = fmt.Sprintf("&#9889;<b>Бонус исчез</b> в движке&#128465;.\n")
	}
	webToBot <- msgBot
	msgBot.ChannelMessage = ""

	for numberNew, bonusNew := range new {
		// проверяем на длину, чтобы не выйти за пределы если добавили новую подсказку
		if numberNew < len(old) {
			// проверка на описание
			if old[numberNew].Task != bonusNew.Task {
				str = fmt.Sprintf("&#11088;<b>Изменение в описании</b> бонуса %s №%d:\n%s\n", bonusNew.Name, bonusNew.Number, replaceTag(bonusNew.Task, gameConf.SubUrl))
				if bonusNew.SecondsLeft != 0 {
					str += fmt.Sprintf("&#9200;Время на выполнение бонуса ограниченно: %s\n", convertTimeSec(bonusNew.SecondsLeft))
				}
				msgBot.ChannelMessage = str
				webToBot <- msgBot
				msgBot.ChannelMessage = ""
				//if text have location
				sentLocation(searchLocation(bonusNew.Task), webToBot)
				//if text have img
				sentPhoto(searchPhoto(bonusNew.Task), webToBot)
			}
			// проверка на изменение текста
			if old[numberNew].Help != bonusNew.Help {
				if bonusNew.IsAnswered {
					msgBot.ChannelMessage = fmt.Sprintf("&#11088;<b>Изменение</b> в бонусе %s №%d (награда %s):\n%s\n", bonusNew.Name, bonusNew.Number, convertTimeSec(bonusNew.AwardTime), replaceTag(bonusNew.Help, gameConf.SubUrl))
				} else {
					msgBot.ChannelMessage = fmt.Sprintf("&#11088;<b>Изменение</b> в бонусе %s №%d:\n%s\n", bonusNew.Name, bonusNew.Number, replaceTag(bonusNew.Help, gameConf.SubUrl))
				}

				webToBot <- msgBot
				msgBot.ChannelMessage = ""
				//if text have location
				sentLocation(searchLocation(bonusNew.Help), webToBot)
				//if text have img
				sentPhoto(searchPhoto(bonusNew.Help), webToBot)
			}

			// проверка, что уже отправляли
			if bonusNew.SecondsToStart != old[numberNew].SecondsToStart {
				// будет доступен через
				str := ""
				switch bonusNew.SecondsToStart {
				case 60:
					str = fmt.Sprintf("&#10004;<b>Бонус</b> %s №%d доступен через 1&#8419; минуту.\n", bonusNew.Name, bonusNew.Number)
				case 300:
					str = fmt.Sprintf("&#10004;<b>Бонус</b> %s №%d доступен через 5&#8419; минут.\n", bonusNew.Name, bonusNew.Number)
				}
				// истечет через
				switch bonusNew.SecondsLeft {
				case 60:
					str += fmt.Sprintf("&#10004;<b>Бонус</b> %s №%d исчезнет через 1&#8419; минуту.\n", bonusNew.Name, bonusNew.Number)
				case 300:
					str += fmt.Sprintf("&#10004;<b>Бонус</b> %s №%d исчезнет через 5&#8419; минут.\n", bonusNew.Name, bonusNew.Number)
				}
				msgBot.ChannelMessage = str
				webToBot <- msgBot
				msgBot.ChannelMessage = ""
			}
		} else {
			str = fmt.Sprintf("&#11088;<b>Новый бонус</b> %s №%d\n%s\n", bonusNew.Name, bonusNew.Number, replaceTag(bonusNew.Task, gameConf.SubUrl))
			if bonusNew.Help != "" {
				str += fmt.Sprintf("Награда %s\n<b>Бонусная подсказка:</b>\n%s\n", convertTimeSec(bonusNew.AwardTime), replaceTag(bonusNew.Help, gameConf.SubUrl))
			}
			msgBot.ChannelMessage = str
			webToBot <- msgBot
			msgBot.ChannelMessage = ""
			// задание
			//if text have location
			sentLocation(searchLocation(bonusNew.Task), webToBot)
			//if text have img
			sentPhoto(searchPhoto(bonusNew.Task), webToBot)
			// ответ
			//if text have location
			sentLocation(searchLocation(bonusNew.Help), webToBot)
			//if text have img
			sentPhoto(searchPhoto(bonusNew.Help), webToBot)

			// будет доступен через
			str := ""
			switch bonusNew.SecondsToStart {
			case 60:
				str = fmt.Sprintf("&#10004;<b>Бонус</b> %s №%d доступен через 1&#8419; минуту.\n", bonusNew.Name, bonusNew.Number)
			case 300:
				str = fmt.Sprintf("&#10004;<b>Бонус</b> %s №%d доступен через 5&#8419; минут.\n", bonusNew.Name, bonusNew.Number)
			}
			// истечет через
			switch bonusNew.SecondsLeft {
			case 60:
				str += fmt.Sprintf("&#10004;<b>Бонус</b> %s №%d исчезнет через 1&#8419; минуту.\n", bonusNew.Name, bonusNew.Number)
			case 300:
				str += fmt.Sprintf("&#10004;<b>Бонус</b> %s №%d исчезнет через 5&#8419; минут.\n", bonusNew.Name, bonusNew.Number)
			}
			msgBot.ChannelMessage = str
			webToBot <- msgBot
			msgBot.ChannelMessage = ""
		}
	}
}
func compareMessages(newMessages []MessagesStruct, oldMessages []MessagesStruct, gameConf ConfigGameJSON, webToBot chan MessengerStyle) {
	// если у нас всё нулевой длины, то нафиг нам идти дальше...
	if (len(newMessages) == 0) && (len(oldMessages) == 0) {
		return
	}

	str := ""

	if len(newMessages) > len(oldMessages) {
		for number, model := range newMessages {
			if number < len(oldMessages) {
				if model.MessageText != oldMessages[number].MessageText {
					str += fmt.Sprintf("&#128495;<b>Сообщение изменено:</b>\n%s\n", model.MessageText)
				}
			} else {
				str += "&#128495;<b>Появилось сообщение</b>:\n"
				str += fmt.Sprintf("&#128172;%s\n", replaceTag(model.MessageText, gameConf.SubUrl))
			}
		}
	}

	if len(newMessages) < len(oldMessages) {
		str += fmt.Sprintf("&#128495;<b>Сообщение удалено&#128465;.</b>\n")
	}

	if len(newMessages) == len(oldMessages) {
		for number, model := range newMessages {
			if model.MessageText != oldMessages[number].MessageText {
				str += fmt.Sprintf("&#128495;<b>Сообщение изменено:</b>\n%s\n", model.MessageText)
			}
		}
	}

	msgBot := MessengerStyle{}
	msgBot.ChannelMessage = str
	msgBot.Type = "text"
	msgBot.MsgId = 0
	webToBot <- msgBot
}
func compareTasks(newTasks []TaskStruct, oldTasks []TaskStruct, gameConf ConfigGameJSON, webToBot chan MessengerStyle) {
	// если у нас всё нулевой длины, то нафиг нам идти дальше...
	if (len(newTasks) == 0) && (len(oldTasks) == 0) {
		return
	}

	msgBot := MessengerStyle{}
	msgBot.Type = "text"
	str := ""

	if len(newTasks) > len(oldTasks) {
		msgBot.ChannelMessage = "&#10004;<b>Появилось новое задание!</b>\n"
		webToBot <- msgBot
	}

	if len(newTasks) < len(oldTasks) {
		msgBot.ChannelMessage = "&#10060;<b>Задание удалено&#128465;!</b>\n"
		webToBot <- msgBot
	}

	for numberNew, newTask := range newTasks {
		str = ""
		if numberNew < len(oldTasks) {
			if newTask.TaskText != oldTasks[numberNew].TaskText {
				if newTask.TaskText != "" {
					str += fmt.Sprintf("&#10060;<b>Задание изменено</b>:\n%s", replaceTag(newTask.TaskText, gameConf.SubUrl))
				} else {
					str += fmt.Sprintf("&#10060;<b>Задание удалено&#128465;!</b>")
				}
			} else {
				continue
			}
		} else {
			str += fmt.Sprintf("&#10004;<b>Новое задание:</b>\n%s", replaceTag(newTask.TaskText, gameConf.SubUrl))
		}
		msgBot.ChannelMessage = str
		webToBot <- msgBot
	}
}

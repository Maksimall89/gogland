package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"testing"
	"time"
)

const pathTest = "config.json"

func TestEngineEnterGame(t *testing.T) {
	cookieJar, _ := cookiejar.New(nil)
	clientTEST := &http.Client{
		Jar: cookieJar,
	}

	var confGameENJSON ConfigGameJSON
	confGameENJSON.initTest("")
	exitGameTest(clientTEST, confGameENJSON.SubUrl)

	result := enterGame(clientTEST, confGameENJSON)
	if result != fmt.Sprintf("&#10004;<b>Авторизация прошла успешно</b> на игру: %s", confGameENJSON.URLGame) {
		t.Errorf("Unable to log in for %v", confGameENJSON)
	}
}
func TestEngineGameEngineModel(t *testing.T) {
	cookieJar, _ := cookiejar.New(nil)
	clientTEST := &http.Client{
		Jar: cookieJar,
	}

	confGameENJSON := ConfigGameJSON{}
	confGameENJSON.initTest(pathTest)

	exitGameTest(clientTEST, confGameENJSON.SubUrl)
	enterGame(clientTEST, confGameENJSON)
	result := gameEngineModel(clientTEST, confGameENJSON)
	if result.Event != 0 {
		t.Error("Impossible to get game state. Error code: ", result.Event)
	}
}
func TestEngineSendCode(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano()) // real random

	webToBotTEST := make(chan MessengerStyle, 10)
	var msgChanel MessengerStyle

	cookieJar, _ := cookiejar.New(nil)
	clientTEST := &http.Client{
		Jar: cookieJar,
	}

	var confGameENJSON ConfigGameJSON
	confGameENJSON.initTest(pathTest)

	exitGameTest(clientTEST, confGameENJSON.SubUrl)
	enterGame(clientTEST, confGameENJSON)
	confGameENJSON.LevelNumber = gameEngineModel(clientTEST, confGameENJSON).Level.Number

	code := fmt.Sprintf("НЕВЕРНЫЙ%d", rand.Int())
	isBonus := new(bool)
	*isBonus = false

	sendCode(clientTEST, &confGameENJSON, code, isBonus, webToBotTEST, 0)
	select {
	// В канал msgChanel будут приходить все новые сообщения from web
	case msgChanel = <-webToBotTEST:
		// конец
		if msgChanel.ChannelMessage == "" {
			return
		}
		if msgChanel.ChannelMessage != fmt.Sprintf("Код %s &#10060;<b>НЕВЕРНЫЙ</b>", code) {
			t.Errorf("For %s\nexpected %s\ngot %s", code, fmt.Sprintf("\n Мы ждали: Код %s &#10060;<b>НЕВЕРНЫЙ</b>", code), msgChanel.ChannelMessage)
		}
	default:
	}

	*isBonus = true
	sendCode(clientTEST, &confGameENJSON, code+"_bonus", isBonus, webToBotTEST, 0)
	select {
	// В канал msgChanel будут приходить все новые сообщения from web
	case msgChanel = <-webToBotTEST:
		// конец
		if msgChanel.ChannelMessage == "" {
			return
		}
		if msgChanel.ChannelMessage != "Бонусный код &#10060;<b>НЕВЕРНЫЙ</b>" {
			t.Errorf("For %s\nexpected %s\ngot %s", code, "\n Бонусный код &#10060;<b>НЕВЕРНЫЙ</b>", msgChanel.ChannelMessage)
		}
	default:
	}
}
func TestEngineGetPenalty(t *testing.T) {
	t.Parallel()
	t.Skip("Test need real penaltyID")

	penaltyID := "1111"

	rand.Seed(time.Now().UTC().UnixNano()) // real random
	webToBotTEST := make(chan MessengerStyle, 10)
	var msgChanel MessengerStyle

	cookieJar, _ := cookiejar.New(nil)
	clientTEST := &http.Client{
		Jar: cookieJar,
	}

	var confGameENJSON ConfigGameJSON
	confGameENJSON.initTest(pathTest)
	exitGameTest(clientTEST, confGameENJSON.SubUrl)
	enterGame(clientTEST, confGameENJSON)
	getPenalty(clientTEST, &confGameENJSON, penaltyID, webToBotTEST)

	select {
	// В канал msgChanel будут приходить все новые сообщения from web
	case msgChanel = <-webToBotTEST:
		// конец
		if msgChanel.ChannelMessage == "" {
			return
		}
		if !strings.Contains(msgChanel.ChannelMessage, "&#9889;Штрафная подсказка") {
			t.Errorf("For %s\nexpected %s\ngot %s", penaltyID, "\n Мы ждали: &#9889;Штрафная подсказка", msgChanel.ChannelMessage)
		}
	default:
	}
}
func TestEngineGetFirstBonuses(t *testing.T) {
	t.Parallel()

	var confGameENJSON ConfigGameJSON
	confGameENJSON.initTest(pathTest)

	bonuses := []BonusesStruct{
		{148016, "123e", 1, "sdfsdf", "werwerwersdfsdf", true, false, 0, 0, 0},
		{148017, "newss", 2, "sdfsdf", "werwerdfgwersdfsdf", true, false, 300, 60, 0},
		{148018, "newss1", 3, "sdfsd342f", "werwerdfg234wersdfsdf", true, false, 0, 0, 0},
		{148019, "neddd2", 4, "sdf42f", "wersdf", true, false, 0, 0, 0},
	}

	checkBonuses := "&#10004;<b>Бонус №1</b> 123e (<b>выполнен</b>, награда: 0 секунд)\nwerwerwersdfsdf\n" +
		"&#128488;<b>Бонус №2</b> newss будет доступен через 5 минут.\nwerwerdfgwersdfsdf\n" +
		"&#10004;<b>Бонус №3</b> newss1 (<b>выполнен</b>, награда: 0 секунд)\nwerwerdfg234wersdfsdf\n" +
		"&#10004;<b>Бонус №4</b> neddd2 (<b>выполнен</b>, награда: 0 секунд)\nwersdf\n"

	result := getFirstBonuses(bonuses, confGameENJSON)
	if getFirstBonuses(bonuses, confGameENJSON) != checkBonuses {
		t.Errorf("For %v\nexpected %s\ngot %s", bonuses, checkBonuses, result)
	}
}
func TestEngineFirstSectors(t *testing.T) {
	t.Parallel()
	type testPair struct {
		input  LevelStruct
		output string
	}

	var tests = []testPair{
		{LevelStruct{270768, "rr1rr", 9, 0, 0, 0, false, false, StartTimeStruct{5}, false, 1, 1, 1, 1, 3, 1, 1, nil, nil, nil, []SectorsStruct{{1, 2, "bb", AnswerStruct{}, false}, {2, 3, "bbb", AnswerStruct{}, false}, {3, 3, "bbb", AnswerStruct{}, false}}, nil, nil, nil}, "&#128269;Вам нужно найти <b>1 из 3</b> секторов.\n"},
		{LevelStruct{270768, "rr2rr", 9, 100, 10, 60, false, false, StartTimeStruct{5}, false, 0, 1, 1, 1, 0, 0, 0, nil, nil, nil, nil, nil, nil, nil}, "&#128269;На уровне 1 код.\n"},
	}

	for _, pair := range tests {
		result := getFirstSector(pair.input)
		if result != pair.output {
			t.Errorf("For %v\nexpected %s\ngot %s", pair.input, pair.output, result)
		}
	}
}
func TestEngineFirstTimer(t *testing.T) {
	t.Parallel()

	type testPair struct {
		input  LevelStruct
		output string
	}

	var tests = []testPair{
		{LevelStruct{270768, "rr1rr", 9, 0, 0, 0, false, false, StartTimeStruct{5}, false, 0, 1, 0, 1, 7, 2, 5, nil, nil, nil, nil, nil, nil, nil}, "&#9203;<b>Автоперехода нет</b>\n"},
		{LevelStruct{270768, "rr2rr", 9, 100, 10, 60, false, false, StartTimeStruct{5}, false, 0, 1, 0, 1, 7, 2, 5, nil, nil, nil, nil, nil, nil, nil}, "&#9200;До автоперехода 1 минута\n&#128276;Штраф за автопереход -10 секунд\n"},
	}

	for _, pair := range tests {
		result := getFirstTimer(pair.input)
		if result != pair.output {
			t.Errorf("For %v\nexpected %s\ngot %s", pair.input, pair.output, result)
		}
	}
}
func TestEngineFirstTask(t *testing.T) {
	t.Parallel()

	var confGameENJSON ConfigGameJSON
	confGameENJSON.initTest(pathTest)
	type testPair struct {
		input  []TaskStruct
		output string
	}

	var tests = []testPair{
		{[]TaskStruct{}, "&#10060;<b>Задания отсуствуют!</b>\n"},
		{[]TaskStruct{
			{true, "sdfsdfsdfsdfsdf", "rr1rr"},
			{true, "", ""}}, "\n&#9889;<b>Задание</b>:\nrr1rr\n\n&#10060;Текст задания <b>отсуствует</b>!\n"},
	}

	for _, pair := range tests {
		result := getFirstTask(pair.input, confGameENJSON)
		if result != pair.output {
			t.Errorf("For %v\nexpected %s\ngot %s", pair.input, pair.output, result)
		}
	}
}
func TestEngineGetLeftCodes(t *testing.T) {
	t.Parallel()
	type testPair struct {
		input  bool
		output string
	}

	var tests = []testPair{
		{true, "Вам осталось снять сектора:\n\n&#10060;Сектор <b>bb №1</b> не отгадан.\n&#10060;Сектор <b>bbb №3</b> не отгадан.\n"},
		{false, "&#10060;Сектор <b>bb №1</b> не отгадан.\n&#10004;Сектор <b>bbb №2</b> отгадан, ответ: \n&#10060;Сектор <b>bbb №3</b> не отгадан.\n"},
	}
	sectors := []SectorsStruct{{1, 1, "bb", AnswerStruct{}, false}, {2, 2, "bbb", AnswerStruct{}, true}, {3, 3, "bbb", AnswerStruct{}, false}}

	for _, pair := range tests {
		result := getLeftCodes(sectors, pair.input)
		if result != pair.output {
			t.Errorf("For %v\nexpected %s\ngot %s", pair.input, pair.output, result)
		}
	}
}
func TestEngineFirstHelps(t *testing.T) {
	t.Parallel()

	var confGameENJSON ConfigGameJSON
	confGameENJSON.initTest(pathTest)

	type testPair struct {
		input  []HelpsStruct
		output string
	}

	var tests = []testPair{
		{[]HelpsStruct{}, ""},
		{[]HelpsStruct{
			{243732, 1, "sdfsdfdsfsdf", false, 0, "", false, 0, 10},
			{243732, 2, "sdfsdf234dsfsdf", false, 0, "", false, 0, 300},
			{243732, 3, "sdfs1111sfsdf", false, 10, "", false, 0, 0},
		}, "&#10004;<b>Подсказка</b> №1 будет через 10 секунд.\n\n&#10004;<b>Подсказка</b> №2 будет через 5 минут.\n\n&#10004;<b>Подсказка</b> №3 \nsdfs1111sfsdf\n\n"},
		{[]HelpsStruct{{243732, 4, "sdfs3233241sfsdf", true, 10, "ddwww", true, 0, 0}}, "&#10004;<b>Штрафная подсказка</b> №4:\nsdfs3233241sfsdf\n\n"},
	}

	for _, pair := range tests {
		result := getFirstHelps(pair.input, confGameENJSON)
		if result != pair.output {
			t.Errorf("For %v\nexpected %s\ngot %s", pair.input, pair.output, result)
		}
	}
}
func TestEngineFirstMessages(t *testing.T) {
	t.Parallel()

	var confGameENJSON ConfigGameJSON
	confGameENJSON.initTest(pathTest)

	type testPair struct {
		input  []MessagesStruct
		output string
	}

	var tests = []testPair{
		{[]MessagesStruct{}, ""},
		{[]MessagesStruct{{158796, "Maksisdfsdfmal_adst", 2419, "   Magn  Телефоны организаторов:895143Максим\r\n", "otum:  Телефоны организаторов:8915 Мам\r<br/> ", true}}, "&#128495;<b>Сообщения на уровне:</b>\n&#128172;   Magn  Телефоны организаторов:895143Максим\r\n"},
	}

	for _, pair := range tests {
		result := getFirstMessages(pair.input, confGameENJSON)
		if result != pair.output {
			t.Errorf("For %v\nexpected %s\ngot %s", pair.input, pair.output, result)
		}
	}
}
func TestEngineCompareHelps(t *testing.T) {
	t.Parallel()

	confGameENJSON := ConfigGameJSON{}
	confGameENJSON.initTest(pathTest)
	webToBotTEST := make(chan MessengerStyle, 10)
	var msgChanel MessengerStyle

	oldHelps := []HelpsStruct{
		{243732, 1, "sdfsdfdsfsdf", false, 0, "", false, 0, 10},
		{243732, 2, "sdfsdfdsfsdf", false, 0, "", false, 0, 0},
	}

	newHelps := []HelpsStruct{
		{243732, 1, "sdfsdfdsfsdf", false, 0, "", false, 0, 10},
		{243732, 2, "sdfsdf234dsfsdf", false, 0, "", false, 0, 300},
		{243732, 3, "sdfs1111sfsdf", false, 10, "", false, 0, 0},
	}

	arrHelps := []string{
		"&#9889;<b>Подсказка</b> появилась в движке.\n",
		"&#11088;<b>Изменение</b> в Подсказка №2:\nsdfsdf234dsfsdf\n",
		"&#10004;<b>Подсказка</b> №2 через 5&#8419; минут.\n",
		"&#11088;<b>Новая подсказка</b> №3\nsdfs1111sfsdf\n",
		"&#11088;<b>Новая подсказка</b> №3\nsdfs1111sfsdf\n",
	}

	compareHelps(newHelps, oldHelps, confGameENJSON, webToBotTEST)
	for _, helps := range arrHelps {
		select {
		// В канал msgChanel будут приходить все новые сообщения from web
		case msgChanel = <-webToBotTEST:
			// конец
			if msgChanel.ChannelMessage == "" {
				return
			}
			// проверка сообщений
			if msgChanel.ChannelMessage != helps {
				t.Errorf("For new: %v\n old: %v\n\nexpected %s\ngot %s", newHelps, oldHelps, helps, msgChanel.ChannelMessage)
			}
		default:
		}
	}
}
func TestEngineCompareHelpsPenalty(t *testing.T) {
	t.Parallel()

	var confGameENJSON ConfigGameJSON
	confGameENJSON.initTest(pathTest)
	webToBotTEST := make(chan MessengerStyle, 10)
	var msgChanel MessengerStyle

	oldHelps := []HelpsStruct{
		{243732, 1, "sdfsdfdsfsdf", true, 0, "Комментарий", false, 0, 10},
		{243733, 2, "sdfsdfdsfsdf", true, 0, "", false, 0, 0},
		{243734, 3, "sdfsdfdsfsdf", true, 0, "", false, 0, 0},
	}

	newHelps := []HelpsStruct{
		{243732, 1, "sdfsdfdsfsdf", true, 0, "Комментарий", false, 0, 10},
		{243733, 2, "sdfsdf234dsfsdf", true, 0, "k1", false, 0, 300},
		{243734, 3, "sdfsdf234dsfsdf", true, 0, "", false, 0, 60},
		{243735, 4, "sdfs1111sfsdf", true, 0, "к2", false, 0, 0},
	}

	arrHelps := []string{
		"&#9889;<b>Штрафная подсказка</b> появилась в движке.\n",
		"&#11088;<b>Изменение в описании</b> штрафной подсказки №2:\nk1\n",
		"&#11088;<b>Изменение</b> в Штрафная подсказка №2:\nsdfsdf234dsfsdf\n",
		"&#10004;<b>Штрафная подсказка</b> №2 через 5&#8419; минут.\n",
		"&#11088;<b>Изменение</b> в Штрафная подсказка №3:\nsdfsdf234dsfsdf\n",
		"&#10004;<b>Штрафная подсказка</b> №3 через 1&#8419; минуту.\n",
		"&#11088;<b>Новая штрафная подсказка</b> №4:\nк2\n<b>Штрафная подсказка:</b>\nsdfs1111sfsdf\n",
		"&#11088;<b>Новая штрафная подсказка</b> №4:\nк2\n<b>Штрафная подсказка:</b>\nsdfs1111sfsdf\n",
	}

	compareHelps(newHelps, oldHelps, confGameENJSON, webToBotTEST)
	for _, help := range arrHelps {
		select {
		// В канал msgChanel будут приходить все новые сообщения from web
		case msgChanel = <-webToBotTEST:
			// конец
			if msgChanel.ChannelMessage == "" {
				return
			}
			// проверка сообщений
			if msgChanel.ChannelMessage != help {
				t.Errorf("For new: %v\nold: %v\nexpected %s\ngot %s", newHelps, oldHelps, help, msgChanel.ChannelMessage)
			}
		default:
		}
	}
}
func TestEngineCompareBonuses(t *testing.T) {
	t.Parallel()

	confGameENJSON := ConfigGameJSON{}
	confGameENJSON.initTest(pathTest)
	webToBotTEST := make(chan MessengerStyle, 10)
	var msgChanel MessengerStyle

	oldBonus := []BonusesStruct{
		{148016, "123e", 1, "sdfsdf", "werwerwersdfsdf", true, false, 0, 0, 0},
		{148017, "newss", 2, "sdfsdf", "werwerdfgwersdfsdf", true, false, 0, 0, 0},
		{148018, "newss1", 3, "sdfsdf", "werwerdfgwersdfsdf", true, false, 0, 0, 0},
	}

	newBonus := []BonusesStruct{
		{148016, "123e", 1, "sdfsdf", "werwerwersdfsdf", true, false, 0, 0, 0},
		{148017, "newss", 2, "sdfsdf", "werwerdfgwersdfsdf", true, false, 300, 60, 0},
		{148018, "newss1", 3, "sdfsd342f", "werwerdfg234wersdfsdf", true, false, 0, 0, 0},
		{148019, "neddd2", 4, "sdf42f", "wersdf", true, false, 0, 0, 0},
	}

	arrBonuses := []string{
		"&#9889;<b>Бонус появился</b> в движке.\n",
		"&#10004;<b>Бонус</b> newss №2 доступен через 5&#8419; минут.\n&#10004;<b>Бонус</b> newss №2 исчезнет через 1&#8419; минуту.\n",
		"&#11088;<b>Изменение в описании</b> бонуса newss1 №3:\nsdfsd342f\n",
		"&#11088;<b>Изменение</b> в бонусе newss1 №3 (награда 0 секунд):\nwerwerdfg234wersdfsdf\n",
		"&#11088;<b>Новый бонус</b> neddd2 №4\nsdf42f\nНаграда 0 секунд\n<b>Бонусная подсказка:</b>\nwersdf\n",
	}

	compareBonuses(newBonus, oldBonus, confGameENJSON, webToBotTEST)
	for _, bonus := range arrBonuses {
		select {
		// В канал msgChanel будут приходить все новые сообщения from web
		case msgChanel = <-webToBotTEST:
			// конец
			if msgChanel.ChannelMessage == "" {
				return
			}
			// проверка сообщений
			if msgChanel.ChannelMessage != bonus {
				t.Errorf("For new: %v\nold: %v\nexpected %s\ngot %s", newBonus, oldBonus, bonus, msgChanel.ChannelMessage)
			}
		default:
		}
	}
}
func TestEngineCompareMessages(t *testing.T) {
	t.Parallel()

	var confGameENJSON ConfigGameJSON
	confGameENJSON.initTest(pathTest)
	webToBotTEST := make(chan MessengerStyle, 10)
	var msgChanel MessengerStyle

	oldMessages := []MessagesStruct{
		{158796, "Maksisdfsdfmal_adst", 2419, "Magn  Телефоны организаторов:895143Максим\r\n", "otum:  Телефоны организаторов:8915 Мам\r<br/> ", true},
		{158796, "Maksfmal_adst", 2420, "Magn  Телефоны организаторов:895143Максим\r\n", "otum:  Телефоны организаторов:8915 Мам\r<br/> ", true},
	}

	newMessages := []MessagesStruct{
		{158796, "Maksisdfsdfmal_adst", 2419, "Magn  Телефонsdfsdfsdfы организаторов:895143Максим\r\n", "otum:  Телефоны организаторов:8915 Мам\r<br/> ", true},
		{158796, "Maksfвввввввввввввmal_adst", 2420, "Magn  Телефоны организаторов:895143Максим\r\n", "otum:  Телефоны органиjjjjjjjjjjuchangeddddddзаторов:8915 Мам\r<br/> ", true},
		{158796, "newMaksfsdfmal_adst", 2421, "new   Magn  Телефоны органыыыыыыыыыыыизаторов:895143Максим\r\n", "otum:  Телефоны органиыыыыыыыыыыыыыыызаторов:8915 Мам\r<br/> ", true},
	}

	messages := "&#128495;<b>Сообщение изменено:</b>\nMagn  Телефонsdfsdfsdfы организаторов:895143Максим\r\n\n" +
		"&#128495;<b>Появилось сообщение</b>:\n&#128172;new   Magn  Телефоны органыыыыыыыыыыыизаторов:895143Максим\r\n"

	compareMessages(newMessages, oldMessages, confGameENJSON, webToBotTEST)

	select {
	// В канал msgChanel будут приходить все новые сообщения from web
	case msgChanel = <-webToBotTEST:
		// конец
		if msgChanel.ChannelMessage == "" {
			return
		}
		// проверка сообщений
		if msgChanel.ChannelMessage != messages {
			t.Errorf("For new: %v\nold: %v\nexpected %s\ngot %s", newMessages, oldMessages, messages, msgChanel.ChannelMessage)
		}
	default:
	}
}
func TestEngineCompareTasks(t *testing.T) {
	t.Parallel()

	var confGameENJSON ConfigGameJSON
	confGameENJSON.initTest(pathTest)
	webToBotTEST := make(chan MessengerStyle, 10)
	var msgChanel MessengerStyle

	oldTasks := []TaskStruct{
		{true, "sdfsdfsdfsdfsdf", "rr1rr"},
		{true, "", ""},
	}

	newTasks := []TaskStruct{
		{true, "sdfsdf1231231231231sdfsdf", "rrddddddddddddddddd1rr"},
		{true, "", ""},
		{true, "sdfsdfsssssssssssssss", "ddddddddddddd"},
	}

	tasks := []string{"&#10004;<b>Появилось новое задание!</b>\n", "&#10060;<b>Задание изменено</b>:\nsdfsdf1231231231231sdfsdf"}

	compareTasks(newTasks, oldTasks, confGameENJSON, webToBotTEST)
	for _, task := range tasks {
		select {
		// В канал msgChanel будут приходить все новые сообщения from web
		case msgChanel = <-webToBotTEST:
			// конец
			if msgChanel.ChannelMessage == "" {
				return
			}
			// проверка сообщений
			if msgChanel.ChannelMessage != task {
				t.Errorf("For new: %v\nold: %v\nexpected %s\ngot %s", newTasks, oldTasks, task, msgChanel.ChannelMessage)
			}
		default:
		}
	}
}
func TestEngineAddUser(t *testing.T) {
	t.Parallel()
	t.Skip("Test need real new member for team")

	cookieJar, _ := cookiejar.New(nil)
	clientTEST := &http.Client{
		Jar: cookieJar,
	}
	confGameENJSON := ConfigGameJSON{}
	confGameENJSON.initTest(pathTest)
	exitGameTest(clientTEST, confGameENJSON.SubUrl)

	type testPair struct {
		input  string
		output string
	}

	var tests = []testPair{
		{"ERROR_sedfsdfsdf0f5we4fwerwjiejf", "&#10134;Не смогли найти игрока <b>ERROR_sedfsdfsdf0f5we4fwerwjiejf</b>"},
		{"158796", "&#10133;Добавили игрока <b>Maksimal_bot (158796)</b> в команду <b>banda_game (7274)</b>"},
		{"Maksimal_bot", "&#10133;Добавили игрока <b>Maksimal_bot (158796)</b> в команду <b>banda_game (7274)</b>"},
	}

	for _, pair := range tests {
		result := addUser(clientTEST, &confGameENJSON, pair.input)
		if result != pair.output {
			t.Errorf("For %v\nexpected %s\ngot %s", pair.input, pair.output, result)
		}
	}
}
func TestEngineTimeToBonuses(t *testing.T) {
	t.Parallel()

	type testPair struct {
		input  BonusesStruct
		output string
	}

	var tests = []testPair{
		{BonusesStruct{1, "", 1, "", "", true, false, 0, 0, 0}, ""},
		{BonusesStruct{2, "", 2, "", "", true, false, 60, 0, 0}, "&#10004;<b>Бонус</b>  №2 доступен через 1&#8419; минуту.\n"},
		{BonusesStruct{3, "", 3, "", "", true, false, 300, 0, 0}, "&#10004;<b>Бонус</b>  №3 доступен через 5&#8419; минут.\n"},
		{BonusesStruct{4, "", 4, "", "", true, false, 0, 60, 0}, "&#10004;<b>Бонус</b>  №4 исчезнет через 1&#8419; минуту.\n"},
		{BonusesStruct{5, "", 5, "", "", true, false, 0, 300, 0}, "&#10004;<b>Бонус</b>  №5 исчезнет через 5&#8419; минут.\n"},
		{BonusesStruct{6, "", 6, "", "", true, false, 60, 300, 0}, "&#10004;<b>Бонус</b>  №6 доступен через 1&#8419; минуту.\n&#10004;<b>Бонус</b>  №6 исчезнет через 5&#8419; минут.\n"},
	}

	for _, pair := range tests {
		result := timeToBonuses(pair.input)
		if result != pair.output {
			t.Errorf("For %v\nexpected %s\ngot %s", pair.input, pair.output, result)
		}
	}
}

// Exit from site
func exitGameTest(client *http.Client, subUrl string) {
	_, _ = client.Get(fmt.Sprintf("http://%s/Login.aspx?action=logout", subUrl))
}

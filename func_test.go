package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func (conf *ConfigGameJSON) initTest() {
	var configuration ConfigBot
	configuration.init("config.json")

	if value, exists := os.LookupEnv("TestNickName"); exists {
		conf.NickName = value
	} else {
		conf.NickName = configuration.TestNickName
	}

	if value, exists := os.LookupEnv("TestPassword"); exists {
		conf.Password = value
	} else {
		conf.Password = configuration.TestPassword
	}

	if value, exists := os.LookupEnv("TestURLGame"); exists {
		conf.URLGame = value
	} else {
		conf.URLGame = configuration.TestURLGame
	}

	if value, exists := os.LookupEnv("TestLevelNumber"); exists {
		conf.LevelNumber, _ = strconv.ParseInt(value, 10, 64)
	} else {
		conf.LevelNumber = configuration.TestLevelNumber
	}

	conf.separateURL()
}

func TestCompareBonuses(t *testing.T) {
	t.Parallel()
	confGameENJSON := ConfigGameJSON{}
	confGameENJSON.initTest()
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
				t.Error("Получен не тот ответ: ", msgChanel.ChannelMessage, "\nМы ждали: ", bonus)
			}
		default:
		}
	}
}
func TestCompareHints(t *testing.T) {
	t.Parallel()
	confGameENJSON := ConfigGameJSON{}
	confGameENJSON.initTest()
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
				t.Error("Получен не тот ответ: ", msgChanel.ChannelMessage, "\nМы ждали: ", helps)
			}
		default:
		}
	}
}
func TestCompareHintsPenalty(t *testing.T) {
	t.Parallel()
	var confGameENJSON ConfigGameJSON
	confGameENJSON.initTest()
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
				t.Error("Получен не тот ответ:\n", msgChanel.ChannelMessage, "\nМы ждали:\n", help)
			}
		default:
		}
	}
}
func TestCompareMessages(t *testing.T) {
	t.Parallel()
	var confGameENJSON ConfigGameJSON
	confGameENJSON.initTest()
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
			t.Error("Получен не тот ответ:\n", msgChanel.ChannelMessage, "\nМы ждали:\n", messages)
		}
	default:
	}
}
func TestCompareTasks(t *testing.T) {
	t.Parallel()
	var confGameENJSON ConfigGameJSON
	confGameENJSON.initTest()
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
				t.Error("Получен не тот ответ:\n", msgChanel.ChannelMessage, "\nМы ждали:\n", task)
			}
		default:
		}
	}
}

func TestGetLeftCodes(t *testing.T) {
	t.Parallel()
	sectors := []SectorsStruct{{1, 1, "bb", AnswerStruct{}, false}, {2, 2, "bbb", AnswerStruct{}, true}, {3, 3, "bbb", AnswerStruct{}, false}}
	checkSectors := []string{
		"Вам осталось снять сектора:\n\n&#10060;Сектор <b>bb №1</b> не отгадан.\n&#10060;Сектор <b>bbb №3</b> не отгадан.\n",
		"&#10060;Сектор <b>bb №1</b> не отгадан.\n&#10004;Сектор <b>bbb №2</b> отгадан, ответ: \n&#10060;Сектор <b>bbb №3</b> не отгадан.\n"}

	if getLeftCodes(sectors, true) != checkSectors[0] {
		t.Error("Получен не тот ответ: ", getLeftCodes(sectors, true), "\nМы ждали: ", checkSectors[0])
	}

	if getLeftCodes(sectors, false) != checkSectors[1] {
		t.Error("Получен не тот ответ: ", getLeftCodes(sectors, false), "\nМы ждали: ", checkSectors[1])
	}
}
func TestGetFirstBonuses(t *testing.T) {
	t.Parallel()
	var confGameENJSON ConfigGameJSON
	confGameENJSON.initTest()

	bonuses := []BonusesStruct{
		{148016, "123e", 1, "sdfsdf", "werwerwersdfsdf", true, false, 0, 0, 0},
		{148017, "newss", 2, "sdfsdf", "werwerdfgwersdfsdf", true, false, 300, 60, 0},
		{148018, "newss1", 3, "sdfsd342f", "werwerdfg234wersdfsdf", true, false, 0, 0, 0},
		{148019, "neddd2", 4, "sdf42f", "wersdf", true, false, 0, 0, 0},
	}

	checkBonuses :=
		"&#10004;<b>Бонус №1</b> 123e (<b>выполнен</b>, награда: 0 секунд)\nwerwerwersdfsdf\n" +
			"&#128488;<b>Бонус №2</b> newss будет доступен через 5 минут.\nwerwerdfgwersdfsdf\n" +
			"&#10004;<b>Бонус №3</b> newss1 (<b>выполнен</b>, награда: 0 секунд)\nwerwerdfg234wersdfsdf\n" +
			"&#10004;<b>Бонус №4</b> neddd2 (<b>выполнен</b>, награда: 0 секунд)\nwersdf\n"

	if getFirstBonuses(bonuses, confGameENJSON) != checkBonuses {
		t.Error("Не верно взяты бонусы. Получен ответ:", getFirstBonuses(bonuses, confGameENJSON), "\nМы ждали:", checkBonuses)
	}
}
func TestFirstSectors(t *testing.T) {
	t.Parallel()
	tasks := []LevelStruct{
		{270768, "rr1rr", 9, 0, 0, 0, false, false, StartTimeStruct{5}, false, 1, 1, 1, 1, 3, 1, 1, nil, nil, nil, []SectorsStruct{{1, 2, "bb", AnswerStruct{}, false}, {2, 3, "bbb", AnswerStruct{}, false}, {3, 3, "bbb", AnswerStruct{}, false}}, nil, nil, nil},
		{270768, "rr2rr", 9, 100, 10, 60, false, false, StartTimeStruct{5}, false, 0, 1, 1, 1, 0, 0, 0, nil, nil, nil, nil, nil, nil, nil}}

	sectors := []string{"&#128269;Вам нужно найти <b>1 из 3</b> секторов.\n", "&#128269;На уровне 1 код.\n"}

	for number, task := range tasks {
		if getFirstSector(task) != sectors[number] {
			t.Error("Не верно взяты сектора уровня. Получен ответ:", getFirstSector(task), "\nМы ждали:", sectors[number])
		}
	}
}
func TestFirstTimer(t *testing.T) {
	t.Parallel()
	tasks := []LevelStruct{
		{270768, "rr1rr", 9, 0, 0, 0, false, false, StartTimeStruct{5}, false, 0, 1, 0, 1, 7, 2, 5, nil, nil, nil, nil, nil, nil, nil},
		{270768, "rr2rr", 9, 100, 10, 60, false, false, StartTimeStruct{5}, false, 0, 1, 0, 1, 7, 2, 5, nil, nil, nil, nil, nil, nil, nil}}
	times := []string{"&#9203;<b>Автоперехода нет</b>\n", "&#9200;До автоперехода 1 минута\n&#128276;Штраф за автопереход -10 секунд\n"}

	for number, task := range tasks {
		if getFirstTimer(task) != times[number] {
			t.Error("Не верно взято время уровня. Получен ответ:", getFirstTimer(task), "\nМы ждали:", times[number])
		}
	}
}
func TestFirstTask(t *testing.T) {
	t.Parallel()
	var confGameENJSON ConfigGameJSON
	confGameENJSON.initTest()

	task := []TaskStruct{
		{true, "sdfsdfsdfsdfsdf", "rr1rr"},
		{true, "", ""}}

	rightTask := []string{"&#10060;<b>Задания отсуствуют!</b>\n", "\n&#9889;<b>Задание</b>:\nrr1rr\n\n&#10060;Текст задания <b>отсуствует</b>!\n"}

	if getFirstTask([]TaskStruct{}, confGameENJSON) != rightTask[0] {
		t.Error("Не верно взяты подсказки. Получен ответ:", getFirstTask([]TaskStruct{}, confGameENJSON), "\nМы ждали:", rightTask[0])
	}

	if getFirstTask(task, confGameENJSON) != rightTask[1] {
		t.Error("Не верно взяты подсказки. Получен ответ:", getFirstTask(task, confGameENJSON), "\nМы ждали:", rightTask[1])
	}

}
func TestFirstHelps(t *testing.T) {
	t.Parallel()
	var confGameENJSON ConfigGameJSON
	confGameENJSON.initTest()

	helps := []HelpsStruct{
		{243732, 1, "sdfsdfdsfsdf", false, 0, "", false, 0, 10},
		{243732, 2, "sdfsdf234dsfsdf", false, 0, "", false, 0, 300},
		{243732, 3, "sdfs1111sfsdf", false, 10, "", false, 0, 0},
	}

	rightHelps := []string{"", "&#10004;<b>Подсказка</b> №1 будет через 10 секунд.\n\n&#10004;<b>Подсказка</b> №2 будет через 5 минут.\n\n&#10004;<b>Подсказка</b> №3 \nsdfs1111sfsdf\n\n", "&#10004;<b>Штрафная подсказка</b> №4:\nsdfs3233241sfsdf\n\n"}
	penaltyHelps := []HelpsStruct{{243732, 4, "sdfs3233241sfsdf", true, 10, "ddwww", true, 0, 0}}

	if getFirstHelps([]HelpsStruct{}, confGameENJSON) != rightHelps[0] {
		t.Error("Не верно взяты подсказки. Получен ответ:", getFirstHelps([]HelpsStruct{}, confGameENJSON), "\nМы ждали:", rightHelps[0])
	}
	if getFirstHelps(helps, confGameENJSON) != rightHelps[1] {
		t.Error("Не верно взяты подсказки. Получен ответ:", getFirstHelps(helps, confGameENJSON), "\nМы ждали:", rightHelps[1])
	}
	if getFirstHelps(penaltyHelps, confGameENJSON) != rightHelps[2] {
		t.Error("Не верно взяты штрафные подсказки. Получен ответ:", getFirstHelps(penaltyHelps, confGameENJSON), "\nМы ждали:", rightHelps[2])
	}
}
func TestFirstMessages(t *testing.T) {
	t.Parallel()
	var confGameENJSON ConfigGameJSON
	confGameENJSON.initTest()

	messages := []MessagesStruct{{158796, "Maksisdfsdfmal_adst", 2419, "   Magn  Телефоны организаторов:895143Максим\r\n", "otum:  Телефоны организаторов:8915 Мам\r<br/> ", true}}
	rightMessages := []string{"", "&#128495;<b>Сообщения на уровне:</b>\n&#128172;   Magn  Телефоны организаторов:895143Максим\r\n"}

	if getFirstMessages([]MessagesStruct{}, confGameENJSON) != rightMessages[0] {
		t.Error("Не верно взяты пустые сообщения. Получен ответ:", getFirstMessages([]MessagesStruct{}, confGameENJSON), "\nМы ждали:", rightMessages[0])
	}
	if getFirstMessages(messages, confGameENJSON) != rightMessages[1] {
		t.Error("Не верно взяты сообщения. Получен ответ:", getFirstMessages(messages, confGameENJSON), "\nМы ждали:", rightMessages[0])
	}
}

func TestGameEngineModel(t *testing.T) {
	cookieJar, _ := cookiejar.New(nil)
	clientTEST := &http.Client{
		Jar: cookieJar,
	}

	confGameENJSON := ConfigGameJSON{}
	confGameENJSON.initTest()

	// ВЫХОД
	str := fmt.Sprintf("http://%s/Login.aspx?action=logout", confGameENJSON.SubUrl)
	_, _ = clientTEST.Get(str)

	enterGameJSON(clientTEST, confGameENJSON)
	answer := gameEngineModel(clientTEST, confGameENJSON)
	if answer.Event != 0 {
		t.Error("Игра НЕ в нормальном состоянии! Код ошибки: ", answer.Event)
	}
}
func TestSentCodeJSON(t *testing.T) {
	t.Parallel()
	t.SkipNow()
	rand.Seed(time.Now().UTC().UnixNano()) // real random

	webToBotTEST := make(chan MessengerStyle, 10)
	var msgChanel MessengerStyle

	cookieJar, _ := cookiejar.New(nil)
	clientTEST := &http.Client{
		Jar: cookieJar,
	}

	confGameENJSON := ConfigGameJSON{}
	confGameENJSON.initTest()

	// ВЫХОД
	str := fmt.Sprintf("http://%s/Login.aspx?action=logout", confGameENJSON.SubUrl)
	_, _ = clientTEST.Get(str)

	enterGameJSON(clientTEST, confGameENJSON)
	confGameENJSON.LevelNumber = gameEngineModel(clientTEST, confGameENJSON).Level.Number

	code := fmt.Sprintf("НЕВЕРНЫЙ%d", rand.Int())
	var isBonus *bool
	*isBonus = true
	sentCodeJSON(clientTEST, &confGameENJSON, code, isBonus, webToBotTEST, 0)
	select {
	// В канал msgChanel будут приходить все новые сообщения from web
	case msgChanel = <-webToBotTEST:
		// конец
		fmt.Println(msgChanel.ChannelMessage)
		if msgChanel.ChannelMessage == "" {
			return
		}
		fmt.Println(msgChanel.ChannelMessage)
		if msgChanel.ChannelMessage != "Код &#10060;<b>НЕВЕРНЫЙ</b>" {
			t.Error("КОД ОТПРАВЛЕН С ОШИБКОЙ! Мы получили: ", msgChanel.ChannelMessage, "\n Мы ждали: Код &#10060;<b>НЕВЕРНЫЙ</b>")
		}
	default:
	}
}
func TestEnterGameENJSON(t *testing.T) {

	cookieJar, _ := cookiejar.New(nil)
	clientTEST := &http.Client{
		Jar: cookieJar,
	}

	var confGameENJSON ConfigGameJSON
	confGameENJSON.initTest()

	// ВЫХОД
	str := fmt.Sprintf("http://%s/Login.aspx?action=logout", confGameENJSON.SubUrl)
	_, _ = clientTEST.Get(str)

	answer := enterGameJSON(clientTEST, confGameENJSON)

	if answer != fmt.Sprintf("&#10004;<b>Авторизация прошла успешно</b> на игру: %s", confGameENJSON.URLGame) {
		t.Error(
			"Невозможно авторизоваться. Под:", confGameENJSON,
		)
	}
}

func TestGetPenaltyJSON(t *testing.T) {
	t.Parallel()
	t.SkipNow()
	rand.Seed(time.Now().UTC().UnixNano()) // real random

	webToBotTEST := make(chan MessengerStyle, 10)
	var msgChanel MessengerStyle

	cookieJar, _ := cookiejar.New(nil)
	clientTEST := &http.Client{
		Jar: cookieJar,
	}

	confGameENJSON := ConfigGameJSON{}
	confGameENJSON.initTest()

	// ВЫХОД
	str := fmt.Sprintf("http://%s/Login.aspx?action=logout", confGameENJSON.SubUrl)
	_, _ = clientTEST.Get(str)

	enterGameJSON(clientTEST, confGameENJSON)

	penaltyID := 1111
	getPenaltyJSON(clientTEST, confGameENJSON, penaltyID, webToBotTEST)
	select {
	// В канал msgChanel будут приходить все новые сообщения from web
	case msgChanel = <-webToBotTEST:
		// конец
		fmt.Println(msgChanel.ChannelMessage)
		if msgChanel.ChannelMessage == "" {
			return
		}
		if !strings.Contains(msgChanel.ChannelMessage, "&#9889;Штрафная подсказка") {
			t.Error("Ошибка при взятии штрафной подсказки. Мы получили: ", msgChanel.ChannelMessage, "\n Мы ждали: &#9889;Штрафная подсказка")
		}
	default:
	}
}

func TestConvertTimeSec(t *testing.T) {
	t.Parallel()

	type testPair struct {
		original int
		replaced string
	}

	var tests = []testPair{
		{0, "0 секунд"},
		{1, "1 секунда"},
		{60, "1 минута"},
		{66, "1 минута 6 секунд"},
		{120, "2 минуты"},
		{122, "2 минуты 2 секунды"},
		{600, "10 минут"},
		{3600, "1 час"},
		{3601, "1 час 1 секунда"},
		{3661, "1 час 1 минута 1 секунда"},
		{86400, "1 день"},
		{90061, "1 день 1 час 1 минута 1 секунда"},
		{36045645, "417 дней 4 часа 45 минут 45 секунд"},
	}

	for _, pair := range tests {
		v := convertTimeSec(pair.original)
		if v != pair.replaced {
			t.Error(
				"For", pair.original,
				"\nexpected", pair.replaced,
				"\ngot", v,
			)
		}
	}
}
func TestReplaceTag(t *testing.T) {
	t.Parallel()

	type testPair struct {
		original string
		replaced string
	}

	var tests = []testPair{
		{"<span> sdf </span>", " sdf "},
		{"< /p >", "\n"},
		{"q</h3>w", "q</b>\nw"},
		{"< H4>", "<b>"},
		{"sdfs", "sdfs"},
		{`s<df>s`, "ss"},
		{`qqqqqqqqqqq 40.167841 58.410761 or 40,167841 58,410761 fgdfgdg`, `qqqqqqqqqqq <code>40.167841,40.167841</code> <a href="https://maps.google.com/?q=40.167841,40.167841">[G]</a> <a href="https://yandex.ru/maps/?source=serp_navig&text=40.167841,40.167841">[Y]</a>, or <code>40.167841,40.167841</code> <a href="https://maps.google.com/?q=40.167841,40.167841">[G]</a> <a href="https://yandex.ru/maps/?source=serp_navig&text=40.167841,40.167841">[Y]</a>, fgdfgdg`},
		{`<iframe width="560" height="315" src="https://www.youtube.com/embed/PNf54u-T-s0?start=600" frameborder="0" allow="autoplay; encrypted-media" allowfullscreen></iframe>`, `[Вставлен iframe с другого сайта: <a href="https://www.youtube.com/embed/PNf54u-T-s0?start=600">https://www.youtube.com/embed/PNf54u-T-s0?start=600</a>]`},
		{`<a href="http://d1.endata.cx/data/games/58762/34vfmdshvf_b.jpg" border="0"><img src="http://d1.endata.cx/data/games/58762/34vfmdshvf_s.jpg" title="fds"></a>`, `[Спрятанная ссылка: <a href="http://d1.endata.cx/data/games/58762/34vfmdshvf_b.jpg">http://d1.endata.cx/data/games/58762/34vfmdshvf_b.jpg</a> под картинкой: <a href="http://d1.endata.cx/data/games/58762/34vfmdshvf_s.jpg">http://d1.endata.cx/data/games/58762/34vfmdshvf_s.jpg</a>]`},
		{`<img src="http://d1.endata.cx/data/games/58762/34vfmdshvf_s.jpg" title="fds"></a>`, `[Картинка: <a href="http://d1.endata.cx/data/games/58762/34vfmdshvf_s.jpg">http://d1.endata.cx/data/games/58762/34vfmdshvf_s.jpg</a>]`},
		{`<img src="http://d1.endata.cx/data/games/58762/34vfmdshvf_s.jpg"></a>`, `[Картинка: <a href="http://d1.endata.cx/data/games/58762/34vfmdshvf_s.jpg">http://d1.endata.cx/data/games/58762/34vfmdshvf_s.jpg</a>]`},
		{`<p><img src="http://cdn.endata.cx/data/games/60896/6-efcfcbvg.jpg">`, `[Картинка: <a href="http://cdn.endata.cx/data/games/60896/6-efcfcbvg.jpg">http://cdn.endata.cx/data/games/60896/6-efcfcbvg.jpg</a>]`},
		{`<img src="http://d1.endata.cx/data/games/58762/34vfmdshvf_s.jpg"title="fds"></a>`, `[Картинка: <a href="http://d1.endata.cx/data/games/58762/34vfmdshvf_s.jpg">http://d1.endata.cx/data/games/58762/34vfmdshvf_s.jpg</a>]`},
		{`<img  id="image_" src="http://d1.endata.cx/data/games/58762/34vfmdshvf_s.jpg"title="fds"></a>`, `[Картинка: <a href="http://d1.endata.cx/data/games/58762/34vfmdshvf_s.jpg">http://d1.endata.cx/data/games/58762/34vfmdshvf_s.jpg</a>]`},
		{`<a href="http://d1.endata.cx/data/games/58762/34vfmdshvf_b.jpg" ><img src="http://d1.endata.cx/data/games/58762/34vfmdshvf_s.jpg"></a>`, `[Спрятанная ссылка: <a href="http://d1.endata.cx/data/games/58762/34vfmdshvf_b.jpg">http://d1.endata.cx/data/games/58762/34vfmdshvf_b.jpg</a> под картинкой: <a href="http://d1.endata.cx/data/games/58762/34vfmdshvf_s.jpg">http://d1.endata.cx/data/games/58762/34vfmdshvf_s.jpg</a>]`},
		{`<a href="http://d1.endata.cx/data/games/60754/%D0%BC2%D0%B0%D0%BF%D0%B2%D0%B0%D0%B2%D0%B0%D0%BD%D1%80.jpeg"><img id="image_" src="http://d1.endata.cx/data/games/60754/%D0%BC2%D0%B0%D0%BF%D0%B2%D0%B0%D0%B2%D0%B0%D0%BD%D1%80.jpeg"></a>`,
			`[Спрятанная ссылка: <a href="http://d1.endata.cx/data/games/60754/м2апваванр.jpeg">http://d1.endata.cx/data/games/60754/м2апваванр.jpeg</a> под картинкой: <a href="http://d1.endata.cx/data/games/60754/м2апваванр.jpeg">http://d1.endata.cx/data/games/60754/м2апваванр.jpeg</a>]`},
	}

	for _, pair := range tests {
		v := replaceTag(pair.original, "RxR")
		if v != pair.replaced {
			t.Error(
				"For", pair.original,
				"\nexpected", pair.replaced,
				"\ngot", "|"+v+"|",
			)
		}
	}
}
func TestDeleteMap(t *testing.T) {
	t.Parallel()

	maps := make(map[int]string)
	maps[0] = "asdas1"
	maps[1] = "asdas2"
	maps[3] = "asdas3"
	maps[3] = "asdas4"

	deleteMap(maps)

	var v int
	v = len(maps)
	if v != 0 {
		t.Error(
			"For delete len map ", 3,
			"expected", len(maps),
			"got", v,
		)
	}
}
func TestDeleteMapFloat(t *testing.T) {
	t.Parallel()

	maps := make(map[string]float64)
	maps["asdas1"] = 45.4
	maps["asdas2"] = 145.423
	maps["asdas3"] = 415.42
	maps["asdas4"] = 45.43

	deleteMapFloat(maps)

	var v int
	v = len(maps)
	if v != 0 {
		t.Error(
			"For delete len map ", 3,
			"expected", len(maps),
			"got", v,
		)
	}
}
func TestSearchPhoto(t *testing.T) {
	t.Parallel()

	PhotoMapGood := make(map[int]string)
	PhotoMapGood[0] = `<img src="http://ya.ru/iads.png">`
	PhotoMapGood[1] = `<img src="http://sdfsd.png" title="sdfsf">`

	text := `<img src="http://ya.ru/iads.png">jdsfjklsjklds</b>fjklsdfjklsdfjklsd<img src="http://sdfsd.png" title="sdfsf">dfkjfjklsdfkjlsd`
	photoMap := searchLocation(text)

	i := 0
	for item := range photoMap {
		if item != PhotoMapGood[i] {
			t.Error(
				"For", item,
				"Expected", PhotoMapGood[i],
			)
		}
		i++
	}
}
func TestConvertUTF(t *testing.T) {
	type testPair struct {
		original string
		replaced string
	}

	var tests = []testPair{
		{`http://d1.endata.cx/data/games/60070/%D0%B0%D0%B3%D0%B5%D0%BD%D1%82+%D1%8F%D0%B3%D1%83%D0%BB.jpg`, `http://d1.endata.cx/data/games/60070/агент ягул.jpg`},
		{`http://en.cx.gameengines/encounter/play/27543/?pid=237876&amp;pact=2`, `http://en.cx.gameengines/encounter/play/27543/?pid=237876&amp;pact=2`},
		{`http://d1.endata.cx/data/games/63892/%d0%bc%d0%b5%d1%82%d0%ba%d0%b0+copy.jpg`, `http://d1.endata.cx/data/games/63892/метка copy.jpg`},
		{`http://d1.endata.cx/data/games/63892/%d0%bc%d0%b5%d1%82%d0%ba%d0%b0+copy.jpg`, `http://d1.endata.cx/data/games/63892/метка copy.jpg`},
	}

	for _, pair := range tests {
		v, _ := convertUTF(pair.original)
		if v != pair.replaced {
			t.Error(
				"For", pair.original,
				"\nexpected", pair.replaced,
				"\ngot", v,
			)
		}
	}

}
func TestAutoCode(t *testing.T) {
	t.Parallel()

	type testPair struct {
		original string
		code     string
	}

	var tests = []testPair{
		//{"/ac", "ЕЛОП ЛssОПЕ ПОЛЕ "},
		{"18", "<b>18</b> регион = Удмуртская Республика;\n"},
		{"18 45", "<b>18</b> регион = Удмуртская Республика;\n<b>45</b> регион = Курганская область;\n"},
		{"электростанция", "Региона электростанция в России нет."},
	}

	for _, pair := range tests {
		v := autoCode(pair.original)
		if v != pair.code {
			t.Error(
				"Для", pair.original,
				"/nОжидали", pair.code,
				"/nПолучили", v,
			)
		}
	}
}
func TestSearchLocation(t *testing.T) {
	t.Parallel()

	locMapGood := make(map[string]float64)

	locMapGood["Latitude0"] = 40.2522222
	locMapGood["Longitude0"] = 58.4363889

	locMapGood["Latitude1"] = 40.167841
	locMapGood["Longitude1"] = 58.410761

	locMapGood["Latitude2"] = 40.167845
	locMapGood["Longitude2"] = 58.410765

	locMapGood["Latitude3"] = 59.413521
	locMapGood["Longitude3"] = 58.410761

	locMapGood["Latitude4"] = 56.84471667
	locMapGood["Longitude4"] = 53.19626667

	text := `dasdada 40°15′08″ 58°26′23″ 
			в. д.HGЯO sdfsdsdasd <b>dfsfsdf</b>\n 40.167841<br/>58.410761 \n 
			40,167845<br/>58,410765 //
			fsdfsd 59°24'48.6756", 58°24'38.7396"
			sdfsdgi :: 56°50.683, 53°11.776`
	locationMap := searchLocation(text)

	for i := 0; i == len(locationMap)/2; i++ {
		if (locationMap["Latitude"+strconv.FormatInt(int64(i), 10)] != locMapGood["Latitude"+strconv.FormatInt(int64(i), 10)]) && (locationMap["Longitude"+strconv.FormatInt(int64(i), 10)] != locMapGood["Longitude"+strconv.FormatInt(int64(i), 10)]) {
			t.Error(
				"Latitude for", locationMap["Latitude"+strconv.FormatInt(int64(i), 10)],
				"\nLatitude expected", locMapGood["Latitude"+strconv.FormatInt(int64(i), 10)],
				"\nLongitude for", locationMap["Longitude"+strconv.FormatInt(int64(i), 10)],
				"\nLongitude expected", locMapGood["Longitude"+strconv.FormatInt(int64(i), 10)],
			)
		}
	}
}
func TestAnagram(t *testing.T) {
	t.Parallel()

	type testPair struct {
		original string
		replaced string
	}

	var tests = []testPair{
		{"/ano плео", "ЕЛОП ЛОПЕ ПОЛЕ "},
		{"/ano фывфывйуфывфы", "<b>Слов не обнаружено.</b>"},
		{"/ano электростанция", "ЭЛЕКТРОСТАНЦИЯ "},
	}

	for _, pair := range tests {
		v, _ := anagram(pair.original)
		if v != pair.replaced {
			t.Error(
				"Для", pair.original,
				"/nОжидали", pair.replaced,
				"/nПолучили", v,
			)
		}
	}
}
func TestSearchForMask(t *testing.T) {

	t.Parallel()

	type testPair struct {
		original string
		replaced string
	}

	var tests = []testPair{
		{"по?е", "ПОБЕ ПОЛЕ "},
		{"фывфывйуфывфы", "<b>Слов не обнаружено.</b>"},
		{"пол*", "ПОЛ ПОЛА ПОЛАБЗИТЬСЯ ПОЛАБЫ ПОЛАВИРУЕМ ПОЛАВОЧНИК ПОЛАВОЧНЫЙ ПОЛАГАНИЕ ПОЛАД ПОЛАДКА ПОЛАЗ ПОЛАЗНА ПОЛАЗНИК ПОЛАЙКА ПОЛАК ПОЛАКАНТ ПОЛАКОМИЛ ПОЛАКСИС ПОЛАН ПОЛАНДРА ПОЛАНИ ПОЛАНСКИ ПОЛАНСКИЙ ПОЛАНЬИ ПОЛАР ПОЛАРОИД ПОЛАТИ ПОЛАТКА ПОЛАШ ПОЛБА ПОЛБИГ ПОЛБИН ПОЛБОТИНКА ПОЛБУТЫЛКА ПОЛБУТЫЛКИ ПОЛВА ПОЛВЕДРА ПОЛВЕКА ПОЛВЕРШКА ПОЛВЕТРА ПОЛГАР ПОЛГАРНЕЦ ПОЛГОДА ПОЛГОРЯ ПОЛДЕЛА ПОЛДЕНЬ ПОЛДНИК ПОЛДНИЧАНИЕ ПОЛДНИЧАНЬЕ ПОЛДНЯ ПОЛДОРОГИ ПОЛДЮЖИНЫ ПОЛДЮЙМА ПОЛЕ ПОЛЕВАНОВ ПОЛЕВАНЬЕ ПОЛЕВЕНИЕ ПОЛЕВИК ПОЛЕВИЦА ПОЛЕВИЧКА ПОЛЕВКА ПОЛЕВОД ПОЛЕВОДСТВО ПОЛЕВОЙ ПОЛЕВСКОЙ ПОЛЕВЩИК ПОЛЕГАДА ПОЛЕГАЕМОСТЬ ПОЛЕГАНИЕ ПОЛЕЖАЕВ ПОЛЕЖАЕВСКАЯ ПОЛЕЖАЕВСКИЙ ПОЛЕЖАЛОЕ ПОЛЕЗНО ПОЛЕЗНОЕ ПОЛЕЗНОСТЬ ПОЛЕЗНЫЙ ПОЛЕЛЕЯЛА ПОЛЕЛЬ "},
		{"электростанция", "ЭЛЕКТРОСТАНЦИЯ "},
	}

	for _, pair := range tests {
		v, _ := searchForMask(pair.original)
		if v != pair.replaced {
			t.Error(
				"Для", pair.original,
				"Ожидали", pair.replaced,
				"Получили", v,
			)
		}
	}
}
func TestAssociations(t *testing.T) {
	t.Parallel()

	type testPair struct {
		original string
		replaced string
	}

	var tests = []testPair{
		//		{"/ass поле", "пшеница луг простор трава лес футбол цветы воин пашня рожь сено степь трактор комбайн кукуруза колос конь ветер ромашки гольф васильки мак нива битва деревня пахарь магнит василёк сенокос природа ковыль огород поляна подсолнух зерно урожай игра лето клевер бабочки равнина земля перекати-поле один овёс стог колхоз колосья пахота комбайнер травы жатва брань ромашка коровы просторы хлеб сеялка русь даль кони село ширь коса корова колосок ягода плуг жизнь чудеса солома дорога небо русское поле посев песня борозда война раздолье работа пастух трактора цветок мина труженик одуванчики конопля пространство земляника злак лютики жёлтый цвет бахча мяч сила лужайка шум путь жаворонок россия лён река дерево "},
		{"/ass фывфывйуфывфы", "<b>Слов не обнаружено.</b>"},
		//		{"/ass электростанция", "провода гоэлро лэп энергия тэс электрификация аэс коммунизм свет ток генератор топливо гэс турбина лампочка пуск течение афганистан прилив обнинск тесла стройка ветер напряжение солнечные батареи двигатель плотина ротор атом курчатов атомы электротехника ленин фотосинтез "},
	}

	for _, pair := range tests {
		v, _ := associations(pair.original)
		if v != pair.replaced {
			t.Error(
				"Для", pair.original,
				"\nОжидали", pair.replaced,
				"\nПолучили", v,
			)
		}
	}
}
func TestTransferToAlphabetToWord(t *testing.T) {
	t.Parallel()

	type testPair struct {
		original  string
		replaced  string
		attribute bool
	}

	var testsN2W = []testPair{
		{"4", "4 = <b>RU</b> г <b>EN</b> d;\n", true},
		{"4 q", "4 = <b>RU</b> г <b>EN</b> d;\nq = недопустимое число!\n", true},
		{"30", "30 = <b>RU</b> ь <b>EN</b> -;\n", true},
		{"33", "33 = недопустимое число!\n", true},
		{"900", "900 = недопустимое число!\n", true},
		{"900 4", "900 = недопустимое число!\n4 = <b>RU</b> г <b>EN</b> d;\n", true},
		{"3 4 5 6", "3 = <b>RU</b> в <b>EN</b> c;\n4 = <b>RU</b> г <b>EN</b> d;\n5 = <b>RU</b> д <b>EN</b> e;\n6 = <b>RU</b> е <b>EN</b> f;\n", true},
	}
	var testsW2N = []testPair{
		{"a", "Номер A в <b>EN</b> алфавите 1;\n", false},
		{"Я", "Номер Я в <b>RU</b> алфавите 33;\n", false},
		{"Я 48 u", "Номер Я в <b>RU</b> алфавите 33;\nНомер U в <b>EN</b> алфавите 21;\n", false},
		{"а б в", "Номер А в <b>RU</b> алфавите 1;\nНомер Б в <b>RU</b> алфавите 2;\nНомер В в <b>RU</b> алфавите 3;\n", false},
	}

	for _, pair := range testsN2W {
		v := transferToAlphabet(pair.original, true)
		//v := transferToAlphabet(pair.original, pair.attribute)
		if v != pair.replaced {
			t.Error(
				"Для", pair.original,
				"\nОжидали", pair.replaced,
				"\nПолучили", v,
			)
		}
	}
	for _, pair := range testsW2N {
		v := transferToAlphabet(pair.original, false)
		//v := transferToAlphabet(pair.original, pair.attribute)
		if v != pair.replaced {
			t.Error(
				"Для", pair.original,
				"\nОжидали", pair.replaced,
				"\nПолучили", v,
			)
		}
	}
}
func TestMendeleevTable(t *testing.T) {
	t.Parallel()

	type testPair struct {
		original string
		replaced string
	}

	var tests = []testPair{
		{"3", "Номер 3 = Li - Литий - масса 6,941;\n"},
		{"3 9 ", "Номер 3 = Li - Литий - масса 6,941;\nНомер 9 = F - Фтор - масса 18,9984;\n"},
		{"48 в", "Номер 48 = Cd - Кадмий - масса 112,41;\nв = не допустимое число!\n"},
		{"в", "в = не допустимое число!\n"},
		{"1999", "Номера 1999 не существует в таблице.\n"},
		{"Zr", "Номер Zr = Zr - Цирконий - масса 91,22;\n"},
	}

	for _, pair := range tests {
		v := tableMendeleev(pair.original)
		if v != pair.replaced {
			t.Error(
				"Для", pair.original,
				"/nОжидали", pair.replaced,
				"/nПолучили", v,
			)
		}
	}
}
func TestMorse(t *testing.T) {
	t.Parallel()

	type testPair struct {
		original string
		replaced string
	}

	var tests = []testPair{
		{".-.-", ".-.- = <b>RU</b> Я <b>EN</b> Ä;\n"},
		{"sf", "Символа sf азбуке морзе нет.\n"},
		{".. ---", ".. = <b>RU</b> И <b>EN</b> I;\n--- = <b>RU</b> О <b>EN</b> O;\n"},
	}

	for _, pair := range tests {
		v := morse(pair.original)
		if v != pair.replaced {
			t.Error(
				"Для", pair.original,
				"\nОжидали", pair.replaced,
				"\nПолучили", v,
			)
		}
	}
}
func TestBrail(t *testing.T) {
	//t.Parallel()

	type testPair struct {
		original string
		replaced string
	}

	var tests = []testPair{
		{"100000", "100000 = <b>RU</b> А <b>EN</b> A <b>№</b> 1;\n"},
		{"sf", "Символа sf шрифте Браиля нет.\n"},
		{"100000 101000", "100000 = <b>RU</b> А <b>EN</b> A <b>№</b> 1;\n101000 = <b>RU</b> Б <b>EN</b> B <b>№</b> 2;\n"},
	}

	for _, pair := range tests {
		v := braille(pair.original)
		if v != pair.replaced {
			t.Error(
				"Для", pair.original,
				"\nОжидали", pair.replaced,
				"\nПолучили", v,
			)
		}
	}
}
func TestBin(t *testing.T) {
	t.Parallel()

	type testPair struct {
		original  string
		replaced  string
		attribute bool
	}

	var tests = []testPair{
		{"999", "999 (10) = 1111100111 (2);\n", true},
		{"999 19", "999 (10) = 1111100111 (2);\n19 (10) = 10011 (2);\n", true},
		{"fff", "fff = недопустимое число!\n", true},
		{"1111100111", "1111100111 - изначально двоичка;\n", true},

		{"1111100111", "1111100111 (2) = 999 (10);\n", false},
		{"aaa", "aaa = недопустимое число!\n", false},
		{"985", "985 (2) = 57 (10);\n", false},
		{"100 11111", "100 (2) = 4 (10);\n11111 (2) = 31 (10);\n", false},
	}

	for _, pair := range tests {
		v := bin(pair.original, pair.attribute)
		if v != pair.replaced {
			t.Error(
				"Для", pair.original,
				"\nОжидали", pair.replaced,
				"\nПолучили", v,
			)
		}
	}
}

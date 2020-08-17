package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

func (conf *ConfigGameJSON) init(str string) string {
	args := strings.Split(str, " ")
	if len(args) < 3 {
		return "Need more arguments! <code>/start login password http://DEMO.en.cx/GameDetails.aspx?gid=1</code>"
	}
	if len(args) > 3 {
		return "Слишком много аргументов!"
	}
	conf.NickName = args[0]
	conf.Password = args[1]
	conf.URLGame = args[2]

	pathUrl := strings.Split(conf.URLGame, "/")
	if len(pathUrl) > 1 {
		conf.SubUrl = pathUrl[2]
	}
	pathUrl = strings.Split(conf.URLGame, "=")
	if len(pathUrl) > 0 {
		conf.Gid = pathUrl[1]
	}
	return ""
}

func (conf *ConfigBot) init(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conf)
	if err != nil {
		log.Println(err)
	}
	if value, exists := os.LookupEnv("TelegramBotToken"); exists {
		conf.TelegramBotToken = value
	}
	if value, exists := os.LookupEnv("OwnName"); exists {
		conf.OwnName = value
	}
}

/*
	Config our bot
*/
type ConfigBot struct {
	TelegramBotToken string
	OwnName          string
	TestNickName     string
	TestPassword     string
	TestURLGame      string
	TestLevelNumber  int64
	Jokes            []string
}

/*
	Style all message between bot and telegram
*/
type MessengerStyle struct {
	Latitude       float64
	Longitude      float64
	ChannelMessage string
	Type           string
	MsgId          int
	ChatId         int64
}

type ConfigGameJSON struct {
	NickName    string
	Password    string
	URLGame     string
	SubUrl      string
	Gid         string
	LevelNumber int64
	Postfix     string
	Prefix      string
}

type Coordinate struct {
	Latitude  float64
	Longitude float64
}

// Оценка текущего состояния игры
type Model struct {
	Event         int                `json:"Event"`      // Отражает в каком состоянии находится игра
	GameId        int                `json:"GameId"`     //	ID игры
	GameNumber    int                `json:"GameNumber"` //	номер игры
	GameTitle     string             `json:"GameTitle"`  //название игры
	GameTypeID    int                `json:"GameTypeId"`
	GameZoneID    int                `json:"GameZoneId"`
	LevelSequence int                `json:"LevelSequence"` //	тип последовательнсти: 0 – линейная, 1 –указанная, 2 – случайная, 3 – штурмовая, 4 – динам. случайная
	UserId        int                `json:"UserId"`        //	ID игрока
	TeamId        int                `json:"TeamId"`        //	ID команды игрока
	EngineAction  EngineActionStruct `json:"EngineAction"`  //	информация о результате последнего запроса игрока
	Level         LevelStruct        `json:"Level"`         //	информация о текущем уровне
	Levels        []LevelStruct      `json:"Levels"`        //	список всех уровней
}

// Оценка текущего состояния уровня
type LevelStruct struct {
	LevelId              int                  // ID Уровня
	Name                 string               // Имя уровня
	Number               int64                // Номер уровня
	Timeout              int                  // время в секундах  срабатывания автоперехода, 0 – если нет
	TimeoutAward         int                  // штраф за автопереход в секундах, 0 – если нет
	TimeoutSecondsRemain int                  // осталось времени до срабатывания автоперехода в секундах)
	IsPassed             bool                 // уровень пройден
	Dismissed            bool                 // уровень снять администратором
	StartTime            StartTimeStruct      //struct{} `json:"StartTime"` // время начала уровня для игрока
	HasAnswerBlockRule   bool                 // есть ли на уровне блокировка ответов
	BlockDuration        int                  // осталось секунду блокировки; 0 – не активна
	BlockTargetId        int                  // блокировка установлена для: 0,1 – для игрока; 2 – для команды
	AttemtsNumber        int                  // количество попыток разрешенных в рамках AttemtsPeriod
	AttemtsPeriod        int                  // период срабатывания блокировки в секундах)
	RequiredSectorsCount int                  // Количество секторов, которые необходимо отгадать
	PassedSectorsCount   int                  // Количество отгаданных секторов
	SectorsLeftToClose   int                  // Количество неотгаданных секторов
	MixedActions         []MixedActionsStruct // История введенных ответов
	Messages             []MessagesStruct     // Сообщения администратора
	Tasks                []TaskStruct         `json:"Tasks"` // Текст задания
	Sectors              []SectorsStruct      // Сектора
	Helps                []HelpsStruct        // Подсказки
	PenaltyHelps         []HelpsStruct        // Штрафные подсказки
	Bonuses              []BonusesStruct      `json:"Bonuses"` // Бонусные задания
}

type StartTimeStruct struct {
	Value float64
}
type EngineActionStruct struct {
	GameId      int          //	ID игры
	LevelId     int          //	ID уровня
	LevelNumber int          //	Номер уровня на который был введен ответ
	LevelAction ActionStruct //	инфо о результате отправки ответа на уровень и бонус;
	BonusAction ActionStruct //	инфо о результате отправки ответа на бонус;
}

// Отправка ответов
type ActionStruct struct {
	Answer          string // Введенный ответ
	IsCorrectAnswer bool   // null – ответа не было, false – неправильный ответ; true – правильный ответ;
}

// История введенных ответов
type MixedActionsStruct struct {
	ActionId      int         //
	LevelId       int         //	ID уровня к которому был введен ответ
	LevelNumber   int         //	Номер уровня к которому был введен ответ
	UserId        int         //	ID игрока который ввел ответ
	Kind          int         //	1 – ответ к уровню, 2 – ответ к бонусу
	Login         string      //	Логин игрока который ввел ответ
	Answer        string      //	Текст ответа
	AnswForm      interface{} //string //	Текст ответа с подсветкой русских букв
	EnterDateTime interface{} // struct {  //	Время ввода ответа UTC+0 //
	//	Value int64 `json:"Value"`
	//	} `json:"EnterDateTime"`
	LocDateTime string //	Локализованное время ввода ответа
	IsCorrect   bool   //	Верен/неверен
}

// Сообщения администратора
type MessagesStruct struct {
	OwnerId      int    //	ID Администратора
	OwnerLogin   string //	Логин администратора
	MessageId    int    //	ID Сообщения
	MessageText  string //	Оригинальный текст сообщения
	WrappedText  string //	Отформатированный текст сообщения с учетом ReplaceNl2Br
	ReplaceNl2Br bool   //	Заменять ли \n на  <BR>
}

// Задание к уровню
type TaskStruct struct {
	ReplaceNlToBr     bool   `json:"ReplaceNlToBr"`     //Заменять ли \n на  <BR>
	TaskText          string `json:"TaskText"`          // Оригинальный текст задания
	TaskTextFormatted string `json:"TaskTextFormatted"` //	Отформатированный текст задания с учетом ReplaceNlToBr
}

// Сектора
type SectorsStruct struct {
	SectorId   int          //	ID сектора
	Order      int          //	№ пп
	Name       string       //	Название зектора
	Answer     AnswerStruct //	Отгаданный ответ
	IsAnswered bool         //	отгадан / не отгадан interface{} //
}

// Ответы
type AnswerStruct struct {
	Answer         string `json:"Answer"`
	AnswerDateTime struct {
		Value float64 `json:"Value"`
	} `json:"AnswerDateTime"`
	Login       string `json:"Login"`
	UserID      int    `json:"UserId"`
	LocDateTime string `json:"LocDateTime"`
}

// Информация о подсказках
type HelpsStruct struct {
	HelpId           int    //	ID подсказки
	Number           int    //	Номер пп
	HelpText         string //	Текст подсказки
	IsPenalty        bool   //	Штрафная/Обычная
	Penalty          int    //	Штраф в секундах для штрафной //
	PenaltyComment   string //	Описание подсказки для штрафной //
	RequestConfirm   bool   //	Требует дополнительного подтверждения для штрафной //
	PenaltyHelpState int    //	Состояние, 0 - не открыта; 2 -  открыта для штрафной //
	RemainSeconds    int    //	осталось секуд до того как подскзка будет доступна для игрока
}

// Бонусы
type BonusesStruct struct {
	BonusId        int    //	ID Бонуса
	Name           string //	Название
	Number         int    //	Номер пп
	Task           string //	Задание
	Help           string //	Бонусная подсказка
	IsAnswered     bool   //	Разгадан / не разгадан
	Expired        bool   //	Время на выполнение истекло
	SecondsToStart int    //	Будет доступен через
	SecondsLeft    int    //	Будет еще доступен
	AwardTime      int    //	Начисленный бонус в секндах
}

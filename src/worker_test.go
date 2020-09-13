package src

import (
	"errors"
	"net/http"
	"net/http/cookiejar"
	"testing"
)

func TestWorkerStartGame(t *testing.T) {
	webToBotTEST := make(chan MessengerStyle, 10)
	botToWebTEST := make(chan MessengerStyle, 10)
	var msgChanel MessengerStyle

	cookieJar, _ := cookiejar.New(nil)
	clientTEST := &http.Client{
		Jar: cookieJar,
	}

	var confGameENJSON ConfigGameJSON
	confGameENJSON.InitTest(pathTestConf)

	exitGameTest(clientTEST, confGameENJSON.SubUrl)
	EnterGame(clientTEST, confGameENJSON)
	confGameENJSON.LevelNumber = GameEngineModel(clientTEST, confGameENJSON).Level.Number

	isWork := true
	rightAnswer := "Игра уже идёт!"
	err := StartGame(clientTEST, &confGameENJSON, &isWork, botToWebTEST, webToBotTEST)
	if err != nil {
		t.Errorf("For start game expected 'nil' got %s", err)
	}

	select {
	// В канал msgChanel будут приходить все новые сообщения from web
	case msgChanel = <-webToBotTEST:
		// конец
		if msgChanel.ChannelMessage == "" {
			return
		}
		if msgChanel.ChannelMessage != rightAnswer {
			t.Errorf("For %s\nexpected %s\ngot %s", confGameENJSON.URLGame, rightAnswer, msgChanel.ChannelMessage)
		}
	default:
	}

	msgChanel.ChannelMessage = "stop"
	botToWebTEST <- msgChanel
	isWork = true
	rightAnswer = "Бот выключен. Мы даже не играли &#128546; \nДля перезапуска используйте /restart"
	err = StartGame(clientTEST, &confGameENJSON, &isWork, botToWebTEST, webToBotTEST)
	select {
	// В канал msgChanel будут приходить все новые сообщения from web
	case msgChanel = <-webToBotTEST:
		// конец
		if msgChanel.ChannelMessage == "" {
			return
		}
		if msgChanel.ChannelMessage != rightAnswer {
			t.Errorf("For %s\nexpected %s\ngot %s", confGameENJSON.URLGame, rightAnswer, msgChanel.ChannelMessage)
		}
	default:
	}
	if errors.Is(err, errors.New("Bot stop")) {
		t.Errorf("For stop game expected 'Bot stop' got %s", err)
	}
}

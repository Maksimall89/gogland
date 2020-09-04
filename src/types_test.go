package src

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"os"
	"testing"
)

const pathTestConf = "config_test.json"

func TestTypesConfigGameInit(t *testing.T) {
	var confJSON ConfigGameJSON
	type testPair struct {
		input  string
		output string
	}
	var tests = []testPair{
		{"", "Need more arguments! <code>/start login password http://DEMO.en.cx/GameDetails.aspx?gid=1</code>"},
		{"login password url ErrorsAgr", "Слишком много аргументов!"},
		{"login password http://demo.en.cx/GameDetails.aspx?gid=1", ""},
	}
	for _, pair := range tests {
		result := confJSON.Init(pair.input)
		if result != pair.output {
			t.Errorf("For %s\nexpected %s\ngot %s", pair.input, pair.output, result)
		}
	}
}
func TestTypesSeparateURL(t *testing.T) {
	t.Parallel()
	type testPair struct {
		input  ConfigGameJSON
		expect ConfigGameJSON
	}
	var tests = []testPair{
		{ConfigGameJSON{"", "", "http://demo.en.cx/GameDetails.aspx?gid=1", "demo.en.cx", "1", 0, "", ""}, ConfigGameJSON{"", "", "", "demo.en.cx", "1", 0, "", ""}},
		{ConfigGameJSON{"", "", "http://demo.en.cx/GameDetails.aspx?gid0", "demo.en.cx", "", 0, "", ""}, ConfigGameJSON{"", "", "", "demo.en.cx", "", 0, "", ""}},
		{ConfigGameJSON{"", "", "demo.en.cx/GameDetails.aspx?gid=1", "d", "1", 0, "", ""}, ConfigGameJSON{"", "", "", "", "1", 0, "", ""}},
		{ConfigGameJSON{"", "", "demo.en.cxGameDetails.aspx?gid0", "d", "1", 0, "", ""}, ConfigGameJSON{"", "", "", "", "1", 0, "", ""}},
	}
	for _, pair := range tests {
		pair.input.SeparateURL()
		if (pair.input.SubUrl != pair.expect.SubUrl) && (pair.input.Gid != pair.expect.Gid) {
			t.Errorf("For input expected subUrl %s gid %s\ngot subUrl %s gid %s", pair.input.SubUrl, pair.input.Gid, pair.expect.SubUrl, pair.expect.Gid)
		}
	}
}
func TestTypesConfigBotInit(t *testing.T) {
	type testPair struct {
		input  string
		output ConfigBot
	}
	var tests = []testPair{
		{pathTestConf, ConfigBot{"tokenTEST", "nickOwnTEST", "userTEST", "passTEST", "http://demo.en.cx/GameDetails.aspx?gid=1", 0, []string{"Какой капитан, такая и команда!"}}},
		{"", ConfigBot{}},
	}

	for _, pair := range tests {
		os.Clearenv()
		var Configuration ConfigBot
		Configuration.Init(pair.input)
		if !cmp.Equal(pair.output, Configuration) {
			t.Errorf("For %s\nexpected %v\ngot %v", pair.input, pair.output, Configuration)
		}
	}
}
func TestTypesConfigTestInit(t *testing.T) {
	type testPair struct {
		input  string
		output ConfigGameJSON
	}
	var tests = []testPair{
		{pathTestConf, ConfigGameJSON{"userTEST", "passTEST", "http://demo.en.cx/GameDetails.aspx?gid=1", "demo.en.cx", "1", 0, "", ""}},
		{"", ConfigGameJSON{}},
	}

	for _, pair := range tests {
		os.Clearenv()
		var Configuration ConfigGameJSON
		Configuration.InitTest(pair.input)
		fmt.Println(Configuration)
		if !cmp.Equal(pair.output, Configuration) {
			t.Errorf("For %s\nexpected %v\ngot %v", pair.input, pair.output, Configuration)
		}
	}
}
func TestTypesSetEnv(t *testing.T) {
	var tests = []ConfigBot{
		{"tokenTEST", "nickOwnTEST", "userTEST", "passTEST", "http://demo.en.cx/GameDetails.aspx?gid=1", 0, []string{"Какой капитан, такая и команда!"}},
		{"", "", "", "", "", 0, []string{}},
	}

	for _, pair := range tests {
		var configuration ConfigBot
		_ = os.Setenv("TelegramBotToken", "TestToken")
		_ = os.Setenv("OwnName", "TestOwnName")
		configuration.SetEnv()
		os.Clearenv()
		if configuration.TelegramBotToken != "TestToken" || configuration.OwnName != "TestOwnName" {
			t.Errorf("For %v\nexpected TestToken and TestOwnName\ngot %v", pair, configuration)
		}
	}
}

// config.json
func (conf *ConfigGameJSON) InitTest(path string) {
	var configuration ConfigBot
	configuration.Init(path)

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
	conf.SeparateURL()
}

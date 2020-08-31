package main

import (
	"github.com/google/go-cmp/cmp"
	"os"
	"testing"
)

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
		result := confJSON.init(pair.input)
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
		pair.input.separateURL()
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
		{"config_test.json", ConfigBot{"tokenTEST", "nickOwnTEST", "userTEST", "passTEST", "http://demo.en.cx/GameDetails.aspx?gid=1", 0, []string{"Какой капитан, такая и команда!"}}},
		{"", ConfigBot{}},
	}

	for _, pair := range tests {
		os.Clearenv()
		var configuration ConfigBot
		configuration.init(pair.input)
		if !cmp.Equal(pair.output, configuration) {
			t.Errorf("For %s\nexpected %v\ngot %v", pair.input, pair.output, configuration)
		}
	}
}
func TestTypesConfigTestInit(t *testing.T) {
	type testPair struct {
		input  string
		output ConfigGameJSON
	}
	var tests = []testPair{
		{"config_test.json", ConfigGameJSON{"userTEST", "passTEST", "http://demo.en.cx/GameDetails.aspx?gid=1", "demo.en.cx", "1", 0, "", ""}},
		{"", ConfigGameJSON{}},
	}

	for _, pair := range tests {
		os.Clearenv()
		var configuration ConfigGameJSON
		configuration.initTest(pair.input)
		if !cmp.Equal(pair.output, configuration) {
			t.Errorf("For %s\nexpected %v\ngot %v", pair.input, pair.output, configuration)
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
		configuration.setEnv()
		os.Clearenv()
		if configuration.TelegramBotToken != "TestToken" || configuration.OwnName != "TestOwnName" {
			t.Errorf("For %v\nexpected TestToken and TestOwnName\ngot %v", pair, configuration)
		}
	}
}

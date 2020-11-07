package help

import (
	"strings"
	"testing"
)

func TestHelpSearchAnagramAndMaskWord(t *testing.T) {
	type testPair struct {
		input      string
		output     string
		typeSearch bool
	}

	var tests = []testPair{
		{"плео", " ПОЛЕ ", true},
		{"фывфывйуфывфы", strNotWords, true},
		{"электростанция", "ЭЛЕКТРОСТАНЦИЯ ", true},
		{"по?е", "ПОБЕ ПОЛЕ ", false},
		{"фывфывйуфывфы", strNotWords, false},
		{"пол*", "ПОЛ ПОЛА ПОЛАБЗИТЬСЯ ПОЛАБЫ ПОЛАВИРУЕМ ПОЛАВОЧНИК ПОЛАВОЧНЫЙ ПОЛАГАНИЕ", false},
		{"электростанция", "ЭЛЕКТРОСТАНЦИЯ ", false},
	}

	for _, pair := range tests {
		result := SearchAnagramAndMaskWord(pair.input, pair.typeSearch)
		if !strings.ContainsAny(result, pair.output) {
			t.Errorf("For %s\nexpected %s\ngot %s", pair.input, pair.output, result)
		}
	}
}
func TestHelpAssociations(t *testing.T) {
	t.Parallel()
	type testPair struct {
		input  string
		output string
	}

	var tests = []testPair{
		{"поле", "пшеница луг простор трава"},
		{"фывфывйуфывфы", strNotWords},
		{"электростанция", "аэс тэс провода энергия"},
	}

	for _, pair := range tests {
		result := Associations(pair.input)
		if !strings.Contains(result, pair.output) {
			t.Errorf("For %s\nexpected %s\ngot %s", pair.input, pair.output, result)
		}
	}
}
func TestHelpTransferToAlphabetToWord(t *testing.T) {
	t.Parallel()

	type testPair struct {
		input  string
		output string
		bool
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
		result := TransferToAlphabet(pair.input, true)
		if result != pair.output {
			t.Errorf("For %s\nexpected %s\ngot %s", pair.input, pair.output, result)
		}
	}
	for _, pair := range testsW2N {
		result := TransferToAlphabet(pair.input, false)
		if result != pair.output {
			t.Errorf("For %s\nexpected %s\ngot %s", pair.input, pair.output, result)
		}
	}
}
func TestHelpMendeleevTable(t *testing.T) {
	t.Parallel()

	type testPair struct {
		input  string
		output string
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
		result := TableMendeleev(pair.input)
		if result != pair.output {
			t.Errorf("For %s\nexpected %s\ngot %s", pair.input, pair.output, result)
		}
	}
}
func TestHelpMorse(t *testing.T) {
	t.Parallel()

	type testPair struct {
		input  string
		output string
	}

	var tests = []testPair{
		{".-.-", ".-.- = <b>RU</b> Я <b>EN</b> Ä;\n"},
		{"sf", "Символа sf азбуке морзе нет.\n"},
		{".. ---", ".. = <b>RU</b> И <b>EN</b> I;\n--- = <b>RU</b> О <b>EN</b> O;\n"},
	}

	for _, pair := range tests {
		result := Morse(pair.input)
		if result != pair.output {
			t.Errorf("For %s\nexpected %s\ngot %s", pair.input, pair.output, result)
		}
	}
}
func TestHelpBrail(t *testing.T) {
	t.Parallel()

	type testPair struct {
		input  string
		output string
	}

	var tests = []testPair{
		{"100000", "100000 = <b>RU</b> А <b>EN</b> A <b>№</b> 1;\n"},
		{"sf", "Символа sf шрифте Браиля нет.\n"},
		{"100000 101000", "100000 = <b>RU</b> А <b>EN</b> A <b>№</b> 1;\n101000 = <b>RU</b> Б <b>EN</b> B <b>№</b> 2;\n"},
	}

	for _, pair := range tests {
		result := Braille(pair.input)
		if result != pair.output {
			t.Errorf("For %s\nexpected %s\ngot %s", pair.input, pair.output, result)
		}
	}
}
func TestHelpBin(t *testing.T) {
	t.Parallel()

	type testPair struct {
		input     string
		output    string
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
		result := Bin(pair.input, pair.attribute)
		if result != pair.output {
			t.Errorf("For %s\nexpected %s\ngot %s", pair.input, pair.output, result)
		}
	}
}
func TestHelpAutoCode(t *testing.T) {
	t.Parallel()

	type testPair struct {
		input  string
		output string
	}

	var tests = []testPair{
		{"18", "<b>18</b> регион = Удмуртская Республика;\n"},
		{"18 45", "<b>18</b> регион = Удмуртская Республика;\n<b>45</b> регион = Курганская область;\n"},
		{"электростанция", "Региона электростанция в России нет."},
	}

	for _, pair := range tests {
		result := AutoCode(pair.input)
		if result != pair.output {
			t.Errorf("For %s\nexpected %s\ngot %s", pair.input, pair.output, result)
		}
	}
}
func TestHelpTranslateQwerty(t *testing.T) {
	t.Parallel()

	type testPair struct {
		input  string
		output string
	}

	var tests = []testPair{
		{"a B c", "a = <b>RU</b> ф <b>EN</b> a;\nb = <b>RU</b> и <b>EN</b> b;\nc = <b>RU</b> с <b>EN</b> c;\n"},
		{"Ц d Ц", "ц = <b>RU</b> ц <b>EN</b> w;\nd = <b>RU</b> в <b>EN</b> d;\nц = <b>RU</b> ц <b>EN</b> w;\n"},
		{"18", "Символа 18 на клавиатуре нет.\n"},
		{"э", "э = <b>RU</b> э <b>EN</b> ';\n"},
		{"z", "z = <b>RU</b> я <b>EN</b> z;\n"},
	}

	for _, pair := range tests {
		result := TranslateQwerty(pair.input)
		if result != pair.output {
			t.Errorf("For %s\nexpected %s\ngot %s", pair.input, pair.output, result)
		}
	}
}

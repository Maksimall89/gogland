package main

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
		{"плео", "ЕЛОП ЛОПЕ ПОЛЕ ", true},
		{"фывфывйуфывфы", "<b>Слов не обнаружено.</b>", true},
		{"электростанция", "ЭЛЕКТРОСТАНЦИЯ ", true},
		{"по?е", "ПОБЕ ПОЛЕ ", false},
		{"фывфывйуфывфы", "<b>Слов не обнаружено.</b>", false},
		{"пол*", "ПОЛ ПОЛА ПОЛАБЗИТЬСЯ ПОЛАБЫ ПОЛАВИРУЕМ ПОЛАВОЧНИК ПОЛАВОЧНЫЙ ПОЛАГАНИЕ ПОЛАД ПОЛАДКА ПОЛАЗ ПОЛАЗНА ПОЛАЗНИК ПОЛАЙКА ПОЛАК ПОЛАКАНТ ПОЛАКОМИЛ ПОЛАКСИС ПОЛАН ПОЛАНДРА ПОЛАНИ ПОЛАНСКИ ПОЛАНСКИЙ ПОЛАНЬИ ПОЛАР ПОЛАРОИД ПОЛАТИ ПОЛАТКА ПОЛАШ ПОЛБА ПОЛБИГ ПОЛБИН ПОЛБОТИНКА ПОЛБУТЫЛКА ПОЛБУТЫЛКИ ПОЛВА ПОЛВЕДРА ПОЛВЕКА ПОЛВЕРШКА ПОЛВЕТРА ПОЛГАР ПОЛГАРНЕЦ ПОЛГОДА ПОЛГОРЯ ПОЛДЕЛА ПОЛДЕНЬ ПОЛДНИК ПОЛДНИЧАНИЕ ПОЛДНИЧАНЬЕ ПОЛДНЯ ПОЛДОРОГИ ПОЛДЮЖИНЫ ПОЛДЮЙМА ПОЛЕ ПОЛЕВАНОВ ПОЛЕВАНЬЕ ПОЛЕВЕНИЕ ПОЛЕВИК ПОЛЕВИЦА ПОЛЕВИЧКА ПОЛЕВКА ПОЛЕВОД ПОЛЕВОДСТВО ПОЛЕВОЙ ПОЛЕВСКОЙ ПОЛЕВЩИК ПОЛЕГАДА ПОЛЕГАЕМОСТЬ ПОЛЕГАНИЕ ПОЛЕЖАЕВ ПОЛЕЖАЕВСКАЯ ПОЛЕЖАЕВСКИЙ ПОЛЕЖАЛОЕ ПОЛЕЗНО ПОЛЕЗНОЕ ПОЛЕЗНОСТЬ ПОЛЕЗНЫЙ ПОЛЕЛЕЯЛА ПОЛЕЛЬ ", false},
		{"электростанция", "ЭЛЕКТРОСТАНЦИЯ ", false},
	}

	for _, pair := range tests {
		result := searchAnagramAndMaskWord(pair.input, pair.typeSearch)
		if result != pair.output {
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
		{"поле", "пшеница луг простор трава футбол лес рожь воин кукуруза пашня трактор степь цветы сено комбайн конь колос ромашки гольф ветер мак огород пахарь битва нива васильки ромашка раздолье поляна деревня магнит хлеб корова зерно овёс газон василёк коровы сенокос природа ковыль сеялка свобода колхоз подсолнух земля стог брань пахота урожай игра пастух лужайка война лето клевер работа бабочки россия равнина лён дерево полевод перекати-поле один напряжённость танки плуг незабудки колосья коса пастбище комбайнер травы картошка жатва пшено просторы кони русь село даль ширь чудеса бой колосок дорога жизнь солома ягода посев русское поле лошадь небо трактора цветок пространство песня борозда путь"},
		{"фывфывйуфывфы", "<b>Слов не обнаружено.</b>"},
		{"электростанция", "аэс тэс провода энергия гоэлро лэп электрификация коммунизм плотина свет генератор ток топливо гэс тэц молния тесла высоковольтная линия напряжение двигатель ротор пуск атом курчатов ленин турбина лампочка течение пудель электричество обнинск ветер солнечные батареи стройка прилив электротехника атомы импульс фотосинтез афганистан балаково"},
	}

	for _, pair := range tests {
		result := associations(pair.input)
		if strings.TrimSpace(result) != pair.output {
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
		result := transferToAlphabet(pair.input, true)
		if result != pair.output {
			t.Errorf("For %s\nexpected %s\ngot %s", pair.input, pair.output, result)
		}
	}
	for _, pair := range testsW2N {
		result := transferToAlphabet(pair.input, false)
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
		result := tableMendeleev(pair.input)
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
		result := morse(pair.input)
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
		result := braille(pair.input)
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
		result := bin(pair.input, pair.attribute)
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
		result := autoCode(pair.input)
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
		result := translateQwerty(pair.input)
		if result != pair.output {
			t.Errorf("For %s\nexpected %s\ngot %s", pair.input, pair.output, result)
		}
	}
}

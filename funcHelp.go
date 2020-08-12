package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
)

// other functional
// type for https://anagram.poncy.ru
type ObjPoncy struct {
	Index    bool   `json:"index"`
	MainPage bool   `json:"main_page"`
	H1       string `json:"h1"`
	Title    string `json:"title"`
	H1s      string `json:"h1s"`
	DescMain string `json:"desc_main"`
	Desc     string `json:"desc"`
}
type AnswerPoncy struct {
	Url             string   `json:"url"`
	Result          []string `json:"result"`
	PageDescription ObjPoncy `json:"page_description"`
}

func anagram(text string) (string, error) {
	// create cookie
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}

	resp, err := client.Get("https://anagram.poncy.ru/anagram-decoding.cgi?name=anagram_index&inword=" + text + "&answer_type=1")
	if err != nil {
		log.Println(err)
		return "Ошибка отправки запроса.", errors.New("ошибка отправки запроса")
	}
	// read from body
	body, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		log.Println(string(body))
		log.Println(err)
		return "Не могу распарсить.", errors.New("не могу распарсить")
	}

	var answer AnswerPoncy // срез байт входных, куда кладём
	err = json.Unmarshal(body, &answer)
	if err != nil {
		log.Println(err)
		return "Не могу распарсить JSON.", errors.New("не могу распарсить JSON")
	}

	var str string
	if len(answer.Result) > 0 {
		for _, item := range answer.Result {
			if len(str) > 1300 {
				break
			}
			str += item + " "
		}
	} else {
		str = "<b>Слов не обнаружено.</b>"
	}
	return str, nil
}
func searchForMask(text string) (string, error) {
	// create cookie
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}

	// replace input text for the site
	text = strings.Replace(text, "*", "%", 2)
	text = strings.Replace(text, "?", "*", 2)

	resp, err := client.Get("https://anagram.poncy.ru/anagram-decoding.cgi?name=words_by_mask_index&inword=" + text + "&answer_type=4")
	if err != nil {
		log.Println(err)
		return "Ошибка отправки запроса.", errors.New("ошибка отправки запроса")
	}

	// read from body
	body, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		log.Println(string(body))
		log.Println(err)
		return "Не могу распарсить.", errors.New("не могу распарсить")
	}

	var answer AnswerPoncy
	// срез байт входных, куда кладём
	err = json.Unmarshal(body, &answer)
	if err != nil {
		log.Println(err)
		return "Не могу распарсить JSON.", errors.New("не могу распарсить JSON")
	}

	var str string
	if len(answer.Result) > 0 {
		for _, item := range answer.Result {
			if len(str) > 1300 {
				break
			}
			str += item + " "
		}
	} else {
		str = "<b>Слов не обнаружено.</b>"
	}

	return str, nil
}
func associations(text string) (string, error) {

	type ObjSocialistic struct {
		Name              string  `json:"name"`
		PopularityInverse int     `json:"popularity_inverse"`
		AssociationsCount int     `json:"associations_count"`
		PopularityDirect  int     `json:"popularity_direct"`
		Popularity        float64 `json:"popularity"`
		WordPopularity    float64 `json:"word_popularity"`
		Weight            float64 `json:"weight"`
		Positivity        float64 `json:"positivity"`
	}
	type AnswerSocialistic struct {
		Associations []struct{ ObjSocialistic } `json:"associations"`
		Word         string                     `json:"word"`
	}

	// create cookie
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}
	resp, err := client.PostForm("http://sociation.org/ajax/word_associations/", url.Values{"max_count": {"0"}, "back": {"false"}, "word": {text}})
	if err != nil {
		log.Println(err)
		return "Ошибка отправки запроса.", errors.New("ошибка отправки запроса")
	}
	// read from body
	body, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		log.Println(err)
		log.Println(string(body))
		return "Не могу распарсить.", errors.New("не могу распарсить")
	}

	var answer AnswerSocialistic
	// срез байт входных, куда кладём
	err = json.Unmarshal(body, &answer)
	if err != nil {
		log.Println(err)
		return "Не могу распарсить JSON.", errors.New("не могу распарсить JSON")
	}

	var str string
	if len(answer.Associations) > 0 {
		for _, item := range answer.Associations {
			if len(str) > 1300 {
				break
			}
			str += item.Name + " "
		}
	} else {
		str = "<b>Слов не обнаружено.</b>"
	}

	return str, nil
}
func tableMendeleev(text string) string {

	var attributeSearch bool
	type table struct {
		name      string
		shortName string
		mass      string
	}

	var elements = []table{
		{"-", "-", "-"},
		{"H", "Водород", "1,00794"},
		{"He", "Гелий", "4,002602"},
		{"Li", "Литий", "6,941"},
		{"Be", "Бериллий", "9,01218"},
		{"B", "Бор", "10,81"},
		{"C", "Углерод", "12,011"},
		{"N", "Азот", "14,0067"},
		{"O", "Кислород", "15,9994"},
		{"F", "Фтор", "18,9984"},
		{"Ne", "Неон", "20,179"},
		{"Na", "Натрий", "22,98977"},
		{"Mg", "Магний", "24,305"},
		{"Al", "Алюминий", "26,98154"},
		{"Si", "Кремний", "28,0855"},
		{"P", "Фосфор", "30,97376"},
		{"S", "Сера", "32,06"},
		{"Cl", "Хлор", "35,453"},
		{"Ar", "Аргон", "39,948"},
		{"K", "Калий", "39,0983"},
		{"Ca", "Кальций", "40,08"},
		{"Sc", "Скандий", "44,9559"},
		{"Ti", "Титан", "47,88"},
		{"V", "Ванадий", "50,9415"},
		{"Cr", "Хром", "51,996"},
		{"Mn", "Марганец", "54,938"},
		{"Fe", "Железо", "55,847"},
		{"Co", "Кобальт", "58,9332"},
		{"Ni", "Никель", "58,69"},
		{"Cu", "Медь", "63,546"},
		{"Zn", "Цинк", "65,38"},
		{"Ga", "Галлий", "69,72"},
		{"Ge", "Германий", "72,59"},
		{"As", "Мышьяк", "74,9216"},
		{"Se", "Селен", "78,96"},
		{"Br", "Бром", "79,904"},
		{"Kr", "Криптон", "83,8"},
		{"Rb", "Рубидий", "85,4678"},
		{"Sr", "Стронций", "87,62"},
		{"Y", "Иттрий", "88,9059"},
		{"Zr", "Цирконий", "91,22"},
		{"Nb", "Ниобий", "92,9064"},
		{"Mo", "Молибден", "95,94"},
		{"Tc", "Технеций", "98"},
		{"Ru", "Рутений", "101,07"},
		{"Rh", "Родий", "102,9055"},
		{"Pd", "Палладий", "106,42"},
		{"Ag", "Серебро", "107,868"},
		{"Cd", "Кадмий", "112,41"},
		{"In", "Индий", "114,82"},
		{"Sn", "Олово", "118,69"},
		{"Sb", "Сурьма", "121,75"},
		{"Te", "Теллур", "127,6"},
		{"I", "Йод", "126,9045"},
		{"Xe", "Ксенон", "131,29"},
		{"Cs", "Цезий", "132,9054"},
		{"Ba", "Барий", "137,33"},
		{"La", "Лантан", "138,9055"},
		{"Ce", "Церий", "140,116"},
		{"Pr", "Празеодим", "140,9077"},
		{"Nd", "Неодим", "144,24"},
		{"Pm", "Прометий", "145"},
		{"Sm", "Самарий", "150,36"},
		{"Eu", "Европий", "151,96"},
		{"Gd", "Гадолиний", "157,25"},
		{"Tb", "Тербий", "158,9254"},
		{"Dy", "Диспрозий", "162,5"},
		{"Ho", "Гольмий", "164,9304"},
		{"Er", "Эрбий", "167,26"},
		{"Tm", "Тулий", "168,9342"},
		{"Yb", "Иттербий", "173,04"},
		{"Lu", "Лютеций", "174,967"},
		{"Hf", "Гафний", "178,49"},
		{"Ta", "Тантал", "180,9479"},
		{"W", "Вольфрам", "183,85"},
		{"Re", "Рений", "186,207"},
		{"Os", "Осьмий", "190,2"},
		{"Ir", "Иридий", "192,22"},
		{"Pt", "Платина", "195,08"},
		{"Au", "Золото", "196,9665"},
		{"Hg", "Ртуть", "200,59"},
		{"Tl", "Таллий", "204,383"},
		{"Pb", "Свинец", "207,2"},
		{"Bi", "Висмут", "208,9804"},
		{"Po", "Полоний", "209"},
		{"At", "Астат", "210"},
		{"Rn", "Радон", "222"},
		{"Fr", "Франций", "223"},
		{"Ra", "Радий", "226,0254"},
		{"Ac", "Актиний", "227,0278"},
		{"Th", "Торий", "232,0381"},
		{"Pa", "Протактиний", "231,0359"},
		{"U", "Уран", "238,0389"},
		{"Np", "Нептуний", "237,0482"},
		{"Pu", "Плутоний", "244"},
		{"Am", "Америций", "243"},
		{"Cm", "Кюрий", "247"},
		{"Bk", "Берклий", "247"},
		{"Cf", "Калифорний", "251"},
		{"Es", "Эйнштейний", "252"},
		{"Fm", "Фермий", "257"},
		{"Md", "Менделевий", "258"},
		{"No", "Нобелий", "255"},
		{"Lr", "Лоуренсий", "260,1"}, // 103
		{"Rf", "Резерфордий", "261"},
		{"Db", "Дубний", "262"},
		{"Sg", "Сиборгий", "266"},
		{"Bh", "Борий", "267"},
		{"Hs", "Хассий", "269"},
		{"Mt", "Мейтнерий", "276"},
		{"Ds", "Дармштадтий", "227"},
		{"Rg", "Рентгений", "280"},
		{"Cn", "Коперниций", "285"},
		{"Nh", "Нихоний", "284"},
		{"Fl", "Флеровий", "289"},
		{"MC", "Московий", "288"},
		{"Lv", "Ливерморий", "293"},
		{"Ts", "Теннессин", "294"},
		{"Oq", "Оганесон", "294"},
		{"Uue", "Унуненний", "316"},
		{"Ubn", "Унбинилий", "320"},
		{"Ubu", "Унбиуний", "320"},
		{"Ubb", "Унбибий", "-"},
		{"Ubt", "Унбитрий", "-"},
		{"Ubq", "Унбиквадий", "-"},
		{"Ubp", "Унбипентий", "322"},
		{"Ubh", "Унбигексий", "322"},
		{"Ubs", "Унбисептий", "-"},
	}

	arrText := strings.Split(text, " ")
	text = ""
	for _, item := range arrText {
		attributeSearch = true
		if item == "" {
			continue
		}
		number, err := strconv.Atoi(item)
		if err != nil {
			for _, element := range elements {
				if strings.EqualFold(element.name, item) {
					text += "Номер " + item + " = " + element.name + " - " + element.shortName + " - масса " + element.mass + ";\n"
					attributeSearch = false
					break
				}
			}
			if attributeSearch {
				text += item + " = не допустимое число!\n"
			}
			continue
		}
		if (number > 0) && (number < 128) {
			text += "Номер " + item + " = " + elements[number].name + " - " + elements[number].shortName + " - масса " + elements[number].mass + ";\n"
		} else {
			text += "Номера " + item + " не существует в таблице.\n"
		}
	}

	return text
}
func morse(text string) string {
	type table struct {
		name      string
		rusSymbol string
		engSymbol string
	}

	var symbols = []table{
		{".-", "А", "A"},
		{"-...", "Б", "B"},
		{".--", "В", "W"},
		{"--.", "Г", "G"},
		{"-..", "Д", "D"},
		{".", "Е(Ё)", "E"},
		{"...-", "Ж", "V"},
		{"..--", "З", "Z"},
		{"..", "И", "I"},
		{".---", "Й", "J"},
		{"-.-", "К", "K"},
		{".-..", "Л", "L"},
		{"--", "М", "M"},
		{"-.", "Н", "N"},
		{"---", "О", "O"},
		{".--.", "П", "P"},
		{".-.", "Р", "R"},
		{"...", "С", "S"},
		{"-", "Т", "T"},
		{"..-", "У", "U"},
		{"..-.", "Ф", "F"},
		{"....", "Х", "H"},
		{"-.-.", "Ц", "C"},
		{"---.", "Ч", "Ö"},
		{"----", "Ш", "CH"},
		{"--.-", "Щ", "Q"},
		{"--.--", "Ъ", "Ñ"},
		{"-.--", "Ы", "Y"},
		{"-..-", "Ь(Ъ)", "X"},
		{"..-..", "Э", "É"},
		{"..--", "Ю", "Ü"},
		{".-.-", "Я", "Ä"},
		{".----", "1", "1"},
		{"..---", "2", "2"},
		{"...--", "3", "3"},
		{"....-", "4", "4"},
		{".....", "5", "5"},
		{"-....", "6", "6"},
		{"--...", "7", "7"},
		{"---..", "8", "8"},
		{"----.", "9", "9"},
		{"-----", "0", "0"},
		{"......", "точка", "точка"},
		{".-.-.-", "запятая", "запятая"},
		{"---...", "двоеточие", "двоеточие"},
		{"-.-.-.", "точка с запятой", "точка с запятой"},
		{"-.--.-", "скобка", "скобка"},
		{".----.", "апостроф", "апостроф"},
		{".-..-.", "кавычки", "кавычки"},
		{"-....-", "тире", "тире"},
		{"-..-.", "косая черта", "косая черта"},
		{"..--..", "вопросительный знак", "вопросительный знак"},
		{"--..--", "восклицательный знак", "восклицательный знак"},
		{"-...-", "знак раздела", "знак раздела"},
		{"........", "ошибка/перебой", "ошибка/перебой"},
		{".--.-.", "собака(@)", "собака(@)"},
		{"..-.-", "конец связи", "конец связи"},
	}

	arrText := strings.Split(text, " ")
	text = ""
	for _, item := range arrText {
		if item == "" {
			continue
		}

		for _, symbol := range symbols {
			if symbol.name == item {
				text += item + " = <b>RU</b> " + symbol.rusSymbol + " <b>EN</b> " + symbol.engSymbol + ";\n"
			}
		}
		if text == "" {
			text += "Символа " + item + " азбуке морзе нет.\n"
		}
	}

	return text
}
func autoCode(text string) string {

	type table struct {
		code   string
		region string
	}

	var regions = []table{
		{"1", "Республика Адыгея"},
		{"2", "Республика Башкортостан"},
		{"102", "Республика Башкортостан"},
		{"3", "Республика Бурятия"},
		{"103", "Республика Бурятия"},
		{"4", "Республика Алтай (Горный Алтай)"},
		{"5", "Республика Дагестан"},
		{"6", "Республика Ингушетия"},
		{"7", "Кабардино-Балкарская Республика"},
		{"8", "Республика Калмыкия"},
		{"9", "Республика Карачаево-Черкессия"},
		{"10", "Республика Карелия"},
		{"11", "Республика Коми"},
		{"12", "Республика Марий Эл"},
		{"13", "Республика Мордовия"},
		{"113", "Республика Мордовия"},
		{"14", "Республика Саха (Якутия)"},
		{"15", "Республика Северная Осетия — Алания"},
		{"16", "Республика Татарстан"},
		{"116", "Республика Татарстан"},
		{"17", "Республика Тыва"},
		{"18", "Удмуртская Республика"},
		{"19", "Республика Хакасия"},
		{"21", "Чувашская Республика"},
		{"121", "Чувашская Республика"},
		{"22", "Алтайский край"},
		{"23", "Краснодарский край"},
		{"93", "Краснодарский край"},
		{"123", "Краснодарский край"},
		{"24", "Красноярский край"},
		{"84", "Красноярский край"},
		{"88", "Красноярский край"},
		{"124", "Красноярский край"},
		{"25", "Приморский край"},
		{"125", "Приморский край"},
		{"26", "Ставропольский край"},
		{"126", "Ставропольский край"},
		{"27", "Хабаровский край"},
		{"28", "Амурская область"},
		{"29", "Архангельская область"},
		{"30", "Астраханская область"},
		{"31", "Белгородская область"},
		{"32", "Брянская область"},
		{"33", "Владимирская область"},
		{"34", "Волгоградская область"},
		{"134", "Волгоградская область"},
		{"35", "Вологодская область, "},
		{"36", "Воронежская область"},
		{"136", "Воронежская область"},
		{"37", "Ивановская область"},
		{"38", "Иркутская область"},
		{"85", "Иркутская область"},
		{"138", "Иркутская область"},
		{"39", "Калининградская область"},
		{"91", "Калининградская область"},
		{"40", "Калужская область"},
		{"41", "Камчатский край"},
		{"42", "Кемеровская область"},
		{"142", "Кемеровская область"},
		{"43", "Кировская область"},
		{"44", "Костромская область"},
		{"45", "Курганская область"},
		{"46", "Курская область"},
		{"47", "Ленинградская область"},
		{"48", "Липецкая область"},
		{"49", "Магаданская область"},
		{"50", "Московская область"},
		{"90", "Московская область"},
		{"150", "Московская область"},
		{"190", "Московская область"},
		{"750", "Московская область"},
		{"51", "Мурманская область"},
		{"52", "Нижегородская область"},
		{"152", "Нижегородская область"},
		{"53", "Новгородская область"},
		{"54", "Новосибирская область"},
		{"154", "Новосибирская область"},
		{"55", "Омская область"},
		{"56", "Оренбургская область"},
		{"57", "Орловская область"},
		{"58", "Пензенская область"},
		{"59", "Пермский край"},
		{"81", "Пермский край"},
		{"159", "Пермский край"},
		{"60", "Псковская область"},
		{"61", "Ростовская область"},
		{"161", "Ростовская область"},
		{"62", "Рязанская область"},
		{"63", "Самарская область"},
		{"163", "Самарская область"},
		{"64", "Саратовская область"},
		{"164", "Саратовская область"},
		{"65", "Сахалинская область"},
		{"66", "Свердловская область"},
		{"67", "Смоленская область"},
		{"68", "Тамбовская область"},
		{"69", "Тверская область"},
		{"70", "Томская область"},
		{"71", "Тульская область"},
		{"72", "Тюменская область"},
		{"73", "Ульяновская область"},
		{"173", "Ульяновская область"},
		{"74", "Челябинская область"},
		{"174", "Челябинская область"},
		{"8", "Забайкальский край"},
		{"75", "Забайкальский край"},
		{"76", "Ярославская область"},
		{"77", "г. Москва"},
		{"97", "г. Москва"},
		{"99", "г. Москва"},
		{"177", "г. Москва"},
		{"197", "г. Москва"},
		{"199", "г. Москва"},
		{"777", "г. Москва"},
		{"78", "г. Санкт-Петребург"},
		{"98", "г. Санкт-Петребург"},
		{"178", "г. Санкт-Петребург"},
		{"79", "Еврейская автономная область"},
		{"82", "Республика Крым"},
		{"83", "Ненецкий автономный округ"},
		{"86", "Ханты-Мансийский автономный округ — Югра"},
		{"186", "Ханты-Мансийский автономный округ — Югра"},
		{"87", "Чукотский автономный округ"},
		{"89", "Ямало-Ненецкий автономный округ"},
		{"92", "г. Севастополь"},
		{"94", "Территории, находящиеся за пределами РФ и обслуживаемые Департаментом режимных объектов МВД России"},
		{"95", "Чеченская республика"},
	}

	arrText := strings.Split(text, " ")
	text = ""

	for _, item := range arrText {
		if item == "" {
			continue
		}

		for _, region := range regions {
			if region.code == item {
				text += "<b>" + item + "</b> регион = " + region.region + ";\n"
			}
		}
		if text == "" {
			text += "Региона " + item + " в России нет."
		}
	}
	// select all regions
	if text == "" {
		text = "1 - Республика Адыгея\n" +
			"2, 102 - Республика Башкортостан\n" +
			"3, 103 - Республика Бурятия\n" +
			"4 - Республика Алтай (Горный Алтай)\n" +
			"5 - Республика Дагестан\n" +
			"6 - Республика Ингушетия\n" +
			"7 - Кабардино-Балкарская Республика\n" +
			"8 - Республика Калмыкия\n" +
			"9 - Республика Карачаево-Черкессия\n" +
			"10 - Республика Карелия\n" +
			"11 - Республика Коми\n" +
			"12 - Республика Марий Эл\n" +
			"13,113 - Республика Мордовия\n" +
			"14 - Республика Саха (Якутия)\n" +
			"15 - Республика Северная Осетия — Алания\n" +
			"16,116 - Республика Татарстан\n" +
			"17 - Республика Тыва\n" +
			"18 - Удмуртская Республика\n" +
			"19 - Республика Хакасия\n" +
			"21,121 - Чувашская Республика\n" +
			"22 - Алтайский край\n" +
			"23,93,123 - Краснодарский край\n" +
			"24,84,88,124 - Красноярский край\n" +
			"25,125 - Приморский край\n" +
			"26,126 - Ставропольский край\n" +
			"27 - Хабаровский край\n" +
			"28 - Амурская область\n" +
			"29 - Архангельская область\n" +
			"30 - Астраханская область\n" +
			"31 - Белгородская область\n" +
			"32 - Брянская область\n" +
			"33 - Владимирская область\n" +
			"34,134 - Волгоградская область\n" +
			"35 - Вологодская область\n" +
			"36,136 - Воронежская область\n" +
			"37 - Ивановская область\n" +
			"38,85,138 - Иркутская область\n" +
			"39,91 - Калининградская область\n" +
			"40 - Калужская область\n" +
			"41 - Камчатский край\n" +
			"42,142 - Кемеровская область\n" +
			"43 - Кировская область\n" +
			"44 - Костромская область\n" +
			"45 - Курганская область\n" +
			"46 - Курская область\n" +
			"47 - Ленинградская область\n" +
			"48 - Липецкая область\n" +
			"49 - Магаданская область\n" +
			"50,90,150,190,750 - Московская область\n" +
			"51 - Мурманская область\n" +
			"52,152 - Нижегородская область\n" +
			"53 - Новгородская область\n" +
			"54,154 - Новосибирская область\n" +
			"55 - Омская область\n" +
			"56 - Оренбургская область\n" +
			"57 - Орловская область\n" +
			"58 - Пензенская область\n" +
			"59,81,159 - Пермский край\n" +
			"60 - Псковская область\n" +
			"61,161 - Ростовская область\n" +
			"62 - Рязанская область\n" +
			"63,163 - Самарская область\n" +
			"64,164 - Саратовская область\n" +
			"65 - Сахалинская область\n" +
			"66 - Свердловская область\n" +
			"67 - Смоленская область\n" +
			"68 - Тамбовская область\n" +
			"69 - Тверская область\n" +
			"70 - Томская область\n" +
			"71 - Тульская область\n" +
			"72 - Тюменская область\n" +
			"73,173 - Ульяновская область\n" +
			"74,174 - Челябинская область\n" +
			"75,80 - Забайкальский край\n" +
			"76 - Ярославская область\n" +
			"77,97,99,177,197,199,777 - г. Москва\n" +
			"78,98,178 - г. Санкт-Петербург\n" +
			"79 - Еврейская автономная область\n" +
			"82 - Республика Крым\n" +
			"83 - Ненецкий автономный округ\n" +
			"86,186 - Ханты-Мансийский автономный округ — Югра\n" +
			"87 - Чукотский автономный округ\n" +
			"89 - Ямало-Ненецкий автономный округ\n" +
			"92 - г. Севастополь\n" +
			"94 - Территории, находящиеся за пределами РФ и обслуживаемые Департаментом режимных объектов МВД России\n" +
			"95 - Чеченская республика"
	}

	return text

}
func braille(text string) string {
	type table struct {
		code      string
		rusSymbol string
		engSymbol string
		number    string
	}

	var symbols = []table{
		{"100000", "А", "A", "1"},
		{"101000", "Б", "B", "2"},
		{"110000", "Ц", "C", "3"},
		{"110100", "Д", "D", "4"},
		{"100100", "Е", "E", "5"},
		{"111000", "Ф", "F", "6"},
		{"111100", "Г", "G", "7"},
		{"101100", "Х", "H", "8"},
		{"011000", "Л", "I", "9"},
		{"011100", "Ж", "J", "0"},
		{"100010", "К", "K", "-"},
		{"101010", "И", "I", "-"},
		{"110010", "М", "M", "-"},
		{"110110", "Н", "N", "-"},
		{"100110", "О", "O", "-"},
		{"111010", "П", "P", "-"},
		{"111110", "Ч", "Q", "-"},
		{"101110", "Р", "R", "-"},
		{"011010", "С", "S", "-"},
		{"011110", "Т", "T", "-"},
		{"100011", "У", "U", "-"},
		{"101011", "-", "V", "-"},
		{"110011", "Щ", "X", "-"},
		{"110111", "-", "Y", "-"},
		{"100111", "З", "Z", "-"},
		{"011101", "В", "W", "-"},
		{"100100", "Ё", "-", "-"},
		{"111011", "Й", "-", "-"},
		{"100101", "Ш", "-", "-"},
		{"101111", "Ъ", "-", "-"},
		{"011011", "Ы", "-", "-"},
		{"011111", "Ь", "-", "-"},
		{"011001", "Э", "-", "-"},
		{"101101", "Ю", "-", "-"},
		{"111001", "Я", "-", "-"},
		{"000001", "-", "-", "Следующая буква заглавная"},
		{"001100", "-", "-", "."},
		{"001000", "-", "-", ","},
		{"000010", "-", "-", "'"},
		{"000011", "-", "-", "-"},
		{"001100", "-", "-", ":"},
		{"001010", "-", "-", ";"},
		{"001011", "-", "-", "“«"},
		{"000111", "-", "-", "”»"},
		{"001001", "-", "-", "?"},
		{"001110", "-", "-", "!"},
		{"001111", "-", "-", "()"},
		{"010111", "-", "-", "Далее следует цифра"},
	}

	arrText := strings.Split(text, " ")
	text = ""
	for _, item := range arrText {
		if item == "" {
			continue
		}

		for _, symbol := range symbols {
			if symbol.code == item {
				text += item + " = <b>RU</b> " + symbol.rusSymbol + " <b>EN</b> " + symbol.engSymbol + " <b>№</b> " + symbol.number + ";\n"
			}
		}
		if text == "" {
			text += "Символа " + item + " шрифте Браиля нет.\n"
		}
	}

	return text
}
func bin(text string, attribute bool) string {
	var t int
	var d int

	arrText := strings.Split(text, " ")
	text = ""
	for _, item := range arrText {
		//clear
		t = 0
		d = 1
		if item == "" {
			continue
		}

		number, err := strconv.Atoi(item)
		if err != nil {
			text += item + " = недопустимое число!\n"
			continue
		}
		if attribute {
			for {
				t += (number % 2) * d
				number = number / 2
				d = d * 10
				if number == 0 {
					if t < 0 {
						text += fmt.Sprintf("%s - изначально двоичка;\n", item)
					} else {
						text += fmt.Sprintf("%s (10) = %d (2);\n", item, t)
					}
					break
				}
			}
		} else {
			for {
				t += (number % 10) * d
				number = number / 10
				d = d * 2
				if number == 0 {
					text += fmt.Sprintf("%s (2) = %d (10);\n", item, t)
					break
				}
			}
		}
	}

	return text
}
func transferToAlphabet(text string, types bool) string {

	rusAlphabet := []string{"а", "б", "в", "г", "д", "е", "ё", "ж", "з", "и", "й", "к", "л", "м", "н", "о", "п", "р", "с", "т", "у", "ф", "х", "ц", "ч", "ш", "щ", "ъ", "ы", "ь", "э", "ю", "я"}
	engAlphabet := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "-", "-", "-", "-", "-", "-", "-"}

	if types {
		arrText := strings.Split(text, " ")

		text = ""
		for _, item := range arrText {
			if item == "" {
				continue
			}
			number, err := strconv.Atoi(item)
			if err != nil {
				text += item + " = недопустимое число!\n"
				continue
			}
			if (number < 33) && (number > 0) {
				text += item + " = <b>RU</b> " + rusAlphabet[number-1] + " <b>EN</b> " + engAlphabet[number-1] + ";\n"
			} else {
				text += item + " = недопустимое число!\n"
			}
		}
	} else {
		text = strings.ToLower(text)
		arrText := strings.Split(text, " ")
		text = ""
		for _, item := range arrText {
			if item == "" {
				continue
			}
			for i := 0; i < 33; i++ {
				if rusAlphabet[i] == item {
					text += fmt.Sprintf("Номер %s в <b>RU</b> алфавите %d;\n", strings.ToUpper(item), i+1)
				}
				if engAlphabet[i] == item {
					text += fmt.Sprintf("Номер %s в <b>EN</b> алфавите %d;\n", strings.ToUpper(item), i+1)
				}
			}
			if text == "" {
				text += "Буквы " + strings.ToUpper(item) + " с таким номером в алфавите нет.\n"
			}
		}
	}

	return text
}
func translateQwerty(text string) string {

	type table struct {
		eng string
		rus string
	}

	var qwertys = []table{
		{"q", "й"},
		{"w", "ц"},
		{"e", "у"},
		{"r", "к"},
		{"t", "е"},
		{"y", "н"},
		{"u", "г"},
		{"i", "ш"},
		{"o", "щ"},
		{"p", "з"},
		{"[", "х"},
		{"]", "ъ"},
		{"a", "ф"},
		{"s", "ы"},
		{"d", "в"},
		{"f", "а"},
		{"g", "п"},
		{"h", "р"},
		{"j", "о"},
		{"k", "л"},
		{"l", "д"},
		{";", "ж"},
		{"'", "э"},
		{"z", "я"},
		{"x", "ч"},
		{"c", "с"},
		{"v", "м"},
		{"b", "и"},
		{"n", "т"},
		{"m", "ь"},
		{",", "б"},
		{".", "ю"},
		{"/", "."},
	}

	arrText := strings.Split(text, " ")
	var attribute bool

	text = ""
	for _, item := range arrText {
		attribute = false
		for _, value := range qwertys {
			if (item == value.rus) || (item == value.eng) {
				text += item + " = <b>RU</b> " + value.rus + " <b>EN</b> " + value.eng + ";\n"
				attribute = true
				break
			}
		}
		if !attribute {
			text += "Символа " + strings.ToUpper(item) + " на клавиатуре нет.\n"
		}
	}

	return text
}

// TODO Шифр Цезаря
// TODO Шифр Хилла
// TODO Шифр Вижинера
// TODO Шифр Атбаш
// TODO Шифр Бодо
// TODO Шифр Бэкона
// TODO T9
// TODO ascii
// TODO gif [ссылка] — разбитие гифки на кадры.
// TODO r4 — решить полную расчлененку.
// TODO mr4 — решить расчлененку с несколькими словами

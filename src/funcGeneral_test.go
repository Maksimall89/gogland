package src

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestGeneralLogInit(t *testing.T) {
	path := "../log"
	err := os.RemoveAll(path)
	if err != nil {
		t.Errorf("Dir %s did not delete", path)
	}

	LogInit()
	log.Println("TEST")

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		t.Errorf("Dir %s does not exist", path)
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		t.Errorf("Dir is clear %s", path)
	}
	for _, file := range files {
		fContent, err := ioutil.ReadFile(path + "/" + file.Name())
		if err != nil {
			t.Errorf("Can not open %s/%s", path, file.Name())
		}
		if !strings.Contains(string(fContent), "Gogland") {
			t.Errorf("In file %s/%s don't have line", path, file.Name())
		}
	}
}
func TestGeneralSearchLocation(t *testing.T) {
	locMapGood := make(map[string]float64)
	initLocationTest(locMapGood)

	text := `dasdada 40°15′08″ 58°26′23″ 
			в. д.HGЯO sdfsdsdasd <b>dfsfsdf</b>\n 40.167841<br/>58.410761 \n 
			40,167845<br/>58,410765 //
			fsdfsd 59°24'48.6756", 58°24'38.7396"
			sdfsdgi :: 56°50.683, 53°11.776
			sdfsd56.8576 53.3529sdfsdf`
	locationMap := SearchLocation(text)

	for i := 0; i == len(locationMap)/2; i++ {
		if (locationMap["Latitude"+strconv.FormatInt(int64(i), 10)] != locMapGood["Latitude"+strconv.FormatInt(int64(i), 10)]) && (locationMap["Longitude"+strconv.FormatInt(int64(i), 10)] != locMapGood["Longitude"+strconv.FormatInt(int64(i), 10)]) {
			t.Errorf(
				"Latitude for %f\nLatitude expected %f\nLongitude for %f\nLongitude expected %f",
				locationMap["Latitude"+strconv.FormatInt(int64(i), 10)],
				locMapGood["Latitude"+strconv.FormatInt(int64(i), 10)],
				locationMap["Longitude"+strconv.FormatInt(int64(i), 10)],
				locMapGood["Longitude"+strconv.FormatInt(int64(i), 10)],
			)
		}
	}
}
func TestGeneralSendLocation(t *testing.T) {
	webToBotTEST := make(chan MessengerStyle, 10)
	var msgChanel MessengerStyle
	locationMap := make(map[string]float64)
	initLocationTest(locationMap)

	type testPair struct {
		latitude  float64
		longitude float64
	}
	var tests = []testPair{
		{40.2522222, 58.4363889},
		{40.167841, 58.410761},
		{40.167845, 58.410765},
		{59.413521, 58.410761},
		{56.84471667, 53.19626667},
		{56.8576, 53.3529},
		{0.0, 0.0},
	}

	SendLocation(locationMap, webToBotTEST)
	for _, location := range tests {
		select {
		// В канал msgChanel будут приходить все новые сообщения from web
		case msgChanel = <-webToBotTEST:
			// проверка сообщений
			if msgChanel.Type != "location" {
				t.Error("Error message type")
			}
			if msgChanel.Latitude != location.latitude || msgChanel.Longitude != location.longitude {
				t.Errorf("For %s\nexpected %v\ngot %v", msgChanel.ChannelMessage, location, locationMap)
			}
		default:
		}
	}
}
func TestGeneralSearchPhoto(t *testing.T) {
	t.Parallel()

	PhotoMapGood := make(map[int]string)
	PhotoMapGood[0] = `<img src="http://ya.ru/iads.png">`
	PhotoMapGood[1] = `<img src="http://sdfsd.png" title="sdfsf">`

	text := `<img src="http://ya.ru/iads.png">jdsfjklsjklds</b>fjklsdfjklsdfjklsd<img src="http://sdfsd.png" title="sdfsf">dfkjfjklsdfkjlsd`
	photoMap := SearchLocation(text)

	i := 0
	for photo := range photoMap {
		if photo != PhotoMapGood[i] {
			t.Errorf("For %s\nexpected %s\ngot %s", text, photo, PhotoMapGood[i])
		}
		i++
	}
}
func TestGeneralSendPhoto(t *testing.T) {
	webToBotTEST := make(chan MessengerStyle, 10)
	var msgChanel MessengerStyle
	photoMap := make(map[int]string)
	arrPhotos := []string{``, `<img src="y.img">`, `qwerty`}

	for _, photo := range arrPhotos {
		photoMap[0] = photo
		SendPhoto(photoMap, webToBotTEST)
		select {
		// В канал msgChanel будут приходить все новые сообщения from web
		case msgChanel = <-webToBotTEST:
			// проверка сообщений
			if msgChanel.Type != "photo" {
				t.Error("Error message type")
			}
			if msgChanel.ChannelMessage != photo {
				t.Errorf("For %s\nexpected %s\ngot %s", msgChanel.ChannelMessage, photo, photoMap[1])
			}
		default:
		}
	}
}
func TestGeneralReplaceTag(t *testing.T) {
	t.Parallel()

	type testPair struct {
		input  string
		output string
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
		result := ReplaceTag(pair.input, "RxR")
		if result != pair.output {
			t.Errorf("For %s\nexpected %s\ngot %s", pair.input, pair.output, result)
		}
	}
}
func TestGeneralDeleteMapFloat(t *testing.T) {
	t.Parallel()

	maps := make(map[string]float64)
	maps["asdas1"] = 45.4
	maps["asdas2"] = 145.423
	maps["asdas3"] = 415.42
	maps["asdas4"] = 45.43

	DeleteMapFloat(maps)

	NewLenMaps := len(maps)
	if NewLenMaps != 0 {
		t.Errorf("For delete len map 3 expected %d got %v", len(maps), NewLenMaps)
	}
}
func TestGeneralDeleteMap(t *testing.T) {
	t.Parallel()

	maps := make(map[int]string)
	maps[0] = "asdas1"
	maps[1] = "asdas2"
	maps[3] = "asdas3"
	maps[3] = "asdas4"

	DeleteMap(maps)

	NewLenMaps := len(maps)
	if NewLenMaps != 0 {
		t.Errorf("For delete len map 3 expected %d got %v", len(maps), NewLenMaps)
	}
}
func TestGeneralToFixed(t *testing.T) {
	t.Parallel()

	type inputStruct struct {
		num       float64
		precision int
	}
	type testPair struct {
		input  inputStruct
		output float64
	}

	var tests = []testPair{
		{inputStruct{3.2, 1}, 3.2},
		{inputStruct{3.2, 0}, 3},
		{inputStruct{3.2, 6}, 3.2},
		{inputStruct{3.7, 1}, 3.7},
		{inputStruct{0, 1}, 0},
	}

	for _, pair := range tests {
		result := ToFixed(pair.input.num, pair.input.precision)
		if result != pair.output {
			t.Errorf("For %v\nexpected %f\ngot %v", pair.input, pair.output, result)
		}
	}
}
func TestGeneralConvertUTF(t *testing.T) {
	t.Parallel()

	type testPair struct {
		input  string
		output string
	}

	var tests = []testPair{
		{`http://d1.endata.cx/data/games/60070/%D0%B0%D0%B3%D0%B5%D0%BD%D1%82+%D1%8F%D0%B3%D1%83%D0%BB.jpg`, `http://d1.endata.cx/data/games/60070/агент ягул.jpg`},
		{`http://en.cx.gameengines/encounter/play/27543/?pid=237876&amp;pact=2`, `http://en.cx.gameengines/encounter/play/27543/?pid=237876&amp;pact=2`},
		{`http://d1.endata.cx/data/games/63892/%d0%bc%d0%b5%d1%82%d0%ba%d0%b0+copy.jpg`, `http://d1.endata.cx/data/games/63892/метка copy.jpg`},
		{`http://d1.endata.cx/data/games/63892/%d0%bc%d0%b5%d1%82%d0%ba%d0%b0+copy.jpg`, `http://d1.endata.cx/data/games/63892/метка copy.jpg`},
	}

	for _, pair := range tests {
		result, _ := ConvertUTF(pair.input)
		if result != pair.output {
			t.Errorf("For %s\nexpected %s\ngot %s", pair.input, pair.output, result)
		}
	}
}
func TestGeneralRound(t *testing.T) {
	t.Parallel()

	type testPair struct {
		input  float64
		output int
	}

	var tests = []testPair{
		{0.01, 0},
		{2.1, 2},
		{4565.0, 4565},
		{4566.07, 4566},
		{4567.7, 4568},
		{0, 0},
	}

	for _, pair := range tests {
		result := Round(pair.input)
		if result != pair.output {
			t.Errorf("For %f\nexpected %d\ngot %d", pair.input, pair.output, result)
		}
	}
}
func TestGeneralConvertTimeSec(t *testing.T) {
	t.Parallel()

	type testPair struct {
		input  int
		output string
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
		result := ConvertTimeSec(pair.input)
		if result != pair.output {
			t.Errorf("For %d\nexpected %s\ngot %s", pair.input, pair.output, result)
		}
	}
}

func initLocationTest(maps map[string]float64) {
	maps["Latitude0"] = 40.2522222
	maps["Longitude0"] = 58.4363889

	maps["Latitude1"] = 40.167841
	maps["Longitude1"] = 58.410761

	maps["Latitude2"] = 40.167845
	maps["Longitude2"] = 58.410765

	maps["Latitude3"] = 59.413521
	maps["Longitude3"] = 58.410761

	maps["Latitude4"] = 56.84471667
	maps["Longitude4"] = 53.19626667

	maps["Latitude5"] = 56.8576
	maps["Longitude5"] = 53.3529
}

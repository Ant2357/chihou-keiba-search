package web_scraping

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
)

type Race struct {
	Name      string  `json:"name"`
	RaceTrack string  `json:"racetrack"`
	Type      string  `json:"type"`
	Distance  int     `json:"distance"`
	Horses    []Horse `json:"horses"`
}

type Horse struct {
	ImgUrl           string       `json:"img_url"`
	Name             string       `json:"name"`
	PlayGameCount    int          `json:"play_game_count"`
	Win              int          `json:"win"`
	Lose             int          `json:"lose"`
	Wp               float64      `json:"wp"`
	CourseAptitude   string       `json:"course_aptitude"`
	DistanceAptitude string       `json:"distance_aptitude"`
	RunningStyle     string       `json:"running_style"`
	HeavyRacetrack   string       `json:"heavy_racetrack"`
	Pedigree         Pedigree     `json:"pedigree"`
	Results          []RaceResult `json:"results"`
}

type Pedigree struct {
	Father              string `json:"father"`
	PaternalGrandfather string `json:"paternal_grandfather"`
	PaternalGrandmother string `json:"paternal_grandmother"`
	Mother              string `json:"mother"`
	MaternalGrandfather string `json:"maternal_grandfather"`
	MaternalGrandmother string `json:"maternal_grandmother"`
}

type RaceResult struct {
	Date     time.Time `json:"date"`
	RaceName string    `json:"raceName"`
	Result   int       `json:"result"`
	Type     string    `json:"type"`
	Distance int       `json:"distance"`
	Baba     string    `json:"baba"`
	Time     string    `json:"time"`
}

func toInt64(strVal string) int64 {
	rex := regexp.MustCompile("[0-9]+")
	strVal = rex.FindString(strVal)
	intVal, err := strconv.ParseInt(strVal, 10, 64)
	if err != nil {
		panic(err)
	}
	return intVal
}

func readRaceType(raceText string) string {
	if raceText[0:3] == "ダ" {
		return "ダート"
	} else {
		return "芝"
	}
}

func loadDocument(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	buf, _ := ioutil.ReadAll(res.Body)

	det := chardet.NewTextDetector()
	detRes, _ := det.DetectBest(buf)

	bReader := bytes.NewReader(buf)
	utf8Reader, _ := charset.NewReaderLabel(detRes.Charset, bReader)

	doc, err := goquery.NewDocumentFromReader(utf8Reader)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func ReadRace(url string) (Race, error) {
	doc, err := loadDocument(url)
	if err != nil {
		return Race{}, errors.New("ロードに失敗しました")
	}

	// 脚質
	runningStyles := map[string]string{}

	// 逃げ情報
	doc.Find("#Netkeiba_Race_Nar_Shutuba > div.Wrap.fc > div.RaceColumn02 > table > tbody > tr:nth-child(1) > td > div > .UmaName").Each(func(i int, s *goquery.Selection) {
		runningStyles[s.Text()] = "逃げ"
	})

	// 先行情報
	doc.Find("#Netkeiba_Race_Nar_Shutuba > div.Wrap.fc > div.RaceColumn02 > table > tbody > tr:nth-child(2) > td > div > .UmaName").Each(func(i int, s *goquery.Selection) {
		runningStyles[s.Text()] = "先行"
	})

	// 差し情報
	doc.Find("#Netkeiba_Race_Nar_Shutuba > div.Wrap.fc > div.RaceColumn02 > table > tbody > tr:nth-child(3) > td > div > .UmaName").Each(func(i int, s *goquery.Selection) {
		runningStyles[s.Text()] = "差し"
	})

	// 追い込み情報
	doc.Find("#Netkeiba_Race_Nar_Shutuba > div.Wrap.fc > div.RaceColumn02 > table > tbody > tr:nth-child(4) > td > div > .UmaName").Each(func(i int, s *goquery.Selection) {
		runningStyles[s.Text()] = "追い込み"
	})

	horses := make([]Horse, 0)
	doc.Find(".HorseInfo > div > div > span.HorseName > a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		horseDoc, err := loadDocument(href)
		if err != nil {
			panic("馬詳細情報の取得に失敗しました")
		}

		horseTableTbody := horseDoc.Find("#db_main_box > div.db_main_deta > div > div.db_prof_area_02 > table > tbody").Text()

		imgUrl, _ := horseDoc.Find(".db_photo_main").Attr("src")
		name := horseDoc.Find("#db_main_box > div.db_head.fc > div.db_head_name.fc > div.horse_title > h1").Text()
		playGameCount := int(toInt64(regexp.MustCompile(`[0-9]{1,}戦`).FindString(horseTableTbody)))
		win := int(toInt64(regexp.MustCompile(`[0-9]{1,}勝`).FindString(horseTableTbody)))
		lose := playGameCount - win
		wp := float64(float64(win)/float64(win+lose)) * 100.0

		// 適正情報(コース or 距離など)を読み込む関数
		// imgSelectorPath: 適正情報が書かれた画像のCSSセレクタPath
		// trueStr: 適正が真だった時に返す文字列
		// falseStr: 適正が偽だった時に返す文字列
		readAptitude := func(imgSelectorPath string, trueStr string, falseStr string) string {
			isTurfImgLink, _ := horseDoc.Find(imgSelectorPath).Attr("src")
			if strings.Contains(isTurfImgLink, "blue") {
				return trueStr
			} else {
				return falseStr
			}
		}

		// コース適正
		courseAptitude := readAptitude(
			"#db_main_box > div.db_main_deta > div > div.db_prof_area_01 > div.db_prof_box > dl > dd > table > tbody > tr:nth-child(1) > td > img:nth-child(1)",
			"ターフ",
			"ダート")

		// 距離適正
		distanceAptitude := readAptitude(
			"#db_main_box > div.db_main_deta > div > div.db_prof_area_01 > div.db_prof_box > dl > dd > table > tbody > tr:nth-child(2) > td > img:nth-child(1)",
			"スプリンター",
			"ステイヤー")

		// 脚質
		runningStyle := runningStyles[name[:9]]

		// 重馬場
		heavyRacetrack := readAptitude(
			"#db_main_box > div.db_main_deta > div > div.db_prof_area_01 > div.db_prof_box > dl > dd > table > tbody > tr:nth-child(5) > td > img:nth-child(1)",
			"得意",
			"苦手")

		// 血統情報のWebスクレイピング
		father := horseDoc.Find("#db_main_box > div.db_main_deta > div > div.db_prof_area_02 > div > dl > dd > table > tbody > tr:nth-child(1) > td:nth-child(1) > a").Text()
		paternalGrandfather := horseDoc.Find("#db_main_box > div.db_main_deta > div > div.db_prof_area_02 > div > dl > dd > table > tbody > tr:nth-child(1) > td:nth-child(2) > a").Text()
		paternalGrandmother := horseDoc.Find("#db_main_box > div.db_main_deta > div > div.db_prof_area_02 > div > dl > dd > table > tbody > tr:nth-child(2) > td > a").Text()

		mother := horseDoc.Find("#db_main_box > div.db_main_deta > div > div.db_prof_area_02 > div > dl > dd > table > tbody > tr:nth-child(3) > td.b_fml > a").Text()
		MaternalGrandfather := horseDoc.Find("#db_main_box > div.db_main_deta > div > div.db_prof_area_02 > div > dl > dd > table > tbody > tr:nth-child(3) > td.b_ml > a").Text()
		MaternalGrandmother := horseDoc.Find("#db_main_box > div.db_main_deta > div > div.db_prof_area_02 > div > dl > dd > table > tbody > tr:nth-child(4) > td > a").Text()

		pedigree := Pedigree{
			Father:              father,
			PaternalGrandfather: paternalGrandfather,
			PaternalGrandmother: paternalGrandmother,
			Mother:              mother,
			MaternalGrandfather: MaternalGrandfather,
			MaternalGrandmother: MaternalGrandmother,
		}

		// レース情報のWebスクレイピング
		raceResults := make([]RaceResult, 0)
		horseDoc.Find("#contents > div.db_main_race.fc > div > table > tbody > tr").Each(func(j int, selection *goquery.Selection) {
			if j >= 10 {
				return
			}

			jst, err := time.LoadLocation("Asia/Tokyo")
			if err != nil {
				panic("JSTの取得に失敗しました")
			}

			date, err := time.ParseInLocation("2006/01/02", selection.Find("td:nth-child(1)").Text(), jst)
			if err != nil {
				fmt.Println(err)
				panic("JSTへの変換に失敗しました")
			}

			raceName := selection.Find("td:nth-child(5)").Text()

			var result int64
			strResult := selection.Find("td:nth-child(12)").Text()
			if strResult == "" || strResult == "除" || strResult == "取" || strResult == "中" {
				result = -1
			} else {
				result = toInt64(strResult)
			}
			raceText := selection.Find("td:nth-child(15)").Text()

			raceType := readRaceType(raceText)
			distance := toInt64(raceText)

			baba := selection.Find("td:nth-child(16)").Text()
			time := selection.Find("td:nth-child(18)").Text()

			raceResults = append(raceResults, RaceResult{
				Date:     date,
				RaceName: raceName,
				Result:   int(result),
				Type:     raceType,
				Distance: int(distance),
				Baba:     baba,
				Time:     time,
			})
		})

		horses = append(horses, Horse{
			ImgUrl:           imgUrl,
			Name:             name,
			PlayGameCount:    playGameCount,
			Win:              win,
			Lose:             lose,
			Wp:               wp,
			CourseAptitude:   courseAptitude,
			DistanceAptitude: distanceAptitude,
			RunningStyle:     runningStyle,
			HeavyRacetrack:   heavyRacetrack,
			Pedigree:         pedigree,
			Results:          raceResults,
		})
	})

	raceName := strings.TrimSpace(doc.Find("div.RaceName").Text())
	raceTrack := doc.Find("#Netkeiba_Race_Nar_Shutuba > div.Wrap.fc > div.RaceColumn01 > div > div.RaceMainColumn > div.RaceList_NameBox > div.RaceList_Item02 > div.RaceData02 > span:nth-child(2)").Text()

	raceText := strings.TrimSpace(doc.Find("div.RaceData01 > span:nth-child(1)").Text())

	raceType := readRaceType(raceText)
	distance := toInt64(raceText)

	race := Race{
		Name:      raceName,
		RaceTrack: raceTrack,
		Type:      raceType,
		Distance:  int(distance),
		Horses:    horses,
	}

	return race, nil
}

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

type Horse struct {
	Name             string       `json:"name"`
	PlayGameCount    int          `json:"play_game_count"`
	Win              int          `json:"win"`
	Lose             int          `json:"lose"`
	CourseAptitude   string       `json:"course_aptitude"`
	DistanceAptitude string       `json:"distance_aptitude"`
	RunningStyle     string       `json:"running_style"`
	HeavyRacetrack   string       `json:"heavy_racetrack"`
	Results          []RaceResult `json:"results"`
}

type RaceResult struct {
	Date     time.Time `json:"date"`
	RaceName string    `json:"raceName"`
	Result   int       `json:"result"`
	Distance string    `json:"distance"`
	Baba     string    `json:"baba"`
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

func Horses(url string) ([]Horse, error) {
	doc, err := loadDocument(url)
	if err != nil {
		return nil, errors.New("ロードに失敗しました")
	}

	horses := make([]Horse, 0)
	doc.Find(".HorseInfo > div > div > span.HorseName > a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		horseDoc, err := loadDocument(href)
		if err != nil {
			panic("馬詳細情報の取得に失敗しました")
		}

		horseTableTbody := horseDoc.Find("#db_main_box > div.db_main_deta > div > div.db_prof_area_02 > table > tbody").Text()

		name := horseDoc.Find("#db_main_box > div.db_head.fc > div.db_head_name.fc > div.horse_title > h1").Text()
		playGameCount := toInt64(regexp.MustCompile(`[0-9]{1,}戦`).FindString(horseTableTbody))
		win := toInt64(regexp.MustCompile(`[0-9]{1,}勝`).FindString(horseTableTbody))
		lose := playGameCount - win

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
			"turf",
			"dirt")

		// 距離適正
		distanceAptitude := readAptitude(
			"#db_main_box > div.db_main_deta > div > div.db_prof_area_01 > div.db_prof_box > dl > dd > table > tbody > tr:nth-child(2) > td > img:nth-child(1)",
			"sprint",
			"styer")

		// 脚質
		runningStyle := readAptitude(
			"#db_main_box > div.db_main_deta > div > div.db_prof_area_01 > div.db_prof_box > dl > dd > table > tbody > tr:nth-child(3) > td > img:nth-child(1)",
			"front_runner",
			"hold_up_runner")

		// 重馬場
		heavyRacetrack := readAptitude(
			"#db_main_box > div.db_main_deta > div > div.db_prof_area_01 > div.db_prof_box > dl > dd > table > tbody > tr:nth-child(5) > td > img:nth-child(1)",
			"tokui",
			"nigate")

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
			result := toInt64(selection.Find("td:nth-child(12)").Text())
			distance := selection.Find("td:nth-child(15)").Text()
			baba := selection.Find("td:nth-child(16)").Text()

			raceResults = append(raceResults, RaceResult{
				Date:     date,
				RaceName: raceName,
				Result:   int(result),
				Distance: distance,
				Baba:     baba,
			})
		})

		horses = append(horses, Horse{
			Name:             name,
			PlayGameCount:    int(playGameCount),
			Win:              int(win),
			Lose:             int(lose),
			CourseAptitude:   courseAptitude,
			DistanceAptitude: distanceAptitude,
			RunningStyle:     runningStyle,
			HeavyRacetrack:   heavyRacetrack,
			Results:          raceResults,
		})
	})
	return horses, nil
}

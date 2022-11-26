package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
)

type Horse struct {
	Name             string `json:"name"`
	PlayGameCount    int    `json:"play_game_count"`
	Win              int    `json:"win"`
	Lose             int    `json:"lose"`
	CourseAptitude   string `json:"course_aptitude"`
	DistanceAptitude string `json:"distance_aptitude"`
	RunningStyle     string `json:"running_style"`
	HeavyRacetrack   string `json:"heavy_racetrack"`
}

func toInt64(strVal string) int64 {
	rex := regexp.MustCompile("[0-9]+")
	strVal = rex.FindString(strVal)
	intVal, err := strconv.ParseInt(strVal, 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	return intVal
}

func loadDocument(url string) *goquery.Document {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	buf, _ := ioutil.ReadAll(res.Body)

	det := chardet.NewTextDetector()
	detRes, _ := det.DetectBest(buf)

	bReader := bytes.NewReader(buf)
	utf8Reader, _ := charset.NewReaderLabel(detRes.Charset, bReader)

	doc, err := goquery.NewDocumentFromReader(utf8Reader)
	if err != nil {
		log.Fatal(err)
	}

	return doc
}

func horses(doc *goquery.Document) []Horse {
	horses := make([]Horse, 0)
	doc.Find(".HorseInfo > div > div > span.HorseName > a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		horsesDoc := loadDocument(href)
		horseTableTbody := horsesDoc.Find("#db_main_box > div.db_main_deta > div > div.db_prof_area_02 > table > tbody").Text()

		name := horsesDoc.Find("#db_main_box > div.db_head.fc > div.db_head_name.fc > div.horse_title > h1").Text()
		playGameCount := toInt64(regexp.MustCompile(`[0-9]{1,}戦`).FindString(horseTableTbody))
		win := toInt64(regexp.MustCompile(`[0-9]{1,}勝`).FindString(horseTableTbody))
		lose := playGameCount - win

		// 適正情報(コース or 距離など)を読み込む関数
		// imgSelectorPath: 適正情報が書かれた画像のCSSセレクタPath
		// trueStr: 適正が真だった時に返す文字列
		// falseStr: 適正が偽だった時に返す文字列
		readAptitude := func(imgSelectorPath string, trueStr string, falseStr string) string {
			isTurfImgLink, _ := horsesDoc.Find(imgSelectorPath).Attr("src")
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

		horses = append(horses, Horse{
			Name:             name,
			PlayGameCount:    int(playGameCount),
			Win:              int(win),
			Lose:             int(lose),
			CourseAptitude:   courseAptitude,
			DistanceAptitude: distanceAptitude,
			RunningStyle:     runningStyle,
			HeavyRacetrack:   heavyRacetrack,
		})
	})
	return horses
}

func main() {
	netKeibaURL := "https://nar.netkeiba.com/race/shutuba.html?race_id=202244111711&rf=race_submenu"
	horses := horses(loadDocument(netKeibaURL))

	file, err := os.Create(`output.json`)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	bytes, _ := json.MarshalIndent(horses, "", "  ")
	file.Write(([]byte)(string(bytes)))

	fmt.Println("SUCCESS")
}

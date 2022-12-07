package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	web_scraping "changeme/web_scraping"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) JsonOutputKeiba(url string) (string, error) {
	// netKeibaURL := "https://nar.netkeiba.com/race/shutuba.html?race_id=202244111711&rf=race_submenu"

	horses, err := web_scraping.Horses(url)
	if err != nil {
		return "", errors.New("馬情報の取得に失敗しました")
	}

	if err := os.MkdirAll("output", 0777); err != nil {
		return "", errors.New("馬情報フォルダの作成に失敗しました")
	}

	file, err := os.Create(`./output/output.json`)
	if err != nil {
		return "", errors.New("JSONの出力に失敗しました")
	}
	defer file.Close()

	bytes, err := json.MarshalIndent(horses, "", "  ")
	if err != nil {
		return "", errors.New("JSONの出力に失敗しました")
	}

	file.Write(([]byte)(string(bytes)))

	return "ローカルにjsonファイルを出力しました", nil
}

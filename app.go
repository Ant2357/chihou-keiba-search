package main

import (
	"context"
	"errors"

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

func (a *App) WebScrapingRace(url string) (web_scraping.Race, error) {
	// netKeibaURL := "https://nar.netkeiba.com/race/shutuba.html?race_id=202244111711&rf=race_submenu"

	race, err := web_scraping.ReadRace(url)
	if err != nil {
		return web_scraping.Race{}, errors.New("馬情報の取得に失敗しました")
	}

	return race, nil
}

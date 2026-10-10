package main

import (
	"embed"
	"lol-teammate-helper/internal/controller"
	"lol-teammate-helper/internal/logging"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	_, closeLog := logging.Setup()
	defer closeLog()

	// Create an instance of the app structure
	app := NewApp()
	mc := controller.NewMatchHistory()
	// Create application with options
	err := wails.Run(&options.App{
		Title:     "lol-teammate-helper",
		Width:     1280,
		Height:    820,
		MinWidth:  1100,
		MinHeight: 640,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 1, G: 10, B: 19, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
			mc,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

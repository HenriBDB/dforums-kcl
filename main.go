package main

import (
	"dforum-app/configuration"
	"dforum-app/network"
	"dforum-app/storage"
	"dforum-app/view"
	_ "embed"
	"os"

	"github.com/wailsapp/wails"
)

//go:embed frontend/build/static/js/main.js
var js string

//go:embed frontend/build/static/css/main.css
var css string

var vHandle *view.ViewHandler
var networkHandle *network.NetworkModule

func main() {
	wd, _ := os.Getwd()
	configuration.InitConfigs(wd)
	logPath := wd + string(os.PathSeparator) + "dfd.log"
	logFile, _ := os.Create(logPath)
	logFile.Close() // Create a new empty file or truncate existing
	configuration.InitLogger(logPath)

	storageModule := storage.NewStorageModule(configuration.GetDatabasePath())
	defer storageModule.TearDown()

	vHandle = view.NewViewHandler(storageModule)

	networkHandle = network.NewNetworkModule(storageModule)
	networkHandle.CreateAndStartHost()
	defer networkHandle.TearDown()

	app := wails.CreateApp(&wails.AppConfig{
		Width:     1024,
		Height:    768,
		Title:     "dforums-app",
		JS:        js,
		CSS:       css,
		Colour:    "#131313",
		Resizable: true,
	})
	app.Bind(configuration.GetJsonConfigs)
	app.Bind(configuration.UpdateConfig)
	app.Bind(vHandle)
	app.Run()
}

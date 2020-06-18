package main

import (
	"os"
	"path"

	"github.com/grvlle/molly_installer/backend/install"
	"github.com/leaanthony/mewn"
	log "github.com/sirupsen/logrus"
	"github.com/wailsapp/wails"
)

var installer *install.Install

func init() {
	var err error

	initLogger() // log to .dag/install.log

	installer, err = install.Init()
	if err != nil {
		panic(err)
	}
}

func main() {

	js := mewn.String("./frontend/dist/app.js")
	css := mewn.String("./frontend/dist/app.css")

	app := wails.CreateApp(&wails.AppConfig{
		Width:     1024,
		Height:    768,
		Title:     "install",
		JS:        js,
		CSS:       css,
		Colour:    "#131313",
		Resizable: true,
	})

	app.Bind(installer)
	app.Bind(runInstaller)
	app.Bind(runUninstaller)
	app.Run()

}

// Called from frontend
func runInstaller() {
	installer.Run()
}

// Called from frontend
func runUninstaller() {
	installer.Uninstall()
}

func initLogger() {

	userDir, _ := os.UserHomeDir()

	// initialize update.log file and set log output to file
	file, err := os.OpenFile(path.Join(userDir, ".dag", "install.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)

	log.Infoln("--------------------------------- Logger Initialized -------------------------------------")
}

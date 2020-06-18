package main

import (
	"github.com/grvlle/molly_installer/backend/install"
	"github.com/leaanthony/mewn"
	"github.com/wailsapp/wails"
)

var installer *install.Install

func init() {
	var err error

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

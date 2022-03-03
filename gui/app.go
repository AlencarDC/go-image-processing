package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type App struct {
	instance fyne.App
	windows  map[string]fyne.Window
}

func (app *App) Start() {
	window, hasMainWindow := app.windows["main"]

	if !hasMainWindow {
		log.Println("app: main window is not defined")
		return
	}

	log.Println("app: started")
	window.ShowAndRun()
}

func (app *App) NewWindow(id, name string, width, height float32) fyne.Window {
	window := app.instance.NewWindow(name)
	window.Resize(fyne.NewSize(width, height))

	app.windows[id] = window

	return window
}

func NewApp(name string, width, height float32) (*App, error) {

	application := new(App)
	application.instance = app.New()
	application.windows = make(map[string]fyne.Window)

	window := application.NewWindow("main", name, width, height)
	screen := NewMainScreen(*application, window)
	window.SetContent(screen.Content())

	return application, nil
}

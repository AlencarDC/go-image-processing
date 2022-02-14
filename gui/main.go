package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {

	application := app.New()
	window := application.NewWindow("Photochopp v1.0")
	window.Resize(fyne.NewSize(1280, 720))

	mainScreen := NewMainScreen(&window)

	window.SetContent(mainScreen.ctnMain)
	window.ShowAndRun()
}

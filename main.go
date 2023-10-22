package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/storage"
)

var cfg App

var filter = storage.NewExtensionFileFilter([]string{".md", ".MD"})

func main() {

	cfg.CurrentDir = getHomeDirOrFail()

	app := app.New()
	window := app.NewWindow("MarkdownEditor")

	cfg.win = &window

	cfg.CreateUiElements()
	cfg.UpdateUi()

	window.Resize(fyne.Size{Width: 1024, Height: 1024})
	window.ShowAndRun()
}

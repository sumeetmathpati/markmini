package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
)

var cfg Config

var filter = storage.NewExtensionFileFilter([]string{".md", ".MD"})

func main() {

	app := app.New()

	window := app.NewWindow("MarkdownEditor")

	edit, preview, explorer := cfg.MakeUI()
	cfg.CreateMenuItems(window)

	explorerSplit := container.NewHSplit(explorer, container.NewHSplit(edit, preview))
	explorerSplit.SetOffset(0.15)
	window.SetContent(explorerSplit)
	window.Resize(fyne.Size{Width: 1024, Height: 1024})
	window.ShowAndRun()
}

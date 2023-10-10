package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Config struct {
	EditWidget    *widget.Entry
	PreviewWidget *widget.RichText
	CurrentFile   fyne.URI
	SaveMenuItem  *fyne.MenuItem
}

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

func (app *Config) MakeUI() (*widget.Entry, *widget.RichText, *widget.List) {

	data := []string{"1", "2", "3"}

	explorer := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			button := widget.NewButton("Do Something", nil)
			button.SetIcon(theme.FileIcon())
			return button
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Button).OnTapped = func() {
				fmt.Println("I am button " + data[i])
			}

		})

	edit := widget.NewMultiLineEntry()
	preview := widget.NewRichTextFromMarkdown("")
	preview.Scroll = container.ScrollBoth

	app.EditWidget = edit
	app.PreviewWidget = preview

	edit.OnChanged = preview.ParseMarkdown

	return edit, preview, explorer
}

func (app *Config) CreateMenuItems(window fyne.Window) {
	openMenuItem := fyne.NewMenuItem("Open", app.openFunc(window))
	saveMenuItem := fyne.NewMenuItem("Save", app.saveFunc(window))
	saveAsMenuItem := fyne.NewMenuItem("Save as", app.saveAsFunc(window))

	app.SaveMenuItem = saveMenuItem
	app.SaveMenuItem.Disabled = true

	fileMenu := fyne.NewMenu("File", openMenuItem, saveMenuItem, saveAsMenuItem)
	menu := fyne.NewMainMenu(fileMenu)
	window.SetMainMenu(menu)
}

func (app *Config) saveFunc(win fyne.Window) func() {
	return func() {
		if app.CurrentFile != nil {
			write, err := storage.Writer(app.CurrentFile)
			if err != nil {
				dialog.ShowError(err, win)
			}

			write.Write([]byte(app.EditWidget.Text))
			defer write.Close()
		}
	}
}

func (app *Config) saveAsFunc(win fyne.Window) func() {
	return func() {
		saveDialog := dialog.NewFileSave(func(write fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if write == nil {
				return
			}

			if !strings.HasSuffix(strings.ToLower(write.URI().String()), ".md") {
				dialog.ShowInformation("Error", "Please make sure your filename has '.md' extension.", win)
			}
			write.Write([]byte(app.EditWidget.Text))
			app.CurrentFile = write.URI()

			defer write.Close()

			win.SetTitle(win.Title() + " - " + write.URI().Name())

			app.SaveMenuItem.Disabled = false
		}, win)

		saveDialog.SetFileName("untitled.md")
		saveDialog.SetFilter(filter)
		saveDialog.Show()
	}
}

func (app *Config) openFunc(win fyne.Window) func() {
	return func() {
		openDialog := dialog.NewFileOpen(func(read fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			if read == nil {
				return
			}

			defer read.Close()

			data, err := ioutil.ReadAll(read)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			app.EditWidget.SetText(string(data))
			app.CurrentFile = read.URI()
			win.SetTitle(win.Title() + " - " + read.URI().Name())
			app.SaveMenuItem.Disabled = false
		}, win)

		openDialog.SetFilter(filter)
		openDialog.Show()
	}
}

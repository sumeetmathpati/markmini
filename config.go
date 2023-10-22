package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type App struct {
	Explorer      *widget.List
	EditWidget    *widget.Entry
	PreviewWidget *widget.RichText
	CurrentFile   fyne.URI
	CurrentDir    string // This will be shown in expolorer.
	SaveMenuItem  *fyne.MenuItem
	win           *fyne.Window
}

func (app *App) CreateUiElements() (*widget.Entry, *widget.RichText, *widget.List) {

	// File explorer
	explorer := app.makeExplorer()
	app.Explorer = explorer

	// MarkdownEditor
	edit := widget.NewMultiLineEntry()
	app.EditWidget = edit

	// Rich text output
	preview := widget.NewRichTextFromMarkdown("")
	preview.Scroll = container.ScrollBoth
	app.PreviewWidget = preview

	edit.OnChanged = preview.ParseMarkdown

	app.createMenuItems()

	return edit, preview, explorer
}

func (app *App) UpdateUi() {
	explorerSplit := container.NewHSplit(app.Explorer, container.NewHSplit(app.EditWidget, app.PreviewWidget))
	explorerSplit.SetOffset(0.15)
	(*app.win).SetContent(explorerSplit)
}

func (app *App) makeExplorer() *widget.List {
	files := app.getCurrentDirFiles()

	explorer := widget.NewList(
		func() int {
			return len(files)
		},
		func() fyne.CanvasObject {
			button := widget.NewButton("Do Something", nil)
			button.SetIcon(theme.FileIcon())
			button.Alignment = widget.ButtonAlignLeading
			return button
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Button).SetText(files[i].Name())
			if files[i].IsDir() {
				o.(*widget.Button).SetIcon(theme.FolderIcon())
			}
			o.(*widget.Button).OnTapped = app.onFileClick(o.(*widget.Button))
		})
	return explorer
}

func (app *App) getCurrentDirFiles() []os.DirEntry {
	/*
		Get the files and folders in the CurrentDir attribute, and
		return list of strings.
	*/

	files, err := os.ReadDir(app.CurrentDir)
	if err != nil {
		fmt.Fprintln(os.Stdout, "Error reading directory:", err)
		os.Exit(1)
	}

	return files
}

func (app *App) onFileClick(btn *widget.Button) func() {
	return func() {

		absolutePath := filepath.Join(app.CurrentDir, btn.Text)
		isDir, err := isDir(absolutePath)
		if err != nil {
			dialog.ShowError(err, *app.win)
		}

		if isDir {
			app.CurrentDir = absolutePath
			app.Explorer = app.makeExplorer()
			app.UpdateUi()
		} else {

		}
	}
}

func (app *App) createMenuItems() {
	openMenuItem := fyne.NewMenuItem("Open", app.openFunc(*app.win))
	saveMenuItem := fyne.NewMenuItem("Save", app.saveFunc(*app.win))
	saveAsMenuItem := fyne.NewMenuItem("Save as", app.saveAsFunc(*app.win))

	app.SaveMenuItem = saveMenuItem
	app.SaveMenuItem.Disabled = true

	fileMenu := fyne.NewMenu("File", openMenuItem, saveMenuItem, saveAsMenuItem)
	menu := fyne.NewMainMenu(fileMenu)
	(*app.win).SetMainMenu(menu)
}

func (app *App) openFunc(win fyne.Window) func() {
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

			data, err := io.ReadAll(read)
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

func (app *App) saveFunc(win fyne.Window) func() {
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

func (app *App) saveAsFunc(win fyne.Window) func() {
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

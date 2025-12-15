package main

import "C"

import (
	"image"
	"fmt"
//	"regexp"
//	"strconv"
//	"strings"
	"os"
//	"bufio"
//	"time"

	"fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
)

var w fyne.Window
var a fyne.App

// input keys and keypress checks for canvas/ window
// since this is all that is called without other handler / timers
// - this is where maze update and edits will vector

func typedRune(r rune) {

//	fmt.Printf("in keys event - %x\n",r)
	if r == 'q' { os.Exit(0) }

if deskCanvas, ok := w.Canvas().(desktop.Canvas); ok {
        deskCanvas.SetOnKeyDown(func(key *fyne.KeyEvent) {
            fmt.Printf("Desktop key down: %h\n", key.Name)
        })
        deskCanvas.SetOnKeyUp(func(key *fyne.KeyEvent) {
//            fmt.Printf("Desktop key up: %v\n", key)
			if key.Name == "Escape" { os.Exit(0) }
       })
    }
}

// init app and win

func aw_init() {

    a = app.New()
    w = a.NewWindow("G¹G²ved")

	menuItemExit := fyne.NewMenuItem("Exit", func() {
		os.Exit(0)
	})
	menuExit := fyne.NewMenu("Exit ", menuItemExit)
	menuItemAbout := fyne.NewMenuItem("About", func() {
		dialog.ShowInformation("About G¹G²ved", "Gauntlet / Gauntlet 2 visual editor\nAuthor: Six [a programmer]\n\ngithub.com/six-of-one/", w)
	})
	menuItemLIC := fyne.NewMenuItem("License", func() {
		dialog.ShowInformation("G¹G²ved License", "Gauntlet visual editor\n\n(c) 2025 Six [a programmer]\n\nGPLv3.0\n\nhttps://www.gnu.org/licenses/gpl-3.0.html", w)
	})
	menuHint := fyne.NewMenu("?, q, fFgG, wWeE, rRt, hm, s, il, u, v, #a")

	menuHelp := fyne.NewMenu("Help ", menuItemAbout, menuItemLIC)
	mainMenu := fyne.NewMainMenu(menuExit, menuHelp, menuHint)
	w.SetMainMenu(mainMenu)
	w.Canvas().SetOnTypedRune(typedRune)
}

// update contents

func upwin(simg *image.NRGBA, mazeN int) {

	bimg := canvas.NewRasterFromImage(simg)
	w.Canvas().SetContent(bimg)
	w.Resize(fyne.NewSize(1024, 1024))
	w.Show()

	til := fmt.Sprintf("G¹G²ved Maze: %d",mazeN + 1)
	w.SetTitle(til)
}
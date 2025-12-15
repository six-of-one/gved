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

// pad for dialog page

func cpad(st string, d int) string {

	spout := st+"                                                                          "
	return string(spout[:d])
}

// init app and win

func aw_init() {

    a = app.New()
    w = a.NewWindow("G¹G²ved")

	menuItemExit := fyne.NewMenuItem("Exit", func() {
		os.Exit(0)
	})
	menuExit := fyne.NewMenu("Exit ", menuItemExit)
	menuItemKeys := fyne.NewMenuItem("Keys ?", func() {

		strp := cpad("single letter commands",36)
		strp += cpad("\n\n? - this list",38)
		strp += cpad("\nq - quit program",37)
		strp += fmt.Sprintf("\nf - floor pattern+\nF - floor pattern-\n" +
			"g - floor color+\nG - floor color-\nw - wall pattern+\nW - wall pattern-\n" +
			"e - wall color+\nE - wall color-\nr - rotate maze +90°\nR - rotate maze -90°\n" +
			"t - turn off rotate\nh - mirror maze horizontal toggle\nm - mirror maze vertical toggle\ns - toggle rnd special potion\n" +
			"i - gauntlet mazes r1 - r9\nl - use gauntlet rev 14\nu - gauntlet 2 mazes\nv - valid address list\n" +
			"{n}umeric of valid maze\n - load maze 1 - 127 g1\n - load maze 1 - 117 g2\n - load address 229376 - 262143\n" +
			" * note some address causes crash\n" +
			"    commands can be chained:\ni.e. i5a switch to g1, load maze 5\n" +
			"G%d ",opts.Gtp)
		if opts.R14 { strp += "(r14)"
			} else { strp += "(r1-9)" }

		dialog.ShowInformation("Command Keys", strp, w)
	})
	menuItemAbout := fyne.NewMenuItem("About", func() {
		dialog.ShowInformation("About G¹G²ved", "Gauntlet / Gauntlet 2 visual editor\nAuthor: Six [a programmer]\n\ngithub.com/six-of-one/", w)
	})
	menuItemLIC := fyne.NewMenuItem("License", func() {
		dialog.ShowInformation("G¹G²ved License", "Gauntlet visual editor\n\n(c) 2025 Six [a programmer]\n\nGPLv3.0\n\nhttps://www.gnu.org/licenses/gpl-3.0.html", w)
	})
	menuHint := fyne.NewMenu("?, q, fFgG, wWeE, rRt, hm, s, il, u, v, #a")

	menuHelp := fyne.NewMenu("Help ", menuItemKeys, menuItemAbout, menuItemLIC)
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
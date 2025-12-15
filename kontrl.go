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
		strp += "\n–—–—–—–—–—–—–—–—–—–—–—"
//		strp += cpad("\n\n? - this list",52)
		strp += cpad("\nq - quit program",42)
		strp += cpad("\nf - floor pattern+",43)
		strp += cpad("\ng - floor color+",45)
		strp += cpad("\nw - wall pattern+",43)
		strp += cpad("\ne - wall color+",46)
		strp += cpad("\nr - rotate maze +90°",41)
		strp += cpad("\nR - rotate maze -90°",42)
		strp += cpad("\nh - mirror maze horizontal toggle",31)
		strp += "\nm - mirror maze vertical toggle"
		strp += cpad("\ns - toggle rnd special potion",34)
		strp += cpad("\ni - gauntlet mazes r1 - r9",38)
		strp += cpad("\nl - use gauntlet rev 14",40)
		strp += cpad("\nu - gauntlet 2 mazes",39)
		strp += cpad("\nv - valid address list",42)
		strp += cpad("\n{n}umeric of valid maze",37)
		strp += cpad("\n - load maze 1 - 127 g1",42)
		strp += cpad("\n - load maze 1 - 117 g2",42)
		strp += "\n - load address 229376 - 262143 "
		strp += "\n * note some address will crash"
//		strp += cpad("\n    commands can be chained:",38)
//		strp += "\n- i5a switch to g1, load maze 5"
		strp += "\n–—–—–—–—–—–—–—–—–—–—–—"
		strb := fmt.Sprintf("\nG%d ",opts.Gtp)
		if opts.R14 { strb += "(r14)"
			} else { strb += "(r1-9)" }
		strp += cpad(strb,50)

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
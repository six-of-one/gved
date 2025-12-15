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

var anum int
var shift bool

func typedRune(r rune) {

// special aux string - put ops in title after maze #
	spau := ""
// relod
	relod := false
	relodsub := false

//	fmt.Printf("in keys event - %x\n",r)
	if r == 'q' { os.Exit(0) }

// new maze
	if r == 'a' {
		if (anum > 0 && anum <= 127 || anum >= 229376 && anum < 262145) {

			if anum <= 127 {
				opts.mnum = anum - 1
				Aov = 0
			} else {
				Aov = anum
				opts.mnum = 0
				spau = fmt.Sprintf("addr = %d",anum)
			}
			anum = 0
// clear these when load new maze
			Ovwallpat = -1
			relod = true
		}
	}

// (almost) blind numeric input
	switch r {
	case '0':
		anum = (anum * 10)
	case '1':
		anum = (anum * 10) + 1
	case '2':
		anum = (anum * 10) + 2
	case '3':
		anum = (anum * 10) + 3
	case '4':
		anum = (anum * 10) + 4
	case '5':
		anum = (anum * 10) + 5
	case '6':
		anum = (anum * 10) + 6
	case '7':
		anum = (anum * 10) + 7
	case '8':
		anum = (anum * 10) + 8
	case '9':
		anum = (anum * 10) + 9
	case '`':
		anum = 0
	}
	if r >= '0' && r <= '9' || r == '`' {
		til := fmt.Sprintf("numeric: %d", anum)
		uptitl(opts.mnum, til)
	}

if deskCanvas, ok := w.Canvas().(desktop.Canvas); ok {
        deskCanvas.SetOnKeyDown(func(key *fyne.KeyEvent) {
            fmt.Printf("Desktop key down: %h\n", key.Name)
			if key.Name == "LeftShift" { shift = true }
			if key.Name == "RightShift" { shift = true }
        })
        deskCanvas.SetOnKeyUp(func(key *fyne.KeyEvent) {
            fmt.Printf("Desktop key up: %v\n", key)
			if key.Name == "Escape" { os.Exit(0) }
			if key.Name == "LeftShift" { shift = false }
			if key.Name == "RightShift" { shift = false }
       })
    }
	fmt.Printf("r %v shift %v\n",r,shift)

		relodsub = true
		switch r {
		case 'z':
			Ovwallpat = -1
			opts.mnum -= 1
			if opts.mnum < 0 { opts.mnum = maxmaze }
		case 'x':
			Ovwallpat = -1
			opts.mnum += 1
			if opts.mnum > maxmaze { opts.mnum = 0 }
		case 'w':
			Ovwallpat += 1
			if anum > 0 { Ovwallpat = anum - 1; anum = 0 }
			if Ovwallpat > 7 { Ovwallpat = 0 }
			spau = fmt.Sprintf("cmd: w - wallp: %d\n",Ovwallpat)
		case 87:		// W
			Ovwallpat -= 1
			if Ovwallpat < 0 { Ovwallpat = 7 }
			spau = fmt.Sprintf("cmd: w - wallp: %d\n",Ovwallpat)
			relod = true
		case 'e':
			Ovwallcol += 1
			if anum > 0 { Ovwallcol = anum - 1; anum = 0 }
			if Ovwallcol > 16 { Ovwallcol = 0 }
			spau = fmt.Sprintf("cmd: e - wallc: %d\n",Ovwallcol)
		case 69:		// E
			Ovwallcol -= 1
			if Ovwallcol < 0 { Ovwallcol = 16 }
			spau = fmt.Sprintf("cmd: e - wallc: %d\n",Ovwallcol)
			relod = true
		case 'f':
			Ovflorpat += 1
			if anum > 0 { Ovflorpat = anum - 1; anum = 0 }
			if Ovflorpat > 8 { Ovflorpat = 0 }
			spau = fmt.Sprintf("cmd: f - floorp: %d\n",Ovflorpat)
		case 70:		// F
			Ovflorpat -= 1
			if Ovflorpat < 0 { Ovflorpat = 8 }
			spau = fmt.Sprintf("cmd: f - floorp: %d\n",Ovflorpat)
			relod = true
		case 'g':
			Ovflorcol += 1
			if anum > 0 { Ovflorcol = anum - 1; anum = 0 }
			if Ovflorcol > 15 { Ovflorcol = 0 }
			spau = fmt.Sprintf("cmd: g - floorc: %d\n",Ovflorcol)
		case 71:		// G
			Ovflorcol -= 1
			if Ovflorcol < 0 { Ovflorcol = 15 }
			spau = fmt.Sprintf("cmd: g - floorc: %d\n",Ovflorcol)
			relod = true
		case 'r':
			opts.MRP = true
			opts.MRM = false
			spau = fmt.Sprintf("cmd: r - mr+: %t mr-: %t\n",opts.MRP,opts.MRM)
		case 82:		// R
			opts.MRP = false
			opts.MRM = true
			spau = fmt.Sprintf("cmd: r - mr+: %t mr-: %t\n",opts.MRP,opts.MRM)
			relod = true
		case 't':
			opts.MRP = false
			opts.MRM = false
			spau = fmt.Sprintf("cmd: t - mr+: %t mr-: %t\n",opts.MRP,opts.MRM)
		case 'm':
			opts.MV = !opts.MV
			spau = fmt.Sprintf("cmd: m - mv: %t\n",opts.MV)
		case 'h':
			opts.MH = !opts.MH
			spau = fmt.Sprintf("cmd: h - mh: %t\n",opts.MH)
		case 'i':
			opts.Gtp = 1
			opts.R14 = false
			G1 = true
			G2 = false
			maxmaze = 126
			spau = "G¹ "
		case 'l':
			opts.Gtp = 1
			opts.R14 = !opts.R14
			G1 = true
			G2 = false
			maxmaze = 126
			spau = "G¹ "
		case 's':
			opts.SP = !opts.SP
		case 'u':
			opts.Gtp = 2
			G1 = false
			G2 = true
			maxmaze = 116
			spau = "G² mazes"
		case 'v':
			lx := 116
			if G1 { lx = 126 }
			fmt.Printf("\n valid maze address for Gauntlet %d\nmaze   dec -    hex\n",opts.Gtp)
			for x := 0; x <= lx;x ++ {
					ad := slapsticMazeGetRealAddr(x)
					fmt.Printf("%03d:%d - x%X  ",x + 1,ad,ad)
					if (x + 1) % 7 == 0 { fmt.Printf("\n") }
			}
			fmt.Printf("\n")
			dialog.ShowInformation("G¹G²ved", "Gauntlet / Gauntlet 2 valid maze address list\nplease check terminal where gved command was issued\n\ngithub.com/six-of-one/", w)
		case 63:
			keyhints()
		default:
			relodsub = false
		}

		if spau == "G¹ " {
			if opts.R14 { spau += "rv14" } else { spau += "rv1-9" }
		}
	if (relod || relodsub) {
		maze := mazeDecompress(slapsticReadMaze(opts.mnum), false)
		mazeloop(maze)
		Ovimg := genpfimage(maze, opts.mnum)
		upwin(Ovimg)
		uptitl(opts.mnum, spau)
	}
}

// pad for dialog page

func cpad(st string, d int) string {

	spout := st+"                                                                          " // jsut guess at a pad fill
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
	menuItemKeys := fyne.NewMenuItem("Keys ?", keyhints)
	menuItemAbout := fyne.NewMenuItem("About", func() {
		dialog.ShowInformation("About G¹G²ved", "Gauntlet / Gauntlet 2 visual editor\nAuthor: Six [a programmer]\n\ngithub.com/six-of-one/", w)
	})
	menuItemLIC := fyne.NewMenuItem("License", func() {
		dialog.ShowInformation("G¹G²ved License", "Gauntlet visual editor\n\n(c) 2025 Six [a programmer]\n\nGPLv3.0\n\nhttps://www.gnu.org/licenses/gpl-3.0.html", w)
	})
	menuHint := fyne.NewMenu("cmds: ?, q, fFgG, wWeE, rRt, hm, s, il, u, v, #a")

	menuHelp := fyne.NewMenu("Help ", menuItemKeys, menuItemAbout, menuItemLIC)
	mainMenu := fyne.NewMainMenu(menuExit, menuHelp, menuHint)
	w.SetMainMenu(mainMenu)
	w.Canvas().SetOnTypedRune(typedRune)
	anum = 0
	shift = false
}

// update contents

func upwin(simg *image.NRGBA) {

	bimg := canvas.NewRasterFromImage(simg)
	w.Canvas().SetContent(bimg)
	w.Resize(fyne.NewSize(1024, 1024))
//	w.Show()

	uptitl(opts.mnum, "")
}

// title special info update

func uptitl(mazeN int, spaux string) {

	til := fmt.Sprintf("G¹G²ved Maze: %d",mazeN + 1)
	if spaux != "" { til += " -- " + spaux }
	w.SetTitle(til)
}

func keyhints() {
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
//		strp += cpad("\nv - valid address list",42)
	strp += "\nv - all maze addr (in termninal)"
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
}
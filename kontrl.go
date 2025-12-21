package main

import (
	"image"
	"fmt"
	"math"
	"os"
	"io/ioutil"
	"time"

	"fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
//    "fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// kontrol is for fyne window ops & input management

var w fyne.Window
var a fyne.App

// input keys and keypress checks for canvas/ window
// since this is all that is called without other handler / timers
// - this is where maze update and edits will vector

var anum int
var shift bool
var ctrl bool

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
				Aov = addrver(anum, 0)
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
//            fmt.Printf("Desktop key down: %h\n", key.Name)
			if key.Name == "BackSpace" {
				anum = (anum / 10);
				til := fmt.Sprintf("numeric: %d", anum)
				uptitl(opts.mnum, til)
			}
			if key.Name == "LeftShift" { shift = true }
			if key.Name == "RightShift" { shift = true }
			if key.Name == "LeftControl" { ctrl = true }
			if key.Name == "RightControl" { ctrl = true }
        })
        deskCanvas.SetOnKeyUp(func(key *fyne.KeyEvent) {
//            fmt.Printf("Desktop key up: %v\n", key)
			if key.Name == "Escape" { os.Exit(0) }
			if key.Name == "LeftShift" { shift = false }
			if key.Name == "RightShift" { shift = false }
			if key.Name == "LeftControl" { ctrl = false }
			if key.Name == "RightControl" { ctrl = false }
			if key.Name == "S" && ctrl { menu_sav() }
			if key.Name == "L" && ctrl  { menu_lod() }
			if key.Name == "R" && ctrl  { menu_res() }
       })
    }
//	fmt.Printf("r %v shift %v\n",r,shift)

		relodsub = true
		switch r {
		case 65:		// A
			if Aov > 0 { Aov = 0 } else {
				Aov = addrver(slapsticMazeGetRealAddr(opts.mnum), 0)
			}
		case 'z':
			Ovwallpat = -1
// allow step parse through valid address
			if Aov > 0 {
				nav := addrver(Aov, -1)
				Aov = nav
			} else {
				opts.mnum -= 1
			}
			if opts.mnum < 0 { opts.mnum = maxmaze }
		case 'x':
			Ovwallpat = -1
			if Aov > 0 {
				nav := addrver(Aov, 1)
				Aov = nav
			} else {
				opts.mnum += 1
			}
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
		case 'p':
			nothing = nothing ^ NOFLOOR
			spau = fmt.Sprintf("no floors: %d\n",nothing & NOFLOOR)
		case 80:
			nothing = nothing ^ NOWALL
			spau = fmt.Sprintf("no walls: %d\n",nothing & NOWALL)
		case 84:
			nt := (nothing & 511) + 1
			nothing = (nothing & 1536) + (nt & 511)
			if anum > 0 { nothing = (nothing & 1536) + anum; anum = 0 }		// set lower 9 bits of no-thing mask [ but not walls or floors ]
			spau = fmt.Sprintf("no things: %d\n",nothing & 511)				// display no things mask
		case 's':
			opts.SP = !opts.SP
		case 76:
			opts.Nogtop = !opts.Nogtop
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
		case 'd':
			fmt.Printf("editor on, maze: %03d\n",opts.mnum+1)
			if opts.edat != 1 {
				opts.edat = 1
				stor_maz(opts.mnum+1)	// this does not auto store new edit mode to buffer save file, unless it creates the file
			}
		case 68:		// D
			fmt.Printf("editor off, maze: %03d\n",opts.mnum+1)
			if opts.edat != 0 {
				opts.edat = 0
				ed_sav(opts.mnum+1)		// this deactivates edit mode on this buffer
			}
		default:
			relodsub = false
		}

		if spau == "G¹ " {
			if opts.R14 { spau += "rv14" } else { spau += "rv1-9" }
		}
	if (relod || relodsub) {
		remaze(opts.mnum)
		uptitl(opts.mnum, spau)
	}
}

// pad for dialog page

func cpad(st string, d int) string {

	spout := st+"                                                                          " // jsut guess at a pad fill
	return string(spout[:d])
}

func menu_savit(y bool) {
	if y { ed_sav(opts.mnum+1) }
}

func menu_sav() {
	if opts.edat == 1 {
		dia := fmt.Sprintf("Save buffer for maze %d in .ed/g%dmaze%03d.ed ?",opts.mnum+1,opts.Gtp,opts.mnum+1)
		dialog.ShowConfirm("Saving",dia, menu_savit, w)
	} else { dialog.ShowInformation("Save Fail","edit mode is not active!",w) }
}

func menu_lodit(y bool) {
	fil := fmt.Sprintf(".ed/g%dmaze%03d.ed",opts.Gtp,opts.mnum+1)
	if y {
		Ovwallpat = -1
		lod_maz(fil)
		remaze(opts.mnum)
	}
}

func menu_lod() {
	if opts.edat == 1 {
		dia := fmt.Sprintf("Load buffer for maze %d from .ed/g%dmaze%03d.ed ?:",opts.mnum+1,opts.Gtp,opts.mnum+1)
		dialog.ShowConfirm("Loading",dia, menu_lodit, w)
	} else { dialog.ShowInformation("Load Fail","edit mode is not active!",w) }
}

func menu_rst(y bool) {
	if y {
		opts.edat = -1	// code to tell maze decompress not to load buffer file
		Ovwallpat = -1
		remaze(opts.mnum)
		opts.edat = 1
		//ed_sav(opts.mnum+1)	// reset does not overwrite file buffer, still need to save
	}
}

func menu_res() {
	if opts.edat == 1 {
		dia := fmt.Sprintf("Reset buffer for maze %d from G%d ROM ?\n - reset does not save to file",opts.mnum+1,opts.Gtp)
		dialog.ShowConfirm("Loading",dia, menu_rst, w)
	} else { dialog.ShowInformation("Reset Fail","edit mode is not active!",w) }
}

// init app and win

func aw_init() {

    a = app.New()
    w = a.NewWindow("G¹G²ved")

// quit menu option does not exit to term!
	menuItemExit := fyne.NewMenuItem("Exit", func() {
		os.Exit(0)
	})
	menuExit := fyne.NewMenu("Exit ", menuItemExit)

	menuItemSave := fyne.NewMenuItem("Save buffer <ctrl>-s", menu_sav)
	menuItemLoad := fyne.NewMenuItem("Load buffer <ctrl>-l", menu_lod)
	menuItemReset := fyne.NewMenuItem("Reset buffer <ctrl>-r", menu_res)
	menuItemEdhin := fyne.NewMenuItem("Edit hints", func() {
		dialog.ShowInformation("Edit hints", "Save - store buffer in file .ed/g{#}maze{###}.ed\n - where g# is 1 or 2 for g1/g2\n - and ### is the maze number e.g. 003\n\nLoad - overwrite current file contents this maze\n\nReset - reload buffer from rom read\n\ngved - G¹G² visual editor\ngithub.com/six-of-one/", w)
	})
	editMenu := fyne.NewMenu("Edit", menuItemSave, menuItemLoad, menuItemReset, menuItemEdhin)

	menuItemKeys := fyne.NewMenuItem("Keys ?", keyhints)
	menuItemAbout := fyne.NewMenuItem("About", func() {
		dialog.ShowInformation("About G¹G²ved", "Gauntlet / Gauntlet 2 visual editor\nAuthor: Six [a programmer]\n\ngithub.com/six-of-one/", w)
	})
	menuItemLIC := fyne.NewMenuItem("License", func() {
		dialog.ShowInformation("G¹G²ved License", "Gauntlet visual editor - gved\n\n(c) 2025 Six [a programmer]\n\nGPLv3.0\n\nhttps://www.gnu.org/licenses/gpl-3.0.html", w)
	})
	menuHelp := fyne.NewMenu("Help ", menuItemKeys, menuItemAbout, menuItemLIC)

	menuHint := fyne.NewMenu("cmds: ?, q, dD, fFgG, wWeE, rRt, hm, pPT, sL, il, u, v, A #a")

	mainMenu := fyne.NewMainMenu(menuExit, editMenu, menuHelp, menuHint)
	w.SetMainMenu(mainMenu)
	w.Canvas().SetOnTypedRune(typedRune)
	anum = 0
	shift = false

// get default win size

	if opts.Geow == 1024 && opts.Geoh == 1024 {		// defs set

		data, err := ioutil.ReadFile(".wstats")
		if err != nil {
			return
		}
		var geow float64
		var geoh float64
		fmt.Sscanf(string(data),"%v %v", &geow, &geoh)
		opts.Geow = math.Max(560,geow)
		opts.Geoh = math.Max(594,geoh)
	fmt.Printf("Load window size: %v x %v\n",geow,geoh)

	} else {
		file, err := os.Create(".wstats")
		if err == nil {
			wfs := fmt.Sprintf("%d %d",int(opts.Geow),int(opts.Geoh))
			file.WriteString(wfs)
			file.Close()
		}
	}
}

// test

type tappableIcon struct {
	widget.Icon
}

func newTappableIcon(res fyne.Resource) *tappableIcon {
	icon := &tappableIcon{}
	icon.ExtendBaseWidget(icon)
	icon.SetResource(res)

	return icon
}

func (t *tappableIcon) Tapped(e *fyne.PointEvent) {
	fmt.Printf("tapped - pos:%v, shf:%v, ctrl:%v\n",e.Position,shift,ctrl)

}
// update contents

func upwin(simg *image.NRGBA) {

// ration required by edit win y = x * ratio
	ratio := 1.0316529
//	bimg := canvas.NewRasterFromImage(simg)
//	w.Canvas().SetContent(bimg)
	geow := int(math.Max(560,opts.Geow))	// 556 is min, maze doesnt seem to fit or shrink smaller
	geoh := int(math.Max(594,opts.Geoh))	// 594 min
	if opts.edat > 0 { geoh = int(float64(geow) * ratio) }
	w.Resize(fyne.NewSize(float32(geow), float32(geoh)))

tres,err := fyne.LoadResourceFromPath("output.png")
if err == nil { w.SetContent(newTappableIcon(tres)) } else {
	fmt.Printf("Error on mouse clickable surface: %v",err)
}

	uptitl(opts.mnum, "")
}

// title special info update

func uptitl(mazeN int, spaux string) {

	til := fmt.Sprintf("G¹G²ved Maze: %d addr: %X edit: %d",mazeN + 1, slapsticMazeGetRealAddr(mazeN),opts.edat)
	if Aov > 0 { til = fmt.Sprintf("G¹G²ved Override addr: %X - %d",Aov,Aov) }
	if spaux != "" { til += " -- " + spaux }
	w.SetTitle(til)
}

// window resize control

func wizecon() {

	time.Sleep(3 * time.Second)		// some hang time to allow win to display & size, otherwise w x h is 1 x 1
	bgeow := int(opts.Geow)
	bgeoh := int(opts.Geoh)
	for {
// dont know why the +8, +36 needed, dont know if it will ever vary ??
		width := int(w.Content().Size().Width) + 8
		height := int(w.Content().Size().Height) + 36
//x					fmt.Printf("Window was resized! st: %d x %d n: %v x %v delta: %d, %d\n",bgeow,bgeoh,w.Content().Size().Width,w.Content().Size().Height,dw,dh)
		if width != bgeow || height != bgeoh {
				// window was resized
// provide live resize so other vis ops dont bounce it back
// for some reason maze updates resize the window down w -= 8 & h -= 36 to minimun
			opts.Geow = float64(width)
			opts.Geoh = float64(height)
// save stat
			file, err := os.Create(".wstats")
			if err == nil {
				wfs := fmt.Sprintf("%d %d",width,height)
				file.WriteString(wfs)
				file.Close()
//q	fmt.Printf("saving .wstats file\n")
			}
		}
		bgeow = int(opts.Geow)
		bgeoh = int(opts.Geoh)
		time.Sleep(2 * time.Second)
	}
}

func keyhints() {
	strp := cpad("single letter commands",36)
	strp += "\n–—–—–—–—–—–—–—–—–—–—–—"
//		strp += cpad("\n\n? - this list",52)
	strp += cpad("\nq - quit program",42)
	strp += cpad("\nd - editor mode",43)
	strp += cpad("\nf - floor pattern+",43)
	strp += cpad("\ng - floor color+",45)
	strp += cpad("\nw - wall pattern+",43)
	strp += cpad("\ne - wall color+",46)
	strp += cpad("\nr - rotate maze +90°",41)
	strp += cpad("\nR - rotate maze -90°",42)
	strp += cpad("\nh - mirror maze horizontal toggle",31)
	strp += "\nm - mirror maze vertical toggle"
	strp += cpad("\np - toggle floor invis",41)
	strp += cpad("\nP - toggle wall invis",42)
	strp += cpad("\nT - loop invis things",42)
	strp += cpad("\ns - toggle rnd special potion",34)
	strp += cpad("\nL - generator indicate letter",35)
	strp += cpad("\ni - gauntlet mazes r1 - r9",38)
	strp += cpad("\nl - use gauntlet rev 14",40)
	strp += cpad("\nu - gauntlet 2 mazes",39)
//		strp += cpad("\nv - valid address list",42)
	strp += "\nv - all maze addr (in termninal)"
	strp += cpad("\nA - toggle a override",41)
	strp += cpad("\n{n}a numeric of valid maze",35)
	strp += cpad("\n - load maze 1 - 127 g1",42)
	strp += cpad("\n - load maze 1 - 117 g2",42)
	strp += "\n - load address 229376 - 262143 "
//	strp += "\n * note some address will crash"
	strp += "\n–—–—–—–—–—–—–—–—–—–—–—"
	strb := fmt.Sprintf("\nG%d ",opts.Gtp)
	if G1 {
	if opts.R14 { strb += "(r14)"
		} else { strb += "(r1-9)" }}
	strp += cpad(strb,50)

	dialog.ShowInformation("Command Keys", strp, w)
}
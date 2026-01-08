package main

import (
	"image"
	"fmt"
	"math"
	"os"
	"io/ioutil"
	"time"
//		"image/color"

	"fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
)

// kontrol is for fyne window ops & input management

var w fyne.Window
var a fyne.App

// status keeper

var statup *fyne.Menu
var hintup *fyne.Menu
var mainMenu *fyne.MainMenu
var smod string

// input keys and keypress checks for canvas/ window
// since this is all that is called without other handler / timers
// - this is where maze update and edits will vector

// disable some cmd keys for edit mode click opts
var cmdoff bool

var anum int
var shift bool
var ctrl bool
var del bool
var logo bool		// other wise labeld "win" key
var edkey int		// for passing edit keys to clicker

func typedRune(r rune) {

// special aux string - put ops in title after maze #
	spau := ""
// relod
	relod := false
	relodsub := false

//	fmt.Printf("in keys event - %x\n",r)
// <ESC> also exits
	if r == 81 {
		needsav()
		os.Exit(0)
	}

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
			needsav()
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
			if key.Name == "Delete" { del = true }
			if key.Name == "BackSpace" { del = true }
			if key.Name == "LeftSuper" { logo = true }
			if key.Name == "LeftShift" { shift = true }
			if key.Name == "RightShift" { shift = true }
			if key.Name == "LeftControl" { ctrl = true }
			if key.Name == "RightControl" { ctrl = true }
        })
        deskCanvas.SetOnKeyUp(func(key *fyne.KeyEvent) {
//            fmt.Printf("Desktop key up: %v\n", key)
			if key.Name == "Escape" { os.Exit(0) }
			if key.Name == "Delete" { del = false }
			if key.Name == "BackSpace" { del = false }
			if key.Name == "LeftSuper" { logo = false }
			if key.Name == "LeftShift" { shift = false }
			if key.Name == "RightShift" { shift = false }
			if key.Name == "LeftControl" { ctrl = false }
			if key.Name == "RightControl" { ctrl = false }
			if key.Name == "S" && ctrl { menu_sav() }
			if key.Name == "L" && ctrl  { menu_lod() }
			if key.Name == "R" && ctrl  { menu_res() }
			if key.Name == "U" && ctrl  { uswap() }
			if key.Name == "Z" && ctrl  { undo() }
			if key.Name == "Y" && ctrl  { redo() }
       })
    }
	fmt.Printf("r %v shift %v\n",r,shift)
		edkey = int(r)

		cmdhin := "cmds: ?\\, Q, dD, fFgG, wWeE, rRt, hm, pPT, sL, S, il, u, v, A #a"

// keys that '\' doesnt block, no maze reloads
		relodsub = false
		switch r {
		case 92:		// \
			ska := "cmd keys mode"
			cmdoff = !cmdoff
// a,d only lower case not avail for edit hotkey
			if cmdoff && opts.edat > 0 {
				cmdhin = "cmds: ? '\\' - enable cmds, Q, #c C, A #a, dD, L, S"
				ska = "edit keys mode"
			}
			fmt.Printf("hint: %s\n", cmdhin)
			statlin(cmdhin,ska)
			opts.dntr = true
			relod = true
		case 63:		// ?
			keyhints()
		case 65:		// A
			if Aov > 0 { Aov = 0 } else {
				Aov = addrver(slapsticMazeGetRealAddr(opts.mnum), 0)
			}
		case 76:		// L
// with anum != 0, this becomes load s[1] buffer into ebuf, if in edit
			if opts.edat > 0 && anum > 0 && anum < sdmax {
fmt.Printf("L, anum: %05d, sdb: %d\n",anum, sdb)
				if sdb == -1 {
					fil := fmt.Sprintf(".ed/ebuf.ed")				// save ebuf for relod
					sav_maz(fil, ebuf, eflg, opts.DimX, opts.DimY)
				} else { needsav() }
				fil := fmt.Sprintf(".ed/sd%05d_g%d.ed",anum,opts.Gtp)
				cnd := lod_maz(fil, ebuf, false)
				if cnd >= 0 { sdb = anum; for y := 0; y < 11; y++ { eflg[y] =  tflg[y] }; ed_maze() }
				anum = 0
			} else { opts.Nogtop = !opts.Nogtop; relod = true }
		case 'd':
			if opts.Aob { dialog.ShowInformation("Edit mode", "Error: can not edit with border around maze!", w) } else {
				if opts.edat != 1 {
					smod = "Edit mode: "
					fmt.Printf("editor on, maze: %03d\n",opts.mnum+1)
					opts.edat = 1
					stor_maz(opts.mnum+1)	// this does not auto store new edit mode to buffer save file, unless it creates the file
					statlin(cmdhin,"on")
				}
				relod = true
			}
		case 68:		// D
			if opts.edat != 0 {
				smod = "View mode: "
				fmt.Printf("editor off, maze: %03d\n",opts.mnum+1)
				opts.edat = 0
				ed_sav(opts.mnum+1)		// this deactivates edit mode on this buffer
				statlin(cmdhin,"on")
				relod = true
			}
		case 'c':
			if anum > 0 && anum < 65 {
				cycl = anum - 1
				anum = 0
			} else { cycl--; cycl-- }
			fallthrough
		case 67:		// C
				cycl++
				fmt.Printf("cyc %d \n",cycl)
				if cycl > 64 { cycl = 0 }
				if G1 {
					g1edit_keymap[cycloc] = cycl
				} else {
					g2edit_keymap[cycloc] = cycl
				}
				stu := fmt.Sprintf("cyc: %d",cycl)
				statlin(cmdhin,stu)
				edkey = 99						// pre set store cycl when cycling
				relod = true					// needed to refresh indicate text
				opts.dntr = true				// ... but dont kill the ebuf
		case 83:		// S
// have anum !=0, save that buffer
				if anum > 0 && anum < sdmax && opts.edat > 0 {
					fil := fmt.Sprintf(".ed/sd%05d_g%d.ed",anum,opts.Gtp)
					sav_maz(fil, ebuf, eflg, opts.DimX, opts.DimY)
					anum = 0
				} else {
// with no anum, rotate curr ebuf thru s[1] - s[27], store eb in s[0]
					cnd := -1
					ldb := sdb
					if sdb == -1 {
						fil := fmt.Sprintf(".ed/ebuf.ed")				// save ebuf for relod
						sav_maz(fil, ebuf, eflg, opts.DimX, opts.DimY)
					}
					for cnd < 0 && ldb < sdmax {
						ldb++
						fil := fmt.Sprintf(".ed/sd%05d_g%d.ed",ldb,opts.Gtp)
						cnd = lod_maz(fil, ebuf, false)
						if cnd >= 0 { sdb = ldb; for y := 0; y < 11; y++ { eflg[y] =  tflg[y] }; ed_maze() }
					}
					if cnd < 0 {
						sdb = -1
						if opts.edat > 0 {
							fil := fmt.Sprintf(".ed/ebuf.ed")			// cycle back out
							cnd = lod_maz(fil, ebuf, true)
						if cnd >= 0 { for y := 0; y < 11; y++ { eflg[y] =  tflg[y] } }
						} else { remaze(opts.mnum) }
					}
				}
		default:
			relodsub = false
		}
// view cmd keys - also on edit, but blockable
	  if !cmdoff || opts.edat < 1 {
		relodsub = true
		switch r {
		case 'z':
			needsav()
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
			needsav()
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
			opts.bufdrt = true
			opts.dntr = true
		case 87:		// W
			Ovwallpat -= 1
			if Ovwallpat < 0 { Ovwallpat = 7 }
			spau = fmt.Sprintf("cmd: w - wallp: %d\n",Ovwallpat)
			opts.bufdrt = true
			opts.dntr = true
		case 'e':
			Ovwallcol += 1
			if anum > 0 { Ovwallcol = anum - 1; anum = 0 }
			if Ovwallcol > 16 { Ovwallcol = 0 }
			spau = fmt.Sprintf("cmd: e - wallc: %d\n",Ovwallcol)
			opts.bufdrt = true
			opts.dntr = true
		case 69:		// E
			Ovwallcol -= 1
			if Ovwallcol < 0 { Ovwallcol = 16 }
			spau = fmt.Sprintf("cmd: e - wallc: %d\n",Ovwallcol)
			opts.bufdrt = true
			opts.dntr = true
		case 'f':
			Ovflorpat += 1
			if anum > 0 { Ovflorpat = anum - 1; anum = 0 }
			if Ovflorpat > 8 { Ovflorpat = 0 }
			spau = fmt.Sprintf("cmd: f - floorp: %d\n",Ovflorpat)
			opts.bufdrt = true
			opts.dntr = true
		case 70:		// F
			Ovflorpat -= 1
			if Ovflorpat < 0 { Ovflorpat = 8 }
			spau = fmt.Sprintf("cmd: f - floorp: %d\n",Ovflorpat)
			opts.bufdrt = true
			opts.dntr = true
		case 'g':
			Ovflorcol += 1
			if anum > 0 { Ovflorcol = anum - 1; anum = 0 }
			if Ovflorcol > 15 { Ovflorcol = 0 }
			spau = fmt.Sprintf("cmd: g - floorc: %d\n",Ovflorcol)
			opts.bufdrt = true
			opts.dntr = true
		case 71:		// G
			Ovflorcol -= 1
			if Ovflorcol < 0 { Ovflorcol = 15 }
			spau = fmt.Sprintf("cmd: g - floorc: %d\n",Ovflorcol)
			opts.bufdrt = true
			opts.dntr = true
		case 'r':
			opts.MRP = true
			opts.MRM = false
			spau = fmt.Sprintf("cmd: r - mr+: %t mr-: %t\n",opts.MRP,opts.MRM)
//			opts.dntr = true
		case 82:		// R
			opts.MRP = false
			opts.MRM = true
			spau = fmt.Sprintf("cmd: r - mr+: %t mr-: %t\n",opts.MRP,opts.MRM)
//			opts.dntr = true
		case 't':
			opts.MRP = false
			opts.MRM = false
			spau = fmt.Sprintf("cmd: t - mr+: %t mr-: %t\n",opts.MRP,opts.MRM)
//			opts.dntr = true
		case 'm':
			opts.MV = !opts.MV
			spau = fmt.Sprintf("cmd: m - mv: %t\n",opts.MV)
//			opts.dntr = true
		case 'h':
			opts.MH = !opts.MH
			spau = fmt.Sprintf("cmd: h - mh: %t\n",opts.MH)
//			opts.dntr = true
		case 'i':
			opts.Gtp = 1
			opts.R14 = false
			G1 = true
			G2 = false
			maxmaze = 126
			spau = "G¹ "
			needsav()
		case 'l':
			opts.Gtp = 1
			opts.R14 = !opts.R14
			G1 = true
			G2 = false
			maxmaze = 126
			spau = "G¹ "
			needsav()
		case 'p':
			nothing = nothing ^ NOFLOOR
			spau = fmt.Sprintf("no floors: %d\n",nothing & NOFLOOR)
			opts.dntr = true
		case 80:		// P
			nothing = nothing ^ NOWALL
			spau = fmt.Sprintf("no walls: %d\n",nothing & NOWALL)
			opts.dntr = true
		case 84:		// T
			nt := (nothing & 511) + 1
			nothing = (nothing & 1536) + (nt & 511)
			if anum > 0 { nothing = (nothing & 1536) + anum; anum = 0 }		// set lower 9 bits of no-thing mask [ but not walls or floors ]
			spau = fmt.Sprintf("no things: %d\n",nothing & 511)				// display no things mask
			opts.dntr = true
		case 's':
			opts.SP = !opts.SP
			opts.dntr = true
		case 'u':
			opts.Gtp = 2
			G1 = false
			G2 = true
			maxmaze = 116
			spau = "G² mazes"
			needsav()
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
			opts.dntr = true
		default:
			relodsub = false
		}
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

func menu_disc(y bool) {
	if y {
		opts.bufdrt = false
	}
}

func dumpbuf() {
	if opts.bufdrt {
		dia := fmt.Sprintf("Unsaved changes for previous maze\nDiscard them to load new maze?\n\n(rejecting discard exits load)")
		dialog.ShowConfirm("Discard?",dia, menu_disc, w)
	}
}

func menu_savit(y bool) {
	if y {
		if sdb < 0 {
			ed_sav(opts.mnum+1)
		} else {
			fil := fmt.Sprintf(".ed/sd%05d_g%d.ed",sdb,opts.Gtp)
			sav_maz(fil, ebuf, eflg, opts.DimX, opts.DimY)
		}
	}
}

func needsav() {
	if opts.bufdrt {
		dia := fmt.Sprintf("Unsaved changes for maze %d in .ed/g%dmaze%03d.ed ?\n\nWARNING:\nif not saved, changes will be discarded",opts.mnum+1,opts.Gtp,opts.mnum+1)
		if sdb >= 0 { dia = fmt.Sprintf("Unsaved changes in buffer .ed/sd%05d_g%d.ed\n\nWARNING:\nif not saved, changes will be discarded",sdb,opts.Gtp) }
		dialog.ShowConfirm("Save?",dia, menu_savit, w)
	}
}

func menu_sav() {
	if opts.edat == 1 {
		dia := fmt.Sprintf("Save buffer for maze %d in .ed/g%dmaze%03d.ed ?",opts.mnum+1,opts.Gtp,opts.mnum+1)
		if sdb >= 0 { dia = fmt.Sprintf("Save store buffer .ed/sd%05d_g%d.ed",sdb,opts.Gtp) }
		dialog.ShowConfirm("Saving",dia, menu_savit, w)
	} else { dialog.ShowInformation("Save Fail","edit mode is not active!",w) }
}

func menu_lodit(y bool) {
	fil := fmt.Sprintf(".ed/g%dmaze%03d.ed",opts.Gtp,opts.mnum+1)
	if y {
		Ovwallpat = -1
		cnd := lod_maz(fil, ebuf, true)
		sdb = -1
		if cnd >= 0 { for y := 0; y < 11; y++ { eflg[y] =  tflg[y] } }
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

func undo() {
	if delstak > 0 {
		delstak--
		sw := ebuf[xy{delbuf.mx[delstak], delbuf.my[delstak]}]
		ebuf[xy{delbuf.mx[delstak], delbuf.my[delstak]}] = delbuf.elem[delstak]
		delbuf.elem[delstak] = sw
		opts.bufdrt = true
		ed_maze()
	}
}

func redo() {
	if delbuf.elem[delstak] != -1 {
		sw := ebuf[xy{delbuf.mx[delstak], delbuf.my[delstak]}]
		ebuf[xy{delbuf.mx[delstak], delbuf.my[delstak]}] = delbuf.elem[delstak]
		delbuf.elem[delstak] = sw
		delstak++
		opts.bufdrt = true
		ed_maze()
	}
//	ed_maze()
}

func uswap() {
	for y := 0; y <= opts.DimX; y++ {
		for x := 0; x <= opts.DimY; x++ {
			sw := ebuf[xy{x,y}]
			ebuf[xy{x,y}] = ubuf[xy{x,y}]
			ubuf[xy{x,y}] = sw
	}}
	ed_maze()
}

// set menus

func st_menu() {
// quit menu option does not exit to term!
	menuItemExit := fyne.NewMenuItem("Exit", func() {
		needsav()
		os.Exit(0)
	})
	menuExit := fyne.NewMenu("Exit ", menuItemExit)

	menuItemSave := fyne.NewMenuItem("Save buffer <ctrl>-s", menu_sav)
	menuItemLoad := fyne.NewMenuItem("Load buffer <ctrl>-l", menu_lod)
	menuItemReset := fyne.NewMenuItem("Reset buffer <ctrl>-r", menu_res)
	menuItemEdhin := fyne.NewMenuItem("Edit hints", func() {
		dialog.ShowInformation("Edit hints", "Save - store buffer in file .ed/g{#}maze{###}.ed\n - where g# is 1 or 2 for g1/g2\n - and ### is the maze number e.g. 003\n"+
			"\nLoad - overwrite current file contents this maze\n\nReset - reload buffer from rom read\n\nedit keys:\nd: turn editor on, init maze store in .ed/\n"+
			"D: turn editor off, saves edits to file\ndel, backspace - set floor *\nC: cycle edit item #++, c: cycle item #-- *\n#c enter number {1-64}c, all set place item *\n"+
			"b - horiz door, B - vert door, w - wall *\nk - key, t - transporter *\n"+
			"edit keys lock when pressed, hit 'b' and place doors\nLogo/ Super key - click to reassing current key\n"+
			"* most edit keys require '\\' mode\n\n\ngved - G¹G² visual editor\ngithub.com/six-of-one/", w)
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

	hintup = fyne.NewMenu("cmds: ?\\, Q, dD, fFgG, wWeE, rRt, hm, pPT, sL, S, il, u, v, A #a")

	statup = fyne.NewMenu("view mode:")

	mainMenu = fyne.NewMainMenu(menuExit, editMenu, menuHelp, hintup, statup)
	w.SetMainMenu(mainMenu)
}

// init app and win

func aw_init() {

    a = app.New()
    w = a.NewWindow("G¹G²ved")

	st_menu()
	w.Canvas().SetOnTypedRune(typedRune)
	anum = 0
// ed stuff, consider moving
	delstak = 0
	sdb = -1
	cycl = 0
	edmaze = mazeDecompress(slapsticReadMaze(1), false)

// get default win size

	if opts.Geow == 1024 && opts.Geoh == 1050 {		// defs set

		data, err := ioutil.ReadFile(".wstats")
		if err != nil {
			return
		}
		var geow float64
		var geoh float64
		fmt.Sscanf(string(data),"%v %v", &geow, &geoh)
		opts.Geow = math.Max(560,geow)
		opts.Geoh = math.Max(586,geoh)
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

// update stat line

func statlin(hs string,ss string) {

	hintup.Label = hs
	statup.Label = smod + ss
	mainMenu.Refresh()
}

// click area for edits

// button we can detect click and release areas for rubberband area & fills

type holdableButton struct {
    widget.Button
}

func newHoldableButton() *holdableButton {

    button := &holdableButton{}
    button.ExtendBaseWidget(button)
	return button
}

// store x & y when mouse button goes down - to start rubberband area
var sxmd float64
var symd float64

// &{{{387 545} {379 509.92188}} 4 0}

func (h *holdableButton) MouseDown(mm *desktop.MouseEvent){
	ax := 0.0	// absolute x & y
	ay := 0.0
	mb := 0		// mb 1 = left, 2 = right, 4 = middle
	mk := 0		// mod key 1 = sh, 2 = ctrl, 4 = alt, 8 = logo
	pos := fmt.Sprintf("%v",mm)
	fmt.Sscanf(pos,"&{{{%f %f} {%f %f}} %d %d",&ax,&ay,&sxmd,&symd,&mb,&mk)
	fmt.Printf("%d down: %.2f x %.2f \n",mb,sxmd,symd)
}


var repl int		// replace will be by ctrl-h in select area or entire maze, by match
var cycl int		// cyclical set - C cycles, c sets - using c loc in keymap
var cycloc = 99

// edkey 'locks' on when pressed

func (h *holdableButton) MouseUp(mm *desktop.MouseEvent){
 //   fmt.Printf("up %v\n",mm)
	if opts.edat > 0 {
		ax := 0.0	// absolute x & y
		ay := 0.0
		ix := 0.0	// rel x & y interm float32
		iy := 0.0
		mb := 0		// mb 1 = left, 2 = right, 4 = middle
		mk := 0		// mod key 1 = sh, 2 = ctrl, 4 = alt, 8 = logo
		pos := fmt.Sprintf("%v",mm)
		fmt.Sscanf(pos,"&{{{%f %f} {%f %f}} %d %d",&ax,&ay,&ix,&iy,&mb,&mk)
		fmt.Printf("%d up: %.2f x %.2f \n",mb,ix,iy)
		ex := int(ix / opts.dtec)
		ey := int(iy / opts.dtec)
		sx := int(sxmd / opts.dtec)
		sy := int(symd / opts.dtec)
		var setcode int			// code to store given edit hotkey
		if G1 {
			setcode = g1edit_keymap[edkey]
		} else {
			setcode = g2edit_keymap[edkey]
		}

// no access, keys: ? Q, A #a, dD, L, S
		fmt.Printf(" dtec: %f maze: %d x %d - element:%d\n",opts.dtec,ex,ey,ebuf[xy{ex, ey}])
		if mb == 4 && cmdoff {		// middle mb, do a reassign
			if G1 {
				g1edit_keymap[edkey] = ebuf[xy{ex, ey}]
			} else {
				g2edit_keymap[edkey] = ebuf[xy{ex, ey}]
			}
		} else {
		if del || cmdoff {
			if ex < sx { t := ex; ex = sx; sx = t }		// swap if end smaller than start
			if ey < sy { t := ey; ey = sy; sy = t }
		 for my := sy; my <= ey; my++ {
		   for mx := sx; mx <= ex; mx++ {
// looped now
			if del { undo_buf(mx, my); ebuf[xy{mx, my}] = 0; opts.bufdrt = true } else {	// delete anything for now makes a floor
			if setcode > 0 { undo_buf(mx, my); ebuf[xy{mx, my}] = setcode; opts.bufdrt = true }
			if edkey == 182 { ebuf[xy{mx, my}] = repl; opts.bufdrt = true }		//
			if edkey == 214 { repl = ebuf[xy{mx, my}] }		// just placeholder until new repl done
			}
		  }}
//			fmt.Printf(" chg elem: %d maze: %d x %d\n",ebuf[xy{mx, my}],mx,my)
		}}
		ed_maze()
	}

}

// update contents

func upwin(simg *image.NRGBA) {

//								                 ┌» un-borded maze is 528 x 528 for a 33 x 33 cell maze
	geow := int(math.Max(560,opts.Geow))	// 560 is min, maze doesnt seem to fit or shrink smaller
	geoh := int(math.Max(586,opts.Geoh))	// 586 min
	if opts.edat > 0 {
//		geow = geow & 0xfe0	+ 13			// lock to multiples of 32
		ngeoh := geow + 26					// square maze area + 26 for menu bar - window is still 4 wider than maze content
		if ngeoh != geoh { dialog.ShowInformation("Edit mode","set window ratio to edit",w) }
		geoh = ngeoh
	}
	opts.dtec = 16.0 * (float64(geow - 4) / 528.0)				// the size of a tile, odd window size may cause issues
	fmt.Printf(" dtec: %f\n",opts.dtec)
	bimg := canvas.NewRasterFromImage(simg)
	w.Resize(fyne.NewSize(float32(geow), float32(geoh)))

// turns display into clickable edit area
	btn := newHoldableButton()
	box := container.NewPadded(btn, bimg)		// key to seeing maze & having the click button will full mouse sense
	w.SetContent(box)
	fmt.Printf("btn sz %v\n",btn.Size())

	uptitl(opts.mnum, "")
}

// title special info update

func uptitl(mazeN int, spaux string) {

	til := fmt.Sprintf("G¹G²ved Maze: %d addr: %X",mazeN + 1, slapsticMazeGetRealAddr(mazeN))
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
	strp += cpad("\nQ - quit program",42)
	strp += cpad("\n\\ - toggle cmd keys",41)
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
	strp += cpad("\n{n}S save curr to buffer #",35)
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
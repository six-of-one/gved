package main

import (
	"fmt"
	"math"
	"os"
	"io/ioutil"
	"time"
	"image"
//		"image/color"

	"fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

// fyne menu & window system isolated from keyboard & control now

var w fyne.Window
var a fyne.App

// status keeper - appears in spare menu item on mbar for now

var statup *fyne.Menu
var hintup *fyne.Menu
var mainMenu *fyne.MainMenu
var smod string

// control ops called from menus

func menu_savit(y bool) {
	if y {
		if sdb < 0 {
			ed_sav(opts.mnum+1)
		} else {
			fil := fmt.Sprintf(".ed/sd%05d_g%d.ed",sdb,opts.Gtp)
			sav_maz(fil, ebuf, eflg, opts.DimX, opts.DimY, 0 - sdb)
		}
	}
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

func menu_rst(y bool) {
	if y {
		opts.edat = -1	// code to tell maze decompress not to load buffer file
		Ovwallpat = -1
		remaze(opts.mnum)
		opts.edat = 1
		//ed_sav(opts.mnum+1)	// reset does not overwrite file buffer, still need to save
	}
}

func menu_sav() {
	if opts.edat == 1 {
		dia := fmt.Sprintf("Save buffer for maze %d in .ed/g%dmaze%03d.ed ?",opts.mnum+1,opts.Gtp,opts.mnum+1)
		if sdb >= 0 { dia = fmt.Sprintf("Save buffer %d to .ed/sd%05d_g%d.ed",sdb,sdb,opts.Gtp) }
		dialog.ShowConfirm("Saving",dia, menu_savit, w)
	} else { dialog.ShowInformation("Save Fail","edit mode is not active!",w) }
}


func menu_lod() {
	if opts.edat == 1 {
		dia := fmt.Sprintf("Load buffer for maze %d from .ed/g%dmaze%03d.ed ?:",opts.mnum+1,opts.Gtp,opts.mnum+1)
		dialog.ShowConfirm("Loading",dia, menu_lodit, w)
	} else { dialog.ShowInformation("Load Fail","edit mode is not active!",w) }
}

func menu_res() {
	if opts.edat == 1 {
		dia := fmt.Sprintf("Reset buffer for maze %d from G%d ROM ?\n - reset does not save to file",opts.mnum+1,opts.Gtp)
		dialog.ShowConfirm("Loading",dia, menu_rst, w)
	} else { dialog.ShowInformation("Reset Fail","edit mode is not active!",w) }
}

// save as
func menu_savas() {

	fileDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
        if err != nil {
            fmt.Println("Save as Error:", err)
            return
        }
        if writer == nil {
            fmt.Println("No file selected")
            return
        }

        fmt.Println("Selected:", writer.URI().Path())
		fil := writer.URI().Path()

		mazn := opts.mnum+1
		if anum > 0 { mazn = anum }
		sav_maz(fil, ebuf, eflg, opts.DimX, opts.DimY, mazn)

    }, w)
	fileDialog.Show()
	fileDialog.Resize(fyne.NewSize(float32(opts.Geow - 10), float32(opts.Geoh - 30)))
}

// load maze file
func menu_laodf() {

	fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
        if err != nil {
            fmt.Println("Save as Error:", err)
            return
        }
        if reader == nil {
            fmt.Println("No file selected")
            return
        }

        fmt.Println("Selected:", reader.URI().Path())
		fil := reader.URI().Path()

		if opts.bufdrt { menu_savit(true) }		// autosave
		Ovwallpat = -1
		cnd := lod_maz(fil, ebuf, true)
		sdb = -1
		if cnd >= 0 { for y := 0; y < 11; y++ { eflg[y] =  tflg[y] } }
		remaze(opts.mnum)
    }, w)
	fileDialog.Show()
	fileDialog.Resize(fyne.NewSize(float32(opts.Geow - 10), float32(opts.Geoh - 30)))
}

// insert blank maze into buffer
// pr true == preserve decor, walls & floors, exit and start
// called with anum set, preserve items in reverse of hide items #T

func menu_blank(pr bool) {
	if opts.bufdrt { menu_savit(true) }		// autosave
	if !pr {
		eflg[4] = eflg[4] & 0xcf			// turn off H & V
		eflg[5] = 0							// default floor & wall
		eflg[6] = 0
	}
	for ty := 0; ty <= opts.DimY; ty++ {
	for tx := 0; tx <= opts.DimX; tx++ {
		clr := true
		if pr {
			if G1 {
				if ebuf[xy{tx, ty}] == G1OBJ_EXIT { clr = false }
				if ebuf[xy{tx, ty}] == G1OBJ_EXIT4 { clr = false }
				if ebuf[xy{tx, ty}] == G1OBJ_EXIT8 { clr = false }
				if ebuf[xy{tx, ty}] == G1OBJ_PLAYERSTART { clr = false }
			} else { // G2
				if ebuf[xy{tx, ty}] == MAZEOBJ_EXIT { clr = false }
				if ebuf[xy{tx, ty}] == MAZEOBJ_EXITTO6 { clr = false }
				if ebuf[xy{tx, ty}] == MAZEOBJ_PLAYERSTART { clr = false }
			}
		}
// anum as item hide flags, but keep those elements
		flg := anum & 4095			// item mask from #T
		if G1 {
			if g1mask[ebuf[xy{tx, ty}]] & flg > 0 { clr = false }
//fmt.Printf(" flg %d elem: %d test: %d\n",flg,g1mask[ebuf[xy{tx, ty}]],g1mask[ebuf[xy{tx, ty}]] & flg)
		} else {
			if g2mask[ebuf[xy{tx, ty}]] & flg > 0 { clr = false }
//fmt.Printf(" flg %d elem: %d test: %d\n",flg,g2mask[ebuf[xy{tx, ty}]],g2mask[ebuf[xy{tx, ty}]] & flg)
		}
		if clr {
			ebuf[xy{tx, ty}] = 0
			if tx == 0 || tx == opts.DimX { ebuf[xy{tx, ty}] = MAZEOBJ_WALL_REGULAR }
			if ty == 0 || ty == opts.DimY { ebuf[xy{tx, ty}] = MAZEOBJ_WALL_REGULAR }
		}
	}}
	pr = false
	opts.dntr = true
	remaze(opts.mnum)
}

func menu_copy() { if opts.edat > 0 { ccp_tog(COPY); if ccp > 0 { smod = "Edit COPY: "}; statlin(cmdhin,"") }}
func menu_cut() { if opts.edat > 0 { ccp_tog(CUT); if ccp > 0 { smod = "Edit CUT: "}; statlin(cmdhin,"") }}
func menu_paste() { if opts.edat > 0 { ccp_tog(PASTE); if ccp > 0 { smod = "Edit PASTE: "}; statlin(cmdhin,"") }}

// set menus

func st_menu() {
// default 'quit' menu option does not call needsav !
	menuItemExit := fyne.NewMenuItem("Exit", func() {
		exitsel = true
		needsav()
	})
	menuItemLodf := fyne.NewMenuItem("Load maze from <shift-ctrl>-l",menu_laodf)
	menuItemSava := fyne.NewMenuItem("Save maze as <shift-ctrl>-s",menu_savas)
	menuItemBlan := fyne.NewMenuItem("Blank maze",func() { menu_blank(false) })
	menuItemBlnK := fyne.NewMenuItem("Blank maze, keep decor",func() { menu_blank(true) })
	menuItemLin1 := fyne.NewMenuItem("═══════════════",nil)
	menuFile := fyne.NewMenu("File", menuItemLodf, menuItemSava, menuItemBlan, menuItemBlnK, menuItemLin1, menuItemExit)

	menuItemSave := fyne.NewMenuItem("Save buffer <ctrl>-s", menu_sav)
	menuItemLoad := fyne.NewMenuItem("Load buffer <ctrl>-l", menu_lod)
	menuItemReset := fyne.NewMenuItem("Reset buffer <ctrl>-r", menu_res)
	menuItemLin2 := fyne.NewMenuItem("═══════════════",nil)
	menuItemCopy := fyne.NewMenuItem("Copy <ctrl>-c", menu_copy)
	menuItemCut := fyne.NewMenuItem("Cut <ctrl>-x", menu_cut)
	menuItemPaste := fyne.NewMenuItem("Paste <ctrl>-p", menu_paste)
	menuItemUndo := fyne.NewMenuItem("Undo <ctrl>-z", undo)
	menuItemRedo := fyne.NewMenuItem("Redo <ctrl>-y", redo)
	menuItemUswp := fyne.NewMenuItem("Ult buf <ctrl>-u", uswap)
	menuItemEdhin := fyne.NewMenuItem("Edit hints", func() {
		strp := ""
		if opts.edat == 1 {
			strp = "Edit mode: "
			if cmdoff { strp += "edit keys" } else { strp += "cmd keys" }
		} else {
			strp = "View mode: cmd keys only"
		}
		dialog.ShowInformation("Edit hints", strp+"\n══════════════════════════════\nSave - store buffer in file .ed/g{#}maze{###}.ed\n - where g# is 1 or 2 for g1/g2\n - and ### is the maze number e.g. 003\n"+
			"\nLoad - overwrite current file contents this maze\n\nReset - reload buffer from rom read\n\nedit keys:\ne: turn editor on, init maze store in .ed/\n"+
			"E: turn editor off, check unsaved buf\ndel, backspace - set floor *\nC: cycle edit item #++, c: cycle item #-- *\n#c enter number {1-64}c, all set place item *\n"+
			"H: toggle horiz wrap, V: toggle vert wrap\n"+
			"d - horiz door, D - vert door, w, W - walls *\nf, F - foods, k - key, t - treasure *\np, P - potions, T - teleporter\n"+
			"edit keys lock when pressed, hit 'b' and place doors\nmiddle click - click to reassign current key\n"+
			"* most edit keys require '\\' mode\n\n\ngved - G¹G² visual editor\ngithub.com/six-of-one/", w)
	})
	editMenu := fyne.NewMenu("Edit", menuItemSave, menuItemLoad, menuItemReset, menuItemLin2, menuItemCopy, menuItemCut, menuItemPaste, menuItemUndo, menuItemRedo, menuItemUswp, menuItemEdhin)

	menuItemKeys := fyne.NewMenuItem("Keys ?", keyhints)
	menuItemAbout := fyne.NewMenuItem("About", func() {
		dialog.ShowInformation("About G¹G²ved", "Gauntlet / Gauntlet 2 visual editor\nAuthor: Six [a programmer]\n\ngithub.com/six-of-one/", w)
	})
	menuItemLIC := fyne.NewMenuItem("License", func() {
		dialog.ShowInformation("G¹G²ved License", "Gauntlet visual editor - gved\n\n(c) 2025 Six [a programmer]\n\nGPLv3.0\n\nhttps://www.gnu.org/licenses/gpl-3.0.html", w)
	})
	menuHelp := fyne.NewMenu("Help ", menuItemKeys, menuItemAbout, menuItemLIC)

	hintup = fyne.NewMenu("cmds: ?, eE, fFgG, wWqQ, rRt, hm, pPT, sL, S, il, u, v, A #a")

	statup = fyne.NewMenu("view mode:")

	mainMenu = fyne.NewMainMenu(menuFile, editMenu, menuHelp, hintup, statup)
	w.SetMainMenu(mainMenu)
}

// init app and main win

func aw_init() {

    a = app.NewWithID("0777")
    w = a.NewWindow("G¹G²ved")
	w.SetCloseIntercept(func() {
		if wpbop { wpb.Close() }
		if wpalop { wpal.Close() }
	})

	st_menu()
	w.Canvas().SetOnTypedRune(typedRune)
	anum = 0
// ed stuff, consider moving
	ccp = NOP
	wpbop = false
	wpalop = false
	sdb = -1
	cycl = 0
	edmaze = mazeDecompress(slapsticReadMaze(1), false)
	cmdhin = "cmds: ?, eE, fFgG, wWqQ, rRt, hm, pPT, sL, S, il, u, v, A #a"
	delbset(0)
	restak = 0
	specialKey()

// get default win size

	if opts.Geow == 1024 && opts.Geoh == 1050 {		// defs set

		data, err := ioutil.ReadFile(".wstats")
		if err == nil {
			var geow float64
			var geoh float64
			fmt.Sscanf(string(data),"%v %v", &geow, &geoh)
			opts.Geow = math.Max(560,geow)
			opts.Geoh = math.Max(586,geoh)
	fmt.Printf("Load window size: %v x %v\n",geow,geoh)
		}

	} else {
		file, err := os.Create(".wstats")
		if err == nil {
			wfs := fmt.Sprintf("%d %d",int(opts.Geow),int(opts.Geoh))
			file.WriteString(wfs)
			file.Close()
		}
	}
	get_pbcnt()

}

// sub win switch G¹ / G²

func subsw() {
	if wpbop { pbmas_cyc(0) }
	if wpalop { palete() }
}

// make clickable image wimg in window cw with given size

func clikwin(cw fyne.Window, wimg *image.NRGBA, wx int, wy int) {

	bimg := canvas.NewRasterFromImage(wimg)
	cw.Resize(fyne.NewSize(float32(wx), float32(wy)))

// turns display into clickable edit area
	btn := newHoldableButton()
	btn.title = cw.Title()
	box := container.NewPadded(btn, bimg)		// key to seeing maze & having the click button will full mouse sense
	cw.SetContent(box)
fmt.Printf("btn sz %v\n",btn.Size())

}
// update contents of main edit window, includes title

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
	clikwin(w, simg, geow, geoh)

	spx := ""
	if sdb > 0 { spx = fmt.Sprintf("sdbuf: %d",sdb) }
	if anum != 0 { spx += fmt.Sprintf("| numeric: %d", anum) }
	uptitl(opts.mnum, spx)
}

// title special info update

func uptitl(mazeN int, spaux string) {

	til := fmt.Sprintf("G¹G²ved Maze: %d addr: %X",mazeN + 1, slapsticMazeGetRealAddr(mazeN))
	if Aov > 0 { til = fmt.Sprintf("G¹G²ved Override addr: %X - %d",Aov,Aov) }
	if spaux != "" { til += " -- " + spaux }
	w.SetTitle(til)
}

// update stat line

func statlin(hs string,ss string) {

	hintup.Label = hs
	statup.Label = smod + ss
	mainMenu.Refresh()
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

// pad for dialog page

func cpad(st string, d int) string {

	spout := st+"                                                                          " // jsut guess at a pad fill
	return string(spout[:d])
}

// dialog called from kby or menu

func keyhints() {
	strp := ""
	kys := 1
	dk := "Command Keys"
	if opts.edat == 1 {
		strp = "Edit mode: "
		if cmdoff { strp += "edit keys"; kys = 2; dk = "Editor Keys" } else { strp += "cmd keys" }
	} else {
		strp = "View mode: cmd keys only"
	}
	strp += cpad("\nsingle letter commands",36)
	strp += "\n–—–—–—–—–—–—–—–—–—–—–—"
//		strp += cpad("\n\n? - this list",52)
	strp += cpad("\nctrl-q - quit program",40)
	if opts.edat == 1 {
		strp += cpad("\nESC> exit editor ╗",40)
		strp += cpad("\n\\ - toggle cmd keys*",40)
	} else { strp += cpad("\nESC> editor mode",40) }
	if kys == 1 {
		strp += cpad("\nf - floor pattern+",44)
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
	} else {
		strp += cpad("\nH - toggle horiz maze wrap",35)
		strp += cpad("\nV - toggle vert maze wrap",36)
		strp += cpad("\nC - cycle item++, key c",40)
		strp += cpad("\nc - cycle item--, key c",43)
		strp += cpad("\n{n}c set item 1 - 64, key c",39)
		strp += cpad("\nL - generator indicate letter",35)
		strp += cpad("\nS - cycle sd buffers",42)
		strp += cpad("\n{n}S save curr to buffer #",35)
	}
	if kys == 1 {
		strp += cpad("\nL - generator indicate letter",35)
		strp += cpad("\n{n}S save curr to buffer #",35)
		strp += cpad("\ni - gauntlet mazes r1 - r9",38)
		strp += cpad("\nl - use gauntlet rev 14",40)
		strp += cpad("\nu - gauntlet 2 mazes",39)
//		strp += cpad("\nv - valid address list",42)
		strp += "\nv - all maze addr (in termninal)"
	}
	strp += cpad("\nA - toggle a override",41)
	strp += cpad("\n{n}a numeric of valid maze",35)
	strp += cpad("\n - load maze 1 - 127 g1",42)
	strp += cpad("\n - load maze 1 - 117 g2",42)
	strp += "\n - load address 229376 - 262143 "
	strp += "\n–—–—–—–—–—–—–—–—–—–—–—"
	strb := fmt.Sprintf("\nG%d ",opts.Gtp)
	if G1 {
	if opts.R14 { strb += "(r14)"
		} else { strb += "(r1-9)" }}
	strp += cpad(strb,50)
	if kys == 2 {
		strp += "\n\ntypical: key selects item,\n L-click place, M-click assign"
		strp += "\n–—–—–—–—–—–—–—"
		strp += cpad("\n<DEL> (hold down) set floor",32)
		strp += cpad("\nw - standard walls",42)
		strp += cpad("\nW - shootable walls",41)
		strp += cpad("\nd - horizontal door",41)
		strp += cpad("\nD - vertical door",44)
		strp += cpad("\nf - shootable food",42)
		strp += cpad("\nF - indestructabl food",39)
		strp += cpad("\np - shootable potion",39)
		strp += cpad("\nP - indestructabl potion",37)
		strp += cpad("\ni - invisible power",42)
		strp += cpad("\nx - exit",51)
		strp += cpad("\nz - Death",49)
		strp += cpad("\nt - treasure box",44)
		strp += cpad("\nT - teleporter pad",42)
	}
//	strp += "\n * note some address will crash"

	dialog.ShowInformation(dk, strp, w)
}
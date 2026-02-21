package main

import (
	"fmt"
	"math"
	"io/ioutil"
	"time"
	"image"
	"strings"
// /	"image/color"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/data/binding"
)

// fyne menu & window system isolated from keyboard & control now

var w fyne.Window
var a fyne.App
var cwt string		// current window title if detected by mouse move

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
			sav_maz(fil, xbuf, ebuf, eflg, opts.DimX, opts.DimY, 0 - sdb, true)
		}
	}
}

func menu_lodit(y bool) {
	fil := fmt.Sprintf(".ed/g%dmaze%03d.ed",opts.Gtp,opts.mnum+1)
	if y {
		if sdb < 0 {
			Ovwallpat = -1
			cnd := lod_maz(fil, xbuf, ebuf, true, true)
			if cnd >= 0 { fax(&eflg,&tflg,11) }
			remaze(opts.mnum)
		} else {
			sdbit(0)
			ed_maze(true)
		}
	}
}

func menu_rst(y bool) {
	if y {
		sv := opts.edat
		opts.edat = -1	// code to tell maze decompress not to load buffer file
		Ovwallpat = -1
		remaze(opts.mnum)
		opts.edat = sv
		ed_maze(true)
		//ed_sav(opts.mnum+1)	// reset does not overwrite file buffer, still need to save
	}
}

func menu_sav() {
	if opts.edat > 0 {
		dia := fmt.Sprintf("Save buffer for maze %d in .ed/g%dmaze%03d.ed ?",opts.mnum+1,opts.Gtp,opts.mnum+1)
		if sdb >= 0 { dia = fmt.Sprintf("Save buffer sd(%d) to .ed/sd%05d_g%d.ed",sdb,sdb,opts.Gtp) }
		dialog.ShowConfirm("Saving",dia, menu_savit, w)
	} else { dialog.ShowInformation("Save Fail","edit mode is not active!",w) }
}


func menu_lod() {
	if opts.edat > 0 {
		dia := fmt.Sprintf("Load buffer for maze %d from .ed/g%dmaze%03d.ed ?:",opts.mnum+1,opts.Gtp,opts.mnum+1)
		if sdb >= 0 { dia = fmt.Sprintf("Load buffer sd(%d) from .ed/sd%05d_g%d.ed",sdb,sdb,opts.Gtp) }
		dialog.ShowConfirm("Loading",dia, menu_lodit, w)
	} else { dialog.ShowInformation("Load Fail","edit mode is not active!",w) }
}

func menu_res() {
	if opts.edat > 0 {
		dia := fmt.Sprintf("Reset buffer for maze %d from G%d ROM ?\n - reset does not save to file",opts.mnum+1,opts.Gtp)
		dialog.ShowConfirm("Reseting",dia, menu_rst, w)
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
		sav_maz(fil, xbuf, ebuf, eflg, opts.DimX, opts.DimY, mazn, true)

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
		cnd := lod_maz(fil, xbuf, ebuf, true, true)
		sdb = -1
		if cnd >= 0 { fax(&eflg,&tflg,11) }
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
			if tx == 0 { ebuf[xy{tx, ty}] = MAZEOBJ_WALL_REGULAR }
			if ty == 0 { ebuf[xy{tx, ty}] = MAZEOBJ_WALL_REGULAR }
		}
	}}
	pr = false
	opts.dntr = true
	remaze(opts.mnum)
}

func menu_copy() { if opts.edat > 0 { ccp_tog(COPY); if ccp > 0 { smod = "Edit COPY: "}; statlin(cmdhin,sshin) }}
func menu_cut() { if opts.edat > 0 { ccp_tog(CUT); if ccp > 0 { smod = "Edit CUT: "}; statlin(cmdhin,sshin) }}
func menu_paste() { if opts.edat > 0 { ccp_tog(PASTE); if ccp > 0 { smod = "Edit PASTE: "}; statlin(cmdhin,sshin) }}

func menu_option() {

	wc := a.NewWindow("Option controls")
	wc.Resize(fyne.NewSize(400, 500))
	wc.SetContent(optCont(wc))
	wc.Show()
}

/*
// save test code for later
func map_test() {
	for i := 1; i < 6; i++ {
	go func() {
			time.Sleep(time.Duration(i * 240) * time.Millisecond)
   fyne.Do(func() {
		opts.mnum++
		remaze(opts.mnum)
   })
	}()
	}

}
*/

// set menus

func st_menu() {
// sub rune calls
	sr := func(r rune) {}
	if wpalop { sr = palRune }
// default 'quit' menu option does not call needsav !
	menuItemExit := fyne.NewMenuItem("Exit", func() {
		exitsel = true
		needsav()
	})
	menuItemLodf := fyne.NewMenuItem("Load maze from <ctrl-shift>-l",menu_laodf)
	menuItemSava := fyne.NewMenuItem("Save maze as <ctrl-shift>-s",menu_savas)
	menuItemBlan := fyne.NewMenuItem("Blank maze",func() { menu_blank(false) })
	menuItemBlnK := fyne.NewMenuItem("Blank maze, keep decor",func() { menu_blank(true) })
	menuItemRand := fyne.NewMenuItem("Random load",func() { rload(ebuf); ed_maze(true) })
	menuItemFmap := fyne.NewMenuItem("Mapper test",func() { map_fargoal(ebuf); ed_maze(true) })
	menuItemFmapb := fyne.NewMenuItem("Mapper 2 test",func() { map_sword(ebuf); ed_maze(true) })
	menuItemFmapc := fyne.NewMenuItem("Mapper 3 test",func() { map_wide(ebuf); ed_maze(true) })
//	menuItemTmap := fyne.NewMenuItem("Map++ test",func() { map_test() })
	menuItemLin1 := fyne.NewMenuItem("═══════════════",nil)
	menuItemGvs := fyne.NewMenuItem("Gaunlet view sim toggle",func() { gvs = !gvs })
	menuItemNosec := fyne.NewMenuItem("Dont show secret walls toggle",func() { opts.Nosec = !opts.Nosec; remaze(opts.mnum) })
	menuItemWob := fyne.NewMenuItem("Wall border right & bottom",func() { opts.Wob = !opts.Wob; remaze(opts.mnum) })
	menuItemPalf := fyne.NewMenuItem("Palette; map decore toggle",func() { if wpalop {palfol = !palfol}; palete(0) })
	menuItemMute := fyne.NewMenuItem("Mute audio toggle",func() { opts.Mute = !opts.Mute })
	menuFile := fyne.NewMenu("File", menuItemLodf, menuItemSava, menuItemBlan, menuItemBlnK, menuItemFmap, menuItemFmapb, menuItemFmapc, menuItemRand, menuItemLin1, menuItemMute, menuItemExit)

	menuItemSave := fyne.NewMenuItem("Save buffer <ctrl>-s", menu_sav)
	menuItemLoad := fyne.NewMenuItem("Load buffer <ctrl>-l", menu_lod)
	menuItemReset := fyne.NewMenuItem("Reset buffer <ctrl>-r", menu_res)
	menuItemLin2 := fyne.NewMenuItem("═══════════════",nil)
	menuItemCopy := fyne.NewMenuItem("Copy <ctrl>-c", menu_copy)
	menuItemCut := fyne.NewMenuItem("Cut <ctrl>-x", menu_cut)
	menuItemPaste := fyne.NewMenuItem("Paste <ctrl>-p", menu_paste)
	menuItemPb := fyne.NewMenuItem("Paste buffer window", nil)
	menuItemPbshw := fyne.NewMenuItem("Show paste buffer", func() {pbmas_cyc(0)})
	menuItemPbmnx := fyne.NewMenuItem("Next Master pb <ctrl-shft>-O", func() {pbmas_cyc(1)})
	menuItemPbsnx := fyne.NewMenuItem("Next Session pb <ctrl-shft>-P", func() {pbsess_cyc(1)})
	menuItemPbmpr := fyne.NewMenuItem("Prior Master pb", func() {pbmas_cyc(-1)})
	menuItemPbspr := fyne.NewMenuItem("Prior Session pb", func() {pbsess_cyc(-1)})
	menuItemPb.ChildMenu = fyne.NewMenu("",menuItemPbshw,menuItemPbmnx,menuItemPbmpr,menuItemPbsnx,menuItemPbspr)
	menuItemUndo := fyne.NewMenuItem("Undo <ctrl>-z", undo)
	menuItemRedo := fyne.NewMenuItem("Redo <ctrl>-y", redo)
	menuItemUswp := fyne.NewMenuItem("Ult buf <ctrl>-u", uswap)
	menuItemEdKey := fyne.NewMenuItem("Edit Key list", func() { listK = dboxtx("Edit key assignments","",400,800,close_keys,sr); list_keys() })
	menuItemStats := fyne.NewMenuItem("Maze statistics", func() { statsB = dboxtx("Maze stats","",400,700,close_stats,sr); calc_stats() })
	menuItemEdhin := fyne.NewMenuItem("Edit hints", func() {
		strp := ""
		if opts.edat > 0 {
			strp = "Edit mode: "
			if cmdoff { strp += "edit keys" } else { strp += "cmd keys" }
		} else {
			strp = "View mode: cmd keys only"
		}
		dboxtx("Edit hints", strp+"\n══════════════════════════════\nSave - store buffer in file .ed/g{#}maze{###}.ed\n - where g# is 1 or 2 for G¹/G²\n - and ### is the maze number e.g. 003\n"+
			"\nLoad - overwrite current file contents this maze\n\nReset - reload buffer from rom read\n\nedit keys:\nESC: turn editor on, init maze store in .ed/\n"+
			"ESC:	turn editor off, check unsaved buf\n\\	┈ toggle edit keys / command keys\n"+
			"del	┈ set floor\nctrl-del - sticky delete\nC: cycle edit item #++, c: cycle item #-- *\n#c enter number {1-64}c, all set place item *\n"+
			"H: toggle horiz wrap, V: toggle vert wrap\n–—–—–—–—–—–—–—\ntypical key assignment:\n\n"+
			"d, D	┈ horiz, vert door, w, W - walls *\nf, F	┈ foods, k - key, t - treasure *\np, P	┈ potions, T - teleporter *\n"+
			"q	┈ trap wall, r - trap tile *\ni	┈ invisible power *\nx	┈ exit, z - Death *\n"+
			"edit keys lock when pressed, hit 'b' and place doors\nmiddle click - click to reassign current key\n(middle click also activates edit mode,\n and uses default key 'y' if not set)\n"+
			"logo key* + mouse: paint curr key or ctrl-del\n* these edit keys require '\\' mode\n\n\ngved - G¹G² visual editor\ngithub.com/six-of-one/", 400,755,nil,typedRune)
	})
	editMenu := fyne.NewMenu("Edit", menuItemSave, menuItemLoad, menuItemReset, menuItemEdhin, menuItemLin2, menuItemPb, menuItemCopy, menuItemCut, menuItemPaste, menuItemUndo, menuItemRedo, menuItemUswp)

	menuItemKeys := fyne.NewMenuItem("Keys ?", keyhints)
	menuItemOpt := fyne.NewMenuItem("Options", menu_option)
	menuItemOps := fyne.NewMenuItem("Operation", func() {
		data, err := ioutil.ReadFile("ops.txt")
		if err == nil {
			txt := fmt.Sprintf("%s",data)
			dboxtx("Operations", txt, 700, 1000,nil,typedRune)
		}
	})
	menuItemAbout := fyne.NewMenuItem("About", func() {
		dialog.ShowInformation("About G¹G²ved", "Gauntlet / Gauntlet 2 visual editor\nAuthor: Six [a programmer]\n\ngithub.com/six-of-one/", w)
	})
	menuItemLIC := fyne.NewMenuItem("License", func() {
		dialog.ShowInformation("G¹G²ved License", "Gauntlet visual editor - gved\n\n(c) 2025 Six [a programmer]\n\nGPLv3.0\n\nhttps://www.gnu.org/licenses/gpl-3.0.html", w)
	})
	menuHelp := fyne.NewMenu("Help ", menuItemKeys, menuItemEdKey, menuItemOps, menuItemAbout, menuItemLIC)

// list of active main keys for view / edit modes
	hintup = fyne.NewMenu("cmds↓    ?, eE, fFgG, wWqQ, rRt, hm, pPT, sL, S, ilu, A #a",menuItemOpt, menuItemStats, menuItemGvs, menuItemNosec, menuItemWob, menuItemPalf)

// some extra hint info about what the editor is doing
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

	ld_config()			// prog config stuff
	st_menu()			// start the menu
	w.Canvas().SetOnTypedRune(typedRune)	// enable plain key handler for main win
	specialKey(w)		// key handlers for specials
	ed_init()			// initialized the editor package
	get_pbcnt()			// paste buffer cnt (per gauntlet)

// test bit shift op
p := 128 << 2	// 512
q := 128 >> 2	// 32
fmt.Printf("p 128 << 2: %d\nq 128 >> 2: %d\n",p,q)

// setup main tabs
	cmain = container.NewStack()
	splash = container.NewStack()
	spexpl = container.NewStack()
	sprview = container.NewStack()
	maintab = container.NewAppTabs(
		container.NewTabItemWithIcon("Maze view",theme.SearchIcon(),
			cmain,
	),
	container.NewTabItemWithIcon("Game",theme.SettingsIcon(),splash),
	container.NewTabItemWithIcon("Sprites",theme.SearchReplaceIcon(),spexpl),
	)

	w.SetContent(cmain)
	maintab.Refresh()
	maintab.OnSelected = func(t *container.TabItem) {
//fmt.Printf("tab: %s\n",t.Text)
		actab = t.Text
		if actab == "Sprites" { sprite_view() }
	}
	actab = "Maze view"
	gif_lodr("splash/splash1.gif", splash, splim, "")		// pre-load intro gif
	go func() {
		splashrot()
	}()
}

// refresh maze tabs, edit unit
var maintab *container.AppTabs		// tabs unit
var cmain *fyne.Container			// content maze viewer
var spexpl *fyne.Container			// sprite explorer / mem viewer
var pmaz *fyne.Container			// box with image, button & blot
var pimg *canvas.Raster				// current maze image
var actab string					// active tab

func maz_tab(tabcon *fyne.Container, maz *image.NRGBA, mbut *holdableButton, pblot *canvas.Image) {

	tabcon.Remove(pmaz)
	pimg = canvas.NewRasterFromImage(maz)
	pmaz = container.NewStack(mbut, pimg, pblot)
	tabcon.Add(pmaz)
   fyne.Do(func() {
		w.SetContent(maintab)
		pmaz.Refresh()
   })
}
// sub win switch G¹ / G²

func subsw() {

	ccp = NOP
	statlin(cmdhin,"")
	if wpbop { pbmas_cyc(0) }
	if wpalop { palete(0) }
}

// make clickable image wimg in window cw with given size

var rbimg *image.NRGBA			// for the pb paste image dealy
var rbtn *holdableButton
var blotup bool

// blot win short ver for pb image to set blotter

func blotwup(cw fyne.Window, limg *image.NRGBA) {
	t := cw.Title()
	if strings.Contains(t, "G¹G²ved") {
		if blot == ccblot { blot.Resize(fyne.Size{0, 0}) }
		blot = canvas.NewImageFromImage(limg)
	//	box := container.NewStack(rbtn, rbimg, blot)
	//	cw.SetContent(box)
		maz_tab(cmain, rbimg, rbtn, blot)
		blotup = false
fmt.Printf("blotwup\n")
	}
}

// sub win upd

func clikwins(cw fyne.Window, wimg *image.NRGBA, wx int, wy int) {

	bimg := canvas.NewRasterFromImage(wimg)

// turns display into clickable edit area
	btn := newHoldableButton()
	btn.title = cw.Title()
	btn.bw = cw
fmt.Printf("clwin-s tl: %s\n",btn.title)
	if !strings.Contains(btn.title, "G¹G²ved") {
		box := container.NewStack(btn, bimg)
		cw.SetContent(box)
	}
}

// main win updater - sets master rbimg for blot overlay

func clikwinm(cw fyne.Window, wimg *image.NRGBA, wx int, wy int) {

	rbimg = wimg //canvas.NewRasterFromImage(wimg)

// turns display into clickable edit area
	rbtn = newHoldableButton()
	rbtn.title = cw.Title()
	rbtn.bw = cw
fmt.Printf("clwin-m tl: %s\n",rbtn.title)

	cw.Resize(fyne.NewSize(float32(wx), float32(wy)))
//	box := container.NewStack(rbtn, rbimg)		// key to seeing maze & having the click button with full mouse sense
//	cw.SetContent(box)								// and blot coming last is shown on top... huh?
	maz_tab(cmain, rbimg, rbtn, blant)

//fmt.Printf("btn sz %v\n",rbtn.Size())
}

// update contents of main edit window, includes title

func upwin(simg *image.NRGBA, lvp int) {

//												 ┌» un-borded maze is 528 x 528 for a 33 x 33 cell maze
	geow := int(math.Max(560,opts.Geow))	// 560 is min, maze doesnt seem to fit or shrink smaller
	geoh := int(math.Max(586,opts.Geoh))	// 586 min
	dtp := 512.0
	vp := viewp
	rez := false
	rmsg := ""
	if lvp > 0 { vp = lvp; dtp = float64(vp) * 16 }		// a local viewport was passed, likely sb buf view as non edit mazes
	if opts.Wob { dtp = 528.0 }
	if opts.edat > 0 {
//		geow = geow & 0xfe0	+ 13			// lock to multiples of 32
		ngeow := geoh - 71					// square maze area + 26 (tabs go to 71 it seems) for menu bar - window is still 4 wider than maze content
		if ngeow != geow {
			rmsg = "set window ratio to edit"
			rez = true
		}
		geow = ngeow
		dtp = float64(vp) * 16
	}											// having an edit viewport will change 528 - will have to be vport wid (same as high) * 16
	opts.dtec = 16.0 * (float64(geow - 4) / dtp)				// the size of a tile, odd window size may cause issues
fmt.Printf("\nupwin %d x %d dtec: %f (vp: %d dtp %.1f) geom: %d x %d, p:%d\n",opts.DimX,opts.DimY,opts.dtec,vp,dtp,geow,geoh,ccp)
if opts.Verbose { fmt.Printf(" dtec: %f\n",opts.dtec) }			// detected size of a single maze tile in pixels, used for click id of cell x,y
	if lvp < 0 || rez || mbd || ccp == PASTE || gvs { clikwinm(w, simg, geow, geoh) }		// one time init shot

	spx := ""
	if sdb > 0 { spx = fmt.Sprintf("%s sdbuf: %d",rmsg,sdb) }
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

// update trick menu items status line

var sshin string

func statlin(hs string,ss string) {

	hintup.Label = hs
	sshin = ss
	statup.Label = smod + ss
	mainMenu.Refresh()
}

// window resize control

func wizecon() {

	time.Sleep(3 * time.Second)		// some hang time to allow win to display & size, otherwise w x h is 1 x 1
	for {
		bgeow := int(opts.Geow)
		bgeoh := int(opts.Geoh)
// why the +8, +36 needed, will it will ever vary ??
// thus gathered: this is the edges around the main edit and the menu bar and title on top
		width := int(w.Content().Size().Width) + 8
		height := int(w.Content().Size().Height) + 36
//					dw := bgeow - width; dh := bgeoh - height
		if width != bgeow || height != bgeoh {
//					fmt.Printf("Window was resized! st: %d x %d n: %v x %v delta: %d, %d\n",bgeow,bgeoh,w.Content().Size().Width,w.Content().Size().Height,dw,dh)
				// window was resized
// provide live resize so other vis ops dont bounce it back
// for some reason maze updates resize the window down w -= 8 & h -= 36 to minimun
			opts.Geow = float64(width)
			opts.Geoh = float64(height)
			sv_config()			// prog config stuff
		}
		time.Sleep(2 * time.Second)
// run here for the timer
		xbl_typ()
	}
}

// dialog called from kby or menu

func keyhints() {
	strp := ""
	kys := 1
	var lenb float32
	dk := "Command Keys"
	strp = eid+"\n\n"
	if opts.edat > 0 {
		strp += "Edit mode: "
		if cmdoff { strp += "edit keys"; kys = 2; dk = "Editor Keys" } else { strp += "cmd keys" }
	} else {
		strp += "View mode: cmd keys only"
	}
	strp += "\nsingle letter commands\n–—–—–—–—–—–—–—–—–—–—–—"
//		strp += "\n\n?	┈ this list"
	strp += "\nctrl-q	┈ quit program"
	if opts.edat > 0 {
		strp += "\nESC>	┈ exit editor ╗\n\\	┈ toggle cmd keys*"
	} else { strp += "\nESC>	┈ editor mode" }
	if kys == 1 {
		strp += "\nf,F	┈ floor pattern+,-\ng,G	┈ floor color+,-"+
				"\nw,W	┈ wall pattern+,-\ne,E	┈ wall color+,-"+
				"\nr	┈ rotate maze +90°\nR	┈ rotate maze -90°"+
				"\nh	┈ mirror maze horizontal toggle"+
				"\nm	┈ mirror maze vertical toggle"+
				"\np	┈ toggle floor invis\nP	┈ toggle wall invis"+
				"\nT	┈ loop invis things"+
				"\ns	┈ toggle rnd special potion"
	} else {
		strp += "\nH	┈ toggle horiz maze wrap"+
				"\nV	┈ toggle vert maze wrap"+
				"\nC	┈ cycle item++, key c"+
				"\nc	┈ cycle item--, key c"+
				"\n{n}c	┈ set item 1 - 64, key c"+
				"\nL	┈ generator indicate letter"+
				"\nS	┈ cycle sd buffers"+
				"\n{n}S	┈ save curr to buffer #"
	}
	if kys == 1 {
		strp += "\nL	┈ generator indicate letter"+
				"\n{n}S	┈ save curr to buffer #"+
				"\ni	┈ gauntlet mazes r1 - r9"+
				"\nl	┈ use gauntlet rev 14"+
				"\nu	┈ gauntlet 2 mazes"+
//				"\nv	┈ valid address list"+
				"\nv	┈ all maze addr (in termninal)"
	}
	strp += "\nA	┈ toggle a override"+
			"\n{n}a	┈ numeric of valid maze"+
			"\n - load maze 1 - 127 G¹"+
			"\n - load maze 1 - 117 G²"+
			"\n - load address 229376 - 262143 "+
			"\n–—–—–—–—–—–—–—–—–—–—–—"
/*	strb := fmt.Sprintf("\nG%d ",opts.Gtp)
	if G1 {
	if opts.R14 { strb += "(r14)"
		} else { strb += "(r1-9)" }}
	strp += strb */
	if kys == 2 {
		lenb = 142
		strp += "\n\ntypical: key selects item,\n L-click place, M-click assign"+
				"\n–—–—–—–—–—–—–—"+
				"\nDEL>	┈ (hold down) set floor"+
				"\nw	┈ standard walls"+
				"\nW	┈ shootable walls"+
				"\nq	┈ trap wall\nr	┈ trap tile"+
				"\nd	┈ horizontal door"+
				"\nD	┈ vertical door"+
				"\nf	┈ shootable food"+
				"\nF	┈ indestructabl food"+
				"\np	┈ shootable potion"+
				"\nP	┈ indestructabl potion"+
				"\ni	┈ invisible power"+
				"\nx	┈ exit\nz	┈ Death"+
				"\nt	┈ treasure box"+
				"\nT	┈ teleporter pad"
	}
//	strp += "\n * note some address will crash"

//	dialog.ShowInformation(dk, strp, w)
	dboxtx(dk, strp, 400, 700 + lenb,nil,typedRune)
}

// text dialog boxes for all hint sets
// title, content, w, h
// return text box point for updating contents live

var wwlup int

func dboxtx(dt string, dbc string, w float32, h float32, cf func(), sbr func(r rune)) binding.Item[string] {

	ww := a.NewWindow(dt)

	txtB := binding.NewString()
	txtWid := widget.NewEntryWithData(txtB)
	txtWid.MultiLine = true
	txtWid.Disabled()

	// we can disable the Entry field so the user can't modify the text:
	txtB.Set(dbc)
	cn := container.NewBorder(nil, nil, nil, nil, txtWid)

	ww.SetContent(cn)
	ww.Resize(fyne.Size{w, h})
	ww.Show()
	specialKey(ww)
	wwlup++; if wwlup > 8 { wwlup = 1 }
//xb line is boigah knigget
	if dt == "x-line" { txtWid.MultiLine = false
	} else {}
	switch wwlup {

	case 1: wwa = ww; ww.Canvas().SetOnTypedRune(generalRune1); subra = sbr	// this is my mess to allow 'q' 'Q' to close any of these dialogs
	case 2: wwb = ww; ww.Canvas().SetOnTypedRune(generalRune2); subrb = sbr	// there may be some way to add-widget or use a struct to pass ww to generalRune
	case 3: wwc = ww; ww.Canvas().SetOnTypedRune(generalRune3); subrc = sbr	// i have other things to do instead of finding it - so... this
	case 4: wwd = ww; ww.Canvas().SetOnTypedRune(generalRune4); subrd = sbr
	case 5: wwe = ww; ww.Canvas().SetOnTypedRune(generalRune5); subre = sbr
	case 6: wwf = ww; ww.Canvas().SetOnTypedRune(generalRune6); subrf = sbr
	case 7: wwg = ww; ww.Canvas().SetOnTypedRune(generalRune7); subrg = sbr
	case 8: wwh = ww; ww.Canvas().SetOnTypedRune(generalRune8); subrh = sbr
	}

	if cf != nil {
		ww.SetCloseIntercept(func() {			// if cf is passed, assign it to close intercept
			cf()
			ww.Close()
		})
	}
	return txtB
}
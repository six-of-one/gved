package main

import (
	"fmt"
	"image"
	"golang.org/x/image/draw"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/canvas"
	"github.com/fogleman/gg"
// /	"fyne.io/fyne/v2/driver/desktop"
)

var chkg1rom *widget.Check		// these cant even self refer inside the interal inline... huh
var chkg2rom *widget.Check
var filerom *widget.Check
var spsheet *widget.Check

func spchks(c1,c2,c3,c4 bool){

	chkg1rom.Checked = c1; chkg1rom.Refresh()
	chkg2rom.Checked = c2; chkg2rom.Refresh()
	filerom.Checked = c3; filerom.Refresh()
	spsheet.Checked = c4; spsheet.Refresh()
}

func radr_bounds() {						// find a way to get size of rom from file info / loading
// bounds addr
	prcadr = maxint(0,minint(prcadr,65536))	// 0x1000000/slashout - now 64K how large can a rom be? it will prob be read as absolute
	lasadr = fmt.Sprintf("%d",prcadr)
	radr.SetText(lasadr)
	radr.Refresh()
}

func pnum_bounds() {
// bounds pnum sel
		pnumsel = maxint(0,minint(pnumsel,pallim))
		laspnume = fmt.Sprintf("%d",pnumsel)
		pnumen.SetText(laspnume)
		pnumen.Refresh()
}

func xsiz_bounds() {
// bounds sprite x size
		svx = maxint(1,minint(svx,32))	// stamp 32 (8 bit units) takes up 256, seems reasonable, prob have issues if ew proceed past end of rom file
		lasx = fmt.Sprintf("%d",svx)
		xsiz.SetText(lasx)
		xsiz.Refresh()
}

func ysiz_bounds() {
// bounds sprite y size
		svy = maxint(1,minint(svy,32))	// stamp 32 (8 bit units) takes up 256, seems reasonable, prob have issues if ew proceed past end of rom file
		lasy = fmt.Sprintf("%d",svy)
		ysiz.SetText(lasy)
		ysiz.Refresh()
}

var sprview *fyne.Container
var bld_btn *widget.Button
var radr *widget.Entry
var pnumen *widget.Entry
var xsiz *widget.Entry
var ysiz *widget.Entry
var lasadr string = "2048"
var prcadr int = 2048			// process from this addr, 2048 is ghosts
var paltype string = "base"
var pallim int = 0				// each palete list has a # lim, which exceeding causes a crash
var laspnume string = "4"
var pnumsel int = 4				// base pnum 1 - most common items, treasure, foods, potions are in palete 1 of base, 4 = 3rd level ghosts
var svx,svy int = 3,3			// xy size fo stamp
var lasx,lasy string = "3","3"
var pixx = 380					// pixel size to fill - makes square canvas
var lpixx string = "380"
var trnc = 8					// trnech space between sprites
var ltrnc string = "8"

func sprite_view() {

var lim *fyne.Container

	chkg1rom = widget.NewCheck("Gauntlet / G² rom", func(gr bool) {
		fmt.Printf("Gauntlet rom %t\n", gr)
		spchks(gr,false,false,false)
	})
	chkg2rom = widget.NewCheck("Gauntlet II rom", func(gr bool) {
		fmt.Printf("Gauntlet 2 rom %t\n", gr)
		spchks(false,gr,false,false)
	})
	filerom = widget.NewCheck("File rom   ", func(fr bool) {
		fmt.Printf("File rom %t\n", fr)
		spchks(false,false,fr,false)
	})
	spsheet = widget.NewCheck("Sprite sheet   file:", func(ss bool) {
		fmt.Printf("Sprite sheet %t\n", ss)
		spchks(false,false,false,ss)
	})
// g2 mode enable
	g2m := widget.NewCheck("G2 mode", func(g bool) {
		fmt.Printf("Gauntlet 2 gfx mode %t\n", g)
	})
// keep fixed address
	keepr := widget.NewCheck("keep", func(k bool) {
		fmt.Printf("keep addr %t\n", k)
	})
// show address on sheet for each sprite
	showr := widget.NewCheck("show", func(sh bool) {
		fmt.Printf("show addr %t\n", sh)
	})
	showr.Checked = true
// use lvl1 & lvl2 colors for bkg checkerboard
	lvlcol := widget.NewCheck("lvl color", func(k bool) {
		fmt.Printf("use level color custom %t\n", k)
	})
// gauntlet palete type - vars lists in palettes.go
	selptype := widget.NewSelect([]string{"teleff","floor","gfloor","wall","gwall","base","gbase","warrior","valkyrie","wizard","elf","trap","stun","secret","shrub","forcefield"}, func(str string) {
		fmt.Printf("Select ptype: %s\n", str)
		paltype = str
		pallim = ptyp_lim[str]
	})
	selptype.SetSelected("base")
	ptyp_label := widget.NewLabelWithStyle("Pal type:", fyne.TextAlignLeading, fyne.TextStyle{Monospace: false})

// select palete num, limited for each palete type
	pnumen = widget.NewEntry()
	pnumen.Resize(fyne.Size{60, optht})
	pnumen.SetText(laspnume)
	pnumen.OnChanged = func(s string) {

		fmt.Sscanf(s,"%d",&pnumsel)

	}
	pnum_label := widget.NewLabelWithStyle("pnum:", fyne.TextAlignLeading, fyne.TextStyle{Monospace: false})
// get a "stamp" size too, controls how rom is read into sprites
	xsiz = widget.NewEntry()
	xsiz.OnChanged = func(s string) {

		fmt.Sscanf(s,"%d",&svx)
	}
	xsiz.SetText(lasx)
	ysiz = widget.NewEntry()
	ysiz.OnChanged = func(s string) {

		fmt.Sscanf(s,"%d",&svy)
	}
	ysiz.SetText(lasy)
// size of stamp, x by y
	ssiz_label := widget.NewLabelWithStyle("size:", fyne.TextAlignLeading, fyne.TextStyle{Monospace: false})
	x_label := widget.NewLabelWithStyle(" x ", fyne.TextAlignLeading, fyne.TextStyle{Monospace: false})
// pxiel size of image to build
	xpxz := widget.NewEntry()
	xpxz.OnChanged = func(s string) {

		fmt.Sscanf(s,"%d",&pixx)
	}
	xpxz.SetText(lpixx)
	pixs_label := widget.NewLabelWithStyle("pixel sz:", fyne.TextAlignLeading, fyne.TextStyle{Monospace: false})
// space between sprites
	trench := widget.NewEntry()
	trench.OnChanged = func(s string) {

		fmt.Sscanf(s,"%d",&trnc)
	}
	trench.SetText(ltrnc)
	trench_label := widget.NewLabelWithStyle("trench:", fyne.TextAlignLeading, fyne.TextStyle{Monospace: false})
// address to start rom read
	if lasadr == "" { lasadr = "0" }
	adr_label := widget.NewLabelWithStyle("Address: ", fyne.TextAlignLeading, fyne.TextStyle{Monospace: false})
	adr_spc := widget.NewLabelWithStyle("            ", fyne.TextAlignLeading, fyne.TextStyle{Monospace: false})
	radr = widget.NewEntry()
	radr.OnChanged = func(s string) {

		fmt.Sscanf(s,"%d",&prcadr)
	}
	radr.SetText(lasadr)
	radr.Resize(fyne.Size{100, optht})
// build button
// need - g1/g2 flag check, tranpar flag
// adjust so it fills test area w/ gx,gy
	bld_btn = widget.NewButton("BUILD", func() {
		var bstamp Stamp
		ova,ovb = HRGB{0xff1f1f1f},HRGB{0xff2f2f2f}
		if lvlcol.Checked { ova,ovb = lvl1col,lvl2col }
// change ops so bad inputs default here

// bounds pnum sel
		pnum_bounds()
// bounds sprite x size
		xsiz_bounds()
// bounds sprite y size
		ysiz_bounds()
// bounds pixel size
		pixx = maxint(128,minint(pixx,1200))	// stamp 32 (8 bit units) takes up 256, seems reasonable, prob have issues if ew proceed past end of rom file
		lpixx = fmt.Sprintf("%d",pixx)
		xpxz.SetText(lpixx)
		xpxz.Refresh()
// bounds addr
		radr_bounds()

		bas := loadfail(pixx,pixx)
		gtop := gg.NewContext(32, 12)
// gtop font
		if err := gtop.LoadFontFace(".font/VrBd.ttf", 7); err != nil {
			panic(err)
			}
		if !chkg1rom.Checked && !chkg2rom.Checked { spchks(true,false,false,false) }
		gsv := G1
		if g2m.Checked { G1 = false }
		bstamp = Stamp{} //itemGetStamp("key")
		gx,gy := svx*8+trnc, svy*8+trnc
		suby := 65 / gy

		fx,fy := pixx / gx, (pixx / gy) - suby
		fmt.Sscanf(lasadr,"%d",&prcadr)
		for y := 0; y <= fy; y++ {
		for x := 0; x <= fx; x++ {
			bstamp.numbers = tilerange(prcadr, svx * svy)
			st := fmt.Sprintf("%d",prcadr)
			prcadr += svx * svy
			bstamp.width = svx
			bstamp.trans0 = false
			bstamp.pnum = pnumsel
			bstamp.ptype = paltype
//fmt.Printf("Write sprite : %s: %d, %d x %d adr: %X - @%d, %d\n",paltype,pnumsel,fx,fy,prcadr,x*gx, y*gy)
			fillstamp(&bstamp)
			writestamptoimage(G1,bas, &bstamp, x*gx, y*gy)

			if showr.Checked {
				gtop.Clear()
				gtop.SetRGB(0.5, 0.5, 0.5)
				gtop.DrawStringAnchored(st, 0, 6, 0, 0.5)
				gtop.SetRGB(0.12, 0.12, 0.12)
				gtopim := gtop.Image()
				offset := image.Pt(x*gx, y*gy+(svy*8)-2)
				draw.Draw(bas, gtopim.Bounds().Add(offset), gtopim, image.ZP, draw.Over)
			}
		}}
		if keepr.Checked { fmt.Sscanf(lasadr,"%d",&prcadr) }
//fmt.Printf("dis sprite gxy: %d x %d fxy %d, %d svxy %d - %d, suby %d\n",gx,gy,fx,fy,svx,svy,suby)
		bld := canvas.NewRasterFromImage(bas)
		gif_blnk(lim)
		savetopng("tst.png", bas)
		sprview.Remove(lim)
		lim = container.NewWithoutLayout(bld)
		sprview.Add(lim)
		bld.Resize(fyne.Size{800, 800})
		G1 = gsv
//		func(b *Button) TypedRune(rune)
	})
//	bld_btn.SetOnTypedRune(typedRune)

	fnent := widget.NewEntry()
	fnent.Resize(fyne.Size{370, optht})
	ld := container.New(
		layout.NewVBoxLayout(),
		container.New(layout.NewHBoxLayout(),
			chkg1rom, ptyp_label, selptype, pnum_label,pnumen,g2m,trench_label,trench,
		),
		container.New(layout.NewHBoxLayout(),
			filerom, spsheet, container.NewWithoutLayout(fnent),
		),
		container.New(layout.NewHBoxLayout(),
			bld_btn,keepr, pixs_label, xpxz, ssiz_label, xsiz, x_label, ysiz, adr_label, container.NewWithoutLayout(radr), adr_spc,showr,lvlcol,
		),
		sprview,
	)
	spexpl.Remove(lim)
	lim = container.NewStack(ld)
	spexpl.Add(lim)
	fyne.Do(func() {
		lim.Refresh()
	})
// blank view on launch
	ova,ovb = HRGB{0xff1f1f1f},HRGB{0xff2f2f2f}
	bas := loadfail(400, 400)
	bld := canvas.NewRasterFromImage(bas)
	savetopng("tst.png", bas)
	sprview.Remove(lim)
	lim = container.NewWithoutLayout(bld)
	sprview.Add(lim)
	bld.Resize(fyne.Size{1000, 1000})
}
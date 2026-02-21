package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/canvas"
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

var sprview *fyne.Container
var laspnume string
var lasadr string = "0"
var paltype string = "base"
var pallim int = 0				// each palete list has a # lim, which exceeding causes a crash
var pnumsel int = 1				// base pnum 1 - most common items, treasure, foods, potions are in palete 1 of base
var sx,sy int = 2,2				// xy size fo stamp
var lasx,lasy string = "2","2"

func sprite_view() {

	chkg1rom = widget.NewCheck("Gauntlet ", func(gr bool) {
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
// gauntlet palete type - vars lists in palettes.go
	selptype := widget.NewSelect([]string{"teleff","floor","gfloor","wall","gwall","base","gbase","warrior","valkyrie","wizard","elf","trap","stun","secret","shrub","forcefield"}, func(str string) {
		fmt.Printf("Select ptype: %s\n", str)
		paltype = str
		pallim = ptyp_lim[str]
	})
	selptype.SetSelected("base")
	ptyp_label := widget.NewLabelWithStyle("Pal type:", fyne.TextAlignLeading, fyne.TextStyle{Monospace: false})

// select palete num, limited for each palete type
	pnumen := widget.NewEntry()
	pnumen.Resize(fyne.Size{60, optht})
	if laspnume == "" { laspnume = "0" }
	pnumen.SetText(laspnume)
	pnumen.OnChanged = func(s string) {

		fmt.Sscanf(s,"%d",&pnumsel)
		pnumsel = maxint(0,minint(pnumsel,pallim))
		laspnume = fmt.Sprintf("%d",pnumsel)
		pnumen.SetText(laspnume)
		pnumen.Refresh()
	}
	pnum_label := widget.NewLabelWithStyle("pnum:", fyne.TextAlignLeading, fyne.TextStyle{Monospace: false})
// get a "stamp" size too, controls how rom is read into sprites
	xsiz := widget.NewEntry()
	xsiz.OnChanged = func(s string) {

		fmt.Sscanf(s,"%d",&sx)
		sx = maxint(1,minint(sx,32))	// stamp 32 (8 bit units) takes up 256, seems reasonable, prob have issues if ew proceed past end of rom file
		lasx = fmt.Sprintf("%d",sx)
		xsiz.SetText(lasx)
		xsiz.Refresh()
	}
	xsiz.SetText(lasx)
	ysiz := widget.NewEntry()
	ysiz.SetText(lasy)
// size of stamp, x by y
	ssiz_label := widget.NewLabelWithStyle("size:", fyne.TextAlignLeading, fyne.TextStyle{Monospace: false})
	x_label := widget.NewLabelWithStyle(" x ", fyne.TextAlignLeading, fyne.TextStyle{Monospace: false})
// address to start rom read
	if lasadr == "" { lasadr = "0" }
	adr_label := widget.NewLabelWithStyle("Address: ", fyne.TextAlignLeading, fyne.TextStyle{Monospace: false})
	radr := widget.NewEntry()
	radr.SetText(lasadr)
	radr.Resize(fyne.Size{100, optht})
// build button
	bld_btn := widget.NewButton("BUILD", func() {

	})

	fnent := widget.NewEntry()
	fnent.Resize(fyne.Size{370, optht})
	ld := container.New(
		layout.NewVBoxLayout(),
		container.New(layout.NewHBoxLayout(),
			chkg1rom, chkg2rom, ptyp_label, selptype, pnum_label,pnumen,
		),
		container.New(layout.NewHBoxLayout(),
			filerom, spsheet, container.NewWithoutLayout(fnent),
		),
		container.New(layout.NewHBoxLayout(),
			bld_btn, ssiz_label, xsiz, x_label, ysiz, adr_label, container.NewWithoutLayout(radr),
		),
		sprview,
	)
var lim *fyne.Container
	spexpl.Remove(lim)
	lim = container.NewStack(ld)
	spexpl.Add(lim)
	fyne.Do(func() {
		lim.Refresh()
//fmt.Printf("Splash load: %s\n",fn)
	})
	bas := loadfail(400, 400)
	bld := canvas.NewRasterFromImage(bas)
	savetopng("tst.png", bas)
	sprview.Remove(lim)
	lim = container.NewWithoutLayout(bld)
	sprview.Add(lim)
	bld.Resize(fyne.Size{800, 800})
}
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
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/dialog"
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

func xysiz_bounds() {
// bounds sprite x size
		svx = maxint(1,minint(svx,32))	// stamp 32 (8 bit units) takes up 256, seems reasonable, prob have issues if ew proceed past end of rom file
		lasx = fmt.Sprintf("%d",svx)
		xsiz.SetText(lasx)
		xsiz.Refresh()
// bounds sprite y size
		svy = maxint(1,minint(svy,32))	// stamp 32 (8 bit units) takes up 256, seems reasonable, prob have issues if ew proceed past end of rom file
		lasy = fmt.Sprintf("%d",svy)
		ysiz.SetText(lasy)
		ysiz.Refresh()
}

func shsiz_bounds() {
		shx = maxint(1,minint(shx,1024))	// sheet read x & y, would a sprite sheet ever have > 1024 x 1024
		shlasx = fmt.Sprintf("%d",shx)
		shxsiz.SetText(shlasx)
		shxsiz.Refresh()
// bounds sprite y size
		shy = maxint(1,minint(shy,1024))
		shlasy = fmt.Sprintf("%d",shy)
		shysiz.SetText(shlasy)
		shysiz.Refresh()
}

func rowcol_bounds() {
		shc = maxint(0,minint(shc,1024))	// sheet read c & r, would a sprite sheet ever have > 1024 x 1024
		lasc = fmt.Sprintf("%d",shc)
		redc.SetText(lasc)
		redc.Refresh()
// bounds sprite y size
		shr = maxint(0,minint(shr,1024))
		lasr = fmt.Sprintf("%d",shr)
		redr.SetText(lasr)
		redr.Refresh()
}

func adr_mode() {
	shxsiz.Hide(); shysiz.Hide()
	redc.Hide(); redr.Hide()
	xsiz.Show(); ysiz.Show()
	adr_label.SetText("Address: ")
	radr.Show()
	adr_spc.Show()
}

var sprview *fyne.Container
var bld_btn *widget.Button
var radr *widget.Entry
var pnumen *widget.Entry
var xsiz *widget.Entry
var ysiz *widget.Entry
var shxsiz *widget.Entry
var shysiz *widget.Entry
var redc *widget.Entry
var redr *widget.Entry
var adr_label *widget.Label
var adr_spc *widget.Label
var lasadr string = "2048"
var prcadr int = 2048			// process from this addr, 2048 is ghosts
var paltype string = "base"
var pallim int = 0				// each palete list has a # lim, which exceeding causes a crash
var laspnume string = "4"
var pnumsel int = 4				// base pnum 1 - most common items, treasure, foods, potions are in palete 1 of base, 4 = 3rd level ghosts
var svx,svy int = 3,3			// xy size fo stamp
var shx,shy int = 16,16			// xy size for sheet read
var lasx,lasy string = "3","3"
var shlasx,shlasy string = "16","16"
var shr,shc int = 0,0			// row col coord for sheet read
var lasr,lasc string = "0","0"
var pixx = 380					// pixel size to fill - makes square canvas
var lpixx string = "380"
var trnc = 8					// trnech space between sprites
var ltrnc string = "8"
var sheet_read bool

func sprite_view() {

var lim *fyne.Container

	chkg1rom = widget.NewCheck("Gauntlet / G² rom", func(gr bool) {
		fmt.Printf("Gauntlet rom %t\n", gr)
		spchks(gr,false,false,false)
		adr_mode()
	})
	chkg2rom = widget.NewCheck("Gauntlet II rom", func(gr bool) {
		fmt.Printf("Gauntlet 2 rom %t\n", gr)
		spchks(false,gr,false,false)
	})
	filerom = widget.NewCheck("File rom   ", func(fr bool) {
		fmt.Printf("File rom %t\n", fr)
		spchks(false,false,fr,false)
		adr_mode()
	})
	spsheet = widget.NewCheck("Sprite sheet   file:", func(ss bool) {
		fmt.Printf("Sprite sheet %t\n", ss)
		spchks(false,false,false,ss)
		shxsiz.Show(); shysiz.Show()
		redc.Show(); redr.Show()
		xsiz.Hide(); ysiz.Hide()
		radr.Hide()
		adr_spc.Hide()
		adr_label.SetText("Read c/r:")
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
// sprite sheet sizes
	shxsiz = widget.NewEntry()
	shxsiz.OnChanged = func(s string) {

		fmt.Sscanf(s,"%d",&shx)
	}
	shxsiz.SetText(shlasx)
	shxsiz.Hide()
	shysiz = widget.NewEntry()
	shysiz.OnChanged = func(s string) {

		fmt.Sscanf(s,"%d",&shy)
	}
	shysiz.SetText(shlasy)
	shysiz.Hide()
//  sheet read r/c
	redc = widget.NewEntry()
	redc.OnChanged = func(s string) {

		fmt.Sscanf(s,"%d",&shc)
	}
	redc.SetText(lasc)
	redc.Hide()
	redr = widget.NewEntry()
	redr.OnChanged = func(s string) {

		fmt.Sscanf(s,"%d",&shr)
	}
	redr.SetText(lasr)
	redr.Hide()

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
// transparent sprites
	xpar := widget.NewCheck("no bkg", func(k bool) {
		fmt.Printf("show transparent sprites %t\n", k)
	})
// address to start rom read
	if lasadr == "" { lasadr = "0" }
	adr_label = widget.NewLabelWithStyle("Address: ", fyne.TextAlignLeading, fyne.TextStyle{Monospace: false})
	adr_spc = widget.NewLabelWithStyle("            ", fyne.TextAlignLeading, fyne.TextStyle{Monospace: false})
	radr = widget.NewEntry()
	radr.OnChanged = func(s string) {

		fmt.Sscanf(s,"%d",&prcadr)
	}
	radr.SetText(lasadr)
	radr.Resize(fyne.Size{100, optht})
// file name for rom/sheet
	fnent := widget.NewEntry()
	fnent.Resize(fyne.Size{370, optht})
	fnload := widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), func() {

		fileDiag := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				fmt.Println("Save as Error:", err)
				return
			}
			if reader == nil {
				fmt.Println("No file selected")
				return
			}

			fmt.Println("Selected:", reader.URI().Path())
			fnent.SetText(reader.URI().Path())

		}, w)
		fileDiag.Show()
		fileDiag.Resize(fyne.NewSize(float32(opts.Geow - 10), float32(opts.Geoh - 30)))
	})

// build button
// need - g1/g2 flag check, tranpar flag
// adjust so it fills test area w/ gx,gy
	bld_btn = widget.NewButton("BUILD", func() {
		var bstamp Stamp
		parimg = nil
		ova,ovb = HRGB{0xff1f1f1f},HRGB{0xff2f2f2f}
		if lvlcol.Checked { ova,ovb = lvl1col,lvl2col }
// change ops so bad inputs default here

// bounds pnum sel
		pnum_bounds()
// bounds sprite x size
		xysiz_bounds()
// bounds sprite y size
		shsiz_bounds()
// bounds pixel size
		pixx = maxint(70,minint(pixx,1200))	// stamp 32 (8 bit units) takes up 256, seems reasonable, prob have issues if ew proceed past end of rom file
		lpixx = fmt.Sprintf("%d",pixx)
		xpxz.SetText(lpixx)
		xpxz.Refresh()
// bounds trench size
		trnc = maxint(0,minint(trnc,32))
		ltrnc = fmt.Sprintf("%d",trnc)
		trench.SetText(ltrnc)
		trench.Refresh()
// bounds addr
		radr_bounds()

		bas := blankimage(pixx,pixx)
		if ova != (HRGB{0}) && ovb != (HRGB{0}) { bas = loadfail(pixx,pixx) }

		gtop := gg.NewContext(32, 12)
// gtop font
		if err := gtop.LoadFontFace(".font/VrBd.ttf", 7); err != nil {
			panic(err)
			}
		if !chkg1rom.Checked && !spsheet.Checked  { spchks(true,false,false,false) }
		sheet_read = spsheet.Checked
		uroms := !sheet_read
		gsv := G1
		if g2m.Checked { G1 = false }
		fx,fy,gx,gy := 0,0,0,0
		subf := int((float64(pixx) / (opts.Geoh-190))* 116)
		if uroms {
	// calc how many rows & cols of sprites will fit in pixel area
			gx,gy = svx*8+trnc, svy*8+trnc
//fmt.Printf("subf: %d, %f, %f\n",subf, float64(pixx) / (opts.Geoh-190),(float64(pixx) / (opts.Geoh-190)) * 118)
			bstamp = Stamp{} //itemGetStamp("key")
			fmt.Sscanf(lasadr,"%d",&prcadr)
		} else {
			gx,gy = shx+trnc,shy+trnc
			_,_,parimg = itemGetPNG(fnent.Text)
		}
		for x := 1; x <= 64; x++ { if x * gx < (pixx+7) { fx = x-1 } }
		for y := 1; y <= 64; y++ { if y * gy < (pixx - subf) { fy = y } }

		st := ""
		usvy := svy*8
fmt.Printf("dis sprite gxy: %d x %d fxy %d, %d svxy %d - %d\n",gx,gy,fx,fy,svx,svy)
		for y := 0; y <= fy; y++ {
		for x := 0; x <= fx; x++ {
		  if uroms {
			bstamp.numbers = tilerange(prcadr, svx * svy)
			st = fmt.Sprintf("%d",prcadr)
			prcadr += svx * svy
			bstamp.width = svx
			bstamp.trans0 = xpar.Checked
			bstamp.pnum = pnumsel
			bstamp.ptype = paltype
//fmt.Printf("Write sprite : %s: %d, %d x %d adr: %X - @%d, %d\n",paltype,pnumsel,fx,fy,prcadr,x*gx, y*gy)
			fillstamp(&bstamp)
			writestamptoimage(G1,bas, &bstamp, x*gx, y*gy)
		  } else {
			if parimg != nil {
				 writepngtoimage(bas, shx,shy,0,0,shc+x,shr+y,x*gx, y*gy,0)
				 st = fmt.Sprintf("%d,%d",shc+x,shr+y)
			}
			usvy = shy
		  }

			if showr.Checked {
				gtop.Clear()
				gtop.SetRGB(0.5, 0.5, 0.5)
				gtop.DrawStringAnchored(st, 0, 6, 0, 0.5)
				gtop.SetRGB(0.12, 0.12, 0.12)
				gtopim := gtop.Image()
				offset := image.Pt(x*gx, y*gy+(usvy)-2)
				draw.Draw(bas, gtopim.Bounds().Add(offset), gtopim, image.ZP, draw.Over)
			}
		}}
		if keepr.Checked { if uroms { fmt.Sscanf(lasadr,"%d",&prcadr) }}
		bld := canvas.NewRasterFromImage(bas)
		gif_blnk(lim)
		savetopng("sheet.png", bas)
		sprview.Remove(lim)
		lim = container.NewWithoutLayout(bld)
		sprview.Add(lim)
		bld.Resize(fyne.Size{800, 800})
		G1 = gsv
//		func(b *Button) TypedRune(rune)
	})
//	bld_btn.SetOnTypedRune(typedRune)

	ld := container.New(
		layout.NewVBoxLayout(),
		container.New(layout.NewHBoxLayout(),
			chkg1rom, ptyp_label, selptype, pnum_label,pnumen,g2m,trench_label,trench,xpar,
		),
		container.New(layout.NewHBoxLayout(),
			filerom, spsheet, fnload, container.NewWithoutLayout(fnent),
		),
		container.New(layout.NewHBoxLayout(),
			bld_btn,keepr, pixs_label, xpxz, ssiz_label, xsiz,shxsiz, x_label, ysiz,shysiz, adr_label, container.NewWithoutLayout(radr),redc,redr,adr_spc,showr,lvlcol,
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

func sprites_keys() {
//	strp := "t,r		- palette type +,-\n"+
	strp := "p,o		- pnum +,-\n"+
			"x,z		- sprite size x (col) +,-\n"+
			"y,u		- sprite size y (row) +,-\n"+
			"<SH>X,Z	- sheet size col +,- 8		<LOGO> for +,- 1\n"+
			"<SH>Y,U	- sheet size row +,- 8		<LOGO> for +,- 1\n"+
			"b,v		- address +,-{sprite size = x*y}\n"+
			"←→		- address -,+              shift modify -,+ 4\n"+
			"↑↓		- address -10,+10     shift modify -,+ 40\n"+
			"pgup		- address -100           shift modify -,+ 1000\n"+
			"pgdn		- address +100\nq		- close key hints\n"+
			"\nsprite sheet mode:\n"+
			"r,e		- read row +,-\n"+
			"c,x		- read col +,-\n"+
			"\n build writes output 'sheet.png'"
fmt.Print(strp)
	dboxtx("Sprite viewer", strp, 480, 370,nil,typedRune)
}
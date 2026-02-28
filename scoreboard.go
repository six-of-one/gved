package main

import (
	"fmt"
	"image"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/canvas"
//	"fyne.io/fyne/v2/layout"
	"github.com/fogleman/gg"
	"golang.org/x/image/draw"
)

// score board (seen in splash screen rotate) and gameplay info

// display scores while in game, and some status

var tsb *fyne.Container				// title small brand
var sbtl *fyne.Container			// title upd for score board
var scorec *fyne.Container			// scoreboard contains
var scors *fyne.Container			// image for scores

func dlg_scboard(stsb string) {

	wc := a.NewWindow("High Score!")
	wc.Resize(fyne.NewSize(270, 600))
	tsb = container.NewStack()
	sbtl = container.NewStack()

	ova,ovb = HRGB{0xff010101},HRGB{0xff010101}
	bas := loadfail(270, 600)
	bld := canvas.NewRasterFromImage(bas)
	scorec = container.NewWithoutLayout(bld)

	osb := container.NewWithoutLayout(
		scorec,
		tsb,
	)
	gif_lodr(stsb, tsb, sbtl, "")
	wc.SetContent(osb)
	wc.Show()
	bld.Resize(fyne.Size{270, 600})
	tsb.Resize(fyne.NewSize(270, 120))
	tsb.Refresh()
}

// 1 cycle to update scoreboard
// ---- think about only updating if data changed

func scor_post(ntsb string) {

// REPLACE: samples only, need real vars for player data
wxtr := []int{0,0,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,}
wpotsmp := 10
wkeysmp := 5
vmode := 0		// 0 for  G¹ / G², 1 for se
usbk := false

	if ntsb != "" {
		gif_lodr(ntsb, tsb, sbtl, "")
//		tsb.Resize(fyne.NewSize(270, 120))
		tsb.Refresh()
	}
	ova,ovb = HRGB{0xff000001},HRGB{0xff000001}
	img := loadfail(270, 600)
//	p,q,r := 0.0,0.0,0.0
	lfont := ".font/VPPxl.ttf"
	sfont := 8.0
	mlen := 42
	c := ""
sb_loop := func(iv int, sbv []dysb) {
		x := 26
		y := iv * 18 + 112
	//	c = fmt.Sprintf("%02d GAUNTLET, 7653428901: WIZARD Level 7",font_tst)
	if sbv[iv].fnr > 0 {
		c = sbv[iv].msb
		mlen = len(c) * 14
		lfont = fmt.Sprintf(".font/%s",ld_font[sbv[iv].fnr])
		sfont = sbv[iv].sz
//		p,q,r = sbv[i].br,sbv[i].bg,sbv[i].bb
		x2,y2 := sbv[iv].xov,sbv[iv].yov
//fmt.Printf("#: %d font: %s, x,y: %d,%d, l:%d, bcol: %0X ox,oy %d %d, msg: %s\n",iv,lfont,x,y,mlen,sbv[iv].bkg,x2,y2, c)
		if x2 > 0 { x = x2 }
		if y2 > 0 { y = y2 }

	gtop := gg.NewContext(mlen, 16)
	if err := gtop.LoadFontFace(lfont, sfont); err == nil {
		gtop.Clear()
		gtop.SetRGB(sbv[iv].r/255.0, sbv[iv].g/255.0, sbv[iv].b/255.0)
		cpos := 0.0
		gtop.DrawStringAnchored(c, 6, 6, cpos, 0.5)
		bc := HRGB{sbv[iv].bkg}
		if bc != (HRGB{0}) && usbk {
	//		cc := HRGB{sbv[iv].bkg}
			ova,ovb = bc,bc
			bimg := loadfail(270, 18)
			offset := image.Pt(x-14, y-2)
			draw.Draw(img, bimg.Bounds().Add(offset), bimg, image.ZP, draw.Over)
		}
//		if p+q+r > 0 { gtop.SetRGB(p/255.0, q/255.0, r/255.0);fmt.Printf("bkg col\n")}
		gtopim := gtop.Image()
		offset := image.Pt(x+sbv[iv].adj, y)
		draw.Draw(img, gtopim.Bounds().Add(offset), gtopim, image.ZP, draw.Over)
	}}
	if sbv[iv].fnr < 0 {
fmt.Printf("v: %d ox,oy %d %d,== ? %t\n",sbv[iv].fnr,sbv[iv].xov,sbv[iv].yov, (sbv[iv].fnr == -4))
		a := sbv[iv].adj
		if sbv[iv].fnr == -3 {		// keys
			err,_,wp := itemGetPNG(xpwr_gly[0][vmode])
			if err == nil {
			x,y := sbv[iv].xov,sbv[iv].yov
			for k := 0; k < wkeysmp; k++ {
				offset := image.Pt(x+k*a, y)
//fmt.Printf("#: %d font: %s, x,y: %d,%d, l:%d, bcol: %d ox,oy %d %d, msg: %s\n",iv,lfont,x,y,mlen,k*a	,x,y, c)
				draw.Draw(img, wp.Bounds().Add(offset), wp, image.ZP, draw.Over)
			}}else { fmt.Printf("lod issue: gfx/sb/key.png\n"); fmt.Println(err)}
		}
		if sbv[iv].fnr == -4 {		// keys
			err,_,wp := itemGetPNG(xpwr_gly[1][vmode])
			if err == nil {
			x,y := sbv[iv].xov,sbv[iv].yov
			for k := 0; k < wpotsmp; k++ {
				offset := image.Pt(x+k*a, y)
				draw.Draw(img, wp.Bounds().Add(offset), wp, image.ZP, draw.Over)
			}}else { fmt.Printf("lod issue: gfx/sb/potion.png\n"); fmt.Println(err)}
		}
		if sbv[iv].fnr == -6 {		// xpwr
		for p := 2; p <= max_xpwr; p++ {
		  if wxtr[p] > 0 {									// later this will handle levels of pwr in expand
			err,_,wp := itemGetPNG(xpwr_gly[p][vmode])			// make all these little gets an array of some sorts
			if err == nil {
			x,y := xpwr_pos[p][0],sbv[iv].yov+xpwr_pos[p][1]
			offset := image.Pt(x, y)
			draw.Draw(img, wp.Bounds().Add(offset), wp, image.ZP, draw.Over)
			} else { fmt.Printf("lod issue: %s\n",xpwr_gly[p][0]); fmt.Println(err)}
		}}}
	}
}

	for i := 1; i <= max_sb; i++ {
		sb_loop(i,sb)
	}
	for i := 1; i <= max_sb2; i++ {
		sb_loop(i,sb2)
	}

	scorec.Remove(scors)
	bld := canvas.NewRasterFromImage(img)
savetopng("tst.png", img)
	scors = container.NewStack(bld)
	scorec.Add(scors)
	scors.Resize(fyne.NewSize(270, 600))
	scors.Refresh()

}
// to change tsb
/*
	gif_lodr("splash/sanc_tsb.gif", tsb, sbtl, "")
	tsb.Resize(fyne.NewSize(270, 120))
	tsb.Refresh()
*/

// to update scores

/*	ova,ovb = HRGB{0xff000001},HRGB{0xff000001}
	bas := loadfail(270, 600)
	bld := canvas.NewRasterFromImage(bas)

	scorec.Remove(scors)
	scors = container.NewStack(bld)
	scorec.Add(scors)
	scors.Resize(fyne.NewSize(270, 600))
	scors.Refresh()
*/

// upd from test file

/*
	err,tst,_ := itemGetPNG("splash/tmp.jpg")
if err != nil { fmt.Printf("jpg err:\n"); fmt.Println(err)}
	scorec.Remove(scors)
	scors = container.NewStack(tst)
	scorec.Add(scors)
	scors.Resize(fyne.NewSize(270, 480))
	scors.Refresh()
	*/
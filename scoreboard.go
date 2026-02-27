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

func scor_post() {

	ova,ovb = HRGB{0xff000001},HRGB{0xff000001}
	img := loadfail(270, 600)
	p,q,r := 0.0,0.0,0.0
	lfont := ".font/VPPxl.ttf"
	sfont := 8.0
	x := 26
	mlen := 42
	c := ""
	for i := 1; i <= max_font; i++ {
		y := i * 18 + 112
	//	c = fmt.Sprintf("%02d GAUNTLET, 7653428901: WIZARD Level 7",font_tst)
	if sb[i].fnr > 0 {
		c = sb[i].msb
		mlen = len(c) * 14
		lfont = fmt.Sprintf(".font/%s",ld_font[sb[i].fnr])
		sfont = sb[i].sz
		p,q,r = sb[i].br,sb[i].bg,sb[i].bb
fmt.Printf("#: %d font: %s, x,y: %d,%d, l:%d, bcol: %0.1f %0.1f %0.1f, msg: %s\n",i,lfont,x,y,mlen,p,q,r, c)

	gtop := gg.NewContext(mlen, 14)
	if err := gtop.LoadFontFace(lfont, sfont); err == nil {
		gtop.Clear()
		gtop.SetRGB(sb[i].r/255.0, sb[i].g/255.0, sb[i].b/255.0)
		cpos := 0.0
		gtop.DrawStringAnchored(c, 6, 6, cpos, 0.5)
		if p+q+r > 0 { gtop.SetRGB(p/255.0, q/255.0, r/255.0);fmt.Printf("bkg col\n")}
		gtopim := gtop.Image()
		offset := image.Pt(x+sb[i].adj, y)
		draw.Draw(img, gtopim.Bounds().Add(offset), gtopim, image.ZP, draw.Over)
	}}}
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
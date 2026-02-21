package main

import (
	"fmt"
	"math/rand"
//	"strconv"
//	"strings"
	"image"
	"image/draw"
	"time"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/canvas"
)

// splash screen rotator
// no video for the time being... he doesnt like them

var splRot = 6000
var splCyc = 0
var splLoop = "0123456789ABCDEFKc"

type Video struct {
	Height     int
	Width      int
	Src        string
	Muted      bool
	Visibility string
	// Simulate play, pause, load methods
}

func (v *Video) Play()   {}
func (v *Video) Pause()  {}
func (v *Video) Load()   {}
func (v *Video) SetVisibility(vis string) {
	v.Visibility = vis
}

type Image struct {
	Src string
}

type Document struct {
	IntroVid  *Video
	Splash    *Image
	ScorDiv   bool // true = visible, false = hidden
	Splashrot *Image
}

var document = &Document{
	IntroVid:  &Video{},
	Splash:    &Image{},
	ScorDiv:   false,
	Splashrot: &Image{},
}

var AudioFX = struct {
	Mute bool
}{
	Mute: false,
}

func hideScorDiv() {
	document.ScorDiv = false
}

func showScorDiv() {
	document.ScorDiv = true
}

// blank bkg display
// not removing previous spc blackout - does this low key leak mem, or does garbage collect clear it?

func gif_blnk(lod *fyne.Container) {
	img := image.NewNRGBA(image.Rect(0, 0, 2000, 2000))
	draw.Draw(img, img.Bounds(), &image.Uniform{HRGB{0}}, image.ZP, draw.Src)
	cimg := canvas.NewRasterFromImage(img)
//	splash.Remove(splim)
	spc := container.NewStack(cimg)
	lod.Add(spc)
}

// load a gif from file fn, set lim as content of container lod, play mus if given
// return true if gif loads

func gif_lodr(fn string, lod, lim *fyne.Container, mus string) bool {

	lded := false
	gif_blnk(lod)	// reg png files are expected to fill splash area, so no need to blank
	gif, err := NewAnimatedGif(storage.NewFileURI(fn))
	if err == nil {
		lod.Remove(lim)
		lim = container.NewStack(gif)
		lod.Add(lim)
	fyne.Do(func() {
		lim.Refresh()
//fmt.Printf("Splash load: %s\n",fn)
	})
		gif.Start()
		lded = true
		if mus != "" { play_sfx(mus) }
	}
	return lded
}

var splash *fyne.Container			// splash intro screens, layout in menu.go
var splim *fyne.Container			// image to splash

func splashrot() {

	sec := false	// first time in play g1 scroller intro w/music
	smpl := ""		// sample play item
	mus := ""		// music with anim, or static even
	srot := 0		// sample play rot
	splashsrc := ""
  for {
	rot := splRot		// def 6000 millis

  if actab == "Game" {		// tab loaded where this happen

	upng := true
// sample play if it didnt play after title, these screens are already done
	if (splCyc == 11 || splCyc == 9) && smpl != "" {
		gif_lodr(smpl, splash, splim, mus)
		smpl = ""
		rot = srot
//fmt.Printf("smpl2: %s\n",rot)
	} else {

	if splCyc == 9 {		// done with g1 splash, load g1 score tbl gfx
		splCyc = 13
	} else {
		if splCyc < 1 || splCyc >= 12 { splCyc = 0 }
		splCyc++
	}

//	if splCyc != 12 { hideScorDiv() }

/*
			vid.Src = "splash/g2samply_q.ogv"
			rot = 119700
			if rand.Float64() < 0.4 {
				vid.Src = "splash/gII_intro.ogv"
				rot = 25200
			}
			if rand.Float64() < 0.27 {
				vid.Src = "splash/gIV_intro.ogv"
				rot = 20650
			}
			if rand.Float64() < 0.22 {
				vid.Src = "splash/gN_intro.ogv"
				rot = 34210  */

	if sec && splCyc == 1 && rand.Float64() > 0.65 { splCyc = 10 }	// after 1st cycle chance to skip from g1 to g2
//	if !sec && splCyc == 1 && rand.Float64() > 0.05 { splCyc = 10 }	// after 1st cycle chance to skip from g1 to g2

// add g1 & 2 smpl gifs & musics, later other intro sets

	if (splCyc == 2 || splCyc == 11) && smpl != "" && rand.Float64() < 0.47 {	// chance for sample play after scroller
		splCyc--	// go back one, hold advance for sample
		gif_lodr(smpl, splash, splim, mus)
		smpl = ""
		rot = srot
//fmt.Printf("smpl1: %s\n",rot)
		upng = false
	} else {		// skip anim splash since cyc goes back to 1 or 10
	if splCyc == 1 || splCyc == 10 || splCyc == 11 {
		splashsrc = fmt.Sprintf("splash/splash%s.gif",string(splLoop[splCyc]))
		rot = 9700			// unless playing 18 secs of music g1, or 25.14 secs g2, or 14 secs ...B.gif
		smpl = "splash/g1smpl.gif"; srot = 43900
		if splCyc == 10 { smpl = "splash/g2smpl.gif"; srot = 122200 }
		if splCyc == 11 { rot = 15000; smpl = "" }
		if (splCyc == 1 && rand.Float64() < 0.71) || !sec { rot = 18700; mus = "sfx/music.title_sf.ogg" }
		if (splCyc == 10 && rand.Float64() < 0.73) { rot = 25160; mus = "sfx/music.g2.title.ogg" }
		upng = !gif_lodr(splashsrc, splash, splim, mus)
		mus = ""
	} else {
		splashsrc = "splash/splash" + string(splLoop[splCyc]) + ".png"
//fmt.Printf("Splash disp: %s\n",splashsrc)
	}}
	if upng {
	err, spl, _ := itemGetPNG(splashsrc)
		if err == nil {
			splash.Remove(splim)
			splim = container.NewStack(spl)
			splash.Add(splim)
		fyne.Do(func() {
			splim.Refresh()
		})
		} else { fmt.Printf("Splash screen fail: %s\n",splashsrc);fmt.Print(err) }
	}
// show score tbl on 12, 13
	if splCyc >= 12 {
//fmt.Printf("Scores: %s\n",splashsrc)
		if rand.Float64() > 0.9 {		// this skips title scroller, strait into g1 ghosts pg
			splCyc = 1
		}
		showScorDiv()
	}}

	sec = true		// second loop+
  }
	time.Sleep(time.Duration(rot) * time.Millisecond)
  }
}

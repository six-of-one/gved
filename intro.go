package main

import (
	"fmt"
	"math/rand"
	"image"
	"image/draw"
	"time"
	"strings"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/canvas"
)

// splash screen rotator
// no video for the time being... he doesnt like them. time biengs are like that

/*
org seq: G¹ sampl, leg, monst, cred, scores, scrolr
		 G² scrolr, sampl, leg, monst, cred, scores

option controls:
slow or fast sample play or none
random chance for each splash set (unless # 5)
1. orig
2. sampl after scrolr
3. sampl between monsters & score
4. mixed up splash set
5. entire load randomized
*/

var splRot = 6000
var splCyc = -1
var splsubCyc = 0

var splLoop = []string{

	"sfx/music.title_sf.ogg",	// 0 unit of loop is music
	"splash/splash1.gif",		// 1
	"splash/splash2.png",
	"splash/splash3.png",
	"splash/splash4.png",
	"splash/splash5.png",
	"splash/splash6.png",
	"splash/splash7.png",
	"splash/splash8.png",
	"splash/splash9.png",
	"splash/g1smplsf.gif",		// B demo play, suporfaster
	"splash/splashD.png",		// scores
//	"splash/g1smplf.gif",		// faster demo play
//	"splash/g1smpl.gif",		// normal speed demo play
	"",		// end of splash set

	"sfx/music.g2.title.ogg",
	"splash/splashA.gif",		// title scroller
	"splash/splashB.gif",		// legend & monsters combine
	"splash/g2smplsf.gif",		// demo play suporfaster
	"splash/splashC.png",		// scores
//	"splash/g2smplf.gif",
//	"splash/g2smpl.gif",
	"",		// end of splash set

	"sfx/z_elec1.ogg",
	"splash/splashSE1.gif",
	"splash/splashSE2.png",
	"splash/splashSE3.png",
	"splash/splashSE4.png",
	"splash/splashSE5.png",
	"splash/splashSE6.png",
	"splash/splashSE7.png",
	"splash/splashSE8.png",
// no demo play yet
	"",		// end of splash set

//	"splash/splashE.png",
//	"splash/splashF.png",
//	"splash/splashK.png",

}

// timing for loops

var splTim = []int{
	18700,				// 0 unit music
	3000,//9700,				// 1 unit time without music
	2000,//6000,				// 2 unit - legend
	1000,//6000,
	1000,//6000,
	1000,//6000,
	1000,//6000,
	1000,//6000,
	1000,//6000,
	1000,//6000,				// 9 unit - theif closes out G¹ monsters
	1000,//26550,				// B unit - demo play (suporfaster)
//	38970,				// faster demo play
//	43930,				// normal speed demo play
	9000,				// A unit - scores
	-1,

	25160,				// unit '0' music	... 13
	3000,//9700,
	15000,				// time for legend + monsters
	1000,//72610,				// demo play suporfaster
//	108490,
//	122200,
	9000,				// scores
	-1,

	66020,				// unit '0' music sanctuary relec1		19
	9700,
	7000,
	7000,
	7000,
	7000,
	7000,
	7000,
	7000,
	-1,
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
		gif.Start()					// sometime - search this out in gif.go and find out what it does
		lded = true
		if mus != "" { play_sfx(mus) }
	} else { fmt.Printf("gif lod issue: %s\n",fn); fmt.Println(err)}
	return lded
}

var splash *fyne.Container			// splash intro screens, layout in menu.go
var splim *fyne.Container			// image to splash

func splashrot() {

	sec := false	// first time in play G¹ scroller intro w/music
//	smpl := ""		// sample play item
	mus := ""		// music with anim, or static even
//	srot := 0		// sample play rot
	splashsrc := ""
	ip := -1		// splash set in play
var	sset = []int{0,13,19}			// start of unit
var	pmus = []float64{0.71,0.33,0.33}	// music percent play
  for {
	rot := splRot		// def 6000 millis
// TESTING
//	if splCyc == -11 { splCyc = 11; smpl = "splash/g2smplsf.gif"; srot = 72390 }
// TESTING

  if actab == "Game" {		// tab loaded where this happen, set cyc to -1 for game run

	upng := true

// new sequence player
	if ip < 0 {		// select a set
		rs := rng.Intn(2)
		ip = sset[rs]
		if !sec { ip = 0; rs = 0 }
		if rand.Float64() < pmus[rs] || !sec { mus = splLoop[ip]; rot = splTim[ip] }
	}
// do splsubCyc here...
  if splsubCyc == 0 {
	ip++	// get next splash, or incr past music
	if splLoop[ip] == "skip" { ip++ }	// not doing sample play
	splashsrc = splLoop[ip]
	if mus == "" { rot = splTim[ip] }
	if rot > 0 {
		if strings.Contains(splashsrc, ".gif") {
			upng = !gif_lodr(splashsrc, splash, splim, mus)
			mus = ""
		}
		if upng {
		err, spl, hsc := itemGetPNG(splashsrc)
			if err == nil {
				if rot == 9000 {
					highscores(hsc,splash,splim)
					if ip == 17 && splsubCyc == 0 { splsubCyc = 18 }
				} else {
					splash.Remove(splim)
					splim = container.NewStack(spl)
					splash.Add(splim)
				fyne.Do(func() {
					splim.Refresh()
				})
				}
			} else { fmt.Printf("Splash screen fail: %s\n",splashsrc);fmt.Print(err) }
		}
	} else {
		ip = -1
	}}
/*
// sample play if it didnt play after title, these screens are already done
	if (splCyc == 11 || splCyc == 9) && smpl != "" {
		gif_lodr(smpl, splash, splim, mus)
		smpl = ""
		rot = srot
//fmt.Printf("smpl2: %s\n",rot)
	} else {

	if splsubCyc == 0 {
	if splCyc == 9 {		// done with G¹ splash, load G¹ score tbl gfx
		hsct[1].msb = fmt.Sprintf("./splash/splD1.png")
		splCyc = 13
	} else {
		if splCyc < 1 || splCyc >= 12 { splCyc = 0 }
		splCyc++
	}}

if splCyc == 17 && splsubCyc == 0 { splsubCyc = 18 }		// replace with ops

	if sec && splCyc == 1 && rand.Float64() > 0.65 { splCyc = 10 }	// after 1st cycle chance to skip from G¹ to G²
// testing
//	if !sec && splCyc == 1 && rand.Float64() > 0.05 { splCyc = 10 }	// after 1st cycle chance to skip from G¹ to G² - test mode
//fmt.Printf("cyc: %d subcyc: %d\n",splCyc,splsubCyc)

// add G¹ & 2 smpl gifs & musics, later other intro sets

	if (splCyc == 2 || splCyc == 11) && smpl != "" && rand.Float64() < 0.47 {	// chance for sample play after scroller
		splCyc--	// go back one, hold advance for sample
		gif_lodr(smpl, splash, splim, mus)
		smpl = ""
		rot = srot
//fmt.Printf("smpl1: %s\n",rot)
		upng = false
	} else {		// skip anim splash since cyc goes back to 1 or 10
	if splCyc == 1 || splCyc == 10 || splCyc == 11 {
		splashsrc = splLoop[splCyc]	//fmt.Sprintf("splash/splash%s.gif",string(splLoop[splCyc]))
		rot = 9700			// unless playing 18 secs of music G¹, or 25.14 secs G², or 14 secs ...B.gif
//		smpl = "splash/g1smpl.gif"; srot = 43930		// std time
//		if splCyc == 10 { smpl = "splash/g2smpl.gif"; srot = 122200 }
//		smpl = "splash/g1smplf.gif"; srot = 38970		// faster samples play
//		if splCyc == 10 { smpl = "splash/g2smplf.gif"; srot = 108490 }
		smpl = "splash/g1smplsf.gif"; srot = 26550		// suporfaster samples play
		if splCyc == 10 { smpl = "splash/g2smplsf.gif"; srot = 72510 }
		if splCyc == 11 { rot = 15000; smpl = "" }
		if (splCyc == 1 && rand.Float64() < 0.71) || !sec { rot = 18700; mus = "sfx/music.title_sf.ogg" }
		if (splCyc == 10 && rand.Float64() < 0.73) { rot = 25160; mus = "sfx/music.g2.title.ogg" }
		upng = !gif_lodr(splashsrc, splash, splim, mus)
		mus = ""
	} else {
		splashsrc = splLoop[splCyc]		//"splash/splash" + string(splLoop[splCyc]) + ".png"
//fmt.Printf("Splash disp: %s\n",splashsrc)
	}}
	if upng {
	err, spl, hsc := itemGetPNG(splashsrc)
		if err == nil {
			if splCyc >= 12 {
				highscores(hsc,splash,splim)
			} else {
				splash.Remove(splim)
				splim = container.NewStack(spl)
				splash.Add(splim)
			fyne.Do(func() {
				splim.Refresh()
			})
			}
		} else { fmt.Printf("Splash screen fail: %s\n",splashsrc);fmt.Print(err) }
	}
// show score tbl on 12, 13
	if splCyc >= 12 && splsubCyc == 0 {
		if rand.Float64() > 0.9 {		// this skips title scroller, strait into G¹ ghosts pg
			splCyc = 1
		}
	}}
*/
	sec = true		// second loop+
  } else {
		splsubCyc = 0		// not in game tab
	}
	if splsubCyc > 0 {			// G² flash high score colors test
		splsubCyc--
		rot = 111
		hsct[1].msb = fmt.Sprintf("./splash/splD%1d.png",(splsubCyc & 3)+1)		// these need to be splD?g2.png and scores need rearranged to match
		time.Sleep(333 * time.Millisecond)
	} else {
		time.Sleep(time.Duration(rot) * time.Millisecond)
	}
  }
}

/* stat save from se vars vids
			vid.Src = "splash/g2samply_q.ogv"
			rot = 119700
			vid.Src = "splash/gII_intro.ogv"
			rot = 25200
			vid.Src = "splash/gIV_intro.ogv"
			rot = 20650
			vid.Src = "splash/gN_intro.ogv"
			rot = 34210
*/

// key press will not take effect until timeout
// the only way to aborgate this is a shorter timeout and some kind of counter

func splash_keytyp(r rune) {

	switch r {

// call up high score table
	case 'S','s':
		splCyc = 0
		err, _, hsc := itemGetPNG("splash/splashD.png")
		if err == nil {
			highscores(hsc,splash,splim)
		}
	}
// start a game!

}
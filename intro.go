package main

import (
	"fmt"
	"math/rand"
//	"strconv"
//	"strings"
	"time"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
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

func splashrot() {
// 	vid := document.IntroVid
//	spl := document.Splash

	// shutdown until game over
/*	if spl.Visibility == "hidden" {
		return
	}
*/
	sec := false
  for {
	rot := splRot		// def 6000 millis

	if splCyc < 1 || splCyc > 12 { splCyc = 0 }
	splCyc++

	if splCyc != 12 { hideScorDiv() }

/*
		vid.Src = "splash/g1samply_q.ogv"
		rot = 44400 // 44 secs - len of g1 sample play

		if rand.Float64() < 0.3 {
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

		if sec && splCyc == 1 && rand.Float64() > 0.65 { splCyc = 10 }

// add g1 & 2 smpl gifs & musics, later other intro sets

		upng := true
		mus := ""
		if splCyc == 1 || splCyc == 10 || splCyc == 11 {
			document.Splashrot.Src = fmt.Sprintf("splash/splash%s.gif",string(splLoop[splCyc]))
			rot = 7700			// unless playing 18 secs of music g1, or 25.14 secs g2, or 14 secs ...B.gif
			if splCyc == 11 { rot = 15000 }
			if (splCyc == 1 && rand.Float64() < 0.21) || !sec { rot = 18100; mus = "sfx/music.title_sf.ogg" }
			if (splCyc == 10 && rand.Float64() < 0.13) || sec { rot = 25160; mus = "sfx/music.g2.title.ogg" }
			gif, err := NewAnimatedGif(storage.NewFileURI(document.Splashrot.Src))
			if err == nil {
				splash.Remove(splim)
				splim = container.NewStack(gif)
				splash.Add(splim)
			fyne.Do(func() {
				splim.Refresh()
fmt.Printf("Splash load: %s\n",document.Splashrot.Src)
			})
				gif.Start()
				upng = false
			}
		} else {
			document.Splashrot.Src = "splash/splash" + string(splLoop[splCyc]) + ".png"
		}

		if splCyc == 12 {
			if rand.Float64() > 0.9 {
				splCyc = 2
			} else {
				if rand.Float64() < 0.38 {
					document.Splashrot.Src = "splash/splash" + string(splLoop[splCyc]) + "2.png"
				}
				showScorDiv()
			}
		}
		if upng {
		err, spl, _ := itemGetPNG(document.Splashrot.Src)
			if err == nil {
				splash.Remove(splim)
				splim = container.NewStack(spl)
				splash.Add(splim)
			fyne.Do(func() {
				splim.Refresh()
			})
			} else { fmt.Printf("Splash screen fail: %s\n",document.Splashrot.Src);fmt.Print(err) }
		}
//	}

	sec = true		// second loop
//	time.AfterFunc(time.Duration(rot)*time.Millisecond, splashrot)
	time.Sleep(time.Duration(rot) * time.Millisecond)
  }
}
/*
func main() {
	// Initialize splash visibility for testing
	document.Splash.Visibility = "visible"
	// Start splash rotation
	splashrot()

	// Prevent main from exiting immediately
	select {}
}
*/
package main

/*
notes:

a. update mazedumps with new visuals from gved - and is done
- consider removing these due to technicality of ownership of visual

*/

import (
	"fmt"
	"image/png"
	"os"
	"os/exec"
	"math/rand"
	"time"

	"git.kirsle.net/go/audio/sdl"			// audio package
	"github.com/veandco/go-sdl2/mix"
)

type TileLinePlane []byte

type TileLinePlaneSet [][]byte

type TileLineMerged []byte

type TileData []TileLineMerged

// sound system

var sfx *sdl.Engine
var aud bool		// true if audio loads

func play_sfx(snd string) {
	if aud && !opts.Mute {
		music, err := sfx.LoadMusic(snd)
		if err == nil {
				music.Play(0)		// arg 0 is loops cnt
		} else {
			if opts.Verbose { fmt.Printf("Audio failure\n%v\n",err) }
		}
	}
}

// indicate which maze set to decode
var G1 bool
var G2 bool
// override the maze address selection by slapstic table
// this is mostly for research, some address will crash gved
var Aov int

// for the user select demo
var Ovwallpat int
var Ovflorpat int
var Ovwallcol int
var Ovflorcol int

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func dotile(tile int) {
	img := genimage(tile, opts.DimX, opts.DimY)
	f, _ := os.OpenFile(opts.Output, os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	// gif.Encode(f, img, &gif.Options{NumColors: 16})
	png.Encode(f, img)
}

func main() {
	args := gevinit()

// audio stuffs
sfx, err := sdl.New(mix.INIT_MP3 | mix.INIT_OGG)
    if err != nil {
        fmt.Printf("Audio fail, Error %v\n",err)
    } else {
		aud = true
		sfx.Setup()
		defer sfx.Teardown()
	}

// new retool - G1 is gauntlet maze, G2 is gauntlet 2 maze
	G1 = false
	G2 = false
// override slapstic maze table lookup address
	Aov = 0

	if opts.Gtp == 2 {
		eid = "Gauntlet II"
		G2 = true
	} else {
		eid = "Gauntlet "
		if opts.R14 { eid += fmt.Sprintf("(rev 14)")
		} else { eid += fmt.Sprintf("(rev 1 - 9)") }
		G1 = true
	}
	fmt.Print(eid+"\n")

	if opts.Intr { fmt.Printf("Audio: %t\n\n",aud) }

	if opts.Addr > 0x37fff && opts.Addr < 0x40000 { Aov = opts.Addr }

	source := rand.NewSource(time.Now().UnixNano()) // random #s
	rng = rand.New(source)

	switch runType {
	case TypeNone:
		if opts.Tile > 0 {
			dotile(opts.Tile)
			fmt.Println("dotile \n")
		} else {
			lv := maxint(0,minint(opts.Lvl,8))
			if lv == 0 {
				if opts.Lvl != -1 {
					ld_config()			// if not forced rnd (-1), see if there is a setting
					lv = opts.Lvl
				}
				if opts.Lvl < 1 {
// put select research here
				lv = rng.Intn(6) + 1 }	// 1 - 7, or 0 select = rnd
			}
			st := fmt.Sprintf("maze%d",lv)
			if opts.Intr { domaze(st) } else {		// set interactive but left out maze# - do it by default
				fmt.Println("nothing selected - more options required, try:\n./gved -i maze1\n./gved floor0\n./gved wall0\n./gved item-dragon-ipotion\n       (./gved item-all for list)\nnote: non-interactive generates output.png\n")
// do a 'help'
				a := "./gved"
				a0 := "-h"
				cmd := exec.Command(a, a0)
				stdout, err := cmd.Output()
				if err != nil {
					fmt.Println(err.Error())
				}
				fmt.Printf("\n")
				fmt.Println(string(stdout))
				os.Exit(1)
			}
		}
	case TypeFloor:
		dofloor(args[0])
	case TypeWall:
		dowall(args[0])
	case TypeMonster:
		domonster(args[0])
	case TypeItem:
		doitem(args[0])
	case TypeMaze:
		domaze(args[0])
	}

	// if opts.Floor >= 0 {
	// 	t := floorStamps[opts.Floor]
	// 	img := genimage_fromarray(t, 2, 2)
	// 	f, _ := os.OpenFile(opts.Output, os.O_WRONLY|os.O_CREATE, 0600)
	// 	defer f.Close()
	// 	gif.Encode(f, img, &gif.Options{NumColors: 16})
	// } else if opts.Wall >= 0 {
	// 	t := wallStamps[opts.Wall]
	// 	img := genimage_fromarray(t, 2, 2)
	// 	f, _ := os.OpenFile(opts.Output, os.O_WRONLY|os.O_CREATE, 0600)
	// 	defer f.Close()
	// 	gif.Encode(f, img, &gif.Options{NumColors: 16})
	// } else if opts.Animate == true {
	// 	t := monsters[opts.Monster].anims["walk"]["upright"]
	// 	x := monsters[opts.Monster].xsize
	// 	y := monsters[opts.Monster].ysize
	// 	imgs := genanim(t, x, y)

	// 	f, _ := os.OpenFile(opts.Output, os.O_WRONLY|os.O_CREATE, 0600)
	// 	defer f.Close()

	// 	var delays []int
	// 	for i := 0; i < len(t); i++ {
	// 		delays = append(delays, 15)
	// 	}

	// 	gif.EncodeAll(f,
	// 		&gif.GIF{
	// 			Image: imgs,
	// 			Delay: delays,
	// 		},
	// 	)
	// } else {
	//  t := opts.Tile
	//  img := genimage(t, opts.DimX, opts.DimY)
	//  f, _ := os.OpenFile(opts.Output, os.O_WRONLY|os.O_CREATE, 0600)
	//  defer f.Close()
	//  gif.Encode(f, img, &gif.Options{NumColors: 16})
	//	}
}

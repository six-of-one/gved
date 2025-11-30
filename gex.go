package main

import (
	"fmt"
	"image/png"
	"os"
)

type TileLinePlane []byte

type TileLinePlaneSet [][]byte

type TileLineMerged []byte

type TileData []TileLineMerged

var G1 int
var G2 int

// for the user select demo
var Ovwallpat int
var Ovflorpat int
var Ovwallcol int
var Ovflorcol int

// FIXME: change name to something not "numbers"
type Stamp struct {
	width   int
	numbers []int
	data    []TileData
	ptype   string
	pnum    int
	trans0  bool
	nudgex  int
	nudgey  int
}

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
	args := gexinit()

// if gauntlet 1 is target, intended to be maze address in rom
	G1 = 0
// can override g2 address load from maze slapstic
	G2 = 0

	switch runType {
	case TypeNone:
		if opts.Tile > 0 {
			dotile(opts.Tile)
		} else {
			fmt.Println("Missing or incorrect identity line.")
			os.Exit(1)
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

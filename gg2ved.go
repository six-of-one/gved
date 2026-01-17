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
)

type TileLinePlane []byte

type TileLinePlaneSet [][]byte

type TileLineMerged []byte

type TileData []TileLineMerged

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

type Stamp struct {
	width   int
	numbers []int
	data    []TileData
	ptype   string
	pnum    int
	trans0  bool
	nudgex  int
	nudgey  int
	mask	int
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
	args := gevinit()

// new retool - G1 is gauntlet maze, G2 is gauntlet 2 maze
	G1 = false
	G2 = false
// override slapstic maze table lookup address
	Aov = 0

	if opts.Gtp == 2 {
		fmt.Printf("Gauntlet II\n")
		G2 = true
	} else {
		fmt.Printf("Gauntlet ")
		if opts.R14 { fmt.Printf("(rev 14)")
		} else { fmt.Printf("(rev 1 - 9)") }
		fmt.Printf("\n")
		G1 = true
	}

	if opts.Addr > 0x37fff && opts.Addr < 0x40000 { Aov = opts.Addr }

	switch runType {
	case TypeNone:
		if opts.Tile > 0 {
			dotile(opts.Tile)
			fmt.Println("dotile \n")
		} else {
			if opts.Intr { domaze("maze1") } else {		// set interactive but left out maze# - do it by default
				fmt.Println("nothing selected - more options required, try:\n./gved -i maze1\n./gved floor0\n./gved wall0\nnote: non-interactive generates output.png\n")
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

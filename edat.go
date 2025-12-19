package main

import (
	"fmt"
	"strings"
	"os"
	"io/ioutil"
)

/*
this is the start of a basic buffer editor
more complexity will be required for:
 a. undo system
 b. sanctuary (g3) mazes
*/

var ebuf MazeData
var eflg [11]int

// save maze to file in .ed

func sav_maz(fil string, mdat MazeData, fdat [11]int, mx int, my int) {
// edit settings
// 1. edit status (1|0) max_x max_y
// 2. 11 bytes of compressed maze lead in - all stats
// 3+ maze data

	file, err := os.Create(fil)
	if err == nil {
//	wfs := fmt.Sprintf("%d\n%d %d %d %d\n%0x\n%#b\n%d %d\n",1,Ovwallpat,Ovflorpat,Ovwallcol,Ovflorcol,maze.secret,maze.flags,lastx,lasty)
		wfs := fmt.Sprintf("%d %d %d\n",opts.edat,mx,my)

		for y := 0; y < 11; y++ {
			wfs += fmt.Sprintf(" %02X", fdat[y])
		}
		wfs += "\n"
		for y := 0; y <= my; y++ {
			for x := 0; x <= mx; x++ {

				wfs += fmt.Sprintf(" %02d", mdat[xy{x, y}])
			}
			wfs += "\n"
		}
		file.WriteString(wfs)
		file.Close()
	} else {
		fmt.Printf("saving maze %s, %d x %d, error:\n",fil,mx,my)
		fmt.Print(err)
	}
}

// load stored maze data into ebuf / eflg

func lod_maz(fil string) int {
	data, err := ioutil.ReadFile(fil)
	edp := 0
	if err == nil {
		var esc int
		dscan := fmt.Sprintf("%s",data)
		fmt.Sscanf(dscan,"%d %d %d\n",&edp,&opts.DimX,&opts.DimY)
		for y := 0; y < 11; y++ {
			fmt.Sscanf(dscan," %02X", &eflg[y])
		}
		fmt.Sscanf(dscan,"\n")
		for y := 0; y <= opts.DimX; y++ {
			for x := 0; x <= opts.DimY; x++ {

				fmt.Sscanf(dscan," %02d", &esc)
				ebuf[xy{x, y}] = esc
			}
			fmt.Sscanf(dscan,"\n")
		}
	} else {
		fmt.Printf("loading maze %s, error:\n",fil)
		fmt.Print(err)
	}
	return edp
}

func stor_maz(mazn int) {

	var lastx int
	var lasty int
	var maze *Maze
	fmt.Printf("buffer maze entry\n")

// if g1 or g2 edit, get size & control bytes
// g3 will be edit of sanctuary mazes
	if opts.Gtp < 3 {
		maze = mazeDecompress(slapsticReadMaze(mazn - 1), false)
		lastx = 32
		if maze.flags&LFLAG4_WRAP_H > 0 {
			lastx = 31
		}
		lasty = 32
		if maze.flags&LFLAG4_WRAP_V > 0 {
			lasty = 31
		}
	}

	fil := fmt.Sprintf(".ed/g%dmaze%03d.ed",opts.Gtp,mazn)

	data, err := ioutil.ReadFile(fil)
	if err != nil {
		errs := fmt.Sprintf("%v",err)
		fmt.Print(errs)
// file does not exist yet
		if strings.Contains(errs, "no such file") {
// editor overs
			maze.optbyts[5] = (Ovflorpat & 0x0f) << 4 + (Ovwallpat & 0x0f)
			maze.optbyts[6] = (Ovflorcol & 0x0f) << 4 + (Ovwallcol & 0x0f)
			for y := 0; y < 11; y++ {
				eflg[y] = maze.optbyts[y]
			}
			ebuf = make(map[xy]int)

			opts.DimX = lastx
			opts.DimY = lasty
			sav_maz(fil, ebuf, eflg, lastx, lasty)
		} else {
		fmt.Print(err)
		return
		}
	}

	fmt.Printf("buffer: %s\n",data)

}

func ed_sav(mazn int) {

	fil := fmt.Sprintf(".ed/g%dmaze%03d.ed",opts.Gtp,mazn)
	sav_maz(fil, ebuf, eflg, opts.DimX, opts.DimY)
}
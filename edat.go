package main

import (
	"fmt"
	"strings"
	"os"
	"io/ioutil"
)

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
	}
}

func stor_maz(mazn int) {

	var lastx int
	var lasty int
	var maze *Maze
	fmt.Printf("buffer maze entry\n")

// if g1 or g2 edit, get size & control bytes
	if opts.Gtp < 3 {
		maze = mazeDecompress(slapsticReadMaze(opts.mnum), false)
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
		if strings.Contains(errs, "no such file") {

// editor overs
			maze.optbyts[5] = (Ovflorpat & 0x0f) << 4 + (Ovwallpat & 0x0f)
			maze.optbyts[6] = (Ovflorcol & 0x0f) << 4 + (Ovwallcol & 0x0f)
			for y := 0; y < 11; y++ {
				eflg[y] = maze.optbyts[y]
			}
			ebuf = make(map[xy]int)

			sav_maz(fil, ebuf, eflg, lastx, lasty)
		}
	} else {
		fmt.Print(err)
		return
	}

	fmt.Printf("buffer: %v\n",data)

}
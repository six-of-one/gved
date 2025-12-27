package main

import (
	"fmt"
	"strings"
	"os"
	"io/ioutil"
	"bufio"
)

/*
this is the start of a basic buffer editor
more complexity will be required for:
 a. undo system
 b. sanctuary (g3) mazes
*/

var edmaze *Maze
var ebuf MazeData
var ubuf MazeData	// initial load from file, swappable with ebuf on <ctrl-u>
var sd [27]MazeData	// save data buffers - save off maze copies
var sdfl [27][13]int
var sdmax = 27
var sdb int			// current sd selected, -1 when on ebuf
var eflg [11]int

// deleted elements buffer

type Deletebuf struct {
	mx     [1000]int
	my     [1000]int
	elem   [1000]int
}

var delbuf = &Deletebuf{}
var delstak int

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

				wfs += fmt.Sprintf("%02d\n", mdat[xy{x, y}])
			}
//			wfs += "\n"
		}
		file.WriteString(wfs)
		file.Close()
	} else {
		fmt.Printf("saving maze %s, %d x %d, error:\n",fil,mx,my)
		fmt.Print(err)
	}
}

// load stored maze data into ebuf / eflg or other data stores

func lod_maz(fil string, mdat MazeData, fdat [11]int, ud bool) int {

	data, err := ioutil.ReadFile(fil)
	edp := 0
	if err == nil {
		esc := 0
		dscan := fmt.Sprintf("%s",data)
// may not be the optimal way, but it works for now
	    scanr := bufio.NewScanner(strings.NewReader(dscan))
		l := "0 32 32"	// the default on scan failure will produce a solid block of wall 32 x 32
		if scanr.Scan() { l = scanr.Text() }
		fmt.Sscanf(l,"%d %d %d",&edp,&opts.DimX,&opts.DimY)
// keeping the verbose scan track for now
	if opts.Verbose { fmt.Printf("\nscanned:\ned %d, %02d x %02d\n", edp,opts.DimX,opts.DimY) }
		l = " 00 00 00 00 00 00 00 0B 5A 5B 49"
		if scanr.Scan() { l = scanr.Text() }
		fmt.Sscanf(l," %02X %02X %02X %02X %02X %02X %02X %02X %02X %02X %02X\n", &fdat[0], &fdat[1], &fdat[2], &fdat[3], &fdat[4], &fdat[5], &fdat[6], &fdat[7], &fdat[8], &fdat[9], &fdat[10])
	if opts.Verbose {
			for y := 0; y < 11; y++ { fmt.Printf(" %02X", fdat[y]) }
			fmt.Printf("\n")
		}

		if mdat == nil { mdat = make(map[xy]int) }
		if ubuf == nil { ubuf = make(map[xy]int) }
		for y := 0; y <= opts.DimX; y++ {
			for x := 0; x <= opts.DimY; x++ {
				l = "02"
				if scanr.Scan() { l = scanr.Text() }
	if opts.Verbose { fmt.Printf("%02s ",l) }
				fmt.Sscanf(l,"%02d", &esc)
				mdat[xy{x, y}] = esc
				if ud { ubuf[xy{x, y}] = esc }		// store ubuf data on flag
				edp = 1		// tell sender we loaded some maze part
			}
	if opts.Verbose { fmt.Printf("\n") }
		}
	} else {
// this warning will issue if a maze buffer save (maze not being edited) has not happened because and the maze is viewed
		fmt.Printf("loading maze %s, warning:\n",fil)
		fmt.Print(err)
		fmt.Printf("\n")
		fmt.Printf("Note: 'no such file' if maze is not being edited and the maze is viewed when editor is on\n")
	}
	return edp
}

func stor_maz(mazn int) {

	var lastx int
	var lasty int
	var maze *Maze
//	fmt.Printf("buffer maze entry\n")

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
			opts.DimX = lastx
			opts.DimY = lasty
			if ebuf == nil { ebuf = make(map[xy]int) }
			for y := 0; y <= lasty; y++ {
				for x := 0; x <= lastx; x++ {
				ebuf[xy{x, y}] = maze.data[xy{x, y}]
			}}
			sav_maz(fil, ebuf, eflg, lastx, lasty)
		} else {
			fmt.Print(err)
		}
		return
	}

	if false { fmt.Printf("buffer: %s\n",data) }
	
// handle g3 mazes here ?
}

func ed_sav(mazn int) {

	eflg[5] = (Ovflorpat & 0x0f) << 4 + (Ovwallpat & 0x0f)
	eflg[6] = (Ovflorcol & 0x0f) << 4 + (Ovwallcol & 0x0f)
	fil := fmt.Sprintf(".ed/g%dmaze%03d.ed",opts.Gtp,mazn)
	sav_maz(fil, ebuf, eflg, opts.DimX, opts.DimY)
}

// udpate maze from edits
func ed_maze() {
	for y := 0; y <= opts.DimX; y++ {
		for x := 0; x <= opts.DimY; x++ {
		edmaze.data[xy{x, y}] = ebuf[xy{x, y}]
	}}
	Ovimg := genpfimage(edmaze, opts.mnum)
	upwin(Ovimg)
}

// reload maze while editing & update window - generates output.png

func remaze(mazn int) {
fmt.Printf("in remaze\n")
	edmaze = mazeDecompress(slapsticReadMaze(mazn), false)
	mazeloop(edmaze)
	Ovimg := genpfimage(edmaze, mazn)
	upwin(Ovimg)
}
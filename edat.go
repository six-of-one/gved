package main

import (
	"fmt"
	"strings"
	"os"
	"io/ioutil"
	"bufio"
	"encoding/binary"
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

var sdmax = 1000
var sdb int			// current sd selected, -1 when on ebuf
var eflg [11]int
var tflg [14]int	// transfer flags - because they dont pass as a parm?

// deleted elements buffer

type Deletebuf struct {
	mx     [10001]int
	my     [10001]int
	elem   [10002]int
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
	opts.bufdrt = false
}

// load stored maze data into ebuf / eflg or other data stores

func lod_maz(fil string, mdat MazeData, ud bool) int {

	data, err := ioutil.ReadFile(fil)
	edp := 0
	if err == nil {
		esc := 0
// setup for rejecting the load because of dirty buffer flag
//		for y := 0; y < 11; y++ { tflg[y] = eflg[y] }
//		dumpbuf()		// check buffer dirty flag for edits needing saved, only opt here is discard or dont load
//		if opts.bufdrt { return 1 }
		dscan := fmt.Sprintf("%s",data)
// may not be the optimal way, but it works for now
	    scanr := bufio.NewScanner(strings.NewReader(dscan))
		l := "0 32 32"	// the default on scan failure will produce a solid block of wall 32 x 32
		if scanr.Scan() { l = scanr.Text() }
		fmt.Sscanf(l,"%d %d %d",&edp,&opts.DimX,&opts.DimY)
		tflg[12] = opts.DimX
		tflg[13] = opts.DimY
// keeping the verbose scan track for now
	if opts.Verbose { fmt.Printf("\nscanned:\ned %d, %02d x %02d\n", edp,opts.DimX,opts.DimY) }
		l = " 00 00 00 00 00 00 00 0B 5A 5B 49"
		if scanr.Scan() { l = scanr.Text() }
		fmt.Sscanf(l," %02X %02X %02X %02X %02X %02X %02X %02X %02X %02X %02X\n", &tflg[0], &tflg[1], &tflg[2], &tflg[3], &tflg[4], &tflg[5], &tflg[6], &tflg[7], &tflg[8], &tflg[9], &tflg[10])
	if opts.Verbose {
			for y := 0; y < 11; y++ { fmt.Printf(" %02X", tflg[y]) }
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
		edp = -1
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

	fil := fmt.Sprintf(".ed/g%dmaze%03d.ed",opts.Gtp,mazn)
	sav_maz(fil, ebuf, eflg, opts.DimX, opts.DimY)
}

func upd_edmaze() {
	for y := 0; y <= opts.DimX; y++ {
		for x := 0; x <= opts.DimY; x++ {
		edmaze.data[xy{x, y}] = ebuf[xy{x, y}]
	}}
	for y := 0; y < 11; y++ {
		edmaze.optbyts[y] = eflg[y]
	}
	flagbytes := make([]byte, 4)
	flagbytes[0] = byte(eflg[1])
	flagbytes[1] = byte(eflg[2])
	flagbytes[2] = byte(eflg[3])
	flagbytes[3] = byte(eflg[4])
	edmaze.flags = int(binary.BigEndian.Uint32(flagbytes))
	edmaze.wallpattern = Ovwallpat
	edmaze.floorpattern = Ovflorpat
	edmaze.wallcolor = Ovwallcol
	edmaze.floorcolor = Ovflorcol

}
// udpate maze from edits
func ed_maze() {
	upd_edmaze()
	Ovimg := genpfimage(edmaze, opts.mnum)
	upwin(Ovimg)
}

// replaceing or deleting - store for ctrl-z / ctrl-y

func undo_buf(sx int, sy int) {
	maxdel := 10000
	if delstak >= maxdel {		//	we hit max, shift back 1, losing the start of undo
		for i := 0; i < maxdel; i++ {
			delbuf.mx[i] = delbuf.mx[i+1]
			delbuf.my[i] = delbuf.my[i+1]
			delbuf.elem[i] = delbuf.elem[i+1]
		}
		delstak--
	}
	delbuf.mx[delstak] = sx
	delbuf.my[delstak] = sy
	delbuf.elem[delstak] = ebuf[xy{sx, sy}]
	fmt.Printf(" del elem: %d maze: %d x %d\n",delbuf.elem[delstak],delbuf.mx[delstak],delbuf.my[delstak])
	delstak++
	delbuf.elem[delstak] = -1 	// when undeleting this is the end
}

// same as mazeloop, but called by Rr, h, m while cmd keys active in edit mode

func rotmirbuf(rmmaze *Maze) {

	fmt.Printf("in rotmirbuf\n")

// to transform maze, array copy
	xform := make(map[xy]int)
// manual mirror, flip
	if opts.MH || opts.MV || opts.MRP || opts.MRM {

		sx := 1
		lastx := 32
		if rmmaze.flags&LFLAG4_WRAP_H > 0 {
			sx = 0
			lastx = 31
		}

		sy := 1
		lasty := 32
		if rmmaze.flags&LFLAG4_WRAP_V > 0 {
			sy = 0		// otherwise it wont MV correct
			lasty = 31
		}

	fmt.Printf("wraps -- hw: %d vw: %d\n", rmmaze.flags&LFLAG4_WRAP_H,rmmaze.flags&LFLAG4_WRAP_V)
	fmt.Printf(" fx: %d lx %d fy %d ly %d\n", sx,lastx,sy,lasty)


// note it
/*		fmt.Printf("init\n")
	for y := 0; y <= lasty; y++ {
		for x := 0; x <= lastx; x++ {

			fmt.Printf(" %02d", maze.data[xy{x, y}])
		}
		fmt.Printf("\n")
	}
		fmt.Printf("\n")
*/
// transform																										 - rotating sq. wall mazes will always work
// rotate +90 degrees				-- * there is the issue of gauntlet arcade NEEDING the y = 0 wall *always* intact, rotating looper mazes wont work
		if opts.MRP {
			for ty := sy; ty <= lasty; ty++ {
			for tx := sx; tx <= lastx; tx++ {
				xform[xy{lastx - tx, ty}] = rmmaze.data[xy{ty, tx}]
// g1 - must transform all dors on a rotat since they have horiz & vert dependent
				if xform[xy{lastx - tx, ty}] == G1OBJ_DOOR_HORIZ { xform[xy{lastx - tx, ty}] = G1OBJ_DOOR_VERT } else {
				if xform[xy{lastx - tx, ty}] == G1OBJ_DOOR_VERT { xform[xy{lastx - tx, ty}] = G1OBJ_DOOR_HORIZ } }
// g2
				if xform[xy{lastx - tx, ty}] == MAZEOBJ_DOOR_HORIZ { xform[xy{lastx - tx, ty}] = MAZEOBJ_DOOR_VERT } else {
				if xform[xy{lastx - tx, ty}] == MAZEOBJ_DOOR_VERT { xform[xy{lastx - tx, ty}] = MAZEOBJ_DOOR_HORIZ } }
			}}
		} else {
		if opts.MRM {
			for ty := sy; ty <= lasty; ty++ {
			for tx := sx; tx <= lastx; tx++ {
				xform[xy{tx, lasty - ty}] = rmmaze.data[xy{ty, tx}]
// g1
				if xform[xy{tx, lasty - ty}] == G1OBJ_DOOR_HORIZ { xform[xy{tx, lasty - ty}] = G1OBJ_DOOR_VERT } else {
				if xform[xy{tx, lasty - ty}] == G1OBJ_DOOR_VERT { xform[xy{tx, lasty - ty}] = G1OBJ_DOOR_HORIZ } }
// g2
				if xform[xy{tx, lasty - ty}] == MAZEOBJ_DOOR_HORIZ { xform[xy{tx, lasty - ty}] = MAZEOBJ_DOOR_VERT } else {
				if xform[xy{tx, lasty - ty}] == MAZEOBJ_DOOR_VERT { xform[xy{tx, lasty - ty}] = MAZEOBJ_DOOR_HORIZ } }
			}}
		}
		}

// mirror x
		if opts.MH {
			for ty := sy; ty <= lasty; ty++ {
			for tx := sx; tx <= lastx; tx++ {
				xform[xy{lastx - tx, ty}] = rmmaze.data[xy{tx, ty}]
			}}
		}

// mirror y: flip
		if opts.MV {
			for ty := sy; ty <= lasty; ty++ {
			for tx := sx; tx <= lastx; tx++ {
				xform[xy{tx, lasty - ty}] = rmmaze.data[xy{tx, ty}]
			}}
			if rmmaze.flags&LFLAG4_WRAP_V > 0 {	// fix wall not allowed being at bottom for arcade gauntlet
				for ty := lasty - 1; ty >= sy ; ty-- {
				for tx := sx; tx <= lastx; tx++ {
					xform[xy{tx, ty + 1}] = xform[xy{tx, ty}]
				}}
				for tx := sx; tx <= lastx; tx++ { xform[xy{tx, 0}] = G1OBJ_WALL_REGULAR }
			}
		}

// copy back
		for y := sy; y <= lasty; y++ {
			for x := sx; x <= lastx; x++ {
				rmmaze.data[xy{x, y}] = xform[xy{x, y}]
				ebuf[xy{x, y}] = xform[xy{x, y}]
			}
		}
// TEMP maze dmp
		fmt.Printf("dun\n")
	for y := 0; y <= lasty; y++ {
		for x := 0; x <= lastx; x++ {

			fmt.Printf(" %02d", rmmaze.data[xy{x, y}])
		}
		fmt.Printf("\n")
	}
		fmt.Printf("\n")
// REM TEMP

	}
// clear all in edit mode
	opts.MRP = false
	opts.MRM = false
	opts.MV = false
	opts.MH = false
}

// reload maze while editing & update window - generates output.png

func remaze(mazn int) {
fmt.Printf("in remaze\n")
	sdb = -1
	if !opts.dntr {
		edmaze = mazeDecompress(slapsticReadMaze(mazn), false)
		mazeloop(edmaze)
		opts.bufdrt = false
	}
	opts.dntr = false
	Ovimg := genpfimage(edmaze, mazn)
	upwin(Ovimg)
}
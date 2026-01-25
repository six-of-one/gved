package main

import (
	"fmt"
	"strings"
	"os"
	"io/ioutil"
	"bufio"
	"image"
	"encoding/binary"
	"image/color"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

/*
edit and transfer system

this is the start of a basic buffer editor
more complexity will be required for:
 a. undo system
 b. sanctuary (g3) mazes
*/

var edmaze *Maze
var ebuf MazeData	// main edit buffer and corresponding flags
var ecolor color.Color		// master color for maze elements

var sdmax = 1000
var sdb int			// current sd selected, -1 when on ebuf
var eflg [11]int
var tflg [14]int	// transfer flags - because they dont pass as a parm for scan from file?
					//					so after a file load, these have to be copied to the appropriate flags
var din [33]int		// set to be 1 line per std gauntlet maze (gved encoding) of 0 - 32 elements [ with H wrap being 0 - 31 ]

// deleted elements / undo storage

type Deletebuf struct {
	mx     []int
	my     []int
	elem   []int
	revc   []int		// this maze element is part of a multiplace, allow one click removal/ restore
}

var delbuf = &Deletebuf{}
var delstak int
var restak int		// keep track of redo chain

// ulternate buffer - copy of maze from load
var ubuf MazeData	// initial load from file, swappable with ebuf on <ctrl-u>
var uflg [11]int
var udb = &Deletebuf{}	// and with the way the delbuf operates now, ubuf must also swap that
var udstak int
var urstak int

// initialize edit control

func ed_init() {
	anum = 0			// vars numeric inputs
	wpalop = false		// pallete, win active
	palfol = true		// - follows maze wall & floor decor
	ccp = NOP			// paste buffer
	prcl = 1			// paint multi undo ops
	wpbop = false		// window for pb active
	blotter(nil,0,0,0,0)	// init blotter
	ccblot = blot			// rubber band blot, saved so pb can display contents then return to rb
	lg1cnt = 1
	lg2cnt = 1
	sdb = -1			// sd buffer
	cycl = 0			// edit key 'c' cycle ops
	cmdhin = "cmds: ?, eE, fFgG, wWqQ, rRt, hm, pPT, sL, S, il, u, v, A #a"
	delbset(0)			// init undo (delbuf)
	restak = 0			// restor position in delbuf
}

func udbck(ct int, t int){

//fmt.Printf("udb len %d, test: %d\n",len(udb.elem),t)

	if len(udb.elem) <= t {
		for y := 0; y < ct; y++ {
			udb.elem = append(udb.elem,-1)
			udb.mx = append(udb.mx,0)
			udb.my = append(udb.my,0)
			udb.revc = append(udb.revc,1)
		}
	}
}

// set head of delbuf to u, add space if needed & init to -1

func delbset(u int) {
	delstak = u
	delbck(3, u + 1)		// initialize with ? units if len < u + 1
	delbuf.mx[delstak] = 0
	delbuf.my[delstak] = 0
	delbuf.revc[delstak] = 1
	delbuf.elem[delstak] = -1
}

// if delbuf len < t, add ct units

func delbck(ct int, t int){

//fmt.Printf("delbuf st: %d len %d, test: %d\n",delstak,len(delbuf.elem),t)

	if len(delbuf.elem) <= t {
		for y := 0; y < ct; y++ {
			delbuf.elem = append(delbuf.elem,-1)
			delbuf.mx = append(delbuf.mx,0)
			delbuf.my = append(delbuf.my,0)
			delbuf.revc = append(delbuf.revc,1)
		}
	}
}

// pre-pend .db_ to save filename for delete buffer save & load

func prep(fn string) string {
	fl := len(fn)
	rfl := fn
	lstb := 0
// find last /, if any
	for y := 0; y < fl; y++ {
		if rfl[y:y+1] == "/" { lstb = y+1 }
	}
	rfl = fn[0:lstb]+".db_"+fn[lstb:fl]
//fmt.Printf("sv fil: %s last bld: %d\n",rfl, lstb)
	return rfl
}

// init vars buffers

func init_buf() {
	if ebuf == nil { ebuf = make(map[xy]int) }
	if ubuf == nil { ubuf = make(map[xy]int) }
	if cpbuf == nil { cpbuf = make(map[xy]int) }
	if plbuf == nil { plbuf = make(map[xy]int) }
}
// save maze to file in .ed
// add a maze # to saves
// svdb - save undo buf

func sav_maz(fil string, mdat MazeData, fdat [11]int, mx int, my int, smazn int, svdb bool) {
// edit settings
// 1. edit status (1) max_x max_y
// 2. 11 bytes of compressed maze lead in - all stats
// 3+ maze data
if opts.Verbose { fmt.Printf("saving maze %s\n",fil) }

	file, err := os.Create(fil)
	if err == nil {
//	wfs := fmt.Sprintf("%d\n%d %d %d %d\n%0x\n%#b\n%d %d\n",1,Ovwallpat,Ovflorpat,Ovwallcol,Ovflorcol,maze.secret,maze.flags,lastx,lasty)
		wfs := fmt.Sprintf("%d %d %d %d\n",smazn,mx,my,opts.Gtp)

		for y := 0; y < 11; y++ {
			wfs += fmt.Sprintf(" %02X", fdat[y])
		}
		wfs += "\n"
		parse := 32
		for y := 0; y <= my; y++ {
			for x := 0; x <= mx; x++ {
//				wfs += fmt.Sprintf("%02d\n", mdat[xy{x, y}])
				wfs += fmt.Sprintf(" %03d", mdat[xy{x, y}])
				if parse < 1 { wfs += "\n"; parse = 32 } else {
					parse--
				}
			}
		}
// /fmt.Printf("parse: %02d \n",parse)
		if parse != 32 { for y := 0; y <= parse; y++ { wfs += " 999" }}		// pad out to end of 33 unit, so read line back in wont cause crash
		wfs += "\n"
		file.WriteString(wfs)
		file.Close()
	} else {
		fmt.Printf("saving maze %s, %d x %d, error:\n",fil,mx,my)
		fmt.Print(err)
		fmt.Printf("\n")
	}

// now save deleted elements -- set mazn 0 for buffers like paste
	if delstak > 0 && svdb {
		dbf := prep(fil) //fil[0:4]+".db_"+fil[4:len(fil)]
if opts.Verbose { fmt.Printf("saving maze undo %s\n",dbf) }
		file, err := os.Create(dbf)
		if err == nil {
			wfs := fmt.Sprintf("%d %d\n",delstak, restak)

			for y := 0; y < delstak; y++ {
				if delbuf.elem[y] < 0 { break }
				wfs += fmt.Sprintf("%d %d %d %d\n", delbuf.elem[y],delbuf.mx[y],delbuf.my[y],delbuf.revc[y])
			}
			wfs += "\n"
			file.WriteString(wfs)
			file.Close()
		} else {
			fmt.Printf("saving maze undo %s, error:\n",dbf)
			fmt.Print(err)
			fmt.Printf("\n")
		}
	}
//	}
	opts.bufdrt = false
}

// load stored maze data into ebuf / eflg or other data stores
var mazln int		// maze load # stored

func lod_maz(fil string, mdat MazeData, ud bool, ldb bool) int {

if opts.Verbose { fmt.Printf("loading maze %s\n",fil) }

	data, err := ioutil.ReadFile(fil)
	edp := 0
	if err == nil {

		dscan := fmt.Sprintf("%s",data)
// may not be the optimal way, but it works for now
	    scanr := bufio.NewScanner(strings.NewReader(dscan))
		l := "0 32 32"	// the default on scan failure will produce a solid block of wall 32 x 32
		if scanr.Scan() { l = scanr.Text() }
		fmt.Sscanf(l,"%d %d %d",&mazln,&opts.DimX,&opts.DimY)
		tflg[12] = opts.DimX
		tflg[13] = opts.DimY
// keeping the verbose scan track for now
	if opts.Verbose { fmt.Printf("\nscanned:\ned %d, %02d x %02d\n", mazln,opts.DimX,opts.DimY) }
		l = " 00 00 00 00 00 00 00 0B 5A 5B 49"
		if scanr.Scan() { l = scanr.Text() }
		fmt.Sscanf(l," %02X %02X %02X %02X %02X %02X %02X %02X %02X %02X %02X\n", &tflg[0], &tflg[1], &tflg[2], &tflg[3], &tflg[4], &tflg[5], &tflg[6], &tflg[7], &tflg[8], &tflg[9], &tflg[10])
		if ud { for y := 0; y < 11; y++ { uflg[y] = tflg[y] }}
	if opts.Verbose {
			for y := 0; y < 11; y++ { fmt.Printf(" %02X", tflg[y]) }
			fmt.Printf("\n")
		}

		if mdat == nil { mdat = make(map[xy]int) }		// init most bufs used by edit system, most come here anyway

// loop to load - note issue with scans of formatted data
		parse := 33
		for y := 0; y <= opts.DimY; y++ {
			for x := 0; x <= opts.DimX; x++ {

// seems working, now to read all the old and rewrite
// new method to parse line of 33 units
				if parse > 32 { //  1				5					A					F					0					5					A
						l = " 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002"
						if scanr.Scan() { l = scanr.Text() }
						//        0    1                        6                            12                       17                                 24                       29             32
				fmt.Sscanf(l," %03d %03d %03d %03d %03d %03d %03d %03d %03d %03d %03d %03d %03d %03d %03d %03d %03d %03d %03d %03d %03d %03d %03d %03d %03d %03d %03d %03d %03d %03d %03d %03d %03d\n",
					&din[0], &din[1], &din[2], &din[3], &din[4], &din[5], &din[6], &din[7], &din[8], &din[9], &din[10], &din[11], &din[12], &din[13], &din[14], &din[15], &din[16], &din[17], &din[18],
					&din[19], &din[20], &din[21], &din[22], &din[23], &din[24], &din[25], &din[26], &din[27], &din[28], &din[29], &din[30], &din[31], &din[32])
					parse = 0
				}
				if din[parse] < 999 {				// max value is end of buffer fill
					mdat[xy{x, y}] = din[parse]
					if ud { ubuf[xy{x, y}] = din[parse] }		// store ubuf data on flag
	if opts.Verbose { fmt.Printf("%03d ",din[parse]) }
				}
				parse++
				edp = 1		// tell sender we loaded some maze part
			}
	if opts.Verbose { fmt.Printf("\n") }
		}
	} else {
// this warning will issue if a maze buffer save (maze not being edited) has not happened because and the maze is viewed
		if opts.Verbose {
			fmt.Printf("loading maze %s, warning:\n",fil)
			fmt.Print(err)
			fmt.Printf("\n")
			fmt.Printf("Note: 'no such file' if maze is not being edited and the maze is viewed when editor is on\n")
		}
		edp = -1
	}
// now load deleted elements - but not for pb or pal
  if ldb {
	dbf := prep(fil) //fil[0:4]+".db_"+fil[4:len(fil)]
	data, err = ioutil.ReadFile(dbf)
	delstak = 0
	restak = 0
	delbset(0)
	if err == nil {
		dscan := fmt.Sprintf("%s",data)
	    scanr := bufio.NewScanner(strings.NewReader(dscan))
		l := "0"	// the default on scan failure will produce a solid block of wall 32 x 32
		if scanr.Scan() { l = scanr.Text() }
		fmt.Sscanf(l,"%d %d",&delstak, &restak)

		for y := 0; y < delstak; y++ {
			l = "-1 0 0 1"
			if scanr.Scan() { l = scanr.Text() }
			delbck(6, y)
			fmt.Sscanf(l, "%d %d %d %d\n", &delbuf.elem[y],&delbuf.mx[y],&delbuf.my[y],&delbuf.revc[y])
			if ud {
				udbck(6,y)
				udb.mx[y] = delbuf.mx[y]
				udb.my[y] = delbuf.my[y]
				udb.revc[y] = delbuf.revc[y]
				udb.elem[y] = delbuf.elem[y]
			}
			if delbuf.elem[y] < 0 { delstak = y; break }
		}
		delbset(delstak)
		if ud {
			udstak = delstak; urstak = restak
			udbck(3,udstak)
			udb.elem[udstak] = -1
		}

	} else {
		if opts.Verbose {
			fmt.Printf("edp %d failed < 0 or loading maze deleted %s, warning:\n",edp,dbf)
			fmt.Print(err)
			fmt.Printf("\n")
		}
	}
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
		fmt.Printf("\n")
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
			for y := 0; y <= lasty; y++ {
				for x := 0; x <= lastx; x++ {
				ebuf[xy{x, y}] = maze.data[xy{x, y}]
			}}
			sav_maz(fil, ebuf, eflg, lastx, lasty, mazn, false)
			delstak = 0
			restak = 0
			delbset(0)
		} else {
			fmt.Print(err)
			fmt.Printf("\n")
		}
		return
	}

	if false { fmt.Printf("buffer: %s\n",data) }
	
// handle g3 mazes here ?
}

func ed_sav(mazn int) {

	upd_edmaze(false)
	fil := fmt.Sprintf(".ed/g%dmaze%03d.ed",opts.Gtp,mazn)
	sav_maz(fil, ebuf, eflg, opts.DimX, opts.DimY, mazn, true)
}

func upd_edmaze(ovrm bool) {
fmt.Printf("upd_edmaze: x,y: %d, %d\n",opts.DimX,opts.DimY)
	for y := 0; y <= opts.DimY; y++ {
		for x := 0; x <= opts.DimX; x++ {
		edmaze.data[xy{x, y}] = ebuf[xy{x, y}]
	}}
	for y := 0; y < 11; y++ {
		edmaze.optbyts[y] = eflg[y]
	}
	if wpalop && palfol { palete() }
	flagbytes := make([]byte, 4)
	flagbytes[0] = byte(eflg[1])
	flagbytes[1] = byte(eflg[2])
	flagbytes[2] = byte(eflg[3])
	flagbytes[3] = byte(eflg[4])
	edmaze.flags = int(binary.BigEndian.Uint32(flagbytes))
	if ovrm {
		edmaze.wallpattern = eflg[5] & 0x0f
		edmaze.floorpattern = (eflg[5] & 0xf0) >> 4
		edmaze.wallcolor = eflg[6] & 0x0f
		edmaze.floorcolor = (eflg[6] & 0xf0) >> 4
	} else {
		edmaze.wallpattern = Ovwallpat
		edmaze.floorpattern = Ovflorpat
		edmaze.wallcolor = Ovwallcol
		edmaze.floorcolor = Ovflorcol
	}
}
// udpate maze from edits - rld false to keep overload colors / pats
func ed_maze(rld bool) {
	upd_edmaze(rld)
	Ovimg := genpfimage(edmaze, opts.mnum)
	upwin(Ovimg)
	calc_stats()
}

// replaceing or deleting - store for ctrl-z / ctrl-y

func undo_buf(sx int, sy int, rc int) {
//	fmt.Printf(" del %d elem: %d\n",delstak,delbuf.elem[delstak])
	delbuf.mx[delstak] = sx
	delbuf.my[delstak] = sy
	delbuf.revc[delstak] = rc					// revoke count for the loop
	delbuf.elem[delstak] = ebuf[xy{sx, sy}]
// append the next unit blank if needed
//fmt.Printf(" del %d elem: %d maze: %d x %d - rloop: %d\n",delstak,delbuf.elem[delstak],delbuf.mx[delstak],delbuf.my[delstak],rc)
	delstak++
	restak = delstak		// placing or deleting one breaks restore chain
	delbset(delstak)
//fmt.Printf(" del %d elem: %d\n",delstak,delbuf.elem[delstak])
}

/*
actual rotate maths

0	1	2	3

1	2	3	4		0
5	6	7	8		1
9	10	11	12		2

0	1	2

9	5	1		0		+90
10	6	2		1
11	7	3		2
12	8	4		3

4	8	12		0		-90
3	7	11		1
2	6	10		2
1	5	9		3
		from  =	to
			y = x		y = rx
			x = ry		x = y

0,0 -> 2,0	+90	-> 0,3  -90
1,0 -> 2,1		-> 0,2
2,0 -> 2,2		-> 0,1
3,0 -> 2,3		-> 0,0

0,1 -> 1,0		-> 1,3
1,1 -> 1,1		-> 1,2
2,1 -> 1,2		-> 1,1
3,1 -> 1,3		-> 1,0

0,2 -> 0,0		-> 2,3
1,2 -> 0,1		-> 2,2
2,2 -> 0,2		-> 2,1
3,2 -> 0,3		-> 2,0

*/
// the actual werker, so we can use it on cpbuf, etc

func rotmirmov(mdat MazeData, sx int, sy int, lastx int, lasty int, flg int) (int, int) {

// to transform maze, array copy
	xform := make(map[xy]int)
// transform																										 - rotating sq. wall mazes will always work
// rotate +90 degrees				-- * there is the issue of gauntlet arcade NEEDING the y = 0 wall *always* intact, rotating looper mazes wont work
		if opts.MRP {
			for ty := sy; ty <= lasty; ty++ {
			for tx := sx; tx <= lastx; tx++ {
				xform[xy{lasty - ty, tx}] = mdat[xy{tx, ty}]
// g1 - must transform all dors on a rotat since they have horiz & vert dependent
				if xform[xy{lasty - ty, tx}] == G1OBJ_DOOR_HORIZ { xform[xy{lasty - ty, tx}] = G1OBJ_DOOR_VERT } else {
				if xform[xy{lasty - ty, tx}] == G1OBJ_DOOR_VERT { xform[xy{lasty - ty, tx}] = G1OBJ_DOOR_HORIZ } }
// g2
				if xform[xy{lasty - ty, tx}] == MAZEOBJ_DOOR_HORIZ { xform[xy{lasty - ty, tx}] = MAZEOBJ_DOOR_VERT } else {
				if xform[xy{lasty - ty, tx}] == MAZEOBJ_DOOR_VERT { xform[xy{lasty - ty, tx}] = MAZEOBJ_DOOR_HORIZ } }
			}}
			if lastx != lasty { sw := lastx; lastx = lasty; lasty = sw }		// on a rotate in edit when size x != size y, they must swap after the rot
		} else {
		if opts.MRM {
			for ty := sy; ty <= lasty; ty++ {
			for tx := sx; tx <= lastx; tx++ {
				xform[xy{ty, lastx - tx}] = mdat[xy{tx, ty}]
// g1
				if xform[xy{ty, lastx - tx}] == G1OBJ_DOOR_HORIZ { xform[xy{ty, lastx - tx}] = G1OBJ_DOOR_VERT } else {
				if xform[xy{ty, lastx - tx}] == G1OBJ_DOOR_VERT { xform[xy{ty, lastx - tx}] = G1OBJ_DOOR_HORIZ } }
// g2
				if xform[xy{ty, lastx - tx}] == MAZEOBJ_DOOR_HORIZ { xform[xy{ty, lastx - tx}] = MAZEOBJ_DOOR_VERT } else {
				if xform[xy{ty, lastx - tx}] == MAZEOBJ_DOOR_VERT { xform[xy{ty, lastx - tx}] = MAZEOBJ_DOOR_HORIZ } }
			}}
			if lastx != lasty { sw := lastx; lastx = lasty; lasty = sw }
		}
		}

// mirror x
		if opts.MH {
			for ty := sy; ty <= lasty; ty++ {
			for tx := sx; tx <= lastx; tx++ {
				xform[xy{lastx - tx, ty}] = mdat[xy{tx, ty}]
			}}
		}

// mirror y: flip
		if opts.MV {
			for ty := sy; ty <= lasty; ty++ {
			for tx := sx; tx <= lastx; tx++ {
				xform[xy{tx, lasty - ty}] = mdat[xy{tx, ty}]
			}}
			if flg&LFLAG4_WRAP_V > 0 {	// fix wall not allowed being at bottom for arcade gauntlet
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
				mdat[xy{x, y}] = xform[xy{x, y}]
			}
		}

// clear all in edit mode
	opts.MRP = false
	opts.MRM = false
	opts.MV = false
	opts.MH = false

	return lastx, lasty
}

// same as mazeloop, but called by Rr, h, m while cmd keys active in edit mode
// 	╚══> except in this buffer is changed by ops

func rotmirbuf(rmmaze *Maze) (int, int) {

	fmt.Printf("in rotmirbuf\n")

// manual mirror, flip
	sx := 0
	sy := 0

	lastx := 32
	if rmmaze.flags&LFLAG4_WRAP_H > 0 {
		lastx = 31
	}

	lasty := 32
	if rmmaze.flags&LFLAG4_WRAP_V > 0 {
		lasty = 31
	}

	fmt.Printf("wraps -- hw: %d vw: %d\n", rmmaze.flags&LFLAG4_WRAP_H,rmmaze.flags&LFLAG4_WRAP_V)
	fmt.Printf("rotmirbuf fx: %d lx %d fy %d ly %d\n", sx,lastx,sy,lasty)

	lastx, lasty = rotmirmov(rmmaze.data,sx,sy,lastx,lasty,rmmaze.flags)

	for y := sy; y <= lasty; y++ {
		for x := sx; x <= lastx; x++ {
			ebuf[xy{x, y}] = rmmaze.data[xy{x, y}]
		}
	}
	return lastx, lasty
}

// reload maze while editing & update window - generates output.png

func remaze(mazn int) {
fmt.Printf("\n\nin remaze dntr: %t edat:%d sdb: %d, delstk: %d\n",opts.dntr,opts.edat,sdb,delstak)
	if !opts.dntr {
		sdb = -1
		delbset(0)
		edmaze = mazeDecompress(slapsticReadMaze(mazn), false)
		mazeloop(edmaze)
		opts.bufdrt = false
	}
	opts.dntr = false
	nsremaze = false
	if opts.edat > 0 { ed_maze(true) } else {
		Ovimg := genpfimage(edmaze, mazn)
		upwin(Ovimg)
		calc_stats()
	}
}

// turn on edit mode

func edit_on(k int) {
if opts.edat == 0 {
		smod = "Edit mode: "
fmt.Printf("editor on, maze: %03d\n",opts.mnum+1)
		opts.edat = 1
		stor_maz(opts.mnum+1)	// this does not auto store new edit mode to buffer save file, unless it creates the file
		statlin(cmdhin,"")
// these all deactivate as override during edit
		opts.MRM = false
		opts.MRP = false
		opts.MV = false
		opts.MH = false
	}
// activate keys & select k (edkdef from mb click)
	if k > 0 {
		if !cmdoff { typedRune('\\') }	// turn cmd keys off
		typedRune(rune(k))
	}
}
// valid check, edit key

func valid_keys(ek int) int {
	if ek > maxkey || ek < minkey { return edkdef }		// 33 to 126, outside this return 121 'y'
	return ek
}
// palette

// statistics on mazes, already set for partial sanctuary expansion

var g1stat [1000]int
var g2stat [1000]int
var stonce [1000]int	// on;y display a stat once

// bring up edit palette after saving
var wpalop bool		// is the pb win open?
var wpal fyne.Window // is the pal win open?
var plbuf MazeData	// initial load from file, swappable with ebuf on <ctrl-u>
var plflg [11]int
var palxs int
var palys int
var palfol bool		// palette decor follows map chg

func palete() {

	nm := 0
	pmx := opts.DimX; pmy := opts.DimY
	for my := 0; my <= pmy; my++ {
	for mx := 0; mx <= pmx; mx++ { plbuf[xy{mx, my}] = 0 }}
	fil := fmt.Sprintf(".ed/sd%05d_g%d.ed",nm,opts.Gtp)
	cnd := lod_maz(fil, plbuf, false, false)
	cpx = opts.DimX; cpy = opts.DimY

	if cnd >= 0 { for y := 0; y < 11; y++ { plflg[y] =  tflg[y] };
		if palfol { for y := 0; y < 11; y++ { plflg[y] =  eflg[y] }}
		bwin(cpx+1, cpy+1, 0, plbuf, plflg, "pal") }
	opts.DimX = pmx; opts.DimY = pmy
	wpal.SetCloseIntercept(func() {
		statsB = nil
	})
}

// typer for pal win

var statsB binding.Item[string]

func palRune(r rune) {

	switch r {
		case '?': dboxtx("palette ops", "in gved main window:\n"+
						"middle click an element\n(click activates edit + default key)\n... or ...\n"+
						"select maze\nhit <ESC> - activate edit mode\nhit '\\' for edit keys\n"+
						"hit a key to map: 'y'\nmove mouse to maze or palette\nand middle click an element\n"+
						"\nin palette window:\nmiddle click an element\n\n"+
						"edit hint on menu bar give status"+
						"\n\npal win keys:\nq,Q - quit\nt,T - hide flags info\n\n(only when window active)\n"+
						"*stats in terminal if palette open\ns, S - stats window open, auto updates\n"+
						"f, F - gauntlet 2 maze flags list\nx, X - gauntlet 2 secret tricks", 350,540)
		case 't': fallthrough
		case 'T': dboxtx("T hide flags", "in gved main window:\n\n"+
						"invisible flag set - hide vars maze elements:\n"+
						" T - cycle through a flag set (loop 0 - 511)\n"+
						" #T - set flags = # ---- <ctrl>-T reset flags to 0\n\n"+
						" s  - (show only) random 'special potions' &\n"+
						"      gold bags on emply floor tiles (not EDIT)\n"+
						" L  - toggle generator indicator letters [ DGLS ]\n"+
						"      showing box gen monster\n"+
						" p  - toggle floor invisible *\n"+
						" P  - toggle walls invisible *\n\n"+
						"NOGEN = 1	// all generators\n"+
						"NOMON = 2		// all monster, dragon\n"+
						"NOFUD = 4		// all food\n"+
						"NOTRS = 8		// treas, locked\n"+
						"NOPOT = 16		// pots & t.powers\n"+
						"NODOR = 32	// doors, keys\n"+
						"NOTRAP = 64	// trap & floor dots, stun, ff tiles\n"+
						"NOEXP = 128	// exit, push wall\n"+
						"NOTHN = 256	// anything else left\n"+
						"NOFLOOR = 512\n"+
						"NOWALL = 1024	// g2 *walls\n"+
						"NOG1W = 2048	// g1 std wall only\n\n"+
						"set # with:\nBlank maze (file menu)\n- keep items flags cover\n\n"+
						"Random profile load\n- only load items flags cover"+
						"\n\n* hide items disabled when edit keys active",440,700)
		case 'f': fallthrough
		case 'F': dboxtx("G2 maze flags","     Gauntlet 2 flags                hex value        bit pos\n"+
						"ODDANGLE_GHOSTS	= 0x01000000 - 00000001\n"+
						"ODDANGLE_GRUNTS	= 0x02000000 - 00000010\n"+
						"ODDANGLE_DEMONS	= 0x04000000 - 00000100\n"+
						"ODDANGLE_LOBBERS	= 0x08000000 - 00001000\n"+
						"ODDANGLE_SORCERERS	= 0x10000000 - 00010000\n"+
						"ODDANGLE_Aux_Grnts	= 0x20000000 - 00100000\n"+
						"ODDANGLE_DEATHS	= 0x40000000 - 01000000\n"+
						"INVIS_TRAPWALLS	= 0x80000000 - 10000000\n\n"+

						"FAST_GHOSTS		= 0x010000      - 00000001\n"+
						"FAST_GRUNTS		= 0x020000      - 00000010\n"+
						"FAST_DEMONS		= 0x040000      - 00000100\n"+
						"FAST_LOBBERS		= 0x080000      - 00001000\n"+
						"FAST_SORCERERS		= 0x100000      - 00010000\n"+
						"FAST_AUX_GRUNTS	= 0x200000      - 00100000\n"+
						"FAST_DEATHS		= 0x400000      - 01000000\n"+
						"INVIS_ALLWALLS		= 0x800000      - 10000000\n\n"+

						"RANDOMFOOD_MASK	= 0x0700           - 00000111\n"+
						"WALLS_CYCLIC		= 0x0800           - 00001000\n"+
						"WALLS_DELETABLE1	= 0x1000           - 00010000\n"+
						"WALLS_DELETABLE2	= 0x2000           - 00100000\n"+
						"EXIT_MOVES		= 0x4000           - 01000000\n"+
						"EXIT_CHOOSEONE	= 0x8000           - 10000000\n\n"+

						"SHOTS_STUN		= 0x01                - 00000001\n"+
						"SHOTS_HURT		= 0x02                - 00000010\n"+
						"TRAPS_LOCAL		= 0x04                - 00000100\n"+
						"TRAPS_RANDOM		= 0x08                - 00001000\n"+
						"WRAP_V			= 0x10                - 00010000\n"+
						"WRAP_H			= 0x20                - 00100000\n"+
						"EXIT_FAKE			= 0x40                - 01000000\n"+
						"PLAYER_OFFSCREEN	= 0x80                - 10000000\n\n"+
						"flag math is binary 'or' \"|\" together\n"+
						"bit position is within byte of that flag portion",420,744)
		case 'x': fallthrough
		case 'X': dboxtx("G2 secret tricks","     Gauntlet 2 secret room tricks\n"+
						"NONE 		= 0x00  \"No trick\"\n"+
						"TRANSPORT1 	= 0x01  \"Try Transportability (onto death)\"\n"+
						"TRANSPORT2	= 0x02  \"Try Transportability (onto death)\"\n"+
						"TRANSPORT3	= 0x03  \"Try Transportability (into exit)\"\n"+
						"TRANSPORT4	= 0x04  \"Try Transportability (onto secret wall)\"\n"+
						"WATCHSHOOT1	= 0x05  \"Watch What You Shoot (shoot foods)\"\n"+
						"WATCHSHOOT	= 0x06  \"Watch What You Shoot (shoot secret walls)\"\n"+
						"SAVESUPERSHOTS = 0x07 \"Save Super Shots\"\n"+
						"NOUSEINVUL	= 0x08  \"Don't Use Invulnerability\"\n"+
						"NOGETHIT		= 0x09  \"Don't Get Hit (while killing a dragon)\"\n"+
						"PUSHWALL		= 0x0a  \"Don't Be Fooled\"\n"+
						"NOFOOLED		= 0x0b  \"Don't Be Fooled\"\n"+
						"NOGREEDY1	= 0x0c  \"Don't Be Greedy (no keys or potions)\"\n"+
						"NOGREEDY2	= 0x0d  \"Don't Be Greedy (no treasure)\"\n"+
						"DIET			= 0x0e  \"Go On a Diet (no food)\"\n"+
						"BEPUSHY		= 0x0f   \"Be Pushy\"\n"+
						"IT 			= 0x10  \"IT Could Be Nice\"\n"+
						"NOHURTFRIENDS	= 0x11  \"Don't Hurt Friends\"",535,435)
		case 's': fallthrough
		case 'S': statsB = dboxtx("Maze stats","",340,700)
		case 'q': fallthrough
		case 'Q': if wpalop { statsB = nil; wpalop = false; wpal.Close() }
		default:
	}
}

// assign keys

func key_asgn(buf MazeData, ax int, ay int) {

	if G1 {
		g1edit_keymap[edkey] = buf[xy{ax, ay}]
		kys := g1mapid[g1edit_keymap[edkey]]
		keyst := fmt.Sprintf("G¹ assn key: %s = %03d, %s",map_keymap[edkey],g1edit_keymap[edkey],kys)
		statlin(cmdhin,keyst)
		play_sfx(g1auds[g1edit_keymap[edkey]])
		if edkey == cycloc { cycl = g1edit_keymap[cycloc] }		// when reassign 'c' key, set cycl as well
	} else {
		g2edit_keymap[edkey] = buf[xy{ax, ay}]
		kys := g2mapid[g2edit_keymap[edkey]]
		keyst := fmt.Sprintf("G² assn key: %s = %03d, %s",map_keymap[edkey],g2edit_keymap[edkey],kys)
		statlin(cmdhin,keyst)
		play_sfx(g2auds[g2edit_keymap[edkey]])
		if edkey == cycloc { cycl = g1edit_keymap[cycloc] }
	}
}

// stat package for palette win

func zero_stat() {

	for y := 0; y < 1000; y++ { g1stat[y] = 0; g2stat[y] = 0; stonce[y] = 1 }
}

// count stuff

func stats(elm int) {

	if G1 { g1stat[elm]++ }
	if G2 { g2stat[elm]++ }
}

// and displate it

func calc_stats() {
	if wpalop {
		if palfol { palete() }
		zero_stat()
fmt.Printf("get stats: %d %d\n",opts.DimX,opts.DimY)
		for y := 0; y <= opts.DimY; y++ {
			for x := 0; x <= opts.DimX; x++ {
			stats(ebuf[xy{x, y}])
		}}
// stats during palette
			if opts.Verbose { fmt.Printf("stats:\n") }
		stl := ""
		if G1 {
		for y := 0; y <= 65; y++ { if g1stat[y] > 0 {
			if opts.Verbose { fmt.Printf("  %s: %d\n",g1mapid[y],g1stat[y]) }
			stl += fmt.Sprintf("  %s: %d\n",g1mapid[y],g1stat[y])
		}}}
		if G2 {
		for y := 0; y <= 65; y++ { if g2stat[y] > 0 {
			if opts.Verbose { fmt.Printf("  %s: %d\n",g2mapid[y],g2stat[y]) }
			stl += fmt.Sprintf("  %s: %d\n",g2mapid[y],g2stat[y])
		}}}
		stl += "══════════════════════\n"
		stl += mazeMetaPrint(edmaze, true)
		if statsB != nil { statsB.Set(stl) }
//		bwin(palxs, palys, 0, plbuf, plflg)
	}
}

// cut / copy & paste

var cpbuf MazeData	// c/c/p buffer
var pbcnt int		// master count of c/c/p buffers saved
var lpbcnt int		// sesssion count of c/c/p buffers - reset every
var cpx int			// max paste buf, start is always 0, 0
var cpy int
// roll thru pb
var masbcnt int		// run thru master pb
var sesbcnt int		// run thru local ses pb
var lg1cnt int		// ses pb save for g1 maps
var lg2cnt int		// ses pb save for g2 maps

// i've discovered a 'local' in function version of these will crash, this prob needs to be a struct
var wpbop bool		// is the pb win open?
var wpb fyne.Window	// win to view pastbuf contents
var wpbimg *image.NRGBA

// get paste buffer cnt each init

func get_pbcnt() {
	pbcnt = 1
	lpbcnt = 1
	if G1 { lpbcnt = lg1cnt }
	if G2 { lpbcnt = lg2cnt }
	masbcnt = 1	// loop thru
	sesbcnt = 1
	fil := fmt.Sprintf(".pb/cnt_g%d",opts.Gtp)
	data, err := ioutil.ReadFile(fil)
	if err == nil {
		fmt.Sscanf(string(data),"%d", &pbcnt)
if opts.Verbose { fmt.Printf("pbcnt: %d, ses cnt: %d\n",pbcnt,lpbcnt) }
	}
}

func pb_upd(id string, nt string, vl int) {
// clear old buf
	pmx := opts.DimX; pmy := opts.DimY		// preserve these
	for my := 0; my <= pmy; my++ {
	for mx := 0; mx <= pmx; mx++ { cpbuf[xy{mx, my}] = 0 }}
	fil := fmt.Sprintf(".pb/%s_%07d_g%d.ed",id,vl,opts.Gtp)
	lod_maz(fil, cpbuf, false, false)
	cpx = opts.DimX; cpy = opts.DimY
fmt.Printf("%spb dun: px %d py %d, %s\n",nt,cpx,cpy,fil)
	opts.DimX = pmx; opts.DimY = pmy
/*
for my := 0; my <= cpy; my++ {
	for mx := 0; mx <= cpx; mx++ {
fmt.Printf("%03d ",cpbuf[xy{mx, my}])
	}
fmt.Printf("\n")
}*/

	bwin(cpx+1, cpy+1, vl, cpbuf, eflg, id)		// draw the buffer
	bl := fmt.Sprintf("paste buf: %d", vl)
	statlin(cmdhin,bl)
}

// transforms on cpbuf in smol window

func pb_loced(cnt int) {
	fil := fmt.Sprintf(".pb/pb_%07d_g%d.ed",cnt,opts.Gtp)
	sav_maz(fil, cpbuf, eflg, cpx, cpy, 0, false)
	pb_upd("pb", "mas", cnt)
}

// typer for pb win

func pbRune(r rune) {

	switch r {
		case '?': dboxtx("paste buffer viewer keys", "q,Q - quit\n"+
						"o,p - cycle session pb - / +\nO,P - cycle master pb - / +\n"+
						", . - cycle master pb - / +\n══════════════════"+
						"\nr - rotate pb +90°\nR - rotate pb -90°\nh - horiz flip\nm - vert mirror"+
						"\n══════════════════\nmiddle click selects element\n"+
						"left click sets element\n(pb has basic edit support)\n"+
						"- pb edit autosaves\n══════════════════\n(* only when window active)", 250,350)			// showing in main win because pb win is usually too small
		case ',': pbmas_cyc(-1)
		case '.': pbmas_cyc(1)
		case 'O': pbmas_cyc(-1)
		case 'P': pbmas_cyc(1)
		case 'o': pbsess_cyc(-1)
		case 'p': pbsess_cyc(1)
// chunk size control
		case 'a': cpx--; if cpx < 1 { cpx = 1 }; pb_loced(masbcnt)
		case 'd': cpx++; pb_loced(masbcnt)
		case 'w': cpy--; if cpy < 1 { cpy = 1 }; pb_loced(masbcnt)
		case 's': cpy++; pb_loced(masbcnt)
// some extras
		case 'r': opts.MRP = true
			cpx,cpy = rotmirmov(cpbuf,0,0,cpx,cpy,0)
			pb_loced(masbcnt)
		case 'R': opts.MRM = true
			cpx,cpy = rotmirmov(cpbuf,0,0,cpx,cpy,0)
			pb_loced(masbcnt)
		case 'h': opts.MH = true
			rotmirmov(cpbuf,0,0,cpx,cpy,0)
			pb_loced(masbcnt)
		case 'm': opts.MV = true
			rotmirmov(cpbuf,0,0,cpx,cpy,0)
			pb_loced(masbcnt)
		case 'q': fallthrough
		case 'Q': if wpbop { wpbop = false; wpb.Close() }
		default:
	}
}

// display a  buffer window with buffer contents - no edit on palette
// px, py - size of paste buffer from 0, 0
// bn - buffer #

func bwin(px int, py int, bn int, mbuf MazeData, fdat [11]int, id string) {

var lw fyne.Window	// local cpy win to view buf contents
  wt := "palette selector"
  nimg := segimage(mbuf,fdat,px,py,false) // (bn == 0)) stats off for now
  dt := float32(opts.dtec)
	if (bn > 0) {
	if !wpbop {
		wpbop = true
		wpb = a.NewWindow(" pbf")
		wpb.Canvas().SetOnTypedRune(pbRune)
		wpb.SetCloseIntercept(func() {
			if blot != ccblot { blot.Hide(); blot = ccblot };	// rb blotter back
			wpbop = false; wpb.Close()})
		wpb.Resize(fyne.NewSize(float32(px)*dt, float32(py)*dt))		// have to do this on new win
		wpb.Show()
	}
// change pb blotter if active
	wpbimg = nimg										// for blotter overlay on ctrl-p
	if wpbop && ccp == PASTE {
		blotup = true
		blot.Resize(fyne.NewSize(float32(px)*dt, float32(py)*dt))
	}
	lw = wpb
	wt = fmt.Sprintf("%d %sf",bn,id)

  } else {	// palette or some other win
	if !wpalop  && bn == 0 {
		wpalop = true
		wpal = a.NewWindow("palette selector")
		wpal.Canvas().SetOnTypedRune(palRune)
		wpal.SetCloseIntercept(func() {wpalop = false;wpal.Close()})
		wpal.Resize(fyne.NewSize(float32(px*32), float32(py*32)))		// have to do this on new win
		wpal.Show()
	}
	palxs = px
	palys = py
	lw = wpal
  }
	clikwins(lw, nimg, px, py)
	lw.SetTitle(wt)
	lw.Resize(fyne.NewSize(float32(px)*dt, float32(py)*dt))

// if opts.Verbose {
fmt.Printf("clkwin sz: %v\n",fyne.NewSize(float32(px)*dt, float32(py)*dt))
}

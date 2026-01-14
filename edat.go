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
var ebuf MazeData	// main edit buffer and corresponding flags

var sdmax = 1000
var sdb int			// current sd selected, -1 when on ebuf
var eflg [11]int
var tflg [14]int	// transfer flags - because they dont pass as a parm for scan from file?
					//					so after a file load, these have to be copied to the appropriate flags
var din [33]int		// set to be 1 line per std gauntlet maze (gved encoding) of 0 - 32 elements [ with H wrap being 0 - 31 ]

// cut / copy & paste

var cpbuf MazeData	// c/c/p buffer
var pbcnt int		// count of c/c/p buffers saved
var cpx int			// max paste buf, start is always 0, 0
var cpy int

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

func udbck(ct int, t int){

fmt.Printf("udb len %d, test: %d\n",len(udb.elem),t)

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
fmt.Printf("sv fil: %s last bld: %d\n",rfl, lstb)
	return rfl
}

// save maze to file in .ed
// add a maze # to saves

func sav_maz(fil string, mdat MazeData, fdat [11]int, mx int, my int, smazn int) {
// edit settings
// 1. edit status (1) max_x max_y
// 2. 11 bytes of compressed maze lead in - all stats
// 3+ maze data

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

// now save deleted elements
	if delstak > 0 {
		dbf := prep(fil) //fil[0:4]+".db_"+fil[4:len(fil)]
fmt.Printf("saving maze delete %s\n",dbf)
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
			fmt.Printf("saving maze deleted %s\n",dbf)
			fmt.Print(err)
			fmt.Printf("\n")
		}
	}
//	}
	opts.bufdrt = false
}

// load stored maze data into ebuf / eflg or other data stores
var mazln int		// maze load # stored

func lod_maz(fil string, mdat MazeData, ud bool) int {

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
		if ubuf == nil { ubuf = make(map[xy]int) }
		if cpbuf == nil { cpbuf = make(map[xy]int) }
// loop to load - note issue with scans of formatted data
		parse := 33
		for y := 0; y <= opts.DimX; y++ {
			for x := 0; x <= opts.DimY; x++ {

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
		fmt.Printf("loading maze %s, warning:\n",fil)
		fmt.Print(err)
		fmt.Printf("\n")
		fmt.Printf("Note: 'no such file' if maze is not being edited and the maze is viewed when editor is on\n")
		edp = -1
	}
// now load deleted elements
	dbf := prep(fil) //fil[0:4]+".db_"+fil[4:len(fil)]
	data, err = ioutil.ReadFile(dbf)
	delstak = 0
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
		fmt.Printf("edp %d failed < 0 or loading maze deleted %s, warning:\n",edp,dbf)
		fmt.Print(err)
		fmt.Printf("\n")
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
			if ebuf == nil { ebuf = make(map[xy]int) }
			for y := 0; y <= lasty; y++ {
				for x := 0; x <= lastx; x++ {
				ebuf[xy{x, y}] = maze.data[xy{x, y}]
			}}
			sav_maz(fil, ebuf, eflg, lastx, lasty, mazn)
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
	sav_maz(fil, ebuf, eflg, opts.DimX, opts.DimY, mazn)
}

func upd_edmaze(ovrm bool) {
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
// udpate maze from edits - rld false to keep colors / pats
func ed_maze(rld bool) {
	upd_edmaze(rld)
	Ovimg := genpfimage(edmaze, opts.mnum)
	upwin(Ovimg)
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

// same as mazeloop, but called by Rr, h, m while cmd keys active in edit mode
// 	╚══> except in this buffer is changed by ops

func rotmirbuf(rmmaze *Maze) {

	fmt.Printf("in rotmirbuf\n")

// to transform maze, array copy
	xform := make(map[xy]int)
// manual mirror, flip
	sx := 0
	lastx := 32
	if rmmaze.flags&LFLAG4_WRAP_H > 0 {
		sx = 0
		lastx = 31
	}

	sy := 0
	lasty := 32
	if rmmaze.flags&LFLAG4_WRAP_V > 0 {
		sy = 0		// otherwise it wont MV correct
		lasty = 31
	}

	fmt.Printf("wraps -- hw: %d vw: %d\n", rmmaze.flags&LFLAG4_WRAP_H,rmmaze.flags&LFLAG4_WRAP_V)
	fmt.Printf("rotmirbuf fx: %d lx %d fy %d ly %d\n", sx,lastx,sy,lasty)


// note it
/*		fmt.Printf("init\n")
	for y := 0; y <= lasty; y++ {
		for x := 0; x <= lastx; x++ {

			fmt.Printf(" %02d", rmmaze.data[xy{x, y}])
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
/*		fmt.Printf("rm dun\n")
	for y := 0; y <= lasty; y++ {
		for x := 0; x <= lastx; x++ {

			fmt.Printf(" %02d", rmmaze.data[xy{x, y}])
		}
		fmt.Printf("\n")
	}
		fmt.Printf("\n")*/
// REM TEMP

// clear all in edit mode
	opts.MRP = false
	opts.MRM = false
	opts.MV = false
	opts.MH = false
}

// bring up edit palette after saving
var swsdb int		// palette on a sd edit

func palete() {

	if opts.bufdrt { menu_savit(true) }		// autosave
// editing an sd, save & palette
	if sdb > 0 { swsdb = sdb; sdb = 0 } else {
// palette up and we got back to sdb
	if swsdb > 0 { sdb = swsdb; swsdb = 0 } else {
// palette maze swaps
	if sdb == 0 {
		menu_lodit(true)		// also also wik - sets sdb -1
	} else {
		sdb = 0
	}}}

	if sdb > -1 {
		nm := sdb
		fil := fmt.Sprintf(".ed/sd%05d_g%d.ed",nm,opts.Gtp)
		cnd := lod_maz(fil, ebuf, false)
		if cnd >= 0 { sdb = nm; for y := 0; y < 11; y++ { eflg[y] =  tflg[y] }; ed_maze(true) }
	}
}

// reload maze while editing & update window - generates output.png

func remaze(mazn int) {
fmt.Printf("in remaze dntr: %t edat:%d sdb: %d, delstk: %d\n",opts.dntr,opts.edat,sdb,delstak)
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
	}
}
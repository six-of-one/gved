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
	"math/rand"
	"time"
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
var ebuf MazeData		// main edit buffer and corresponding flags
type Xdat map[xy]string	// extra data store
var xbuf Xdat
var ecolor color.Color	// master color for maze elements
var eid string			// id string for titles

var sdmax = 1000
var sdb int			// current sd selected, -1 when on ebuf
var eflg [14]int
var tflg [14]int	// transfer flags - because they dont pass as a parm for scan from file?
					//					so after a file load, these have to be copied to the appropriate flags
var din [33]int		// set to be 1 line per std gauntlet maze (gved encoding) of 0 - 32 elements [ with H wrap being 0 - 31 ]

// deleted elements / undo storage

type Deletebuf struct {
	mx     []int
	my     []int
	elem   []int
	xbfd   []string
	revc   []int		// this maze element is part of a multiplace, allow one click removal/ restore
}

var delbuf = &Deletebuf{}
var delstak int
var restak int		// keep track of redo chain

// ulternate buffer - copy of maze from load
var ubuf MazeData	// initial load from file, swappable with ebuf on <ctrl-u>
var xubf Xdat
var uflg [14]int
var udb = &Deletebuf{}	// and with the way the delbuf operates now, ubuf must also swap that
var udstak int
var urstak int

// viewport
var viewp int
var vpx int
var vpy int
var maxvp = 32
var minvp = 12

var unpinx bool		// in gauntlet - whether you approach and edge wall, or stay in the center of screen
var unpiny bool		//				 this doesnt happen until much later, somewhere after level 50, or 80 or so

var rng *rand.Rand

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
//	diff_level = 1.0	// default diff, for now only rload uses - in options now
	source := rand.NewSource(time.Now().UnixNano()) // random #s
	rng = rand.New(source)
	zmod = 0			// play test mode
	zm_x = -1
	zm_y = -1
}

// turn on edit mode

func edit_on(k int) {
if opts.edat == 0 {
		smod = "Edit mode: "
fmt.Printf("editor on, maze: %03d or sd: %d\n",opts.mnum+1, sdb)
		opts.edat = 1
		statlin(cmdhin,"viewport")
// these all deactivate as override during edit
		opts.MRM = false
		opts.MRP = false
		opts.MV = false
		opts.MH = false
	}
	if xbline == nil {
		xbline = dboxtx("xb-line", "0000000", 512, 60,xbl_cls,nil)
	}
// activate keys & select k (edkdef from mb click)
	if k > 0 {
		if !cmdoff { typedRune('\\') }	// turn cmd keys off
		typedRune(rune(k))
	}
}

// txt dialog to store & retr xbuf data

var xbline binding.Item[string]
var xblchg string					// detect changes in xbline

// box typer so edits go into xb edit key

func xbl_typ() {

	if xbline != nil {		// this open input needs validated to hex string, no spaces
		nv, _ := xbline.Get()
		if nv != xblchg {
			g1edit_xbmap[valid_keys(edkey)] = nv
			sv_config()
			xblchg = nv
fmt.Printf("xbline key: %d = %s\n",valid_keys(edkey),g1edit_xbmap[valid_keys(edkey)])
		}
	}
}

// close out edit line for se exp

func xbl_cls() {

	xbline = nil
}

func udbck(ct int, t int){

//fmt.Printf("udb len %d, test: %d\n",len(udb.elem),t)

	if len(udb.elem) <= t {
		for y := 0; y < ct; y++ {
			udb.elem = append(udb.elem,-1)
			udb.xbfd = append(udb.xbfd,"00")
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
	delbuf.xbfd[delstak] = "00"
	delbuf.elem[delstak] = -1
}

// if delbuf len < t, add ct units

func delbck(ct int, t int){

//fmt.Printf("delbuf st: %d len %d, test: %d\n",delstak,len(delbuf.elem),t)

	if len(delbuf.elem) <= t {
		for y := 0; y < ct; y++ {
			delbuf.elem = append(delbuf.elem,-1)
			delbuf.xbfd = append(delbuf.xbfd,"00")
			delbuf.mx = append(delbuf.mx,0)
			delbuf.my = append(delbuf.my,0)
			delbuf.revc = append(delbuf.revc,1)
		}
	}
}

// pre-pend $prp to save filename for delete buffer save & load

func prep(fn string, prp string) string {
	fl := len(fn)
	rfl := fn
	lstb := 0
// find last /, if any
	for y := 0; y < fl; y++ {
		if rfl[y:y+1] == "/" { lstb = y+1 }
	}
	rfl = fn[0:lstb]+prp+fn[lstb:fl]
//fmt.Printf("sv fil: %s last bld: %d\n",rfl, lstb)
	return rfl
}

// init vars buffers

func init_buf() {
	if ebuf == nil { ebuf = make(map[xy]int) }
	if xbuf == nil { xbuf = make(map[xy]string) }
	if ubuf == nil { ubuf = make(map[xy]int) }
	if xubf == nil { xubf = make(map[xy]string) }
	if cpbuf == nil { cpbuf = make(map[xy]int) }
	if xcpbuf == nil { xcpbuf = make(map[xy]string) }
	if plbuf == nil { plbuf = make(map[xy]int) }
	if xplb == nil { xplb = make(map[xy]string) }
}

// clear mazedata buf, max size mx x my, fill with z, unless wh is set > -66, then only replace wh

func clr_buf(buf MazeData, xdat Xdat, mx int, my int, z int, wh int) {
	if wh < -65 {		// if we dont have a when = wh, set size for possible larger edited mazes
		de := 256
		if mx < de { mx = de }
		if my < de { my = de }
	}
	for y := 0; y <= my; y++ {
		for x := 0; x <= mx; x++ {
			if wh < -65 { buf[xy{x, y}] = z; xdat[xy{x, y}] = "0"
		} else {
			if buf[xy{x, y}] == wh { buf[xy{x, y}] = z; xdat[xy{x, y}] = "0" }
		}
	}}
}

// save maze to file in .ed
// add a maze # to saves
// svdb - save undo buf

func sav_maz(fil string, xdat Xdat, mdat MazeData, fdat [14]int, mx int, my int, smazn int, svdb bool) {
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

// save the xtra buffer sanctuary engine edit/perf data
	xbf := prep(fil, "xb_")
	file, err = os.Create(xbf)
		if err == nil {
			wfs := fmt.Sprintf("%d %d\n",mx,my)	// size of buf

			for i := 0; i < curwf; i++ { wfs += fmt.Sprintf("%s %s\n",wlfl.florn[i],wlfl.walln[i]); fmt.Printf("%d: %s %s\n",i,wlfl.florn[i],wlfl.walln[i]) }
			wfs += fmt.Sprintf("xwfdn\n")	// end of floors / walls
// custom walls here, read until done
			for y := 0; y <= my; y++ {	// store one line per element due to ops
			for x := 0; x <= mx; x++ {
				wfs += fmt.Sprintf("%s\n", xdat[xy{x, y}])
			}}
			file.WriteString(wfs)
			file.Close()
		} else {
			fmt.Printf("saving maze xtra buf %s, error:\n",xbf)
			fmt.Print(err)
			fmt.Printf("\n")
		}

// now save deleted elements -- set mazn 0 for buffers like paste
	if delstak > 0 && svdb {
		dbf := prep(fil, ".db_") //fil[0:4]+".db_"+fil[4:len(fil)]
if opts.Verbose { fmt.Printf("saving maze undo %s\n",dbf) }
		file, err := os.Create(dbf)
		if err == nil {
			wfs := fmt.Sprintf("%d %d\n",delstak, restak)

			for y := 0; y < delstak; y++ {
				if delbuf.elem[y] < 0 { break }
				wfs += fmt.Sprintf("%d %d %d %d %s\n", delbuf.elem[y],delbuf.mx[y],delbuf.my[y],delbuf.revc[y],delbuf.xbfd[y])
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

func lod_maz(fil string, xdat Xdat, mdat MazeData, ud bool, ldb bool) int {

if opts.Verbose { fmt.Printf("loading maze %s\n",fil) }

	data, err := ioutil.ReadFile(fil)
	edp := 0
	if err == nil {
		clr_buf(mdat, xdat, opts.DimX, opts.DimY, -1, -66)		// erase old data now

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
		if ud { fax(&uflg,&tflg,11) }
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
						l = " 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002 002"	// if scan fails, defaults to empty maze 32x32
						if y != 0 && y != opts.DimY { 
						l = " 002 000 000 000 000 000 000 000 000 000 000 000 000 000 000 000 000 000 000 000 000 000 000 000 000 000 000 000 000 000 000 000 002" }
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
// check for xtra buffer store
	xbf := prep(fil, "xb_")
	data, err = ioutil.ReadFile(xbf)
	if err == nil {
		dscan := fmt.Sprintf("%s",data)
		scanr := bufio.NewScanner(strings.NewReader(dscan))
		ix, iy := 0, 0
		l := "0 0"
		if scanr.Scan() { l = scanr.Text() }
		fmt.Sscanf(l,"%d %d",&ix, &iy)		// buffer size
fmt.Printf("xbuf %s -- %d x %d\n",xbf,ix,iy)
		if ix > 0 && iy > 0 {
			l, fin, wal := "gfx/floor016.jpg gfx/wall_jsgv_A.b.png", "", ""		// defaults on fail - this happens not...
			i := 0
			lsv := 500
			for fin != "xwfdn" && lsv > 0 {
				if scanr.Scan() { l = scanr.Text() }
				fin, wal = "xwfdn",""
				fmt.Sscanf(l,"%s %s",&fin, &wal)		// this loop will read cust walls & floor pairs until xwfdn
				if fin != "xwfdn" {
					if i <= maxwf { nwalflor() }
					wlfl.florn[i] = fin
					wlfl.walln[i] = wal
					err, _, wlfl.ftamp[i] = itemGetPNG(fin)
					if err != nil { wlfl.ftamp[i] = blankimage(64, 64) }
					err, _, wlfl.wtamp[i] = itemGetPNG(wal)
					if err != nil { wlfl.wtamp[i] = blankimage(832, 16) }
fmt.Printf("%d: %s %s\n",i,wlfl.florn[i],wlfl.walln[i])
					i++
				}
				lsv--
			}
			curwf = i	// current walls & floors for save
			for i := 0; i < maxwf; i++ { wlfl.flrblt[i] = false }		// clear all built floors
			for y := 0; y <= iy; y++ {	// read one line per element due to ops
			for x := 0; x <= ix; x++ {
				l, fin = "00", "000"
				if scanr.Scan() { l = scanr.Text() }
				fmt.Sscanf(l,"%s",&fin)
				xdat[xy{x, y}] = fin
				if ud { xubf[xy{x, y}] = fin }
fmt.Printf("%s ",fin)
			}
fmt.Printf("\n")
		}
		}
	} else {
		if opts.Verbose {
			fmt.Printf("edp %d failed < 0 or loading maze xtra buffer %s, warning:\n",edp,xbf)
			fmt.Print(err)
			fmt.Printf("\n")
		}
	}

// now load deleted elements - but not for pb or pal
  if ldb {
	dbf := prep(fil, ".db_") //fil[0:4]+".db_"+fil[4:len(fil)]
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
			fmt.Sscanf(l, "%d %d %d %d %s\n", &delbuf.elem[y],&delbuf.mx[y],&delbuf.my[y],&delbuf.revc[y],&delbuf.xbfd[y])
			if ud {
				udbck(6,y)
				udb.mx[y] = delbuf.mx[y]
				udb.my[y] = delbuf.my[y]
				udb.revc[y] = delbuf.revc[y]
				udb.elem[y] = delbuf.elem[y]
				udb.xbfd[y] = delbuf.xbfd[y]
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

func ed_sav(mazn int) {

	upd_edmaze(false)
	fil := fmt.Sprintf(".ed/g%dmaze%03d.ed",opts.Gtp,mazn)
	sav_maz(fil, xbuf, ebuf, eflg, opts.DimX, opts.DimY, mazn, true)
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
	if wpalop && palfol { palete(0) }
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
	lviewp := viewp
	lvpp := 0		// local viewport pass thru to upwin, if sd sim view over blows vp
// viewport ops
	if unpinx && vpx > opts.DimX + 1 { vpx = vpx - opts.DimX - 1 }
	if unpinx && vpx <= 0 - lviewp { vpx = vpx + opts.DimX + 1 }
	if unpiny && vpy > opts.DimY + 1 { vpy = vpy - opts.DimY - 1 }
	if unpiny && vpy <= 0 - lviewp { vpy = vpy + opts.DimY + 1 }
	fx := vpx + lviewp
	fy := vpy + lviewp
	if fx > opts.DimX && !unpinx { vpx = opts.DimX - lviewp + 1 }		// test scroll over endpoint, dont pass end of maze
	if fy > opts.DimY && !unpiny  { vpy = opts.DimY - lviewp + 1 }
	if vpx < 0 && !unpinx { vpx = 0 }
	if vpy < 0 && !unpiny { vpy = 0 }
	fx = vpx + lviewp
	fy = vpy + lviewp
fmt.Printf("viewport: %d sx,sy: %d, %d - ex,ey: %d, %d\n",lviewp,vpx,vpy,fx,fy)
	if opts.edat < 1 || opts.edat == 2 {	// simulate view mode for sd bufs if not in edit
		vpx, vpy = 0, 0						// & allow full view of edit regular mazes
		fx = opts.DimX
		if opts.DimY > fx { fx = opts.DimY }
		if fx < 30 { fx = 30 }
		fx++
		fy = fx
		lvpp = fx
fmt.Printf("lvpp: %d sx,sy: %d, %d - ex,ey: %d, %d\n",lvpp,vpx,vpy,fx,fy)
	}
	Ovimg := segimage(ebuf, xbuf, eflg, vpx, vpy, fx,fy, false)
	upwin(Ovimg, lvpp)
	calc_stats()
}

// replaceing or deleting - store for ctrl-z / ctrl-y

func undo_buf(sx int, sy int, rc int) {
//	fmt.Printf(" del %d elem: %d\n",delstak,delbuf.elem[delstak])
	delbuf.mx[delstak] = sx
	delbuf.my[delstak] = sy
	delbuf.revc[delstak] = rc					// revoke count for the loop
	delbuf.elem[delstak] = ebuf[xy{sx, sy}]
	delbuf.xbfd[delstak] = xbuf[xy{sx, sy}]
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

func rotmirmov(mdat MazeData, xdat Xdat, sx int, sy int, lastx int, lasty int, flg int) (int, int) {

// to transform maze, array copy
	xform := make(map[xy]int)
	xbxf := make(map[xy]string)
// transform																										 - rotating sq. wall mazes will always work
// rotate +90 degrees				-- * there is the issue of gauntlet arcade NEEDING the y = 0 wall *always* intact, rotating looper mazes wont work
		if opts.MRP {
			for ty := sy; ty <= lasty; ty++ {
			for tx := sx; tx <= lastx; tx++ {
				xform[xy{lasty - ty, tx}] = mdat[xy{tx, ty}]
				xbxf[xy{lasty - ty, tx}] = xdat[xy{tx, ty}]
// g1 - must transform all dors on a rotat since they have horiz & vert dependent
// xb doors will have to use g1 codes cause of this
				if xform[xy{lasty - ty, tx}] == G1OBJ_DOOR_HORIZ { xform[xy{lasty - ty, tx}] = G1OBJ_DOOR_VERT } else {
				if xform[xy{lasty - ty, tx}] == G1OBJ_DOOR_VERT { xform[xy{lasty - ty, tx}] = G1OBJ_DOOR_HORIZ } }
// g2
				if xform[xy{lasty - ty, tx}] == MAZEOBJ_DOOR_HORIZ { xform[xy{lasty - ty, tx}] = MAZEOBJ_DOOR_VERT } else {
				if xform[xy{lasty - ty, tx}] == MAZEOBJ_DOOR_VERT { xform[xy{lasty - ty, tx}] = MAZEOBJ_DOOR_HORIZ } }
			}}
			if lastx != lasty { lastx, lasty = is(lastx, lasty) }		// on a rotate in edit when size x != size y, they must swap after the rot
		} else {
		if opts.MRM {
			for ty := sy; ty <= lasty; ty++ {
			for tx := sx; tx <= lastx; tx++ {
				xform[xy{ty, lastx - tx}] = mdat[xy{tx, ty}]
				xbxf[xy{ty, lastx - tx}] = xdat[xy{tx, ty}]
// g1
				if xform[xy{ty, lastx - tx}] == G1OBJ_DOOR_HORIZ { xform[xy{ty, lastx - tx}] = G1OBJ_DOOR_VERT } else {
				if xform[xy{ty, lastx - tx}] == G1OBJ_DOOR_VERT { xform[xy{ty, lastx - tx}] = G1OBJ_DOOR_HORIZ } }
// g2
				if xform[xy{ty, lastx - tx}] == MAZEOBJ_DOOR_HORIZ { xform[xy{ty, lastx - tx}] = MAZEOBJ_DOOR_VERT } else {
				if xform[xy{ty, lastx - tx}] == MAZEOBJ_DOOR_VERT { xform[xy{ty, lastx - tx}] = MAZEOBJ_DOOR_HORIZ } }
			}}
			if lastx != lasty { lastx, lasty = is(lastx, lasty) }
		}
		}

// mirror x
		if opts.MH {
			for ty := sy; ty <= lasty; ty++ {
			for tx := sx; tx <= lastx; tx++ {
				xform[xy{lastx - tx, ty}] = mdat[xy{tx, ty}]
				xbxf[xy{lastx - tx, ty}] = xdat[xy{tx, ty}]
			}}
		}

// mirror y: flip
		if opts.MV {
			for ty := sy; ty <= lasty; ty++ {
			for tx := sx; tx <= lastx; tx++ {
				xform[xy{tx, lasty - ty}] = mdat[xy{tx, ty}]
				xbxf[xy{tx, lasty - ty}] = xdat[xy{tx, ty}]
			}}
			if flg&LFLAG4_WRAP_V > 0 {	// fix wall not allowed being at bottom for arcade gauntlet
				for ty := lasty - 1; ty >= sy ; ty-- {
				for tx := sx; tx <= lastx; tx++ {
					xform[xy{tx, ty + 1}] = xform[xy{tx, ty}]
					xbxf[xy{tx, ty + 1}] = xbxf[xy{tx, ty}]
				}}
				for tx := sx; tx <= lastx; tx++ { xform[xy{tx, 0}] = G1OBJ_WALL_REGULAR; xbxf[xy{tx, 0}] = "00" }
			}
		}

// copy back
		for y := sy; y <= lasty; y++ {
			for x := sx; x <= lastx; x++ {
				mdat[xy{x, y}] = xform[xy{x, y}]
				xdat[xy{x, y}] = xbxf[xy{x, y}]
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

func rotmirbuf(rmmaze *Maze, rmxdat Xdat, lastx,lasty int) (int, int) {

	fmt.Printf("in rotmirbuf\n")

// manual mirror, flip
	sx := 0
	sy := 0

	fmt.Printf("rmb wraps -- hw: %d vw: %d\n", rmmaze.flags&LFLAG4_WRAP_H,rmmaze.flags&LFLAG4_WRAP_V)
	fmt.Printf("rotmirbuf fx: %d lx %d fy %d ly %d\n", sx,lastx,sy,lasty)

	lastx, lasty = rotmirmov(rmmaze.data,rmxdat,sx,sy,lastx,lasty,rmmaze.flags)

	for y := sy; y <= lasty; y++ {
		for x := sx; x <= lastx; x++ {
			ebuf[xy{x, y}] = rmmaze.data[xy{x, y}]
		}
	}
	return lastx, lasty
}

// reload maze while editing & update window - generates output.png

func remaze(mazn int) {
fmt.Printf("\nin remaze dntr: %t edat:%d sdb: %d, delstk: %d, DIMS: %d - %d\n",opts.dntr,opts.edat,sdb,delstak,opts.DimX,opts.DimY)
	if !opts.dntr {
		sdb = -1
		delbset(0)
		clr_buf(ebuf, xbuf, 32, 32, -1, -66)
		edmaze = mazeDecompress(slapsticReadMaze(mazn), false)
		mazeloop(edmaze)
		opts.bufdrt = false
	}
	opts.dntr = false
	nsremaze = false
	if opts.edat > 0 || sdb > 0 { ed_maze(true) } else {
		Ovimg := genpfimage(edmaze, mazn)
		upwin(Ovimg, 0)
		calc_stats()
	}
}

// test location

func loc(mdat MazeData, x, y int) int {
	if x >= 0 && x <= opts.DimX && y >= 0 && y <= opts.DimY {
		return mdat[xy{x, y}]
	}
	return -1
}

// sx, sy allow re-entry to find element multiple times

func find(mdat MazeData, wh int, sx, sy int) (int, int) {

	rx, ry := -1, -1
	if sx < 0 || sx > opts.DimX { sx = 0 }
	if sy < 0 || sx > opts.DimY { sy = 0 }
	for y := sx; y <= opts.DimY; y++ {
		for x := sy; x <= opts.DimX; x++ {
			if loc(mdat, x, y) == wh { rx = x; ry = y; break }
	}}
	return rx, ry
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
var xplb Xdat
var plflg [14]int
var palxs int
var palys int
var palfol bool		// palette decor follows map chg

func palete(p int) {

	pmx := opts.DimX; pmy := opts.DimY
	clr_buf(plbuf, xplb, pmx, pmy, -1, -66)
	fil := fmt.Sprintf(".ed/pal%03d_g%d.ed",p,opts.Gtp)
	cnd := lod_maz(fil, xplb, plbuf, false, false)
	cpx = opts.DimX; cpy = opts.DimY

	if cnd >= 0 { fax(&plflg,&tflg,11)
		if palfol { fax(&plflg,&eflg,11)}
		bwin(cpx+1, cpy+1, 0, plbuf, xplb, plflg, "pal") }
	opts.DimX = pmx; opts.DimY = pmy
}

// typer for pal win

var statsB binding.Item[string]		// statistics win update
var listK binding.Item[string]		// editkey win update

func palRune(r rune) {

	switch r {
		case '?': dboxtx("palette ops", "in gved main window:\n"+
						"middle click an element\n(click activates edit + default key)\n... or ...\n"+
						"select maze\nhit <ESC> - activate edit mode\nhit '\\' for edit keys\n"+
						"hit a key to map: 'y'\nmove mouse to maze or palette\nand middle click an element\n"+
						"\nin palette window:\nmiddle click an element\n\n"+
						"edit hint on menu bar give status"+
						"\n\npal win keys:\nq,Q	┈ quit\nl,L	┈ list edit keys, auto updates\nt,T	┈ hide flags info\n"+
						"s, S	┈ stats window open, auto updates\n"+
						"f, F	┈ gauntlet 2 maze flags list\nx, X	┈ gauntlet 2 secret tricks", 350,540,nil,palRune)
		case 't': fallthrough
		case 'T': hide_flags(palRune)
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
						"bit position is within byte of that flag portion",420,744,nil,palRune)
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
						"NOHURTFRIENDS	= 0x11  \"Don't Hurt Friends\"",535,435,nil,palRune)
		case 'l': fallthrough
		case 'L': listK = dboxtx("Edit key assignments","",400,800, close_keys,palRune); list_keys()
		case 's': fallthrough
		case 'S': statsB = dboxtx("Maze stats","",340,700,close_stats,palRune); calc_stats()
		case 'q': fallthrough
		case 'Q': if wpalop { wpalop = false; wpal.Close() }
		default:
	}
}

// hide / select flags

func hide_flags(sr func(r rune)) {

dboxtx("T hide flags", "in gved main window:\n\n"+
		"invisible flag set - hide vars maze elements:\n"+
		" T	┈ cycle through a flag set (loop 0 - 511)\n"+
		" #T	┈ set flags = # ---- <ctrl>-T reset flags to 0\n\n"+
		" s 	┈ (show only) random 'special potions' &\n"+
		"	   gold bags on emply floor tiles (not EDIT)\n"+
		" L 	┈ toggle generator indicator letters [ DGLS ]\n"+
		"	   showing box gen monster type\n"+
		" p 	┈ toggle floor invisible *\n"+
		" P 	┈ toggle walls invisible *\n(use key in main window)\n\n"+
		"NOGEN	= 1		// all generators\n"+
		"NOMON	= 2		// all monster, dragon\n"+
		"NOFUD	= 4		// all food\n"+
		"NOTRS	= 8		// treas, locked\n"+
		"NOPOT	= 16		// pots & t.powers\n"+
		"NODOR	= 32		// doors, keys\n"+
		"NOTRAP	= 64		// trap & floor dots, stun, ff tiles\n"+
		"NOEXP	= 128		// exit, push wall\n"+
		"NOTHN	= 256		// anything else left\n"+
		"NOFLOOR	= 512\n"+
		"NOWALL	= 1024	// g2 *walls\n"+
		"NOG1W	= 2048	// g1 std wall only\n\n"+
		"set # with:\nBlank maze (file menu)\n- keep items flags cover\n\n"+
		"Random profile load\n- only load items flags cover"+
		"\n\n* T hide items disabled when edit keys active",440,700,nil,sr)
}

// when closing key lister panel, shut down updater
func close_keys() {
	listK = nil
}

// walls & floors g1/g2 color valid

var gmaxcol = 16

func colvld() int {

	return gmaxcol
}

// wall valid

var g1maxwal = 11		// only 6 walls, Se exp has shadowless walls
var g2maxwal = 12

func walvld() int {

	max := g1maxwal
	if G2 { max = g2maxwal }
	return max
}

// floors

var g1maxflor = 10		// floors 9,10 were not in game, a bit sketchy as well
var g2maxflor = 10

func florvld() int {

	max := g1maxflor
	if G2 { max = g2maxflor }
	return max
}

// supervalid - bounds check all 4 assignments

func suprval(wp,wc,fp,fc int) (int,int,int,int) {

	wip, wic, fip, fic := minint(walvld(),wp),minint(colvld(),wc),minint(florvld(),fp),minint(colvld(),fc)

	return wip, wic, fip, fic
}

// valid check, mapid
var g1maxid = 166
var g2maxid = 66

func maxid() int {

	max := g1maxid
	if G2 { max = g2maxid }
	return max
}

func valid_id(tid int) int {
	id := G1OBJ_WALL_REGULAR
	if tid >= 0 && tid <= maxid() { id = tid }
	return id
}

// assign keys

func key_asgn(buf MazeData, xdat Xdat, ax int, ay int) {

	edk := valid_keys(edkey)
	if G1 {
		g1edit_keymap[edk] = buf[xy{ax, ay}]
		g1edit_xbmap[edk] = xdat[xy{ax, ay}]
		if xbline != nil { xblchg = g1edit_xbmap[edk]; xbline.Set(xblchg) }
		kys := g1mapid[valid_id(g1edit_keymap[edk])]
		keyst := fmt.Sprintf("G¹ assn key: %s = %03d, %s",map_keymap[edk],g1edit_keymap[edk],kys)
		statlin(cmdhin,keyst)
		play_sfx(g1auds[g1edit_keymap[edk]])
		if edk == cycloc { cycl = g1edit_keymap[cycloc] }		// when reassign 'c' key, set cycl as well
	} else {
		g2edit_keymap[edk] = buf[xy{ax, ay}]
		kys := g2mapid[valid_id(g2edit_keymap[edk])]
		keyst := fmt.Sprintf("G² assn key: %s = %03d, %s",map_keymap[edk],g2edit_keymap[edk],kys)
		statlin(cmdhin,keyst)
		play_sfx(g2auds[g2edit_keymap[edk]])
		if edk == cycloc { cycl = g1edit_keymap[cycloc] }
	}
	if listK != nil { list_keys() }
}

// valid check, edit key

func valid_keys(ek int) int {

	if ek > maxkey || ek < minkey { return edkdef }		// 33 to 126, outside this return 121 'y'
	if G1 && g1edit_keymap[ek] < 0 { return edkdef }
	if G2 && g2edit_keymap[ek] < 0 { return edkdef }
	return ek
}

// list assigned keys when dialog open

func list_keys() {

	kl := "assigned edit keys:\n══════════════════════\n"
	for y := minkey; y <= maxkey; y++ {
		kv := 0
		if G1 { kv = g1edit_keymap[y] }
		if G2 { kv = g2edit_keymap[y] }
		sta := "┈┈┈┈┈┈┈┈  ᵏᵉʸ ⁿᵒᵗ ᵃˢˢᶦᵍⁿᵉᵈ"		// .................... ᴀssɪɢɴ ᴛʜɪs ᴋᴇʏ
														   //   ┈┈┈┈┈┈┈┈  ᵃˢˢᶦᵍⁿ ᵗʰᶦˢ ᵏᵉʸ
		if kv < 0 { sta = " ┈┈┈┈┈┈┈┈  ɴᴏ ᴀᴄᴄᴇss" }	//   ┈┈┈┈┈┈┈┈┈  n̵o̵t̵ ̵a̵v̵a̵i̵l̵a̵b̵l̵e
															// ɴᴏᴛ ᴀᴠᴀɪʟ		n̵o̵t̵ ̵a̵v̵a̵i̵l̵a̵b̵l̵e̵
		if kv > 0 && G1 { sta = g1mapid[valid_id(kv)] }
		if kv > 0 && G2 { sta = g2mapid[valid_id(kv)] }
		kl += fmt.Sprintf("%s\t=\t(%03d)   %s\n",map_keymap[y],kv,sta)
	}
	kl += "──────────────────────\nall keys (000) may be assigned\n"+
		  "with edit keys active, select the key\nand mouse middle click any maze item\n"+
		  "de-assign to (000) middle click any floor tile\n"+
		  "\nno keys (-01) may be assigned,\nthese are reserved system keys\n"+
		  "pressing a (-01) key will either:\n• select default key 'y'\n• operate the key function"
	if listK != nil { listK.Set(kl) }
}

// when closing stats panel, shut down updater
func close_stats() {
	statsB = nil
}

// stat package for palette win

func zero_stat() {

	for y := 0; y < 1000; y++ { g1stat[y] = 0; g2stat[y] = 0; stonce[y] = 1 }
}

// count stuff

func stats(elm int) {

	if elm >= 0 {
		if G1 { g1stat[elm]++ }
		if G2 { g2stat[elm]++ }
	}
}

// blotter select stats, or pb
var stl_str string

func mini_stat (buf MazeData, sx int, sy int, ex int, ey int, hed string) {

	stl := stl_str
	stl += "\n══════════════════════\n"+hed+"\n"
	zero_stat()
	for y := sy; y <= ey; y++ {
		for x := sx; x <= ex; x++ {
		stats(ebuf[xy{x, y}])
	}}
	tot := 0
	totnf := 0
	if G1 {
	for y := 0; y <= 65; y++ { if g1stat[y] > 0 {
		if opts.Verbose { fmt.Printf("  %s: %d\n",g1mapid[valid_id(y)],g1stat[y]) }
		stl += fmt.Sprintf("  %s: %d\n",g1mapid[valid_id(y)],g1stat[y])
		tot += g1stat[y]
		if y > 0 { totnf += g1stat[y] }
	}}}
	if G2 {
	for y := 0; y <= 65; y++ { if g2stat[y] > 0 {
		if opts.Verbose { fmt.Printf("  %s: %d\n",g2mapid[valid_id(y)],g2stat[y]) }
		stl += fmt.Sprintf("  %s: %d\n",g2mapid[valid_id(y)],g2stat[y])
		tot += g2stat[y]
		if y > 0 { totnf += g2stat[y] }
	}}}
	stl += "──────────────────────\n"
	stl += fmt.Sprintf("Total: %d\nTot, no floor: %d",tot, totnf)
	if statsB != nil { statsB.Set(stl) }

}

// and displate it

func calc_stats() {

	if wpalop { if palfol { palete(0) }}

	zero_stat()
//fmt.Printf("get stats: %d %d\n",opts.DimX,opts.DimY)
	for y := 0; y <= opts.DimY; y++ {
		for x := 0; x <= opts.DimX; x++ {
		stats(ebuf[xy{x, y}])
	}}
// stats during palette
	stl := fmt.Sprintf("%s\nmaze: %d x %d = %d cells\n",eid,opts.DimX+1,opts.DimY+1,(opts.DimX+1)*(opts.DimY+1))
		if opts.Verbose { fmt.Printf("%s\nstats:\n",stl) }

	if G1 {
	for y := 0; y <= 65; y++ { if g1stat[y] > 0 {
		if opts.Verbose { fmt.Printf("  %s: %d\n",g1mapid[valid_id(y)],g1stat[y]) }
		stl += fmt.Sprintf("  %s: %d\n",g1mapid[valid_id(y)],g1stat[y])
	}}}
	if G2 {
	for y := 0; y <= 65; y++ { if g2stat[y] > 0 {
		if opts.Verbose { fmt.Printf("  %s: %d\n",g2mapid[valid_id(y)],g2stat[y]) }
		stl += fmt.Sprintf("  %s: %d\n",g2mapid[valid_id(y)],g2stat[y])
	}}}
	stl += "══════════════════════\n"
	stl += mazeMetaPrint(edmaze, true)
	if statsB != nil { statsB.Set(stl) }
	stl_str = stl	// for mini stat append
//	bwin(palxs, palys, 0, plbuf, plflg)
}

// cut / copy & paste

var cpbuf MazeData	// c/c/p buffer
var xcpbuf Xdat		// exp c/c/p buffer
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
	for mx := 0; mx <= pmx; mx++ { cpbuf[xy{mx, my}] = 0; xcpbuf[xy{mx, my}] = "0" }}
	fil := fmt.Sprintf(".pb/%s_%07d_g%d.ed",id,vl,opts.Gtp)
	lod_maz(fil, xcpbuf, cpbuf, false, false)
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

	bwin(cpx+1, cpy+1, vl, cpbuf, xcpbuf, eflg, id)		// draw the buffer
	bl := fmt.Sprintf("paste buf: %d", vl)
	statlin(cmdhin,bl)
}

// transforms on cpbuf in smol window

func pb_loced(cnt int) {
	fil := fmt.Sprintf(".pb/pb_%07d_g%d.ed",cnt,opts.Gtp)
	sav_maz(fil, xcpbuf, cpbuf, eflg, cpx, cpy, 0, false)
	pb_upd("pb", "mas", cnt)
}

// typer for pb win

func pbRune(r rune) {

	switch r {
		case '?': dboxtx("paste buffer viewer keys", "q,Q	┈ quit\n"+
						"o,p	┈ cycle session pb - / +\n[,]	┈ cycle master pb - / +\n"+
						"l,L	┈ list edit keys\nt,T	┈ hide flags info\n"+
						"S  	┈ stats window open\n"+
						"══════════════════"+
						"\nr	┈ rotate pb +90°\nR	┈ rotate pb -90°\nh	┈ horiz flip\nm	┈ vert mirror"+
						"\na,d	┈ -/+ buffer horiz size\nw,s	┈ -/+ buffer vert size"+
						"\n══════════════════\nmiddle click selects element\n"+
						"left click sets element\n(pb has basic edit support)\n"+
						"- pb edit autosaves\n══════════════════\n(* only when window active)", 250,420,nil,pbRune)			// showing in main win because pb win is usually too small
		case '[': pbmas_cyc(-1)
		case ']': pbmas_cyc(1)
		case 'o': pbsess_cyc(-1)
		case 'p': pbsess_cyc(1)
// chunk size control
		case 'a': cpx--; if cpx < 1 { cpx = 1 }; pb_loced(masbcnt)
		case 'd': cpx++; pb_loced(masbcnt)
		case 'w': cpy--; if cpy < 1 { cpy = 1 }; pb_loced(masbcnt)
		case 's': cpy++; pb_loced(masbcnt)
// some extras
		case 'r': opts.MRP = true
			cpx,cpy = rotmirmov(cpbuf,xcpbuf,0,0,cpx,cpy,0)
			pb_loced(masbcnt)
		case 'R': opts.MRM = true
			cpx,cpy = rotmirmov(cpbuf,xcpbuf,0,0,cpx,cpy,0)
			pb_loced(masbcnt)
		case 'h': opts.MH = true
			rotmirmov(cpbuf,xcpbuf,0,0,cpx,cpy,0)
			pb_loced(masbcnt)
		case 'm': opts.MV = true
			rotmirmov(cpbuf,xcpbuf,0,0,cpx,cpy,0)
			pb_loced(masbcnt)
// std dialogs
		case 'l': fallthrough
		case 'L': listK = dboxtx("Edit key assignments","",400,800, close_keys,pbRune); list_keys()
		case 'S': statsB = dboxtx("Maze stats","",340,700,close_stats,pbRune); calc_stats()
		case 't': fallthrough
		case 'T': hide_flags(pbRune)
		case 'q': fallthrough
		case 'Q': if wpbop { wpbop = false; wpb.Close() }
		default:
	}
}

// display a  buffer window with buffer contents - no edit on palette
// px, py - size of paste buffer from 0, 0
// bn - buffer #

func bwin(px int, py int, bn int, mbuf MazeData, xdat Xdat, fdat [14]int, id string) {

var lw fyne.Window	// local cpy win to view buf contents
  wt := "palette selector"
  nimg := segimage(mbuf,xdat,fdat,0,0,px,py,false) // (bn == 0)) stats off for now
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
	dt = 16.0 * float32(opts.Geow - 4) / 528.0
  }
	clikwins(lw, nimg, px, py)
	lw.SetTitle(wt)
	lw.Resize(fyne.NewSize(float32(px)*dt, float32(py)*dt))
	specialKey(lw)

// if opts.Verbose {
fmt.Printf("clkwin sz: %v\n",fyne.NewSize(float32(px)*dt, float32(py)*dt))
}

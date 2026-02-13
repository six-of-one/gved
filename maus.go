package main

import (
	"fmt"
	"os"
	"image"
//	"image/color"
	"image/draw"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
    "fyne.io/fyne/v2/canvas"
)

// main mouse handler

// rubber banded

var blot *canvas.Image
var ccblot *canvas.Image
var blotimg string		// replace blotter with png image - blotter is stretched, so design must be right for outlines
var blotcol uint32		// with no image, this controls color & transparency in hex 0xAARRGGBB
var gvs bool			// use blotter to simulate view of gauntlet viewport

func blotter(img *image.NRGBA,px float32, py float32, sx float32, sy float32) {

	if img == nil {
		img = image.NewNRGBA(image.Rect(0, 0, 1, 1))
		draw.Draw(img, img.Bounds(), &image.Uniform{HRGB{blotcol}}, image.ZP, draw.Src)
	}
// config override for default blotter with png image
	if blotimg != "" {
			err, bll, _ := itemGetPNG(blotimg)
			if err == nil {
				blot = bll
			} else { blotimg = "" }
	}
	if blotimg == "" {
		blot = canvas.NewImageFromImage(img)
	}
	blot.Move(fyne.Position{px, py})
	blot.Resize(fyne.Size{sx, sy})
}

// turn off blotter after a window update
// because the window update...
// a. turns it on full maze for no reason
// b. refuses to turn it off, even with a delay in fn()
// and...
// c. resize window also covers the maze in blotter, which needs a fix
//		- blot.Hide() works, however blot.Show() flikers the entire maze with momentary blotter

func blotoff() {

	go func() {
			time.Sleep(5 * time.Millisecond)
   fyne.Do(func() {
			blot.Resize(fyne.Size{0, 0})
   })
	}()
}

// click area for edits

// button we can detect click and release areas for rubberband area & fills
// title tells us what window the button is on, assigned on btn creation if win is titled

type holdableButton struct {
    widget.Button
	title string
	bw fyne.Window
}

func newHoldableButton() *holdableButton {

    button := &holdableButton{}
    button.ExtendBaseWidget(button)

	return button
}

// no negative, because math.Min complained: cannot use float32() as float64 value in argument to math.Min
func nong(tv float32) float32 {

	v := tv
	if tv < 0.0 { v = 0.0 }
	return v
}

// store x & y when mouse button goes down - to start rubberband area
// 		and when released for other ops like cup & paste
var sxmd float64
var symd float64
var exmd float64
var eymd float64
// maze x & y mouse down
var mxmd int
var mymd int
var mbd bool			// true when mouse button 1 is held down, false otherwise
// mouse move pos global track
var rxm float32
var rym float32
// painter counter on undo, x,y
var prcl int
var pmx int
var pmy int

// &{{{387 545} {379 509.92188}} 4 0}

func (h *holdableButton) MouseMoved(mm *desktop.MouseEvent){
	ax := 0.0       // absolute x & y
	ay := 0.0
	rx := 0.0
	ry := 0.0
	mb := 0         // mb 1 = left, 2 = right, 4 = middle
	mk := 0         // mod key 1 = sh, 2 = ctrl, 4 = alt, 8 = logo
	pos := fmt.Sprintf("%v",mm)
	fmt.Sscanf(pos,"&{{{%f %f} {%f %f}} %d %d",&ax,&ay,&rx,&ry,&mb,&mk)
	cwt = h.title	// current window title by btn establish
//fmt.Printf("a: %f x %f rp: %f x %f\n",ax,ay,rx,ry)
	dt := float32(opts.dtec)
	sx := float32(sxmd)
	sy := float32(symd)
	ex := float32(rx)
	ey := float32(ry)
	if logo { mk = 8 } else { prcl = 1 }		// mod keys not picked up here ?
//	mbdi := 0; if mbd { mbdi = 1 }	// this is part of beef
//beef := fmt.Sprintf("a: %.0f x %.0f r: %.0f x %.0f dt: %.0f, mb/d %d/%d mk %d",sx,sy,ex,ey,dt,mb,mbdi,mk)
//statlin(cmdhin,beef)

	if strings.Contains(h.title, "G¹G²ved") {		// only in main win
		rxm = float32(rx)
		rym = float32(ry)
		lvpx, lvpy := 0, 0
		if opts.edat > 0 { lvpx, lvpy = vpx, vpy }
	if gvs {
		sx := nong(float32(int(ex / dt)) * dt - 7 * dt)
		sy := nong(float32(int(ey / dt)) * dt - 7 * dt)
		lx := 15 * dt
		ly := 15 * dt
		whlim := float32(opts.Geoh - 30)
		if sx + lx > whlim { sx = whlim - lx }
		if sy + ly > whlim { sy = whlim - ly }

		blot.Move(fyne.Position{sx, sy})
		blot.Resize(fyne.Size{lx, ly})
	} else {
	if ccp == PASTE {
//		ex = float32(float32(rx) + dt)
//		ey = float32(float32(ry) + dt)
		sx := nong(float32(int(ex / dt)) * dt - 3)
		sy := nong(float32(int(ey / dt)) * dt - 3)
		lx := float32(cpx) * dt + dt
		ly := float32(cpy) * dt + dt

		if blotup { blotwup(w, wpbimg) }
		blot.Move(fyne.Position{sx, sy})
		blot.Resize(fyne.Size{lx, ly})
	} else {
	tcmdhn := cmdhin
	tsshn := sshin
	if ex < sx { t := sx; sx = ex; ex = t }		// swap if end smaller than start
	if ey < sy { t := sy; sy = ey; ey = t }
// blotter size hinter, before pushing 1 past
		mxme := int(ex / dt)
		myme := int(ey / dt)
	ex = float32(float32(ex) + dt)					// click in 1 tile selects the tile
	ey = float32(float32(ey) + dt)
	if mbd {
// blotter size hinter
		mxmd = int(sx / dt) // redo as start / end can swap
		mymd = int(sy / dt)
// optimize blotter to cover selected cells
		sx = nong(float32(int(sx / dt)) * dt - 3)				// blotter selects tiles with original unit of 16 x 16
		sy = nong(float32(int(sy / dt)) * dt - 4)
		ex = float32(int(ex / dt)) * dt - 1
		ey = float32(int(ey / dt)) * dt - 2
		blot.Move(fyne.Position{sx, sy})
		blot.Resize(fyne.Size{ex - sx, ey - sy})
// blotter size hinter
		if mxmd == mxme && mymd == myme {
			mid := g1mapid[valid_id(ebuf[xy{mxmd+lvpx, mymd+lvpy}])]
			if G2 { mid = g2mapid[valid_id(ebuf[xy{mxmd+lvpx, mymd+lvpy}])] }
			pos = fmt.Sprintf("r: %.0f,%.0f+ %.0f cell: %d, %d elem: %d %s",sx,sy,dt,mxmd+lvpx,mymd+lvpy,ebuf[xy{mxmd+lvpx, mymd+lvpy}],mid)
		} else {
			dx := mxme-mxmd+1
			dy := myme-mymd+1
			pos = fmt.Sprintf("r: %.0f,%.0f - %.0f,%.0f mz: %d, %d to %d, %d... %d by %d = %d cells",sx,sy,ex,ey,mxmd+lvpx,mymd+lvpy,mxme+lvpx,myme+lvpy,dx,dy,dy*dx)
			hdr := fmt.Sprintf("r: %.0f,%.0f - %.0f,%.0f\nmz: %d, %d to %d, %d...\n%d by %d = %d cells\n",sx,sy,ex,ey,mxmd+lvpx,mymd+lvpy,mxme+lvpx,myme+lvpy,dx,dy,dy*dx)
			mini_stat(ebuf, mxmd+lvpx,mymd+lvpy,mxme+lvpx,myme+lvpy,hdr)
		}
		statlin(pos,tsshn)
//		fmt.Printf("st: %f x %f pos: %f x %f\n",sx,sy,ex,ey)
	} else {
		if mk == 8 {			// logo key = paint for any stored ops
			mxmd = int(rxm / dt)+lvpx // redo as start / end can swap
			mymd = int(rym / dt)+lvpy
			mid := g1mapid[valid_id(ebuf[xy{mxmd, mymd}])]
			if G2 { mid = g2mapid[valid_id(ebuf[xy{mxmd, mymd}])] }
			pos = fmt.Sprintf("paint: %.0f,%.0f+ %.0f cell: %d, %d elem: %d %s",rx,ry,dt,mxmd,mymd,ebuf[xy{mxmd, mymd}],mid)
			if pmx != mxmd || pmy != mymd {

				var setcode int			// code to store given edit hotkey
				var xstcode string		// code to store given edit hotkey
				if cmdoff {
				if G1 {
					setcode = g1edit_keymap[edkey]
					xstcode = g1edit_xbmap[edkey]
				} else {
					setcode = g2edit_keymap[edkey]
					xstcode = "00"
				}}
				if del { undo_buf(mxmd, mymd,prcl); ebuf[xy{mxmd, mymd}] = 0; xbuf[xy{mxmd, mymd}] = "0"; opts.bufdrt = true } else {	// delete anything for now makes a floor
				if setcode > 0 { undo_buf(mxmd, mymd,prcl); ebuf[xy{mxmd, mymd}] = setcode; xbuf[xy{mxmd, mymd}] = xstcode; opts.bufdrt = true }
				}
// FX: a loop is happening including this when mouse exits main win
fmt.Printf("prc: %d r: %.0f x %.0f cel: %d x %d - ls: %d x %d\n",prcl,rx,ry,mxmd,mymd,pmx,pmy)
				prcl++
				pmx = mxmd; pmy = mymd
				ed_maze(true)
			}
			flordirt = opts.bufdrt
			statlin(pos,tsshn)
		} else {				// no op on mouse move here
			statlin(tcmdhn,tsshn)
			blot.Resize(fyne.Size{0, 0})
	}}}}}
}

func (h *holdableButton) MouseDown(mm *desktop.MouseEvent){
	ax := 0.0	// absolute x & y
	ay := 0.0
	mb := 0		// mb 1 = left, 2 = right, 4 = middle
	mk := 0		// mod key 1 = sh, 2 = ctrl, 4 = alt, 8 = logo
	prcl = 1
	pos := fmt.Sprintf("%v",mm)
fmt.Printf("%v\n",mm)
	fmt.Sscanf(pos,"&{{{%f %f} {%f %f}} %d %d",&ax,&ay,&sxmd,&symd,&mb,&mk)
	mxmd = int(sxmd / opts.dtec)
	mymd = int(symd / opts.dtec)
if opts.Verbose {
if strings.Contains(h.title, "G¹G²ved") {
fmt.Printf("%d down - rel: %.0f x %.0f maze cell: %d x %d: %d\n",mb,sxmd,symd,mxmd,mymd,ebuf[xy{mxmd, mymd}])
} else {
fmt.Printf("%d down - rel: %.0f x %.0f maze cell: %d x %d\n",mb,sxmd,symd,mxmd,mymd)
}}
	mbd = (mb == 1)
	if mbd { h.MouseMoved(mm) }		// engage 1 tile click
}

var repl int		// replace will be by ctrl-h in select area or entire maze, by match
var cycl int		// cyclical set - C cycles, c sets - using c loc in keymap
var cycloc = 99

// edkey 'locks' on when pressed

func (h *holdableButton) MouseUp(mm *desktop.MouseEvent){

	mb := 0		// mb 1 = left, 2 = right, 4 = middle
	mbd = false
	h.MouseMoved(mm)				// disengage blotter
	ax := 0.0	// absolute x & y
	ay := 0.0
	exmd = 0.0	// rel x & y interm float32
	eymd = 0.0
	mk := 0		// mod key 1 = sh, 2 = ctrl, 4 = alt, 8 = logo
	pos := fmt.Sprintf("%v",mm)
	fmt.Sscanf(pos,"&{{{%f %f} {%f %f}} %d %d",&ax,&ay,&exmd,&eymd,&mb,&mk)
	dt := opts.dtec
	edkey = valid_keys(edkey)

// pal win
	inpal := false
	if wpalop { if h.bw == wpal {
		inpal = true
		dt = 16.0 * float64(opts.Geow - 4) / 528.0		// palette dtec is locked at orig win size
	}}
	ex := int(exmd / dt) + vpx
	ey := int(eymd / dt) + vpy
// middle mouse click anywhere activates edit mode & pulls up def key
	if mb == 4 {
		if opts.edat == 0 || !cmdoff { edit_on(edkdef) }
		if wpalop {					// palette element selector
		if inpal {
				if cmdoff { key_asgn(plbuf, xplb, int(exmd / dt), int(eymd / dt)); sv_config() }
				return
			}
	}}
// right mb functions
	if mb == 2 {
		if strings.Contains(h.title, " pbf") { pbmas_cyc(1) } else {
		if pgdir != 0 {
			if sdb > 0 {
				sdbit(pgdir)
			} else {
				lrelod := pagit(pgdir)
				upd_edmaze(false)
				if lrelod { remaze(opts.mnum) }
			}
		}}
		return
	}

 //   fmt.Printf("up %v\n",mm)
	if opts.edat > 0 {
		opbuf := ebuf
		xopbf := xbuf
		pbe := false		// paste buf edit
		if strings.Contains(h.title, " pbf") {			// simple edit on pb win content
			opbuf = cpbuf
			xopbf = xcpbuf
			ccp = NOP
			pbe = true
		}
//fmt.Printf("%d up: %.0f x %.0f \n",mb,exmd,eymd)

		sx := int(sxmd / dt) + vpx
		sy := int(symd / dt) + vpy
		if ex < sx { t := ex; ex = sx; sx = t }		// swap if end smaller than start
		if ey < sy { t := ey; ey = sy; sy = t }
		var setcode int			// code to store given edit hotkey
		var xstcode string
		if G1 {
			setcode = g1edit_keymap[edkey]
			xstcode = g1edit_xbmap[edkey]
		} else {
			setcode = g2edit_keymap[edkey]
			xstcode = "00"
		}
// a cut / copy / paste is active
		pasty := false
		if ccp != NOP && !inpal {
		if ccp == PASTE { pasty = true }
		if mb != 1 { ccp_NOP(); fmt.Printf("mb: ccp to NOP\n") }
		if ccp != NOP {
			px :=0
			if ccp == COPY || ccp == CUT {
				py :=0
			for my := sy; my <= ey; my++ {
				px =0
			for mx := sx; mx <= ex; mx++ {
				cpbuf[xy{px, py}] = opbuf[xy{mx, my}]
				xcpbuf[xy{px, py}] = xopbf[xy{mx, my}]
fmt.Printf("%03d ",cpbuf[xy{px, py}])
				px++
				}
fmt.Printf("\n")
			py++
			}
			cpx = px - 1; if cpx < 0 { cpx = 0 }		// if these arent 1 less, the paste is 1 over
			cpy = py - 1; if cpy < 0 { cpy = 0 }
// if opts.Verbose {
fmt.Printf("cc dun: px %d py %d\n",px,py)
// saving paste buffer now
			fil := fmt.Sprintf(".pb/pb_%07d_g%d.ed",pbcnt,opts.Gtp)
			pbcnt++
			sav_maz(fil, xcpbuf, cpbuf, eflg, cpx, cpy, 0, false)
// local for short range
			fil = fmt.Sprintf(".pb/ses_%07d_g%d.ed",lpbcnt,opts.Gtp)
			lpbcnt++
			if G1 { lg1cnt = lpbcnt}
			if G2 { lg2cnt = lpbcnt}
			sav_maz(fil, xcpbuf, cpbuf, eflg, cpx, cpy, 0, false)

// pb sort of doesnt end
				fil = fmt.Sprintf(".pb/cnt_g%d",opts.Gtp)
				file, err := os.Create(fil)
				if err == nil {
					wfs := fmt.Sprintf("%d",pbcnt)
					file.WriteString(wfs)
					file.Close()
				}
			}
			del = false						// copy or paste should not have del on
			if ccp == CUT { del = true }
			if pasty {
// if opts.Verbose {
fmt.Printf("in pasty\n")
				setcode = 0; 
				ex = sx + cpx
				ey = sy + cpy
				if ex < 0 || ey < 0 { fmt.Printf("paste fail x\n"); return }
				if ex > opts.DimX { ex = opts.DimX }
				if ey > opts.DimY { ey = opts.DimY }
			} else {
				bwin(cpx+1, cpy+1, pbcnt - 1, cpbuf, xcpbuf, eflg, "ses")		// draw the buffer
				if ccp == COPY { return }
			}
		}}
// no access for keys: ?, \, C, A #a, eE, L, S, H, V
// if opts.Verbose {
mid := g1mapid[valid_id(opbuf[xy{ex, ey}])]
if G2 { mid = g2mapid[valid_id(opbuf[xy{ex, ey}])] }
fmt.Printf(" dtec: %.2f maze: %d x %d - element:%d - %s --- XB: %s\n",dt,ex,ey,opbuf[xy{ex, ey}],mid,xopbf[xy{ex, ey}])
		if mb == 4 && cmdoff {		// middle mb, do a reassign
			 key_asgn(opbuf, xopbf, ex, ey); sv_config()
		} else {
		if inpal { return }		// L clicks on palete should not place on main
		if del || cmdoff || pasty {
			rcl := 1		// loop count for undoing multi ops
		 for my := sy; my <= ey; my++ {
		   for mx := sx; mx <= ex; mx++ {
			rop := true		// run ops
			if ctrl {		// with ctrl held on drag op, only do outline
				rop = false
				if my == sy || my == ey { rop = true }
				if mx == sx || mx == ex { rop = true }
			}
// looped now, with ctrl op
			if rop {
				delstr := 0
				if shift { delstr = -1 }
				if del { undo_buf(mx, my,rcl); opbuf[xy{mx, my}] = delstr; xopbf[xy{mx, my}] = "0"; opts.bufdrt = true } else {	// delete anything for now makes a floor
				if pasty { undo_buf(mx, my,rcl); opbuf[xy{mx, my}] = cpbuf[xy{mx - sx, my - sy}]; xopbf[xy{mx, my}] = xcpbuf[xy{mx - sx, my - sy}]; opts.bufdrt = true }	// cant use setcode below, it wont set floors
				if setcode > 0 {						// i think i found the bug where paste doesnt work right
					undo_buf(mx, my,rcl);
					if ! shift { opbuf[xy{mx, my}] = setcode; fmt.Printf("!shift set\n")}
					if ! ctrl { xopbf[xy{mx, my}] = xstcode; }
					opts.bufdrt = true
				}
fmt.Printf("---stored %03d ::%s",opbuf[xy{mx, my}],xopbf[xy{mx, my}])
				}
				rcl++
			}
			if edkey == 314 { repl = opbuf[xy{mx, my}] }		// just placeholder until new repl done -- yes, NOT being used
//			if edkey == 182 { opbuf[xy{mx, my}] = repl; opts.bufdrt = true }		//
		  }
fmt.Printf("\n")
		}
//			fmt.Printf(" chg elem: %d maze: %d x %d\n",opbuf[xy{mx, my}],mx,my)
		}}
		flordirt = opts.bufdrt
		if pbe && opts.bufdrt {			// paste buf edit
			pb_loced(masbcnt)
			opts.bufdrt = false
		} else { ed_maze(true) }
	}

}

package main

import (
	"image"
//	"image/png"
	"math"
	"math/rand"
	"fmt"
//	"image/draw"
	"github.com/fogleman/gg"
	"image/color"
	"encoding/binary"
	"golang.org/x/image/draw"
)


// arrays for item masks
var g1mask [256]int
var g2mask [256]int

// for maze output to se -- outputter is in pfrender
func ParseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff
	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("invalid length, must be 7 or 4")

	}
	return
}

var foods = []string{"ifood1", "ifood2", "ifood3"}
var nothing int

// scan maze data - handle unpins & wraps

func whatis(maze *Maze, x, y int) int {
	return maze.data[xy{x, y}]
}

// isolating loop over target code
// sx,y - starting point
// tx,y - test point

func lot(sx, sy, tx, ty int) (bool, int, int) {

	rf := true				// return over flows
	if tx < 0 {
		if !unpinx && tx != -1  { rf = false }
		if opts.edat == 0 && tx != -1  { rf = false }		// not entirely sure - border wall should always render correct
		tx = opts.DimX + tx + 1
	}

	if tx > opts.DimX {
		if !unpinx && tx != opts.DimX + 1 { rf = false }
		if opts.edat == 0 && tx != opts.DimX + 1 { rf = false }
		tx = tx - opts.DimX - 1
	}

	if ty < 0 {
		if !unpiny && ty != -1 { rf = false }
		if opts.edat == 0 && ty != -1  { rf = false }
		ty = opts.DimY + ty + 1
	}

	if ty > opts.DimY {
		if !unpiny && ty != opts.DimY + 1 { rf = false }
		if opts.edat == 0 && ty != opts.DimY + 1 { rf = false }
		ty = ty - opts.DimY - 1
	}

	if tx < 0 { tx = 0 }
	if ty < 0 { ty = 0 }

	return rf, tx, ty
}

// scan buffer data same,
// sx,y - starting point
// tx,y - test point
// asgn - if > -2, this is assign value
//		  when testing shadows, etc, tells where we started
//		  so slip calc (past maze edge to other side) math works

func scanbuf (mdat MazeData, sx, sy, tx, ty, asgn int) int {

	i := -1
	txe, tye := tx, ty		// for debug fmt so we know how test is adjusted

	rf, ux, uy := lot(sx, sy, tx, ty)
	if rf { i = mdat[xy{ux, uy}] }

if false && vpx < 0 {
fmt.Printf("scan: %d s-e: %d x %d, %d x %d test: %d x %d\n",i,sx,sy,txe,tye,tx,ty)
}
// the assigner, for when we need it
//		if asgn > -2 { mdat[xy{tx, ty}] = asgn }
	return i
}
/* scnbuf test out:
scnbuf s-e: 0 x 31, 1 x 30 dif: 1, 1 test: 1 x 30
scnbuf s-e: 0 x 31, -1 x 31 dif: 1, 0 test: 31 x 31
scnbuf s-e: 0 x 31, 1 x 31 dif: 1, 0 test: 1 x 31
scnbuf s-e: 0 x 31, -1 x 32 dif: 1, 1 test: 31 x 0
scnbuf s-e: 0 x 31, 0 x 32 dif: 0, 1 test: 0 x 0

scnbuf s-e: 29 x 31, 29 x 32 dif: 0, 1 test: 29 x 0
scnbuf s-e: 29 x 31, 28 x 31 dif: 1, 0 test: 28 x 31
*/

// also need to scan xbuf data same coords system

func scanxb (xdat Xdat, sx, sy, tx, ty int, asgn string) string {

	v := "0"
	txe, tye := tx, ty		// for debug fmt so we know how test is adjusted

	rf, ux, uy := lot(sx, sy, tx, ty)
	if rf { v = xdat[xy{ux, uy}] }

if false && vpx < 0 {
fmt.Printf("scan: %s s-e: %d x %d, %d x %d test: %d x %d\n",v,sx,sy,txe,tye,tx,ty)
}
// the assigner, for when we need it
//		if asgn != "" { xdat[xy{tx, ty}] = asgn }
	return v
}

// door check

func isdoor(t int) bool {
	if t == MAZEOBJ_DOOR_HORIZ || t == MAZEOBJ_DOOR_VERT {
		return true
	} else {
		return false
	}
}

// G¹ version
func isdoorg1(t int) bool {
	if t == G1OBJ_DOOR_HORIZ || t == G1OBJ_DOOR_VERT {
		return true
	} else {
		return false
	}
}

// check to see if there's walls adjacent left, left/down, and down
// used to set wall shadows, G¹ engine darkens floor pixels with palette shift
// horizontal wall += 4
// diagonal wall += 8
// vertical wall += 16

// shadows by wall pattern - G²

func shad_wallpat() int {
	wp := 6
	if G2 { wp = 11 }
	return wp
}

// G¹ version
func checkwalladj3g1(maze *Maze, xdat Xdat, x int, y int) int {
	adj := 0
	wp := maze.wallpattern
	wpsha := shad_wallpat()		// what wall patterns have shadows, G2 < 11, SE < 6

	if !iswallg1(scanbuf(maze.data, x, y, x, y, -2)) {	// no need for a shadow under a wall

		xp := scanxb(xdat, x, y, x-1, y, "")
		p,_,_ := parser(xp, SE_WALL)			// set wall pat
		if p >= 0 { wp,_,_,_ = suprval(p,0,0,0) }
		if iswallg1(scanbuf(maze.data, x, y, x-1, y, -2)) && wp < wpsha {
			adj += 4
		}

		xp = scanxb(xdat, x, y, x, y+1, "")
		p,_,_ = parser(xp, SE_WALL)			// set wall pat
		if p >= 0 { wp,_,_,_ = suprval(p,0,0,0) }
		if iswallg1(scanbuf(maze.data, x, y, x, y+1, -2)) && wp < wpsha  {
			adj += 16
		}

		xp = scanxb(xdat, x, y, x-1, y+1, "")
		p,_,_ = parser(xp, SE_WALL)			// set wall pat
		if p >= 0 { wp,_,_,_ = suprval(p,0,0,0) }
		if iswallg1(scanbuf(maze.data, x, y, x-1, y+1, -2)) && wp < wpsha  {
			adj += 8
		}
	}

	return adj
}

// check to see if there's walls on any side of location, for picking
// which wall tile needs ot be used
//
// Values to use:
//    up left:  0x01      up:         0x02      up right:  0x04
//    left:     0x08      right:      0x10      down left: 0x20
//    down:     0x40      down right: 0x80

// added in Se wally value for picking walls from Se wall set

// G¹ version -- G² has more walls
func checkwalladj8g1(maze *Maze, x int, y int) (int, int) {
	adj := 0
	wally := 0					// sanctuary expanded wall tester mk II
	wp1, wp2 := false, false	// pointers for findwall(mpixel(tx, ty, tx+1, ty+1, m)),
								//				findwall(mpixel(tx, ty, tx-1, ty+1, m))

	if iswallg1(scanbuf(maze.data, x, y, x-1, y-1, -2)) {
		adj += 0x01
	}
	if iswallg1(scanbuf(maze.data, x, y, x, y-1, -2)) {
		adj += 0x02
		wally |= 1
	}
	if iswallg1(scanbuf(maze.data, x, y, x+1, y-1, -2)) {
		adj += 0x04
	}
	if iswallg1(scanbuf(maze.data, x, y, x-1, y, -2)) {
		adj += 0x08
		wally |= 8
	}
	if iswallg1(scanbuf(maze.data, x, y, x+1, y, -2)) {
		adj += 0x010
		wally |= 2
	}
	if iswallg1(scanbuf(maze.data, x, y, x-1, y+1, -2)) {
		adj += 0x20
		wp2 = true
	}
	if iswallg1(scanbuf(maze.data, x, y, x, y+1, -2)) {
		adj += 0x40
		wally |= 4
	}
	if iswallg1(scanbuf(maze.data, x, y, x+1, y+1, -2)) {
		adj += 0x80
		wp1 = true
	}
// Se logics extending wall set
	if wally > 13 {
		if wp2 {		// adj = 0x20
			wally += 6
			if wp1 {	// adj = 0x80
				wally += 4
			}
		} else if wp1 {	// adj = 0x80
			wally += 8
		}
	}

	if wally == 6 || wally == 7 {
		if wp1 {		// adj = 0x80
			wally += 10
		}
	}

	if wally == 12 || wally == 13 {
		if wp2 {		// adj = 0x20
			wally += 6
		}
	}

	return adj, wally
}

// Look and see what doors are adjacent to this door
//
// Values to use:
//    up:  0x01    right:  0x02    down:  0x04    left:  0x08

// G¹ version
func checkdooradj4g1(maze *Maze, x int, y int) int {
	adj := 0

	if isdoorg1(scanbuf(maze.data, x, y, x, y-1, -2)) {
		adj += 0x01
	}
	if isdoorg1(scanbuf(maze.data, x, y, x+1, y, -2)) {
		adj += 0x02
	}
	if isdoorg1(scanbuf(maze.data, x, y, x, y+1, -2)) {
		adj += 0x04
	}
	if isdoorg1(scanbuf(maze.data, x, y, x-1, y, -2)) {
		adj += 0x08
	}

	return adj
}

// Below lies the stuff for figuring out where forcefield ground tiles
// should go. It's not particularly efficient or elegant, but it works.
var ffLoopDirs = []xy{
	xy{0, -1}, // "up"
	xy{1, 0},  // right
	xy{0, 1},  // "down"
	xy{-1, 0}, // left
}

var adjvalues = []int{0x01, 0x02, 0x04, 0x08}

func checkffadj4(maze *Maze, x int, y int) int {
	adj := 0
	for i := 0; i < 4; i++ {
		for j := 1; j <= 15; j++ {
			t := scanbuf(maze.data, x, y, x+(j*ffLoopDirs[i].x), y+(j*ffLoopDirs[i].y), -2)
			if j > 1 && isforcefield(t) {
				adj += adjvalues[i]
				break
			} else if iswall(t) {
				break
			}
		}
	}

	return adj
}

type FFMap map[xy]bool
type AtMap map[xy]int		// animate tiles
var anmap AtMap
var anmapt AtMap			// animated timer cnt
var svanim bool				// if true save animated map

func ffMark(ffmap FFMap, maze *Maze, x int, y int, dir int) {
	for i := 1; i < 90000; i++ {		// this had no upper limit and could inf loop if ff were skunky
		d := ffLoopDirs[dir]			// -- 90k reps a maze 300 x 300, this may already cause a delay on a bad ff placement
		nx := x + (d.x * i)
		ny := y + (d.y * i)

		if isforcefield(scanbuf(maze.data, nx, ny, nx, ny, -2)) {		// maze.data[xy{nx, ny}]) {
			// done with this direction
			return
		}

		// mark our map
		ffmap[xy{nx, ny}] = true
	}

}

func ffMakeMap(maze *Maze) FFMap {
	ffmap := FFMap{}

	for k, v := range maze.data {
		if !isforcefield(v) {
			if svanim {
				anmap[xy{k.x, k.y}] = isanimtil(v)
				anmapt[xy{k.x, k.y}] = 0
if anmap[xy{k.x, k.y}] > 0 {fmt.Printf("det anim %d: %d x %d\n",v,k.x, k.y)}
			}
			continue
		}

		// Only check for 'right' or 'down' adjacencies, since up and left
		// are just the same tiles from the other end
		adj := checkffadj4(maze, k.x, k.y)
		if (adj & 0x02) > 0 { // adj to the right
			ffMark(ffmap, maze, k.x, k.y, 1)
		}
		if (adj & 0x04) > 0 { // adj down
			ffMark(ffmap, maze, k.x, k.y, 2)
		}
	}

	return ffmap
}

func isforcefield(t int) bool {
	if t == MAZEOBJ_FORCEFIELDHUB || t == SEOBJ_FORCEFIELDHUB {
		return true
	} else {
		return false
	}
}

// returns item if animated, all arstamp[] have a frames count in them

func isanimtil(t int) int {
	r := 0
	for i := 0; animcyc[i] > 0; i++ {
		if animcyc[i] == t { r = t; manim = true }
	}
	return r
}

func dotat(img *image.NRGBA, xloc int, yloc int) {
	c := IRGB{0xffff}

	for y := 0; y < 2; y++ {
		for x := 0; x < 2; x++ {
			img.Set(xloc+x, yloc+y, c)
		}
	}
}

func renderdots(img *image.NRGBA, xloc int, yloc int, count int) {
	switch count {
	case 1:
		dotat(img, xloc+7, yloc+7)
	case 2:
		dotat(img, xloc+9, yloc+5)
		dotat(img, xloc+5, yloc+9)
	case 3:
		dotat(img, xloc+7, yloc+7)
		dotat(img, xloc+9, yloc+5)
		dotat(img, xloc+5, yloc+9)
	case 4:
		dotat(img, xloc+9, yloc+5)
		dotat(img, xloc+5, yloc+9)
		dotat(img, xloc+5, yloc+5)
		dotat(img, xloc+9, yloc+9)
	}
}

// write a 16x16 tile of any color onto img @x,y, can be fed hex tripl 0xrrggbb or 0xaarrggbb

func coltil(img *image.NRGBA, col uint32, xloc int, yloc int) {
	c := HRGB{col}
//	b := HRGB{0xffffff-col}

	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
//			if y & 3 == 0 && x & 3 == 0 { img.Set(xloc+x, yloc+y, b) } else { // this is a dot field
				img.Set(xloc+x, yloc+y, c)
//			}
		}
	}
}

// viewport going neg for loops needs coord adjust to write stamps on canvas
// coord, coord begin, bias adj

func vcoord(c, cb, ba int) int {

	i := c-cb+ba
	if cb > 0 { return i }		// main ajust > 0, do std
	i = c+ba
	return i
}
//writestamptoimage(G1,img, stamp, (x-xb+xba)*16, (y-yb+yba)*16)

// write sqaure off png image grid to img
// img - image to write on, if nil no write
// xw, yw - x & y loc to write on img
// ptamp - png image stamp
// rx, ry - pixel size of rectaNGLE to copy
// ax, ay - add this to rectangle
// clm rw - col (x) & row (y) tile of ptamp		-- rx, ry used for row / col select unless st set as sub tile size, rx,ry always used for sizing
// also return extracted area
var dbgwrt bool
//								dbgwrt = true
var parimg image.Image			// keep large images off stacks

func writepngtoimage(img *image.NRGBA, rx,ry,ax,ay,cl,rw, xw, yw, st int) {

	bnds := parimg.Bounds()
	iw, ih := bnds.Dx(), bnds.Dy()
	tsx, tsy := rx,ry
	if st > 0 { tsx, tsy = st,st }
	rec := image.Rect((cl)*tsx, (rw)*tsy, ax+(cl+1)*rx, ay+(rw+1)*ry)			// this pegs the intended rect on sprite sheet
	rrr := image.Rect(0,0,iw,ih)
	cpy := image.NewRGBA(rrr)
	draw.Copy(cpy, image.Pt(0,0), parimg, rec, draw.Over, nil)
if dbgwrt {
fmt.Printf("sz %d %d c, r %d, %d vp abs %d x %d\n",rx,ry,cl,rw,xw,yw)
}
	offset := image.Pt(xw, yw)
	if img != nil {
		draw.Draw(img, cpy.Bounds().Add(offset), cpy, image.ZP, draw.Over)
	}
//	return cpy
}

// teh quickening - stop sending large images thru parms, all floor & wall slice store img comes here

func writewftoimage(img *image.NRGBA, ftmp,fflm,wtmp, rx,ry,ax,ay,cl,rw, xw, yw int) {

//fmt.Printf("wwf %d %d %d - x,y %d %d\n",ftmp,fflm,wtmp, xw, yw)
	var lptamp image.Image
	if ftmp >= 0 { lptamp = wlfl.ftamp[ftmp] }
	if fflm >= 0 { lptamp = wlfl.flim[fflm] }
	if wtmp >= 0 { lptamp = wlfl.wtamp[wtmp] }

	if lptamp != nil {
		bnds := lptamp.Bounds()
		iw, ih := bnds.Dx(), bnds.Dy()
		rec := image.Rect((cl)*rx, (rw)*ry, ax+(cl+1)*rx, ay+(rw+1)*ry)			// this pegs the intended rect on sprite sheet
		rrr := image.Rect(0,0,iw,ih)
		cpy := image.NewRGBA(rrr)
		draw.Copy(cpy, image.Pt(0,0), lptamp, rec, draw.Over, nil)
if dbgwrt {
fmt.Printf("sz %d %d c, r %d, %d vp abs %d x %d\n",rx,ry,cl,rw,xw,yw)
}
		offset := image.Pt(xw, yw)
		if img != nil {
//if ftmp >= 0 { fmt.Printf("wrote ftmp %d, %d\n",xw, yw)}
			draw.Draw(img, cpy.Bounds().Add(offset), cpy, image.ZP, draw.Over)
		}
	}
//	return cpy
}

/*
				writepngtoimage(img, shtamp, 16,16,na,0, (x-xb)*16, (y-yb)*16)

				r := image.Rect((na)*16, 0, (na+1)*16, 16)
				rr := image.Rect(0,0,256,16)
				shd := image.NewRGBA(rr)
				draw.Copy(shd, image.Pt(0,0), shtamp, r, draw.Over, nil)
//fmt.Printf("shadow %d: %d x %d \n",na,(x-xb)*16, (y-yb)*16)
				offset := image.Pt((x-xb)*16, (y-yb)*16)
				draw.Draw(img, shd.Bounds().Add(offset), shd, image.ZP, draw.Over)	*/

type walflr struct {
	ftamp	[]image.Image
	flim	[]*image.NRGBA
	wtamp	[]image.Image
	florn   []string
	walln   []string
	flrtls	[]bool
	totw	[]int				// total w & h of a floor built in flim, may need expanded if a maze gets larger
	toth	[]int
}
var wref []int					// ref pntrs to walflr array, only going to load / build floors once
var fref []int

var maxwf int
var curwf int
var wlfl = &walflr{}

// master floor/wall replace
var Se_mflor int
var Se_mwal int
var Se_rwal int
var Se_rrnd int

// when map is loaded, store floors & walls as designated in xb_*.ed file after X Y size and before "xwfdn" marker

func nwalflor(){

//fmt.Printf("delbuf st: %d len %d, test: %d\n",delstak,len(delbuf.elem),t)
	maxwf++
	wlfl.florn = append(wlfl.florn,"")
	wlfl.walln = append(wlfl.walln,"")
	wlfl.ftamp = append(wlfl.ftamp,nil)				// floor tile loaded from fil
	wlfl.flim  = append(wlfl.flim,nil)				// floor panel made for maze writepngtoimage
	wlfl.wtamp = append(wlfl.wtamp,nil)				// wall tiles, should be 26 16x16 segments + 26 shot wall 16x16 segs, as many rows as desired
	wlfl.flrtls = append(wlfl.flrtls,false)			// flag indicates set of floor tiles not rendering into flim
	wlfl.totw = append(wlfl.totw,0)
	wlfl.toth = append(wlfl.toth,0)
	wref = append(wref,0)							// ref pointer in wlfl slice, floors & walls, only load once
	fref = append(fref,0)
}

// find a wall or floor in the slice struct

func findwf(fl,wl string) (int, int) {
	f,w := -1,-1
	for i := 0; i <= maxwf; i++ {
		if f < 0 && fl == wlfl.florn[i] { f = i }
		if w < 0 && wl == wlfl.walln[i] { w = i }
//fmt.Printf("srch %d: %s %s\n",i,wlfl.florn[i],wlfl.walln[i])
	}
fmt.Printf("fnd f%d, w%d = ",f,w)
if f >= 0 {fmt.Printf(" %s",wlfl.florn[f])}
if w >= 0 {fmt.Printf(" %s",wlfl.walln[w])}
fmt.Printf("\n")
	return f,w
}

// build each loaded flim

func florflim(p int) {

	if wlfl.flrtls[p] { return }	// dont render tile set into flim here
	bnds := wlfl.ftamp[p].Bounds()
	iw, ih := bnds.Dx(), bnds.Dy()		// in theory this image does not HAVE to be square anymore
	totw :=  int(math.Ceil(float64(((opts.DimX+1)*16))/float64(iw))) * iw		// round up so images not divinding easily into maze size cover entire maze
	toth :=  int(math.Ceil(float64(((opts.DimY+1)*16))/float64(ih))) * ih
	if totw <= wlfl.totw[p] && toth <= wlfl.toth[p] { return }

	if wlfl.totw[p] == 0 { wlfl.flim[p] = blankimage(totw, toth) }
	for ty := 0; ty < toth ; ty=ty+ih {
	for tx := 0; tx < totw ; tx=tx+iw {
		offset := image.Pt(tx, ty)
		draw.Draw(wlfl.flim[p], wlfl.ftamp[p].Bounds().Add(offset), wlfl.ftamp[p], image.ZP, draw.Over)
	}}
	 wlfl.totw[p], wlfl.toth[p] = totw, toth
fmt.Printf("flim %s entry %d t:%d x %d, src %d x %d\n",wlfl.florn[p],p,totw, toth,iw,ih)
}

// make base floor, of: null space, SE_COLRT, SE_CFLOR, SE_TFLOR, SE_NOFLOR, Se_mflor, std floor, adj/wly shadows, ff beams

var florb *image.NRGBA
var flordirt int			// whether or not an edit could dirty the flor, pb & palete set to -1
var fldrsv int				// pb & pal save flordirt state
var tmanim bool

func florbas(img *image.NRGBA, maze *Maze, xdat Xdat, xs, ys int, one bool) {

	xb, yb := 0,0
//	img = blankimage(16*(xs-xb), 16*(ys-yb))
// one - render single tile at xs,ys
	if one { xb, yb = xs, ys;  xs, ys = xs+1, ys+1}
	// Map out where forcefield floor tiles are, so we can lay those down first
	ffmap := ffMakeMap(maze)

// ** this causes a bug with traps & ff on custom floors, it needs to be done every wp, wc, fp, fc re-assign where there is a trap/ff and should be in animate
	paletteMakeSpecial(maze.floorpattern, maze.floorcolor, maze.wallpattern, maze.wallcolor)

//	if G2 {			removed G² render

	_, _, shtamp := itemGetPNG("gfx/shadows.16.png")		// no error block on this
	xp := scanxb(xdat, 0, 0, 0, 0, "")
	Se_mflor, _,_ = parser(xp, SE_MFLR)
	if Se_mflor >= curwf { Se_mflor = -1 }
// G¹ checks
// building the ENTIRE floor everytime we come here as main maze (not palete, or pb), which is much slower
	for y := yb; y < ys; y++ {
		for x := xb; x < xs; x++ {
			adj := 0
			nwt := NOWALL | NOG1W
// Se can override these on individual tiles
			sb := scanbuf(maze.data, x, y, x, y, -2)
			xp := scanxb(xdat, x, y, x, y, "")
			wp, fp, fc := maze.wallpattern, maze.floorpattern, maze.floorcolor
			gt := G1
			p,q,_ := parser(xp, SE_G2)			// turn off G¹ if G² selected for a cell
			if p == 1 { gt = false }
			p,q,_ = parser(xp, SE_FLOR)			// set flor pat, flor col, G¹ or G²
			if p >= 0 { _,_,fp,fc = suprval(0,0,p,q) }
			p,q,_ = parser(xp, SE_WALL)			// set wall pat
			if p >= 0 { wp,_,_,_ = suprval(p,0,0,0) }

			if sb == G1OBJ_WALL_TRAP1 { nwt = NOWALL }
			if sb == G1OBJ_WALL_DESTRUCTABLE { nwt = NOWALL }

			if (nothing & nwt) == 0 {			// std wall shadows here
				adj = checkwalladj3g1(maze, xdat, x, y)	// this sets adjust for shadows, floorGetStamp sets shadows by darkening floor parts
			}

			stamp := floorGetStamp(fp, adj+rand.Intn(4), fc)
			if sb < 0 {
				coltil(img,0,x*16, y*16)		// null cell, black tile
			}
			if sb >= 0 {
			if (nothing & NOFLOOR) == 0 {
				var r int
				p,q,r = parser(xp, SE_COLRT)
				var cl uint32
				cl = 0
				if p >= 0 {
					cl = uint32(0xff000000 + r + q * 256 + p * 65536)
					coltil(img,cl,x*16, y*16)
				}
				p2,_,_ := parser(xp, SE_CFLOR)
				if p2 >= 0 && p2 < curwf {			// cust floor from png - laded by lod_maz from xb file
//					_, ux, uy := lot(x, y, x, y)
//fmt.Printf("SE_CFLOR %d - %d m: %df\n",p2,curwf,maxwf)//,wlfl.florn[fref[p2]])
					writewftoimage(img, -1,fref[p2],-1, 16,16,0,0,x,y,x*16, y*16)
				}
				p3,c,_ := parser(xp, SE_TFLOR)
				if p3 >= 0 && p3 < curwf {			// cust floor tiled in png (select tile with 'c' val) - laded by lod_maz from xb file
					bnds :=  wlfl.ftamp[fref[p3]].Bounds()
					ih := bnds.Dy()
//fmt.Printf("SE_TFLOR %d c:%d - %s, x: %d\n",p3,c,wlfl.florn[fref[p3]],ih)
					writewftoimage(img, fref[p3],-1,-1, ih,ih,0,0,c,0,x*16, y*16)
				}
				p4,_,_ := parser(xp, SE_NOFLOR)			// note: for now SEOBJ_FLOORNODRAW only works where players & monsters dont cross the tile, e.g. use SE_NOFLOR
				if p3 < 0 && p2 < 0 && p < 0 && p4 < 0 && sb != SEOBJ_FLOORNODRAW {
				if Se_mflor >= 0 {
					stamp = nil
//					_, ux, uy := lot(x, y, x, y)
					writewftoimage(img, -1,fref[Se_mflor],-1, 16,16,0,0,x, y,x*16, y*16)		// master floor replace SE_MFLR
				 } else {
					writestamptoimage(gt,img, stamp, x*16, y*16)		// G¹ floors & overrides SE_FLOR
				}}
				if p >= 0 || p2 >= 0 || p3 >= 0 || p4 >= 0 || Se_mflor >= 0 {				// cust floor or colortiles req this shadow set (for no shadow, set wp cust to 7)
					na := (adj >> 2)		// div 4
					if na > 0 && wp < shad_wallpat() {
						parimg = shtamp
						writepngtoimage(img, 16,16,0,0,na,0, x*16, y*16, 0)
					}
				}
			}}
			if ffmap[xy{x,y}] {		// are we on a forcefield beam area
				if nothing & NOTRAP == 0 {
//fmt.Printf("ffbeam %d x %d, vc: %d x %d\n ",x,y,vcx, vcy)
					stamp.ptype = "forcefield"								// this is writter over: void tiles, color tiles, cust floor
					stamp.pnum = 0
					writestamptoimage(G1,img, stamp, x*16, y*16)
				}
			}
		}
	}
fmt.Printf("rebuilt florb: %d\n",flordirt)
	flordirt = 0
}

// make walls base 

var walsb *image.NRGBA
var walsdirt int			// whether or not an edit could dirty the flor, pb & palete set to -1

func walbas(img *image.NRGBA, maze *Maze, xdat Xdat, xs, ys int, one bool) {

	xb, yb := 0,0
// one - render single tile at xs,ys
	if one { xb, yb = xs, ys;  xs, ys = xs+1, ys+1}

// seperating walls from other ents so walls dont overwrite 24 x 24 ents
// unless emu is wrong, this is the way g & G² draw walls, see screens

	xp := scanxb(xdat, 0, 0, 0, 0, "")
	Se_mwal, Se_rwal,_ = parser(xp, SE_MWAL)
//fmt.Printf("Se_mwal %d row %d\n",Se_mwal, Se_rwal)
	Se_rrnd = 0
	if Se_mwal < 0 { Se_mwal, Se_rwal, Se_rrnd = parser(xp, SE_MWALRND) }		// randomly select from wall row Se_rwal + rnd 0 - Se_rrnd val
//fmt.Printf("Se_mwalrnd %d row %d Se_rrnd %d\n",Se_mwal, Se_rwal,Se_rrnd)

	_, _, dvw := itemGetPNG("gfx/g1door_overlp.png")			// door over wall std
	for y := yb; y <= ys; y++ {
		for x := xb; x <= xs; x++ {
			var stamp *Stamp
			var dots int // dot count
			wp, wc := maze.wallpattern, maze.wallcolor
			gt := G1
			xp := scanxb(xdat, x, y, x, y, "")
			p,q,_ := parser(xp, SE_G2)
			if p == 1 { gt = false }
			p,q,_ = parser(xp, SE_WALL)
			if p >= 0 { wp,wc,_,_ = suprval(p,q,0,0) }

				//			if G2 {		removed G² render
				//	}		removed G² render
			wly, adj, walop := 0,0,0

			if !G2 {
				if wp > 5 { wp -= 6 }		// Se enhance that allows shadowless G¹ walls
			}
			nwt := NOWALL | NOG1W		// reg G¹ walls taken out by themselves (no traps, cycs etc) by NOG1W flags
			wbd := scanbuf(maze.data, x, y, x, y, -2)

			switch wbd {
			case G1OBJ_WALL_DESTRUCTABLE:
				adj, wly = checkwalladj8g1(maze, x, y)
			if (nothing & NOWALL) == 0 {
				p,q,_ = parser(xp, SE_CWAL)
				if p >= 0 && p < curwf {
					stamp = nil
					writewftoimage(img, -1,-1,wref[p], 16,16,0,0,wly+26,q, x*16, y*16)
				} else {
					if Se_mwal >= 0 {
							stamp = nil
							rn := rndr(0, Se_rrnd)
							writewftoimage(img, -1,-1,wref[Se_mwal], 16,16,0,0,wly+26,Se_rwal + rn, x*16, y*16)		// in new Se, destruct is 26 past regylar
					} else {
					stamp = wallGetDestructableStamp(wp, adj, wc)
					}
				}
				walop = wbd
			}

			case SEOBJ_SECRTWAL:
				adj, wly = checkwalladj8g1(maze, x, y)
			if (nothing & NOWALL) == 0 {
				p,q,_ = parser(xp, SE_CWAL)
				if p >= 0 && p < curwf {
					stamp = nil
					wlt := wlfl.wtamp[wref[p]]
					if !opts.Nosec {
						wlt = AdjustHue(wlfl.wtamp[wref[p]], 41.0)
					}
					parimg = wlt
					writepngtoimage(img, 16,16,0,0,wly,q, x*16, y*16,0)
				} else {
					stamp = wallGetStamp(wp, adj, wc)
					if !opts.Nosec {
						ppn := stamp.pnum + 1;		// shift secret wall display color so it cant match any wall spec
						if ppn > 16 { ppn = 0 }
						paletteSecret(ppn)
						stamp.ptype = "secret"
						stamp.pnum = 0
					}
				}
				walop = wbd
			}
			case G1OBJ_WALL_TRAP1:
				fallthrough
			case SEOBJ_WAL_TRAPCYC1:
				dots = 1; nwt = NOWALL
				fallthrough
			case SEOBJ_WAL_TRAPCYC2:
				if dots == 0 { dots = 2 }; nwt = NOWALL
				fallthrough
			case SEOBJ_WAL_TRAPCYC3:
				if dots == 0 { dots = 3 }; nwt = NOWALL
				fallthrough
			case SEOBJ_RNDWAL:
				if dots == 0 { dots = 4 }; nwt = NOWALL
				fallthrough
			case G1OBJ_WALL_REGULAR:
				adj, wly = checkwalladj8g1(maze, x, y)
				if (nothing & nwt) == 0 {
				p,q,_ = parser(xp, SE_CWAL)
				if p >= 0 && p < curwf {
					stamp = nil
					writewftoimage(img, -1,-1,wref[p], 16,16,0,0,wly,q, x*16, y*16)
				} else {
					if Se_mwal >= 0 {
							stamp = nil
							rn := rndr(0,Se_rrnd)
							writewftoimage(img, -1,-1,wref[Se_mwal], 16,16,0,0,wly,Se_rwal + rn, x*16, y*16)
					} else {
					stamp = wallGetStamp(wp, adj, wc)
					}
				}
				walop = wbd
			}
// test of some items not place in mazes - place in empty floor tile @random
			case SEOBJ_FLOOR:
				fallthrough
			case G1OBJ_TILE_FLOOR:
				p,q,r := parser(xp, SE_LETR)
				c := ""
				len := 12
				if p < 0 {
					p,q,r = parser(xp, SE_MSG)		// letter, msg mutually exclude
					if p >= 0 {
						for i := 0; i < 32; i++ {
							if xpar[i] < 130 { if xpar[i] == 0 {break}; c += map_keymap[xpar[i]]; len += 14 }
						}
					}
				} else {
					l := xpar[0]
					if l < 130 { c = map_keymap[l] }
				}
				if p >= 0 {
						gtop := gg.NewContext(len, 12)
						if err := gtop.LoadFontFace(".font/VrBd.ttf", 10); err == nil {
						gtop.Clear()
						fp, fq, fr := float64(p)/256,float64(q)/256,float64(r)/256
						gtop.SetRGB(fp, fq, fr)
						cpos := 0.5
						if len > 16 { cpos = 0.0 }
						gtop.DrawStringAnchored(c, 6, 6, cpos, 0.5)
						gtopim := gtop.Image()
						offset := image.Pt(x*16+4, y*16)
						draw.Draw(img, gtopim.Bounds().Add(offset), gtopim, image.ZP, draw.Over)
					}}
				if opts.SP {
					ts := rand.Intn(670)
					if ts == 2 { maze.data[xy{x, y}] = G1OBJ_TREASURE_BAG }
					if ts == 11 { maze.data[xy{x, y}] = MAZEOBJ_HIDDENPOT }
					if ts == 311 { maze.data[xy{x, y}] = MAZEOBJ_HIDDENPOT }
				}
			}
			if stamp != nil {
				writestamptoimage(gt,img, stamp, x*16+stamp.nudgex, y*16+stamp.nudgey)
			}
	// check door -> wall overlaps
			if wly > 0 || walop > 0 {
	//fmt.Printf("wall seg %d adj %d, type %d, dor: ",wly,adj,walop)
				for i := 0; i < 4; i++ {
					s := scanbuf(maze.data, x + dirs[i].x, y + dirs[i].y, x + dirs[i].x, y + dirs[i].y, -2)
					if (s == G1OBJ_DOOR_HORIZ && dirs[i].x != 0) || (s == G1OBJ_DOOR_VERT && dirs[i].y != 0) {
	//fmt.Printf("i(%d) %d.%d ",i, dirs[i].x, dirs[i].y)
							ovlp := dorvwal[wly][i]
							parimg = dvw
							if ovlp > 0 { writepngtoimage(img, 16,16,0,0,15+ovlp,0,x*16, y*16,0) }
					}
				}
	//fmt.Printf("\n")
			}
			if dots != 0 && nothing & NOWALL == 0 {
				renderdots(img, x*16, y*16, dots)
			}
		}
	}

	walsdirt = 0
}
// image from buffer segment			- stat: display stats 'On image' if true
// segment of buffer from xb,yb to xs,ys (begin to stop)

var fimg *image.NRGBA
var wimg *image.NRGBA
var mimg *image.NRGBA

func segimage(mdat MazeData, xdat Xdat, fdat [14]int, xb int, yb int, xs int, ys int, stat bool) *image.NRGBA {

	vlock = true
//if opts.Verbose {
fmt.Printf("segimage %dx%d - %dx%d: %t, vp: %d\n",xb,yb,xs,ys,stat,viewp)

	var err error
	var ptamp image.Image		// png stamp

// dummy maze for ops that require it
	var maze = &Maze{}
// G² edit & game will now translate to SE mode
	var skp bool
	if G2 {
		maze.data = make(map[xy]int)
		for y := 0; y <= opts.DimY; y++ {
			for x := 0; x <= opts.DimX; x++ {
				c := g2tose[mdat[xy{x, y}]]
				g1stat[c] = g2stat[mdat[xy{x, y}]]
				if mdat[xy{x, y}] > G1OBJ_EXTEND { skp = true }
				maze.data[xy{x, y}] = c
			}}
	}
	if skp || !G2 { maze.data = mdat }			// whats really wild is this just translates for the seg render system - edit still works normal

// get flags when passed
	flagbytes := make([]byte, 4)
	flagbytes[0] = byte(fdat[1])
	flagbytes[1] = byte(fdat[2])
	flagbytes[2] = byte(fdat[3])
	flagbytes[3] = byte(fdat[4])
	maze.flags = int(binary.BigEndian.Uint32(flagbytes))

	maze.wallpattern = fdat[5] & 0x0f
	maze.floorpattern = (fdat[5] & 0xf0) >> 4
	maze.wallcolor = fdat[6] & 0x0f
	maze.floorcolor = (fdat[6] & 0xf0) >> 4

	walsdirt = flordirt		// cheet for now - later these should seperate

	// unpin issue - -vals flummox canvas writes
	xba, yba := vpc_adj(xb, yb)

	fimg = blankimage(16*(xs-xb), 16*(ys-yb))		// pre-set for viewport: floors, walls, mobs
	wimg = blankimage(16*(xs-xb), 16*(ys-yb))
	mimg = blankimage(16*(xs-xb), 16*(ys-yb))

	if flordirt > 0 {
		florb = blankimage(16*(opts.DimX+1), 16*(opts.DimY+1))
		florbas(florb, maze, xdat, opts.DimX+1, opts.DimY+1,false)		//rebuild floor on load or when edit dirties it
	}

	if opts.edat < 0 || opts.edat == 2 {
		parimg = florb
		writepngtoimage(fimg, opts.DimX*16+16,opts.DimY*16+16,0,0,0,0,0,0,0)
	} else {
		parimg = florb
		sf := true
		for y := yb; y < ys; y++ {
			for x := xb; x < xs; x++ {
				_, ux, uy := lot(x, y, x, y)
				fxs, fys := xs, ys
				if fxs > opts.DimX { fxs = opts.DimX+1 }
				if fys > opts.DimY { fys = opts.DimY+1 }
				if x >= 0 && y >= 0 && x < fxs && y < fys {		// when bulk of main render is in std bounds, do super floor copy
					if sf {
fmt.Printf(" flor x,y,xs,ys %d %d %d %d ux,y %d %d, vc,y %d %d\n",(fxs-x)*16,(fys-y)*16,xs,ys,ux,uy,vcoord(x,xb,xba)*16,vcoord(y,yb,yba)*16)
						writepngtoimage(fimg,(fxs-x)*16,(fys-y)*16,0,0,ux,uy,vcoord(x,xb,xba)*16, vcoord(y,yb,yba)*16,16)
						sf = false
					}
				} else {
					writepngtoimage(fimg, 16,16,0,0,ux,uy,vcoord(x,xb,xba)*16, vcoord(y,yb,yba)*16,0)
				}
			}}
	}

	if walsdirt > 0 {
fmt.Printf("walldirt, cleen em up\n")
		walsb = blankimage(16*(opts.DimX+1), 16*(opts.DimY+1))
		walbas(walsb, maze, xdat, opts.DimX+1, opts.DimY+1,false)		//rebuild walls on load or when edit dirties it
	}

	if opts.edat < 0 || opts.edat == 2 {
		parimg = walsb
		writepngtoimage(wimg, opts.DimX*16+16,opts.DimY*16+16,0,0,0,0,0,0,0)
	} else {
		parimg = walsb
		sf := true
		for y := yb; y < ys; y++ {
			for x := xb; x < xs; x++ {
				_, ux, uy := lot(x, y, x, y)
				fxs, fys := xs, ys
				if fxs > opts.DimX { fxs = opts.DimX+1 }
				if fys > opts.DimY { fys = opts.DimY+1 }
				if x >= 0 && y >= 0 && x < fxs && y < fys {		// when bulk of main render is in std bounds, do super wall copy
					if sf {
fmt.Printf(" wals x,y,xs,ys %d %d %d %d ux,y %d %d, vc,y %d %d\n",(fxs-x)*16,(fys-y)*16,xs,ys,ux,uy,vcoord(x,xb,xba)*16,vcoord(y,yb,yba)*16)
						writepngtoimage(wimg,(fxs-x)*16,(fys-y)*16,0,0,ux,uy,vcoord(x,xb,xba)*16, vcoord(y,yb,yba)*16,16)
						sf = false
					}
				} else {
					writepngtoimage(wimg, 16,16,0,0,ux,uy,vcoord(x,xb,xba)*16, vcoord(y,yb,yba)*16,0)
				}
			}}
	}

fmt.Printf(" xb,yb,xs,ys %d %d %d %d xba,yba %d %d, dimX,y %d %d\n",xb,yb,xs,ys,xba, yba,opts.DimX,opts.DimY)

	opr := 3		// G² hack to present specials on scoreboard / info maze 104
//	_, _, sents := itemGetPNG("gfx/se_ents.16.png")			// sanct engine ent sheet
	for y := yb; y <= ys; y++ {
if opts.Verbose { fmt.Printf("\n") }
		for x := xb; x <= xs; x++ {
			var stamp *Stamp
			var dots int // dot count

			ptamp = nil
//			psx, psy, szx, szy := -1, -1, 0 ,0

			vcx, vcy := vcoord(x,xb,xba), vcoord(y,yb,yba)
			sb := scanbuf(maze.data, x, y, x, y, -2)
			xp := scanxb(xdat, x, y, x, y, "")
			gtp := G1
			p,_,_ := parser(xp, SE_G2)			// turn off G¹ if G² selected for a cell
			if p == 1 { G1 = false }			// have to literally false G¹, gtp preserves state in loop

// gen type op - letter to draw
			gtopl := ""
			gtopcol := false	// disable gen letter seperate colors
// gen type op - the context to draw
			gtop := gg.NewContext(32, 12)

// gtop font
			if err := gtop.LoadFontFace(".font/VrBd.ttf", 10); err != nil {
				panic(err)
				}

if opts.Verbose { fmt.Printf("%03d ",scanbuf(maze.data, x, y, x, y, -2)) }

			//	if G2 {			removed G² render

	if G2 {
 // hack for score table map display of: gold bag after treasure box, special potions
	if x < (ys - 1) && opts.mnum == 103 {	// dont hit past end of array & only do on score table maze

		tt := scanbuf(maze.data, x, y, x+1, y, -2)
		if sb == G1OBJ_TREASURE && tt == G1OBJ_TREASURE { maze.data[xy{x+1, y}] = G1OBJ_TREASURE_BAG }
		switch opr {
		case 1:
			if sb == G1OBJ_KEY && tt == G1OBJ_KEY {
				maze.data[xy{x, y}] = G1OBJ_X_SHTSPD
				maze.data[xy{x+1, y}] = G1OBJ_X_FIGHT
				opr--
			}
		case 2:
			if sb == G1OBJ_KEY && tt == G1OBJ_KEY {
				maze.data[xy{x, y}] = G1OBJ_X_MAGIC
				maze.data[xy{x+1, y}] = G1OBJ_X_SHOTPW
				opr--
			}
		case 3:
			if sb == G1OBJ_KEY && tt == G1OBJ_KEY {
				maze.data[xy{x, y}] = G1OBJ_X_ARMOR
				maze.data[xy{x+1, y}] = G1OBJ_X_SPEED
				opr--
			}
		}
	}}

					 // }			removed G² render
// G¹ decodes
				//	if !G2 {
// gen type op - put a letter on up left corner of every gen to indicate monsters
//		brw G - grunts
//		red D - demons
//		yel L - lobbers
//		pur S - sorceror
			gtop.Clear()
			gtopl = ""// make sure G² code (if it runs with G¹) doesnt set extra dots on non walls
			dots = 0

		if opts.edat < 1 || opts.edat == 2 {
			if x > opts.DimX || y > opts.DimY { sb = SEOBJ_FLOORNUL }
		}
// /fmt.Printf("G¹ dec: %x -- ", scanbuf(maze.data, x, y, x, y, -2))
		switch sb {

		case G1OBJ_TILE_FLOOR:
		// adj := checkwalladj3(maze, x, y) + rand.Intn(4)
		// stamp = floorGetStamp(maze.floorpattern, adj, maze.floorcolor)

		case SEOBJ_STUN:
			adj := checkwalladj3g1(maze, xdat, x, y) + rand.Intn(4)
			stamp = floorGetStamp(maze.floorpattern, adj, maze.floorcolor)
			stamp.ptype = "stun" // use trap palette (FIXME: consider moving)
			stamp.pnum = 0

		case SEOBJ_TILE_TRAP1:
			fallthrough
		case G1OBJ_TILE_TRAP1:
			dots = 1
			fallthrough
		case SEOBJ_TILE_TRAP2:
			if dots == 0 { dots = 2 }
			fallthrough
		case SEOBJ_TILE_TRAP3:
			if dots == 0 { dots = 3 }

			adj := checkwalladj3(maze, x, y) + rand.Intn(4)
			if (nothing & NOTRAP) == 0 {
				stamp = floorGetStamp(maze.floorpattern, adj, maze.floorcolor)
				stamp.ptype = "trap" // use trap palette (FIXME: consider moving)
				stamp.pnum = 0
			}

		case SEOBJ_DOOR_H:
			G1 = false; fallthrough
		case G1OBJ_DOOR_HORIZ:
			adj := checkdooradj4g1(maze, x, y)
			stamp = doorGetStamp(DOOR_HORIZ, adj)
		case SEOBJ_DOOR_V:
			G1 = false; fallthrough
		case G1OBJ_DOOR_VERT:
			adj := checkdooradj4g1(maze, x, y)
			stamp = doorGetStamp(DOOR_VERT, adj)

		case G1OBJ_PLAYERSTART:
//			arstamp[lk] = itemGetStamp("plusg1")
			if G2 { stamp = itemGetStamp("plus") }
		case G1OBJ_EXIT:
//			arstamp[lk] = itemGetStamp("exitg1")
			if G2 { stamp = itemGetStamp("exit") }
		case G1OBJ_TRANSPORTER:
//			arstamp[lk] = itemGetStamp("tportg1")
			if G2 { stamp = itemGetStamp("tport") }
		case SEOBJ_FORCEFIELDHUB:
			G1 = false
			adj := checkffadj4(maze, x, y)
			if nothing & NOEXP == 0 { stamp = ffGetStamp(adj) }
		case GORO_TEST:
			err, _, ptamp = itemGetPNG("gfx/goro.16.png")

		default:
			if opts.Verbose && false { fmt.Printf("G¹ WARNING: Unhandled obj id 0x%02x\n", sb) }
		}
// set mask flag in array
		nugetx, nugety := -4, -4
		if sb > 0 {

		if stamp != nil {
		if stamp.mask & nothing == 0 {
			g1mask[sb] = stamp.mask
// note G¹ here, opposite of other writes using gt - here gt preserves true G¹ state due to complex tile rom extract and pallet select
			writestamptoimage(G1,mimg, stamp, vcx*16+stamp.nudgex, vcy*16+stamp.nudgey)
			nugetx, nugety = stamp.nudgex, stamp.nudgey
		}} else {
//fmt.Printf("star ld %d, %v\n",sb)
		if arstamp[sb].mask & nothing == 0 {
			if arstamp[sb].pnum > -1 || arstamp[sb].pnum == -7 {
				gtopl = arstamp[sb].gtopl
//				writestamptoimage(G1,img, arstamp[sb], vcx*16+arstamp[sb].nudgex, vcy*16+arstamp[sb].nudgey)
				offset := image.Pt(vcx*16+arstamp[sb].nudgex, vcy*16+arstamp[sb].nudgey)
//if sb < 99 || sb > 100 { fmt.Printf("star ld %d, %v %v\n",sb,arstamp[sb].mimg.Bounds(),offset) }
				if arstamp[sb].mask & NOFLOOR != 0 {
					draw.Draw(fimg, arstamp[sb].mimg.Bounds().Add(offset), arstamp[sb].mimg, image.ZP, draw.Over)	// this will work, but may not be ideal
				} else {
					draw.Draw(mimg, arstamp[sb].mimg.Bounds().Add(offset), arstamp[sb].mimg, image.ZP, draw.Over)
				}
				if arstamp[sb].pnum != -7 { nugetx, nugety = arstamp[sb].nudgex, arstamp[sb].nudgey }
			}
		}}}

// Six: end G¹ decode
// if !G1 { fmt.Printf("stamp # %d - p: %s\n",scanbuf(maze.data, x, y, x, y, -2),stamp.ptype)}
// stats on palette
			if stat {			// on palette screen, show stats for loaded maze
				st := ""
				mel := sb
				st = fmt.Sprintf("%d",g1stat[mel])
//				if G2 { st = fmt.Sprintf("%d",g2stat[mel]) }
				if st != "" && stonce[mel] > 0 {
					gtop.Clear()
					gtop.SetRGB(0.5, 0.5, 0.5)
					gtop.SetRGB(1, 0, 0)
					gtop.DrawStringAnchored(st, 6, 6, 0, 0.5)
					gtopim := gtop.Image()
					if mel == G1OBJ_WALL_REGULAR { nugetx += 16; nugety += 240 }		// hackety mchakerson
					if mel == G1OBJ_TILE_FLOOR { nugetx += 16; nugety += 240 }
					offset := image.Pt(vcx*16+nugetx-5, vcy*16+nugety-5)
					draw.Draw(mimg, gtopim.Bounds().Add(offset), gtopim, image.ZP, draw.Over)
					gtopl = ""		// these seem to conflict and the palette id's box gens with monsters nearby
					stonce[mel] = 0
				}
			}
// generator monster type letter draw - only do when set
			if gtopl != "" && !opts.Nogtop {
// while each monsters gen has a letter color, some are hard to read - resetting to red
				gtop.Clear()
				if !gtopcol { gtop.SetRGB(1, 0, 0) }
				if nothing & NOGEN == 0 {
					gtop.DrawStringAnchored(gtopl, 6, 6, 0.5, 0.5)
				}
				gtopim := gtop.Image()
				offset := image.Pt(vcx*16+nugetx-4, vcy*16+nugety-4)
				draw.Draw(mimg, gtopim.Bounds().Add(offset), gtopim, image.ZP, draw.Over)
			}
// expand and sanctuary -- this is a test item that is very out of place here
			if err == nil && ptamp != nil {
				parimg = ptamp
				writepngtoimage(mimg, 16,16,0,0,0,0,vcx*16, vcy*16,0)
			}

			if dots != 0 && nothing & NOWALL == 0 {
				renderdots(mimg, (x-xb)*16, (y-yb)*16, dots)
			}
			G1 = gtp			// restore G¹ for any SE using G² turning it off
		}
	}

	g2mask[G1OBJ_WALL_REGULAR] = 2048
	g2mask[SEOBJ_SECRTWAL] = 1024
	g2mask[G1OBJ_WALL_DESTRUCTABLE] = 1024
	g2mask[SEOBJ_RNDWAL] = 1024
	g2mask[SEOBJ_WAL_TRAPCYC1] = 1024
	g2mask[SEOBJ_WAL_TRAPCYC2] = 1024
	g2mask[SEOBJ_WAL_TRAPCYC3] = 1024
	g2mask[SEOBJ_TILE_TRAP1] = 64
	g2mask[SEOBJ_TILE_TRAP2] = 64
	g2mask[SEOBJ_TILE_TRAP3] = 64
//	g2mask[] =
	g1mask[G1OBJ_WALL_REGULAR] = 2048
	g1mask[G1OBJ_WALL_DESTRUCTABLE] = 1024
	g1mask[G1OBJ_WALL_TRAP1] = 1024
	g1mask[G1OBJ_TILE_TRAP1] = 64
//	g1mask[] =
	rimg := blankimage(16*(xs-xb), 16*(ys-yb))
//savetopng("tst-img-seg.png", img)
	draw.Draw(rimg, fimg.Bounds(), fimg, image.ZP, draw.Over)
	draw.Draw(rimg, wimg.Bounds(), wimg, image.ZP, draw.Over)
	draw.Draw(rimg, mimg.Bounds(), mimg, image.ZP, draw.Over)

//	savetopng(opts.Output, img)
	vlock = false
	nobld = false
	return rimg
}
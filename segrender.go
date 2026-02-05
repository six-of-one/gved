package main

import (
	"image"
//	"image/png"
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

// scan buffer data same, 
// sx,y - starting point
// tx,y - test point
// asgn - if > -2, this is assign value
//		  when testing shadows, etc, tells where we started
//		  so slip calc (past maze edge to other side) math works

func scanbuf (mdat MazeData, sx, sy, tx, ty, asgn int) int {

	i := -1
	rf := true				// return over flows
	txe, tye := tx, ty		// for debug fmt so we know how test is adjusted

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

		if rf { i = mdat[xy{tx, ty}] }

if false && vpx < 0 {
fmt.Printf("scan: %d s-e: %d x %d, %d x %d test: %d x %d\n",i,sx,sy,txe,tye,tx,ty)
}
// the assigner, for when we need it
//		if asgn > -2 { mdat[xy{tx, ty}] = asgn }
	return i
}
/* scanbuf test out:
scanbuf s-e: 0 x 31, 1 x 30 dif: 1, 1 test: 1 x 30
scanbuf s-e: 0 x 31, -1 x 31 dif: 1, 0 test: 31 x 31
scanbuf s-e: 0 x 31, 1 x 31 dif: 1, 0 test: 1 x 31
scanbuf s-e: 0 x 31, -1 x 32 dif: 1, 1 test: 31 x 0
scanbuf s-e: 0 x 31, 0 x 32 dif: 0, 1 test: 0 x 0

scanbuf s-e: 29 x 31, 29 x 32 dif: 0, 1 test: 29 x 0
scanbuf s-e: 29 x 31, 28 x 31 dif: 1, 0 test: 28 x 31
*/

// door check

func isdoor(t int) bool {
	if t == MAZEOBJ_DOOR_HORIZ || t == MAZEOBJ_DOOR_VERT {
		return true
	} else {
		return false
	}
}

// g1 version
func isdoorg1(t int) bool {
	if t == G1OBJ_DOOR_HORIZ || t == G1OBJ_DOOR_VERT {
		return true
	} else {
		return false
	}
}

// check to see if there's walls adjacent left, left/down, and down
// used to set wall shadows, g1 engine darkens floor pixels with palette shift
// horizontal wall += 4
// diagonal wall += 8
// vertical wall += 16

// g1 version
func checkwalladj3g1(maze *Maze, x int, y int) int {
	adj := 0

	if iswallg1(scanbuf(maze.data, x, y, x-1, y, -2)) {
		adj += 4
	}

	if iswallg1(scanbuf(maze.data, x, y, x, y+1, -2)) {
		adj += 16
	}

	if iswallg1(scanbuf(maze.data, x, y, x-1, y+1, -2)) {
		adj += 8
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

// g1 version -- g2 has more walls
func checkwalladj8g1(maze *Maze, x int, y int) int {
	adj := 0

	if iswallg1(scanbuf(maze.data, x, y, x-1, y-1, -2)) {
		adj += 0x01
	}
	if iswallg1(scanbuf(maze.data, x, y, x, y-1, -2)) {
		adj += 0x02
	}
	if iswallg1(scanbuf(maze.data, x, y, x+1, y-1, -2)) {
		adj += 0x04
	}
	if iswallg1(scanbuf(maze.data, x, y, x-1, y, -2)) {
		adj += 0x08
	}
	if iswallg1(scanbuf(maze.data, x, y, x+1, y, -2)) {
		adj += 0x010
	}
	if iswallg1(scanbuf(maze.data, x, y, x-1, y+1, -2)) {
		adj += 0x20
	}
	if iswallg1(scanbuf(maze.data, x, y, x, y+1, -2)) {
		adj += 0x40
	}
	if iswallg1(scanbuf(maze.data, x, y, x+1, y+1, -2)) {
		adj += 0x80
	}

	return adj
}

// Look and see what doors are adjacent to this door
//
// Values to use:
//    up:  0x01    right:  0x02    down:  0x04    left:  0x08

// g1 version
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

func ffMark(ffmap FFMap, maze *Maze, x int, y int, dir int) {
	for i := 1; ; i++ {
		d := ffLoopDirs[dir]
		nx := x + (d.x * i)
		ny := y + (d.y * i)

		if isforcefield(maze.data[xy{nx, ny}]) {
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
	if t == MAZEOBJ_FORCEFIELDHUB {
		return true
	} else {
		return false
	}
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
//writestamptoimage(img, stamp, (x-xb+xba)*16, (y-yb+yba)*16)

// image from buffer segment			- stat: display stats if true
// segment of buffer from xb,yb to xs,ys (begin to stop)

// testing cust floor
var Se_cflr_cnt int

func segimage(mdat MazeData, fdat [11]int, xb int, yb int, xs int, ys int, stat bool) *image.NRGBA {

//if opts.Verbose {
fmt.Printf("segimage %dx%d - %dx%d: %t, vp: %d\n ",xb,yb,xs,ys,stat,viewp)

	var err error
	var ptamp image.Image		// png stamp

// dummy maze for ops that require it
	var maze = &Maze{}
	maze.data = mdat

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

	// unpin issue - - vals flummox canvas writes
	xba, yba := 0, 0
	if xb < 0 { xba = absint(xb) }
	if yb < 0 { yba = absint(yb) }

fmt.Printf("xb,yb,xs,ys %d %d %d %d xba,yba %d %d, dimX,y %d %d\n",xb,yb,xs,ys,xba, yba,opts.DimX,opts.DimY)

	// 8 pixels * 2 tiles * x,y stamps
	img := blankimage(8*2*(xs-xb), 8*2*(ys-yb))

	// Map out where forcefield floor tiles are, so we can lay those down first
	ffmap := ffMakeMap(maze)

	paletteMakeSpecial(maze.floorpattern, maze.floorcolor, maze.wallpattern, maze.wallcolor)

	if G2 {
// g2 checks
	for y := yb; y < ys; y++ {
		for x := xb; x < xs; x++ {
			adj := 0
			if maze.wallpattern < 11 {
				if (nothing & NOWALL) == 0 {		// wall shadows here
				adj = checkwalladj3(maze, x, y)
				}
			}

			stamp := floorGetStamp(maze.floorpattern, adj+rand.Intn(4), maze.floorcolor)
			if ffmap[xy{x, y}] {
				if nothing & NOTRAP == 0 {
					stamp.ptype = "forcefield"
					stamp.pnum = 0
					writestamptoimage(img, stamp, vcoord(x,xb,xba)*16, vcoord(y,yb,yba)*16)
				}
			}
			if scanbuf(maze.data, x, y, x, y, -2) < 0 {		// dont floor a null area
				coltil(img,0,vcoord(x,xb,xba)*16, vcoord(y,yb,yba)*16)
			} else {
			if (nothing & NOFLOOR) == 0 {
				writestamptoimage(img, stamp, vcoord(x,xb,xba)*16, vcoord(y,yb,yba)*16)
			}}
		}
	}} else {
// tesing Se, xpanded floor
		stdfl := true
		Se_cflr_cnt++
		if Se_cflr_cnt > 11 { Se_cflr_cnt = 1 }
		err, _, ptamp = itemGetPNG(Se_cflr[Se_cflr_cnt])
// resizing test
//		smol := image.NewRGBA(image.Rect(0, 0, ptamp.Bounds().Max.X/2, ptamp.Bounds().Max.Y/2))
//		draw.BiLinear.Scale(smol, smol.Rect, ptamp, ptamp.Bounds(), draw.Over, nil)

		bnds := ptamp.Bounds()
		iw, ih := bnds.Dx(), bnds.Dy()		// in theory this image does not HAVE to be square anymore

		tw := int(opts.Geow - 4)
		th := int(opts.Geoh - 30)

		for ty := 0; ty < th ; ty=ty+ih {
			for tx := 0; tx < tw ; tx=tx+iw {
				offset := image.Pt(tx, ty)
//				draw.Draw(img, smol.Bounds().Add(offset), smol, image.ZP, draw.Over)
				draw.Draw(img, ptamp.Bounds().Add(offset), ptamp , image.ZP, draw.Over)
			}}
// g1 checks
	for y := yb; y < ys; y++ {
		for x := xb; x < xs; x++ {
			adj := 0
			nwt := NOWALL | NOG1W
			if scanbuf(maze.data, x, y, x, y, -2) == G1OBJ_WALL_TRAP1 { nwt = NOWALL }
			if scanbuf(maze.data, x, y, x, y, -2) == G1OBJ_WALL_DESTRUCTABLE { nwt = NOWALL }
			if maze.wallpattern < 11 {
				if (nothing & nwt) == 0 {			// wall shadows here
				adj = checkwalladj3g1(maze, x, y)	// this sets adjust for shadows, floorGetStamp sets shadows by darkening floor parts
				}
			}

		  if stdfl {	// do std floor stamps
			stamp := floorGetStamp(maze.floorpattern, adj+rand.Intn(4), maze.floorcolor)
			if scanbuf(maze.data, x, y, x, y, -2) < 0 {
				coltil(img,0,(x-xb)*16, (y-yb)*16)
			} else {
			if (nothing & NOFLOOR) == 0 {
				writestamptoimage(img, stamp, vcoord(x,xb,xba)*16, vcoord(y,yb,yba)*16)
			}}
		  }
// testing
//			coltil(img,0x770077,(x-xb)*16, (y-yb)*16)
		}

	}}

// seperating walls from other ents so walls dont overwrite 24 x 24 ents
// unless emu is wrong, this is the way g & g2 draw walls, see screens
	for y := yb; y <= ys; y++ {
		for x := xb; x <= xs; x++ {
			var stamp *Stamp
			var dots int // dot count

			if G2 {
				switch whatis(maze, x, y) {
				case MAZEOBJ_WALL_DESTRUCTABLE:
					adj := checkwalladj8(maze, x, y)
				if (nothing & NOWALL) == 0 {
					stamp = wallGetDestructableStamp(maze.wallpattern, adj, maze.wallcolor)
				}
				case MAZEOBJ_WALL_SECRET:
					adj := checkwalladj8(maze, x, y)
				if (nothing & NOWALL) == 0 {
					stamp = wallGetStamp(maze.wallpattern, adj, maze.wallcolor)
					stamp.ptype = "secret"
					stamp.pnum = 0
				}
				case MAZEOBJ_WALL_TRAPCYC1:
					dots = 1
					fallthrough
				case MAZEOBJ_WALL_TRAPCYC2:
					if dots == 0 {
						dots = 2
					}
					fallthrough
				case MAZEOBJ_WALL_TRAPCYC3:
					if dots == 0 {
						dots = 3
					}
					fallthrough
				case MAZEOBJ_WALL_RANDOM:
					if dots == 0 {
						dots = 4
					}
					fallthrough
				case MAZEOBJ_WALL_REGULAR:
					adj := checkwalladj8(maze, x, y)
					if (nothing & NOWALL) == 0 {
						stamp = wallGetStamp(maze.wallpattern, adj, maze.wallcolor)
				}
// test of some items not place in mazes
				case MAZEOBJ_TILE_FLOOR:
					if opts.SP {
						ts := rand.Intn(470)
						if ts == 2 { mdat[xy{x, y}] = 	MAZEOBJ_TREASURE_BAG }		// do this with mdat
						if ts == 111 { mdat[xy{x, y}] = MAZEOBJ_HIDDENPOT }
						if ts == 311 {mdat[xy{x, y}] = MAZEOBJ_HIDDENPOT }
					}
			}}
			if G1 {
				nwt := NOWALL | NOG1W
				switch scanbuf(maze.data, x, y, x, y, -2) {
				case G1OBJ_WALL_DESTRUCTABLE:
					adj := checkwalladj8g1(maze, x, y)
				if (nothing & NOWALL) == 0 {
					stamp = wallGetDestructableStamp(maze.wallpattern, adj, maze.wallcolor)
				}

				case G1OBJ_WALL_TRAP1:
					dots = 1
					nwt = NOWALL
					fallthrough
				case G1OBJ_WALL_REGULAR:
					adj := checkwalladj8g1(maze, x, y)
					if (nothing & nwt) == 0 {
						stamp = wallGetStamp(maze.wallpattern, adj, maze.wallcolor)
				}
// test of some items not place in mazes - place in empty floor tile @random
				case MAZEOBJ_TILE_FLOOR:
					if opts.SP {
						ts := rand.Intn(470)
						if ts == 2 { mdat[xy{x, y}] = G1OBJ_TREASURE_BAG }
						if ts == 11 { mdat[xy{x, y}] = MAZEOBJ_HIDDENPOT }
						if ts == 311 {mdat[xy{x, y}] = MAZEOBJ_HIDDENPOT }
					}
				}}
			if stamp != nil {
				writestamptoimage(img, stamp, vcoord(x,xb,xba)*16+stamp.nudgex, vcoord(y,yb,yba)*16+stamp.nudgey)
			}

			if dots != 0 && nothing & NOWALL == 0 {
				renderdots(img, (x-xb)*16, (y-yb)*16, dots)
			}
		}
	}

	for y := yb; y <= ys; y++ {
if opts.Verbose { fmt.Printf("\n") }
		for x := xb; x <= xs; x++ {
			var stamp *Stamp
			var dots int // dot count

			ptamp = nil

// gen type op - letter to draw
			gtopl := ""
			gtopcol := false	// disable gen letter seperate colors
// gen type op - the context to draw
			gtop := gg.NewContext(12, 12)
// gtop font
			if err := gtop.LoadFontFace(".font/VrBd.ttf", 10); err != nil {
				panic(err)
				}

if opts.Verbose { fmt.Printf("%03d ",scanbuf(maze.data, x, y, x, y, -2)) }
// g2 decodes
			if G2 {

			// We should do better
			switch whatis(maze, x, y) {
// specials are jammed in somewhere in G2 code, we just do this
			case 70:
				stamp = itemGetStamp("speedpotion")
			case 71:
				stamp = itemGetStamp("shotpowerpotion")
			case 72:
				stamp = itemGetStamp("shotspeedpotion")
			case 73:
				stamp = itemGetStamp("shieldpotion")
			case 74:
				stamp = itemGetStamp("fightpotion")
			case 75:
				stamp = itemGetStamp("magicpotion")
			case 76:
				stamp = itemGetStamp("goldbag")

			case MAZEOBJ_TILE_FLOOR:
			// adj := checkwalladj3(maze, x, y) + rand.Intn(4)
			// stamp = floorGetStamp(maze.floorpattern, adj, maze.floorcolor)
			case MAZEOBJ_TILE_STUN:
				adj := checkwalladj3(maze, x, y) + rand.Intn(4)
				if (nothing & NOTRAP) == 0 {
					stamp = floorGetStamp(maze.floorpattern, adj, maze.floorcolor)
					stamp.ptype = "stun" // use trap palette (FIXME: consider moving)
					stamp.pnum = 0
				}

				// Tried to simplify these a bit with a goto, but golang didn't
				// like it ('jump into block'). I should figure out why.
			case MAZEOBJ_TILE_TRAP1:
				dots = 1
				fallthrough
			case MAZEOBJ_TILE_TRAP2:
				if dots == 0 {
					dots = 2
				}
				fallthrough
			case MAZEOBJ_TILE_TRAP3:
				if dots == 0 {
					dots = 3
				}
				adj := checkwalladj3(maze, x, y) + rand.Intn(4)
				if (nothing & NOTRAP) == 0 {
					stamp = floorGetStamp(maze.floorpattern, adj, maze.floorcolor)
					stamp.ptype = "trap" // use trap palette (FIXME: consider moving)
					stamp.pnum = 0
				} else { dots = 0 }
			case MAZEOBJ_WALL_MOVABLE:
				stamp = itemGetStamp("pushwall")
			case MAZEOBJ_KEY:
				stamp = itemGetStamp("key")

			case MAZEOBJ_POWER_INVIS:
				stamp = itemGetStamp("invis")
			case MAZEOBJ_POWER_REPULSE:
				stamp = itemGetStamp("repulse")
			case MAZEOBJ_POWER_REFLECT:
				stamp = itemGetStamp("reflect")
			case MAZEOBJ_POWER_TRANSPORT:
				stamp = itemGetStamp("transportability")
			case MAZEOBJ_POWER_SUPERSHOT:
				stamp = itemGetStamp("supershot")
			case MAZEOBJ_POWER_INVULN:
				stamp = itemGetStamp("invuln")

			case MAZEOBJ_DOOR_HORIZ:
				adj := checkdooradj4(maze, x, y)
				stamp = doorGetStamp(DOOR_HORIZ, adj)
			case MAZEOBJ_DOOR_VERT:
				adj := checkdooradj4(maze, x, y)
				stamp = doorGetStamp(DOOR_VERT, adj)

			case MAZEOBJ_PLAYERSTART:
				stamp = itemGetStamp("plus")
			case MAZEOBJ_EXIT:
				stamp = itemGetStamp("exit")
			case MAZEOBJ_EXITTO6:
				stamp = itemGetStamp("exit6")

			case MAZEOBJ_MONST_GHOST:
				stamp = itemGetStamp("ghost")
			case MAZEOBJ_MONST_GRUNT:
				stamp = itemGetStamp("grunt")
			case MAZEOBJ_MONST_DEMON:
				stamp = itemGetStamp("demon")
			case MAZEOBJ_MONST_LOBBER:
				stamp = itemGetStamp("lobber")
			case MAZEOBJ_MONST_SORC:
				stamp = itemGetStamp("sorcerer")
			case MAZEOBJ_MONST_AUX_GRUNT:
				stamp = itemGetStamp("auxgrunt")
			case MAZEOBJ_MONST_DEATH:
				stamp = itemGetStamp("death")
			case MAZEOBJ_MONST_ACID:
				stamp = itemGetStamp("acid")
			case MAZEOBJ_MONST_SUPERSORC:
				stamp = itemGetStamp("supersorc")
			case MAZEOBJ_MONST_IT:
				stamp = itemGetStamp("it")
			case MAZEOBJ_MONST_DRAGON:
				stamp = itemGetStamp("dragon")
			case MAZEOBJ_MONST_THIEF:
				stamp = itemGetStamp("thief")
			case MAZEOBJ_MONST_MUGGER:
				stamp = itemGetStamp("mugger")

			case MAZEOBJ_GEN_GHOST1:
				stamp = itemGetStamp("ghostgen1")
			case MAZEOBJ_GEN_GHOST2:
				stamp = itemGetStamp("ghostgen2")
			case MAZEOBJ_GEN_GHOST3:
				stamp = itemGetStamp("ghostgen3")

			case MAZEOBJ_GEN_AUX_GRUNT1:
				gtopl = "G`"
				if gtopcol { gtop.SetRGB(0.65, 0.3, 0.1) }
				stamp = itemGetStamp("generator1")
			case MAZEOBJ_GEN_GRUNT1:
				gtopl = "G"
				if gtopcol { gtop.SetRGB(0.65, 0.3, 0.1) }
				stamp = itemGetStamp("generator1")
			case MAZEOBJ_GEN_DEMON1:
				gtopl = "D"
				if gtopcol { gtop.SetRGB(1, 0, 0) }
				stamp = itemGetStamp("generator1")
			case MAZEOBJ_GEN_LOBBER1:
				gtopl = "L"
				if gtopcol { gtop.SetRGB(0.7, 0.5, 0.2) }
				stamp = itemGetStamp("generator1")
			case MAZEOBJ_GEN_SORC1:
				gtopl = "S"
				if gtopcol { gtop.SetRGB(0.37, 0.2, 0.7) }
				stamp = itemGetStamp("generator1")

			case MAZEOBJ_GEN_AUX_GRUNT2:
				gtopl = "G`"
				if gtopcol { gtop.SetRGB(0.65, 0.3, 0.1) }
				stamp = itemGetStamp("generator2")
			case MAZEOBJ_GEN_GRUNT2:
				gtopl = "G"
				if gtopcol { gtop.SetRGB(0.65, 0.3, 0.1) }
				stamp = itemGetStamp("generator2")
			case MAZEOBJ_GEN_DEMON2:
				gtopl = "D"
				if gtopcol { gtop.SetRGB(1, 0, 0) }
				stamp = itemGetStamp("generator2")
			case MAZEOBJ_GEN_LOBBER2:
				gtopl = "L"
				if gtopcol { gtop.SetRGB(0.7, 0.5, 0.2) }
				stamp = itemGetStamp("generator2")
			case MAZEOBJ_GEN_SORC2:
				gtopl = "S"
				if gtopcol { gtop.SetRGB(0.37, 0.2, 0.7) }
				stamp = itemGetStamp("generator2")

			case MAZEOBJ_GEN_AUX_GRUNT3:
				gtopl = "G`"
				if gtopcol { gtop.SetRGB(0.65, 0.3, 0.1) }
				stamp = itemGetStamp("generator3")
			case MAZEOBJ_GEN_GRUNT3:
				gtopl = "G"
				if gtopcol { gtop.SetRGB(0.65, 0.3, 0.1) }
				stamp = itemGetStamp("generator3")
			case MAZEOBJ_GEN_DEMON3:
				gtopl = "D"
				if gtopcol { gtop.SetRGB(1, 0, 0) }
				stamp = itemGetStamp("generator3")
			case MAZEOBJ_GEN_LOBBER3:
				gtopl = "L"
				if gtopcol { gtop.SetRGB(0.7, 0.5, 0.2) }
				stamp = itemGetStamp("generator3")
			case MAZEOBJ_GEN_SORC3:
				gtopl = "S"
				if gtopcol { gtop.SetRGB(0.37, 0.2, 0.7) }
				stamp = itemGetStamp("generator3")

			case MAZEOBJ_TREASURE:
				stamp = itemGetStamp("treasure")
			case MAZEOBJ_TREASURE_LOCKED:
				stamp = itemGetStamp("treasurelocked")
			case MAZEOBJ_TREASURE_BAG:
				stamp = itemGetStamp("goldbag")
			case MAZEOBJ_FOOD_DESTRUCTABLE:
				stamp = itemGetStamp("food")
			case MAZEOBJ_FOOD_INVULN:
				stamp = itemGetStamp(foods[rand.Intn(3)])
			case MAZEOBJ_POT_DESTRUCTABLE:
				stamp = itemGetStamp("potion")
			case MAZEOBJ_POT_INVULN:
				stamp = itemGetStamp("ipotion")

			case MAZEOBJ_FORCEFIELDHUB:
				adj := checkffadj4(maze, x, y)
				if nothing & NOEXP == 0 { stamp = ffGetStamp(adj) }
			case MAZEOBJ_TRANSPORTER:
				stamp = itemGetStamp("tport")

			default:
				if opts.Verbose && false { fmt.Printf("G² WARNING: Unhandled obj id 0x%02x\n", scanbuf(maze.data, x, y, x, y, -2)) }
			}
// set mask flag in array
			if scanbuf(maze.data, x, y, x, y, -2) > 0 && stamp != nil { g2mask[scanbuf(maze.data, x, y, x, y, -2)] = stamp.mask }
			}
// g1 decodes
			if G1 {
// gen type op - put a letter on up left corner of every gen to indicate monsters
//		brw G - grunts
//		red D - demons
//		yel L - lobbers
//		pur S - sorceror
				gtop.Clear()
				gtopl = ""// make sure g2 code (if it runs with g1) doesnt set extra dots on non walls
				dots = 0
// /fmt.Printf("g1 dec: %x -- ", scanbuf(maze.data, x, y, x, y, -2))
			switch scanbuf(maze.data, x, y, x, y, -2) {

			case G1OBJ_TILE_FLOOR:
			// adj := checkwalladj3(maze, x, y) + rand.Intn(4)
			// stamp = floorGetStamp(maze.floorpattern, adj, maze.floorcolor)
/*
// dont think g1 has stun tile
			case G1OBJ_TILE_STUN:
				adj := checkwalladj3g1(maze, x, y) + rand.Intn(4)
				stamp = floorGetStamp(maze.floorpattern, adj, maze.floorcolor)
				stamp.ptype = "stun" // use trap palette (FIXME: consider moving)
				stamp.pnum = 0
*/
			case G1OBJ_TILE_TRAP1:
				dots = 1

				adj := checkwalladj3(maze, x, y) + rand.Intn(4)
				if (nothing & NOTRAP) == 0 {
					stamp = floorGetStamp(maze.floorpattern, adj, maze.floorcolor)
					stamp.ptype = "trap" // use trap palette (FIXME: consider moving)
					stamp.pnum = 0
				}
			case G1OBJ_KEY:
				stamp = itemGetStamp("key")

			case G1OBJ_DOOR_HORIZ:
				adj := checkdooradj4g1(maze, x, y)
				stamp = doorGetStamp(DOOR_HORIZ, adj)
			case G1OBJ_DOOR_VERT:
				adj := checkdooradj4g1(maze, x, y)
				stamp = doorGetStamp(DOOR_VERT, adj)

			case G1OBJ_PLAYERSTART:
				stamp = itemGetStamp("plusg1")
			case G1OBJ_EXIT:
				stamp = itemGetStamp("exitg1")
			case G1OBJ_EXIT4:
				stamp = itemGetStamp("exit4")
			case G1OBJ_EXIT8:
				stamp = itemGetStamp("exit8")

			case G1OBJ_MONST_GHOST1:
				stamp = itemGetStamp("ghost1")
			case G1OBJ_MONST_GHOST2:
				stamp = itemGetStamp("ghost2")
			case G1OBJ_MONST_GHOST3:
				stamp = itemGetStamp("ghost")
			case G1OBJ_MONST_GRUNT1:
				stamp = itemGetStamp("grunt1")
			case G1OBJ_MONST_GRUNT2:
				stamp = itemGetStamp("grunt2")
			case G1OBJ_MONST_GRUNT3:
				stamp = itemGetStamp("grunt")
			case G1OBJ_MONST_DEMON1:
				stamp = itemGetStamp("demon1")
			case G1OBJ_MONST_DEMON2:
				stamp = itemGetStamp("demon2")
			case G1OBJ_MONST_DEMON3:
				stamp = itemGetStamp("demon")
			case G1OBJ_MONST_LOBBER1:
				stamp = itemGetStamp("lobber1")
			case G1OBJ_MONST_LOBBER2:
				stamp = itemGetStamp("lobber2")
			case G1OBJ_MONST_LOBBER3:
				stamp = itemGetStamp("lobber")
			case G1OBJ_MONST_SORC1:
				stamp = itemGetStamp("sorcerer1")
			case G1OBJ_MONST_SORC2:
				stamp = itemGetStamp("sorcerer2")
			case G1OBJ_MONST_SORC3:
				stamp = itemGetStamp("sorcerer")
			case G1OBJ_MONST_DEATH:
				stamp = itemGetStamp("death")

			case G1OBJ_MONST_THIEF:
				stamp = itemGetStamp("thief")

			case G1OBJ_GEN_GHOST1:
				stamp = itemGetStamp("ghostgen1")
			case G1OBJ_GEN_GHOST2:
				stamp = itemGetStamp("ghostgen2")
			case G1OBJ_GEN_GHOST3:
				stamp = itemGetStamp("ghostgen3")

// if a clear is done after, this SetRGB set bkg somehow
			case G1OBJ_GEN_GRUNT1:
				gtopl = "G"
				if gtopcol { gtop.SetRGB(0.65, 0.3, 0.1) }
				stamp = itemGetStamp("generator1")
			case G1OBJ_GEN_DEMON1:
				gtopl = "D"
				if gtopcol { gtop.SetRGB(1, 0, 0) }
				stamp = itemGetStamp("generator1")
			case G1OBJ_GEN_LOBBER1:
				gtopl = "L"
				if gtopcol { gtop.SetRGB(0.7, 0.5, 0.2) }
				stamp = itemGetStamp("generator1")
			case G1OBJ_GEN_SORC1:
				gtopl = "S"
				if gtopcol { gtop.SetRGB(0.37, 0.2, 0.7) }
				stamp = itemGetStamp("generator1")

			case G1OBJ_GEN_GRUNT2:
				gtopl = "G"
				if gtopcol { gtop.SetRGB(0.65, 0.3, 0.1) }
				stamp = itemGetStamp("generator2")
			case G1OBJ_GEN_DEMON2:
				gtopl = "D"
				if gtopcol { gtop.SetRGB(1, 0, 0) }
				stamp = itemGetStamp("generator2")
			case G1OBJ_GEN_LOBBER2:
				gtopl = "L"
				if gtopcol { gtop.SetRGB(0.7, 0.5, 0.2) }
				stamp = itemGetStamp("generator2")
			case G1OBJ_GEN_SORC2:
				gtopl = "S"
				if gtopcol { gtop.SetRGB(0.37, 0.2, 0.7) }
				stamp = itemGetStamp("generator2")

			case G1OBJ_GEN_GRUNT3:
				gtopl = "G"
				if gtopcol { gtop.SetRGB(0.65, 0.3, 0.1) }
				stamp = itemGetStamp("generator3")
			case G1OBJ_GEN_DEMON3:
				gtopl = "D"
				if gtopcol { gtop.SetRGB(1, 0, 0) }
				stamp = itemGetStamp("generator3")
			case G1OBJ_GEN_LOBBER3:
				gtopl = "L"
				if gtopcol { gtop.SetRGB(0.7, 0.5, 0.2) }
				stamp = itemGetStamp("generator3")
			case G1OBJ_GEN_SORC3:
				gtopl = "S"
				if gtopcol { gtop.SetRGB(0.37, 0.2, 0.7) }
				stamp = itemGetStamp("generator3")

			case G1OBJ_TREASURE:
				stamp = itemGetStamp("treasure")
			case G1OBJ_TREASURE_BAG:
				stamp = itemGetStamp("goldbag")
			case G1OBJ_FOOD_DESTRUCTABLE:
				stamp = itemGetStamp("food")
			case G1OBJ_FOOD_INVULN:
				stamp = itemGetStamp(foods[rand.Intn(3)])
			case G1OBJ_POT_DESTRUCTABLE:
				stamp = itemGetStamp("potion")
			case G1OBJ_POT_INVULN:
				stamp = itemGetStamp("ipotion")
			case G1OBJ_INVISIBL:
				stamp = itemGetStamp("invis")
// specials added after convert to se id'ed them on maze 115, score table block
			case G1OBJ_X_SPEED:
				stamp = itemGetStamp("speedpotion")
			case G1OBJ_X_SHOTPW:
				stamp = itemGetStamp("shotpowerpotion")
			case G1OBJ_X_SHTSPD:
				stamp = itemGetStamp("shotspeedpotion")
			case G1OBJ_X_ARMOR:
				stamp = itemGetStamp("shieldpotion")
			case G1OBJ_X_FIGHT:
				stamp = itemGetStamp("fightpotion")
			case G1OBJ_X_MAGIC:
				stamp = itemGetStamp("magicpotion")

			case G1OBJ_TRANSPORTER:
				stamp = itemGetStamp("tportg1")
			case GORO_TEST:
				err, _, ptamp = itemGetPNG("gfx/goro.16.png")

			default:
				if opts.Verbose && false { fmt.Printf("G¹ WARNING: Unhandled obj id 0x%02x\n", scanbuf(maze.data, x, y, x, y, -2)) }
			}
// set mask flag in array
			if scanbuf(maze.data, x, y, x, y, -2) > 0 && stamp != nil { g1mask[scanbuf(maze.data, x, y, x, y, -2)] = stamp.mask }
		}
// Six: end G1 decode
			if stamp != nil {
				writestamptoimage(img, stamp, vcoord(x,xb,xba)*16+stamp.nudgex, vcoord(y,yb,yba)*16+stamp.nudgey)
// stats on palette
				if stat {			// on palette screen, show stats for loaded maze
					st := ""
					mel := scanbuf(maze.data, x, y, x, y, -2)
					if G1 { st = fmt.Sprintf("%d",g1stat[mel]) }
					if G2 { st = fmt.Sprintf("%d",g2stat[mel]) }
					if st != "" && stonce[mel] > 0 {
						gtop.Clear()
						gtop.SetRGB(0.5, 0.5, 0.5)
						gtop.SetRGB(1, 0, 0)
						gtop.DrawStringAnchored(st, 6, 6, 0.5, 0.5)
						gtopim := gtop.Image()
						offset := image.Pt((x-xb)*16+stamp.nudgex+15, (y-yb)*16+stamp.nudgey-5)
						draw.Draw(img, gtopim.Bounds().Add(offset), gtopim, image.ZP, draw.Over)
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
					offset := image.Pt((x-xb)*16+stamp.nudgex-4, (y-yb)*16+stamp.nudgey-4)
					draw.Draw(img, gtopim.Bounds().Add(offset), gtopim, image.ZP, draw.Over)
				}
			}
// expand and sanctuary
			if err == nil && ptamp != nil {

//				r := image.Rect(0, 0, 11, 11)
//				dst := image.NewRGBA(r)
//				draw.Copy(dst, image.Pt(-4,-4), ptamp, r, draw.Over, nil)

				offset := image.Pt((x-xb)*16, (y-yb)*16)
//				draw.Draw(img, dst.Bounds().Add(offset), dst, image.ZP, draw.Over)
				draw.Draw(img, ptamp.Bounds().Add(offset), ptamp, image.ZP, draw.Over)
			}

			if dots != 0 && nothing & NOWALL == 0 {
				renderdots(img, (x-xb)*16, (y-yb)*16, dots)
			}
		}
	}

	g2mask[MAZEOBJ_WALL_REGULAR] = 2048
	g2mask[MAZEOBJ_WALL_SECRET] = 1024
	g2mask[MAZEOBJ_WALL_DESTRUCTABLE] = 1024
	g2mask[MAZEOBJ_WALL_RANDOM] = 1024
	g2mask[MAZEOBJ_WALL_TRAPCYC1] = 1024
	g2mask[MAZEOBJ_WALL_TRAPCYC2] = 1024
	g2mask[MAZEOBJ_WALL_TRAPCYC3] = 1024
	g2mask[MAZEOBJ_TILE_TRAP1] = 64
	g2mask[MAZEOBJ_TILE_TRAP2] = 64
	g2mask[MAZEOBJ_TILE_TRAP3] = 64
//	g2mask[] =
	g1mask[G1OBJ_WALL_REGULAR] = 2048
	g1mask[G1OBJ_WALL_DESTRUCTABLE] = 1024
	g1mask[G1OBJ_WALL_TRAP1] = 1024
	g1mask[G1OBJ_TILE_TRAP1] = 64
//	g1mask[] =

	savetopng(opts.Output, img)
// for user select
	return img
}
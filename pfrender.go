package main

import (
	"image"
	"image/png"
	"math/rand"
	"fmt"
	"os"
	"image/draw"
	"github.com/fogleman/gg"
	"image/color"
	"encoding/binary"
)


// arrays for item masks
var g1mask [256]int
var g2mask [256]int

// Okay, so we have a maze. We need to adjust the edges to take care of
// wrap or no wrap.
//
// If we're not wrapping in a direction, we need to duplicate the left and
// top walls, so that the maze will be enclosed.
// six - note: this is only for appearance of viewer...
// - and needs to be done in a way the editor 1. doesnt save, 2. no editing, 3. drawn differently
func copyedges(maze *Maze) {
	for i := 0; i <= 32; i++ {
		if (maze.flags & LFLAG4_WRAP_H) == 0 {
			maze.data[xy{32, i}] = maze.data[xy{0, i}]
		} else {
			maze.data[xy{32, i}] = 0
		}
			if opts.edat < 1 || opts.edip == 0 { ebuf[xy{32, i}] = maze.data[xy{32, i}] } else {
				maze.data[xy{32, i}] = ebuf[xy{32, i}]
			}	// have to do edit buffer as well
	}

	for i := 0; i <= 32; i++ {
		if (maze.flags & LFLAG4_WRAP_V) == 0 {
			maze.data[xy{i, 32}] = maze.data[xy{i, 0}]
		} else {
			maze.data[xy{i, 32}] = 0
		}
			if opts.edat < 1 || opts.edip == 0 { ebuf[xy{i, 32}] = maze.data[xy{i, 32}] } else {
				maze.data[xy{i, 32}] = ebuf[xy{i, 32}]
			}
	}
}

// for maze output to se
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

// six tile dumper fn
func writile(stamp *Stamp, tbas int, tbaddr int, sz int , ada int) {

 //	fmt.Printf("tbas pass %d\n",tbas)
 // exit tile is special build
	if ada != 0x7f0 { stamp.numbers = tilerange(tbas, tbaddr) }
	fillstamp(stamp)

 // file name with addr
	wnam := fmt.Sprintf(".p%d/tl_%05d_%04X.png",stamp.pnum,tbas + ada,tbas + ada)
 // -sz = use (s)pecial file designation
	if sz < 0 {
		sz = max(sz,-sz)
		wnam = fmt.Sprintf(".p%d/tl_s%05d_%04X.png",stamp.pnum,tbas + ada,tbas + ada)
	}
 // for 8x8 single tile, place is sub color dirs sep from .p*
	if sz == 8 {
		wnam = fmt.Sprintf(".8x8/c%d/i%05d_%04X.png",stamp.pnum,tbas + ada,tbas + ada)
	}
 // 24 pixels * 24 pixels - temp write out of all tiles
 // impl: 16 x 16 for the 2 x 2 tiles, and dragon size for hims (4 x 4)
 // special 8 x 8 tiles for unit list
	wimg := blankimage(sz, sz)
	writestamptoimage(wimg, stamp, 0, 0)
	wrfile, err := os.Create(wnam)
	if err == nil {
		png.Encode(wrfile,wimg)
	}
	wrfile.Close()
}

var foods = []string{"ifood1", "ifood2", "ifood3"}
var nothing int

func genpfimage(maze *Maze, mazenum int) *image.NRGBA {
	extrax, extray := 0, 0	// this becomes the space for copyedges walls...
	if opts.Wob {			// and this extra space is an issue to blotter & measure, it either always has to be or not
		extrax = 16			// - of course inimical to edit system, so it has to be accounted
		extray = 16			// - but now only in view mode
	}

// no {floor|wall} - only things
// option to generate image with no floors or walls (say for color correcting g1 mazes we dont have color codes for)
// nofloor = 1, nowall = 2, no both = 3
// maybe make a cli switch?

	// 8 pixels * 2 tiles * 32 stamps, plus extra space on edges
	xspc := 32		// this where old viewer drew passage 'arrows'
	xpad := 16
// dont draw the arrow space border
	if !opts.Aob { xspc = 0; xpad = 0 }
	img := blankimage(8*2*32+xspc+extrax, 8*2*32+xspc+extray)

	// Map out where forcefield floor tiles are, so we can lay those down first
	ffmap := ffMakeMap(maze)

	// mazes will always be the same size, so just use constants
	// maze := mazeDecompress(mazedata)
	copyedges(maze)
	paletteMakeSpecial(maze.floorpattern, maze.floorcolor, maze.wallpattern, maze.wallcolor)

	if G2 {
// g2 checks
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
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
					writestamptoimage(img, stamp, x*16+xpad, y*16+xpad)
				}
			}
			if (nothing & NOFLOOR) == 0 {
				writestamptoimage(img, stamp, x*16+xpad, y*16+xpad)
			}
		}
	}} else {
// g1 checks
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			adj := 0
			nwt := NOWALL | NOG1W
			if whatis(maze, x, y) == G1OBJ_WALL_TRAP1 { nwt = NOWALL }
			if whatis(maze, x, y) == G1OBJ_WALL_DESTRUCTABLE { nwt = NOWALL }
			if maze.wallpattern < 11 {
				if (nothing & nwt) == 0 {		// wall shadows here
				adj = checkwalladj3g1(maze, x, y)
				}
			}

			stamp := floorGetStamp(maze.floorpattern, adj+rand.Intn(4), maze.floorcolor)
			if (nothing & NOFLOOR) == 0 {
				writestamptoimage(img, stamp, x*16+xpad, y*16+xpad)
			}
		}

	}}

	lastx := 32
	lasty := 32

	if maze.flags&LFLAG4_WRAP_H > 0 {
		lastx = 31
	}

	if maze.flags&LFLAG4_WRAP_V > 0 {
		lasty = 31
	}

// seperating walls from other ents so walls dont overwrite 24 x 24 ents
// unless emu is wrong, this is the way g & g2 draw walls, see screens
	for y := 0; y <= lasty; y++ {
		for x := 0; x <= lastx; x++ {
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
						if ts == 2 { maze.data[xy{x, y}] = 	MAZEOBJ_TREASURE_BAG }
						if ts == 111 { maze.data[xy{x, y}] = MAZEOBJ_HIDDENPOT }
						if ts == 311 { maze.data[xy{x, y}] = MAZEOBJ_HIDDENPOT }
					}
			}}
			if G1 {
				nwt := NOWALL | NOG1W
				switch whatis(maze, x, y) {
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
						if ts == 2 { maze.data[xy{x, y}] = G1OBJ_TREASURE_BAG }
						if ts == 11 { maze.data[xy{x, y}] = MAZEOBJ_HIDDENPOT }
						if ts == 311 { maze.data[xy{x, y}] = MAZEOBJ_HIDDENPOT }
					}
				}}
			if stamp != nil {
				writestamptoimage(img, stamp, x*16+xpad+stamp.nudgex, y*16+xpad+stamp.nudgey)
			}

			if dots != 0 && nothing & NOWALL == 0 {
				renderdots(img, x*16+xpad, y*16+xpad, dots)
			}
		}
	}

	opr := 3				// 3 sets of special potions
// main maze decode loop - op over all maze cells
	for y := 0; y <= lasty; y++ {
		for x := 0; x <= lastx; x++ {
			var stamp *Stamp
			var dots int // dot count
// gen type op - letter to draw
			gtopl := ""
			gtopcol := false	// disable gen letter seperate colors
// gen type op - the context to draw
			gtop := gg.NewContext(12, 12)
// gtop font
			if err := gtop.LoadFontFace(".font/VrBd.ttf", 14); err != nil {
				panic(err)
				}
// g2 decodes
			if G2 {
 // hack for score table map display of: gold bag after treasure box, special potions
	if x < (lastx - 1) && mazenum == 103 {	// dont hit past end of array & only do on score table maze
		ts := maze.data[xy{x, y}]
		tt := maze.data[xy{x+1, y}]
		if ts == MAZEOBJ_TREASURE && tt == MAZEOBJ_TREASURE { maze.data[xy{x+1, y}] = 76 }
		switch opr {
		case 1:
			if ts == MAZEOBJ_KEY && tt == MAZEOBJ_KEY {
				maze.data[xy{x, y}] = 72
				maze.data[xy{x+1, y}] = 74
				opr--
			}
		case 2:
			if ts == MAZEOBJ_KEY && tt == MAZEOBJ_KEY {
				maze.data[xy{x, y}] = 75
				maze.data[xy{x+1, y}] = 71
				opr--
			}
		case 3:
			if ts == MAZEOBJ_KEY && tt == MAZEOBJ_KEY {
				maze.data[xy{x, y}] = 73
				maze.data[xy{x+1, y}] = 70
				opr--
			}
		}
	}
			// We should do better
			switch whatis(maze, x, y) {
 // specials are jammed in somewhere in G2 code, we just do this
 // specials added after convert to se id'ed them on maze 115, score table block
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
 // testing special potions
			case MAZEOBJ_HIDDENPOT:
				if opts.SP {
					ts := rand.Intn(6)
					switch ts {
					case 1:
						stamp = itemGetStamp("speedpotion")
					case 2:
						stamp = itemGetStamp("shotpowerpotion")
					case 3:
						stamp = itemGetStamp("shotspeedpotion")
					case 4:
						stamp = itemGetStamp("shieldpotion")
					case 5:
						stamp = itemGetStamp("fightpotion")
					case 6:
						stamp = itemGetStamp("magicpotion")
					}
				}
			default:
				if opts.Verbose && false { fmt.Printf("G² WARNING: Unhandled obj id 0x%02x\n", whatis(maze, x, y)) }
			}
 // set mask flag in array
			if whatis(maze, x, y) > 0 && stamp != nil { g2mask[whatis(maze, x, y)] = stamp.mask }
			}
// end G2 decodes

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
 // /fmt.Printf("g1 dec: %x -- ", whatis(maze, x, y))
			fn := whatis(maze, x, y)
			if x == zm_x && y == zm_y { fn = G1OBJ_WIZARD }
			switch fn {

			case G1OBJ_TILE_FLOOR:
			// adj := checkwalladj3(maze, x, y) + rand.Intn(4)
			// stamp = floorGetStamp(maze.floorpattern, adj, maze.floorcolor)
 // dont think g1 has stun tile
			case G1OBJ_TILE_STUN:
				adj := checkwalladj3g1(maze, x, y) + rand.Intn(4)
				stamp = floorGetStamp(maze.floorpattern, adj, maze.floorcolor)
				stamp.ptype = "stun" // use trap palette (FIXME: consider moving)
				stamp.pnum = 0

			case G1OBJ_TILE_TRAP1:
				dots = 1
 //				fallthrough
	/*
			case G1OBJ_TILE_TRAP2:
				if dots == 0 {
					dots = 2
				}
				fallthrough
			case G1OBJ_TILE_TRAP3:
				if dots == 0 {
					dots = 3
				}
	*/
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
 // when death shows up on score board maze, hack theif into sample area
				if mazenum == 115 {
					maze.data[xy{x, y + 3}] = G1OBJ_MONST_THIEF
				}
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
 // testing special potions - seperate from actual placed in maze above, this is the view tester option
			case MAZEOBJ_HIDDENPOT:
				if opts.SP {
					ts := rand.Intn(6)
					switch ts {
					case 1:
						stamp = itemGetStamp("speedpotion")
					case 2:
						stamp = itemGetStamp("shotpowerpotion")
					case 3:
						stamp = itemGetStamp("shotspeedpotion")
					case 4:
						stamp = itemGetStamp("shieldpotion")
					case 5:
						stamp = itemGetStamp("fightpotion")
					case 6:
						stamp = itemGetStamp("magicpotion")
					}
				}
			case G1OBJ_WIZARD:
				stamp = itemGetStamp("wizard")
			default:
				if opts.Verbose && false { fmt.Printf("G¹ WARNING: Unhandled obj id 0x%02x\n", whatis(maze, x, y)) }
			}
 // set mask flag in array
			if whatis(maze, x, y) > 0 && stamp != nil { g1mask[whatis(maze, x, y)] = stamp.mask }
		}
// Six: end G1 decode
			if stamp != nil {
				writestamptoimage(img, stamp, x*16+xpad+stamp.nudgex, y*16+xpad+stamp.nudgey)
// generator monster type letter draw - only do when set
				if gtopl != "" && !opts.Nogtop {
// while each monsters gen has a letter color, some are hard to read - resetting to red
					gtop.Clear()
					if !gtopcol { gtop.SetRGB(1, 0, 0) }
					if nothing & NOGEN == 0 {
						gtop.DrawStringAnchored(gtopl, 6, 6, 0.5, 0.5)
					}
					gtopim := gtop.Image()
					offset := image.Pt(x*16+xpad+stamp.nudgex-4, y*16+xpad+stamp.nudgey-4)
					draw.Draw(img, gtopim.Bounds().Add(offset), gtopim, image.ZP, draw.Over)
				}
			}

			if dots != 0 && nothing & NOWALL == 0 {
				renderdots(img, x*16+xpad, y*16+xpad, dots)
			}
		}
	}

// set #T hider/includer flags that arent set above
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



/// Six - individual tile dumper
// - these can be used to make sprite sheets and the like

/*******************************************
 TTT           ___  DDD                __
  T   i  L    E     D  d  u  u  m\/m  P  p
  T   i  L    E--   D  d  u  u  m  m  p__p
  T   i  LLL  E___  DDD    uu   m  m  p

 *******************************************/
 // lock out with this - set to true to run dump
 // run decode for maze116
 // written to .p[0-31]/tl_%05d_%04X.png
 // where %d and %X are tile start addr

 // loops thru all palette nums for all valid tiles, that are known
 // atempts to correct for various 2x2, 3x3, 4x4 (dragon) and odd size doors
if false {
 // counter for tiles - imprv - dont write dups
	wcnt := 1
 // tb adder controls size of tile render, and mem skip to next tile
 // this could also control the render out 16x16, 24x24 or 32x32
	tbaddr := 9
	var stamp *Stamp
	stamp = itemGetStamp("ghost")
	stamp.pnum = 0

 // 0000 - 1FFF dump

	for stamp != nil && mazenum == 116 {

		wcnt++
		for i := 0x800; i < 0xfff; i += tbaddr {

			tbaddr = 9
			stamp.width = 3
			if i & 0xff == 0xfc {
				stamp.width = 2
				writile(stamp, i, 4, 16 ,-0x800)
				tbaddr = 4

			} else {
				writile(stamp, i, tbaddr, 24 ,-0x800)
			}
		}
		tbaddr = 9
		stamp.width = 3
		for i := 0x7b3; i < 0x7e2; i += tbaddr {

			writile(stamp, i, 9, 24 ,0x800)
		}
		stamp.width = 2
		tbaddr = 4
		for i := 0x011; i < 0x1e0; i += tbaddr {

			writile(stamp, i, 4, 16 ,0x800)
		}
		for i := 0x1e0; i < 0x39c; i += tbaddr {

			writile(stamp, i, 4, 16 ,0x800)
		}
		for i := 0x39e; i < 0x4b0; i += tbaddr {

			writile(stamp, i, 4, 16 ,0x800)
		}
		for i := 0x73b; i < 0x7b2; i += tbaddr {

			writile(stamp, i, 4, 16 ,0x800)
		}
		for i := 0x1800; i < 0x1c47; i += tbaddr {

			tbaddr = 9
			stamp.width = 3
			if i > 0x1b50 && i < 0x1c47 { tbaddr = 6 }
			if i == 0x1bc3 { tbaddr = 4 }
			if i & 0xff == 0xfc || tbaddr < 6 {
				stamp.width = 2
				writile(stamp, i, tbaddr, 16 ,-0x800)
				tbaddr = 4

			} else {
				writile(stamp, i, tbaddr, 24 ,-0x800)
			}
			if i == 0x1bfd { i = i - 3 }
		}
		for i := 0x1000; i < 0x17ff; i += tbaddr {

			tbaddr = 9
			stamp.width = 3
			if i & 0xff == 0xfc {
				stamp.width = 2
				writile(stamp, i, 4, 16 ,0x800)
				tbaddr = 4

			} else {
				writile(stamp, i, tbaddr, 24 ,0x800)
			}
		}

		stamp.width = 2
		tbaddr = 4
		for i := 0x1c48; i < 0x1c87; i += tbaddr {
			writile(stamp, i, 4, 16 ,-0x800)
		}
		for i := 0x1c8b; i < 0x1d3b; i += tbaddr {
			writile(stamp, i, 4, 16 ,-0x800)
			if i == 0x1cfb { i++ }
		}
		writile(stamp, 0x1d48, 4, 16 ,-0x800)
 // so doors are an absolute mess, with 2, 3, and 4 wide all mixed up 
		stamp.width = 3
		writile(stamp, 0x1d3c, 6, 24,-0x800)
		writile(stamp, 0x1d42, 6, 24,-0x800)
		writile(stamp, 0x1d4c, 6, 24,-0x800)
		writile(stamp, 0x1d5a, 6, 24,-0x800)
		writile(stamp, 0x1d68, 6, 24,-0x800)
		writile(stamp, 0x1d76, 6, 24,-0x800)

		stamp.width = 4
		writile(stamp, 0x1d52, 8, 32,-0x800)
		writile(stamp, 0x1d60, 8, 32,-0x800)
		writile(stamp, 0x1d6e, 8, 32,-0x800)
 // some of these doors are wrong
		for i := 0x1d7c; i < 0x1db3; i += tbaddr {
			writile(stamp, i, 4, 16 ,-0x800)
		}
		tbaddr = 9
		stamp.width = 3
		for i := 0x1dcf; i < 0x1dfb; i += tbaddr {
			writile(stamp, i, tbaddr, 24 ,-0x800)
		}
		writile(stamp, 0x1e00, 9, 24,-0x800)
		writile(stamp, 0x1e0d, 6, 24,-0x800)
		stamp.width = 2
		writile(stamp, 0x1e09, 4, 16 ,-0x800)
		stamp.width = 3
		writile(stamp, 0x1c88, 3, 24,-0x800)
		tbaddr = 3
		for i := 0x1db4; i < 0x1dce; i += tbaddr {
			writile(stamp, i, tbaddr, 24 ,-0x800)
		}

 // every loop, increase palette # to next till end
		stamp.pnum++;

 // done, no further pallets

		if stamp.pnum == 32 {
			stamp = nil
		}
	}

 // gauntlet 2 add ins not handled by 0000 - 1FFF dump
	pnum := 0

	for pnum < 32 && mazenum == 116 {

		stamp = itemGetStamp("pushwall")
		stamp.pnum = pnum
		writile(stamp, 0x20f6, 6, 24 ,0)
		stamp = itemGetStamp("pfood")
		stamp.pnum = pnum
		writile(stamp, 0x25ed, 9, 24 ,0)
		stamp = itemGetStamp("ppotion")
		stamp.pnum = pnum
		writile(stamp, 0x20fc, 4, 16 ,0)
		stamp = itemGetStamp("mfood")
		stamp.pnum = pnum
		writile(stamp, 0x277b, 9, 24 ,0)
		stamp = itemGetStamp("treasurelocked")
		stamp.pnum = pnum
		writile(stamp, 0x25e4, 9, 24 ,0)

 // g2 temp powers
		stamp = itemGetStamp("transportability")
		stamp.pnum = pnum
		writile(stamp, 0x23fc, 4, 16 ,0)
		stamp = itemGetStamp("reflect")
		stamp.pnum = pnum
		writile(stamp, 0x24fc, 4, 16 ,0)
		stamp = itemGetStamp("repulse")
		stamp.pnum = pnum
		writile(stamp, 0x26fc, 4, 16 ,0)
		stamp = itemGetStamp("invuln")
		stamp.pnum = pnum
		writile(stamp, 0x2784, 4, 16 ,0)
		stamp = itemGetStamp("supershot")
		stamp.pnum = pnum
		writile(stamp, 0x2788, 4, 16 ,0)
 // g1 powers
 // handled by loop now

		stamp = itemGetStamp("it")
		tbaddr = 9
		stamp.pnum = pnum
		for i := 0x2600; i < 0x2690; i += tbaddr {

			writile(stamp, i, tbaddr, 24 ,0)
		}
 // pickles
		for i := 0x2300; i < 0x23fb; i += tbaddr {

			writile(stamp, i, tbaddr, 24 ,0)
		}
		writile(stamp, 0x25db, tbaddr, 24 ,0)
 // ?
		for i := 0x2400; i < 0x24fb; i += tbaddr {

			writile(stamp, i, tbaddr, 24 ,0)
		}
		for i := 0x2690; i < 0x26fb; i += tbaddr {

			writile(stamp, i, tbaddr, 24 ,0)
		}
 /* < x2000 handled by 1st loop now
		for i := 0x1fab; i < 0x1ffb; i += tbaddr {

			writile(stamp, i, tbaddr, 24 ,0)
		}
		for i := 0x15cf; i < 0x1608; i += tbaddr {

			writile(stamp, i, tbaddr, 24 ,0)
		}
 */

 // dragon breath
		for i := 0x278c; i < 0x27f7; i += tbaddr {

			writile(stamp, i, tbaddr, 24 ,0)
		}

		stamp = itemGetStamp("dragon")
		tbaddr = 16
		stamp.pnum = pnum
		for i := 0x2100; i < 0x2300; i += tbaddr {

			writile(stamp, i, tbaddr, 32 ,0)
		}
		for i := 0x2500; i < 0x2560; i += tbaddr {

			writile(stamp, i, tbaddr, 32 ,0)
		}
		for i := 0x2740; i < 0x2760; i += tbaddr {

			writile(stamp, i, tbaddr, 32 ,0)
		}
		tbaddr = 2
		stamp.width = 2
		for i := 0x25f6; i < 0x25ff; i += tbaddr {

			writile(stamp, i, tbaddr, 16 ,0)
		}

 // only needed once - pnum 0, even tho this isnt the palette num on these
		if pnum == 0 {
			tbaddr = 4
			stamp = itemGetStamp("exit4")
			writile(stamp, 0x4fc, tbaddr, -16 ,0x800)
			stamp = itemGetStamp("exit8")
			writile(stamp, 0x5fc, tbaddr, -16 ,0x800)
			writile(stamp, 0x3fc, tbaddr, -16 ,0x800)
			stamp = itemGetStamp("exit")
			writile(stamp, 0x39e, tbaddr, -16 ,0x7f0)
			for i := 0x39e; i < 0x49d; i += tbaddr {

				writile(stamp, i, tbaddr, -16 ,0x800)
			}
			stamp = itemGetStamp("tport")
			for i := 0x49e; i < 0x4b2; i += tbaddr {

				writile(stamp, i, tbaddr, -16 ,0x800)
			}
		}

 // the single 8x8 tile set is locked out for normal tile writing ops
 // set true to run
 if false {
 // single tile, for all the issues
 // this is written to .8x8/c[0-31]/i%05d_%04X.png
		stamp = itemGetStamp("ghost")
		tbaddr = 1
		stamp.pnum = pnum
		stamp.width = 1
		for i := 0x0; i < 0x7ff; i += tbaddr {

			writile(stamp, i, tbaddr, 8 ,0x800)
		}
		for i := 0x800; i < 0xfff; i += tbaddr {

			writile(stamp, i, tbaddr, 8 ,-0x800)
		}
		for i := 0x1000; i < 0x17ff; i += tbaddr {

			writile(stamp, i, tbaddr, 8 ,0x800)
		}
		for i := 0x1800; i < 0x1fff; i += tbaddr {

			writile(stamp, i, tbaddr, 8 ,-0x800)
		}
		for i := 0x2000; i < 0x27ff; i += tbaddr {

			writile(stamp, i, tbaddr, 8 ,0)
		}
 // locked out
 }
 // the single 8x8 tile set end
 // this op is a beast to run and it takes a while

		pnum++
	}
 // Six - end of tile dumper
}
//
// tile dumper ending


// Six - maze dumper
// - can be used to dump mazes in vars formats for clones and such
// this bit is focused on sanctuary png encoded mazes, where the hex color 0xFFFFFF provides map data
// and a possible SVRLOAD substitute encoded storage format in javascript
if opts.Verbose || opts.Se {
	i := 0
	mz := mazenum + 1
	wimg := blankimage(33, 33)
	if opts.Se {
 // paste in sanctuary converter
		if mz > 116 { mz = mz - 2 }	else {	// sanctuary does not have 115 as demo or 116 as score table
 // just a convention for se - these are already numbered in the system and not at 115 and 116
			if mz == 115 { mz = 0 }
			if mz == 116 { mz = 127 }		// 1 past se end
		}
		fmt.Printf("	SVRLOAD[1] = [ ];\n	SVRLOAD[1][1] = \"levels/glevel%d.png\"\n	SVRLOAD[1][2] = \"Level %d\";\n	SVRLOAD[1][3] = [ ];\n	SVRLOAD[1][4] =\"1089\";\n", mz, mz)
	}

	for y := 0; y <= lasty; y++ {
		for x := 0; x <= lastx; x++ {

			if opts.Verbose { fmt.Printf(" %02d", maze.data[xy{x, y}]) }
			if opts.Se {
 //				fmt.Printf("	SVRLOAD[1][3][%d] = \"0x%x\";\n", i, sanct_vrt[maze.data[xy{x, y}]])
				i++
				var vr int
				var hexc string
				if G1 {
					vr = sanct_vrt[maze.data[xy{x, y}]]
					hexc = fmt.Sprintf("#%06x",sanct_vrt[maze.data[xy{x, y}]])
				} else {
					vr = sanct_vrt2[maze.data[xy{x, y}]]
					hexc = fmt.Sprintf("#%06x",sanct_vrt2[maze.data[xy{x, y}]])	// CHANGE: sanct_vrt2
				}
				if vr < 0x1000 { fmt.Printf("// error used - %X\n",vr) }
				mcol, err := ParseHexColor(hexc)
				if err == nil {
					wimg.Set(x,y,mcol)
				}
			}
		}
		if opts.Verbose { fmt.Printf("\n") }
	}
	if opts.Se {
		wnam := fmt.Sprintf("selvls/glevel%d.png",mz)
		if G2 { wnam = fmt.Sprintf("selvls/g2level%d.png",mz) }
		wrfile, err := os.Create(wnam)
		if err == nil {
			png.Encode(wrfile,wimg)
		}
		wrfile.Close()
	}
 // Six end maze dumper
}
// maze dumper ending

	if xspc == 32 {		// write wrap dir arrows
		if maze.flags&LFLAG4_WRAP_H > 0 {
			l := itemGetStamp("arrowleft")
			r := itemGetStamp("arrowright")
			for i := 2; i <= 32; i++ {
				writestamptoimage(img, l, 0, i*16)
				writestamptoimage(img, r, 32*16+16, i*16)
			}
		}

		if maze.flags&LFLAG4_WRAP_V > 0 {
			u := itemGetStamp("arrowup")
			d := itemGetStamp("arrowdown")
			for i := 1; i < 32; i++ {
				writestamptoimage(img, u, i*16+16, 0)
				writestamptoimage(img, d, i*16+16, 32*16+16)
			}
		}
	}
	savetopng(opts.Output, img)
// for user select
	return img
}

// check to see if there's walls adjacent left, left/up, and up
// FIXME: This might should be a different set of directions
// horizontal wall += 4
// diagonal wall += 8
// vertical wall += 16

func whatis(maze *Maze, x int, y int) int {
	return maze.data[xy{x, y}]
}

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

// g2 version - wall shadow, etc
func checkwalladj3(maze *Maze, x int, y int) int {
	adj := 0

	if iswall(whatis(maze, x-1, y)) {
		adj += 4
	}

	if iswall(whatis(maze, x, y+1)) {
		adj += 16
	}

	if iswall(whatis(maze, x-1, y+1)) {
		adj += 8
	}

	return adj
}

// g1 version
func checkwalladj3g1(maze *Maze, x int, y int) int {
	adj := 0

	if iswallg1(whatis(maze, x-1, y)) {
		adj += 4
	}

	if iswallg1(whatis(maze, x, y+1)) {
		adj += 16
	}

	if iswallg1(whatis(maze, x-1, y+1)) {
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
//
// FIXME: Our sense of up/down here is probably confused

func checkwalladj8(maze *Maze, x int, y int) int {
	adj := 0

	if iswall(whatis(maze, x-1, y-1)) {
		adj += 0x01
	}
	if iswall(whatis(maze, x, y-1)) {
		adj += 0x02
	}
	if iswall(whatis(maze, x+1, y-1)) {
		adj += 0x04
	}
	if iswall(whatis(maze, x-1, y)) {
		adj += 0x08
	}
	if iswall(whatis(maze, x+1, y)) {
		adj += 0x010
	}
	if iswall(whatis(maze, x-1, y+1)) {
		adj += 0x20
	}
	if iswall(whatis(maze, x, y+1)) {
		adj += 0x40
	}
	if iswall(whatis(maze, x+1, y+1)) {
		adj += 0x80
	}

	return adj
}

// g1 version -- g2 has more walls
func checkwalladj8g1(maze *Maze, x int, y int) int {
	adj := 0

	if iswallg1(whatis(maze, x-1, y-1)) {
		adj += 0x01
	}
	if iswallg1(whatis(maze, x, y-1)) {
		adj += 0x02
	}
	if iswallg1(whatis(maze, x+1, y-1)) {
		adj += 0x04
	}
	if iswallg1(whatis(maze, x-1, y)) {
		adj += 0x08
	}
	if iswallg1(whatis(maze, x+1, y)) {
		adj += 0x010
	}
	if iswallg1(whatis(maze, x-1, y+1)) {
		adj += 0x20
	}
	if iswallg1(whatis(maze, x, y+1)) {
		adj += 0x40
	}
	if iswallg1(whatis(maze, x+1, y+1)) {
		adj += 0x80
	}

	return adj
}

// Look and see what doors are adjacent to this door
//
// Values to use:
//    up:  0x01    right:  0x02    down:  0x04    left:  0x08

func checkdooradj4(maze *Maze, x int, y int) int {
	adj := 0

	if isdoor(whatis(maze, x, y-1)) {
		adj += 0x01
	}
	if isdoor(whatis(maze, x+1, y)) {
		adj += 0x02
	}
	if isdoor(whatis(maze, x, y+1)) {
		adj += 0x04
	}
	if isdoor(whatis(maze, x-1, y)) {
		adj += 0x08
	}

	return adj
}

// g1 version
func checkdooradj4g1(maze *Maze, x int, y int) int {
	adj := 0

	if isdoorg1(whatis(maze, x, y-1)) {
		adj += 0x01
	}
	if isdoorg1(whatis(maze, x+1, y)) {
		adj += 0x02
	}
	if isdoorg1(whatis(maze, x, y+1)) {
		adj += 0x04
	}
	if isdoorg1(whatis(maze, x-1, y)) {
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
			t := whatis(maze, x+(j*ffLoopDirs[i].x), y+(j*ffLoopDirs[i].y))
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
// image from buffer segment			- stat: display stats if true
// segment of buffer from xb,yb to xs,ys (begin to stop)

func segimage(mdat MazeData, fdat [11]int, xb int, yb int, xs int, ys int, stat bool) *image.NRGBA {

//if opts.Verbose {
fmt.Printf("segimage %dx%d - %dx%d: %t, vp: %d\n ",xb,yb,xs,ys,stat,viewp)

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
					writestamptoimage(img, stamp, (x-xb)*16, (y-yb)*16)
				}
			}
			if whatis(maze, x, y) < 0 {		// dont floor a null area
				coltil(img,0,(x-xb)*16, (y-yb)*16)
			} else {
			if (nothing & NOFLOOR) == 0 {
				writestamptoimage(img, stamp, (x-xb)*16, (y-yb)*16)
			}}
		}
	}} else {
// g1 checks
	for y := yb; y < ys; y++ {
		for x := xb; x < xs; x++ {
			adj := 0
			nwt := NOWALL | NOG1W
			if whatis(maze, x, y) == G1OBJ_WALL_TRAP1 { nwt = NOWALL }
			if whatis(maze, x, y) == G1OBJ_WALL_DESTRUCTABLE { nwt = NOWALL }
			if maze.wallpattern < 11 {
				if (nothing & nwt) == 0 {		// wall shadows here
				adj = checkwalladj3g1(maze, x, y)
				}
			}

			stamp := floorGetStamp(maze.floorpattern, adj+rand.Intn(4), maze.floorcolor)
			if whatis(maze, x, y) < 0 {
				coltil(img,0,(x-xb)*16, (y-yb)*16)
			} else {
			if (nothing & NOFLOOR) == 0 {
				writestamptoimage(img, stamp, (x-xb)*16, (y-yb)*16)
			}}
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
						if ts == 2 { maze.data[xy{x, y}] = 	MAZEOBJ_TREASURE_BAG }
						if ts == 111 { maze.data[xy{x, y}] = MAZEOBJ_HIDDENPOT }
						if ts == 311 { maze.data[xy{x, y}] = MAZEOBJ_HIDDENPOT }
					}
			}}
			if G1 {
				nwt := NOWALL | NOG1W
				switch whatis(maze, x, y) {
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
						if ts == 2 { maze.data[xy{x, y}] = G1OBJ_TREASURE_BAG }
						if ts == 11 { maze.data[xy{x, y}] = MAZEOBJ_HIDDENPOT }
						if ts == 311 { maze.data[xy{x, y}] = MAZEOBJ_HIDDENPOT }
					}
				}}
			if stamp != nil {
				writestamptoimage(img, stamp, (x-xb)*16+stamp.nudgex, (y-yb)*16+stamp.nudgey)
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
// gen type op - letter to draw
			gtopl := ""
			gtopcol := false	// disable gen letter seperate colors
// gen type op - the context to draw
			gtop := gg.NewContext(12, 12)
// gtop font
			if err := gtop.LoadFontFace(".font/VrBd.ttf", 10); err != nil {
				panic(err)
				}

if opts.Verbose { fmt.Printf("%03d ",whatis(maze, x, y)) }
// g2 decodes
			if G2 {

			// We should do better
			switch whatis(maze, x, y) {
// specials are jammed in somewhere in G2 code, we just do this
// specials added after convert to se id'ed them on maze 115, score table block
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
				if opts.Verbose && false { fmt.Printf("G² WARNING: Unhandled obj id 0x%02x\n", whatis(maze, x, y)) }
			}
// set mask flag in array
			if whatis(maze, x, y) > 0 && stamp != nil { g2mask[whatis(maze, x, y)] = stamp.mask }
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
// /fmt.Printf("g1 dec: %x -- ", whatis(maze, x, y))
			switch whatis(maze, x, y) {

			case G1OBJ_TILE_FLOOR:
			// adj := checkwalladj3(maze, x, y) + rand.Intn(4)
			// stamp = floorGetStamp(maze.floorpattern, adj, maze.floorcolor)
// dont think g1 has stun tile
			case G1OBJ_TILE_STUN:
				adj := checkwalladj3g1(maze, x, y) + rand.Intn(4)
				stamp = floorGetStamp(maze.floorpattern, adj, maze.floorcolor)
				stamp.ptype = "stun" // use trap palette (FIXME: consider moving)
				stamp.pnum = 0

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

			default:
				if opts.Verbose && false { fmt.Printf("G¹ WARNING: Unhandled obj id 0x%02x\n", whatis(maze, x, y)) }
			}
// set mask flag in array
			if whatis(maze, x, y) > 0 && stamp != nil { g1mask[whatis(maze, x, y)] = stamp.mask }
		}
// Six: end G1 decode
			if stamp != nil {
				writestamptoimage(img, stamp, (x-xb)*16+stamp.nudgex, (y-yb)*16+stamp.nudgey)
// stats on palette
				if stat {			// on palette screen, show stats for loaded maze
					st := ""
					mel := whatis(maze, x, y)
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
package main

import (
	"image"
	"image/png"
	"math/rand"
	"fmt"
	"os"
	"github.com/fogleman/gg"
)

// For testing
var maze0 = []int{
	0x0,  // secret trick
	0x0,  // flags 1
	0x0,  // flags 2
	0x0,  // flags 3
	0x0,  // flags 4
	0x61, // patterns
	0x1D, // colors
	0x6E, 0x54, 0xE, 0x9, 0x35, 0x80, 0x47, 0xC4, 0x31,
	0x42, 0xC9, 0x35, 0x6F, 0xDE, 0x14, 0xC8, 0x1D,
	0xC1, 0x1D, 0xC4, 0x1E, 0xCB, 0x22, 0xDE, 0x14,
	0xDD, 0xC, 0x81, 0xDF, 0xD9, 0x9, 0x81, 0xC3,
	0x9, 0x81, 0xD0, 0x58, 0xDF, 0xC3, 0x24, 0xC6,
	0x24, 0xC2, 0x59, 0xDF, 0xD0, 0x20, 0x40, 0xC0,
	0x74, 0xC0, 0x11, 0xC0, 0x74, 0x40, 0xC0, 0x20,
	0xC9, 0x7, 0x84, 0xB3, 0xA2, 0xB5, 0xA8, 0xB5,
	0xA1, 0xB2, 0x8, 0x86, 0xC4, 0xA, 0x80, 0xCD,
	0xB, 0x80, 0xDF, 0xC7, 0xBF, 0xDF, 0xDF, 0xDF,
	0xDF, 0xC1, 0x1E, 0xDC, 0xA4, 0xDA, 0x10, 0xC9,
	0x1C, 0xC3, 0x1C, 0xDF, 0x45, 0xC2, 0x32, 0xC0,
	0x1E, 0xC5, 0xB5, 0xA6, 0xD, 0x82, 0xA5, 0xB5,
	0xA8, 0xDD, 0x2E, 0xC5, 0x1F, 0xDF, 0xDF, 0x48,
	0xCA, 0xF, 0xDF, 0xDF, 0xC2, 0xBE, 0x35, 0xDC,
	0x2E, 0x0}

// Okay, so we have a maze. We need to adjust the edges to take care of
// wrap or no wrap.
//
// If we're not wrapping in a direction, we need to duplicate the left and
// top walls, so that the maze will be enclosed.
func copyedges(maze *Maze) {
	for i := 0; i <= 32; i++ {
		if (maze.flags & LFLAG4_WRAP_H) == 0 {
			maze.data[xy{32, i}] = maze.data[xy{0, i}]
		}
	}

	for i := 0; i <= 32; i++ {
		if (maze.flags & LFLAG4_WRAP_V) == 0 {
			maze.data[xy{i, 32}] = maze.data[xy{i, 0}]
		}
	}
}

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

func genpfimage(maze *Maze, mazenum int) {
	extrax, extray := 0, 0
	if (maze.flags & LFLAG4_WRAP_H) == 0 {
		extrax = 16
	}
	if (maze.flags & LFLAG4_WRAP_V) == 0 {
		extray = 16
	}

	// 8 pixels * 2 tiles * 32 stamps, plus extra space on edges
	img := blankimage(8*2*32+32+extrax, 8*2*32+32+extray)

	// Map out where forcefield floor tiles are, so we can lay those down first
	ffmap := ffMakeMap(maze)

	// mazes will always be the same size, so just use constants
	// maze := mazeDecompress(mazedata)
	copyedges(maze)
	paletteMakeSpecial(maze.floorpattern, maze.floorcolor, maze.wallpattern, maze.wallcolor)

	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			adj := 0
			if maze.wallpattern < 11 {
				adj = checkwalladj3(maze, x, y)
			}

			stamp := floorGetStamp(maze.floorpattern, adj+rand.Intn(4), maze.floorcolor)
			if ffmap[xy{x, y}] == true {
				stamp.ptype = "forcefield"
				stamp.pnum = 0
			}
			writestamptoimage(img, stamp, x*16+16, y*16+16)
		}
	}

	lastx := 32
	if maze.flags&LFLAG4_WRAP_H > 0 {
		lastx = 31
	}

	lasty := 32
	if maze.flags&LFLAG4_WRAP_V > 0 {
		lasty = 31
	}

	for y := 0; y <= lasty; y++ {
		for x := 0; x <= lastx; x++ {
			var stamp *Stamp
			var dots int // dot count

			// We should do better
			switch whatis(maze, x, y) {
			case MAZEOBJ_TILE_FLOOR:
			// adj := checkwalladj3(maze, x, y) + rand.Intn(4)
			// stamp = floorGetStamp(maze.floorpattern, adj, maze.floorcolor)
			case MAZEOBJ_TILE_STUN:
				adj := checkwalladj3(maze, x, y) + rand.Intn(4)
				stamp = floorGetStamp(maze.floorpattern, adj, maze.floorcolor)
				stamp.ptype = "stun" // use trap palette (FIXME: consider moving)
				stamp.pnum = 0

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
				stamp = floorGetStamp(maze.floorpattern, adj, maze.floorcolor)
				stamp.ptype = "trap" // use trap palette (FIXME: consider moving)
				stamp.pnum = 0
			case MAZEOBJ_WALL_DESTRUCTABLE:
				adj := checkwalladj8(maze, x, y)
				stamp = wallGetDestructableStamp(maze.wallpattern, adj, maze.wallcolor)
			case MAZEOBJ_WALL_SECRET:
				adj := checkwalladj8(maze, x, y)
				stamp = wallGetStamp(maze.wallpattern, adj, maze.wallcolor)
				stamp.ptype = "secret"
				stamp.pnum = 0
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
				stamp = wallGetStamp(maze.wallpattern, adj, maze.wallcolor)
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

			case MAZEOBJ_GEN_GHOST1:
				stamp = itemGetStamp("ghostgen1")
			case MAZEOBJ_GEN_GHOST2:
				stamp = itemGetStamp("ghostgen2")
			case MAZEOBJ_GEN_GHOST3:
				stamp = itemGetStamp("ghostgen3")

			case MAZEOBJ_GEN_GRUNT1:
				fallthrough
			case MAZEOBJ_GEN_DEMON1:
				fallthrough
			case MAZEOBJ_GEN_LOBBER1:
				fallthrough
			case MAZEOBJ_GEN_SORC1:
				fallthrough
			case MAZEOBJ_GEN_AUX_GRUNT1:
				stamp = itemGetStamp("generator1")

			case MAZEOBJ_GEN_GRUNT2:
				fallthrough
			case MAZEOBJ_GEN_DEMON2:
				fallthrough
			case MAZEOBJ_GEN_LOBBER2:
				fallthrough
			case MAZEOBJ_GEN_SORC2:
				fallthrough
			case MAZEOBJ_GEN_AUX_GRUNT2:
				stamp = itemGetStamp("generator2")

			case MAZEOBJ_GEN_GRUNT3:
				fallthrough
			case MAZEOBJ_GEN_DEMON3:
				fallthrough
			case MAZEOBJ_GEN_LOBBER3:
				fallthrough
			case MAZEOBJ_GEN_SORC3:
				fallthrough
			case MAZEOBJ_GEN_AUX_GRUNT3:
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
				stamp = ffGetStamp(adj)
			case MAZEOBJ_TRANSPORTER:
				stamp = itemGetStamp("tport")
			default:
				// fmt.Printf("WARNING: Unhandled obj id 0x%02x\n", whatis(maze, x, y))
			}

// g1 decodes
			if mazenum > 0x37FFF {
// /fmt.Printf("g1 dec: %x -- ", whatis(maze, x, y))
			switch whatis(maze, x, y) {

			case G1OBJ_TILE_FLOOR:
			// adj := checkwalladj3(maze, x, y) + rand.Intn(4)
			// stamp = floorGetStamp(maze.floorpattern, adj, maze.floorcolor)
			case G1OBJ_TILE_STUN:
				adj := checkwalladj3(maze, x, y) + rand.Intn(4)
				stamp = floorGetStamp(maze.floorpattern, adj, maze.floorcolor)
				stamp.ptype = "stun" // use trap palette (FIXME: consider moving)
				stamp.pnum = 0

				// Tried to simplify these a bit with a goto, but golang didn't
				// like it ('jump into block'). I should figure out why.
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
				stamp = floorGetStamp(maze.floorpattern, adj, maze.floorcolor)
				stamp.ptype = "trap" // use trap palette (FIXME: consider moving)
				stamp.pnum = 0
			case G1OBJ_WALL_DESTRUCTABLE:
				adj := checkwalladj8(maze, x, y)
				stamp = wallGetDestructableStamp(maze.wallpattern, adj, maze.wallcolor)

			case G1OBJ_WALL_TRAP1:
				dots = 1
				fallthrough
			case G1OBJ_WALL_REGULAR:
				adj := checkwalladj8(maze, x, y)
				stamp = wallGetStamp(maze.wallpattern, adj, maze.wallcolor)
			case G1OBJ_KEY:
				stamp = itemGetStamp("key")

			case G1OBJ_DOOR_HORIZ:
				adj := checkdooradj4(maze, x, y)
				stamp = doorGetStamp(DOOR_HORIZ, adj)
			case G1OBJ_DOOR_VERT:
				adj := checkdooradj4(maze, x, y)
				stamp = doorGetStamp(DOOR_VERT, adj)

			case G1OBJ_PLAYERSTART:
				stamp = itemGetStamp("plus")
			case G1OBJ_EXIT:
				stamp = itemGetStamp("exit")
			case G1OBJ_EXIT4:
				stamp = itemGetStamp("exit4")
			case G1OBJ_EXIT8:
				stamp = itemGetStamp("exit8")

			case G1OBJ_MONST_GHOST1:
				fallthrough
			case G1OBJ_MONST_GHOST2:
				fallthrough
			case G1OBJ_MONST_GHOST3:
				stamp = itemGetStamp("ghost")
			case G1OBJ_MONST_GRUNT1:
				fallthrough
			case G1OBJ_MONST_GRUNT2:
				fallthrough
			case G1OBJ_MONST_GRUNT3:
				stamp = itemGetStamp("grunt")
			case G1OBJ_MONST_DEMON1:
				fallthrough
			case G1OBJ_MONST_DEMON2:
				fallthrough
			case G1OBJ_MONST_DEMON3:
				stamp = itemGetStamp("demon")
			case G1OBJ_MONST_LOBBER1:
				fallthrough
			case G1OBJ_MONST_LOBBER2:
				fallthrough
			case G1OBJ_MONST_LOBBER3:
				stamp = itemGetStamp("lobber")
			case G1OBJ_MONST_SORC1:
				fallthrough
			case G1OBJ_MONST_SORC2:
				fallthrough
			case G1OBJ_MONST_SORC3:
				stamp = itemGetStamp("sorcerer")
			case G1OBJ_MONST_DEATH:
				stamp = itemGetStamp("death")

			case G1OBJ_GEN_GHOST1:
				stamp = itemGetStamp("ghostgen1")
			case G1OBJ_GEN_GHOST2:
				stamp = itemGetStamp("ghostgen2")
			case G1OBJ_GEN_GHOST3:
				stamp = itemGetStamp("ghostgen3")

			case G1OBJ_GEN_GRUNT1:
				fallthrough
			case G1OBJ_GEN_DEMON1:
				fallthrough
			case G1OBJ_GEN_LOBBER1:
				fallthrough
			case G1OBJ_GEN_SORC1:
				stamp = itemGetStamp("generator1")

			case G1OBJ_GEN_GRUNT2:
				fallthrough
			case G1OBJ_GEN_DEMON2:
				fallthrough
			case G1OBJ_GEN_LOBBER2:
				fallthrough
			case G1OBJ_GEN_SORC2:
				stamp = itemGetStamp("generator2")

			case G1OBJ_GEN_GRUNT3:
				fallthrough
			case G1OBJ_GEN_DEMON3:
				fallthrough
			case G1OBJ_GEN_LOBBER3:
				fallthrough
			case G1OBJ_GEN_SORC3:
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

			case G1OBJ_TRANSPORTER:
				stamp = itemGetStamp("tport")
			default:
			}}
// Six: end G1 decode
			if stamp != nil {
				writestamptoimage(img, stamp, x*16+16+stamp.nudgex, y*16+16+stamp.nudgey)
			}

			if dots != 0 {
				renderdots(img, x*16+16, y*16+16, dots)
			}
		}
	}

/// Six - individual tile dumper

/*******************************************
TTT           ___  DDD                __
 T   i  L    E     D  d  u  u  m\/m  P  p
 T   i  L    E--   D  d  u  u  m  m  p__p
 T   i  LLL  E___  DDD    uu   m  m  p

 *******************************************/
// lock out with this - set to true to run dump
// run decode for maze116
// written to .p[0-11]/tl_%05d_%04X.png
// where %d and %X are tile start addr

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

		if stamp.pnum == 12 {
			stamp = nil
		}
	}

// gauntlet 2 add ins not handled by 0000 - 1FFF dump
	pnum := 0

	for pnum < 12 && mazenum == 116 {

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
// this is written to .8x8/c[0-11]/i%05d_%04X.png
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
		pnum++
	}
// Six - end of tile dumper
}
// Six - maze dumper
if true {
	fmt.Printf("maze: %d\n", mazenum)
	for y := 0; y <= lasty; y++ {
		for x := 0; x <= lastx; x++ {

			fmt.Printf(" %02d", maze.data[xy{x, y}])
		}
		fmt.Printf("\n")
	}
// Six end maze dumper
}
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
	savetopng(opts.Output, img)
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

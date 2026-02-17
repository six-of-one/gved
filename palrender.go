package main

import (
	"image"
	"math/rand"
	"fmt"
	"github.com/fogleman/gg"
	"encoding/binary"
	"golang.org/x/image/draw"
)

// palete and pb are getting sep render, since they are much simpler


func palimage(mdat MazeData, xdat Xdat, fdat [14]int, xb int, yb int, xs int, ys int, stat bool) *image.NRGBA {

fmt.Printf("palimage %dx%d - %dx%d: %t, vp: %d\n",xb,yb,xs,ys,stat,viewp)

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

	// unpin issue - -vals flummox canvas writes
	xba, yba := vpc_adj(xb, yb)

	img := blankimage(16*(xs-xb), 16*(ys-yb))		// not main maze

	svanim = false
	florbas(img, maze, xdat, opts.DimX+1, opts.DimY+1,false)
	walbas(img, maze, xdat, opts.DimX+1, opts.DimY+1,false)
	svanim = true

	for y := yb; y <= ys; y++ {
		for x := xb; x <= xs; x++ {
			var stamp *Stamp
			var dots int // dot count

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

				//	if !G2 {
// gen type op - put a letter on up left corner of every gen to indicate monsters
//		brw G - grunts
//		red D - demons
//		yel L - lobbers
//		pur S - sorceror
			gtop.Clear()
			gtopl = ""// make sure G² code (if it runs with G¹) doesnt set extra dots on non walls
			dots = 0

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

		default:
			if opts.Verbose && false { fmt.Printf("G¹ WARNING: Palete/pb obj id 0x%02x\n", sb) }
		}
// set mask flag in array
		nugetx, nugety := -4, -4
		if sb > 0 {

		if stamp != nil {
		if stamp.mask & nothing == 0 {
			g1mask[sb] = stamp.mask
// note G¹ here, opposite of other writes using gt - here gt preserves true G¹ state due to complex tile rom extract and pallet select
			writestamptoimage(G1,img, stamp, vcx*16+stamp.nudgex, vcy*16+stamp.nudgey)
			nugetx, nugety = stamp.nudgex, stamp.nudgey
		}} else {
//fmt.Printf("star ld %d, %v\n",sb)
		if arstamp[sb].mask & nothing == 0 {
			if arstamp[sb].pnum > -1 || arstamp[sb].pnum == -7 {
				gtopl = arstamp[sb].gtopl
				offset := image.Pt(vcx*16+arstamp[sb].nudgex, vcy*16+arstamp[sb].nudgey)
				draw.Draw(img, arstamp[sb].mimg.Bounds().Add(offset), arstamp[sb].mimg, image.ZP, draw.Over)
				if arstamp[sb].pnum != -7 { nugetx, nugety = arstamp[sb].nudgex, arstamp[sb].nudgey }
			}
		}}}

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
				offset := image.Pt(vcx*16+nugetx-4, vcy*16+nugety-4)
				draw.Draw(img, gtopim.Bounds().Add(offset), gtopim, image.ZP, draw.Over)
			}

			if dots != 0 && nothing & NOWALL == 0 {
				renderdots(img, (x-xb)*16, (y-yb)*16, dots)
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
	g1mask[G1OBJ_WALL_REGULAR] = 2048
	g1mask[G1OBJ_WALL_DESTRUCTABLE] = 1024
	g1mask[G1OBJ_WALL_TRAP1] = 1024
	g1mask[G1OBJ_TILE_TRAP1] = 64

	return img
}
package main

// #include <stdlib.h>
// #include <Tilengine.h>

import "C"

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"os"
//	"time"

	"fyne.io/fyne/v2"
)

func mazeMetaPrint(maze *Maze) {
	fmt.Printf("  Encoded length: %3d bytes\n", maze.encodedbytes)
	fmt.Printf("  Wall pattern: %02d, Wall color: %02d     Floor pattern: %02d, Floor color: %02d\n",
		maze.wallpattern, maze.wallcolor, maze.floorpattern, maze.floorcolor)
	fmt.Printf("  Flags: ")
	for k, v := range mazeFlagStrings {
		if (maze.flags & k) != 0 {
			fmt.Printf("%s ", v)
		}
	}
	if G2 {
		fmt.Printf("\n  Random food adds: %d\n", (maze.flags&LFLAG3_RANDOMFOOD_MASK)>>8)
		fmt.Printf("  Secret trick: %2d - %s\n", maze.secret, mazeSecretStrings[maze.secret])
	}
}

var reMazeNum = regexp.MustCompile(`^maze(\d+)`)
var reMazeMeta = regexp.MustCompile(`^meta$`)
var maxmaze int

func domaze(arg string) {
	split := strings.Split(arg, "-")

	var mazeNum = -1
	var mazeMeta = 0
	maxmaze = 116

// g1 has more mazes, but treasure rooms can only spec from address, for now
	if G1 { maxmaze = 126 }

	for _, ss := range split {
		switch {
		case reMazeNum.MatchString(ss):
			m, _ := strconv.ParseInt(reMazeNum.FindStringSubmatch(ss)[1], 10, 0)
			mazeNum = int(m) - 1

		case reMazeMeta.MatchString(ss):
			mazeMeta = 1
		}
	}
	if mazeNum < 0 || mazeNum > maxmaze {
		panic("Invalid maze number / address specified.")
	}

	opts.mnum = mazeNum
	fmt.Printf("Maze number: %d", mazeNum + 1)
	if Aov > 0 {
		fmt.Printf(", address: 0x%X ", Aov)
	}
	fmt.Printf("\n")

// set 1 override to -1 to set in decoder
	Ovwallpat = -1
	nothing = 0

	maze := mazeDecompress(slapsticReadMaze(mazeNum), false)

	if opts.Verbose || mazeMeta > 0 {
		mazeMetaPrint(maze)
		if mazeMeta > 0 { os.Exit(0) }
	}

// interactive viewer not selected - gen maze, output png & exit
	if !opts.Intr {
		genpfimage(maze, mazeNum)
		os.Exit(0)
	}

// in interactive, start the window
	aw_init()

	Ovimg := genpfimage(maze, mazeNum)
	upwin(Ovimg)
	w.Resize(fyne.NewSize(1024, 1024))
// only run Show once, here - a second time relocates the win to 0,0
// yes... even though fyne can NOT reposition windows, must be a bug
	w.ShowAndRun()

}

// loop called by typedRune in kontrol.go to re-issue maze after viewer parm changes

func mazeloop(maze *Maze) {
// to transform maze, array copy
	xform := make(map[xy]int)
// manual mirror, flip
	if opts.MH || opts.MV || opts.MRP || opts.MRM {
		lastx := 32
		if maze.flags&LFLAG4_WRAP_H > 0 {
			lastx = 31
		}

		lasty := 32
		if maze.flags&LFLAG4_WRAP_V > 0 {
			lasty = 31
		}
// note it
/*
	for y := 0; y <= lasty; y++ {
		for x := 0; x <= lastx; x++ {

			fmt.Printf(" %02d", maze.data[xy{x, y}])
		}
		fmt.Printf("\n")
	}
*/
// transform
// rotate +90 degrees
		if opts.MRP {
			for ty := 1; ty <= lasty; ty++ {
			for tx := 1; tx <= lastx; tx++ {
				xform[xy{lastx - tx, ty}] = maze.data[xy{ty, tx}]
// g1 - must transform all dors on a rotat since they have horiz & vert dependent
				if xform[xy{lastx - tx, ty}] == G1OBJ_DOOR_HORIZ { xform[xy{lastx - tx, ty}] = G1OBJ_DOOR_VERT } else {
				if xform[xy{lastx - tx, ty}] == G1OBJ_DOOR_VERT { xform[xy{lastx - tx, ty}] = G1OBJ_DOOR_HORIZ } }
// g2
				if xform[xy{lastx - tx, ty}] == MAZEOBJ_DOOR_HORIZ { xform[xy{lastx - tx, ty}] = MAZEOBJ_DOOR_VERT } else {
				if xform[xy{lastx - tx, ty}] == MAZEOBJ_DOOR_VERT { xform[xy{lastx - tx, ty}] = MAZEOBJ_DOOR_HORIZ } }
			}}
		} else {
		if opts.MRM {
			for ty := 1; ty <= lasty; ty++ {
			for tx := 1; tx <= lastx; tx++ {
				xform[xy{tx, lasty - ty}] = maze.data[xy{ty, tx}]
// g1
				if xform[xy{tx, lasty - ty}] == G1OBJ_DOOR_HORIZ { xform[xy{tx, lasty - ty}] = G1OBJ_DOOR_VERT } else {
				if xform[xy{tx, lasty - ty}] == G1OBJ_DOOR_VERT { xform[xy{tx, lasty - ty}] = G1OBJ_DOOR_HORIZ } }
// g2
				if xform[xy{tx, lasty - ty}] == MAZEOBJ_DOOR_HORIZ { xform[xy{tx, lasty - ty}] = MAZEOBJ_DOOR_VERT } else {
				if xform[xy{tx, lasty - ty}] == MAZEOBJ_DOOR_VERT { xform[xy{tx, lasty - ty}] = MAZEOBJ_DOOR_HORIZ } }
			}}
		}
		}

// have to copy back if doing with any mirror cmd
		if opts.MRP || opts.MRM {
		if opts.MH || opts.MV {
		for y := 1; y <= lasty; y++ {
			for x := 1; x <= lastx; x++ { maze.data[xy{x, y}] = xform[xy{x, y}] }
		}}}

// mirror x
		if opts.MH {
			for ty := 1; ty <= lasty; ty++ {
			for tx := 1; tx <= lastx; tx++ {
				xform[xy{lastx - tx, ty}] = maze.data[xy{tx, ty}]
			}}
		}
// have to copy back if doing both together
		if opts.MH && opts.MV {
		for y := 1; y <= lasty; y++ {
			for x := 1; x <= lastx; x++ { maze.data[xy{x, y}] = xform[xy{x, y}] }
		}}

// mirror y: flip
		if opts.MV {
			for ty := 1; ty <= lasty; ty++ {
			for tx := 1; tx <= lastx; tx++ {
				xform[xy{tx, lasty - ty}] = maze.data[xy{tx, ty}]
			}}
		}

// copy back
		for y := 1; y <= lasty; y++ {
			for x := 1; x <= lastx; x++ { maze.data[xy{x, y}] = xform[xy{x, y}] }
		}
	}
}

// verify manually selected address, or page keys 'z' 'x' move thru address list

// ad - address to test
// dr - direction to move in array
// a. select next addr in loop (dr = 1, -1)
// b. verify an entered addr (dr = 0)


func addrver(ad int, dr int) int {

	rt := 0					// return
	rn := 0					// return nearest
	i := 0
	if G1 {

	as := g1validaddr[0]	// addr search set
	for as > 0 {
		if ad == as {
			if (i + dr) == -1 { rt = g1validaddr[len(g1validaddr)-1]} else {		// looped array start, get end
				rt = g1validaddr[i + dr]
			}
		} else {
			if dr == 0 {	// for addr entry verify, pick nearest lower if not found
				if ad > g1validaddr[i] { rn = g1validaddr[i] }
			}
		}
		i++
		as = g1validaddr[i]
	}
	if rt == -1 { rt = g1validaddr[0] }		// looped end of array, get 1st
	} else {

	as := g2validaddr[0]
	for as > 0 {
		if ad == as {
			if (i + dr) == -1 { rt =  g2validaddr[len(g2validaddr)-1]} else {
				rt = g2validaddr[i + dr]
			}
		} else {
			if dr == 0 {
				if ad > g2validaddr[i] { rn = g2validaddr[i] }
			}
		}
		i++
		as = g2validaddr[i]
	}
	if rt == -1 { rt = g2validaddr[0] }
	}
	if rt == 0 && dr == 0 && rn > 0 { rt = rn }
	return rt
}
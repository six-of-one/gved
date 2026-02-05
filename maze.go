package main

// #include <stdlib.h>
// #include <Tilengine.h>

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"os"
)

// six: enhance to add flags & such to stats output panel from palette

func mazeMetaPrint(maze *Maze, stp bool) string {

	fmtsp := ""
	if !stp { fmt.Printf("  Encoded length: %3d bytes\n", maze.encodedbytes) }
	if sdb < 0 {
		fmtsp = fmt.Sprintf(" Encoded length: %3d bytes\n", maze.encodedbytes)
	}
	if !stp { fmt.Printf("  Wall pattern: %02d, Wall color: %02d     Floor pattern: %02d, Floor color: %02d\n",
		maze.wallpattern, maze.wallcolor, maze.floorpattern, maze.floorcolor) }
	fmtsp += fmt.Sprintf(" Wall pattern: %02d, Wall color: %02d\n Floor pattern: %02d, Floor color: %02d\n",
		maze.wallpattern, maze.wallcolor, maze.floorpattern, maze.floorcolor)
	if !stp { fmt.Printf("  Flags: ") }
	byt := ""
	if maze.flags > 0 { byt = "   byte breakdown  (8 bits)" }
	fmtsp += fmt.Sprintf("    Flags:%s\n%032b\n",byt,maze.flags)
	if eflg[1] > 0 { fmtsp += fmt.Sprintf("%08b_______________________________\n",eflg[1]) }
	if eflg[2] > 0 { fmtsp += fmt.Sprintf("__________%08b_____________________\n",eflg[2]) }
	if eflg[3] > 0 { fmtsp += fmt.Sprintf("_____________________%08b__________\n",eflg[3]) }
	if eflg[4] > 0 { fmtsp += fmt.Sprintf("_______________________________%08b\n",eflg[4]) }
	if maze.flags > 0 { fmtsp += fmt.Sprintf("    val:{ %v }\n",maze.flags) }
	g1flg := false
	g1mask := 0xB3			// see constants.go, g1 only has these flags type elements, yet I dont think they are controlled by flags
	for k, v := range mazeFlagStrings {
		if (maze.flags & k) != 0 {
			if G1 {
				if k & g1mask != 0 { if !stp { fmt.Printf("%s ", v) }; g1flg = true; fmtsp += fmt.Sprintf("- %s\n", v) }
			} else {
				if !stp { fmt.Printf("%s ", v) }
				fmtsp += fmt.Sprintf("- %s\n", v)
			}
		}
	}
	if G2 {
		if !stp { fmt.Printf("\n  Random food adds: %d\n", (maze.flags&LFLAG3_RANDOMFOOD_MASK)>>8) }
		fmtsp += fmt.Sprintf("\nRandom food adds: %d\n", (maze.flags&LFLAG3_RANDOMFOOD_MASK)>>8)
		if !stp { fmt.Printf("  Secret trick: %2d - %s\n", maze.secret, mazeSecretStrings[maze.secret]) }
		fmtsp += fmt.Sprintf("Secret trick: %2d - %s\n", maze.secret, mazeSecretStrings[maze.secret])
	} else {
		if g1flg {	if !stp { fmt.Printf("\n  ╚══> while gauntlet has these elements, these flags probably do not operate") }
					fmtsp += fmt.Sprintf("\n  ╚══> while gauntlet has these elements,\n       these flags probably do not operate") }
		if !stp { fmt.Printf("\n") }
		fmtsp += fmt.Sprintf("\n")
	}
	return fmtsp
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
		if opts.Intr { mazeNum = 0 } else {			// interactive defauklt to gauntlet maze 1 if no spec
			panic("Invalid maze number / address specified.")
		}
	}

	opts.mnum = mazeNum
	fmt.Printf("Maze number: %d", mazeNum + 1)
	if Aov > 0 {
		fmt.Printf(", address: 0x%X ", Aov)
	}
	fmt.Printf("\n")

// set 1 override to -1 to set in decoder
	Ovwallpat = -1
	nothing = opts.Mask & 0xfff;

	init_buf()	// need buffers, one gets loaded
	edmaze = mazeDecompress(slapsticReadMaze(mazeNum), false)

	if opts.Verbose || mazeMeta > 0 {
		mazeMetaPrint(edmaze, false)
		if mazeMeta > 0 { os.Exit(0) }
	}

// interactive viewer not selected - gen maze, output png & exit
	if !opts.Intr {
		genpfimage(edmaze, mazeNum)
		os.Exit(0)
	}

// in interactive, start the window
	aw_init()

	Ovimg := genpfimage(edmaze, mazeNum)
	upwin(Ovimg, 0)

// call handle window resize lock
	go func() {
		wizecon()
	}()

// only run Show once, here - show() a second time relocates the win to 0,0
// yes... even though fyne can NOT reposition windows, must be a bug
	w.ShowAndRun()

}

// loop called by typedRune in kontrol.go to re-issue maze after viewer parm changes

func mazeloop(maze *Maze) {
// to transform maze, array copy
	xform := make(map[xy]int)
// manual mirror, flip
	if opts.MH || opts.MV || opts.MRP || opts.MRM {

		sx := 1
		lastx := 32
		if maze.flags&LFLAG4_WRAP_H > 0 {
			sx = 0
			lastx = 31
		}

		sy := 1
		lasty := 32
		if maze.flags&LFLAG4_WRAP_V > 0 {
			sy = 0		// otherwise it wont MV correct
			lasty = 31
		}
if opts.Verbose {
	fmt.Printf("mloop wraps -- hw: %d vw: %d\n", maze.flags&LFLAG4_WRAP_H,maze.flags&LFLAG4_WRAP_V)
	fmt.Printf("mazeloop fx: %d lx %d fy %d ly %d\n", sx,lastx,sy,lasty)
}

// note it
/*		fmt.Printf("init\n")
	for y := 0; y <= lasty; y++ {
		for x := 0; x <= lastx; x++ {

			fmt.Printf(" %02d", maze.data[xy{x, y}])
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
			for ty := sy; ty <= lasty; ty++ {
			for tx := sx; tx <= lastx; tx++ {
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
// TEMP maze dmp
/*		fmt.Printf("rots\n")
	for y := 0; y <= lasty; y++ {
		for x := 0; x <= lastx; x++ {

			fmt.Printf(" %02d",xform[xy{x, y}])
		}
		fmt.Printf("\n")
	}
		fmt.Printf("\n") */
// REM TEMP

// have to copy back if doing rot with any mirror cmd
		if opts.MRP || opts.MRM {
		if opts.MH || opts.MV {
		for y := sy; y <= lasty; y++ {
			for x := sx; x <= lastx; x++ { maze.data[xy{x, y}] = xform[xy{x, y}] }
		}}}

// mirror x
		if opts.MH {
			for ty := sy; ty <= lasty; ty++ {
			for tx := sx; tx <= lastx; tx++ {
				xform[xy{lastx - tx, ty}] = maze.data[xy{tx, ty}]
			}}
		}
// have to copy back if doing both together
		if opts.MH && opts.MV {
		for y := sy; y <= lasty; y++ {
			for x := sx; x <= lastx; x++ { maze.data[xy{x, y}] = xform[xy{x, y}] }
		}}

// mirror y: flip
		if opts.MV {
			for ty := sy; ty <= lasty; ty++ {
			for tx := sx; tx <= lastx; tx++ {
				xform[xy{tx, lasty - ty}] = maze.data[xy{tx, ty}]
			}}
			if maze.flags&LFLAG4_WRAP_V > 0 {	// fix wall not allowed being at bottom for arcade gauntlet
				for ty := lasty - 1; ty >= sy ; ty-- {
				for tx := sx; tx <= lastx; tx++ {
					xform[xy{tx, ty + 1}] = xform[xy{tx, ty}]
				}}
				for tx := sx; tx <= lastx; tx++ { xform[xy{tx, 0}] = G1OBJ_WALL_REGULAR }
			}
		}

// copy back
		for y := sy; y <= lasty; y++ {
			for x := sx; x <= lastx; x++ { maze.data[xy{x, y}] = xform[xy{x, y}] }
		}
// TEMP maze dmp
/*		fmt.Printf("dun\n")
	for y := 0; y <= lasty; y++ {
		for x := 0; x <= lastx; x++ {

			fmt.Printf(" %02d", maze.data[xy{x, y}])
		}
		fmt.Printf("\n")
	}
		fmt.Printf("\n") */
// REM TEMP

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
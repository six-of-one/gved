package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
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
	fmt.Printf("\n  Random food adds: %d\n", (maze.flags&LFLAG3_RANDOMFOOD_MASK)>>8)
	fmt.Printf("  Secret trick: %2d - %s\n", maze.secret, mazeSecretStrings[maze.secret])
}

var reMazeNum = regexp.MustCompile(`^maze(\d+)`)
var reMazeMeta = regexp.MustCompile(`^meta$`)

func domaze(arg string) {
	split := strings.Split(arg, "-")

	var mazeNum = -1
	var mazeMeta = 0

	for _, ss := range split {
		switch {
		case reMazeNum.MatchString(ss):
			m, _ := strconv.ParseInt(reMazeNum.FindStringSubmatch(ss)[1], 10, 0)
			mazeNum = int(m)
// Six: g2 maze num 1 - 117 or g1 start address x38000 - x3FFFF
			if mazeNum < 0 || mazeNum > 117 && mazeNum < 229376 || mazeNum > 262144 {
				panic("Invalid maze number / address specified.")
			}

		case reMazeMeta.MatchString(ss):
			mazeMeta = 1
		}
	}

	if mazeNum < 118 {
		fmt.Printf("Gauntlet II\n")
		fmt.Printf("Maze number: %d\n", mazeNum)
	} else {
		fmt.Printf("Gauntlet\n")
		G1 = mazeNum		// G1 mode active, testing
		fmt.Printf("Maze address: 0x%X -- %d\n", mazeNum, mazeNum) }

		G2 = opts.AddrG2

	maze := mazeDecompress(slapsticReadMaze(mazeNum), mazeMeta > 0, mazeNum)
	xform := make(map[xy]int)

// manual mirror, flip
	if opts.MH || opts.MV {
		lastx := 32
		if maze.flags&LFLAG4_WRAP_H > 0 {
			lastx = 31
		}

		lasty := 32
		if maze.flags&LFLAG4_WRAP_V > 0 {
			lasty = 31
		}
// note it
	for y := 0; y <= lasty; y++ {
		for x := 0; x <= lastx; x++ {

			fmt.Printf(" %02d", maze.data[xy{x, y}])
		}
		fmt.Printf("\n")
	}
// transform
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

	if opts.Verbose || mazeMeta > 0 {
		mazeMetaPrint(maze)
	}

	if mazeMeta == 0 {
		genpfimage(maze, mazeNum)
	}
}

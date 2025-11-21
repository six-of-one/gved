package main

import (
	"encoding/binary"
	"fmt"
	"math/rand"
)

// var typeArr = []int{
// 	' ', 'a', 'b', 'c', 'd', 'e' /* 'f', 'g', 'h', 'i', */, 'b', 'b', 'b', 'b',
// 	'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's',
// 	't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C',
// 	'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
// 	'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W',
// 	'X', 'Y', 'Z',
// 	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
// 	'+',
// }

func getbytefortype(t int) int {
	// return typeArr[t]
	return t
}

func index2xy(index int) (x int, y int) {
	if index < 0 {
		fmt.Printf("ERROR: Coordinates requested for index < 0: %d\n", index)
		panic("Coordinates requested for index < 0")
	}

	y = index / 32
	x = index - (y * 32)

	return
}

type xy struct{ x, y int }

type MazeData map[xy]int

type Maze struct {
	data         MazeData
	encodedbytes int
	secret       int
	flags        int
	wallpattern  int
	wallcolor    int
	floorpattern int
	floorcolor   int
}

// Is a maze object a wall?
func iswall(t int) bool {
	if t == MAZEOBJ_WALL_REGULAR || t == MAZEOBJ_WALL_SECRET || t == MAZEOBJ_WALL_DESTRUCTABLE || t == MAZEOBJ_WALL_RANDOM || t == MAZEOBJ_WALL_TRAPCYC1 || t == MAZEOBJ_WALL_TRAPCYC2 || t == MAZEOBJ_WALL_TRAPCYC3 { // || t == MAZEOBJ_FORCEFIELDHUB {
		return true
	} else {
		return false
	}
}

// Is it a floor tile of some type (or similar)
func isspecialfloor(t int) bool {
	if t == MAZEOBJ_TILE_STUN || t == MAZEOBJ_TILE_TRAP1 || t == MAZEOBJ_TILE_TRAP2 || t == MAZEOBJ_TILE_TRAP3 || t == MAZEOBJ_EXIT || t == MAZEOBJ_EXITTO6 || t == MAZEOBJ_TRANSPORTER {
		return true
	} else {
		return false

	}
}

// FIXME: needs to handle vflip and hflip
func expand(maze *Maze, location int, t int, count int) int {
	if t == MAZEOBJ_TILE_FLOOR {
		return (location + count)
	}

	for i := 0; i < count; i++ {
		if iswall(t) {
			x, y := index2xy(location + i)
			maze.data[xy{x, y}] = getbytefortype(t)

		} else if isspecialfloor(t) {
			x, y := index2xy(location + i)
			maze.data[xy{x, y}] = getbytefortype(t)
		} else {
			// things here will need an offset to be completely visible
			/* if t == MAZEOBJ_MONST_DRAGON */

			x, y := index2xy(location + i)
			maze.data[xy{x, y}] = getbytefortype(t)

			if t == MAZEOBJ_MONST_DRAGON {
				i++
			}
		}
	}
	return location + count
}

// FIXME: Needs to handle vflip and hflip
func vexpand(maze *Maze, location int, t int, count int) int {
	if t == MAZEOBJ_TILE_FLOOR {
		return location + 1
	}

	for i := 0; i < count; i++ {
		if iswall(t) {
			x, y := index2xy(location - (i * 32))
			maze.data[xy{x, y}] = getbytefortype(t)
		} else if isspecialfloor(t) {
			x, y := index2xy(location - (i * 32))
			maze.data[xy{x, y}] = getbytefortype(t)
		} else {
			// things here will need a position adjustment to be fully visible
			x, y := index2xy(location - (i * 32))
			maze.data[xy{x, y}] = getbytefortype(t)
		}
	}

	return location + 1
}

// Outoput is maze[y][x]
func mazeDecompress(compressed []int, metaonly bool) *Maze {
	rand.Seed(5)
	//  var m [32][32]int
	var maze = &Maze{}
	maze.data = make(map[xy]int)
	maze.encodedbytes = len(compressed)
	maze.secret = compressed[0] & 0x1f

// Six - maze dumper compresssed
if true {
	fmt.Printf("compresssed: %d\n", maze.encodedbytes)
	y := 0
	for y < maze.encodedbytes {
		for x := 0; x <= 16; x++ {

			if y < maze.encodedbytes {
				fmt.Printf(" %02d", compressed[y])
			}
			y++
		}
		fmt.Printf("\n")
	}
// Six end maze dumper
}
/*
	// This inability to transparently go back and forth between types is
	// obnoxious.
	flagbytes := make([]byte, 4)
	flagbytes[0] = byte(compressed[1])
	flagbytes[1] = byte(compressed[2])
	flagbytes[2] = byte(compressed[3])
	flagbytes[3] = byte(compressed[4])
	maze.flags = int(binary.BigEndian.Uint32(flagbytes))

	maze.wallpattern = compressed[5] & 0x0f
	maze.floorpattern = (compressed[5] & 0xf0) >> 4
	maze.wallcolor = compressed[6] & 0x0f
	maze.floorcolor = (compressed[6] & 0xf0) >> 4

	if metaonly {
		return maze
	}

	htype1 := compressed[7]  // horiz type 1
	htype2 := compressed[8]  // horiz type 2
	vtype1 := compressed[9]  // vert type 1
	vtype2 := compressed[10] // vert type 2

	prev := htype2 // previous value, for 'repeat previous' types
*/
// TEST set - remove
flagbytes := make([]byte, 4)
flagbytes[0] = 0
flagbytes[1] = 0
flagbytes[2] = 0
flagbytes[3] = 0
maze.flags = int(binary.BigEndian.Uint32(flagbytes))
maze.wallpattern = compressed[1] & 0x0f
maze.floorpattern = (compressed[1] & 0xf0) >> 4
maze.wallcolor = compressed[2] & 0x0f
maze.floorcolor = (compressed[2] & 0xf0) >> 4
htype1 := compressed[3]  // horiz type 1
htype2 := compressed[4]  // horiz type 2
vtype1 := compressed[5]  // vert type 1
vtype2 := compressed[6] // vert type 2

prev := htype2 // previous value, for 'repeat previous' types
// end test
	// fmt.Printf("Encoded size: %d\n", maze.encodedbytes)

	// Fill in first row with walls, always
	for i := 0; i < 32; i++ {
		maze.data[xy{i, 0}] = MAZEOBJ_WALL_REGULAR
	}

	// Unpack here starts
	location := 32               // how many spots we've filled
//	compressed = compressed[11:] // pointer to where we are in the input stream
	compressed = compressed[7:] // pointer to where we are in the input stream

	for location < 1024 {
		// fmt.Printf("input remaining: %d, next byte 0x%02x, output remaining: %d\n", len(compressed), compressed[0], 1024-location)
		if compressed[0] == 0 {
			fmt.Printf("WARNING: Read end of maze datastream before maze full.\n")
			break
		}
		var token int
		//      fmt.Printf("Remaining input length: %d, output remaining: %d\n", len(level), 1024-p)
		token, compressed = compressed[0], compressed[1:]
		count := (token & 0x0f) + 1
		longcount := (token & 0x1f) + 1 // used for 'repeat last' and 'skip'

fmt.Printf("Pos: %04d, left: %03d tok 0x%02x: count:%d lcnt: %d\n", location, len(compressed), token, count, longcount)

		switch token & 0xc0 { // look at top two bits
		case 0x00: // place one of literal object
			location = expand(maze, location, token&0x3f, 1)
			prev = token
		case 0x40: // Repeat special type
			switch token & 0x30 {
			case 0x00:
				prev = htype1
			case 0x10:
				prev = vtype1
			case 0x20:
				prev = htype2
			case 0x30:
				prev = vtype2
			}

			previtem := prev & 0x3f
			switch prev & 0xc0 {
			case 0x00: // repeat type
				if (token & 0x10) != 0 {
					location = vexpand(maze, location, previtem, count)
				} else {
					location = expand(maze, location, previtem, count)
				}
			case 0x40: // skip and add
				location = expand(maze, location, MAZEOBJ_TILE_FLOOR, count)
				location = expand(maze, location, previtem, 1)
			case 0x80: // add and skip
				location = expand(maze, location, previtem, 1)
				location = expand(maze, location, MAZEOBJ_TILE_FLOOR, count)
			case 0xc0: // repeat wall and add
				location = expand(maze, location, MAZEOBJ_WALL_REGULAR, count)
				location = expand(maze, location, previtem, 1)
			}
		case 0x80: // repeat wall
			if (token & 0x20) != 0 { // Repeat wall
				if (token & 0x10) != 0 {
					// vertical
					location = vexpand(maze, location, MAZEOBJ_WALL_REGULAR, count)
				} else {
					// horizontal
					location = expand(maze, location, MAZEOBJ_WALL_REGULAR, count)
				}
			} else {
				location = expand(maze, location, prev&0x3f, longcount)
			}
		case 0xc0:
			if (token & 0x20) != 0 {
				// skip and add wall
				location = expand(maze, location, MAZEOBJ_TILE_FLOOR, longcount)
				location = expand(maze, location, MAZEOBJ_WALL_REGULAR, 1)
			} else {
				// just skip
				location = expand(maze, location, MAZEOBJ_TILE_FLOOR, longcount)
			}
		}
	}

	if len(compressed) != 1 || compressed[0] != 0 {
		fmt.Printf("WARNING: Incomplete maze decode? (%d bytes remaining)\n", len(compressed))
	}

	return maze
}

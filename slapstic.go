package main

import (
	"encoding/binary"
	"fmt"
)

/* 
   g2 - bank info precalculated & addr of maze decoded from rom header per maze number
 *110[56].10[ab]
42a7f2b4a456e70319d1a5506341905f  ROMs/136043-1105.10a
c3feebd1f91ae90056351ffbbcee675b  ROMs/136043-1106.10b
-- all g2 sets use same maze roms
*/
var slapsticRoms = []string{
	"ROMs/136043-1105.10a",
	"ROMs/136043-1106.10b",
}
// g1 - since no bank info list provided for g1, we just load from direct addr
/*
 *20[56].10[ab] - r14
d06c71b1cf55cd3f637c94f3570b5450  ROMs-g1/136037-205.10a
8193af138bee2b76720709f42082a343  ROMs-g1/136037-206.10b
-- only r14 uses these
 *10[56].10[ab] - r7
862976922791fda377c23039db74c203  ROMs-g1/gauntletr7/136037-105.10a
16ce166415e8cdc678ef0411371ee004  ROMs-g1/gauntletr7/136037-106.10b
-- all g1 roms 105.10a / 106.10b are the same
*/
var slapsticRomsG1 = []string{
//	"ROMs-g1/gauntletr7/136037-105.10a",
//	"ROMs-g1/gauntletr7/136037-106.10b",
	"ROMs-g1/136037-205.10a",
	"ROMs-g1/136037-206.10b",
}

var slapsticBankInfo = []int{
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x54, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x95,
	0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xFE, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x03, 0xFC, 0x0E,
}

const (
	SLAPSTIC_START = 0x038000
)

// Do this the lazy way -- read an oversized chunk, then keep what we need
func slapsticReadMaze(mazenum int) []int {
	addr := 0x3f354
	if mazenum < 200 {
		addr = slapsticMazeGetRealAddr(mazenum)
	} else {
		addr = mazenum
	}
fmt.Printf("Maze real addr: 0x%06x\n", addr)

	b := slapsticReadBytes(addr, 512, mazenum)

	var intbuf []int
	for i := 0; true; i++ {
		intbuf = append(intbuf, int(b[i]))
		if i >= 11 && int(b[i]) == 0 {
			break
		}
	}

	//	return b[:i+1]
	return intbuf
}

func slapsticMazeGetRealAddr(mazenum int) int {
	bank := slapsticMazeGetBank(mazenum)
	addr := slapsticReadMazeOffset(mazenum) + (0x2000 * bank)

fmt.Printf("Maze real addr: 0x%06x, bank %d, boff: 0x%04x\n", addr, bank, 0x2000 * bank)
	return addr
}

func slapsticMazeGetBank(mazenum int) int {
	if mazenum < 0 || mazenum > 116 {
		panic("Invalid maze number requested (must be 0 <= x <= 116")
	}

	// More or less a direct port of the 68k cohde. Should improve this.
	offset := mazenum / 4
	bi := slapsticBankInfo[offset]
	offset = (mazenum % 4) * 2
	bi = bi >> uint(offset)
	bi = bi & 0x3

	return bi
}

func slapsticReadMazeOffset(mazenum int) int {

	buf := slapsticReadBytes(0x03800c+(4*mazenum), 4, mazenum)
	mazeoffset := binary.BigEndian.Uint32(buf)

fmt.Printf("Offset for maze: 0x%06x\n", mazeoffset)
fmt.Printf("big endian buf: %l\n", buf)

	return int(mazeoffset)
}

// Read bytes from combined ROM. Only works if reading an even address
func slapsticReadBytes(offset int, count int, mazn int) []byte {
	if offset >= SLAPSTIC_START {
		offset -= SLAPSTIC_START
	}
	buf := romSplitRead(slapsticRoms, offset, count)
	if mazn > 0x037FFF {
		buf = romSplitRead(slapsticRomsG1, offset, count)
	}

	return buf
}

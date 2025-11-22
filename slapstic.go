package main

import (
	"encoding/binary"
	"fmt"
)

var slapsticRoms = []string{
//	"ROMs/136043-1105.10a",
//	"ROMs/136043-1106.10b",
// g1 exper
	"ROMs/136037-205.10a",
	"ROMs/136037-206.10b",
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
// /	addr := slapsticMazeGetRealAddr(mazenum)
// TEMP remove - put above line back
addr := 0x38abe;
// TEMP remove
	b := slapsticReadBytes(addr, 512)

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
	buf := slapsticReadBytes(0x03800c+(4*mazenum), 4)
	mazeoffset := binary.BigEndian.Uint32(buf)

fmt.Printf("Offset for maze: 0x%06x\n", mazeoffset)
fmt.Printf("big endian buf: %l\n", buf)

	return int(mazeoffset)
}

// Read bytes from combined ROM. Only works if reading an even address
func slapsticReadBytes(offset int, count int) []byte {
	if offset >= SLAPSTIC_START {
		offset -= SLAPSTIC_START
	}
	buf := romSplitRead(slapsticRoms, offset, count)

	return buf
}

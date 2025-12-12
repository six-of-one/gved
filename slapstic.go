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

	"ROMs-g1/gauntletr7/136037-105.10a",
	"ROMs-g1/gauntletr7/136037-106.10b",

}
// seperate out g1 rev 14 roms, selectable by --r14
var slapsticRomsG1fr = []string{

	"ROMs-g1/136037-205.10a",
	"ROMs-g1/136037-206.10b",
}
var slapsticBankInfo = []int{
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x54, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x95,
	0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xFE, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x03, 0xFC, 0x0E,
}

// there is an odd issue with the following maze address reads
/*							maze start # = 1
52:246030 - 3C10E
53-246202 - 3C1BA
54-246554 - 3C31A
55-246874 - 3C45A

followup - in a side by side tracethru {newcode} with mazedumps_g1, these (new) mismatch: 51, 52, 53, 54

78:254218 - 3E10A
79-254442 - 3E1EA
80-254629 - 3E2A5
81-254852 - 3E384

followup - in a side by side tracethru {newcode} with mazedumps_g1, these (new) mismatch: 78, 79, 80, 81

matchups with _g1 dumps
- mazedumps_g1 start # = 0

mazedumps_g1 vs. r1-9_$N

051			:		078
052			:		079
053			:		080
054			:		081

077			:		052
078			:		053
079			:		054
080			:		055

// where the following values differ from what was manually discovered to load those mazes in g1rv7

246026 - 246030 =  -4	  bank adjust	1FFC
246250 - 246202 =  48					2030
246437 - 246554 =  -117					1F8B
246660 - 246874 =  -214					1F2A
254222 - 254218 =  4					2004
254394 - 254442 =  -48					1FD0
254746 - 254629 =  117					2075
255066 - 254852 =  214					20D6
*/

// manual build of bank # add value to come up with correct maze load addr on each maze num
var slapsticBankInfoG1 = []int{
																		// recorded from testing, ea. val is -1
	0, 0, 0, 0, 0, 0, 0, 0,												//	1, 2, 3, 4, 5, 6, 7, 8,
	0, 0, 0, 0, 0, 0, 0, 0,												//	9, 10, 11, 12, 13, 14, 15, 16,
	0, 0, 0, 0, 0, 0, 0, 0,												//	17, 18, 19, 20, 21, 22, 23, 24,
	0, 0, 0, 0, 0,														//	25, 26, 27, 28, 29,
	0x2000, 0x2000, 0x4000, 0x4000, 0x4000, 0x4000, 0x4000, 0x2000,		//	30, 31, 32, 33, 34, 35, 36, 37,
	0x2000, 0x2000, 0x2000, 0x2000, 0x2000, 0x2000, 0x2000, 0x2000,		//	38, 39, 40, 41, 42, 43, 44, 45,
	0x2000, 0x2000, 0x2000, 0x2000, 0x2000,								//	46, 47, 48, 49, 50,
	0x4000, 0x3ffc, 0x4030, 0x3f8b, 0x3f2a, 0x6000, 0x2000, 0x2000,		//	51, 52, 53, 54, 55, 56, 57, 58,
	0x2000, 0x2000, 0x2000, 0x2000, 0x2000, 0x2000, 0x2000, 0x2000,		//	59, 60, 61, 62, 63, 64, 65, 66,
	0x2000, 0x2000, 0x2000, 0x2000, 0x2000,								//	67, 68, 69, 70, 71,
	0x6000, 0x6000, 0x6000, 0x6000, 0x6000, 0x6000, 0x6004, 0x5fd0,		//	72, 73, 74, 75, 76, 77, 78, 79,
	0x6075, 0x60d6, 0x4000, 0x4000, 0x4000, 0x4000, 0x4000, 0x4000,		//	80, 81, 82, 83, 84, 85, 86, 87,
	0x4000, 0x4000, 0x4000, 0x4000, 0x4000, 0x4000, 0x6000, 0x6000,		//	88, 89, 90, 91, 92, 93, 94, 95,
	0x6000, 0x6000, 0x6000, 0x6000, 0x4000, 0x4000, 0x4000, 0x4000,		//	96, 97, 98, 99, 100, 101, 102, 103,
	0x4000, 0x4000, 0x4000, 0x4000, 0x4000, 0x4000, 0x4000, 0x4000,		//	104, 105, 106, 107, 108, 109, 110, 111,
	0x4000, 0x4000, 0x4000,												//	112, 113, 114,

// g1 address read direct for demo maze, score table maze, treasure rooms - these have no big endian store
	0x3F2F8, 0x3F357, 0x3F3DD, 0x3F479, 0x3F504, 0x3F654,				//	115, 116, 117, 118, 119, 120,
	0x3F6E6, 0x3F813, 0x3F940, 0x3FA22, 0x3FB20, 0x3FBD8, 0x3FCBD,		//	121, 122, 123, 124, 125, 126, 127,
}

const (
	SLAPSTIC_START = 0x038000
)

// Do this the lazy way -- read an oversized chunk, then keep what we need
func slapsticReadMaze(mazenum int) []int {

	addr := slapsticMazeGetRealAddr(mazenum)

// --ad={hex address} overrides maze read address here
	if Aov > 0 { addr = Aov }

if opts.Verbose { fmt.Printf("Maze read from: 0x%06x - %d\n", addr, addr) }

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
	addr := slapsticReadMazeOffset(mazenum,0xc) + (0x2000 * bank)

	if G1 {
// note: manual load of g1 banks, need proper algorithm for bank data with g1 slaps
		bankof := slapsticBankInfoG1[mazenum]
		if mazenum < 114 {
			addr = slapsticReadMazeOffset(mazenum,0x32) + bankof + 3
			bank = bankof / 0x2000
		} else {
			addr = slapsticBankInfoG1[mazenum]
		}
	}

if opts.Verbose { fmt.Printf("G:%d Maze real addr: %d - 0x%06X, bank %d, boff: 0x%04x\n", opts.Gtp, addr, addr, bank, 0x2000 * bank) }
	return addr
}

func slapsticMazeGetBank(mazenum int) int {

	if G1 { return 0 }	// g1 bank info not available here

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

// cohde philosophizing: do i add 3 and use the g2 decode (which works except floor/wall colors)
//                       or not add 3 and put in an adjust in mazedecode for g1 mazes having 3 more lead in bytes
// followup - disassemble trace of g1 & see if the decoder is different

func slapsticReadMazeOffset(mazenum int, x int) int {

	buf := slapsticReadBytes(0x038000+x+(4*mazenum), 4)
	mazeoffset := binary.BigEndian.Uint32(buf)

if opts.Verbose {
	if mazeoffset > 0x37fff && mazeoffset < 0x40000 {
		fmt.Printf("Offset for maze %d: %d - 0x%06x\n", mazenum, mazeoffset, mazeoffset)
		fmt.Printf("big endian buf: %l\n", buf)
	}}

	return int(mazeoffset)
}

// Read bytes from combined ROM. Only works if reading an even address
func slapsticReadBytes(offset int, count int) []byte {
	if offset >= SLAPSTIC_START {
		offset -= SLAPSTIC_START
	}
	buf := romSplitRead(slapsticRoms, offset, count)
	if G1 {
		if opts.R14 {
			buf = romSplitRead(slapsticRomsG1fr, offset, count)
		} else {
			buf = romSplitRead(slapsticRomsG1, offset, count)
		}
	}

	return buf
}

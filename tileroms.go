package main

import (
	"fmt"
)

type Romset struct {
	offset int
	roms   []string
}

/* rom md5s - these are used by all revs of g1 / g2

aef6687efa3a8dd75bbc5af9886bb56e  ROMs-g1/136037-111.1a
20b8d6bd306b258fb7d6dcac237dafa2  ROMs-g1/136037-113.1l

17f5397a1ca3cf35405131f1989da5ba  ROMs/136043-1111.1a
23e2a5ce11261217b048b58cd08a32d7  ROMs/136043-1113.1l

c53da9cf6911899b20e72408fb105188  ROMs-g1/136037-115.2a
2f607e4ddf4abcb31b8d505f082c7c4d  ROMs-g1/136037-117.2l

72bff2446feec00b69eaa73e02e6d8fe  ROMs/136043-1115.2a
9deeeaf33e12ddd7acaa7b750b050d0b  ROMs/136043-1117.2l

33ab85164c347033a7226d32ea4faec1  ROMs-g1/136037-112.1b
d78df568160063ed3c24717379a70090  ROMs-g1/136037-114.1mn

33ab85164c347033a7226d32ea4faec1  ROMs/136037-112.1b
d78df568160063ed3c24717379a70090  ROMs/136037-114.1mn

c2cb2ee155951d16f954813149a3adba  ROMs-g1/136037-116.2b
b2f90dee1b30dae1add893eebf88a26c  ROMs-g1/136037-118.2mn

c2cb2ee155951d16f954813149a3adba  ROMs/136037-116.2b
b2f90dee1b30dae1add893eebf88a26c  ROMs/136037-118.2mn

808a9f4c401f9138f51ccea00ceb5bc5  ROMs/136043-1123.1c
c985c04c6630836148e4a01dea9e61b8  ROMs/136043-1124.1p
9728a2d52192b0c212af0647b0dfb8c9  ROMs/136043-1125.2c
090ff7faf44f73f9cc5cdbc32eaac643  ROMs/136043-1126.2p
*/

// gauntlet II roms
var tileRoms = [][]string{
	{
		"ROMs/136043-1111.1a",
		"ROMs/136043-1113.1l",
		"ROMs/136043-1115.2a",
		"ROMs/136043-1117.2l",
	},
	{
		"ROMs/136037-112.1b",
		"ROMs/136037-114.1mn",
		"ROMs/136037-116.2b",
		"ROMs/136037-118.2mn",
	},
	{
		"ROMs/136043-1123.1c",
		"ROMs/136043-1124.1p",
		"ROMs/136043-1125.2c",
		"ROMs/136043-1126.2p",
	},
}

var tileRomSets = []Romset{
	{
		offset: 0x800,
		roms:   tileRoms[0],
	},
	{
		offset: 0x0,
		roms:   tileRoms[0],
	},
	{
		offset: 0x800,
		roms:   tileRoms[1],
	},
	{
		offset: 0x0,
		roms:   tileRoms[1],
	},
	{
		offset: 0x0,
		roms:   tileRoms[2],
	},
}
// gauntlet roms
var tileRomsG1 = [][]string{
	{
		"ROMs-g1/136037-111.1a",
		"ROMs-g1/136037-113.1l",
		"ROMs-g1/136037-115.2a",
		"ROMs-g1/136037-117.2l",
	},
// yes. these are a repeat of the g2 roms
// no. i dont feel like horking the code to exclude them, i like the structure the way it is
	{
		"ROMs-g1/136037-112.1b",
		"ROMs-g1/136037-114.1mn",
		"ROMs-g1/136037-116.2b",
		"ROMs-g1/136037-118.2mn",
	},
// these are the g2 - only gfx, leaving them, in case a g1 read slips in a g2 item some way
	{
		"ROMs/136043-1123.1c",
		"ROMs/136043-1124.1p",
		"ROMs/136043-1125.2c",
		"ROMs/136043-1126.2p",
	},
}

var tileRomSetsG1 = []Romset{
	{
		offset: 0x800,
		roms:   tileRomsG1[0],
	},
	{
		offset: 0x0,
		roms:   tileRomsG1[0],
	},
	{
		offset: 0x800,
		roms:   tileRomsG1[1],
	},
	{
		offset: 0x0,
		roms:   tileRomsG1[1],
	},
	{
		offset: 0x0,
		roms:   tileRomsG1[2],
	},
}

// returns the actual tile number to use, and the rom set to use it with
// Kind of a mess, since it uses knowledge for calculating the tile number
// that should be contained in the above structs, but isn't
func getromset(tilenum int) (int, []string) {
	whichbank := tilenum / 0x800
	actualtile := (tilenum % 0x800) + tileRomSets[whichbank].offset
	rombk := tileRomSets[whichbank].roms
	if G1 > 0 {
// in g1 mode - select g1 roms, actualtile will have same value, offsets are same
		rombk = tileRomSetsG1[whichbank].roms
	}

if false {
fmt.Printf("G: 0x%X  tn: 0x%x  tile: 0x%x   romset: %s\n", G1, tilenum, actualtile, rombk)  // this doesnt show which romfile used
}
	return actualtile, rombk
}

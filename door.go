package main

const (
	DOOR_HORIZ = iota
	DOOR_VERT
)

func doorGetTiles(doorDir int, doorAdj int) []int {
	t := make([]int, 4)
	m := doorStamps[doorAdj]

	// This is super sloppy
	if m == 0 {
		switch doorDir {
		case DOOR_HORIZ:
			fallthrough
		case DOOR_VERT:
			return nil
		}
	}

	for i := 0; i < 4; i++ {
		t[i] = m + i
	}

	return t
}

func doorGetStamp(doorDir int, doorAdj int) *Stamp {
	var stamp *Stamp
	tiles := doorGetTiles(doorDir, doorAdj)

  if  !svanim || (nothing & NODOR) == 0 {
	if tiles == nil {
		switch doorDir {
		case DOOR_HORIZ:
			stamp = itemGetStamp("hdoor")
		case DOOR_VERT:
			stamp = itemGetStamp("vdoor")
		}
	} else {
		stamp = genstamp_fromarray(tiles, 2, "base", 0)
		stamp.trans0 = true
	}
  }
	return stamp
}

// These are the first tile numbers, the next three are sequential for
// all door types
var doorStamps = []int{
	0x0000, // nothing adjacent
	0x0000, // only adjacent up
	0x0000, // only adjacent right
	0x1d34, // up right

	0x0000, // only adjacent down
	0x0000, // only adjacent down and up
	0x1d2c, // down right
	0x1d1c, // up-right-down

	0x0000, // only adjacent left
	0x1d38, // left-up
	0x0000, // only adjacent left and right
	0x1d24, // up-left-right

	0x1d30, // left-down
	0x1d18, // up--down-left
	0x1d20, // right left down
	0x1d28, // up down left right
}

package main

import (
	"math"
	"math/rand"
	"fmt"
	"time"
)

var RLPROF = [][]int{
	{64,  1,  0,  0,  0,  0,  0,  0,  2,  1,  0,  1,  1},		//		0x008090	// 	G1MP_TREASURE_BAG: 64,
	{53,  1,  1,  1,  0,  3,  1,  0,  1,  2,  0,  1,  0},		//		0x008050	// 	G1MP_KEY: 53,
	{54,  0,  0,  1,  2,  2,  0,  1,  0,  0,  0,  0,  1},		//		0x008061	// 	G1MP_POT_INVULN: 45,
	{44,  2,  2,  0,  3,  1,  2,  0,  1,  0,  1,  1,  0},		//		0x008060	// 	G1MP_POT_DESTRUCTABLE: 44,
	{40,  13, 20, 13, 20, 17, 14, 3, 11,  0,  0,  0,  0},		//		0x008070	// 	G1MP_TREASURE: 40,
	{42,  4,  4,  6,  4,  6,  8,  2,  2,  5,  2,  1,  0},		//		0x008000	// 	G1MP_FOOD_DESTRUCTABLE: 42,
	{43,  1,  0,  2,  2,  1,  0,  0,  1,  5,  0,  0,  1},		//		0x008020	// 	G1MP_FOOD_INVULN: 43,	- food h
	{43,  1,  1,  0,  2,  3,  0,  1,  1,  5,  0,  0,  0},		//		0x008040	food turk
	{46,  1,  0,  0,  1,  1,  0,  0,  0,  0,  0,  0,  0},		//		0x0080F0	// 	G1MP_INVISIBL: 46,
	{ 3,  1,  1,  2,  0,  3,  0,  1,  0, 16,  5,  0,  1},		//		0x0080D0	// 	G2MP_WALL_MOVABLE: 3,
	{62,  3,  0,  0,  4,  0,  4,  8, 44,  5, 10,  0,  0},		//		0x0080C0	// 	G1MP_TILE_STUN: 62,
	{ 6,  1,  1,  0,  2,  1,  2,  1,  3,  2,  2,  0,  0},		//		0x004000	// 	G1MP_EXIT: 6,
	{59,  1,  0,  3,  2,  6,  2,  0,  8, 10,  5,  0,  1},		//		0x0080A0	// 	G1MP_TRANSPORTER: 59,
	{27,  3,  5,  3,  5,  3,  0,  0,  0,  0,  0,  0,  0},		//		0xF00000	// 	G1MP_GEN_GHOST3: 27,
	{33,  2,  0,  0,  3,  3, 10,  0,  0,  0,  0,  0,  0},		//		0xF00010	// 	G1MP_GEN_DEMON3: 33,
	{30,  7,  0,  3,  7,  7,  0, 12,  0,  0,  0,  0,  0},		//		0xF00020	// 	G1MP_GEN_GRUNT3: 30,
	{39,  0,  0,  0,  0,  5,  0,  6,  0,  0,  0,  0,  0},		//		0xF00030	// 	G1MP_GEN_SORC3: 39,
	{36,  0,  0,  0,  2,  4,  0,  0,  0,  0,  0,  0,  0},		//		0xF00050	// 	G1MP_GEN_LOBBER3: 36,
	{26,  3,  7, 11, 10, 10,  0,  0,  0,  0,  0,  0,  1},		//		0xF00060	// 	G1MP_GEN_GHOST2: 26,
	{32,  0,  0,  0,  5,  5, 10,  0,  0,  0,  0,  0,  0},		//		0xF00070	// 	G1MP_GEN_DEMON2: 32,
	{29,  3,  0, 11, 12, 12,  0,  8,  0,  0,  0,  0,  0},		//		0xF00080	// 	G1MP_GEN_GRUNT2: 29,
	{38,  2,  0,  0,  0,  8,  0,  5,  0,  0,  0,  0,  0},		//		0xF00090	// 	G1MP_GEN_SORC2: 38,
	{35,  0,  0,  0,  1,  5,  0,  4,  0,  0,  0,  0,  0},		//		0xF000A0	// 	G1MP_GEN_LOBBER2: 35,
	{25,  3, 15,  5,  2,  2,  0,  0,  0,  0,  0,  0,  0},		//		0xF000B0	// 	G1MP_GEN_GHOST1: 25, 
	{31,  2,  0,  0, 10, 10, 10,  0,  0,  0,  0,  0,  0},		//		0xF000C0	// 	G1MP_GEN_DEMON1: 31,
	{28,  4,  0,  8,  5,  5,  0,  5,  9,  0,  0,  0,  0},		//		0xF000D0	// 	G1MP_GEN_GRUNT1: 28,
	{37,  0,  0,  0,  0,  8,  0,  8,  0,  0,  0,  0,  0},		//		0xF000E0	// 	G1MP_GEN_SORC1: 37,
	{34,  0,  0,  0,  3,  2,  0,  0,  0,  0,  0,  0,  0},		//		0xF000F0	// 	G1MP_GEN_LOBBER1: 34,
	{24,  2,  0,  1,  2,  0,  0,  0,  1, 24,  0,  0,  1},		//		0x400040	// 	G1MP_MONST_DEATH: 24,
	{11,  0,  6,  0,  0,  0,  0,  0,  0,  0,  0,  0,  1},		//		0x400000	// 	G1MP_MONST_GHOST3: 11,
	{17,  0,  0,  0,  0, 12,  7,  0,  0,  0,  0,  0,  0},		//		0x400010	// 	G1MP_MONST_DEMON3: 17,
	{14,  10, 0, 22,  0,  0,  0,  8,  0,  0,  0,  0,  0},		//		0x400020	// 	G1MP_MONST_GRUNT3: 14,
	{23,  0,  0,  0, 20,  0,  0, 12,  0,  0,  0,  0,  0},		//		0x400030	// 	G1MP_MONST_SORC3: 23,
	{20,  0,  0,  6,  0,  0,  0,  0, 15,  0,  0,  0,  0},		//		0x400050	// 	G1MP_MONST_LOBBER3: 20,
// col 10 is treasure room
// col 11,  12 extras
	{47,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  1},		//		0x0080E3	// 	G1MP_X_ARMOR: 47,
	{25,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  1},		//		0x400120	// 	G2MP_MONST_ACID: 25,
	{58,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  1,  1},		//		0x0080F4	// 	G2MP_POWER_SUPERSHOT: 58,
// dont use, not in system yet
	{0x008150,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  1},		//			// Lava - G3
	{0x008106,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  1,  1},		//			// fake food bottle - G3
}

var RLOAD = []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

// difficulty level for rnd load profile
var diff_level float64
var def_diff float64
var max_diff_level float64

// rload is sanctuary code port
// put rnd load in mbuf (with the stuff already there)
// anum is use mask, details on palette "T"

func rload(mbuf MazeData) {

source := rand.NewSource(time.Now().UnixNano())
rng := rand.New(source)

dx := opts.DimX
dy := opts.DimY
cx := 0
cy := 0

diff_level = 10.0
def_diff = 7.0
ldiff := diff_level / def_diff

rlloop := 33
rlline := 9
rprof := rng.Intn(rlline-1)+1 // for now pick a random profile
/*if troomtime > 0 {
	rprof = rlline + 1
}*/
fmt.Printf("%d prof, %f diff %f\n",rprof,float64(diff_level / def_diff),ldiff)
for f := 0; f <= rlloop; f++ {
	RLOAD[f] = int(math.Ceil(float64(RLPROF[f][rprof]) * ldiff)) // get item counts for a profile
fmt.Printf("%d rlprof: %d\n",f,RLOAD[f])
}
for f := 0; f <= rlloop; f++ {
	sft := 6000
	if anum < 1 || g1mask[RLPROF[f][0]] & anum > 0 {				// mask inclusive, see #T vals

		for RLOAD[f] > 0 && sft > 0 {
			fnd := false
			for !fnd && sft > 0 {

				cx = rand.Intn(dx)
				cy = rand.Intn(dy)
				if mbuf[xy{cx, cy}] == 0 { fnd = true }
				sft--
			}
			if fnd {
				mbuf[xy{cx, cy}] = RLPROF[f][0]
//			vartxt.value += fmt.Sprintf("\tSVRLOAD[%d][3][%d] = \"0x%x\";   //  x: %d y:  %d, w: %d h: %d\n",
//				svrcnt, (cell.tx + cell.ty*Mtw), RLPROF[f][0], cx, cy, Mtw, Mth)
				RLOAD[f]--
				// randomize the last 2
				if RLOAD[f] == 2 && rand.Float64() < 0.1 {
					RLOAD[f]--
				}
				if RLOAD[f] == 1 && rand.Float64() < 0.25 {
					RLOAD[f]--
				}
			}
		}
	}
}
}

// far goal rnd mapper

var (
	MAP_H int
	MAP_W int
	spots [100][100]int
)

// random int range

func rndr(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func _room(x1, y1, x2, y2, val int) {
	for y := y1; y < y2; y++ {
		for x := x1; x < x2; x++ {
			spots[y][x] = val
		}
	}
}

func map_fargoal(mbuf MazeData) {

	type point struct{ x, y int }
	rand.Seed(time.Now().UnixNano())

	MAP_H = opts.DimY - 1
	MAP_W = opts.DimX - 1

	room_center := make([]point, 10)

	dirs := []point{
		{1, 0},  // right
		{0, 1},  // down
		{0, -1}, // up
		{-1, 0}, // left
	}

	// Rooms
	for i := 0; i < 10; i++ {
		w := rndr(3, 7)
		h := rndr(3, 7)

		x := rndr(1, MAP_W-2-w)
		y := rndr(1, MAP_H-2-h)

		room_center[i] = point{x + w/2, y + h/2}

		_room(x, y, x+w, y+h, G1OBJ_TILE_FLOOR)
	}

	// Corridors
	for i := 0; i < 10; i++ {
		stone := 1 // 1 for initial, 2 for has hit stone, 0 for hit floor
		x := room_center[i].x
		y := room_center[i].y
		j := 0

		dir := rndr(0, 3)
		last := -1
		skip := 1

		for stone != 0 {
			if skip == 0 {
				if j > 1 {
					last = dir
				}
				for {
					dir = rndr(0, 3)
					if dir != 3-last {
						break
					}
				}
			} else {
				skip = 0
			}

			stone = 1
			run := rndr(0, 8) + 5
			j = 1

			for j != run {
				m_x := x + dirs[dir].x
				m_y := y + dirs[dir].y

				if spots[m_y][m_x] != G1OBJ_TILE_FLOOR {
					stone = 2
				}

				if m_x < 1 || m_x > MAP_W-2 || m_y < 1 || m_y > MAP_H-2 {
					break
				}

				if stone == 2 && spots[m_y][m_x] == G1OBJ_TILE_FLOOR {
					stone = 0
					break
				}

				spots[m_y][m_x] = G1OBJ_TILE_FLOOR

				x = m_x
				y = m_y

				j++
			}
		}
	}
}
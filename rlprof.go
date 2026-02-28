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

dx := opts.DimX
dy := opts.DimY
cx := 0
cy := 0

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

				cx = rng.Intn(dx)
				cy = rng.Intn(dy)
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

// far goal rnd mapper(s) 1, 2, and 3

var (
	MAP_H int
	MAP_W int
	gridb [100][100]int		// grid to build map
)

func _room(x1, y1, x2, y2, val int) {
// user selected a maze (or -1/rnd) to attemp items copy from
	var maze *Maze
	if xcont.Checked {
		maze = &Maze{}
		mazn,mt := rndr(7, maxmaze),0
		if xconsel.Text != "rnd" {
			fmt.Sscanf(xconsel.Text,"%d",&mt)
			if mt > 0 && mt <= maxmaze { mazn = mt - 1 }
		}
		svx, svy := opts.DimX, opts.DimY	// justDec resets these
		if Aov > 0 { Aov = addrver(slapsticMazeGetRealAddr(mazn)) + rndr(0, 15) - 5 }
		maze = justDecompress(slapsticReadMaze(mazn), false)
		opts.DimX, opts.DimY = svx, svy
	}

	for y := y1; y < y2; y++ {
		if y < 0 || y >= MAP_H { continue }
		for x := x1; x < x2; x++ {
			if x < 0 || x >= MAP_W { continue }
			gridb[y][x] = val
// replace some tiles with copy items
// thoughts: rnd on this? also allow T-flags to exclude items
		if xcont.Checked {
			if val == G1OBJ_TILE_FLOOR && maze.data[xy{x, y}] != G1OBJ_WALL_REGULAR {
				gridb[y][x] = maze.data[xy{x, y}]
			}}
		}
	}
}

func _attach(x1, y1, x2, y2 int, what int) (ax, ay int) {
	x, y := x1, y1

	for (x != x2 || y != y2) && gridb[y][x] != what {
		dx, dy := 0, 0

		if x < x2 {
			dx = 1
		}
		if x > x2 {
			dx = -1
		}
		if y < y2 {
			dy = 1
		}
		if y > y2 {
			dy = -1
		}
		if dx != 0 && dy != 0 {
			f := rndr(0, 1)
			if f != 0 {
				dx = 0
			} else {
				dy = 0
			}
		}
		x += dx
		y += dy
	}
	return x, y
}

func _path(x1, y1, x2, y2 int, what int) {
	x, y := x1, y1

	for x != x2 || y != y2 {
		l := rndr(1, 4)
		dx, dy := 0, 0

		if x < x2 {
			dx = 1
		}
		if x > x2 {
			dx = -1
		}
		if y < y2 {
			dy = 1
		}
		if y > y2 {
			dy = -1
		}
		if dx != 0 && dy != 0 {
			f := rndr(0, 1)
			if f != 0 {
				dx = 0
			} else {
				dy = 0
			}
		}
		for i := 0; i < l; i++ {
			gridb[y][x] = what
			x += dx
			y += dy
			if x == x2 && y == y2 {
				break
			}
			if !(x > 0 && y > 0 && x < MAP_W-1 && y < MAP_H-1) {
				x -= dx
				y -= dy
				break
			}
		}
	}
	gridb[y][x] = what
}

func _corridor(x, y int, where, what int) (ex, ey int) {
	ddx := []int{0, 1, 0, -1}
	ddy := []int{1, 0, -1, 0}
	l := rndr(10, 70)
	dx, dy := 0, 0
	s := 0

	for i := 0; i < l; i++ {
		if s > 0 {
			s--
		} else {
			d := rndr(0, 3)
			dx = ddx[d]
			dy = ddy[d]
			s = rndr(2, 8)
		}

		gridb[y][x] = what

		x += dx
		y += dy

		if x > 0 && y > 0 && x < MAP_W-1 && y < MAP_H-1 &&
			gridb[y+dy][x+dx] == where &&
			gridb[y+dy+dx][x+dx+dy] == where &&
			gridb[y+dy-dx][x+dx-dy] == where &&
			gridb[y+dx][x+dy] == where &&
			gridb[y-dx][x-dy] == where {
			// continue
		} else {
			x -= dx
			y -= dy
			s = 0
		}
	}
	return x, y
}

func grid_put(x, y int, val int) {
	if x >= 0 && x < MAP_W && y >= 0 && y < MAP_H {
		gridb[y][x] = val
	}
}

func grid_get(x, y int) int {
	if x >= 0 && x < MAP_W && y >= 0 && y < MAP_H {
		return gridb[y][x]
	}
	return G1OBJ_WALL_REGULAR
}

// 8 ray test from a cell
// for bounds lx,ly - mx,my (low to max) check tspot at tx,ty for tv (test val) if so, return rv

type point struct{ x, y int }

var dirs = []point{
		{1, 0},  // right
		{0, 1},  // down
		{0, -1}, // up
		{-1, 0}, // left
// expand orig 8 ray test around cell
		{-1, -1},// up - lf
		{-1, 1}, // dn - lf
		{1, -1}, // up - rt
		{1, 1},  // dn - rt
	}

func ray(lx, ly, mx, my, tx, ty, tv, rv int,tspot [100][100]int) int {

	r := -1
	if tx >= lx && ty >= ly && tx <= mx && ty <= my {
		if tspot[ty][tx] == tv { r = rv }
	}
	return r
}

// mapper 1: std map rooms + corridors

func map_fargoal(mbuf MazeData) {

//	rand.Seed(time.Now().UnixNano())
/*
	for y := 0; y <= opts.DimY; y++ {
		for x := 0; x <= opts.DimX; x++ {
		mbuf[xy{x, y}] = G1OBJ_WALL_REGULAR
	}}
*/
	MAP_H = opts.DimY
	MAP_W = opts.DimX

	for y := 0; y <= MAP_H; y++ {
		for x := 0; x <= MAP_W; x++ {
		gridb[y][x] = -1
	}}

	room_center := make([]point, 10)

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
fmt.Printf("run %d, stone %d\n",run,stone)

			for j != run {
				m_x := x + dirs[dir].x
				m_y := y + dirs[dir].y

				if gridb[m_y][m_x] != G1OBJ_TILE_FLOOR {
					stone = 2
				}

				if m_x < 1 || m_x > MAP_W-2 || m_y < 1 || m_y > MAP_H-2 {
					break
				}

				if stone == 2 && gridb[m_y][m_x] == G1OBJ_TILE_FLOOR {
					stone = 0
					break
				}

				gridb[m_y][m_x] = G1OBJ_TILE_FLOOR

				x = m_x
				y = m_y

				j++
			}
		}
	}
// wall off floor from null space - pre load room copy is an issue here
	for y := 1; y <= MAP_H; y++ {
		for x := 1; x <= MAP_W; x++ {
		if gridb[y][x] < 0 {
		for i := 0; i < 8; i++ {
			nv := ray(1, 1, MAP_W, MAP_H, x + dirs[i].x, y + dirs[i].y, G1OBJ_TILE_FLOOR, G1OBJ_WALL_REGULAR, gridb)
			if nv >= 0 { gridb[y][x] = nv }
		}
		}
	}}

	for y := 1; y <= MAP_H; y++ {
		for x := 1; x <= MAP_W; x++ {
		mbuf[xy{x, y}] = gridb[y][x]
	}}
}

// mapper 2: more complex maze

func map_sword(mbuf MazeData) {

var sword bool
	SPOT_MARKER := 256

//	rand.Seed(time.Now().UnixNano())

//	opts.DimY = 24		// these seem not working this way now?
//	opts.DimX = 39
/*
	for y := 0; y <= opts.DimY; y++ {
		for x := 0; x <= opts.DimX; x++ {
		mbuf[xy{x, y}] = G1OBJ_WALL_REGULAR
	}}*/

	MAP_H = opts.DimY
	MAP_W = opts.DimX

	for y := 0; y <= MAP_H; y++ {
		for x := 0; x <= MAP_W; x++ {
		gridb[y][x] = G1OBJ_WALL_REGULAR
	}}

	for y := 0; y < MAP_H+1; y++ {
		grid_put(MAP_W, y, G1OBJ_TILE_FLOOR)
	}
	for x := 0; x < MAP_W+1; x++ {
		grid_put(x, MAP_H, G1OBJ_TILE_FLOOR)
	}

	x, y := 1, 2
	_room(17, 10, 21, 14, G1OBJ_TILE_FLOOR)
	grid_put(x, y, SPOT_MARKER+4)

	for {
		for {
			dir := rndr(0, 3)
			last := dir
			for {
				x2 := x + dirs[dir].x*2
				y2 := y + dirs[dir].y*2

				if x2 > 0 && y2 > 0 && x2 < MAP_W-1 && y2 < MAP_H-1 &&
					grid_get(x2, y2) == G1OBJ_WALL_REGULAR {
					grid_put(x2, y2, SPOT_MARKER+dir)
					grid_put(x+dirs[dir].x, y+dirs[dir].y, G1OBJ_TILE_FLOOR)
					x = x2
					y = y2
					break
				}

				dir++
				if dir == 4 {
					dir = 0
				}
				if dir == last {
					goto breakbreak
				}
			}
		}
	breakbreak:
		dir := int(grid_get(x, y))
		grid_put(x, y, G1OBJ_TILE_FLOOR)
		if dir >= SPOT_MARKER && dir <= SPOT_MARKER+3 {
			dir -= SPOT_MARKER
			x -= dirs[dir].x * 2
			y -= dirs[dir].y * 2
		} else {
			break
		}
	}

	dir := rndr(0, 3)
	switch dir {
	case 0:
		grid_put(16, 12, G1OBJ_TREASURE_BAG)
	case 1:
		grid_put(22, 12, G1OBJ_TREASURE_BAG)
	case 2:
		grid_put(19, 9, G1OBJ_TREASURE_BAG)
	case 3:
		grid_put(19, 15, G1OBJ_TREASURE_BAG)
	}

	if !sword {
		sw := rng.Intn(6) + G1OBJ_INVISIBL
		grid_put(19, 12, sw)
		sword = true
	}
/*
	for y := 0; y < 25; y++ {
		grid_put(39, y, G1OBJ_WALL_REGULAR)
	}
	for x := 0; x < 40; x++ {
		grid_put(x, 24, G1OBJ_WALL_REGULAR)
	}
*/
	for y := 1; y <= MAP_H; y++ {
		for x := 1; x <= MAP_W; x++ {
		mbuf[xy{x, y}] = gridb[y][x]
	}}
}

// mapper 3

func map_wide(mbuf MazeData) {

//	opts.DimY = 24
//	opts.DimX = 39
	MAP_H = opts.DimY + 1
	MAP_W = opts.DimX + 1
/*
	for y := 0; y <= opts.DimY; y++ {
		for x := 0; x <= opts.DimX; x++ {
		mbuf[xy{x, y}] = G1OBJ_WALL_REGULAR
	}}
*/
	MAP_H = opts.DimY
	MAP_W = opts.DimX

	for y := 0; y <= MAP_H; y++ {
		for x := 0; x <= MAP_W; x++ {
		gridb[y][x] = -1
	}}

	var sx, sy [15]int
	l, t := 2, 2
	mx, my := MAP_W/2, MAP_H/2
	r, b := MAP_W-3, MAP_H-3

	// Random room positions.
	sx[0], sy[0] = rndr(l, r), rndr(t, b)
	sx[1], sy[1] = rndr(l, mx), rndr(t, my)
	sx[2], sy[2] = rndr(mx, r), rndr(t, my)
	sx[3], sy[3] = rndr(l, mx), rndr(my, b)
	sx[4], sy[4] = rndr(mx, r), rndr(my, b)

	n := rndr(2, 7)

	for i := 5; i < n; i++ {
		sx[i], sy[i] = rndr(l, r), rndr(t, b)
	}

	// Place rooms.
	for i := 0; i < n; i++ {
		_room(sx[i]-rndr(1, 4), sy[i]-rndr(1, 4), sx[i]+rndr(1, 4), sy[i]+rndr(1, 4), G1OBJ_TILE_FLOOR)
	}

	// Connect rooms.
	for i := 0; i < n; i++ {
		j := i
		if i < n-1 {
			j = i + 1
		} else {
			j = 0
		}
		ax, ay := _attach(sx[i], sy[i], sx[j], sy[j], G1OBJ_WALL_REGULAR)
		gridb[ay][ax] = G1OBJ_EXIT4
		ax2, ay2 := _attach(sx[j], sy[j], sx[i], sy[i], G1OBJ_WALL_REGULAR)
		gridb[ay2][ax2] = G1OBJ_EXIT
		_path(ax, ay, ax2, ay2, G1OBJ_TILE_FLOOR)
	}

fmt.Printf("random corridors\n")
	// Some random corridors
	m := rndr(2, 7)
	for i := 0; i < m; i++ {
		var ex, ey, x, y int
		f := rndr(0, 3)

		for {
			x = rndr(1, MAP_W-2)
			y = rndr(4, MAP_H-2)
			if gridb[y][x] == -1 {
				break
			}
		}

		switch f {
		case 0:
			y = MAP_H - 2
		case 1:
			x = 1
		case 2:
			y = 4
		case 3:
			x = MAP_W - 2
		}

		ex, ey = _corridor(x, y, G1OBJ_WALL_REGULAR, G1OBJ_TILE_FLOOR)

		// Avoid dead ends
		j := rndr(0, n-1)
		_path(ex, ey, sx[j], sy[j], G1OBJ_TILE_FLOOR)
		j = rndr(0, n-1)
		_path(x, y, sx[j], sy[j], G1OBJ_TILE_FLOOR)
	}

// wall off floor from null space
	for y := 1; y <= MAP_H; y++ {
		for x := 1; x <= MAP_W; x++ {
		if gridb[y][x] < 0 {
		for i := 0; i < 8; i++ {
			nv := ray(1, 1, MAP_W, MAP_H, x + dirs[i].x, y + dirs[i].y, G1OBJ_TILE_FLOOR, G1OBJ_WALL_REGULAR, gridb)
			if nv >= 0 { gridb[y][x] = nv }
		}
		}
	}}

	for y := 1; y <= MAP_H; y++ {
		for x := 1; x <= MAP_W; x++ {
		mbuf[xy{x, y}] = gridb[y][x]
	}}
}

// all randomizers in G¹G²ved pre-suppose maze has been blanked or walled ahead of time
// the floowing use all walls & carve out

// DFS mapper

type TPoint struct {
	x, y int
}

var DFSdirections = [4]int{0, 1, 2, 3}

func InitializeDFSMaze() {
//	rand.Seed(time.Now().UnixNano())
 //initialize direction vector
  // Shuffle the directions array to randomize the order
	for x := 0; x <= 3; x++ {
		y := rand.Intn(4)
//		t := DFSdirections[x]
		DFSdirections[x], DFSdirections[y] = is(DFSdirections[x], DFSdirections[y])
//		DFSdirections[y] = t
	}
}

func GenerateDFSMaze(mdat MazeData, startX, startY, x, y, BiasCoefficient int) {

	if x < startX {
		x = startX
	}
	if y < startY {
		y = startY
	}

	mdat[xy{x,y}] = G1OBJ_TILE_FLOOR // Mark the current cell as walkable

// Shuffle again with a random bias
	for i := 0; i < BiasCoefficient; i++ {
		j := rand.Intn(BiasCoefficient)
/*		temp := DFSdirections[i]
		DFSdirections[i] = DFSdirections[j]
		DFSdirections[j] = temp */
		DFSdirections[i], DFSdirections[j] = is(DFSdirections[i], DFSdirections[j])
	}

  // Explore each direction
	for i := 0; i <= 3; i++ {
		var dx, dy int
		switch DFSdirections[i] {
		case 0: // Up
			dx, dy = 0, -1
		case 1: // Right
			dx, dy = 1, 0
		case 2: // Down
			dx, dy = 0, 1
		case 3: // Left
			dx, dy = -1, 0
		}

		nx := x + dx*2
		ny := y + dy*2

		if nx >= startX && nx <= opts.DimX && ny >= startY && ny <= opts.DimY && mdat[xy{nx,ny}] == G1OBJ_WALL_REGULAR {
			mdat[xy{x+dx,y+dy}] = G1OBJ_TILE_FLOOR // Carve a path
			GenerateDFSMaze(mdat, startX, startY, nx, ny, BiasCoefficient)
	// Recursively generate the maze
	// uhh....
//fmt.Printf("wut: %d x %d, sxy: %d x %d  b: %d\n",nx, ny,startX, startY,BiasCoefficient)

		}
	}
}

func map_dfs(mdat MazeData) {

	bias := int(time.Now().UnixNano() & 3)
	InitializeDFSMaze()
	GenerateDFSMaze(mdat,1,1,mxmd,mymd,bias)
}

// Prim maze gen

func GeneratePrimMaze(mdat MazeData, startX, startY int) {

	frontier := []TPoint{}
	dirs := []TPoint{{0, -2}, {2, 0}, {0, 2}, {-2, 0}}

	startY--		// this could hit -1
	startX--

// Start with a random cell
	current := TPoint{
		x: 1 + rand.Intn(opts.DimX/2)*2,
		y: 1 + rand.Intn(opts.DimY/2)*2,
	}
	mdat[xy{current.x,current.y}] = G1OBJ_TILE_FLOOR

// Add neighboring walls to the frontier
	for i := 0; i <= 3; i++ {
		nx := current.x + dirs[i].x
		ny := current.y + dirs[i].y
		if nx > startX && nx <= opts.DimX && ny > startY && ny <= opts.DimY && mdat[xy{nx,ny}] == G1OBJ_WALL_REGULAR {
			frontier = append(frontier, TPoint{nx, ny})
		}
	}

// Pick a random frontier cell
	for len(frontier) > 0 {
		index := rand.Intn(len(frontier))
		current = frontier[index]
		frontier[index] = frontier[len(frontier)-1]
		frontier = frontier[:len(frontier)-1]

// Connect it to the maze
		for i := 0; i <= 3; i++ {
			nx := current.x + dirs[i].x
			ny := current.y + dirs[i].y
			if nx > startX && nx <= opts.DimX && ny > startY && ny <= opts.DimY && mdat[xy{nx,ny}] == G1OBJ_TILE_FLOOR {
				mdat[xy{current.x,current.y}] = 0
				mdat[xy{(current.x+nx)/2,(current.y+ny)/2}] = G1OBJ_TILE_FLOOR
				break
			}
		}

 // Add new frontier cells
		for i := 0; i <= 3; i++ {
			nx := current.x + dirs[i].x
			ny := current.y + dirs[i].y
			if nx > startX && nx <= opts.DimX && ny > startY && ny <= opts.DimY && mdat[xy{nx,ny}] == G1OBJ_WALL_REGULAR {
				frontier = append(frontier, TPoint{nx, ny})
			}
		}
	}
}

// wall reducer came with DFS / Prim - sounds like a neet idea

func ReduceWalls(mdat MazeData, startX, startY int) {
	countLiveNeighbours := func(x, y int) int {
		result := 0
		if x > 0 {
			result += mdat[xy{x-1,y}]
			if y > 0 {
				result += mdat[xy{x-1,y-1}]
			}
			if y < opts.DimY {
				result += mdat[xy{x-1,y+1}]
			}
		}
		if x < opts.DimX {
			result += mdat[xy{x+1,y}]
			if y > 0 {
				result += mdat[xy{x+1,y-1}]
			}
			if y < opts.DimY {
				result += mdat[xy{x+1,y+1}]
			}
		}
		if y > 0 {
			result += mdat[xy{x,y-1}]
		}
		if y < opts.DimY {
			result += mdat[xy{x,y+1}]
		}
		return result / G1OBJ_WALL_REGULAR
	}

	for x := startX; x <= opts.DimX; x++ {
		for y := startY; y <= opts.DimY; y++ {
			if countLiveNeighbours(x, y) < 2 {
				mdat[xy{x,y}] = G1OBJ_TILE_FLOOR
			}
		}
	}
}

// take a quarter from 4 random mazes and put together

func Map4quart(mdat MazeData) {

	var maze = &Maze{}

//	opts.DimX, opts.DimY = 31,31		// for simplicity for now
										// caveats - too small maze copys into null space, too large only fills a portion...
	svx, svy := opts.DimX, opts.DimY	// justDec resets these
	aovs := Aov

	mazn := rndr(0, maxmaze)
	if Aov > 0 { Aov = addrver(slapsticMazeGetRealAddr(mazn)) + rndr(0, 15) - 5 }		// this is less productive
	maze = justDecompress(slapsticReadMaze(mazn), false)
	for y := 0; y <= 15; y++ {
		for x := 0; x <=15; x++ {
			mdat[xy{x,y}] = maze.data[xy{x,y}]
	}}
	mazn = rndr(0, maxmaze)
	if Aov > 0 { Aov = addrver(slapsticMazeGetRealAddr(mazn)) + rndr(0, 15) - 5 }
	maze = justDecompress(slapsticReadMaze(mazn), false)
	for y := 15; y <= opts.DimY; y++ {
		for x := 0; x <= 15; x++ {
			mdat[xy{x,y}] = maze.data[xy{x,y}]
	}}
	mazn = rndr(0, maxmaze)
	if Aov > 0 { Aov = addrver(slapsticMazeGetRealAddr(mazn)) + rndr(0, 15) - 5 }
	maze = justDecompress(slapsticReadMaze(mazn), false)
	for y := 0; y <= opts.DimY; y++ {
		for x := 15; x <= opts.DimX; x++ {
			mdat[xy{x,y}] = maze.data[xy{x,y}]
	}}
	mazn = rndr(0, maxmaze)
	if Aov > 0 { Aov = addrver(slapsticMazeGetRealAddr(mazn)) +rndr(0, 15) - 5 }
	maze = justDecompress(slapsticReadMaze(mazn), false)
	for y := 15; y <= opts.DimY; y++ {
		for x := 15; x <= opts.DimX; x++ {
			mdat[xy{x,y}] = maze.data[xy{x,y}]
	}}
	Aov = aovs
	opts.DimX, opts.DimY = svx, svy
}

// Kruskal algo

type DisjointSet struct {
	parent []int
	rank   []int
	size   int
}

func NewDisjointSet(size int) *DisjointSet {
	ds := &DisjointSet{
		parent: make([]int, size),
		rank:   make([]int, size),
		size:   size,
	}
	for i := 0; i < size; i++ {
		ds.parent[i] = i
		ds.rank[i] = 0
	}
	return ds
}

func (ds *DisjointSet) Find(i int) int {
	if ds.parent[i] != i {
		ds.parent[i] = ds.Find(ds.parent[i])
	}
	return ds.parent[i]
}

func (ds *DisjointSet) Union(i, j int) {
	rootI := ds.Find(i)
	rootJ := ds.Find(j)
	if rootI != rootJ {
		if ds.rank[rootI] < ds.rank[rootJ] {
			ds.parent[rootI] = rootJ
		} else if ds.rank[rootI] > ds.rank[rootJ] {
			ds.parent[rootJ] = rootI
		} else {
			ds.parent[rootJ] = rootI
			ds.rank[rootI]++
		}
	}
}

func (ds *DisjointSet) GetSize() int {
	return ds.size
}

type Edge struct {
	u, v, weight int
}

const (
	MAZE_SIZE   = 16 // example size, adjust as needed
	MAX_CELLS   = MAZE_SIZE * MAZE_SIZE
//	WALL_GEN_ID = -1
)

type GauntMap [32][32]int

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func GenerateKruskalMaze(mdat MazeData, startX, startY, weightX, weightY int) {
	rand.Seed(time.Now().UnixNano())

	var FMaze [32][32]int

	// Initialize maze with walls
	for y := 0; y < MAZE_SIZE*2+1; y++ {
		for x := 0; x < MAZE_SIZE*2+1; x++ {
			FMaze[x][y] = G1OBJ_WALL_REGULAR
		}
	}

	// Create edges
	var edges []Edge
	for y := 0; y < MAZE_SIZE; y++ {
		for x := 0; x < MAZE_SIZE; x++ {
			u := y*MAZE_SIZE + x
			if x < MAZE_SIZE-1 {
				edges = append(edges, Edge{
					u:      u,
					v:      u + 1,
					weight: rand.Intn(weightX),
				})
			}
			if y < MAZE_SIZE-1 {
				edges = append(edges, Edge{
					u:      u,
					v:      u + MAZE_SIZE,
					weight: rand.Intn(weightY),
				})
			}
		}
	}

	// Sort edges by weight (Bubble Sort)
	for i := 0; i < len(edges)-1; i++ {
		for j := i + 1; j < len(edges); j++ {
			if edges[i].weight > edges[j].weight {
				edges[i], edges[j] = edges[j], edges[i]
			}
		}
	}

	// Kruskal's algorithm
	ds := NewDisjointSet(MAX_CELLS)
	for _, edge := range edges {
		if ds.Find(edge.u) != ds.Find(edge.v) {
			ds.Union(edge.u, edge.v)
			// Create passage in the maze
			u := edge.u
			v := edge.v
			x := u % MAZE_SIZE
			y := u / MAZE_SIZE
			FMaze[x*2+1][y*2+1] = 0
			x = v % MAZE_SIZE
			y = v / MAZE_SIZE
			FMaze[x*2+1][y*2+1] = 0
			if edge.u+1 == edge.v {
				FMaze[min(edge.u%MAZE_SIZE, edge.v%MAZE_SIZE)*2+2][edge.u/MAZE_SIZE*2+1] = 0
			} else {
				FMaze[edge.u%MAZE_SIZE*2+1][min(edge.u/MAZE_SIZE, edge.v/MAZE_SIZE)*2+2] = 0
			}
		}
	}

	// Copy the temporary array to the actual maze
	for x := startX; x < 32; x++ {
		for y := startY; y < 32; y++ {
			mdat[xy{x,y}] = FMaze[x][y]
		}
	}
}

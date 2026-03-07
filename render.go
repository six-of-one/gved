package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"fmt"
	"time"
	"golang.org/x/image/draw"
)

func gettiledatafromfile(file string, tilenum int) TileLinePlane {
	f, err := os.Open(file)
	check(err)

	databytes := make([]byte, 8)

	f.Seek(int64(tilenum*8), 0)
	cnt, err := f.Read(databytes)
	check(err)

//fmt.Printf("tn: 0x%x  file: %s\n", tilenum, file)

	if cnt != 8 {
		panic("failed to read full tile from file")
	}

	f.Close()
	return databytes
}

func bytetobits(databyte byte) []byte {
	res := make([]byte, 8)

	// databyte = databyte ^ 0xff
	for i := 0; i < 8; i++ {
		if databyte%2 > 0 {
			res = append([]byte{0}, res...)
		} else {
			res = append([]byte{1}, res...)
		}
		databyte = databyte >> 1
	}

	return res
}

func mergeplanes(planes TileLinePlaneSet) TileLineMerged {
	mergedline := TileLineMerged{}

	for i := 0; i < 8; i++ {
		val := (planes[3][i] * byte(math.Pow(2, 3))) + (planes[2][i] * byte(math.Pow(2, 2))) + (planes[1][i] * byte(math.Pow(2, 1))) + (planes[0][i] * byte(math.Pow(2, 0)))
		mergedline = append(mergedline, val)
	}

	return mergedline
}

func blankimage(x int, y int) *image.NRGBA {
	rect := image.Rect(0, 0, x, y)

	// palette 0 (more or less), for exits and such
	//	palette := gauntletPalettes[opts.PalType][opts.PalNum]
	img := image.NewNRGBA(rect)
	return img
}

// image failed to load, make something halflife-esq tex missing
var ova,ovb color.Color = HRGB{0xff8f008f},HRGB{0xff008f8f}

func loadfail(x int, y int) *image.NRGBA {
	img := blankimage(x,y)
//	b := HRGB{0xff008f8f}
//	a := HRGB{0xff8f008f}
	var c color.Color
	for fy := 0; fy < y; fy += 8 {
		if fy & 8 == 8 { c = ova } else { c = ovb }
		for fx := 0; fx < x; fx += 8 {

		for j := 0; j < 8; j++ {
			for i := 0; i < 8; i++ {
				img.Set(fx+i, fy+j, c)
			}}
			if c == ovb { c = ova } else { c = ovb }
		}
	}
	ova,ovb = HRGB{0xff8f008f},HRGB{0xff008f8f}
	return img
}

func getparsedtile(tilenum int) TileData {
	planedata := make([]TileLinePlane, 4)

	realtilenum, roms := getromset(tilenum)

	planedata[0] = gettiledatafromfile(roms[0], realtilenum)
	planedata[1] = gettiledatafromfile(roms[1], realtilenum)
	planedata[2] = gettiledatafromfile(roms[2], realtilenum)
	planedata[3] = gettiledatafromfile(roms[3], realtilenum)
	// fmt.Printf("planedata is: %d\n", planedata)

	// fulltile := Tile{}
	fulltile := make([]TileLineMerged, 8)

	// For each line in tile
	for i := 0; i < 8; i++ {
		linedata := make([][]byte, 4)
		linedata[0] = bytetobits(planedata[0][i])
		linedata[1] = bytetobits(planedata[1][i])
		linedata[2] = bytetobits(planedata[2][i])
		linedata[3] = bytetobits(planedata[3][i])
		// fmt.Printf("line is: %d\n", linedata)

		fulltile[i] = mergeplanes(linedata)
		// fmt.Printf("merged line is: %d\n", fulltile[i])
	}

	// fmt.Printf("tile is: %d\n", fulltile)
	return fulltile
}

// Write an 8x8 tile into a (usually) larger image
func writetiletoimage(img *image.NRGBA, tile TileData, palette []color.Color, trans0 bool, x int, y int) {
	for j := 0; j < 8; j++ {
		for i := 0; i < 8; i++ {
			tc := tile[j][i]
			if tc == 0 && trans0 {
				continue
			}

			c := palette[tc]
			img.Set(x+i, y+j, c)
			// fmt.Printf("%x", tc)
		}
		// fmt.Printf("\n")
	}
}

func genimage(tilenum int, xtiles int, ytiles int) *image.NRGBA {
	t := make([]int, xtiles*ytiles)

	for i := 0; i < (xtiles * ytiles); i++ {
		t[i] = tilenum + i
	}

	return genimage_fromarray(t, xtiles, ytiles)
}

func genstamp_fromarray(tiles []int, width int, ptype string, pnum int) *Stamp {
	stamp := Stamp{
		width: width,
		ptype: ptype,
		pnum:  pnum,
	}

	stamp.numbers = tiles
	stamp.data = make([]TileData, len(tiles))

	fillstamp(&stamp)

	return &stamp
}

func fillstamp(stamp *Stamp) {
	tc := 0
	height := len(stamp.numbers) / stamp.width

	stamp.data = make([]TileData, len(stamp.numbers))
	// spew.Dump("Stamp: ", stamp)
	for y := 0; y < height; y++ {
		for x := 0; x < stamp.width; x++ {
			stamp.data[(stamp.width*y)+x] = getparsedtile(stamp.numbers[tc])
			tc++
		}
	}
}

func genimage_fromarray(tiles []int, xtiles int, ytiles int) *image.NRGBA {
	stamp := genstamp_fromarray(tiles, xtiles, opts.PalType, opts.PalNum)

	img := blankimage(8*xtiles, 8*ytiles)
	writestamptoimage(G1,img, stamp, 0, 0)

	return img
}

func writestamptoimage(g1 bool, img *image.NRGBA, stamp *Stamp, xloc int, yloc int) {
	ptyp := stamp.ptype
// gauntlet has diff base pallet from G², this is easiest way for now
	if ptyp == "base" && g1 { ptyp = "gbase" }
	if ptyp == "wall" && g1 { ptyp = "gwall" }
	if ptyp == "floor" && g1 { ptyp = "gfloor" }
//fmt.Printf("g palettes %s,%d\n", ptyp, stamp.pnum)
	p := gauntletPalettes[ptyp][stamp.pnum]
	for y := 0; y < len(stamp.data)/stamp.width; y++ {
		for x := 0; x < stamp.width; x++ {
			// fmt.Printf("Writing to image at %d,%d\n", xloc, yloc)
			writetiletoimage(img, stamp.data[(stamp.width*y)+x], p,
				stamp.trans0, xloc+(x*8), yloc+(y*8))
		}
	}
if false {
fmt.Printf("stamp (msk) %d @ %d x %d\n ",stamp.mask,xloc, yloc)
}
}

func savetopng(fn string, img *image.NRGBA) {
	f, _ := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE, 0600)
if opts.Verbose {fmt.Println(fn)}
	defer f.Close()
	png.Encode(f, img)
}

/*
// while this compiles, with opts.Animate in monsters.go it produces blak images
func genanim(animarray []int, xtiles int, ytiles int) []*image.Paletted {
	var images []*image.Paletted

	for i := 0; i < len(animarray); i++ {
		fmt.Printf("generating from tile %d\n", animarray[i])
		src := genimage(animarray[i], xtiles, ytiles)
		apimg := image.NewPaletted(src.Bounds(), palette.Plan9)
		images = append(images, apimg)
	}

	return images
}
*/

// Sanctuary expanded maze parser
// parse exp maze string, looking for cmd lc, return parms if found
const (
	SE_FLOR		= 1		// load any gauntlet & G² floors & walls
	SE_WALL		= 2
	SE_G2		= 3		// gauntlet 2 mode - e.g. turn off G²
	SE_CFLOR	= 4		// custom floor unit from xb*.ed file wall & floor lines
	SE_CWAL		= 5		// data: sheet, r, (c is wally), xy size		wal: sheet, row, (xy size?)
	SE_COLRT	= 6		// color tiles under any
	SE_LETR		= 7		// draw a letter index to map_keymap (as gen hints) in color, as R,G,B, ind
	SE_MSG		= 8		// write a null term msg (up to 32 hex byts) onto maze in color as R,G,B, {MSG}, 00			test: 0800FFFF5A206973201E00
						// -- NOT really compatible with any other opcode due to possible embed action

	SE_MWAL		= 10	// master wall replace - must be placed under 0,0 - will be read first, data: sheet, r
	SE_MFLR		= 11	// master floor replace from string list - data: list entry
// once we get to super mazes, can these be a localized area?
	SE_TFLOR	= 12	// cust floor from xb*.ed, where floor pieces are tiled in sheet, data: line of xb, c floor tile col
	SE_NOFLOR	= 13	// display no: regular floor tile, or master override - same as SEOBJ_FLOORNODRAW, but ents can occupy cell
	SE_MWALRND	= 14	// master wall replace - under 0,0 as above, rnd loads of multiwall set, data: sheet, r, cnt rnd rows past - 2 rows on sheet 1 starting row 3 would be 0E000201
// more:
// 'item underlay' - put any item under, say a trap wall, or shootable wall so it appears when wall is gone, or even under a dragon, generator or item

)
// bytes for each cmd
var parms = []int{
	0, 					// item o, not used
	2, 2, 0, 1,
	2, 3, 4, 36,
	0, 2, 1, 2,
	0, 3, 0, 0,
}
var maxparm = 15
var secmd [64]int
var lastsp string		// dont need to reprocess same dats every parse call
var lastprl int			// last parm len
var xpar [64]int		// extra parms past 3... - parms[] can NOT exceed this array size!

func parser(sp string, lc int) (int, int, int) {
	r1, r2, r3 := -1,0,0
//fmt.Printf("parse %s\n ",sp)
	if lc > 0 {
	  if lastsp != sp {
		lastsp = sp
		for i := 0; i < 17; i++ { secmd[i] = 0 }
					//	0	1	2	3	4	5	6	7	8	9	10	11	12	13	14	15	16	17	18	19	20	21	22	23	24	25	26	27	28	29	30	31	32	33	34	35	36
		fmt.Sscanf(sp,"%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X",
					&secmd[0],&secmd[1],&secmd[2],&secmd[3],&secmd[4],&secmd[5],&secmd[6],&secmd[7],&secmd[8],&secmd[9],&secmd[10],&secmd[11],&secmd[12],&secmd[13],&secmd[14],&secmd[15],&secmd[16],
					&secmd[17],&secmd[18],&secmd[19],&secmd[20],&secmd[21],&secmd[22],&secmd[23],&secmd[24],&secmd[25],&secmd[26],&secmd[27],&secmd[28],&secmd[29],&secmd[30],&secmd[31],&secmd[32],
					&secmd[33],&secmd[34],&secmd[35],&secmd[36])
		for i := 0; i < 37; i++ { if secmd[i] != 0 { lastprl = i }}
	  }
		for i := 0; i <= lastprl; i++ {
			prc := parms[minint(secmd[i],maxparm)]
			if lc == secmd[i] {
				r1 =secmd[i+1]; r2 =secmd[i+2]; r3 =secmd[i+3]
				if prc == 0 { r1 = 1 }
				if prc > 3 { k:=0; for k = 0; k < (prc - 3); k++ { xpar[k] = secmd[k+4]
//if lc == 8 {fmt.Printf("%02X ",secmd[k+4])}
			}}
//				fmt.Printf("c:%d p:%d - r1 %d r2 %d r3 %d\n ",lc,prc,r1,r2,r3)
//				break
			}
//fmt.Printf("parm %d of %d= %d\n ",secmd[i],parms[secmd[i]])
			i = i + prc
		}

	}

	return r1, r2, r3
}

// stamp array builder

func bld_star(lk int ) {

	_, _, sents := itemGetPNG("gfx/se_ents.16.png")			// sanct engine ent sheet
	gtopl := ""
//	gtopcol := false	// disable gen letter seperate colors
	psx, psy, azx, azy := -1,-1,0,0
	gsv := G1
	mask := 256 		// set some masks here
	cnt := 0			// animation frames
//	acnt := 0			// attack anim
var dyn [100]int

	switch lk {
	case SEOBJ_PUSHWAL:
		G1 = false
		arstamp[lk] = itemGetStamp("pushwall")
	case G1OBJ_KEY:
		arstamp[lk] = itemGetStamp("key")

	case SEOBJ_POWER_REPULSE:
		G1 = false
		arstamp[lk] = itemGetStamp("repulse")
	case SEOBJ_POWER_REFLECT:
		G1 = false
		arstamp[lk] = itemGetStamp("reflect")
	case SEOBJ_POWER_TRANSPORT:
		G1 = false
		arstamp[lk] = itemGetStamp("transportability")
	case SEOBJ_POWER_SUPERSHOT:
		G1 = false
		arstamp[lk] = itemGetStamp("supershot")
	case SEOBJ_POWER_INVULN:
		G1 = false
		arstamp[lk] = itemGetStamp("invuln")

	case G1OBJ_PLAYERSTART:
		arstamp[lk] = itemGetStamp("plusg1")
//		if G2 { arstamp[lk] = itemGetStamp("plus") }
	case G1OBJ_EXIT:
		arstamp[lk] = itemGetStamp("exitg1")
//		if G2 { arstamp[lk] = itemGetStamp("exit") }
//fmt.Printf("g1exit %d x %d, vc: %d x %d\n ",x,y,vcx, vcy)
	case G1OBJ_EXIT4:
		arstamp[lk] = itemGetStamp("exit4")
	case SEOBJ_EXIT6:
		G1 = false
		arstamp[lk] = itemGetStamp("exit6")
	case G1OBJ_EXIT8:
		arstamp[lk] = itemGetStamp("exit8")

	case G1OBJ_MONST_GHOST1:
		arstamp[lk] = itemGetStamp("ghost1")
		cnt = 32
		dyn = [100]int{
			2048, 2057, 2066, 2075,		// D
			2084, 2093, 2102, 2111,		// DR
			2120, 2129, 2138, 2147,		// R
			2156, 2165, 2174, 2183,		// UR
			2192, 2201, 2210, 2219,		// U
			2228, 2237, 2246, 2255,		// UL
			2264, 2273, 2282, 2291,		// L
			2304, 2313, 2322, 2331,		// DL
			-1}
		arstamp[lk].awalk = [12]int{4,0,4,8,12,16,20,24,28,-1}

	case G1OBJ_MONST_GHOST2:
		arstamp[lk] = itemGetStamp("ghost2")
		cnt = 32
		dyn = [100]int{
			2048, 2057, 2066, 2075,		// D
			2084, 2093, 2102, 2111,		// DR
			2120, 2129, 2138, 2147,		// R
			2156, 2165, 2174, 2183,		// UR
			2192, 2201, 2210, 2219,		// U
			2228, 2237, 2246, 2255,		// UL
			2264, 2273, 2282, 2291,		// L
			2304, 2313, 2322, 2331,		// DL
			-1}
		arstamp[lk].awalk = [12]int{4,0,4,8,12,16,20,24,28,-1}
	case SEOBJ_G2GHOST:
		G1 = false; fallthrough
	case G1OBJ_MONST_GHOST3:
		arstamp[lk] = itemGetStamp("ghost")
		cnt = 32
		dyn = [100]int{
			2048, 2057, 2066, 2075,		// D
			2084, 2093, 2102, 2111,		// DR
			2120, 2129, 2138, 2147,		// R
			2156, 2165, 2174, 2183,		// UR
			2192, 2201, 2210, 2219,		// U
			2228, 2237, 2246, 2255,		// UL
			2264, 2273, 2282, 2291,		// L
			2304, 2313, 2322, 2331,		// DL
			-1}
		arstamp[lk].awalk = [12]int{4,0,4,8,12,16,20,24,28,-1}
	case G1OBJ_MONST_GRUNT1:
		arstamp[lk] = itemGetStamp("grunt1")
		cnt = 40
		dyn = [100]int{
			2529, 2538, 2547,	// D
			2560, 2569, 2578,	// DR
			2587, 2596, 2605,	// R			// 1 extar walk cyc going directly L & R 	- 2614, 2731
			2623, 2632, 2641,	// UR
			2650, 2659, 2668,	// U
			2677, 2686, 2695,	// UL
			2704, 2713, 2722,	// L
			2740, 2749, 2758,	// DL
			2767, 2776,		// D
			2785, 2794,		// DR
			2803, 2816,		// R
			2825, 2834,		// UR
			2843, 2852,		// U
			2861, 2870,		// UL
			2879, 2888,		// L
			2897, 2906,		// DL
			-1}
		arstamp[lk].awalk = [12]int{3,0,3,6,9,12,15,18,21,-1}
		arstamp[lk].amel = [12]int{2,24,26,28,30,32,34,36,38,-1}
	case G1OBJ_MONST_GRUNT2:
		arstamp[lk] = itemGetStamp("grunt2")
		cnt = 40
		dyn = [100]int{
			2529, 2538, 2547,	// D
			2560, 2569, 2578,	// DR
			2587, 2596, 2605,	// R
			2623, 2632, 2641,	// UR
			2650, 2659, 2668,	// U
			2677, 2686, 2695,	// UL
			2704, 2713, 2722,	// L
			2740, 2749, 2758,	// DL
			2767, 2776,		// D
			2785, 2794,		// DR
			2803, 2816,		// R
			2825, 2834,		// UR
			2843, 2852,		// U
			2861, 2870,		// UL
			2879, 2888,		// L
			2897, 2906,		// DL
			-1}
		arstamp[lk].awalk = [12]int{3,0,3,6,9,12,15,18,21,-1}
		arstamp[lk].amel = [12]int{2,24,26,28,30,32,34,36,38,-1}
	case SEOBJ_G2GRUNT:		// grunts & ghosts go D, DR, R, UR, U, UL, L, DL walk, & grunt attack
		G1 = false; fallthrough
	case G1OBJ_MONST_GRUNT3:
		arstamp[lk] = itemGetStamp("grunt")
		cnt = 40
		dyn = [100]int{
			2529, 2538, 2547,	// D
			2560, 2569, 2578,	// DR
			2587, 2596, 2605,	// R
			2623, 2632, 2641,	// UR
			2650, 2659, 2668,	// U
			2677, 2686, 2695,	// UL
			2704, 2713, 2722,	// L
			2740, 2749, 2758,	// DL
			2767, 2776,		// D
			2785, 2794,		// DR
			2803, 2816,		// R
			2825, 2834,		// UR
			2843, 2852,		// U
			2861, 2870,		// UL
			2879, 2888,		// L
			2897, 2906,		// DL
			-1}
		arstamp[lk].awalk = [12]int{3,0,3,6,9,12,15,18,21,-1}
		arstamp[lk].amel = [12]int{2,24,26,28,30,32,34,36,38,-1}
	case G1OBJ_MONST_DEMON1:
		arstamp[lk] = itemGetStamp("demon1")
		cnt = 64
		dyn = [100]int{
			6189, 6198, 6207, 6216, 6225, 		// D
			6279, 6288, 6297, 6306, 6315, 		// DR
			6369, 6378, 6387, 6400, 6409, 		// R
			6463, 6472, 6481, 6490, 6499, 		// UR
			6508, 6517, 6526, 6535, 6544, 		// U
			6418, 6427, 6436, 6445, 6454, 		// UL
			6324, 6333, 6342, 6351, 6360, 		// L
			6234, 6243, 6252, 6261, 6270, 		// DL
			6553, 6571, 6589, 6607, 6616, 6598, 6580, 6562,
			6625, 6634,
			6665, 6674,
			6683, 6692,
			6737, 6746,
			6755, 6764,
			6719, 6728,
			6701, 6710,
			6643, 6656,
			-1}
		arstamp[lk].awalk = [12]int{5,0,5,10,15,20,25,30,35,-1}
		arstamp[lk].ashot = [12]int{2,48,50,52,54,56,58,60,62,-1}
/*
6189, 6198, 6207, 6216, 6225, 		// D
6234, 6243, 6252, 6261, 6270, 		// DL
6279, 6288, 6297, 6306, 6315, 		// DR
6324, 6333, 6342, 6351, 6360, 		// L
6369, 6378, 6387, 6400, 6409, 		// R
6418, 6427, 6436, 6445, 6454, 		// UL
6463, 6472, 6481, 6490, 6499, 		// UR
6508, 6517, 6526, 6535, 6544, 		// U
6553, 6562, 6571, 6580, 6589, 6598, 6607, 6616, 	//? - D DL DR L R UL UR U
6625, 6634, D
6643, 6656, DL
6665, 6674, DR
6683, 6692, R
6701, 6710, L
6719, 6728, UL
6737, 6746, UR
6755, 6764, U
*/
	case G1OBJ_MONST_DEMON2:		// demons go D, DL, DR L, R, UL, UR, U
		arstamp[lk] = itemGetStamp("demon2")
		cnt = 64
		dyn = [100]int{
			6189, 6198, 6207, 6216, 6225, 		// D
			6279, 6288, 6297, 6306, 6315, 		// DR
			6369, 6378, 6387, 6400, 6409, 		// R
			6463, 6472, 6481, 6490, 6499, 		// UR
			6508, 6517, 6526, 6535, 6544, 		// U
			6418, 6427, 6436, 6445, 6454, 		// UL
			6324, 6333, 6342, 6351, 6360, 		// L
			6234, 6243, 6252, 6261, 6270, 		// DL
			6553, 6571, 6589, 6607, 6616, 6598, 6580, 6562,
			6625, 6634,
			6665, 6674,
			6683, 6692,
			6737, 6746,
			6755, 6764,
			6719, 6728,
			6701, 6710,
			6643, 6656,
			-1}
		arstamp[lk].awalk = [12]int{5,0,5,10,15,20,25,30,35,-1}
		arstamp[lk].ashot = [12]int{2,48,50,52,54,56,58,60,62,-1}
	case SEOBJ_G2DEMON:
		G1 = false; fallthrough
	case G1OBJ_MONST_DEMON3:
		arstamp[lk] = itemGetStamp("demon")
		cnt = 64
		dyn = [100]int{
			6189, 6198, 6207, 6216, 6225, 		// D
			6279, 6288, 6297, 6306, 6315, 		// DR
			6369, 6378, 6387, 6400, 6409, 		// R
			6463, 6472, 6481, 6490, 6499, 		// UR
			6508, 6517, 6526, 6535, 6544, 		// U
			6418, 6427, 6436, 6445, 6454, 		// UL
			6324, 6333, 6342, 6351, 6360, 		// L
			6234, 6243, 6252, 6261, 6270, 		// DL
			6553, 6571, 6589, 6607, 6616, 6598, 6580, 6562,
			6625, 6634,
			6665, 6674,
			6683, 6692,
			6737, 6746,
			6755, 6764,
			6719, 6728,
			6701, 6710,
			6643, 6656,
			-1}
		arstamp[lk].awalk = [12]int{5,0,5,10,15,20,25,30,35,-1}
		arstamp[lk].ashot = [12]int{2,48,50,52,54,56,58,60,62,-1}
	case G1OBJ_MONST_LOBBER1:
		arstamp[lk] = itemGetStamp("lobber1")
		cnt = 49
		dyn = [100]int{
6993, 6999, 7005, 7011,		// D
7017, // hand emptys
7023, 7029, 7035, 7041, 	// DR
7047,
7053, 7059, 7065, 7071, 	// R
7077,
7083, 7089, 7095, 7101, 	// UR
7111,
7117, 7123, 7129, 7135, 	// U
7141,
7147, 7153, 7159, 7168, 	// UL
7174,
7180, 7186, 7192, 7198, 	// L
7204,
7210, 7216, 7222, 7228, 	// DL
7234,
7240, 7244, 7248, 7252, 7256, // shot sm, up, big, dn, sm
7252, 7248, 7244, 7240,
			-1}
		arstamp[lk].awalk = [12]int{4,0,5,10,15,20,25,30,35,-1}
		arstamp[lk].ashot = [12]int{2,3,8,13,18,23,28,33,38,-1}
	case G1OBJ_MONST_LOBBER2:
		arstamp[lk] = itemGetStamp("lobber2")
		cnt = 49
		dyn = [100]int{
6993, 6999, 7005, 7011,		// D
7017, // hand emptys
7023, 7029, 7035, 7041, 	// DR
7047,
7053, 7059, 7065, 7071, 	// R
7077,
7083, 7089, 7095, 7101, 	// UR
7111,
7117, 7123, 7129, 7135, 	// U
7141,
7147, 7153, 7159, 7168, 	// UL
7174,
7180, 7186, 7192, 7198, 	// L
7204,
7210, 7216, 7222, 7228, 	// DL
7234,
7240, 7244, 7248, 7252, 7256, // shot sm, up, big, dn, sm
7252, 7248, 7244, 7240,
			-1}
		arstamp[lk].awalk = [12]int{4,0,5,10,15,20,25,30,35,-1}
		arstamp[lk].ashot = [12]int{2,3,8,13,18,23,28,33,38,-1}
	case SEOBJ_G2LOBER:
		G1 = false; fallthrough
	case G1OBJ_MONST_LOBBER3:
		arstamp[lk] = itemGetStamp("lobber")
		cnt = 49
		dyn = [100]int{
6993, 6999, 7005, 7011,		// D
7017, // hand emptys
7023, 7029, 7035, 7041, 	// DR
7047,
7053, 7059, 7065, 7071, 	// R
7077,
7083, 7089, 7095, 7101, 	// UR
7111,
7117, 7123, 7129, 7135, 	// U
7141,
7147, 7153, 7159, 7168, 	// UL
7174,
7180, 7186, 7192, 7198, 	// L
7204,
7210, 7216, 7222, 7228, 	// DL
7234,
7240, 7244, 7248, 7252, 7256, // shot sm, up, big, dn, sm
7252, 7248, 7244, 7240,
			-1}
		arstamp[lk].awalk = [12]int{4,0,5,10,15,20,25,30,35,-1}
		arstamp[lk].ashot = [12]int{2,3,8,13,18,23,28,33,38,-1}
	case G1OBJ_MONST_SORC1:
		arstamp[lk] = itemGetStamp("sorcerer1")
		cnt = 40
		dyn = [100]int{
			5026, 5035, 5044, 	//. walks D
			5219, 5228, 5237, 	// DR
			5192, 5201, 5210, 	// R
			5165, 5174, 5183, 	// UR
			5138, 5147, 5156, 	// U
			5107, 5120, 5129, 	// UL
			5080, 5089, 5098, 	// L
			5053, 5062, 5071, 	// DL
			5246, 5255,		// melee D
			5376, 5385, 	// DR
			5354, 5363, 	// R
			5336, 5345, 	// UR
			5318, 5327, 	// U
			5300, 5309, 	// UL
			5282, 5291, 	// L
			5264, 5273, 	// DL
			-1}
		arstamp[lk].awalk = [12]int{3,0,3,6,9,12,15,18,21,-1}
		arstamp[lk].amel = [12]int{2,24,26,28,30,32,34,36,38,-1}
/*
// also wizzerd - 48 frems
4954, 4963, 4972, 4981, 4990, 4999, 5008, 5017, 	// shots D, DL, L, UL, U, UR, R, DR
5026, 5035, 5044, 	//. walks D
5053, 5062, 5071, 	// DL
5080, 5089, 5098, 	// L
5107, 5120, 5129, 	// UL
5138, 5147, 5156, 	// U
5165, 5174, 5183, 	// UR
5192, 5201, 5210, 	// R
5219, 5228, 5237, 	// DR
5246, 5255,		// melee D
5264, 5273, 	// DL
5282, 5291, 	// L
5300, 5309, 	// UL
5318, 5327, 	// U
5336, 5345, 	// UR
5354, 5363, 	// R
5376, 5385, 	// DR
*/
	case G1OBJ_MONST_SORC2:
		arstamp[lk] = itemGetStamp("sorcerer2")
		cnt = 40
		dyn = [100]int{
			5026, 5035, 5044, 	//. walks D
			5219, 5228, 5237, 	// DR
			5192, 5201, 5210, 	// R
			5165, 5174, 5183, 	// UR
			5138, 5147, 5156, 	// U
			5107, 5120, 5129, 	// UL
			5080, 5089, 5098, 	// L
			5053, 5062, 5071, 	// DL
			5246, 5255,		// melee D
			5376, 5385, 	// DR
			5354, 5363, 	// R
			5336, 5345, 	// UR
			5318, 5327, 	// U
			5300, 5309, 	// UL
			5282, 5291, 	// L
			5264, 5273, 	// DL
			-1}
		arstamp[lk].awalk = [12]int{3,0,3,6,9,12,15,18,21,-1}
		arstamp[lk].amel = [12]int{2,24,26,28,30,32,34,36,38,-1}
	case SEOBJ_G2SORC:
		G1 = false; fallthrough
	case G1OBJ_MONST_SORC3:
		arstamp[lk] = itemGetStamp("sorcerer")
		cnt = 40
		dyn = [100]int{
			5026, 5035, 5044, 	//. walks D
			5219, 5228, 5237, 	// DR
			5192, 5201, 5210, 	// R
			5165, 5174, 5183, 	// UR
			5138, 5147, 5156, 	// U
			5107, 5120, 5129, 	// UL
			5080, 5089, 5098, 	// L
			5053, 5062, 5071, 	// DL
			5246, 5255,		// melee D
			5376, 5385, 	// DR
			5354, 5363, 	// R
			5336, 5345, 	// UR
			5318, 5327, 	// U
			5300, 5309, 	// UL
			5282, 5291, 	// L
			5264, 5273, 	// DL
			-1}
		arstamp[lk].awalk = [12]int{3,0,3,6,9,12,15,18,21,-1}
		arstamp[lk].amel = [12]int{2,24,26,28,30,32,34,36,38,-1}
	case SEOBJ_G2AUXGR:
		G1 = false
		arstamp[lk] = itemGetStamp("auxgrunt")
		cnt = 40
		dyn = [100]int{
			2529, 2538, 2547,	// D
			2560, 2569, 2578,	// DR
			2587, 2596, 2605,	// R
			2623, 2632, 2641,	// UR
			2650, 2659, 2668,	// U
			2677, 2686, 2695,	// UL
			2704, 2713, 2722,	// L
			2740, 2749, 2758,	// DL
			2767, 2776,		// D
			2785, 2794,		// DR
			2803, 2816,		// R
			2825, 2834,		// UR
			2843, 2852,		// U
			2861, 2870,		// UL
			2879, 2888,		// L
			2897, 2906,		// DL
			-1}
		arstamp[lk].awalk = [12]int{3,0,3,6,9,12,15,18,21,-1}
		arstamp[lk].amel = [12]int{2,24,26,28,30,32,34,36,38,-1}

	case G1OBJ_MONST_DEATH:		// D DL L U DR R UR UL
		arstamp[lk] = itemGetStamp("death")
		cnt = 24
		dyn = [100]int{
			6773, 6782, 6791,		// D
			6881, 6890, 6899,		// DR
			6912, 6921, 6930,		// R
			6939, 6948, 6957,		// UR
			6854, 6863, 6872,		// U
			6966, 6975, 6984,		// UL
			6827, 6836, 6845,		// L
			6800, 6809, 6818,		// DL
			 -1}
		arstamp[lk].awalk = [12]int{3,0,3,6,9,12,15,18,21,-1}
/* esperiment - 4 frame deth		-- doesnt look right, tho not yet combine with move
		cnt = 32
		dyn = [100]int{
			6773, 6782, 6791, 6782,		// D
			6881, 6890, 6899, 6890,		// DR
			6912, 6921, 6930, 6921,		// R
			6939, 6948, 6957, 6948,		// UR
			6854, 6863, 6872, 6863,		// U
			6966, 6975, 6984, 6975,		// UL
			6827, 6836, 6845, 6836,		// L
			6800, 6809, 6818, 6809,		// DL
			 -1}
		arstamp[lk].awalk = [12]int{4,0,4,8,12,16,20,24,28,-1}
*/
	case SEOBJ_G2ACID:
		G1 = false
		arstamp[lk] = itemGetStamp("acid")
		cnt = 52
		dyn = [100]int{
			8960, 8969, 8978, 8987, 8996, 9005, 9014, 9023, 	// vert
			9158, 9167, 9176, 9185, 9194, 9203,			//  diag \
			9032, 9041, 9050, 9059, 9068, 9077, 9086, 9095, 	// horiz
			9104, 9113, 9122, 9131, 9140, 9149, 		//  diag /			// acid expanded from base frameset
			9014, 9005, 8996, 8987, 8978, 8960, 		// U				// to make more dirs look diff
			9194, 9185, 9176, 9167, 9158, 9203,			// UL
			9068, 9059, 9050, 9041, 9095, 9086,			// L
			9122, 9113, 9149, 9140, 9131, 9104,			// DL
			 -1}
		arstamp[lk].awalk = [12]int{6,1,8,15,22,28,34,40,46}		// acid could have rnd 4-6 loop, & rnd starts
	case SEOBJ_G2SUPSORC:
		G1 = false
		arstamp[lk] = itemGetStamp("supersorc")
		cnt = 48
		dyn = [100]int{
			5026, 5035, 5044, 	//. walks D
			5219, 5228, 5237, 	// DR
			5192, 5201, 5210, 	// R
			5165, 5174, 5183, 	// UR
			5138, 5147, 5156, 	// U
			5107, 5120, 5129, 	// UL
			5080, 5089, 5098, 	// L
			5053, 5062, 5071, 	// DL
			5246, 5255,		// melee D
			5376, 5385, 	// DR
			5354, 5363, 	// R
			5336, 5345, 	// UR
			5318, 5327, 	// U
			5300, 5309, 	// UL
			5282, 5291, 	// L
			5264, 5273, 	// DL
			4954, 5017, 5008, 4999, 4990, 4981, 4972, 4963, 	// shots D, DR, R, UR, U, UL, L, DL
			-1}
		arstamp[lk].awalk = [12]int{3,0,3,6,9,12,15,18,21,-1}
		arstamp[lk].amel = [12]int{2,24,26,28,30,32,34,36,38,-1}
	case SEOBJ_G2IT:
		G1 = false
		arstamp[lk] = itemGetStamp("it")
		cnt = 15
		dyn = [100]int{
			9728, 9737, 9746, 9755, 9764, 9773, 9782, 9791, 9800, 9809, 9818, 9827, 9836, 9845, 9854, 9863,		// just a cycle
			-1}
		arstamp[lk].awalk = [12]int{-15,0}
	case SEOBJ_MONST_DRAGON:
		G1 = false
		arstamp[lk] = itemGetStamp("dragon")
		cnt = 40
		dyn = [100]int{								// drag fighting 4 dirs
			8704, 8720, 8736, 8752, 8768, 8784, 8800, 8816, 	// L		-> D
			8576, 8592, 8608, 8624, 8640, 8656, 8672, 8688, 	// R
			8832, 8848, 8864, 8880, 8896, 8912, 8928, 8944,		// D		-> U
			8448, 8464, 8480, 8496, 8512, 8528, 8544, 8560, 	// U		-> L
			9472, 9488, 9504, 9520, 	// waking? DL DR UR UL
			9536, 9552, 10048, 10064,	// sleeping D R U L ????
			-1}
		arstamp[lk].awalk = [12]int{8,0,0,8,8,16,16,24,24,-1}
	case G1OBJ_MONST_THIEF:
		arstamp[lk] = itemGetStamp("thief")
		cnt = 72
		dyn = [100]int{
			3562, 3571, 3584, 	// D
			3647, 3656, 3665, 	// DR
			3683, 3692, 3701, 	// R
			3737, 3746, 3755, 	// UR
			3593, 3602, 3611, 	// U
			3764, 3773, 3782,	// UL
			3710, 3719, 3728, 	// L
			3620, 3629, 3638, 	// DL
			3791, 3800, 3809, 	// glowy D
			3930, 3939, 3948, 	// DR
			3876, 3885, 3894, 	// R
			3984, 3993, 4002, 	// UR
			3849, 3858, 3867, 	// U
			3957, 3966, 3975, 	// UL
			3818, 3827, 3840, 	// L
			3903, 3912, 3921, 	// DL
			4011, 4020, 4029, 	// hop	D
			4141, 4150, 4159, 	// DR
			4065, 4074, 4083, 	// R
			4195, 4204, 4213,	// UR
			4038, 4047, 4056, 	// U
			4168, 4177, 4186, 	// UL
			4096, 4105, 4114, 	// L
			4123, 4132, 		// DL ?
			-1}
		arstamp[lk].awalk = [12]int{3,0,3,6,9,12,15,18,21,-1}
		arstamp[lk].amel = [12]int{3,24,27,30,33,36,39,42,45,-1}

	case SEOBJ_MONST_MUGGER:
		G1 = false
		arstamp[lk] = itemGetStamp("mugger")
		cnt = 42
		dyn = [100]int{
			9216, 9225, 9234, // D
			9243, 9252, 9261, // U
			9270, 9279, 9288, // DL
			9297, 9306, 9315, // DR
			9333, 9342, 9351, // R
			9360, 9369, 9378, // L
			9387, 9396, 9405, // UR
			9414, 9423, 9432, // UL
			9441, 9450, 9459, // attacks D
			9872, 9881, 9890, // L ?
			9899, 9908, 9917, // R ?
			9926, 9935, 9944, // U
			9953, 9962, 9971,// UR
			10080,10089,10098, // R
			-1}
		arstamp[lk].awalk = [12]int{3,0,3,6,9,12,15,18,21,-1}

	case SEOBJ_G2GN_GST1:
		G1 = false; fallthrough
	case G1OBJ_GEN_GHOST1:
		arstamp[lk] = itemGetStamp("ghostgen1")
	case SEOBJ_G2GN_GST2:
		G1 = false; fallthrough
	case G1OBJ_GEN_GHOST2:
		arstamp[lk] = itemGetStamp("ghostgen2")
	case SEOBJ_G2GN_GST3:
		G1 = false; fallthrough
	case G1OBJ_GEN_GHOST3:
		arstamp[lk] = itemGetStamp("ghostgen3")

// if a clear is done after, this SetRGB set bkg somehow
	case SEOBJ_G2GN_GR1:
		G1 = false; fallthrough
	case SEOBJ_G2GN_AUXGR1:
		G1 = false; fallthrough
	case G1OBJ_GEN_GRUNT1:
		gtopl = "G"
//		if gtopcol { gtop.SetRGB(0.65, 0.3, 0.1) }
		arstamp[lk] = itemGetStamp("generator1")
		arstamp[lk].theme  = color.NRGBA{0xff, 166, 77, 25}
	case SEOBJ_G2GN_DM1:
		G1 = false; fallthrough
	case G1OBJ_GEN_DEMON1:
		gtopl = "D"
//		if gtopcol { gtop.SetRGB(1, 0, 0) }
		arstamp[lk] = itemGetStamp("generator1")
		arstamp[lk].theme  = color.NRGBA{0xff, 255, 0, 0}
	case SEOBJ_G2GN_LB1:
		G1 = false; fallthrough
	case G1OBJ_GEN_LOBBER1:
		gtopl = "L"
//		if gtopcol { gtop.SetRGB(0.7, 0.5, 0.2) }
		arstamp[lk] = itemGetStamp("generator1")
		arstamp[lk].theme  = color.NRGBA{0xff, 179, 128, 52}
	case SEOBJ_G2GN_SORC1:
		G1 = false; fallthrough
	case G1OBJ_GEN_SORC1:
		gtopl = "S"
//		if gtopcol { gtop.SetRGB(0.37, 0.2, 0.7) }
//		arstamp[lk].theme  = color.NRGBA{0xff, 255 * 0.37, 255 * 0.21, 255 * 0.7}
		arstamp[lk] = itemGetStamp("generator1")
		arstamp[lk].theme  = color.NRGBA{0xff, 95, 52, 179}

	case SEOBJ_G2GN_GR2:
		G1 = false; fallthrough
	case SEOBJ_G2GN_AUXGR2:
		G1 = false; fallthrough
	case G1OBJ_GEN_GRUNT2:
		gtopl = "G"
//		if gtopcol { gtop.SetRGB(0.65, 0.3, 0.1) }
		arstamp[lk] = itemGetStamp("generator2")
		arstamp[lk].theme  = color.NRGBA{0xff, 166, 77, 25}
	case SEOBJ_G2GN_DM2:
		G1 = false; fallthrough
	case G1OBJ_GEN_DEMON2:
		gtopl = "D"
//		if gtopcol { gtop.SetRGB(1, 0, 0) }
		arstamp[lk] = itemGetStamp("generator2")
		arstamp[lk].theme  = color.NRGBA{0xff, 255, 0, 0}
	case SEOBJ_G2GN_LB2:
		G1 = false; fallthrough
	case G1OBJ_GEN_LOBBER2:
		gtopl = "L"
//		if gtopcol { gtop.SetRGB(0.7, 0.5, 0.2) }
		arstamp[lk] = itemGetStamp("generator2")
		arstamp[lk].theme  = color.NRGBA{0xff, 179, 128, 52}
	case SEOBJ_G2GN_SORC2:
		G1 = false; fallthrough
	case G1OBJ_GEN_SORC2:
		gtopl = "S"
//		if gtopcol { gtop.SetRGB(0.37, 0.2, 0.7) }
		arstamp[lk] = itemGetStamp("generator2")
		arstamp[lk].theme  = color.NRGBA{0xff, 95, 52, 179}

	case SEOBJ_G2GN_GR3:
		G1 = false; fallthrough
	case SEOBJ_G2GN_AUXGR3:
		G1 = false; fallthrough
	case G1OBJ_GEN_GRUNT3:
		gtopl = "G"
//		if gtopcol { gtop.SetRGB(0.65, 0.3, 0.1) }
		arstamp[lk] = itemGetStamp("generator3")
		arstamp[lk].theme  = color.NRGBA{0xff, 166, 77, 25}
	case SEOBJ_G2GN_DM3:
		G1 = false; fallthrough
	case G1OBJ_GEN_DEMON3:
		gtopl = "D"
//		if gtopcol { gtop.SetRGB(1, 0, 0) }
		arstamp[lk] = itemGetStamp("generator3")
		arstamp[lk].theme  = color.NRGBA{0xff, 255, 0, 0}
	case SEOBJ_G2GN_LB3:
		G1 = false; fallthrough
	case G1OBJ_GEN_LOBBER3:
		gtopl = "L"
//		if gtopcol { gtop.SetRGB(0.7, 0.5, 0.2) }
		arstamp[lk] = itemGetStamp("generator3")
		arstamp[lk].theme  = color.NRGBA{0xff, 179, 128, 52}
	case SEOBJ_G2GN_SORC3:
		G1 = false; fallthrough
	case G1OBJ_GEN_SORC3:
		gtopl = "S"
//		if gtopcol { gtop.SetRGB(0.37, 0.2, 0.7) }
		arstamp[lk] = itemGetStamp("generator3")
		arstamp[lk].theme  = color.NRGBA{0xff, 95, 52, 179}

	case G1OBJ_TREASURE:
		arstamp[lk] = itemGetStamp("treasure")
		cnt = 3
		dyn = [100]int{
			2439,2448,2457,
			-1}
	case SEOBJ_TREASURE_LOCKED:
		G1 = false
		arstamp[lk] = itemGetStamp("treasurelocked")

	case G1OBJ_TREASURE_BAG:
		arstamp[lk] = itemGetStamp("goldbag")
	case G1OBJ_FOOD_DESTRUCTABLE:
		arstamp[lk] = itemGetStamp("food")
	case G1OBJ_FOOD_INVULN:
		arstamp[lk] = itemGetStamp(foods[rand.Intn(3)])
	case G1OBJ_POT_DESTRUCTABLE:
		arstamp[lk] = itemGetStamp("potion")
	case G1OBJ_POT_INVULN:
		arstamp[lk] = itemGetStamp("ipotion")
	case G1OBJ_INVISIBL:
		arstamp[lk] = itemGetStamp("invis")
// specials added after convert to se id'ed them on maze 115, score table block
	case G1OBJ_X_SPEED:
		arstamp[lk] = itemGetStamp("speedpotion")
	case G1OBJ_X_SHOTPW:
		arstamp[lk] = itemGetStamp("shotpowerpotion")
	case G1OBJ_X_SHTSPD:
		arstamp[lk] = itemGetStamp("shotspeedpotion")
	case G1OBJ_X_ARMOR:
		arstamp[lk] = itemGetStamp("shieldpotion")
	case G1OBJ_X_FIGHT:
		arstamp[lk] = itemGetStamp("fightpotion")
	case G1OBJ_X_MAGIC:
		arstamp[lk] = itemGetStamp("magicpotion")
	case G1OBJ_TRANSPORTER:
		arstamp[lk] = itemGetStamp("tportg1")
//		if G2 { arstamp[lk] = itemGetStamp("tport") }
// player animation
	case G1OBJ_WARRIOR:
		arstamp[lk] = itemGetStamp("warrior")
		cnt = 71
		dyn = [100]int{
			2915, 2924, 2933,	// D
			3108, 3117, 3126,	// DR
			3081, 3090, 3099,	// R
			3050, 3059, 3072,	// UR
			3023, 3032, 3041,	// U
			2996, 3005, 3014,	// UL
			2969, 2978, 2987,	// L
			2942, 2951, 2960,	// DL
			3135, 3144,		// throw
			3261, 3270,
			3243, 3252,
			3225, 3234,
			3207, 3216,
			3189, 3198,
			3171, 3180,
			3153, 3162,
			3279, 3288, 3297,	// melee		these may be 4 parts like valk
			3526, 3535, 3544,
			3499, 3508, 3517,
			3472, 3481, 3490,
			3445, 3454, 3463,
			3418, 3427, 3436,
			3391, 3400, 3409,
			3364, 3373, 3382,
			3337, 3346, 3355,	// unk ?
			3306, 3315, 3328,
			3553,		// unk, unless drain
			4231, 4240, 4249, 4258, 4267, 4276, 4285,		// warr drain
			-1}
		arstamp[lk].awalk = [12]int{3,0,3,6,9,12,15,18,21,-1}
		arstamp[lk].ashot = [12]int{2,24,26,28,30,32,34,36,38,-1}
		arstamp[lk].amel = [12]int{1,40,43,46,49,52,55,58,61,-1}
	case G1OBJ_WIZARD:
		arstamp[lk] = itemGetStamp("wizard")
		cnt = 48
		dyn = [100]int{
			5026, 5035, 5044, 	//. walks D
			5219, 5228, 5237, 	// DR
			5192, 5201, 5210, 	// R
			5165, 5174, 5183, 	// UR
			5138, 5147, 5156, 	// U
			5107, 5120, 5129, 	// UL
			5080, 5089, 5098, 	// L
			5053, 5062, 5071, 	// DL
			5246, 5255,		// melee D
			5376, 5385, 	// DR
			5354, 5363, 	// R
			5336, 5345, 	// UR
			5318, 5327, 	// U
			5300, 5309, 	// UL
			5282, 5291, 	// L
			5264, 5273, 	// DL
			4954, 5017, 5008, 4999, 4990, 4981, 4972, 4963, 	// shots D, DR, R, UR, U, UL, L, DL
			6122, 6131, 6140, 6144, 6153, 6162, 6171, 6180,		// wiz drain
			-1}
		arstamp[lk].awalk = [12]int{3,0,3,6,9,12,15,18,21,-1}
		arstamp[lk].amel = [12]int{2,24,26,28,30,32,34,36,38,-1}
		arstamp[lk].ashot = [12]int{1,40,41,42,43,44,45,46,47,-1}
	case G1OBJ_VALKYRIE:
		arstamp[lk] = itemGetStamp("valkyrie")
		cnt = 72
		dyn = [100]int{

			4370, 4379, 4388,		// R walk
			4397, 4406, 4415,		// DR
			4424, 4433, 4442,		// D
			4451, 4460, 4469,		// DL
			4478, 4487, 4496,		// L
			4505, 4514, 4523,		// UL
			4532, 4541, 4550,		// U
			4559, 4568, 4577,		// UR

			4586, 4595, 4608, 4617,			// melee
			4626, 4635, 4644, 4653,
			4662, 4671, 4680, 4689,
			4698, 4707, 4716, 4725,
			4734, 4743, 4752, 4761,
			4770, 4779, 4788, 4797,
			4806, 4815, 4824, 4833,
			4842, 4851, 4864, 4873,

			4294, 4303, 4312, 4321, 4330, 4339, 4352, 4361, 	// melee ?		R, DR, D, DL, L, UL, U, UR
			4882, 4891, 4900, 4909, 4918, 4927, 4936, 4945,		// thrown 	R, DR, D, DL, L, UL, U, UR
			6059, 6068, 6077, 6086, 6095, 6104, 6113, 			// valk drain
			-1}
		arstamp[lk].awalk = [12]int{3,6,3,0,21,18,15,12,9,-1}
		arstamp[lk].ashot = [12]int{1,66,65,64,71,70,69,68,67,-1}
		arstamp[lk].amel = [12]int{4,32,28,24,52,48,44,40,36,-1}
	case G1OBJ_ELF:
		arstamp[lk] = itemGetStamp("elf")
		cnt = 71
		dyn = [100]int{
			5538, 5547, 5556,	// L
			5565, 5574, 5583,	// UL
			5592, 5601, 5610,	// U
			5619, 5632, 5641,	// UR
			5650, 5659, 5668,	// R
			5677, 5686, 5695,	// DR
			5704, 5713, 5722,	// D
			5731, 5740, 5749,	// DL

			5394, 5466,
			5403, 5475,
			5412, 5484,
			5421, 5493,
			5430, 5502,
			5439, 5511,
			5448, 5520,
			5457, 5529,	// shots?

			5758, 5767, 5776,
			5785, 5794, 5803,
			5812, 5821, 5830,
			5839, 5848, 5857,
			5866, 5875, 5906,
			5915, 5924, 5933,
			5942, 5951, 5960,
			5969, 5978, 5987, 		// melee?

			5996, 6005, 6014, 6023, 6032, 6041, 6050, // elf drain
			-1}
		arstamp[lk].awalk = [12]int{3,18,15,12,9,6,3,0,21,-1}
		arstamp[lk].ashot = [12]int{2,36,34,32,30,28,26,24,38,-1}
		arstamp[lk].amel = [12]int{1,58,55,52,49,46,43,40,61,-1}

// SE expand
	case SEOBJ_SE_ANKH:
		psx, psy = 21, 11
		mask = NOPOT
	case SEOBJ_FIRE_STICK:
		psx, psy = 32, 26
		mask = 256 | ANIM
		cnt = 4
	case SEOBJ_G2_POISPOT:
		arstamp[lk] = itemGetStamp("ppotion")
		psx, psy = 8, 11
		mask = NOPOT
	case SEOBJ_G2_POISFUD:
		arstamp[lk] = itemGetStamp("pfood")
		psx, psy = 1, 11
		mask = NOFUD
	case SEOBJ_G2_QFUD:
		arstamp[lk] = itemGetStamp("mfood")
		psx, psy = 2, 11
		mask = NOFUD
	case SEOBJ_KEYRING:
		psx, psy = 28, 10
		mask = NODOR

	case SEOBJ_MAPPYBDG:
		psx, psy = 32, 22
	case SEOBJ_MAPPYGORO:
		psx, psy = 34, 22

	case SEOBJ_MAPPYARAD:		// 25, 21
		psx, psy, azx = 24, 20, 16
	case SEOBJ_MAPPYATV:		// 27, 21
		psx, psy, azx = 26, 20, 16
	case SEOBJ_MAPPYAPC:		// 29, 21
		psx, psy, azx = 28, 20, 16
	case SEOBJ_MAPPYAART:		// 31, 21
		psx, psy, azx = 30, 20, 16
	case SEOBJ_MAPPYASAF:		// 33, 21
		psx, psy, azx = 32, 20, 16

	case SEOBJ_MAPPYRAD:		// 25, 22
		psx, psy, azx = 24, 21, 16
	case SEOBJ_MAPPYTV:		// 27, 22
		psx, psy, azx = 26, 21, 16
	case SEOBJ_MAPPYPC:		// 29, 22
		psx, psy, azx = 28, 21, 16
	case SEOBJ_MAPPYART:		// 31, 22
		psx, psy, azx = 30, 21, 16
	case SEOBJ_MAPPYSAF:		// 33, 22
		psx, psy, azx = 32, 21, 16

	case SEOBJ_MAPPYBELL:		// 35, 21
		psx, psy = 34, 20
		mask = 256 | ANIM
		cnt = 2
	case SEOBJ_MAPPYBAL:		// 35, 22
		psx, psy = 34, 21
		mask = 256 | ANIM
		cnt = 2

	case SEOBJ_DETHGEN3:		// 34, 8
		gtopl = "D"
		mask = NOGEN
//		gtop.SetRGB(0, 0, 0)
		psx, psy = 34, 8
	case SEOBJ_DETHGEN2:		// 35, 8
		gtopl = "D"
		mask = NOGEN
//		gtop.SetRGB(0, 0, 0)
		psx, psy = 33, 8
	case SEOBJ_DETHGEN1:		// 36, 8
		gtopl = "D"
		mask = NOGEN
//		gtop.SetRGB(0, 0, 0)
		psx, psy = 32, 8
	case SEOBJ_FLOORNUL:
		psx, psy = 34, 10
		mask = 0
// animated floor
	case SEOBJ_WATER_POOL:
		psx, psy = 0, 26
		mask = ANIM | NOFLOOR
		cnt = 4
	case SEOBJ_WATER_TOP:
		psx, psy = 4, 26
		mask = ANIM | NOFLOOR
		cnt = 4
	case SEOBJ_WATER_RT:
		psx, psy = 12, 26
		mask = ANIM | NOFLOOR
		cnt = 4
	case SEOBJ_WATER_COR:
		psx, psy = 8, 26
		mask = ANIM | NOFLOOR
		cnt = 4
	case SEOBJ_SLIME_POOL:
		psx, psy = 16, 27
		mask = ANIM | NOFLOOR
		cnt = 3
	case SEOBJ_SLIME_TOP:
		psx, psy = 19, 27
		mask = ANIM | NOFLOOR
		cnt = 3
	case SEOBJ_SLIME_RT:
		psx, psy = 25, 27
		mask = ANIM | NOFLOOR
		cnt = 3
	case SEOBJ_SLIME_COR:
		psx, psy = 22, 27
		mask = ANIM | NOFLOOR
		cnt = 3
	case SEOBJ_LAVA_POOL:
		psx, psy = 16, 26
		mask = ANIM | NOFLOOR
		cnt = 4
	case SEOBJ_LAVA_TOP:
		psx, psy = 20, 26
		mask = ANIM | NOFLOOR
		cnt = 4
	case SEOBJ_LAVA_RT:
		psx, psy = 28, 26
		mask = ANIM | NOFLOOR
		cnt = 4
	case SEOBJ_LAVA_COR:
		psx, psy = 24, 26
		mask = ANIM | NOFLOOR
		cnt = 4
	case SEOBJ_PULS_FLOR:
		psx, psy = 28, 27
		mask = ANIM | NOFLOOR
		cnt = 8
	default:

			if opts.Verbose && false { fmt.Printf("G¹ WARNING: Unhandled obj id 0x%02x\n", lk) }
	}
//fmt.Printf("star ld %d, %v\n",lk,arstamp[lk])
	if arstamp[lk] == nil {
		arstamp[lk] = itemGetStamp("key")
		arstamp[lk].pnum = -1				// failed assign, no use
	}
	arstamp[lk].animtm = cnt				// frames to animate (in a single set) if any
	if arstamp[lk].pnum >= 0 {
//fmt.Printf("star ld %d, 1\n",lk)
		v := arstamp[lk].width * 8
		arstamp[lk].mimg = blankimage(v,v)
		writestamptoimage(G1,arstamp[lk].mimg, arstamp[lk], 0, 0)
		arstamp[lk].altimg = arstamp[lk].mimg
		if cnt > 0 {
			for i := 0; i < cnt; i++ {
				arstamp[lk].anim = append(arstamp[lk].anim,nil)
				arstamp[lk].anim[i] = blankimage(v,v)
				stc := arstamp[lk]
				if dyn[i] > 0 {
					ld := len(stc.data)
					if dyn[i] > 7239 && dyn[i] < 7269 { ld = 4; stc.ptype = "base"; stc.pnum = 1; stc.width = 2 } // lobber rocks,  shots will prob need rebase
					stc.numbers = tilerange(dyn[i],ld)
					fillstamp(stc)
//fmt.Printf("anim stamp %d ad: %v\n",lk,stc.numbers)
					writestamptoimage(G1,arstamp[lk].anim[i], stc, 0, 0)
//fl := fmt.Sprintf("tst%02d.png",i)
//savetopng(fl, arstamp[lk].anim[i])
				}
		}}
	}
	if psx >= 0 && psy >= 0 {			// supply alt img
		parimg = sents
		arstamp[lk].altimg = blankimage(16+azx,16+azy)
		writepngtoimage(arstamp[lk].altimg,16,16,azx,azy,psx,psy,0,0,0)
		if arstamp[lk].pnum < 0 { arstamp[lk].mimg = arstamp[lk].altimg; arstamp[lk].pnum = -7 }	// no main img, use alt
		arstamp[lk].mask = mask
		if cnt > 0 {		// animation frames
//fmt.Printf("bld anim %d, c%d w%d\n",lk,cnt,16)
			arstamp[lk].anim = append(arstamp[lk].anim,nil)
			arstamp[lk].anim[0] = arstamp[lk].mimg		// main img is always 1st frame ?
			for i := 1; i < cnt; i++ {
//fmt.Printf("bld anim %d, c%d w%d\n",lk,i,16)
				arstamp[lk].anim = append(arstamp[lk].anim,nil)
				arstamp[lk].anim[i] = blankimage(16+azx,16+azy)
				writepngtoimage(arstamp[lk].anim[i],16,16,azx,azy,psx+i,psy,0,0,0)
			}
		}
	}
	arstamp[lk].gtopl = gtopl
	G1 = gsv
}

func vpc_adj(x, y int) (int,int) {

	xba, yba := 0, 0
	if x < 0 { xba = absint(x) }
	if y < 0 { yba = absint(y) }
	return xba, yba
}

// animate tiles
var vlock bool		// viewport lock so maze loads dont blank screen
var nobld bool		// dont build multi layer
var tdir int		// test dir, incr 1 - 8 with logo

func animcon() {

	time.Sleep(2 * time.Second)
	for {
	time.Sleep(100 * time.Millisecond)
// anim not compatible with blot
  if !mbd && ccp != PASTE && !gvs {
	ablot = false
//fmt.Printf("in anim %t, sv %t\n",manim,svanim)
	mobflg := false						// on anim loop found some mimg layer items that need animation
	if manim {								// only run when anim tiles are on map

		xba, yba := vpc_adj(mvpx, mvpy)
	// we need to check bounds of current viewport, set animation of any visible floor tiles
	for y := mvpy; y < mvye; y++ {
		for x := mvpx; x < mvxe; x++ {
			_, ux, uy := lot(x, y, x, y)	// what would be nice when mapping for the vp, is to make a list of all animatables
			for alock {}					// cant access maps while segimg is writing
			tl := anmap[xy{ux, uy}]			// and not have to check 200 to 400 cells every frame
			r := anmapt[xy{ux, uy}]
			if tl > 0 {
//fmt.Printf("in anim %d, r %d @ %d, %d\n",tl,r,x,y)
				r--
				if r <= 0 {
					r = arstamp[tl].animtm
					if arstamp[tl].awalk[0] < 0 { r = -arstamp[tl].awalk[0] }
					if arstamp[tl].awalk[0] > 0 { r = arstamp[tl].awalk[0] }
				}
				if arstamp[tl].awalk[0] > 0 { arstamp[tl].awalk[9] = tdir; arstamp[tl].awalk[10] = arstamp[tl].awalk[tdir] }	// this has to be stored per entity
				anmapt[xy{ux, uy}] = r
				drimg := blankimage(16,16)
				if r < len(arstamp[tl].anim) {
					drimg = arstamp[tl].anim[r - 1]
				}
				offset := image.Pt(vcoord(x,mvpx,xba)*16, vcoord(y,mvpy,yba)*16)
				if arstamp[tl].mask & NOFLOOR != 0 {
					draw.Draw(fimg, drimg.Bounds().Add(offset), drimg, image.ZP, draw.Over)
				} else { mobflg = true }
			}
		}}
	}
	if mobflg { for vlock {}; flordirt, walsdirt = -1,-1; mimg = segimage(ebuf, xbuf, eflg, mvpx, mvpy, mvxe,mvye, false) } 
	rimg := blankimage(16*(mvxe-mvpx), 16*(mvye-mvpy))
	if !nobld {
		draw.Draw(rimg, fimg.Bounds(), fimg, image.ZP, draw.Over)
		draw.Draw(rimg, wimg.Bounds(), wimg, image.ZP, draw.Over)
		draw.Draw(rimg, mimg.Bounds(), mimg, image.ZP, draw.Over)
		rbimg = rimg
	}
	if !vlock  && actab == "Maze view" {

		maz_tab(cmain, rbimg, rbtn, blant)
		if ctrl {time.Sleep(1400 * time.Millisecond)}
	}}}
}

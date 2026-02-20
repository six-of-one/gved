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
	"fyne.io/fyne/v2/canvas"
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

func loadfail(x int, y int) *image.NRGBA {
	img := blankimage(x,y)
	b := HRGB{0xff008f8f}
	a := HRGB{0xff8f008f}
	var c color.Color
	for fy := 0; fy < y; fy += 8 {
		if fy & 8 == 8 { c = a } else { c = b }
		for fx := 0; fx < x; fx += 8 {

		for j := 0; j < 8; j++ {
			for i := 0; i < 8; i++ {
				img.Set(fx+i, fy+j, c)
			}}
			if c == b { c = a } else { c = b }
		}
	}
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
// gauntlet has diff base pallet from g2, this is easiest way for now
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
fmt.Println(fn)
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
	SE_G2		= 3		// gauntlet 2 mode - e.g. turn off g1
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
	case G1OBJ_MONST_GHOST2:
		arstamp[lk] = itemGetStamp("ghost2")
	case SEOBJ_G2GHOST:
		G1 = false; fallthrough
	case G1OBJ_MONST_GHOST3:
		arstamp[lk] = itemGetStamp("ghost")
	case G1OBJ_MONST_GRUNT1:
		arstamp[lk] = itemGetStamp("grunt1")
	case G1OBJ_MONST_GRUNT2:
		arstamp[lk] = itemGetStamp("grunt2")
	case SEOBJ_G2GRUNT:
		G1 = false; fallthrough
	case G1OBJ_MONST_GRUNT3:
		arstamp[lk] = itemGetStamp("grunt")
	case G1OBJ_MONST_DEMON1:
		arstamp[lk] = itemGetStamp("demon1")
	case G1OBJ_MONST_DEMON2:
		arstamp[lk] = itemGetStamp("demon2")
	case SEOBJ_G2DEMON:
		G1 = false; fallthrough
	case G1OBJ_MONST_DEMON3:
		arstamp[lk] = itemGetStamp("demon")
	case G1OBJ_MONST_LOBBER1:
		arstamp[lk] = itemGetStamp("lobber1")
	case G1OBJ_MONST_LOBBER2:
		arstamp[lk] = itemGetStamp("lobber2")
	case SEOBJ_G2LOBER:
		G1 = false; fallthrough
	case G1OBJ_MONST_LOBBER3:
		arstamp[lk] = itemGetStamp("lobber")
	case G1OBJ_MONST_SORC1:
		arstamp[lk] = itemGetStamp("sorcerer1")
	case G1OBJ_MONST_SORC2:
		arstamp[lk] = itemGetStamp("sorcerer2")
	case SEOBJ_G2SORC:
		G1 = false; fallthrough
	case G1OBJ_MONST_SORC3:
		arstamp[lk] = itemGetStamp("sorcerer")
	case SEOBJ_G2AUXGR:
		G1 = false
		arstamp[lk] = itemGetStamp("auxgrunt")

	case G1OBJ_MONST_DEATH:
		arstamp[lk] = itemGetStamp("death")
	case SEOBJ_G2ACID:
		G1 = false
		arstamp[lk] = itemGetStamp("acid")
	case SEOBJ_G2SUPSORC:
		G1 = false
		arstamp[lk] = itemGetStamp("supersorc")
	case SEOBJ_G2IT:
		G1 = false
		arstamp[lk] = itemGetStamp("it")
	case SEOBJ_MONST_DRAGON:
		G1 = false
		arstamp[lk] = itemGetStamp("dragon")

	case G1OBJ_MONST_THIEF:
		arstamp[lk] = itemGetStamp("thief")
	case SEOBJ_MONST_MUGGER:
		G1 = false
		arstamp[lk] = itemGetStamp("mugger")

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
// SE expand
	case SEOBJ_SE_ANKH:
		psx, psy = 21, 11
		mask = NOPOT
	case SEOBJ_FIRE_STICK:
		psx, psy = 32, 26
		mask = 256 | ANIM
		cnt = 4
	case SEOBJ_G2_POISPOT:
		psx, psy = 8, 11
		mask = NOPOT
	case SEOBJ_G2_POISFUD:
		psx, psy = 1, 11
		mask = NOFUD
	case SEOBJ_G2_QFUD:
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

func animcon() {

	for {
	time.Sleep(200 * time.Millisecond)
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
				r--
				if r <= 0 { r = arstamp[tl].animtm }
				anmapt[xy{ux, uy}] = r
				drimg := arstamp[tl].anim[r - 1]
				offset := image.Pt(vcoord(x,mvpx,xba)*16, vcoord(y,mvpy,yba)*16)
				if arstamp[tl].mask & NOFLOOR != 0 {
					draw.Draw(fimg, drimg.Bounds().Add(offset), drimg, image.ZP, draw.Over)
				} else { mobflg = true }
			}
		}}
	}
	if mobflg { for mlock {}; flordirt, walsdirt = -1,-1; mimg = segimage(ebuf, xbuf, eflg, mvpx, mvpy, mvxe,mvye, false) } 
	rimg := blankimage(16*(mvxe-mvpx), 16*(mvye-mvpy))
	if !nobld {
		draw.Draw(rimg, fimg.Bounds(), fimg, image.ZP, draw.Over)
		draw.Draw(rimg, wimg.Bounds(), wimg, image.ZP, draw.Over)
		draw.Draw(rimg, mimg.Bounds(), mimg, image.ZP, draw.Over)
//		rbimg = canvas.NewRasterFromImage(rimg)
		rbimg = rimg
	}
//	box := container.NewStack(rbtn, rbimg)
	if !vlock {
//		w.SetContent(box)
// left this in while options still has animcon maze tab
		cmzw.Remove(mzw)
		mzw = canvas.NewRasterFromImage(rbimg)
		cmzw.Add(mzw)
		cmzw.Refresh()
		maz_tab(cmain, rbimg, rbtn, blant)
	}}}
}

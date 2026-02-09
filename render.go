package main

import (
	"image"
	"image/color"
//	"image/color/palette"
	"image/png"
	"math"
	"os"
	"fmt"
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
	SE_FLOR		= 1		// load any gauntlet & g2 floors & walls
	SE_WALL		= 2
	SE_G2		= 3		// gauntlet 2 mode - e.g. turn off g1
	SE_CFLOR	= 4		// custom floor unit from xb*.ed file wall&floor lines
	SE_CWAL		= 5		// data: flor: sheet, r, (c is wally), xy size		wal: sheet, row, (xy size?)
	SE_COLRT	= 6		// color tiles under any
	SE_LETR		= 7		// draw a letter index to map_keymap (as gen hints) in color, as R,G,B, ind
	SE_MSG		= 8		// write a null term msg (up to 28 hex byts) onto maze in color as R,G,B, {MSG}, 00 
)
// bytes for each cmd
var parms = []int{
	0, 2, 2, 0, 1, 2, 3, 4, 36,
}
var secmd [64]int
var lastst string
var xpar [64]int		// extra parms past 3... - parms[] can NOT exceed this array size!

func parser(sp string, lc int) (int, int, int) {
	r1, r2, r3 := -1,0,0
//fmt.Printf("parse %s\n ",sp)
	var t int
	fmt.Sscanf(sp,"%X",&t) // validate
//	if t > 0 && lastst != sp { lastst = sp; fmt.Printf("t %X\n",t) }
	if t > 0 && lc > 0 {
		for i := 0; i < 17; i++ { secmd[i] = 0 }
					//	0	1	2	3	4	5	6	7	8	9	10	11	12	13	14	15	16	17	18	19	20	21	22	23	24	25	26	27	28	29	30	31	32	33	34	35	36
		fmt.Sscanf(sp,"%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X",
					&secmd[0],&secmd[1],&secmd[2],&secmd[3],&secmd[4],&secmd[5],&secmd[6],&secmd[7],&secmd[8],&secmd[9],&secmd[10],&secmd[11],&secmd[12],&secmd[11],&secmd[12],&secmd[13],&secmd[14],&secmd[15],&secmd[16],
					&secmd[16],&secmd[17],&secmd[18],&secmd[19],&secmd[20],&secmd[21],&secmd[22],&secmd[23],&secmd[24],&secmd[25],&secmd[26],&secmd[27],&secmd[28],&secmd[29],&secmd[30],&secmd[31],&secmd[32],
					&secmd[33],&secmd[34],&secmd[35],&secmd[36])
		for i := 0; i < 37; i++ {
			if lc == secmd[i] {
				prc := parms[secmd[i]]
				r1 =secmd[i+1]; r2 =secmd[i+2]; r3 =secmd[i+3]
				if prc == 0 { r1 = 1 }
				if prc > 3 { for k := 0; k < (prc - 3); k++ { xpar[k] = secmd[k+4]}}
//				fmt.Printf("c:%d p:%d - r1 %d r2 %d r3 %d\n ",lc,prc,r1,r2,r3)
				break
			}
//fmt.Printf("parm %d of %d= %d\n ",secmd[i],parms[secmd[i]])
			i = i + parms[secmd[i]]
		}

	}

	return r1, r2, r3
}
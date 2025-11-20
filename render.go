package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
//	"fmt"
)

type Color interface {
	IRGB() (irgb uint16)
}

type IRGB struct {
	irgb uint16
}

func (c IRGB) RGBA() (r, g, b, a uint32) {
	i := uint32(c.irgb&0xf000) >> 12
	r = uint32(c.irgb&0x0f00) >> 8 * i
	g = uint32(c.irgb&0x00f0) >> 4 * i
	b = uint32(c.irgb&0x000f) * i

	r = r << 8
	g = g << 8
	b = b << 8
	a = 0xffff

	return
}

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
			if tc == 0 && trans0 == true {
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
	writestamptoimage(img, stamp, 0, 0)

	return img
}

// FIXME: Rename later!
func writestamptoimage(img *image.NRGBA, stamp *Stamp, xloc int, yloc int) {
	p := gauntletPalettes[stamp.ptype][stamp.pnum]
	for y := 0; y < len(stamp.data)/stamp.width; y++ {
		for x := 0; x < stamp.width; x++ {
			// fmt.Printf("Writing to image at %d,%d\n", xloc, yloc)
			writetiletoimage(img, stamp.data[(stamp.width*y)+x], p,
				stamp.trans0, xloc+(x*8), yloc+(y*8))
		}
	}
}

func savetopng(fn string, img *image.NRGBA) {
	f, _ := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	png.Encode(f, img)
}

// func genanim(animarray []int, xtiles int, ytiles int) []*image.Paletted {
// 	var images []*image.Paletted

// 	for i := 0; i < len(animarray); i++ {
// 		fmt.Printf("generating from tile %d\n", animarray[i])
// 		images = append(images, genimage(animarray[i], xtiles, ytiles))
// 	}

// 	return images
// }

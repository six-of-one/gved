package main

import (
	"strings"
	"fmt"
	"fyne.io/fyne/v2/canvas"
	"os"
	"image"
)

// type Stamp struct {
//     width   int
//     numbers []int
//     data    []TileData
//     ptype   string
//     pnum    int
//     trans0  bool
// }

var itemStamps = map[string]Stamp{
	"blank": Stamp{
		width:   2,
		numbers: []int{0, 0, 0, 0},
		ptype:   "base",
		pnum:    0,
		trans0:  false,
		mask:    NOTHN,
	},

	"key": Stamp{
		width:   2,
		numbers: tilerange(0xafc, 4),
		ptype:   "base",
		pnum:    1,
		trans0:  true,
		mask:    NODOR,
	},
	"keyring": Stamp{
		width:   3,
		numbers: tilerange(0x1d76, 6),
		ptype:   "base",
		pnum:    1,
		trans0:  true,
		mask:    NODOR,
	},

	"food": Stamp{
		width:   3,
		numbers: tilerange(0x963, 9),
		ptype:   "base",
		pnum:    1,
		trans0:  true, nudgex: -4,
		nudgey: -4,
		mask:    NOFUD,
	},
	"ifood1": Stamp{
		width:   3,
		numbers: tilerange(0x96c, 9),
		ptype:   "base",
		pnum:    1,
		trans0:  true, nudgex: -4,
		nudgey: -8,
		mask:    NOFUD,
	},
	"ifood2": Stamp{
		width:   3,
		numbers: tilerange(0x975, 9),
		ptype:   "base",
		pnum:    1,
		trans0:  true, nudgex: -4,
		nudgey: -8,
		mask:    NOFUD,
	},
	"ifood3": Stamp{
		width:   3,
		numbers: tilerange(0x97e, 9),
		ptype:   "base",
		pnum:    1,
		trans0:  true, nudgex: -4,
		nudgey: -8,
		mask:    NOFUD,
	},
	"mfood": Stamp{
		width:   3,
		numbers: tilerange(0x277b, 9),
		ptype:   "base",
		pnum:    1,
		trans0:  true, nudgex: -4,
		nudgey: -4,
		mask:    NOFUD,
	},
	"pfood": Stamp{
		width:   3,
		numbers: tilerange(0x25ed, 9),
		ptype:   "base",
		pnum:    1,
		trans0:  true, nudgex: -4,
		nudgey: -4,
		mask:    NOFUD,
	},

	"potion": Stamp{
		width:   2,
		numbers: tilerange(0x8fc, 4),
		ptype:   "base",
		pnum:    1,
		trans0:  true,
		mask:    NOPOT,
	},
	"ipotion": Stamp{
		width:   2,
		numbers: tilerange(0x9fc, 4),
		ptype:   "base",
		pnum:    1,
		trans0:  true,
		mask:    NOPOT,
	},
	"ppotion": Stamp{
		width:   2,
		numbers: tilerange(0x20fc, 4),
		ptype:   "base",
		pnum:    1,
		trans0:  true,
		mask:    NOPOT,
	},

	"shieldpotion": Stamp{
		width:   2,
		numbers: tilerange(0x11fc, 4),
		ptype:   "base",
		pnum:    1,
		trans0:  true,
		mask:    NOPOT,
	},
	"speedpotion": Stamp{
		width:   2,
		numbers: tilerange(0x12fc, 4),
		ptype:   "base",
		pnum:    1,
		trans0:  true,
		mask:    NOPOT,
	},
	"magicpotion": Stamp{
		width:   2,
		numbers: tilerange(0x13fc, 4),
		ptype:   "base",
		pnum:    1,
		trans0:  true,
		mask:    NOPOT,
	},
	"shotpowerpotion": Stamp{
		width:   2,
		numbers: tilerange(0x14fc, 4),
		ptype:   "base",
		pnum:    1,
		trans0:  true,
		mask:    NOPOT,
	},
	"shotspeedpotion": Stamp{
		width:   2,
		numbers: tilerange(0x15fc, 4),
		ptype:   "base",
		pnum:    1,
		trans0:  true,
		mask:    NOPOT,
	},
	"fightpotion": Stamp{
		width:   2,
		numbers: tilerange(0x16fc, 4),
		ptype:   "base",
		pnum:    1,
		trans0:  true,
		mask:    NOPOT,
	},

	"invis": Stamp{
		width:   3,
		numbers: tilerange(0x1700, 9),
		ptype:   "base",
		pnum:    1,
		trans0:  true, nudgex: -4,
		nudgey: -4,
		mask:    NOPOT,
	},
	"transportability": Stamp{
		width:   2,
		numbers: tilerange(0x23fc, 4),
		ptype:   "base",
		pnum:    1,
		trans0:  true,
		mask:    NOPOT,
	},
	"reflect": Stamp{
		width:   2,
		numbers: tilerange(0x24fc, 4),
		ptype:   "base",
		pnum:    1,
		trans0:  true,
		mask:    NOPOT,
	},
	"repulse": Stamp{
		width:   2,
		numbers: tilerange(0x26fc, 4),
		ptype:   "base",
		pnum:    1,
		trans0:  true,
		mask:    NOPOT,
	},
	"invuln": Stamp{
		width:   2,
		numbers: tilerange(0x2784, 4),
		ptype:   "base",
		pnum:    1,
		trans0:  true,
		mask:    NOPOT,
	},
	"supershot": Stamp{
		width:   2,
		numbers: tilerange(0x2788, 4),
		ptype:   "base",
		pnum:    1,
		trans0:  true,
		mask:    NOPOT,
	},

	"pushwall": Stamp{
		width:   3,
		numbers: tilerange(0x20f6, 6),
		ptype:   "base",
		pnum:    1,
		trans0:  true,
		nudgex:  -4,
		nudgey:  -4,
		mask:    NOEXP,
	},

	"treasure": Stamp{
		width:   3,
		numbers: tilerange(0x987, 9),
		ptype:   "base",
		pnum:    1,
		trans0:  true,
		nudgex:  -4,
		nudgey:  -8,
		mask:    NOTRS,
	},
	"treasurelocked": Stamp{
		width:   3,
		numbers: tilerange(0x25e4, 9),
		ptype:   "base",
		pnum:    1,
		trans0:  true,
		nudgex:  -4,
		nudgey:  -4,
		mask:    NOTRS,
	},
	"goldbag": Stamp{
		width:   3,
		numbers: tilerange(0x9a2, 9),
		ptype:   "base",
		pnum:    1,
		trans0:  true,
		nudgex:  -4,
		nudgey:  -8,
		mask:    NOTRS,
	},

	"tport": Stamp{
		width:   2,
		numbers: tilerange(0x49e, 4),		// g2
		ptype:   "teleff",
		pnum:    0,
		trans0:  true,
		mask:    NOEXP,
	},
	"tportg1": Stamp{
		width:   2,
		numbers: tilerange(0x3a4, 4),		// g1
		ptype:   "teleff",
		pnum:    0,
		trans0:  true,
		mask:    NOEXP,
	},

	// ?: missing all the various directions
	"ff": Stamp{
		width:   2,
		numbers: tilerange(0x4a2, 4),
		ptype:   "teleff",
		pnum:    0,
		trans0:  true,
		mask:    NOTRAP,
	},

	"exit": Stamp{
		width:   2,
		numbers: []int{0x39e, 0x39f, 0x6, 0x6},
		ptype:   "floor",
		pnum:    0,
		trans0:  false,
		mask:    NOEXP,
	},
	"exit6": Stamp{
		width:   2,
		numbers: tilerange(0x39e, 4),
		ptype:   "floor",
		pnum:    0,
		trans0:  false,
		mask:    NOEXP,
	},
	"exitg1": Stamp{				// g1 exits wont take floor palette it seems
		width:   2,
		numbers: tilerange(0xbfc, 4),
		ptype:   "base",
		pnum:    5,
		trans0:  false,
		mask:    NOEXP,
	},
	"exit4": Stamp{
		width:   2,
		numbers: tilerange(0xcfc, 4),
		ptype:   "base",
		pnum:    5,
		trans0:  false,
		mask:    NOEXP,
	},
	"exit8": Stamp{
		width:   2,
		numbers: tilerange(0xdfc, 4),
		ptype:   "base",
		pnum:    5,
		trans0:  false,
		mask:    NOEXP,
	},

	"vdoor": Stamp{
		width:   2,
		numbers: tilerange(0x1d80, 4),
		ptype:   "base",
		pnum:    0,
		trans0:  true,
		mask:    NODOR,
	},
	"hdoor": Stamp{
		width:   2,
		numbers: tilerange(0x1d48, 4),
		ptype:   "base",
		pnum:    0,
		trans0:  true,
		mask:    NODOR,
	},

	"plus": Stamp{
		width:   2,
		numbers: tilerange(0xbfc, 4),
		ptype:   "base",
		pnum:    1,
		trans0:  true,
		mask:    NOTHN,
	},
	"plusg1": Stamp{
		width:   2,
		numbers: tilerange(0x1e09, 4),
		ptype:   "base",
		pnum:    0,
		trans0:  true,
		mask:    NOTHN,
	},

	// FIXME: Needs to be in monsters, really
	"dragon": Stamp{
		width:   4,
		numbers: tilerange(0x2100, 16),
		ptype:   "base",
		pnum:    8, // or 7 or 6
		trans0:  true,
		nudgex:  0,
		nudgey:  -16, // because we're rendering the maze "upside-down" (I think)
		mask:    NOMON,
	},
	"generator1": Stamp{
		width:   3,
		numbers: tilerange(0x9c6, 9),
		ptype:   "base",
		pnum:    5,
		trans0:  true,
		nudgex:  -4,
		nudgey:  -7,
		mask:    NOGEN,
	},
	"generator2": Stamp{
		width:   3,
		numbers: tilerange(0x9cf, 9),
		ptype:   "base",
		pnum:    5,
		trans0:  true,
		nudgex:  -4,
		nudgey:  -7,
		mask:    NOGEN,
	},
	"generator3": Stamp{
		width:   3,
		numbers: tilerange(0x9d8, 9),
		ptype:   "base",
		pnum:    5,
		trans0:  true,
		nudgex:  -4,
		nudgey:  -7,
		mask:    NOGEN,
	},
	"ghostgen1": Stamp{
		width:   3,
		numbers: tilerange(0x9ab, 9),
		ptype:   "base",
		pnum:    5,
		trans0:  true,
		nudgex:  -4,
		nudgey:  -4,
		mask:    NOGEN,
	},
	"ghostgen2": Stamp{
		width:   3,
		numbers: tilerange(0x9b4, 9),
		ptype:   "base",
		pnum:    5,
		trans0:  true,
		nudgex:  -4,
		nudgey:  -5,
		mask:    NOGEN,
	},
	"ghostgen3": Stamp{
		width:   3,
		numbers: tilerange(0x9bd, 9),
		ptype:   "base",
		pnum:    5,
		trans0:  true,
		nudgex:  -4,
		nudgey:  -7,
		mask:    NOGEN,
	},

	"ghost": Stamp{
		width:   3,
		numbers: tilerange(0x800, 9),
		ptype:   "base",
		pnum:    4,
		trans0:  true,
		nudgex:  -4,
		nudgey:  -4,
		mask:    NOMON,
	},
// encode levels 1, 2 for pre-gen monsters
	"ghost2": Stamp{
		width:   3,
		numbers: tilerange(0x800, 9),
		ptype:   "base",
		pnum:    3,
		trans0:  true,
		nudgex:  -4,
		nudgey:  -4,
		mask:    NOMON,
	},
	"ghost1": Stamp{
		width:   3,
		numbers: tilerange(0x800, 9),
		ptype:   "base",
		pnum:    2,
		trans0:  true,
		nudgex:  -4,
		nudgey:  -4,
		mask:    NOMON,
	},
	"grunt": Stamp{
		width:   3,
		numbers: tilerange(0x9e1, 9),
		ptype:   "base",
		pnum:    4,
		trans0:  true,
		nudgex:  -4,
		nudgey:  -4,
		mask:    NOMON,
	},
	"grunt2": Stamp{
		width:   3,
		numbers: tilerange(0x9e1, 9),
		ptype:   "base",
		pnum:    3,
		trans0:  true,
		nudgex:  -4,
		nudgey:  -4,
		mask:    NOMON,
	},
	"grunt1": Stamp{
		width:   3,
		numbers: tilerange(0x9e1, 9),
		ptype:   "base",
		pnum:    2,
		trans0:  true,
		nudgex:  -4,
		nudgey:  -4,
		mask:    NOMON,
	},
	"demon": Stamp{
		width:   3,
		numbers: tilerange(0x183f, 9),
		ptype:   "base",
		pnum:    8,
		trans0:  true,
		nudgex:  -4,
		nudgey:  -6,
		mask:    NOMON,
	},
	"demon2": Stamp{
		width:   3,
		numbers: tilerange(0x183f, 9),
		ptype:   "base",
		pnum:    7,
		trans0:  true,
		nudgex:  -4,
		nudgey:  -6,
		mask:    NOMON,
	},
	"demon1": Stamp{
		width:   3,
		numbers: tilerange(0x183f, 9),
		ptype:   "base",
		pnum:    6,
		trans0:  true,
		nudgex:  -4,
		nudgey:  -6,
		mask:    NOMON,
	},
	"lobber": Stamp{
		width:   3,
		numbers: tilerange(0x1b57, 6),
		ptype:   "base",
		pnum:    11,
		trans0:  true,
		nudgex:  -4,
		nudgey:  -4,
		mask:    NOMON,
	},
	"lobber2": Stamp{
		width:   3,
		numbers: tilerange(0x1b57, 6),
		ptype:   "base",
		pnum:    10,
		trans0:  true,
		nudgex:  -4,
		nudgey:  -4,
		mask:    NOMON,
	},
	"lobber1": Stamp{
		width:   3,
		numbers: tilerange(0x1b57, 6),
		ptype:   "base",
		pnum:    9,
		trans0:  true,
		nudgex:  -4,
		nudgey:  -4,
		mask:    NOMON,
	},
	"sorcerer": Stamp{
		width:   3,
		numbers: tilerange(0x13a2, 9),
		ptype:   "base",
		pnum:    11,
		trans0:  true,
		nudgex:  -4,
		nudgey:  -4,
		mask:    NOMON,
	},
	"sorcerer2": Stamp{
		width:   3,
		numbers: tilerange(0x13a2, 9),
		ptype:   "base",
		pnum:    10,
		trans0:  true,
		nudgex:  -4,
		nudgey:  -4,
		mask:    NOMON,
	},
	"sorcerer1": Stamp{
		width:   3,
		numbers: tilerange(0x13a2, 9),
		ptype:   "base",
		pnum:    9,
		trans0:  true,
		nudgex:  -4,
		nudgey:  -4,
		mask:    NOMON,
	},
	"auxgrunt": Stamp{
		width:   3,
		numbers: tilerange(0x9e1, 9),
		ptype:   "base",
		pnum:    4,
		trans0:  true,
		nudgex:  -4,
		nudgey:  -4,
		mask:    NOMON,
	},
	"death": Stamp{
		width:   3,
		numbers: tilerange(0x1a75, 9),
		ptype:   "base",
		pnum:    0,
		trans0:  true,
		nudgex:  -4,
		nudgey:  -4,
		mask:    NOMON,
	},
	"thief": Stamp{
		width:   3,
		numbers: tilerange(0xdea, 9),
		ptype:   "base",
		pnum:    0,
		trans0:  true,
		nudgex:  -4,
		nudgey:  -4,
		mask:    NOMON,
	},
	"mugger": Stamp{
		width:   3,
		numbers: tilerange(0x24ea, 9),
		ptype:   "base",
		pnum:    1,
		trans0:  true,
		nudgex:  -4,
		nudgey:  -4,
		mask:    NOMON,
	},
	"acid": Stamp{
		width:   3,
		numbers: tilerange(0x2300, 9),
		ptype:   "base",
		pnum:    1,
		trans0:  true,
		nudgex:  -4,
		nudgey:  -4,
		mask:    NOMON,
	},
	"supersorc": Stamp{
		width:   3,
		numbers: tilerange(0x13a2, 9),
		ptype:   "base",
		pnum:    11,
		trans0:  true,
		nudgex:  -4,
		nudgey:  -4,
		mask:    NOMON,
	},
	"it": Stamp{
		width:   3,
		numbers: tilerange(0x2600, 9),
		ptype:   "base",
		pnum:    8,
		trans0:  true,
		nudgex:  -4,
		nudgey:  -4,
		mask:    NOMON,
	},
	"wizard": Stamp{
		width:   3,
		numbers: tilerange(0x135A, 9),
		ptype:   "wizard",
		pnum:    2,
		trans0:  true,
		nudgex:  -4,
		nudgey:  -4,
		mask:    0,
	},

	"arrowleft": Stamp{
		width: 2,
		// numbers: tilerange(0x1c8f, 4),
		numbers: []int{0, 0x1c8f, 0, 0x1c91},
		ptype:   "elf",
		pnum:    1,
		trans0:  true,
	},
	"arrowright": Stamp{
		width: 2,
		// numbers: tilerange(0x1c7c, 4),
		numbers: []int{0x1c7d, 0, 0x1c7f, 0},
		ptype:   "elf",
		pnum:    1,
		trans0:  true,
	},
	"arrowup": Stamp{
		width: 2,
		// numbers: tilerange(0x1c74, 4),
		numbers: []int{0, 0, 0x1c74, 0x1c75},
		ptype:   "elf",
		pnum:    1,
		trans0:  true,
	},
	"arrowdown": Stamp{
		width: 2,
		// numbers: tilerange(0x1c84, 4),
		numbers: []int{0x1c86, 0x1c87, 0, 0},
		ptype:   "elf",
		pnum:    1,
		trans0:  true,
	},
}

func tilerange(start int, count int) []int {
	r := make([]int, count)
	for i := range r {
		r[i] = start
		start += 1
	}
	return r
}

//type AFrames []Stamp
// type MobAnimFrames []int
// type MobAnimsDir map[string]MobAnimFrames
// type MobAnims map[string]MobAnimsDir

// var ghostAnims = MobAnims{
//     "walk": {
//         "down":      {0x800, 0x809, 0x812, 0x81b},
//         "downright": {0x824, 0x82d, 0x836, 0x83f},
//         "right":     {0x848, 0x851, 0x85a, 0x863},
//         "upright":   {0x86c, 0x875, 0x87e, 0x887},
//         "up":        {0x890, 0x899, 0x8a2, 0x8ab},
//         "upleft":    {0x8b4, 0x8bd, 0x8c6, 0x8cf},
//         "left":      {0x8d8, 0x8e1, 0x8ea, 0x8f3},
//         "downleft":  {0x900, 0x909, 0x912, 0x91b},
//     },
// }

// var monsters = map[string]Monster{
// 	"ghost": {
// 		xsize: 3,
// 		ysize: 3,
// 		ptype: "base",
// 		pnum:  4,
// 		// palette: gauntletPalettes["base"][0],
// 		anims: ghostAnims,
// 	},
// }

// var reMonsterType = regexp.MustCompile(`^(ghost)(\d+)?`)
// var reMonsterAction = regexp.MustCompile(`^(walk|fight|attack)`)
// var reMonsterDir = regexp.MustCompile(`^(up|upright|right|downright|down|downleft|left|upleft)`)

func doitem(arg string) {
	split := strings.Split(arg, "-")

	c := 0
	maxh := 2
	maxw := 0
	all := false

	for _, ss := range split {
		if ss != "item" { c++ }
		stamp := itemGetStamp(ss)
		height := len(stamp.numbers) / stamp.width
		if height > maxh { maxh = height }
		maxw += stamp.width
		if ss == "all" { all = true }
	}
	if all {
		fmt.Printf("blank\nkey\nkeyring\nfood\nifood1\nifood2\nifood3\nmfood\npfood\npotion\nipotion\nppotion\n"+
				   "shieldpotion\nspeedpotion\nmagicpotion\nshotpowerpotion\nshotspeedpotion\nfightpotion\ninvis\n"+
				   "transportability\nreflect\nrepulse\ninvuln\nsupershot\npushwall\ntreasure\ntreasurelocked\ngoldbag\n"+
				   "tport\ntportg1\nff\nexit\nexit4\nexit6\nexit8\nvdoor\nhdoor\nplus\nplusg1\ndragon\n"+
				   "generator1\ngenerator2\ngenerator3\nghostgen1\nghostgen2\nghostgen3\nghost\nghost2\nghost1\n"+
				   "grunt\ngrunt2\ngrunt1\ndemon\ndemon2\ndemon1\nlobber\nlobber2\nlobber1\nsorcerer\nsorcerer2\nsorcerer1\n"+
				   "auxgrunt\ndeath\nthief\nacid\nsupersorc\nit\narrowleft\narrowright\narrowup\narrowdown\n\n")
	} else {
		img := blankimage(16*maxw, 16*maxh)
		pos := 8

		for _, ss := range split {

			if ss != "item" {
				stamp := itemGetStamp(ss)
				writestamptoimage(img, stamp, pos, 8)
				pos += stamp.width*16
			}
		}

		savetopng(opts.Output, img)
	}
}

func itemGetStamp(itemType string) *Stamp {
	stamp, ok := itemStamps[itemType]

	if !ok {
// failed to get that item - just return extra speed as a warning
		stamp, ok = itemStamps["speedpotion"]
		if !ok {
			panic("total fail bad item: " + itemType + "and speedpotion")
		}
	}

// if nothing bit matches mask in item, send blank back
	if stamp.mask & nothing == 0 {
		fillstamp(&stamp)
	}
	return &stamp
}

// expand for sanctuary gfx

func itemGetPNG(fil string) (error, *canvas.Image) {
	rimg := canvas.NewImageFromImage(blankimage(16, 16))

	inf, err := os.Open(fil)
	if err == nil {
		src, _, err := image.Decode(inf)
		if err == nil {
			rimg = canvas.NewImageFromImage(src)
		}
	}
	defer inf.Close()

	return err, rimg
}
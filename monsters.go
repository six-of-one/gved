package main

import (
	"regexp"
	"strconv"
	"strings"
)

type Monster struct {
	xsize int
	ysize int
	ptype string
	pnum  int
	// palette []color.Color
	anims MobAnims
}

type MobAnimFrames []int
type MobAnimsDir map[string]MobAnimFrames
type MobAnims map[string]MobAnimsDir

var monsters = map[string]Monster{
	"ghost": {
		xsize: 3,
		ysize: 3,
		ptype: "base",
		pnum:  0, // FIXME: This is weird and seems wrong
		// palette: gauntletPalettes["base"][0],
		anims: ghostAnims,
	},
}

var reMonsterType = regexp.MustCompile(`^(ghost)(\d+)?`)
var reMonsterAction = regexp.MustCompile(`^(walk|fight|attack)`)
var reMonsterDir = regexp.MustCompile(`^(up|upright|right|downright|down|downleft|left|upleft)`)

func domonster(arg string) {
	split := strings.Split(arg, "-")

	var monsterType string
	var monsterAction = "walk"
	var monsterDir = "up"
	var monsterLevel = 1

	// Still wonder if there's a cleaner way
	for _, ss := range split {
		switch {
		case reMonsterType.MatchString(ss):
			monst := reMonsterType.FindStringSubmatch(ss)
			monsterType = monst[1]
			if monst[2] != "" {
				ml, _ := strconv.ParseInt(monst[2], 10, 0)
				monsterLevel = int(ml)
			}
		case reMonsterAction.MatchString(ss):
			monsterAction = reMonsterAction.FindStringSubmatch(ss)[1]

		case reMonsterDir.MatchString(ss):
			monsterDir = reMonsterDir.FindStringSubmatch(ss)[1]
		}
	}

	opts.PalType = monsters[monsterType].ptype
	opts.PalNum = monsters[monsterType].pnum + (monsterLevel + 1) // FIME: This is weird and seems wrong

	if opts.Animate {
		// t := monsters[monsterType].anims[monsterAction][monsterDir]
		// x := monsters[monsterType].xsize
		// y := monsters[monsterType].ysize
		// imgs := genanim(t, x, y)

		// f, _ := os.OpenFile(opts.Output, os.O_WRONLY|os.O_CREATE, 0600)
		// defer f.Close()

		// var delays []int
		// for i := 0; i < len(t); i++ {
		// 	delays = append(delays, 15)
		// }

		// gif.EncodeAll(f,
		// 	&gif.GIF{
		// 		Image: imgs,
		// 		Delay: delays,
		// 	},
		// )
	} else {
		// fmt.Printf("%#v\n", monsters[monsterType].anims["walk"])
		// fmt.Printf("Action: %#v\n", monsterAction)
		t := monsters[monsterType].anims[monsterAction][monsterDir][0]
		x := monsters[monsterType].xsize
		y := monsters[monsterType].ysize
		img := genimage(t, x, y)
		savetopng(opts.Output, img)
	}
}

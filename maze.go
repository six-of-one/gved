package main

// #include <stdlib.h>
// #include <Tilengine.h>

import "C"

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"os"
	"bufio"
//	"time"

	"fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/canvas"
//	"fyne.io/fyne/v2/driver/desktop"
)

func mazeMetaPrint(maze *Maze) {
	fmt.Printf("  Encoded length: %3d bytes\n", maze.encodedbytes)
	fmt.Printf("  Wall pattern: %02d, Wall color: %02d     Floor pattern: %02d, Floor color: %02d\n",
		maze.wallpattern, maze.wallcolor, maze.floorpattern, maze.floorcolor)
	fmt.Printf("  Flags: ")
	for k, v := range mazeFlagStrings {
		if (maze.flags & k) != 0 {
			fmt.Printf("%s ", v)
		}
	}
	if G2 {
		fmt.Printf("\n  Random food adds: %d\n", (maze.flags&LFLAG3_RANDOMFOOD_MASK)>>8)
		fmt.Printf("  Secret trick: %2d - %s\n", maze.secret, mazeSecretStrings[maze.secret])
	}
}

var reMazeNum = regexp.MustCompile(`^maze(\d+)`)
var reMazeMeta = regexp.MustCompile(`^meta$`)

func domaze(arg string) {
	split := strings.Split(arg, "-")

	var mazeNum = -1
	var mazeMeta = 0
	var maxmaze = 116

// g1 has more mazes, but treasure rooms can only spec from address, for now
	if G1 { maxmaze = 127 }

	for _, ss := range split {
		switch {
		case reMazeNum.MatchString(ss):
			m, _ := strconv.ParseInt(reMazeNum.FindStringSubmatch(ss)[1], 10, 0)
			mazeNum = int(m) - 1

		case reMazeMeta.MatchString(ss):
			mazeMeta = 1
		}
	}
	if mazeNum < 0 || mazeNum > maxmaze {
		panic("Invalid maze number / address specified.")
	}

	fmt.Printf("Maze number: %d", mazeNum + 1)
	if Aov > 0 {
		fmt.Printf(", address: 0x%X ", Aov)
	}
	fmt.Printf("\n")

// set 1 override to -1 to set in decoder
	Ovwallpat = -1

	maze := mazeDecompress(slapsticReadMaze(mazeNum), mazeMeta > 0)
	xform := make(map[xy]int)

	if opts.Verbose || mazeMeta > 0 {
		mazeMetaPrint(maze)
		if mazeMeta > 0 { os.Exit(0) }
	}

// interactive viewer not selected - gen maze, output png & exit
	if !opts.Intr {
		genpfimage(maze, mazeNum)
		os.Exit(0)
	}

// setup kby read
	consoleReader := bufio.NewReaderSize(os.Stdin, 1)
// setup window

    a = app.New()
    w = a.NewWindow("G¹G²ved")

	menuItemExit := fyne.NewMenuItem("Exit...", func() {
		os.Exit(0)
	})
	menuExit := fyne.NewMenu("Exit ", menuItemExit)
	menuItemAbout := fyne.NewMenuItem("About...", func() {
		dialog.ShowInformation("About G¹G²ved", "Gauntlet / Gauntlet 2 visual editor\nAuthor: Six\n\ngithub.com/six-of-one/", w)
	})
	menuHelp := fyne.NewMenu("Help ", menuItemAbout)
	mainMenu := fyne.NewMainMenu(menuExit, menuHelp)
	w.SetMainMenu(mainMenu)


	genpfimage(maze, mazeNum)

// testing gotilengine - leftover
	suser := " "	// user action string


// interactive loop here - lets user tweak vars settings & load new mazes
// user controls loop for tweaking
		noact := false
// input new maze #
		anum := -1
		var ascii byte
		for {


		if !noact {
// redo maze #, colors, walls, rotates, etc
			if (anum > 0 && anum <= 127 || anum >= 229376 && anum < 262145) && ascii == 97 {
				if anum <= 127 {
					fmt.Printf("\nnew maze: %d\n",anum)
					mazeNum = anum - 1
					Aov = 0
					suser = fmt.Sprintf("maze = %d",anum)
				} else {
					fmt.Printf("\nnew addr: %d\n",anum)
					Aov = anum
					suser = fmt.Sprintf("addr = %d",anum)
				}
				anum = -1
// clear these when load new maze
				Ovwallpat = -1
			}
			maze = mazeDecompress(slapsticReadMaze(mazeNum), mazeMeta > 0)

// manual mirror, flip
	if opts.MH || opts.MV || opts.MRP || opts.MRM {
		lastx := 32
		if maze.flags&LFLAG4_WRAP_H > 0 {
			lastx = 31
		}

		lasty := 32
		if maze.flags&LFLAG4_WRAP_V > 0 {
			lasty = 31
		}
// note it
/*
	for y := 0; y <= lasty; y++ {
		for x := 0; x <= lastx; x++ {

			fmt.Printf(" %02d", maze.data[xy{x, y}])
		}
		fmt.Printf("\n")
	}
*/
// transform
// rotate +90 degrees
		if opts.MRP {
			for ty := 1; ty <= lasty; ty++ {
			for tx := 1; tx <= lastx; tx++ {
				xform[xy{lastx - tx, ty}] = maze.data[xy{ty, tx}]
// g1 - must transform all dors on a rotat since they have horiz & vert dependent
				if xform[xy{lastx - tx, ty}] == G1OBJ_DOOR_HORIZ { xform[xy{lastx - tx, ty}] = G1OBJ_DOOR_VERT } else {
				if xform[xy{lastx - tx, ty}] == G1OBJ_DOOR_VERT { xform[xy{lastx - tx, ty}] = G1OBJ_DOOR_HORIZ } }
// g2
				if xform[xy{lastx - tx, ty}] == MAZEOBJ_DOOR_HORIZ { xform[xy{lastx - tx, ty}] = MAZEOBJ_DOOR_VERT } else {
				if xform[xy{lastx - tx, ty}] == MAZEOBJ_DOOR_VERT { xform[xy{lastx - tx, ty}] = MAZEOBJ_DOOR_HORIZ } }
			}}
		} else {
		if opts.MRM {
			for ty := 1; ty <= lasty; ty++ {
			for tx := 1; tx <= lastx; tx++ {
				xform[xy{tx, lasty - ty}] = maze.data[xy{ty, tx}]
// g1
				if xform[xy{tx, lasty - ty}] == G1OBJ_DOOR_HORIZ { xform[xy{tx, lasty - ty}] = G1OBJ_DOOR_VERT } else {
				if xform[xy{tx, lasty - ty}] == G1OBJ_DOOR_VERT { xform[xy{tx, lasty - ty}] = G1OBJ_DOOR_HORIZ } }
// g2
				if xform[xy{tx, lasty - ty}] == MAZEOBJ_DOOR_HORIZ { xform[xy{tx, lasty - ty}] = MAZEOBJ_DOOR_VERT } else {
				if xform[xy{tx, lasty - ty}] == MAZEOBJ_DOOR_VERT { xform[xy{tx, lasty - ty}] = MAZEOBJ_DOOR_HORIZ } }
			}}
		}
		}

// have to copy back if doing with any mirror cmd
		if opts.MRP || opts.MRM {
		if opts.MH || opts.MV {
		for y := 1; y <= lasty; y++ {
			for x := 1; x <= lastx; x++ { maze.data[xy{x, y}] = xform[xy{x, y}] }
		}}}

// mirror x
		if opts.MH {
			for ty := 1; ty <= lasty; ty++ {
			for tx := 1; tx <= lastx; tx++ {
				xform[xy{lastx - tx, ty}] = maze.data[xy{tx, ty}]
			}}
		}
// have to copy back if doing both together
		if opts.MH && opts.MV {
		for y := 1; y <= lasty; y++ {
			for x := 1; x <= lastx; x++ { maze.data[xy{x, y}] = xform[xy{x, y}] }
		}}

// mirror y: flip
		if opts.MV {
			for ty := 1; ty <= lasty; ty++ {
			for tx := 1; tx <= lastx; tx++ {
				xform[xy{tx, lasty - ty}] = maze.data[xy{tx, ty}]
			}}
		}
		if opts.MH || opts.MV || opts.MRP || opts.MRM {
			suser += ","
			if opts.MV { suser += " m-vert" }
			if opts.MH { suser += " m-horz" }
			if opts.MRP { suser += "+90°" }
			if opts.MRM { suser += "-90°" }
		}
// copy back
		for y := 1; y <= lasty; y++ {
			for x := 1; x <= lastx; x++ { maze.data[xy{x, y}] = xform[xy{x, y}] }
		}
	}

			Ovimg := genpfimage(maze, mazeNum)
			bimg := canvas.NewRasterFromImage(Ovimg)
			w.Canvas().SetContent(bimg)
			w.Resize(fyne.NewSize(1024, 1024))
			w.Show()
			w.CenterOnScreen()
			til := fmt.Sprintf("G¹G²ved Maze: %d",mazeNum)
			w.SetTitle(til)
			w.Canvas().SetOnTypedRune(typedRune)

			fmt.Printf("G%d Command (?, q, fFgG, wWeE, rRt, hm, s, il, u, v, #a): ",opts.Gtp)
		}
		a.Run()
// key tester
		if ascii != 17 {
			input, _ := consoleReader.ReadByte()
			ascii = input
		} else {
			ascii = 97
		}
// ESC = 27 and q = 113
			if ascii == 27 || ascii == 113 {
				fmt.Printf("Exiting...\n")
//				gotilengine.TLN_DeleteBitmap(bkg)
//				gotilengine.TLN_Deinit()
				os.Exit(0)
			}
			noact = false
			switch ascii {
			case 10:
// it picks up the <CR> that enters cmd, mask that off here, do nothing
				noact = true
				anum = -1
			case 119:		// w
				Ovwallpat += 1
				if anum >= 0 { Ovwallpat = anum }
				if Ovwallpat > 7 { Ovwallpat = 0 }
				fmt.Printf("cmd: w - wallp: %d\n",Ovwallpat)
			case 87:		// W
				Ovwallpat -= 1
				if Ovwallpat < 0 { Ovwallpat = 7 }
				fmt.Printf("cmd: W - wallp: %d\n",Ovwallpat)
			case 101:		// e
				Ovwallcol += 1
				if anum >= 0 { Ovwallcol = anum }
				if Ovwallcol > 16 { Ovwallcol = 0 }
				fmt.Printf("cmd: e - wallc: %d\n",Ovwallcol)
			case 69:
				Ovwallcol -= 1
				if Ovwallcol < 0 { Ovwallcol = 16 }
				fmt.Printf("cmd: E - wallc: %d\n",Ovwallcol)
			case 102:		// f
				Ovflorpat += 1
				if anum >= 0 { Ovflorpat = anum }
				if Ovflorpat > 8 { Ovflorpat = 0 }
				fmt.Printf("cmd: f - floorp: %d\n",Ovflorpat)
			case 70:
				Ovflorpat -= 1
				if Ovflorpat < 0 { Ovflorpat = 8 }
				fmt.Printf("cmd: F - floorp: %d\n",Ovflorpat)
			case 103:		// g
				Ovflorcol += 1
				if anum >= 0 { Ovflorcol = anum }
				if Ovflorcol > 15 { Ovflorcol = 0 }
				fmt.Printf("cmd: g - floorc: %d\n",Ovflorcol)
			case 71:
				Ovflorcol -= 1
				if Ovflorcol < 0 { Ovflorcol = 15 }
				fmt.Printf("cmd: G - floorc: %d\n",Ovflorcol)
			case 114:		// r
				opts.MRP = true
				opts.MRM = false
				fmt.Printf("cmd: r - mr+: %t mr-: %t\n",opts.MRP,opts.MRM)
			case 82:		// R
				opts.MRP = false
				opts.MRM = true
				fmt.Printf("cmd: R - mr+: %t mr-: %t\n",opts.MRP,opts.MRM)
			case 116:		// t
				opts.MRP = false
				opts.MRM = false
				fmt.Printf("cmd: t - mr+: %t mr-: %t\n",opts.MRP,opts.MRM)
			case 109:		// m
				opts.MV = !opts.MV
				fmt.Printf("cmd: m - mv: %t\n",opts.MV)
			case 104:		// h
				opts.MH = !opts.MH
				fmt.Printf("cmd: h - mh: %t\n",opts.MH)
			case 97:		// a
				noact = false
				Ovwallpat = -1
			case 105:		// i
				opts.Gtp = 1
				G1 = true
				G2 = false
			case 108:		// l
				opts.R14 = !opts.R14
			case 115:		// s
				opts.SP = !opts.SP
			case 117:		// u
				opts.Gtp = 2
				G1 = false
				G2 = true
			case 118:		// v
				lx := 116
				if G1 { lx = 126 }
				fmt.Printf("\n valid maze address for Gauntlet %d\nmaze   dec -    hex\n",opts.Gtp)
				for x := 0; x <= lx;x ++ {
					ad := slapsticMazeGetRealAddr(x)
					fmt.Printf("%03d:%d - x%X  ",x,ad,ad)
					if (x + 1) % 7 == 0 { fmt.Printf("\n") }
				}
				fmt.Printf("\n")
			case 63:
				fmt.Printf("single letter commands\n\n? - this list\nq - quit program\nf - floor pattern+\nF - floor pattern-\n")
				fmt.Printf("g - floor color+\nG - floor color-\nw - wall pattern+\nW - wall pattern-\n")
				fmt.Printf("e - wall color+\nE - wall color-\nr - rotate maze +90°\nR - rotate maze -90°\n")
				fmt.Printf("t - turn off rotate\nh - mirror maze horizontal toggle\nm - mirror maze vertical toggle\ns - toggle rnd special potion\n")
				fmt.Printf("i - gauntlet mazes r1 - r9\nl - use gauntlet rev 14\nu - gauntlet 2 mazes\n")
				fmt.Printf("{n}a - load maze 1 - 127 g1, 1 - 117 g2, or address 229376 - 262143\n")
				fmt.Printf("v - valid address list                      * note some address will cause crash out\n")
				fmt.Printf("    commands can be chained - i.e. i5a will switch to g1 and load maze 5\n")
				fmt.Printf("G%d ",opts.Gtp)
				if opts.R14 { fmt.Printf("(r14)")
				} else { fmt.Printf("(r1-9)") }
				fmt.Printf(" Command (?, q, fFgG, wWeE, rRt, hm, s, il, u, v, #a): ")
				noact = true
			default:
				if ascii < 48 || ascii > 57 { fmt.Printf("unk: %d\n",ascii) }
			}
			if ascii > 47 && ascii < 58 {
				noact = true
				bascii, _ := strconv.Atoi(string(ascii))
				if anum < 0 {
					anum = bascii
				} else {
					anum = (anum * 10) + bascii
				}
			}
//fmt.Printf("ascii: %d\n",ascii)

		}
}
package main

import (
	"os"
	"regexp"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	MV bool   `long:"mv" description:"maze: mirror vertical"`
	MH bool   `long:"mh" description:"maze: mirror horizontal ( --mv --mh will rotate 180 )"`
	MRP bool  `long:"mrp" description:"maze: rotate +90"`
	MRM bool  `long:"mrm" description:"maze: rotate -90"`
	AddrG2  int    `long:"ad" default:"0" base:"10" description:"G2 address override (in dec)"`
	Animate bool   `short:"a" long:"animate" description:"Animate monster"`
	PalType string `long:"pt" default:"base" description:"Palette type"`
	PalNum  int    `long:"pn" default:"0" base:"16" description:"Palette number (in hex)"`
	Tile    int    `short:"t" long:"tile" base:"16" description:"Tile number to render (in hex)"`
	DimX    int    `short:"x" default:"2" description:"X dimension, in tiles"`
	DimY    int    `short:"y" default:"2" description:"Y dimension, in tiles"`
	Output  string `short:"o" long:"output" default:"output.png" description:"Output file"`
	Monster string `short:"m" long:"monster" description:"Monster to render"`
	Floor   int    `short:"f" long:"floor" default:"-1" base:"16" description:"Floor stamp to render (in hex)"`
	Wall    int    `short:"w" long:"wall" default:"-1" base:"16" description:"Wall stamp to render (in hex)"`
	Verbose bool   `short:"v" long:"verbose"`
}

const (
	TypeNone = iota
	TypeMonster
	TypeFloor
	TypeWall
	TypeItem
	TypeMaze
)

var runType = TypeNone

var reMonsters = regexp.MustCompile(`^(ghost)`)
var reFloor = regexp.MustCompile(`^(floor)`)
var reWall = regexp.MustCompile(`^(wall)`)
var reItem = regexp.MustCompile(`^(item)`)
var reMaze = regexp.MustCompile(`^(maze)`)

func gexinit() []string {
	args, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
		// check(err)
	}

	if len(args) > 0 {
		switch {
		case reMonsters.MatchString(args[0]):
			runType = TypeMonster
		case reFloor.MatchString(args[0]):
			runType = TypeFloor
		case reWall.MatchString(args[0]):
			runType = TypeWall
		case reItem.MatchString(args[0]):
			runType = TypeItem
		case reMaze.MatchString(args[0]):
			runType = TypeMaze
		}
	}

	return args
}

package main

import (
	"os"
	"regexp"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	mnum	int    `description:"maze: number"`
	edat	int    `default:"0" description:"maze: edit active"`
// cli options to mirror & rotate mazes
	MV		bool   `long:"mv" description:"maze: mirror vertical"`
	MH		bool   `long:"mh" description:"maze: mirror horizontal ( --mv --mh will rotate 180 )"`
	MRP		bool   `long:"mrp" description:"maze: rotate +90"`
	MRM		bool   `long:"mrm" description:"maze: rotate -90"`
// option to randomly inject special potions and gold bags
	SP		bool   `short:"s" long:"sp" description:"random insert special potions, gold bags"`
// option to load mask for batch
	Mask	int    `long:"mask" default:"0" base:"16" description:"mask to hide elements"`
	Nob		bool   `long:"nb" description:"no border around outer wall "`
// cli option to force an address (originally for g2 force)
	Addr	int    `long:"ad" default:"0" base:"16" description:"load maze rom address x38000 to x3FFFF (in hex)"`
// cli option to use rev 14 maze roms (final release, differs from identical maze roms in r1 - r9)
	R14		bool   `long:"r14" description:"use gauntlet rev14 maze rom"`
// select gauntlet 1 or 2 to process - default is 2
	Gtp 	int    `short:"g" long:"gtp" default:"1" base:"10" description:"Gauntlet to process, 1 or 2"`
	Se		bool   `short:"z" description:"sanctuary engine data output"`
// interactive mode for maze display, select wall & floors, rotates & mirrors, load new mazes, test addresses
// only with maze{n}, if -i not given, prog just exits with maze in output.png
	Intr	bool   `short:"i" description:"maze interactive cli mode and following parms"`
	Geow float64   `long:"xw" default:"1024" description:"window width in pixels"`
	Geoh float64   `long:"xh" default:"1024" description:"window height in pixels"`
// orig options
	Animate bool   `short:"a" long:"animate" description:"Animate monster"`
	PalType string `long:"pt" default:"base" description:"Palette type"`
	PalNum  int    `long:"pn" default:"0" base:"16" description:"Palette number (in hex)"`
	Tile    int    `short:"t" long:"tile" base:"16" description:"Tile number to render (in hex)"`
	DimX    int    `short:"x" default:"6" description:"X dimension, in tiles"`
	DimY    int    `short:"y" default:"6" description:"Y dimension, in tiles"`
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

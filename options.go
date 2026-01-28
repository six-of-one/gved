package main

import (
	"os"
	"regexp"
	"fmt"
	"math"
	"io/ioutil"
	"github.com/jessevdk/go-flags"
	"strings"
	"bufio"
)

var opts struct {
	mnum	int    `description:"maze: number"`
	edat	int    `default:"0" description:"maze: edit active"`
	dtec float64   `default:"16.0" description:"edit: tile size detector for click"`
	edip	int    `default:"0" description:"maze: last load from file"`
	dntr	bool   `description:"dont reload maze ebuf on a refresh"`
	bufdrt	bool   `description:"maze ebuf is unsaved"`
// cli options to mirror & rotate mazes
	MV		bool   `long:"mv" description:"maze: mirror vertical"`
	MH		bool   `long:"mh" description:"maze: mirror horizontal ( --mv --mh will rotate 180 )"`
	MRP		bool   `long:"mrp" description:"maze: rotate +90"`
	MRM		bool   `long:"mrm" description:"maze: rotate -90"`
// option to randomly inject special potions and gold bags
	SP		bool   `short:"s" long:"sp" description:"random insert special potions, gold bags"`
// option to load mask for batch
	Mask	int    `long:"mask" default:"0" base:"16" description:"mask to hide elements"`
// draw a tile size border around maze for horiz & vert loop indicator arrows - must be turned OFF for edit mode
	Aob		bool   `long:"ab" description:"arrow border around outer wall"`
	Wob		bool   `long:"wb" description:"extra walls border maze right and bottom"`
	Nogtop	bool   `long:"ngt" description:"no generator indicate letter"`
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
	SDef    bool   `long:"def" description:"set config defaults on all interactive parms"`
	Geow float64   `long:"xw" default:"1060" description:"window width in pixels"`
	Geoh float64   `long:"xh" default:"1086" description:"window height in pixels, maze disp is 26 px less"`
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
	Mute    bool   `long:"mute" description:"Mute all Audio"`
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

func gevinit() []string {
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

// adding config saver/loader

func ld_config() {

	var geow float64
	var geoh float64
	geow = 1060.0
	geoh = 1086.0
	viewp = 21
// load data attempt, leaves defaults loaded if fails
	fil := ".config"
	if opts.SDef { fil = ".config-def" }	// special config with default restore
	data, err := ioutil.ReadFile(fil)
	if err == nil {
		dscan := fmt.Sprintf("%s",data)
		scanr := bufio.NewScanner(strings.NewReader(dscan))
		l := "1060 1086 23"
		if scanr.Scan() { l = scanr.Text() }
		fmt.Sscanf(string(l),"%v %v %d", &geow, &geoh,&viewp)		// win size w * h, viewport size
		for i := 27; i <= 107; i += 20 {
			if scanr.Scan() { l = scanr.Text()						// edit keys g1
				fmt.Sscanf(l,"%d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d \n",&g1edit_keymap[i+0], &g1edit_keymap[i+1], &g1edit_keymap[i+2], &g1edit_keymap[i+3],
						&g1edit_keymap[i+4], &g1edit_keymap[i+5], &g1edit_keymap[i+6], &g1edit_keymap[i+7], &g1edit_keymap[i+8], &g1edit_keymap[i+9], &g1edit_keymap[i+10],&g1edit_keymap[i+11],
						&g1edit_keymap[i+12], &g1edit_keymap[i+13], &g1edit_keymap[i+14], &g1edit_keymap[i+15], &g1edit_keymap[i+16], &g1edit_keymap[i+17], &g1edit_keymap[i+18], &g1edit_keymap[i+19]) }
		}
		for i := 27; i <= 107; i += 20 {
			if scanr.Scan() { l = scanr.Text()						// edit keys g2
						//        0    1                        6                            12                       17                                 24                       29             32
				fmt.Sscanf(l,"%d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d \n",&g2edit_keymap[i+0], &g2edit_keymap[i+1], &g2edit_keymap[i+2], &g2edit_keymap[i+3],
						&g2edit_keymap[i+4], &g2edit_keymap[i+5], &g2edit_keymap[i+6], &g2edit_keymap[i+7], &g2edit_keymap[i+8], &g2edit_keymap[i+9], &g2edit_keymap[i+10],&g2edit_keymap[i+11],
						&g2edit_keymap[i+12], &g2edit_keymap[i+13], &g2edit_keymap[i+14], &g2edit_keymap[i+15], &g2edit_keymap[i+16], &g2edit_keymap[i+17], &g2edit_keymap[i+18], &g2edit_keymap[i+19]) }
		}
		l = "82cd00cd"
		if scanr.Scan() { l = scanr.Text() }						// blotter color & alpha
		fmt.Sscanf(l,"%08x",&blotcol)
		l = ""
		if scanr.Scan() { l = scanr.Text() }						// blotter replacement image
		fmt.Sscanf(l,"%s",&blotimg)
	}
// get default win size
// main win size is a bit tricky on user adjust as i cheaped out and made click detect of a cell
// - based on a square pixel block (default here is 32x32 pix, and smallest is 16x16 pix)
// - whats more, palette and paste buf will auto-readjust back to the size detected for the main win
// - mainly because i only have one click detect routine, and i wanted a square to keep the math ops less ugly
// - additionally the geometry captured in .config is the maze edit area, the total win is slightly larger
// - the minimums set by the .Max() are based on no shrinkage below min edit size of 16x16 pixel cell
// so a user can make the screen any rectangle they want, but going into edit mode will force sqaure cells!
// - even more obtuse, going 'forced fullscreen' wont play nice (and go doesnt have prog adjustable win stuff...)

	if opts.Geow == 1060 && opts.Geoh == 1086 {		// defs on entry, load from cfg

			opts.Geow = math.Max(560,geow)
			opts.Geoh = math.Max(586,geoh)
fmt.Printf("Load window size: %v x %v\n",geow,geoh)
	}

// do a save back
	sv_config()
}

func sv_config() {

// save stat
	file, err := os.Create(".config")
	if err == nil {
		wfs := fmt.Sprintf("%d %d %d\n",int(opts.Geow),int(opts.Geoh),viewp)
		k1, k2, k3, k4, k5 := "","","","",""
		for i := 27; i < 46; i++ {
			k1 += fmt.Sprintf("%d ",g1edit_keymap[i])
			k2 += fmt.Sprintf("%d ",g1edit_keymap[i+20])
			k3 += fmt.Sprintf("%d ",g1edit_keymap[i+40])
			k4 += fmt.Sprintf("%d ",g1edit_keymap[i+60])
			k5 += fmt.Sprintf("%d ",g1edit_keymap[i+80])
		}
		k1 += "\n"; k2 += "\n"; k3 += "\n"; k4 += "\n"; k5 += "\n"
		wfs += k1+k2+k3+k4+k5
		k1, k2, k3, k4, k5 = "","","","",""
		for i := 27; i < 46; i++ {
			k1 += fmt.Sprintf("%d ",g2edit_keymap[i])
			k2 += fmt.Sprintf("%d ",g2edit_keymap[i+20])
			k3 += fmt.Sprintf("%d ",g2edit_keymap[i+40])
			k4 += fmt.Sprintf("%d ",g2edit_keymap[i+60])
			k5 += fmt.Sprintf("%d ",g2edit_keymap[i+80])
		}
		k1 += "\n"; k2 += "\n"; k3 += "\n"; k4 += "\n"; k5 += "\n"
		wfs += k1+k2+k3+k4+k5
//fmt.Print(wfs)
fmt.Printf("sv_config\n")
		wfs += fmt.Sprintf("%08x\n%s\n",blotcol,blotimg)
		file.WriteString(wfs)
		file.Close()
//	fmt.Printf("saving .wstats file\n")
	}
}

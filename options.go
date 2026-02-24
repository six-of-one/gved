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

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/dialog"
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
	Nosec	bool   `long:"nsec" description:"no secret walls shown"`
// cli option to force an address (originally for g2 force)
	Addr	int    `long:"ad" default:"0" base:"16" description:"load maze rom address x38000 to x3FFF0 (in hex)"`
// cli option to use rev 14 maze roms (final release, differs from identical maze roms in r1 - r9)
	R14		bool   `long:"r14" description:"use gauntlet rev14 maze rom"`
// select gauntlet 1 or 2 to process - default is 2
	Gtp 	int    `short:"g" long:"gtp" default:"1" base:"10" description:"Gauntlet to process, 1 or 2"`
	Lvl		int	   `short:"l" long:"level" default:"0" base:"10" description:"Level start 1 - 8, def (or -1 to force) = random"`
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
// g1 ed keymap read
		for i := 27; i <= 107; i += 20 {
			if scanr.Scan() { l = scanr.Text()						// edit keys g1
				fmt.Sscanf(l,"%d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d \n",&g1edit_keymap[i+0], &g1edit_keymap[i+1], &g1edit_keymap[i+2], &g1edit_keymap[i+3],
						&g1edit_keymap[i+4], &g1edit_keymap[i+5], &g1edit_keymap[i+6], &g1edit_keymap[i+7], &g1edit_keymap[i+8], &g1edit_keymap[i+9], &g1edit_keymap[i+10],&g1edit_keymap[i+11],
						&g1edit_keymap[i+12], &g1edit_keymap[i+13], &g1edit_keymap[i+14], &g1edit_keymap[i+15], &g1edit_keymap[i+16], &g1edit_keymap[i+17], &g1edit_keymap[i+18], &g1edit_keymap[i+19]) }
		}
// g1 xb keymap read
		for i := 27; i <= 107; i += 20 {
			if scanr.Scan() { l = scanr.Text()						// xb keys g1
				fmt.Sscanf(l,"%s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s %s \n",&g1edit_xbmap[i+0], &g1edit_xbmap[i+1], &g1edit_xbmap[i+2], &g1edit_xbmap[i+3],
						&g1edit_xbmap[i+4], &g1edit_xbmap[i+5], &g1edit_xbmap[i+6], &g1edit_xbmap[i+7], &g1edit_xbmap[i+8], &g1edit_xbmap[i+9], &g1edit_xbmap[i+10],&g1edit_xbmap[i+11],
						&g1edit_xbmap[i+12], &g1edit_xbmap[i+13], &g1edit_xbmap[i+14], &g1edit_xbmap[i+15], &g1edit_xbmap[i+16], &g1edit_xbmap[i+17], &g1edit_xbmap[i+18], &g1edit_xbmap[i+19]) }
		}
// g2 ed keymap read
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
		l = "0 1.0"
		if scanr.Scan() { l = scanr.Text() }						// difficulty skill - affects rnd loader
		fmt.Sscanf(l,"%d %f",&opts.Lvl,&diff_level)
		l = "false, false"
		if scanr.Scan() { l = scanr.Text() }						// difficulty skill - affects rnd loader
		fmt.Sscanf(l,"%t, %t\n",&unpinx,&unpiny)
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
// g1 ed keymap
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
// g1 xb keymap
		k1, k2, k3, k4, k5 = "","","","",""
		for i := 27; i < 46; i++ {
			k1 += fmt.Sprintf("%s ",g1edit_xbmap[i])
			k2 += fmt.Sprintf("%s ",g1edit_xbmap[i+20])
			k3 += fmt.Sprintf("%s ",g1edit_xbmap[i+40])
			k4 += fmt.Sprintf("%s ",g1edit_xbmap[i+60])
			k5 += fmt.Sprintf("%s ",g1edit_xbmap[i+80])
		}
		k1 += "\n"; k2 += "\n"; k3 += "\n"; k4 += "\n"; k5 += "\n"
		wfs += k1+k2+k3+k4+k5
// g2 ed keymap
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
		wfs += fmt.Sprintf("%d %.1f\n",opts.Lvl,diff_level)
		wfs += fmt.Sprintf("%t, %t\n",unpinx,unpiny)
		file.WriteString(wfs)
		file.Close()
//	fmt.Printf("saving .wstats file\n")
	}
}

// options dialog
// testing

var optht float32 = 36.0

func optCont(wn fyne.Window) fyne.CanvasObject {

// viewport size
	vp_label := widget.NewLabelWithStyle("View size:       ", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
	vp_entr := widget.NewEntry()
	vp_entr.Resize(fyne.Size{70, optht})
	vp := fmt.Sprintf("%d",viewp)
	vp_entr.SetText(vp)
	vp_entr.OnChanged = func(s string) {
		fmt.Sscanf(s,"%d",&viewp)					// force a canvas refresh here if visible
		if diff_level < 0 { diff_level = 1.0 }
		ns := fmt.Sprintf("Viewport size: %d",viewp)
		statlin(cmdhin,ns)
		sv_config()
	}

// set difficulty level
	diff_label := widget.NewLabelWithStyle("Difficulty:      ", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})

	diff_entr := widget.NewEntry()
	diff_entr.Resize(fyne.Size{70, optht})
	ns := fmt.Sprintf("%.1f",diff_level)
	diff_entr.SetText(ns)
	diff_entr.OnChanged = func(s string) {
		fmt.Sscanf(s,"%f",&diff_level)
		if diff_level < 0 { diff_level = 1.0 }
		ns := fmt.Sprintf("Difficulty: %.1f",diff_level)
		statlin(cmdhin,ns)
		sv_config()
	}

// select start level
	sellvl := widget.NewSelect([]string{"Research", "Level 1", "Level 2", "Level 3", "Level 4", "Level 5", "Level 6", "Level 7", "Level 8...", "Random 1-7"}, func(str string) {
		fmt.Printf("Select level: %s\n", str)
		opts.Lvl = lvl_sel[str]
		sv_config()
	})
	sellvl.SetSelected(lvl_str[opts.Lvl])
	sel_label := widget.NewLabelWithStyle("Start on:", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})

// unpin  controls
	unpx := widget.NewCheck("Unpin X", func(upx bool) {
		fmt.Printf("Unpin X set to %t\n", upx)
		unpinx = upx
		sv_config()
	})
	unpx.Checked = unpinx
	unpy := widget.NewCheck("Unpin Y", func(upy bool) {
		fmt.Printf("Unpin Y set to %t\n", upy)
		unpiny = upy
		sv_config()
	})
	unpy.Checked = unpiny

// enable tutorial messages / announcements: "use keys to open doors"
	tut_label := widget.NewLabelWithStyle("Tutorials: ", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
	g1tut := widget.NewCheck("G¹ ", func(g1t bool) {
		fmt.Printf("Tutorial show / announce G¹ %t\n", g1t)
	//	g1ton = g1t
		sv_config()
	})
	g2tut := widget.NewCheck("G² ", func(g2t bool) {
		fmt.Printf("Tutorial show / announce G² %t\n", g2t)
	//	g2ton = g2t
		sv_config()
	})
	setut := widget.NewCheck("Sanctuary", func(se_t bool) {
		fmt.Printf("Tutorial show Sanctuary %t\n", se_t)
	//	se_ton = se_t
		sv_config()
	})
// game progress (similar to arcade but with a couple extra items)
//controls progression of maze/game enhancements - default settings are sanctuary engine enhanced
/* from sanctuary options sheet
value="0.7" select max % to cap all of these might happen tests -- 10% = 0.1, set to zero 0 = turn ALL off

above (+%) or (-%) below MID skill adjust base % by indicated %
mid skill value="5" middle diffculty for % increment / decrement of test value
value="0.1" decr by % per skill under MID">% value="0.05" incr by % per skill over MID - skill over 9 counts as 9

per level count * %mod added to base % value="0.02" incr by % per level count %

first level these maze progression will be tested and base percent to occur - 25% = 0.25 begins @level
---------------------------------------------------------------------------
mazes will be mirrored in X coordinate, starting on level at base %, set any % to 0 to stop that item
value="10" select level to start rnd mirrors (Horiz) : value="0.3" base % to mirror (X axis) %
mazes will be flipped in Y coordinate, starting on level at base %
value="20" select level to start rnd flips (Vert) : value="0.25" base % to flip (Y axis) %
mazes will be rotated 270 degrees: Se, starting on level at base %
value="30" select level to start rnd rotates value="0.15" base % to rotate (270°) %
mazes will be unpinned from X and/or Y zero line, starting on level at base %">unpin:</td></tr>
value="50" select level to start rnd unpins : value="0.25" base % to unpin X or Y %
mazes will have shots stun other players, starting on level at base %
value="24" select level to start shots stun other players : value="0.2" base % players shots stun players %
mazes will have shots hurt other players, starting on level at base %
value="33" select level to start shots hurt other players : value="0.16" base % players shots hurt players %

all std maze walls will be set invisible - rnd chance: Se, starting on level at base %
value="80" randomly turn some mazes (all std walls) invisible - this is a sanctuary enhance, G2 had perm set invis walls
value="0.1" base % all std maze walls turn invisible
*/

	prog_label := widget.NewLabelWithStyle("Level to start check : % chance: ", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
// note set / save val on these
	pmir_lab :=  widget.NewLabelWithStyle("Maze mirror (X): ", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
	pflip_lab := widget.NewLabelWithStyle("Maze flip (Y):   ", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
	prot_lab :=  widget.NewLabelWithStyle("Maze rot -/+90°: ", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
	unp_lab :=   widget.NewLabelWithStyle("Unpin edges:     ", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
	shst_lab :=  widget.NewLabelWithStyle("Shots stun:      ", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
	shhr_lab :=  widget.NewLabelWithStyle("Shots hurt:      ", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
	invw_lab :=  widget.NewLabelWithStyle("Invisible walls: ", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
	lcol :=  widget.NewLabelWithStyle("  :", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
	lper :=  widget.NewLabelWithStyle(" %", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})

// level the check will start
	prog_mirror := widget.NewEntry()
	prog_mirror.Resize(fyne.Size{60, optht})
	prog_flip := widget.NewEntry()
	prog_flip.Resize(fyne.Size{60, optht})
	prog_rot := widget.NewEntry()
	prog_rot.Resize(fyne.Size{60, optht})
	prog_unpin := widget.NewEntry()
	prog_unpin.Resize(fyne.Size{60, optht})
	prog_sstun := widget.NewEntry()
	prog_sstun.Resize(fyne.Size{60, optht})
	prog_shurt := widget.NewEntry()
	prog_shurt.Resize(fyne.Size{60, optht})
	prog_invw := widget.NewEntry()
	prog_invw.Resize(fyne.Size{60, optht})
// inital percent of check
	per_mirror := widget.NewEntry()
	per_mirror.Resize(fyne.Size{50, optht})
	per_flip := widget.NewEntry()
	per_flip.Resize(fyne.Size{50, optht})
	per_rot := widget.NewEntry()
	per_rot.Resize(fyne.Size{50, optht})
	per_unpin := widget.NewEntry()
	per_unpin.Resize(fyne.Size{50, optht})
	per_sstun := widget.NewEntry()
	per_sstun.Resize(fyne.Size{50, optht})
	per_shurt := widget.NewEntry()
	per_shurt.Resize(fyne.Size{50, optht})
	per_invw := widget.NewEntry()
	per_invw.Resize(fyne.Size{50, optht})

// maze config - moved from file menu for now
	keepdec := widget.NewCheck("keep decor", func(kp bool) {
		fmt.Printf("Keep decore %t\n", kp)
	})
	dim_y := widget.NewEntry()
	dim_y.Resize(fyne.Size{50, optht})
	dim_y.SetText(fmt.Sprintf("%d",opts.DimY))
	dim_x := widget.NewEntry()
	dim_x.Resize(fyne.Size{50, optht})
	dim_x.SetText(fmt.Sprintf("%d",opts.DimX))
	lblnk :=  widget.NewLabelWithStyle("   dims: ", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
	lbby  :=  widget.NewLabelWithStyle(" by ", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})

	blnk_maz := widget.NewButton("Blank maze", func() {
		if opts.edat == 0 {
			dialog.ShowInformation("G¹G²ved", "Edit mode is off                                           \nturn on edit mode with <ESC>                \n"+
			"or by middle mouse click on maze item\n\nBlank current maze load to X by Y dims\ncheck 'keep decore' to save wall/floor style", wn)
		} else {
			fmt.Sscanf(dim_x.Text,"%d",&opts.DimX)
			fmt.Sscanf(dim_y.Text,"%d",&opts.DimY)
fmt.Printf("blnk x,y %d x %d from dims: %s x %s\n", opts.DimX,opts.DimY,dim_x.Text,dim_y.Text)
			menu_blank(keepdec.Checked)
		}
	})
	allwals := widget.NewButton("Walls maze", func() {
		if opts.edat == 0 {
			dialog.ShowInformation("G¹G²ved", "Edit mode is off                                           \nturn on edit mode with <ESC>                \n"+
			"or middle mouse click on maze item\n\nMake a maze of all walls of X by X dims\ncheck 'keep decore' to save wall/floor style", wn)
		} else {
			fmt.Sscanf(dim_x.Text,"%d",&opts.DimX)
			fmt.Sscanf(dim_y.Text,"%d",&opts.DimY)
fmt.Printf("blnk x,y %d x %d from dims: %s x %s\n", opts.DimX,opts.DimY,dim_x.Text,dim_y.Text)
			menu_walls(keepdec.Checked)
		}
	})
	rnd_lod := widget.NewButton("Random item load", func() {
		if opts.edat == 0 {
			dialog.ShowInformation("G¹G²ved", "Edit mode is off                                           \nturn on edit mode with <ESC>                \n"+
			"or by middle mouse click on maze item\n\nLoad a profile of random items into current maze", wn)
		} else {
			rload(ebuf)
			ed_maze(true,1,1)
		}
	})
	reducwal := widget.NewButton("Reduct walls", func() {
		if opts.edat == 0 {
			dialog.ShowInformation("G¹G²ved", "Edit mode is off                                           \nturn on edit mode with <ESC>                \n"+
			"or by middle mouse click on maze item\n\nReduce walls of maze - good with Prim", wn)
		} else {
			ReduceWalls(ebuf,mxmd,mymd)
			ed_maze(true,1,1)
		}
	})
	mapfar := widget.NewButton("Mapper fargoal", func() {
		if opts.edat == 0 {
			dialog.ShowInformation("G¹G²ved", "Edit mode is off                                           \nturn on edit mode with <ESC>                \n"+
			"or by middle mouse click on maze item\n\nFargoal mapper system - std map", wn)
		} else {
			fmt.Sscanf(dim_x.Text,"%d",&opts.DimX)
			fmt.Sscanf(dim_y.Text,"%d",&opts.DimY)
			map_fargoal(ebuf)
			ed_maze(true,1,1)
		}
	})
	mapfar2 := widget.NewButton("Mapper sword", func() {
		if opts.edat == 0 {
			dialog.ShowInformation("G¹G²ved", "Edit mode is off                                           \nturn on edit mode with <ESC>                \n"+
			"or by middle mouse click on maze item\n\nFargoal mapper system - sword map", wn)
		} else {
			fmt.Sscanf(dim_x.Text,"%d",&opts.DimX)
			fmt.Sscanf(dim_y.Text,"%d",&opts.DimY)
			map_sword(ebuf)
			ed_maze(true,1,1)
		}
	})
	mapfar3 := widget.NewButton("Mapper wide", func() {
		if opts.edat == 0 {
			dialog.ShowInformation("G¹G²ved", "Edit mode is off                                           \nturn on edit mode with <ESC>                \n"+
			"or by middle mouse click on maze item\n\nFargoal mapper system - wide map", wn)
		} else {
			fmt.Sscanf(dim_x.Text,"%d",&opts.DimX)
			fmt.Sscanf(dim_y.Text,"%d",&opts.DimY)
			map_wide(ebuf)
			ed_maze(true,1,1)
		}
	})
	mapdfs := widget.NewButton("Mapper DFS", func() {
		if opts.edat == 0 {
			dialog.ShowInformation("G¹G²ved", "Edit mode is off                                           \nturn on edit mode with <ESC>                \n"+
			"or by middle mouse click on maze item\n\nDFS mapper algo", wn)
		} else {
			fmt.Sscanf(dim_x.Text,"%d",&opts.DimX)
			fmt.Sscanf(dim_y.Text,"%d",&opts.DimY)
			map_dfs(ebuf)
			ed_maze(true,1,1)
		}
	})
	mapprim := widget.NewButton("Mapper Prim", func() {
		if opts.edat == 0 {
			dialog.ShowInformation("G¹G²ved", "Edit mode is off                                           \nturn on edit mode with <ESC>                \n"+
			"or by middle mouse click on maze item\n\nPrim mapper algo - nice with reduce walls", wn)
		} else {
			fmt.Sscanf(dim_x.Text,"%d",&opts.DimX)
			fmt.Sscanf(dim_y.Text,"%d",&opts.DimY)
			GeneratePrimMaze(ebuf,mxmd,mymd)
			ed_maze(true,1,1)
		}
	})
/*
	menuItemRedwl := fyne.NewMenuItem("Reduct walls",func() { ReduceWalls(ebuf,mxmd,mymd); ed_maze(true,1,1) })
	menuItemFmap := fyne.NewMenuItem("Mapper fargoal",func() { map_fargoal(ebuf); ed_maze(true,1,1) })
	menuItemFmapb := fyne.NewMenuItem("Mapper 2",func() { map_sword(ebuf); ed_maze(true,1,1) })
	menuItemFmapc := fyne.NewMenuItem("Mapper 3",func() { map_wide(ebuf); ed_maze(true,1,1) })
	menuItemFmapd := fyne.NewMenuItem("Mapper DFS",func() { map_dfs(ebuf); ed_maze(true,1,1) })
	menuItemFmape := fyne.NewMenuItem("Mapper Prim",func() { GeneratePrimMaze(ebuf,mxmd,mymd); ed_maze(true,1,1) })
*/
// the massive block out
	opdlg := container.NewAppTabs(
	container.NewTabItemWithIcon("Game",theme.SettingsIcon(),			// stats VisibilityIcon, scores StorageIcon, edit FileApplicationIcon,
	container.New(
		layout.NewVBoxLayout(),
//		layout.NewSpacer(),
		container.New(
			layout.NewVBoxLayout(),
			layout.NewSpacer(),
			container.New(
				layout.NewHBoxLayout(),
				vp_label,
				container.NewWithoutLayout(vp_entr),
			),
			layout.NewSpacer(),
			container.New(
				layout.NewHBoxLayout(),
				diff_label,
				container.NewWithoutLayout(diff_entr),
			),
			layout.NewSpacer(),
			container.New(
				layout.NewHBoxLayout(),
				sel_label,
				sellvl,
			),
			layout.NewSpacer(),
			container.New(
				layout.NewHBoxLayout(),
				tut_label,
				g1tut, g2tut, setut,
			),
			layout.NewSpacer(),
			container.New(
				layout.NewVBoxLayout(),
				prog_label,
				container.New(
					layout.NewHBoxLayout(),
					pmir_lab,container.NewWithoutLayout(prog_mirror),lcol,container.NewWithoutLayout(per_mirror),lper,
				),
				container.New(
				layout.NewHBoxLayout(),
					pflip_lab,container.NewWithoutLayout(prog_flip),lcol,container.NewWithoutLayout(per_flip),lper,
				),
				container.New(
				layout.NewHBoxLayout(),
					prot_lab,container.NewWithoutLayout(prog_rot),lcol,container.NewWithoutLayout(per_rot),lper,
				),
				container.New(
				layout.NewHBoxLayout(),
					unp_lab,container.NewWithoutLayout(prog_unpin),lcol,container.NewWithoutLayout(per_unpin),lper,
				),
				container.New(
				layout.NewHBoxLayout(),
					shst_lab,container.NewWithoutLayout(prog_sstun),lcol,container.NewWithoutLayout(per_sstun),lper,
				),
				container.New(
				layout.NewHBoxLayout(),
					shhr_lab,container.NewWithoutLayout(prog_shurt),lcol,container.NewWithoutLayout(per_shurt),lper,
				),
				container.New(
				layout.NewHBoxLayout(),
					invw_lab,container.NewWithoutLayout(prog_invw),lcol,container.NewWithoutLayout(per_invw),lper,
				),
			),
		),
		layout.NewSpacer(),
	)),
	container.NewTabItemWithIcon("Dev",theme.WarningIcon(),
	container.New(
		layout.NewVBoxLayout(),
//		layout.NewSpacer(),
		container.New(
			layout.NewVBoxLayout(),
			layout.NewSpacer(),
			container.New(
				layout.NewHBoxLayout(),
				unpx, unpy,
			),
			layout.NewSpacer(),
			container.New(
				layout.NewHBoxLayout(),
				sel_label,
				sellvl,
			),
			layout.NewSpacer(),
		),
		layout.NewSpacer(),
	)),
	container.NewTabItemWithIcon("Edit",theme.FileApplicationIcon(),			// stats VisibilityIcon, scores StorageIcon, edit FileApplicationIcon,
	container.New(
		layout.NewVBoxLayout(),
//		layout.NewSpacer(),
		container.New(
			layout.NewVBoxLayout(),
			layout.NewSpacer(),
			container.New(
				layout.NewHBoxLayout(),
				vp_label,
				container.NewWithoutLayout(vp_entr),
			),
		),
		layout.NewSpacer(),
	)),
	container.NewTabItemWithIcon("Colors",theme.ColorPaletteIcon(),colorCont(wn),
	),
	container.NewTabItemWithIcon("Maze",theme.ViewFullScreenIcon(),
	container.New(
		layout.NewVBoxLayout(),
	container.New(
		layout.NewVBoxLayout(),
		layout.NewSpacer(),
		container.New(
			layout.NewHBoxLayout(),
			blnk_maz,lblnk,dim_x,lbby,dim_y,
		),
		container.New(
			layout.NewHBoxLayout(),
			allwals,
		),
		container.New(
			layout.NewHBoxLayout(),
			keepdec,
		),
		container.New(
			layout.NewHBoxLayout(),
			rnd_lod,
		),
		container.New(
			layout.NewHBoxLayout(),
			mapfar,
		),
		container.New(
			layout.NewHBoxLayout(),
			mapfar2,
		),
		container.New(
			layout.NewHBoxLayout(),
			mapfar3,
		),
		container.New(
			layout.NewHBoxLayout(),
			mapdfs,
		),
		container.New(
			layout.NewHBoxLayout(),
			mapprim,
		),
		container.New(
			layout.NewHBoxLayout(),
			reducwal,
		),

		layout.NewSpacer(),
	),
	)),
	)
	return opdlg
}
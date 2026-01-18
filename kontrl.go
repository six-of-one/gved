package main

import (
	"fmt"
	"os"
	"image"
	"image/color"
	"image/draw"
	"time"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
    "fyne.io/fyne/v2/canvas"
)

// kontrol is for keyboard, mouse & input management

// weirdo needsav exit condition and remaze, to allow dialog to block next action

var exitsel bool
var nsremaze bool
var swnothing int		// store nothing flag when cmd keys off

// input keys and keypress checks for canvas/ window
// since this is all that is called without other handler / timers
// - this is where maze update and edits will vector

// disable some cmd keys for edit mode click opts
var cmdoff bool

var anum int
var shift bool
var ctrl bool
var ctrlwon bool	// ctrl was on for sub wins - the win has to turn it off on ops
var del bool
var home bool
var logo bool		// other wise labeld "win" key
var ccp int			// cut copy and paste - current op
var edkey int		// for passing edit keys to clicker
var cmdhin string

// main keyboard handlers

// special keys not detected by keyrune handler

func specialKey() {
// handle keys down
	if deskCanvas, ok := w.Canvas().(desktop.Canvas); ok {
        deskCanvas.SetOnKeyDown(func(key *fyne.KeyEvent) {
//	fmt.Printf("Desktop key down: %h\n", key.Name)
			if key.Name == "BackSpace" {
				anum = (anum / 10);
				spx := ""
				if anum != 0 { spx = fmt.Sprintf("| numeric: %d", anum) }
				uptitl(opts.mnum, spx)
			}
			if key.Name == "Delete" { del = true; ccp_NOP() }
//			if key.Name == "BackSpace" { del = true }
			if key.Name == "Home" { home = true; if !wpalop { palete() }}
			if key.Name == "LeftSuper" { logo = true }
			if key.Name == "LeftShift" { shift = true }
			if key.Name == "RightShift" { shift = true }
			if key.Name == "LeftControl" { ctrl = true; ctrlwon = true }
			if key.Name == "RightControl" { ctrl = true; ctrlwon = true }
        })
// handle keys up
        deskCanvas.SetOnKeyUp(func(key *fyne.KeyEvent) {
//	fmt.Printf("Desktop key up: %v\n", key)
			srelod := false
			if key.Name == "Escape" {		// now toggle editor on/ off
				if opts.Aob { dialog.ShowInformation("Edit mode", "Error: can not edit with border around maze!", w) } else {
					if opts.edat == 0 {
						smod = "Edit mode: "
				fmt.Printf("editor on, maze: %03d\n",opts.mnum+1)
						opts.edat = 1
						stor_maz(opts.mnum+1)	// this does not auto store new edit mode to buffer save file, unless it creates the file
						statlin(cmdhin,"")
// these all deactivate as override during edit
						opts.MRM = false
						opts.MRP = false
						opts.MV = false
						opts.MH = false
						srelod = true
					} else {
						nsremaze = true
						srelod = needsav()
						if sdb == 0 { menu_lodit(true) }
						smod = "View mode: "
				fmt.Printf("editor off, maze: %03d\n",opts.mnum+1)
						cmdoff = false
						if swnothing > 0 { nothing = swnothing; swnothing = 0 }
						opts.edat = 0
						opts.dntr = false
						ccp_NOP()
						cmdhin = "cmds: ?, <ESC>, fFgG, wWqQ, rRt, hm, pPT, sL, S, il, u, v, A #a"
						statlin(cmdhin,"")
						Ovwallpat = -1
					}
				}
			}
			if key.Name == "Delete" { del = false }
//			if key.Name == "BackSpace" { del = false }
			if key.Name == "Home" { home = false }
			if key.Name == "LeftSuper" { logo = false }
			if key.Name == "LeftShift" { shift = false }
			if key.Name == "RightShift" { shift = false }
			if key.Name == "LeftControl" { ctrl = false }
			if key.Name == "RightControl" { ctrl = false }
			if key.Name == "S" && ctrl  { if shift { menu_savas() } else { menu_sav() }}
			if key.Name == "L" && ctrl  { if shift { menu_laodf() } else { menu_lod() }}
			if key.Name == "R" && ctrl  { menu_res() }
			if key.Name == "U" && ctrl  { uswap() }
			if key.Name == "Z" && ctrl  { undo() }
			if key.Name == "Y" && ctrl  { redo() }
			if key.Name == "C" && ctrl  { menu_copy() }
			if key.Name == "X" && ctrl  { menu_cut() }
			if key.Name == "P" && ctrl  { if shift { pbsess_cyc(1) } else { menu_paste() }}
			if key.Name == "O" && ctrl  { if shift { pbmas_cyc(1) }}
			if key.Name == "Q" && ctrl  { exitsel = true; needsav() }
			if key.Name == "Prior" {
				if sdb > 0 {
					sdbit(-1)
				} else { srelod = pagit(-1) }
			}
			if key.Name == "Next" {
				if sdb > 0 {
					sdbit(1)
				} else { srelod = pagit(1) }
			}
			upd_edmaze(false)
//fmt.Printf("sk cond relod: %t\n",srelod)
			if srelod {
				remaze(opts.mnum)
			}
       })
    }
}

// regular keys

func typedRune(r rune) {

// special aux string - put ops in title after maze #
	spau := ""
// relod
	relod := false
	relodsub := false

//	fmt.Printf("in keys event - %x\n",r)

// new maze
	if r == 'a' {
		if (anum > 0 && anum <= 127 || anum >= 229376 && anum < 262145) {

			nsremaze = true
			relod = needsav()
			if anum <= 127 {
				opts.mnum = anum - 1
				Aov = 0
			} else {
				Aov = addrver(anum, 0)
				opts.mnum = 0
				spau = fmt.Sprintf("addr = %d",anum)
			}
			anum = 0
// clear these when load new maze
			Ovwallpat = -1
		}
	}

// (almost) blind numeric input
	switch r {
	case '0':
		anum = (anum * 10)
	case '1':
		anum = (anum * 10) + 1
	case '2':
		anum = (anum * 10) + 2
	case '3':
		anum = (anum * 10) + 3
	case '4':
		anum = (anum * 10) + 4
	case '5':
		anum = (anum * 10) + 5
	case '6':
		anum = (anum * 10) + 6
	case '7':
		anum = (anum * 10) + 7
	case '8':
		anum = (anum * 10) + 8
	case '9':
		anum = (anum * 10) + 9
	case '`':
		anum = 0
	}

//fmt.Printf("r %v shift %v\n",r,shift)
		edkey = int(r)
		if cmdoff {
			if G1 {
//				if g1edit_keymap[edkey] < 0 { keyst := fmt.Sprintf("locked key: %s not usable",map_keymap[edkey]) }
				if g1edit_keymap[edkey] == 0 { keyst := fmt.Sprintf("G¹ free key: %s middle mouse click to set",map_keymap[edkey]); statlin(cmdhin,keyst) }
				if g1edit_keymap[edkey] > 0 {
					kys := g1mapid[g1edit_keymap[edkey]]
					keyst := fmt.Sprintf("G¹ ed key: %s = %03d, %s",map_keymap[edkey],g1edit_keymap[edkey],kys)
					ccp_NOP()
					statlin(cmdhin,keyst)
fmt.Printf("G¹ ed key: %d - %s\n",edkey,kys)
				}
			} else {
//				if g2edit_keymap[edkey] < 0 { keyst := fmt.Sprintf("locked key: %s not usable",map_keymap[edkey]) }
				if g2edit_keymap[edkey] == 0 { keyst := fmt.Sprintf("G² free key: %s middle mouse click to set",map_keymap[edkey]); statlin(cmdhin,keyst) }
				if g2edit_keymap[edkey] > 0 {
					kys := g2mapid[g2edit_keymap[edkey]]
					keyst := fmt.Sprintf("G² ed key: %s = %03d, %s",map_keymap[edkey],g2edit_keymap[edkey],kys)
					ccp_NOP()
					statlin(cmdhin,keyst)
				}
			}
		}
// keys that '\' doesnt block, no maze reloads
		relodsub = false
		switch r {
		case 92:		// \
			ska := "cmd keys mode"
			if opts.edat > 0 {			// have to be in editor to turn on edit keys
				cmdoff = !cmdoff
// a,e only lower case not avail for edit hotkey
				if cmdoff {
					ska = "edit keys mode"
					if nothing > 0 { swnothing = nothing; nothing = 0 }
				} else {
					if swnothing > 0 { nothing = swnothing; swnothing = 0 }
					ccp_NOP()
				}
				opts.dntr = true
				relod = true
			}
			cmdhin = "cmds: ?, <ESC>, fFgG, wWqQ, rRt, hm, pPT, sL, S, il, u, v, A #a"
			if opts.edat > 0 { cmdhin = "cmds: ?, <ESC>, '\\', fFgG, wWqQ, rRt, hm, pPT, sL, S, il, u, v, A #a" }
			if cmdoff && opts.edat > 0 { cmdhin = "cmds: ?, <ESC>, '\\' - edit keys, #c C, HV, A #a, L, S" }
			fmt.Printf("hint: %s\n", cmdhin)
			statlin(cmdhin,ska)
		case 63:		// ?
			keyhints()
		case 65:		// A
			if Aov > 0 { Aov = 0 } else {
				Aov = addrver(slapsticMazeGetRealAddr(opts.mnum), 0)
			}
		case 76:		// L
// with anum != 0, this becomes load s[1] buffer into ebuf, if in edit
			if opts.edat > 0 && anum > 0 && anum < sdmax {
fmt.Printf("L, anum: %05d, sdb: %d\n",anum, sdb)
				if opts.bufdrt { menu_savit(true) }		// autosave
				fil := fmt.Sprintf(".ed/sd%05d_g%d.ed",anum,opts.Gtp)
				cnd := lod_maz(fil, ebuf, false)
				if cnd >= 0 { sdb = anum; for y := 0; y < 11; y++ { eflg[y] =  tflg[y] }; ed_maze(true) }
				anum = 0
			} else { opts.Nogtop = !opts.Nogtop; opts.dntr = (opts.edat > 0); relod = true }
//		case 69:		// E
		case 'c':
			if anum > 0 && anum < 65 {
				cycl = anum - 1
				anum = 0
			} else { cycl--; cycl-- }
			fallthrough
		case 67:		// C
				cycl++
				if cycl > 64 { cycl = 0 }
				if cycl < 1 { cycl = 64 }		// cause 'c' falls thru here
				kys := "n/a"
				if G1 {
					g1edit_keymap[cycloc] = cycl
					kys = g1mapid[cycl]
				} else {
					g2edit_keymap[cycloc] = cycl
					kys = g2mapid[cycl]
				}
				fmt.Printf("cyc %d - %s\n",cycl,kys)
				stu := fmt.Sprintf("cyc key: %s = %03d, %s",map_keymap[cycloc],cycl,kys)
				statlin(cmdhin,stu)
				edkey = 99						// pre set store cycl when cycling
				relod = true					// needed to refresh indicate text
				opts.dntr = true				// ... but dont kill the ebuf
		case 72:		// H	- horiz wrap
				eflg[4] = eflg[4] ^ LFLAG4_WRAP_H
				opts.dntr = true
				relod = true
		case 86:		// V	- vert wrap
				eflg[4] = eflg[4] ^ LFLAG4_WRAP_V
				opts.dntr = true
				relod = true
//	fmt.Printf("4 flag: %d\n",eflg[4])
		case 83:		// S
// have anum !=0, save ebuf into that buffer
				if anum > 0 && anum < sdmax && opts.edat > 0 {
					fil := fmt.Sprintf(".ed/sd%05d_g%d.ed",anum,opts.Gtp)
					sav_maz(fil, ebuf, eflg, opts.DimX, opts.DimY, 0 - anum)
					anum = 0
				} else {
					spau = sdbit(1)
				}
		default:
			relodsub = false
		}
// view cmd keys - also on edit, but blockable
	  if !cmdoff || opts.edat < 1 {
		relodsub = true
		switch r {
//		case 'z':
//		case 'x':
		case 'w':
			Ovwallpat += 1
			if anum > 0 { Ovwallpat = anum - 1; anum = 0 }
			if Ovwallpat > 7 { Ovwallpat = 0 }
			eflg[5] = (Ovflorpat & 0x0f) << 4 + (Ovwallpat & 0x0f)
			spau = fmt.Sprintf("cmd: w - wallp: %d\n",Ovwallpat)
			opts.bufdrt = (opts.edat > 0)
			opts.dntr = (opts.edat > 0)
		case 87:		// W
			Ovwallpat -= 1
			if Ovwallpat < 0 { Ovwallpat = 7 }
			eflg[5] = (Ovflorpat & 0x0f) << 4 + (Ovwallpat & 0x0f)
			spau = fmt.Sprintf("cmd: w - wallp: %d\n",Ovwallpat)
			opts.bufdrt = (opts.edat > 0)
			opts.dntr = (opts.edat > 0)
		case 'q':
			Ovwallcol += 1
			if anum > 0 { Ovwallcol = anum - 1; anum = 0 }
			if Ovwallcol > 16 { Ovwallcol = 0 }
			eflg[6] = (Ovflorcol & 0x0f) << 4 + (Ovwallcol & 0x0f)
			spau = fmt.Sprintf("cmd: e - wallc: %d\n",Ovwallcol)
			opts.bufdrt = (opts.edat > 0)
			opts.dntr = (opts.edat > 0)
		case 81:		// Q
			Ovwallcol -= 1
			if Ovwallcol < 0 { Ovwallcol = 16 }
			eflg[6] = (Ovflorcol & 0x0f) << 4 + (Ovwallcol & 0x0f)
			spau = fmt.Sprintf("cmd: e - wallc: %d\n",Ovwallcol)
			opts.bufdrt = (opts.edat > 0)
			opts.dntr = (opts.edat > 0)
		case 'f':
			Ovflorpat += 1
			if anum > 0 { Ovflorpat = anum - 1; anum = 0 }
			if Ovflorpat > 8 { Ovflorpat = 0 }
			eflg[5] = (Ovflorpat & 0x0f) << 4 + (Ovwallpat & 0x0f)
			spau = fmt.Sprintf("cmd: f - floorp: %d\n",Ovflorpat)
			opts.bufdrt = (opts.edat > 0)
			opts.dntr = (opts.edat > 0)
		case 70:		// F
			Ovflorpat -= 1
			if Ovflorpat < 0 { Ovflorpat = 8 }
			eflg[5] = (Ovflorpat & 0x0f) << 4 + (Ovwallpat & 0x0f)
			spau = fmt.Sprintf("cmd: f - floorp: %d\n",Ovflorpat)
			opts.bufdrt = (opts.edat > 0)
			opts.dntr = (opts.edat > 0)
		case 'g':
			Ovflorcol += 1
			if anum > 0 { Ovflorcol = anum - 1; anum = 0 }
			if Ovflorcol > 15 { Ovflorcol = 0 }
			eflg[6] = (Ovflorcol & 0x0f) << 4 + (Ovwallcol & 0x0f)
			spau = fmt.Sprintf("cmd: g - floorc: %d\n",Ovflorcol)
			opts.bufdrt = (opts.edat > 0)
			opts.dntr = (opts.edat > 0)
		case 71:		// G
			Ovflorcol -= 1
			if Ovflorcol < 0 { Ovflorcol = 15 }
			eflg[6] = (Ovflorcol & 0x0f) << 4 + (Ovwallcol & 0x0f)
			spau = fmt.Sprintf("cmd: g - floorc: %d\n",Ovflorcol)
			opts.bufdrt = (opts.edat > 0)
			opts.dntr = (opts.edat > 0)
		case 'r':
			if opts.edat > 0 {
				opts.MRP = true
				upd_edmaze(false)
				rotmirbuf(edmaze)
				opts.dntr = true
				opts.bufdrt = true
			} else {
				opts.MRP = true
				opts.MRM = false
				spau = fmt.Sprintf("cmd: r - mr+: %t mr-: %t\n",opts.MRP,opts.MRM)
			}
		case 82:		// R
			if opts.edat > 0 {
				opts.MRM = true
				upd_edmaze(false)
				rotmirbuf(edmaze)
				opts.dntr = true
				opts.bufdrt = true
			} else {
				opts.MRP = false
				opts.MRM = true
				spau = fmt.Sprintf("cmd: r - mr+: %t mr-: %t\n",opts.MRP,opts.MRM)
			}
		case 't':
			if opts.edat == 0 {
				opts.MRP = false
				opts.MRM = false
				spau = fmt.Sprintf("cmd: t - mr+: %t mr-: %t\n",opts.MRP,opts.MRM)
			}
		case 'm':
			if opts.edat > 0 {
				opts.MV = true
				upd_edmaze(false)
				rotmirbuf(edmaze)
				opts.dntr = true
				opts.bufdrt = true
			} else {
				opts.MV = !opts.MV
				spau = fmt.Sprintf("cmd: m - mv: %t\n",opts.MV)
			}
		case 'h':
			if opts.edat > 0 {
				opts.MH = true
				upd_edmaze(false)
				rotmirbuf(edmaze)
				opts.dntr = true
				opts.bufdrt = true
			} else {
				opts.MH = !opts.MH
				spau = fmt.Sprintf("cmd: h - mh: %t\n",opts.MH)
			}
		case 'i':
			nsremaze = true
			relodsub = needsav()
			opts.Gtp = 1
			opts.R14 = false
			og2 := G2
			G1 = true
			G2 = false
			if og2 { get_pbcnt(); subsw() }
			maxmaze = 126
			if opts.mnum > maxmaze { opts.mnum = 114 }
			spau = "G¹ "
		case 'l':
			nsremaze = true
			relodsub = needsav()
			opts.Gtp = 1
			opts.R14 = !opts.R14
			og2 := G2
			G1 = true
			G2 = false
			if og2 { get_pbcnt(); subsw() }
			maxmaze = 126
			if opts.mnum > maxmaze { opts.mnum = 114 }
			spau = "G¹ "
		case 'p':
			if !cmdoff {
				nothing = nothing ^ NOFLOOR
				spau = fmt.Sprintf("no floors: %d\n",nothing & NOFLOOR)
				opts.dntr = true
			}
		case 80:		// P
			if !cmdoff {
				nothing = nothing ^ NOWALL
				spau = fmt.Sprintf("no walls: %d\n",nothing & NOWALL)
				opts.dntr = true
			}
		case 84:		// T
			if !cmdoff {
				nt := (nothing & 511) + 1
				nothing = (nothing & 1536) + (nt & 511)
				if anum > 0 { nothing = (nothing & 1536) + anum; anum = 0 }		// set lower 9 bits of no-thing mask [ but not walls or floors ]
				spau = fmt.Sprintf("no things: %d\n",nothing & 511)				// display no things mask
				opts.dntr = true
			}
		case 's':
			opts.SP = !opts.SP
			opts.dntr = true
		case 'u':
			nsremaze = true
			relodsub = needsav()
			opts.Gtp = 2
			og1 := G1
			G1 = false
			G2 = true
			if og1 { get_pbcnt(); subsw() }
			maxmaze = 116
			if opts.mnum > maxmaze { opts.mnum = 102 }
			spau = "G² mazes"
		case 'v':
			lx := 116
			if G1 { lx = 126 }
			fmt.Printf("\n valid maze address for Gauntlet %d\nmaze   dec -    hex\n",opts.Gtp)
			for x := 0; x <= lx;x ++ {
					ad := slapsticMazeGetRealAddr(x)
					fmt.Printf("%03d:%d - x%X  ",x + 1,ad,ad)
					if (x + 1) % 7 == 0 { fmt.Printf("\n") }
			}
			fmt.Printf("\n")
			dialog.ShowInformation("G¹G²ved", "Gauntlet / Gauntlet 2 valid maze address list\nplease check terminal where gved command was issued\n\ngithub.com/six-of-one/", w)
			opts.dntr = true
		default:
			relodsub = false
		}
	  }
	  upd_edmaze(false)		// store vars view changes like floors or walls
		if spau == "G¹ " {
			if opts.R14 { spau += "rv14" } else { spau += "rv1-9" }
		}
fmt.Printf("kr cond relod: %t\n",relod || relodsub)
	if (relod || relodsub) {
		remaze(opts.mnum)
	}
	spx := ""
	if sdb > 0 { spx = fmt.Sprintf("sdbuf: %d",sdb) }
	if anum != 0 { spx += fmt.Sprintf("| numeric: %d", anum) }
	uptitl(opts.mnum, spau + spx)
}

// data needing preserved by needsav - all this could be changed by the next op while dialog waits on user
var nsxd int
var nsyd int
var nsgg int
var nsmz int
var nssb int

func menu_ndsav(y bool) {
	if y {
		if nssb < 0 {
			fil := fmt.Sprintf(".ed/g%dmaze%03d.ed",nsgg,nsmz)
			sav_maz(fil, ebuf, eflg, nsxd, nsyd, nsmz)
		} else {
			fil := fmt.Sprintf(".ed/sd%05d_g%d.ed",nssb,nsgg)
			sav_maz(fil, ebuf, eflg, nsxd, nsyd, 0 - nssb)
		}
	}
	if exitsel { os.Exit(0) }
	if nsremaze { remaze(opts.mnum) }
}

func needsav() bool {
	ret := true
	if opts.bufdrt {
		nsxd = opts.DimX
		nsyd = opts.DimY
		nsgg = opts.Gtp
		nsmz = opts.mnum+1
		nssb = sdb

		dia := fmt.Sprintf("Save changes for maze %d in .ed/g%dmaze%03d.ed ?\n\nWARNING:\nif not saved, changes will be discarded",nsmz,nsgg,nsmz)
		if nssb >= 0 { dia = fmt.Sprintf("Save changes in buffer %d to .ed/sd%05d_g%d.ed ?\n\nWARNING:\nif not saved, changes will be discarded",nssb,nssb,nsgg) }
		dialog.ShowConfirm("Save?",dia, menu_ndsav, w)
//		opts.bufdrt = false;		// save clears this, clear here in case discard is selected
		ret = false
	} else { if exitsel { os.Exit(0) }}
	return ret
}

func undo() {
	if restak > 0 {
		restak--		// keep track of position in stack for buffer save
		revk := delbuf.revc[restak]	// revoke count - items in loops can undo/redo all at once
fmt.Printf(" stk %d undo %d elem: %d maze: %d x %d - rloop: %d\n",delstak,restak,delbuf.elem[restak],delbuf.mx[restak],delbuf.my[restak],delbuf.revc[restak])
		for revk > 0 && restak >= 0 {
			sw := ebuf[xy{delbuf.mx[restak], delbuf.my[restak]}]
			ebuf[xy{delbuf.mx[restak], delbuf.my[restak]}] = delbuf.elem[restak]
			delbuf.elem[restak] = sw
			revk--
			if revk > 0 && restak > 0 { restak-- }
		}
fmt.Printf(" del %d elem: %d\n",restak,delbuf.elem[restak])
		opts.bufdrt = true
		ed_maze(true)
	}
}

func redo() {
	if delbuf.elem[restak] >= 0 {
fmt.Printf(" stk %d redo %d elem: %d maze: %d x %d - rloop: %d\n",delstak,restak,delbuf.elem[restak],delbuf.mx[restak],delbuf.my[restak],delbuf.revc[restak])
		revk := delbuf.revc[restak]	// revoke count - items in loops can undo/redo all at once
		for revk > 0 && restak >= 0 {
			sw := ebuf[xy{delbuf.mx[restak], delbuf.my[restak]}]
			ebuf[xy{delbuf.mx[restak], delbuf.my[restak]}] = delbuf.elem[restak]
//fmt.Printf(" redo %d elem: %d maze: %d x %d - rloop: %d\n",restak,delbuf.elem[restak],delbuf.mx[restak],delbuf.my[restak],delbuf.revc[restak])
			delbuf.elem[restak] = sw
			revk++
			restak++
			if delbuf.elem[restak] < 0 || delbuf.revc[restak] == 1 { revk = 0}
		}
		opts.bufdrt = true
		ed_maze(true)
	}
//	ed_maze()
}

func uswap() {
	for y := 0; y <= opts.DimY; y++ {
		for x := 0; x <= opts.DimX; x++ {
			sw := ebuf[xy{x,y}]
			ebuf[xy{x,y}] = ubuf[xy{x,y}]
			ubuf[xy{x,y}] = sw
	}}
	for y := 0; y < 11; y++ { sw := eflg[y]; eflg[y] = uflg[y]; uflg[y] = sw }
// also have to swap delete stak
	udbck(delstak+1,delstak)
	su := udstak
	udstak = delstak
	delstak = su
	su = urstak
	urstak = restak
	restak = su
	for y := 0; y <= udstak; y++ {
		su = udb.mx[y]; udb.mx[y] = delbuf.mx[y]; delbuf.mx[y] = su
		su = udb.my[y]; udb.my[y] = delbuf.my[y]; delbuf.my[y] = su
		su = udb.revc[y]; udb.revc[y] = delbuf.revc[y]; delbuf.revc[y] = su
		su = udb.elem[y]; udb.elem[y] = delbuf.elem[y]; delbuf.elem[y] = su
//		if udb.elem[y] < 0 { delstak = y; break }
	}
	ed_maze(true)
}

// cut / copy / paste (c/c/p) controls

func ccp_NOP() { if ccp == PASTE { blotter(nil,0,0,0,0); blotoff() }; ccp = NOP; if opts.edat > 0 { smod = "Edit mode: "; statlin(cmdhin,"") }}
func ccp_tog(op int) { wccp := ccp; ccp_NOP(); if wccp != op { ccp = op }}

func pb_upd(id string, nt string, vl int) {
// clear old buf
	pmx := opts.DimX; pmy := opts.DimY		// preserve these
	for my := 0; my <= pmy; my++ {
	for mx := 0; mx <= pmx; mx++ { cpbuf[xy{mx, my}] = 0 }}
	fil := fmt.Sprintf(".pb/%s_%07d_g%d.ed",id,vl,opts.Gtp)
	lod_maz(fil, cpbuf, false)
	cpx = opts.DimX; cpy = opts.DimY
fmt.Printf("%spb dun: px %d py %d, %s\n",nt,cpx,cpy,fil)
	opts.DimX = pmx; opts.DimY = pmy
	bwin(cpx+1, cpy+1, vl, cpbuf, eflg)		// draw the buffer
	bl := fmt.Sprintf("paste buf: %d", vl)
	statlin(cmdhin,bl)
}

func pbsess_cyc(dr int) {

fmt.Printf("pbses c: %d %d\n",lpbcnt,sesbcnt)
	pb_upd("ses", "ses", sesbcnt)
	sesbcnt += dr
	if sesbcnt >= lpbcnt { sesbcnt = 1 }
	if sesbcnt < 1 { sesbcnt = lpbcnt - 1 }
	if sesbcnt < 1 { sesbcnt = 1 }
}

func pbmas_cyc(dr int) {

fmt.Printf("pbmas c: %d %d\n",pbcnt,masbcnt)
	pb_upd("pb", "mas", masbcnt)
	masbcnt += dr
	if masbcnt >= pbcnt { masbcnt = 1 }
	if masbcnt < 1 { masbcnt = pbcnt - 1 }
	if masbcnt < 1 { masbcnt = 1 }
}

// page thru maze #s, sd buf
// mb = 2 (right button) will call the active op with last dir
// possible: wheel mouse these one day

var pgdir int

func pagit(dir int) bool {

	lrelod := false
	nsremaze = true
	lrelod = needsav()
	Ovwallpat = -1
	pgdir = dir
	if Aov > 0 {
		nav := addrver(Aov, dir)
		Aov = nav
	} else {
		opts.mnum += dir
	}
	if dir > 0 {
		if opts.mnum > maxmaze { opts.mnum = 0 }
	} else {
		if opts.mnum < 0 { opts.mnum = maxmaze }
	}
fmt.Printf("pg %d\n",opts.mnum)
	return lrelod
}

func sdbit(dir int) string {
// rotate curr ebuf thru s[1] - s[?]
	cnd := -1
	ldb := sdb
	spar := ""
	if opts.bufdrt { menu_savit(true) }		// autosave
	for cnd < 0 && ldb < sdmax {
		pgdir = dir
		ldb += dir
		if ldb == 0 { ldb = 1 }
		fil := fmt.Sprintf(".ed/sd%05d_g%d.ed",ldb,opts.Gtp)
		cnd = lod_maz(fil, ebuf, false)
		if cnd >= 0 { sdb = ldb; for y := 0; y < 11; y++ { eflg[y] =  tflg[y] }; ed_maze(true); spar = fmt.Sprintf("cmd: S - ") }
		if dir < 0 && cnd < 0 && ldb == 1 { cnd = 0; break }
	}
	if cnd < 0 {
		menu_lodit(true)
	}
	return spar
}

// rubber banded

var blot *canvas.Image

func blotter(img *image.NRGBA,px float32, py float32, sx float32, sy float32) {

	if img == nil {
		img = image.NewNRGBA(image.Rect(0, 0, 1, 1))
		draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{R: 255, G: 0, B: 255, A: 180}}, image.ZP, draw.Src)
	}
	blot = canvas.NewImageFromImage(img)
	blot.Move(fyne.Position{px, py})
	blot.Resize(fyne.Size{sx, sy})
}

// turn off blotter after a window update
// because the window update...
// a. turns it on full maze for no reason
// b. refuses to turn it off, even with a delay in fn()

func blotoff() {
// restor after pb uses
	time.Sleep(5 * time.Millisecond)
	blot.Resize(fyne.Size{0, 0})
}
// click area for edits

// button we can detect click and release areas for rubberband area & fills
// title tells us what window the button is on, assigned on btn creation if win is titled

type holdableButton struct {
    widget.Button
	title string
}

func newHoldableButton() *holdableButton {

    button := &holdableButton{}
    button.ExtendBaseWidget(button)

	return button
}

// main mouse handler

// store x & y when mouse button goes down - to start rubberband area
// 		and when released for other ops like cup & paste
var sxmd float64
var symd float64
var exmd float64
var eymd float64
var mbd bool

// &{{{387 545} {379 509.92188}} 4 0}

func (h *holdableButton) MouseMoved(mm *desktop.MouseEvent){
	ax := 0.0       // absolute x & y
	ay := 0.0
	rx := 0.0
	ry := 0.0
	mb := 0         // mb 1 = left, 2 = right, 4 = middle
	mk := 0         // mod key 1 = sh, 2 = ctrl, 4 = alt, 8 = logo
	pos := fmt.Sprintf("%v",mm)
	fmt.Sscanf(pos,"&{{{%f %f} {%f %f}} %d %d",&ax,&ay,&rx,&ry,&mb,&mk)
	cwt = h.title	// current window title by btn establish
//fmt.Printf("a: %f x %f rp: %f x %f\n",ax,ay,rx,ry)
	dt := float32(opts.dtec)
	sx := float32(sxmd)
	sy := float32(symd)
	ex := float32(rx)
	ey := float32(ry)
//beef := fmt.Sprintf("a: %.2f x %.2f r: %.2f x %.2f dt: %.2f",sx,sy,ex,ey,dt)
//statlin(cmdhin,beef)

	if strings.Contains(h.title, "G¹G²ved") {		// only in main win
	exmd = rx			// so bwin can locate pb changes if drawn
	eymd = ry
	if ccp == PASTE {
//		ex = float32(float32(rx) + dt)
//		ey = float32(float32(ry) + dt)
		sx := float32(int(ex / dt)) * dt - 3
		sy := float32(int(ey / dt)) * dt - 3
		lx := float32(cpx) * dt + dt
		ly := float32(cpy) * dt + dt
		blot.Move(fyne.Position{sx, sy})
		if !wpbop { blot.Resize(fyne.Size{lx, ly}) }
beef := fmt.Sprintf("rbd o: %.2f x %.2f s: %.2f x %.2f dt: %.2f",sx,sy,lx,ly,dt)
statlin(cmdhin,beef)
	} else {
	if ex < sx { t := sx; sx = ex; ex = t }		// swap if end smaller than start
	if ey < sy { t := sy; sy = ey; ey = t }
	ex = float32(float32(ex) + dt)					// click in 1 tile selects the tile
	ey = float32(float32(ey) + dt)
	if mbd {
		sx = float32(int(sx / dt)) * dt - 3				// blotter selects tiles with original unit of 16 x 16
		sy = float32(int(sy / dt)) * dt - 4
		ex = float32(int(ex / dt)) * dt - 1
		ey = float32(int(ey / dt)) * dt - 2
		blot.Move(fyne.Position{sx, sy})
		blot.Resize(fyne.Size{ex - sx, ey - sy})
//		fmt.Printf("st: %f x %f pos: %f x %f\n",sx,sy,ex,ey)
	} else {
		blot.Resize(fyne.Size{0, 0})
	}}}
}

func (h *holdableButton) MouseDown(mm *desktop.MouseEvent){
	ax := 0.0	// absolute x & y
	ay := 0.0
	mb := 0		// mb 1 = left, 2 = right, 4 = middle
	mk := 0		// mod key 1 = sh, 2 = ctrl, 4 = alt, 8 = logo
	pos := fmt.Sprintf("%v",mm)
	fmt.Sscanf(pos,"&{{{%f %f} {%f %f}} %d %d",&ax,&ay,&sxmd,&symd,&mb,&mk)
	fmt.Printf("%d down: %.2f x %.2f \n",mb,sxmd,symd)
	mbd = (mb == 1)
	if mbd { h.MouseMoved(mm) }		// engage 1 tile click
}

var repl int		// replace will be by ctrl-h in select area or entire maze, by match
var cycl int		// cyclical set - C cycles, c sets - using c loc in keymap
var cycloc = 99

// edkey 'locks' on when pressed

func (h *holdableButton) MouseUp(mm *desktop.MouseEvent){

	mb := 0		// mb 1 = left, 2 = right, 4 = middle
	mbd = false
	h.MouseMoved(mm)				// disengage blotter
	ax := 0.0	// absolute x & y
	ay := 0.0
	exmd = 0.0	// rel x & y interm float32
	eymd = 0.0
	mk := 0		// mod key 1 = sh, 2 = ctrl, 4 = alt, 8 = logo
	pos := fmt.Sprintf("%v",mm)
	fmt.Sscanf(pos,"&{{{%f %f} {%f %f}} %d %d",&ax,&ay,&exmd,&eymd,&mb,&mk)
	ex := int(exmd / opts.dtec)
	ey := int(eymd / opts.dtec)

	if wpalop {
	if h.title == wpal.Title() {
		if mb == 4 && cmdoff {
		if G1 {
			g1edit_keymap[edkey] = plbuf[xy{ex, ey}]
			kys := g1mapid[g1edit_keymap[edkey]]
			keyst := fmt.Sprintf("G¹ assn key: %s = %03d, %s",map_keymap[edkey],g1edit_keymap[edkey],kys)
			statlin(cmdhin,keyst)
		} else {
			g2edit_keymap[edkey] = plbuf[xy{ex, ey}]
			kys := g2mapid[g2edit_keymap[edkey]]
			keyst := fmt.Sprintf("G² assn key: %s = %03d, %s",map_keymap[edkey],g2edit_keymap[edkey],kys)
			statlin(cmdhin,keyst)
		}}
		return
	}}
// right mb functions
	if mb == 2 {
		if pgdir != 0 {
			if sdb > 0 {
				sdbit(pgdir)
			} else {
				lrelod := pagit(pgdir)
				upd_edmaze(false)
				if lrelod { remaze(opts.mnum) }
			}
			return
		}
	}

 //   fmt.Printf("up %v\n",mm)
	if opts.edat > 0 {
fmt.Printf("%d up: %.2f x %.2f \n",mb,exmd,eymd)

		sx := int(sxmd / opts.dtec)
		sy := int(symd / opts.dtec)
		if ex < sx { t := ex; ex = sx; sx = t }		// swap if end smaller than start
		if ey < sy { t := ey; ey = sy; sy = t }
		var setcode int			// code to store given edit hotkey
		if G1 {
			setcode = g1edit_keymap[edkey]
		} else {
			setcode = g2edit_keymap[edkey]
		}
// a cut / copy / paste is active
		pasty := false
		if ccp == PASTE { pasty = true }
		if ccp != NOP {
		if mb != 1 { ccp_NOP(); fmt.Printf("ccp to NOP\n") }
//		if sx == ex && sy == ey { ccp_NOP() }
		if ccp != NOP {
			px :=0
			if ccp == COPY || ccp == CUT {
				py :=0
			for my := sy; my <= ey; my++ {
				px =0
			for mx := sx; mx <= ex; mx++ {
				cpbuf[xy{px, py}] = ebuf[xy{mx, my}]
fmt.Printf("%03d ",cpbuf[xy{px, py}])
				px++
				}
fmt.Printf("\n")
			py++
			}
			cpx = px - 1; if cpx < 0 { cpx = 0 }		// if these arent 1 less, the paste is 1 over
			cpy = py - 1; if cpy < 0 { cpy = 0 }
fmt.Printf("cc dun: px %d py %d\n",px,py)
// saving paste buffer now
			fil := fmt.Sprintf(".pb/pb_%07d_g%d.ed",pbcnt,opts.Gtp)
			pbcnt++
			sav_maz(fil, cpbuf, eflg, cpx, cpy, 0)
// local for short range
			fil = fmt.Sprintf(".pb/ses_%07d_g%d.ed",lpbcnt,opts.Gtp)
			lpbcnt++
			if G1 { lg1cnt = lpbcnt}
			if G2 { lg2cnt = lpbcnt}
			sav_maz(fil, cpbuf, eflg, cpx, cpy, 0)

// pb sort of doesnt end
				fil = fmt.Sprintf(".pb/cnt_g%d",opts.Gtp)
				file, err := os.Create(fil)
				if err == nil {
					wfs := fmt.Sprintf("%d",pbcnt)
					file.WriteString(wfs)
					file.Close()
				}
			}
			del = false						// copy or paste should not have del on
			if ccp == CUT { del = true }
			if pasty {
fmt.Printf("in pasty\n")
				ex = sx + cpx
				ey = sy + cpy
				if ex < 0 || ex > opts.DimX || cpx > opts.DimX { fmt.Printf("paste fail x\n"); return }
				if ey < 0 || ey > opts.DimY || cpy > opts.DimY { fmt.Printf("paste fail y\n"); return }
			} else {
				bwin(cpx+1, cpy+1, pbcnt - 1, cpbuf, eflg)		// draw the buffer
			}
		}}
// no access for keys: ?, \, C, A #a, eE, L, S, H, V
		fmt.Printf(" dtec: %f maze: %d x %d - element:%d\n",opts.dtec,ex,ey,ebuf[xy{ex, ey}])
		if mb == 4 && cmdoff {		// middle mb, do a reassign
			if G1 {
				g1edit_keymap[edkey] = ebuf[xy{ex, ey}]
				kys := g1mapid[g1edit_keymap[edkey]]
				keyst := fmt.Sprintf("G¹ assn key: %s = %03d, %s",map_keymap[edkey],g1edit_keymap[edkey],kys)
				statlin(cmdhin,keyst)
			} else {
				g2edit_keymap[edkey] = ebuf[xy{ex, ey}]
				kys := g2mapid[g2edit_keymap[edkey]]
				keyst := fmt.Sprintf("G² assn key: %s = %03d, %s",map_keymap[edkey],g2edit_keymap[edkey],kys)
				statlin(cmdhin,keyst)
			}
		} else {
		if del || cmdoff || pasty {
			rcl := 1		// loop count for undoing multi ops
		 for my := sy; my <= ey; my++ {
		   for mx := sx; mx <= ex; mx++ {
			rop := true		// run ops
			if ctrl {		// with ctrl held on drag op, only do outline
				rop = false
				if my == sy || my == ey { rop = true }
				if mx == sx || mx == ex { rop = true }
			}
// looped now, with ctrl op
			if rop {
				if del { undo_buf(mx, my,rcl); ebuf[xy{mx, my}] = 0; opts.bufdrt = true } else {	// delete anything for now makes a floor
				if pasty { undo_buf(mx, my,rcl); ebuf[xy{mx, my}] = cpbuf[xy{mx - sx, my - sy}]; opts.bufdrt = true }	// cant use setcode below, it wont set floors
				if setcode > 0 { undo_buf(mx, my,rcl); ebuf[xy{mx, my}] = setcode; opts.bufdrt = true }
fmt.Printf("%03d ",ebuf[xy{mx, my}])
				}
				rcl++
			}
			if edkey == 314 { repl = ebuf[xy{mx, my}] }		// just placeholder until new repl done -- yes, NOT being used
//			if edkey == 182 { ebuf[xy{mx, my}] = repl; opts.bufdrt = true }		//
		  }
fmt.Printf("\n")
		}
//			fmt.Printf(" chg elem: %d maze: %d x %d\n",ebuf[xy{mx, my}],mx,my)
		}}
		ed_maze(true)
	}

}

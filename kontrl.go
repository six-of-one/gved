package main

import (
	"fmt"
	"os"
//		"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// kontrol is for keyboard, mouse & input management


// input keys and keypress checks for canvas/ window
// since this is all that is called without other handler / timers
// - this is where maze update and edits will vector

// disable some cmd keys for edit mode click opts
var cmdoff bool

var anum int
var shift bool
var ctrl bool
var del bool
var logo bool		// other wise labeld "win" key
var ccp int			// cut copy and paste - current op
var edkey int		// for passing edit keys to clicker
var cmdhin string

// main keyboard handler

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
			relod = true
			needsav()
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
	if r >= '0' && r <= '9' || r == '`' {
		spau = fmt.Sprintf("numeric: %d", anum)
	}

	if deskCanvas, ok := w.Canvas().(desktop.Canvas); ok {
        deskCanvas.SetOnKeyDown(func(key *fyne.KeyEvent) {
//            fmt.Printf("Desktop key down: %h\n", key.Name)
			if key.Name == "BackSpace" {
				anum = (anum / 10);
				spau = fmt.Sprintf("numeric: %d", anum)
			}
			if key.Name == "Delete" { del = true }
			if key.Name == "BackSpace" { del = true }
			if key.Name == "LeftSuper" { logo = true }
			if key.Name == "LeftShift" { shift = true }
			if key.Name == "RightShift" { shift = true }
			if key.Name == "LeftControl" { ctrl = true }
			if key.Name == "RightControl" { ctrl = true }
        })
        deskCanvas.SetOnKeyUp(func(key *fyne.KeyEvent) {
//            fmt.Printf("Desktop key up: %v\n", key)
//			if key.Name == "Escape" { os.Exit(0) }
			if key.Name == "Delete" { del = false }
			if key.Name == "BackSpace" { del = false }
			if key.Name == "LeftSuper" { logo = false }
			if key.Name == "LeftShift" { shift = false }
			if key.Name == "RightShift" { shift = false }
			if key.Name == "LeftControl" { ctrl = false }
			if key.Name == "RightControl" { ctrl = false }
			if key.Name == "S" && ctrl { menu_sav() }
			if key.Name == "L" && ctrl  { menu_lod() }
			if key.Name == "R" && ctrl  { menu_res() }
			if key.Name == "U" && ctrl  { uswap() }
			if key.Name == "Z" && ctrl  { undo() }
			if key.Name == "Y" && ctrl  { redo() }
			if key.Name == "C" && ctrl  { if opts.edat > 0 { ccp = COPY }}
			if key.Name == "X" && ctrl  { if opts.edat > 0 { ccp = CUT }}
			if key.Name == "P" && ctrl  { if opts.edat > 0 { ccp = PASTE }}
			if key.Name == "Q" && ctrl  { needsav(); os.Exit(0) }
       })
    }

	fmt.Printf("r %v shift %v\n",r,shift)
		edkey = int(r)
		if cmdoff {
			if G1 {
//				if g1edit_keymap[edkey] < 0 { keyst := fmt.Sprintf("locked key: %s not usable",map_keymap[edkey]) }
				if g1edit_keymap[edkey] == 0 { keyst := fmt.Sprintf("G¹ free key: %s middle mouse click to set",map_keymap[edkey]); statlin(cmdhin,keyst) }
				if g1edit_keymap[edkey] > 0 {
					kys := g1mapid[g1edit_keymap[edkey]]
					keyst := fmt.Sprintf("G¹ ed key: %s = %03d, %s",map_keymap[edkey],g1edit_keymap[edkey],kys)
					statlin(cmdhin,keyst)
fmt.Printf("G¹ ed key: %d - %s\n",edkey,kys)
				}
			} else {
//				if g2edit_keymap[edkey] < 0 { keyst := fmt.Sprintf("locked key: %s not usable",map_keymap[edkey]) }
				if g2edit_keymap[edkey] == 0 { keyst := fmt.Sprintf("G² free key: %s middle mouse click to set",map_keymap[edkey]); statlin(cmdhin,keyst) }
				if g2edit_keymap[edkey] > 0 {
					kys := g2mapid[g2edit_keymap[edkey]]
					keyst := fmt.Sprintf("G² ed key: %s = %03d, %s",map_keymap[edkey],g2edit_keymap[edkey],kys)
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
				if cmdoff && opts.edat > 0 { ska = "edit keys mode" }
				opts.dntr = true
				relod = true
			}
			cmdhin = "cmds: ?, eE, fFgG, wWqQ, rRt, hm, pPT, sL, S, il, u, v, A #a"
			if cmdoff && opts.edat > 0 { cmdhin = "cmds: ? '\\' - edit keys, #c C, HV, A #a, eE, L, S" }
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
				if sdb == -1 {
					fil := fmt.Sprintf(".ed/ebuf.ed")				// save ebuf for relod
					sav_maz(fil, ebuf, eflg, opts.DimX, opts.DimY)
				} else { needsav() }
				fil := fmt.Sprintf(".ed/sd%05d_g%d.ed",anum,opts.Gtp)
				cnd := lod_maz(fil, ebuf, false)
				if cnd >= 0 { sdb = anum; for y := 0; y < 11; y++ { eflg[y] =  tflg[y] }; ed_maze(true) }
				anum = 0
			} else { opts.Nogtop = !opts.Nogtop; opts.dntr = (opts.edat > 0); relod = true }
		case 'e':
			if opts.Aob { dialog.ShowInformation("Edit mode", "Error: can not edit with border around maze!", w) } else {
				if opts.edat != 1 {
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
				}
				relod = true
			}
		case 69:		// E
			if opts.edat != 0 {
				smod = "View mode: "
				fmt.Printf("editor off, maze: %03d\n",opts.mnum+1)
				cmdoff = false
				opts.edat = 0
				opts.dntr = false
				ccp = NOP
				needsav()
				cmdhin = "cmds: ?, eE, fFgG, wWqQ, rRt, hm, pPT, sL, S, il, u, v, A #a"
				statlin(cmdhin,"")
				relod = true
			}
		case 'c':
			if anum > 0 && anum < 65 {
				cycl = anum - 1
				anum = 0
			} else { cycl--; cycl-- }
			fallthrough
		case 67:		// C
				cycl++
				if cycl > 64 { cycl = 0 }
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
// have anum !=0, save that buffer
				if anum > 0 && anum < sdmax && opts.edat > 0 {
					fil := fmt.Sprintf(".ed/sd%05d_g%d.ed",anum,opts.Gtp)
					sav_maz(fil, ebuf, eflg, opts.DimX, opts.DimY)
					anum = 0
				} else {
// with no anum, rotate curr ebuf thru s[1] - s[?], store eb in s[0]
					cnd := -1
					ldb := sdb
					if sdb == -1 {
						fil := fmt.Sprintf(".ed/ebuf.ed")				// save ebuf for relod
						sav_maz(fil, ebuf, eflg, opts.DimX, opts.DimY)
					}
					for cnd < 0 && ldb < sdmax {
						ldb++
						fil := fmt.Sprintf(".ed/sd%05d_g%d.ed",ldb,opts.Gtp)
						cnd = lod_maz(fil, ebuf, false)
						if cnd >= 0 { sdb = ldb; for y := 0; y < 11; y++ { eflg[y] =  tflg[y] }; ed_maze(true); spau = fmt.Sprintf("cmd: S - sdbuf: %d\n",sdb) }
					}
					if cnd < 0 {
						sdb = -1
						if opts.edat > 0 {
							fil := fmt.Sprintf(".ed/ebuf.ed")			// cycle back out
							cnd = lod_maz(fil, ebuf, true)
						if cnd >= 0 { for y := 0; y < 11; y++ { eflg[y] =  tflg[y] }; ed_maze(true); spau = fmt.Sprintf("cmd: S - maze %d\n",opts.mnum+1) }
						} else { remaze(opts.mnum) }
					}
				}
		default:
			relodsub = false
		}
// view cmd keys - also on edit, but blockable
	  if !cmdoff || opts.edat < 1 {
		relodsub = true
		switch r {
		case 'z':
			needsav()
			Ovwallpat = -1
// allow step parse through valid address
			if Aov > 0 {
				nav := addrver(Aov, -1)
				Aov = nav
			} else {
				opts.mnum -= 1
			}
			if opts.mnum < 0 { opts.mnum = maxmaze }
		case 'x':
			needsav()
			Ovwallpat = -1
			if Aov > 0 {
				nav := addrver(Aov, 1)
				Aov = nav
			} else {
				opts.mnum += 1
			}
			if opts.mnum > maxmaze { opts.mnum = 0 }
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
			opts.Gtp = 1
			opts.R14 = false
			G1 = true
			G2 = false
			maxmaze = 126
			spau = "G¹ "
			needsav()
		case 'l':
			opts.Gtp = 1
			opts.R14 = !opts.R14
			G1 = true
			G2 = false
			maxmaze = 126
			spau = "G¹ "
			needsav()
		case 'p':
			nothing = nothing ^ NOFLOOR
			spau = fmt.Sprintf("no floors: %d\n",nothing & NOFLOOR)
			opts.dntr = true
		case 80:		// P
			nothing = nothing ^ NOWALL
			spau = fmt.Sprintf("no walls: %d\n",nothing & NOWALL)
			opts.dntr = true
		case 84:		// T
			nt := (nothing & 511) + 1
			nothing = (nothing & 1536) + (nt & 511)
			if anum > 0 { nothing = (nothing & 1536) + anum; anum = 0 }		// set lower 9 bits of no-thing mask [ but not walls or floors ]
			spau = fmt.Sprintf("no things: %d\n",nothing & 511)				// display no things mask
			opts.dntr = true
		case 's':
			opts.SP = !opts.SP
			opts.dntr = true
		case 'u':
			opts.Gtp = 2
			G1 = false
			G2 = true
			maxmaze = 116
			spau = "G² mazes"
			needsav()
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
	if (relod || relodsub) {
		remaze(opts.mnum)
	}
	uptitl(opts.mnum, spau)
}

// data needing preserved by needsav - all this could be changed by the next op while dialog waits on user
var nsxd int
var nsyd int
var nsgg int
var nsmz int
var nssb int

func menu_ndsav(y bool) {
	if y {
		if sdb < 0 {
			fil := fmt.Sprintf(".ed/g%dmaze%03d.ed",nsgg,nsmz)
			sav_maz(fil, nsbuf, nsflg, nsxd, nsyd)
		} else {
			fil := fmt.Sprintf(".ed/sd%05d_g%d.ed",nssb,nsgg)
			sav_maz(fil, nsbuf, nsflg, nsxd, nsyd)
		}
	}
}

func needsav() {
	if opts.bufdrt {
		nsxd = opts.DimX
		nsyd = opts.DimY
		nsgg = opts.Gtp
		nsmz = opts.mnum+1
		nssb = sdb
// because the dialog doesnt hold back transition away from buffer, this has to immediatley save *everything*
		for y := 0; y <= nsxd; y++ {
		for x := 0; x <= nsyd; x++ {
			nsbuf[xy{x, y}] = ebuf[xy{x, y}]
		}}
		for y := 0; y < 11; y++ {
			nsflg[y] = eflg[y]
		}
		dia := fmt.Sprintf("Save changes for maze %d in .ed/g%dmaze%03d.ed ?\n\nWARNING:\nif not saved, changes will be discarded",nsmz,nsgg,nsmz)
		if nssb >= 0 { dia = fmt.Sprintf("Save changes in buffer %d to .ed/sd%05d_g%d.ed ?\n\nWARNING:\nif not saved, changes will be discarded",nssb,nssb,nsgg) }
		dialog.ShowConfirm("Save?",dia, menu_ndsav, w)
		opts.bufdrt = false;		// save clears this, clear here in case discard is selected
	}
}

func undo() {
	if delstak > 0 {
		delstak--
		revk := delbuf.revc[delstak]	// revoke count - items in loops can undo/redo all at once
fmt.Printf(" undo %d elem: %d maze: %d x %d - rloop: %d\n",delstak,delbuf.elem[delstak],delbuf.mx[delstak],delbuf.my[delstak],delbuf.revc[delstak])
		for revk > 0 && delstak >= 0 {
			sw := ebuf[xy{delbuf.mx[delstak], delbuf.my[delstak]}]
			ebuf[xy{delbuf.mx[delstak], delbuf.my[delstak]}] = delbuf.elem[delstak]
			delbuf.elem[delstak] = sw
			revk--
			if revk > 0 && delstak > 0 { delstak-- }
		}
fmt.Printf(" del %d elem: %d\n",delstak,delbuf.elem[delstak])
		opts.bufdrt = true
		ed_maze(false)
	}
}

func redo() {
	if delbuf.elem[delstak] != -1 {
fmt.Printf(" redo %d elem: %d maze: %d x %d - rloop: %d\n",delstak,delbuf.elem[delstak],delbuf.mx[delstak],delbuf.my[delstak],delbuf.revc[delstak])
		revk := delbuf.revc[delstak]	// revoke count - items in loops can undo/redo all at once
		for revk > 0 && delstak >= 0 {
			sw := ebuf[xy{delbuf.mx[delstak], delbuf.my[delstak]}]
			ebuf[xy{delbuf.mx[delstak], delbuf.my[delstak]}] = delbuf.elem[delstak]
//fmt.Printf(" redo %d elem: %d maze: %d x %d - rloop: %d\n",delstak,delbuf.elem[delstak],delbuf.mx[delstak],delbuf.my[delstak],delbuf.revc[delstak])
			delbuf.elem[delstak] = sw
			revk++
			delstak++
			if delbuf.elem[delstak] == -1 || delbuf.revc[delstak] == 1 { revk = 0}
		}
		opts.bufdrt = true
		ed_maze(false)
	}
//	ed_maze()
}

func uswap() {
	for y := 0; y <= opts.DimX; y++ {
		for x := 0; x <= opts.DimY; x++ {
			sw := ebuf[xy{x,y}]
			ebuf[xy{x,y}] = ubuf[xy{x,y}]
			ubuf[xy{x,y}] = sw
	}}
	for y := 0; y < 11; y++ { sw := eflg[y]; eflg[y] = uflg[y]; uflg[y] = sw }
	ed_maze(true)
}


// click area for edits

// button we can detect click and release areas for rubberband area & fills

type holdableButton struct {
    widget.Button
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

// &{{{387 545} {379 509.92188}} 4 0}

func (h *holdableButton) MouseDown(mm *desktop.MouseEvent){
	ax := 0.0	// absolute x & y
	ay := 0.0
	mb := 0		// mb 1 = left, 2 = right, 4 = middle
	mk := 0		// mod key 1 = sh, 2 = ctrl, 4 = alt, 8 = logo
	pos := fmt.Sprintf("%v",mm)
	fmt.Sscanf(pos,"&{{{%f %f} {%f %f}} %d %d",&ax,&ay,&sxmd,&symd,&mb,&mk)
	fmt.Printf("%d down: %.2f x %.2f \n",mb,sxmd,symd)
}


var repl int		// replace will be by ctrl-h in select area or entire maze, by match
var cycl int		// cyclical set - C cycles, c sets - using c loc in keymap
var cycloc = 99

// edkey 'locks' on when pressed

func (h *holdableButton) MouseUp(mm *desktop.MouseEvent){
 //   fmt.Printf("up %v\n",mm)
	if opts.edat > 0 {
		ax := 0.0	// absolute x & y
		ay := 0.0
		exmd = 0.0	// rel x & y interm float32
		eymd = 0.0
		mb := 0		// mb 1 = left, 2 = right, 4 = middle
		mk := 0		// mod key 1 = sh, 2 = ctrl, 4 = alt, 8 = logo
		pos := fmt.Sprintf("%v",mm)
		fmt.Sscanf(pos,"&{{{%f %f} {%f %f}} %d %d",&ax,&ay,&exmd,&eymd,&mb,&mk)
		fmt.Printf("%d up: %.2f x %.2f \n",mb,exmd,eymd)
		ex := int(exmd / opts.dtec)
		ey := int(eymd / opts.dtec)
		sx := int(sxmd / opts.dtec)
		sy := int(symd / opts.dtec)
		var setcode int			// code to store given edit hotkey
		if G1 {
			setcode = g1edit_keymap[edkey]
		} else {
			setcode = g2edit_keymap[edkey]
		}
// a cut / copy / paste is active
		if ccp != NOP {
		if mb != 1 { ccp = NOP }
		if sx == ex && sy == ey { ccp = NOP }
		if ccp != NOP {


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
		if del || cmdoff {
			if ex < sx { t := ex; ex = sx; sx = t }		// swap if end smaller than start
			if ey < sy { t := ey; ey = sy; sy = t }
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
				if setcode > 0 { undo_buf(mx, my,rcl); ebuf[xy{mx, my}] = setcode; opts.bufdrt = true }
				}
				rcl++
			}
			if edkey == 314 { repl = ebuf[xy{mx, my}] }		// just placeholder until new repl done -- yes, NOT being used
//			if edkey == 182 { ebuf[xy{mx, my}] = repl; opts.bufdrt = true }		//
		  }}
//			fmt.Printf(" chg elem: %d maze: %d x %d\n",ebuf[xy{mx, my}],mx,my)
		}}
		ed_maze(true)
	}

}

package main

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
)

// kontrol is for keyboard input management

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
//			if key.Name == "Tab" { tab = true }
        })
// handle keys up
        deskCanvas.SetOnKeyUp(func(key *fyne.KeyEvent) {
//	fmt.Printf("Desktop key up: %v\n", key)
			srelod := false
			sta := "vp ⊙ %d x %d"
			stu := ""
			px, py := 0, 0
			if key.Name == "Escape" {		// now toggle editor on/ off
				if opts.Aob { dialog.ShowInformation("Edit mode", "Error: can not edit with border around maze!", w) } else {
					if opts.edat == 0 {
						edit_on(0)
						if sdb > 0 { opts.dntr = true }
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
			if key.Name == "Delete" { if !ctrl { del = false }; if opts.edat == 1 { smod = "Edit mode: "; if ctrl { smod = "Edit DEL: " }; statlin(cmdhin,"")}}
//			if key.Name == "BackSpace" { del = false }
			if key.Name == "Home" { home = false }
			if key.Name == "LeftSuper" { logo = false }
			if key.Name == "LeftShift" { shift = false }
			if key.Name == "RightShift" { shift = false }
			if key.Name == "LeftControl" { ctrl = false }
			if key.Name == "RightControl" { ctrl = false }
//			if key.Name == "Tab" { tab = false }
			if key.Name == "S" && ctrl  { if shift { menu_savas() } else { menu_sav() }}
			if key.Name == "L" && ctrl  { if shift { menu_laodf() } else { menu_lod() }}
			if key.Name == "R" && ctrl  { menu_res() }
			if key.Name == "U" && ctrl  { uswap() }
			if key.Name == "T" && ctrl  { if ! cmdoff { nothing = 0; srelod = true }}
			if key.Name == "Z" && ctrl  { undo() }
			if key.Name == "Y" && ctrl  { redo() }
			if key.Name == "C" && ctrl  { menu_copy() }
			if key.Name == "X" && ctrl  { menu_cut() }
			if key.Name == "P" && ctrl  { if shift { pbsess_cyc(1) } else { menu_paste() }}
			if key.Name == "O" && ctrl  { if shift { pbmas_cyc(1) }}
			if key.Name == "Q" && ctrl  { exitsel = true; needsav() }
			if key.Name == "Left" {
				opts.dntr = true; srelod = true
				if ctrl { opts.DimX--; if opts.DimX < 1 { opts.DimX = 1 }
						  opts.bufdrt = true; sta = "maze: %d x %d"; px, py = opts.DimX, opts.DimY
						} else {
							vpx--; if vpx < 0 { vpx = 0 }; px, py = vpx, vpy
						}
					stu = fmt.Sprintf(sta, px, py)
					}
			if key.Name == "Right" {
				opts.dntr = true; srelod = true
				if ctrl { opts.DimX++; opts.bufdrt = true
						  sta = "maze: %d x %d"; px, py = opts.DimX, opts.DimY
						} else {
							//if vpx + viewp < opts.DimX { vpx++ }
							vpx++; px, py = vpx, vpy
						}
					stu = fmt.Sprintf(sta, px, py)
					}
			if key.Name == "Up" {
				opts.dntr = true; srelod = true
				if ctrl { opts.DimY--; if opts.DimY < 1 { opts.DimY = 1 };
						  opts.bufdrt = true; sta = "maze: %d x %d"; px, py = opts.DimX, opts.DimY
						} else {
							vpy--; if vpy < 0 { vpy = 0 }; px, py = vpx, vpy
						}
					stu = fmt.Sprintf(sta, px, py)
					}
			if key.Name == "Down" {
				opts.dntr = true; srelod = true
				if ctrl { opts.DimY++; opts.bufdrt = true
						  sta = "maze: %d x %d"; px, py = opts.DimX, opts.DimY
						} else {
							//if vpy + viewp < opts.DimY { vpy++ }
							vpy++; px, py = vpx, vpy
						}
					stu = fmt.Sprintf(sta, px, py)
					}
			if key.Name == "Prior" {
				if ctrl {
					viewp--; if viewp < minvp { viewp = minvp }
					opts.dntr = true; srelod = true
					stu = fmt.Sprintf("vp size %d x %d", viewp,viewp)
				} else {
				if sdb > 0 {
					sdbit(-1)
				} else { srelod = pagit(-1) }
			}}
			if key.Name == "Next" {
				if ctrl {
					viewp++; if viewp > maxvp { viewp = maxvp }
					opts.dntr = true; srelod = true
					stu = fmt.Sprintf("vp size %d x %d",viewp,viewp)
				} else {
				if sdb > 0 {
					sdbit(1)
				} else { srelod = pagit(1) }
			}}
			if stu != "" { statlin(cmdhin,stu) }
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
		edkey = valid_keys(int(r))
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
// a  only lower case not avail for edit hotkey
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
// with anum != 0, this becomes load s[1] buffer into ebuf
			if anum > 0 && anum < sdmax {
fmt.Printf("Load SD buf, anum: %05d, sdb: %d\n",anum, sdb)
				if opts.bufdrt { menu_savit(true) }		// autosave
				fil := fmt.Sprintf(".ed/sd%05d_g%d.ed",anum,opts.Gtp)
				cnd := lod_maz(fil, ebuf, false, true)
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
				edkey = cyckey					// pre set store cycl when cycling
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
// have anum !=0, save ebuf into that sd buffer
				if anum > 0 && anum < sdmax {		// save buf when not in edit, rand load can go into all mazes
fmt.Printf("Save to SD buf, anum: %05d, sdb: %d\n",anum, sdb)
					fil := fmt.Sprintf(".ed/sd%05d_g%d.ed",anum,opts.Gtp)
					sav_maz(fil, ebuf, eflg, opts.DimX, opts.DimY, 0 - anum, true)
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
		case 'e':
			Ovwallcol += 1
			if anum > 0 { Ovwallcol = anum - 1; anum = 0 }
			if Ovwallcol > 16 { Ovwallcol = 0 }
			eflg[6] = (Ovflorcol & 0x0f) << 4 + (Ovwallcol & 0x0f)
			spau = fmt.Sprintf("cmd: e - wallc: %d\n",Ovwallcol)
			opts.bufdrt = (opts.edat > 0)
			opts.dntr = (opts.edat > 0)
		case 'E':
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
				opts.DimX, opts.DimY = rotmirbuf(edmaze)
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
				opts.DimX, opts.DimY = rotmirbuf(edmaze)
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
			play_sfx("sfx/music.4sec.ogg")
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
			play_sfx("sfx/music.g2.4sec.ogg")
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
//fmt.Printf("kr cond relod: %t\n",relod || relodsub)
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
			sav_maz(fil, ebuf, eflg, nsxd, nsyd, nsmz, true)
		} else {
			fil := fmt.Sprintf(".ed/sd%05d_g%d.ed",nssb,nsgg)
			sav_maz(fil, ebuf, eflg, nsxd, nsyd, 0 - nssb, true)
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

func ccp_NOP() {
	if ccp == PASTE {
		if blot != ccblot { blot.Hide(); blot = ccblot }
	}
	ccp = NOP
	if opts.edat > 0 { smod = "Edit mode: "; statlin(cmdhin,"") }
}
func ccp_tog(op int) {
	wccp := ccp
	ccp_NOP()
	if wccp != op { ccp = op }
	if ccp == PASTE {
		if !wpbop { wpbimg = segimage(cpbuf,eflg,0,0,cpx+1,cpy+1,false) }
		blotup = true
		blotoff()
	}
}

func pbsess_cyc(dr int) {

if opts.Verbose { fmt.Printf("pbses c: %d %d\n",lpbcnt,sesbcnt) }
	sesbcnt += dr
	if sesbcnt >= lpbcnt { sesbcnt = 1 }
	if sesbcnt < 1 { sesbcnt = lpbcnt - 1 }
	if sesbcnt < 1 { sesbcnt = 1 }
	pb_upd("ses", "ses", sesbcnt)
}

func pbmas_cyc(dr int) {

if opts.Verbose { fmt.Printf("pbmas c: %d %d\n",pbcnt,masbcnt) }
	masbcnt += dr
	if masbcnt >= pbcnt { masbcnt = 1 }
	if masbcnt < 1 { masbcnt = pbcnt - 1 }
	if masbcnt < 1 { masbcnt = 1 }
	pb_upd("pb", "mas", masbcnt)
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
		cnd = lod_maz(fil, ebuf, false, true)
		if cnd >= 0 { sdb = ldb; for y := 0; y < 11; y++ { eflg[y] =  tflg[y] }; ed_maze(true); spar = fmt.Sprintf("cmd: S - ") }
		if dir < 0 && cnd < 0 && ldb == 1 { cnd = 0; break }
	}
	if cnd < 0 {
		menu_lodit(true)
	}
	return spar
}

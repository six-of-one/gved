G¹G²ved
=======

- visual editor for gauntlet, gauntlet 2 mazes
    (and eventually sanctuary...)
- gauntlet roms required: G¹ in ./ROMs-g1 and G² in ./ROMs (MUST be unzipped! see slapstic.go for details)
- fyne window module - had issues with gotilengine.TLN_DrawFrame(0) only wanting to draw 2 or 3 mazes

command line suggestions:
* gved -i maze115
    (interactive system)
* gved -v -g2 maze7
    (view output.png)
* gved floor0
    (view output.png)

interactive mode:
* '?' - calls up key hints dialog
* 'v' - lists maze rom addresses to terminal gved was run from
* 'z'. 'x' - previous & next maze - loops (address ov loops too)
* #'a' - type a valid maze number (digits 0-9) followed by 'a' - address also 229376 - 262145
* 'A' - switch between maze # and address override with Aov = {curr maze addr}
* visual ops should be straight forward
* 'r', 'R' - rotate maze +/- 90° are NOT a feature of gauntlet or g2
* 's' - displays special potions / gold bags at random empty locs
* 'L' - toggle generator monster indicator
* 'p', 'P' - floor and wall invisible in output.png
* 'T' - cycle bitmask to hide vars items (ref constants.go), can be set by #'T'
* 'w' 'e' 'f' 'g' - can all be proceded by #, just like #a, and shifted W,E,F,G reverse ops
* 'd' 'D' - turn edit mode on/ off
* writes output.png as each maze is viewed

Notes
=====
- some features are for research only...
- address override was used to verify gauntlet maze reads vs. link table at 0x38032
- rotate 90° is only used in sanctuary, NOT in gauntlet or g2
- the vars visual ops are view ONLY and do not yet allow editing

Research
========
* gaunt_prog.gnumeric (yes, it is gnumeric - libreoffice will open, but without the fancy parts)
* stats compiled from G¹/ G² roms

Issues
======
* editing is very rudimentary
* G¹ floor and wall colors are still being rendered with G² color palette

Attributions
=============
- based on code from [gex](https://github.com/alinsavix/gex/tree/master)

License
=======
[GPL3](https://www.gnu.org/licenses/gpl-3.0.html)

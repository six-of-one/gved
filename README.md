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
* 'A' - switch between maze # and address override with Aov = {curr maze addr}
* visual ops should be straight forward
* 's' - displays special potions at random empty locs
* writes output.png as each maze is viewed
* 'p', 'P' - floor and wall invisible in output.png

Research
========
* gaunt_prog.gnumeric (yes, it is gnumeric - libreoffice will open, but without the fancy parts)
* stats compiled from G¹/ G² roms

Issues
======
* still only a viewer
* G¹ floor and wall colors are still being rendered with G² color palette

Attributions
=============
- based on code from [gex](https://github.com/alinsavix/gex/tree/master)

License
=======
[GPL3](https://www.gnu.org/licenses/gpl-3.0.html)

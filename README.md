G¹G²ved
=======

- visual editor for gauntlet, gauntlet 2 mazes
    (and eventually sanctuary...)
- gauntlet roms required: G¹ in ./ROMs-g1 and G² in ./ROMs (MUST be unzipped! see slapstic.go for details)
- interactive mode still needs code moved out of terminal ops gotilengine display
- had issues with gotilengine.TLN_DrawFrame(0) only wanting to draw 2 or 3 mazes

command line suggestions:
* gved -i maze115
    (interactive system)
* gved -v -g2 maze7
    (view output.png)
* gved floor0
    (view output.png)

Research
========
* gaunt_prog.gnumeric (yes, it is gnumeric - libreoffice will open, but without the fancy parts)
* stats compiled from G¹/ G² roms

Issues
======
* still only a viewer
* possible to crash with valid parms
* G¹ floor and wall colors are still being rendered with G² color palette

Attributions
=============
- based on code from https://github.com/alinsavix/gex/tree/master

License
=======
[GPL3](https://www.gnu.org/licenses/gpl-3.0.html)

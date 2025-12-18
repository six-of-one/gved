G¹G²ved
=======

- visual editor for gauntlet, gauntlet 2 mazes

Notes
=====
update for mazedumps_g1 follows, settling:

* gved handles gauntlet maze read from link table at 0x38032
* gauntlet tile roms are also read, excepting floor and wall color ram palettes
* palette for base tiles is known and used for gauntlet renders now
* current floor and wall colors obtained from F4 option in emu
* trap walls and shootable walls may not color correct, they are gved rendered
* except all horiz scroll mazes - those will be done later
* pre-genned monsters are now the correct level and color palette
* pyramid monster gens now have a red letter: D, G, L, S to indicate monster type

answers to other old issues
---------------------------

1. colors were wrong because g2 palette slightly differs from gauntlet
2. monsters had to render from diff palette number for diff levels - the tiles are all the same
3. special potions are randomly inserted by code and do not appear in maze compressed store
4. there were some bank issues in the link table, fixed now, except 2 mazes that can swap with no issue

*old* Known issues:
  - some colors/textures may be off. Not 100% positive.
  - g1 lets mazes have level 1/2/3 monsters, but they all show as level 3
  - special potions (extra armor, etc) don't show
  - a couple of the mazes are obviously decoded incorrectly


Attributions
=============
- based on code from [gex](https://github.com/alinsavix/gex/tree/master)

License
=======
[GPL3](https://www.gnu.org/licenses/gpl-3.0.html)



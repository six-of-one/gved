package main

import (
	"fmt"
	"strings"
	"os"
	"io/ioutil"
)

func stor_maz(mazn int) {

	fmt.Printf("buffer maze entry\n")

	fil := fmt.Sprintf(".ed/g%dmaze%03d.ed",opts.Gtp,mazn)

	data, err := ioutil.ReadFile(fil)
	if err != nil {
		errs := fmt.Sprintf("%v",err)
		fmt.Print(errs)
		if strings.Contains(errs, "no such file") {
// settings
// 1. edit status 1, 0
// 2. wall_pattern floor_pattern wall_color floor_color
// 3. trick flags		-- see constants.go 
// 4. max_x max_y
// 5+ maze data
			file, err := os.Create(fil)
			if err == nil {
				maze := mazeDecompress(slapsticReadMaze(opts.mnum), false)
					lastx := 32
					if maze.flags&LFLAG4_WRAP_H > 0 {
						lastx = 31
					}
					lasty := 32
					if maze.flags&LFLAG4_WRAP_V > 0 {
						lasty = 31
					}
//				wfs := fmt.Sprintf("%d\n%d %d %d %d\n%0x\n%#b\n%d %d\n",1,Ovwallpat,Ovflorpat,Ovwallcol,Ovflorcol,maze.secret,maze.flags,lastx,lasty)
				wfs := fmt.Sprintf("%d\n",1)
// editor overs
				maze.optbyts[5] = (Ovflorpat & 0x0f) << 4 + (Ovwallpat & 0x0f)
				maze.optbyts[6] = (Ovflorcol & 0x0f) << 4 + (Ovwallcol & 0x0f)

				for y := 0; y < 11; y++ {
					wfs += fmt.Sprintf(" %02X", maze.optbyts[y])
				}
				wfs += "\n"
				for y := 0; y <= lasty; y++ {
					for x := 0; x <= lastx; x++ {

						wfs += fmt.Sprintf(" %02d", maze.data[xy{x, y}])
					}
					wfs += "\n"
				}
				file.WriteString(wfs)
				file.Close()
			}
		} else {
			fmt.Print(err)
			return
		}
	}
	fmt.Printf("buffer: %v\n",data)

}
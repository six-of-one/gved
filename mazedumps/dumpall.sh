#!/bin/bash

if [[ -f "./gved" ]]; then

for i in $(seq -f "%03g" 1 117);do ./gved -g2 -v -o mazedumps/maze$i.png maze$i > mazedumps/maze$i.txt;done

else

echo "run $0 from go build . dir"

fi

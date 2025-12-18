#!/bin/sh
for i in $(seq -f "%03g" 1 117)
do
    ../gved -v -o maze$i.png maze$i > maze$i.txt
done

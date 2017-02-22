#!/bin/bash

readers=(1 16 4 16 64)
askers=(1 2 8 32 64)

if [ -f results.txt ]; then
    cat /dev/null > results.txt
fi

echo "Build starting..."

go build

echo "Build complete!"

for ((i=0;i<5;++i)) do
    echo "Running test ${i}"
    time ./assignment1 -askers=${askers[i]} -readers=${readers[2]} -askdelay=10 -infiles=./data/pg1041.txt,./data/pg1103.txt,./data/pg1107.txt,./data/pg1112.txt,./data/pg1120.txt,./data/pg1128.txt,./data/pg1129.txt,./data/pg1514.txt,./data/pg1524.txt,./data/pg2235.txt,./data/pg2240.txt,./data/pg2242.txt,./data/pg2243.txt,./data/pg2264.txt,./data/pg2265.txt,./data/pg2267.txt >> results.txt
done

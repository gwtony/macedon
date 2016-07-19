#!/bin/bash
gocmd="$GOROOT/bin/go test"

if [ "$1" == "cover" ]; then 
	gocmd="$GOROOT/bin/go test -cover"
fi

for dir in `ls ./ | grep -v "test\|script\|sql\|Makefile\|error\|example\|variable\|utils\|run\|message\|README"`; do
	(cd $dir && $gocmd)
done

for dir in `ls ./modules | grep -v context`; do
	(cd modules/$dir && $gocmd)
done

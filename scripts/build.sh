#!/bin/bash
if [ $# -gt 0 -a "$1" == "debug" ]; then
	echo "Debug version"
	$GOROOT/bin/go build -o dist/bin/macedon -gcflags '-N -l' main/main.go
else 
	echo "Release version"
	$GOROOT/bin/go build -o dist/bin/macedon -ldflags '-s -w' main/main.go
fi
if [ $? -eq 0 ]; then
	cp conf/macedon.conf dist/conf/
	echo "Build done: binary in dist dir"
else
	echo "Build failed"
fi 

#!/bin/bash
$GOROOT/bin/go build -o dist/bin/macedon -gcflags '-N -l' example/macedon_main.go
if [ $? -eq 0 ]; then
	cp example/macedon.conf dist/conf/
	echo "Build done"
	tree dist
else
	echo "Build failed"
fi 

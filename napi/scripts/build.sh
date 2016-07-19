#!/bin/bash
$GOROOT/bin/go build -o dist/bin/napi -gcflags '-N -l' example/napi_main.go
if [ $? -eq 0 ]; then
	cp example/napi.conf dist/conf/
	echo "Build done"
	tree dist
else
	echo "Build failed"
fi 

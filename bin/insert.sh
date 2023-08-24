#!/bin/sh

sleep 5
echo "insert $1 start" 1>&2 
go run kdbfetch/main.go "$1" | go run kdbmining/main.go
echo "insert $1 end" 1>&2 
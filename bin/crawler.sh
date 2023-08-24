#!/bin/sh

go run kdb2csv/main.go | awk -F ',' 'NR>1{print $1}' | tr -d '"' | xargs -t -P3 -n 1 "./bin/insert.sh"
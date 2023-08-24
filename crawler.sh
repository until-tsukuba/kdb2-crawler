#!/bin/sh

SUBJECTS=$(go run kdb2csv/main.go | awk -F ',' '{print $1}' | tr -d '"')

for $SUBJECT in $SUBJECTS; do
    sleep 5
    go run kdbfetch/main.go "$SUBJECT" | go run kdbmining/main.go | go run esinsert/main.go
done
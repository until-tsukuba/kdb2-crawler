#!/bin/sh

SUBJECTS=$(go run kdb2csv/main.go | awk -F ',' '{print $1}' | tr -d '"')

for SUBJECT in $SUBJECTS; do
    sleep 5
    echo "insert $SUBJECT start" 1>&2 
    go run kdbfetch/main.go "$SUBJECT" | go run kdbmining/main.go | go run esinsert/main.go
    echo "insert $SUBJECT end" 1>&2 
done

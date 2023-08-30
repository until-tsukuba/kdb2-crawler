#!/bin/sh

random() {
    echo "$(($(openssl rand 4 | od -vAn -N4 -tu4) % 10 + 3))"
}

waitsec=$(random)
echo "insert $1 start after $waitsec sec." 1>&2
sleep "$waitsec"
go run kdbfetch/main.go "$1" | go run kdbmining/main.go
if [ $? -ne 0 ]; then
    echo $HTML
    exit 1
fi
echo $HTML
echo "$HTML" | go run esinsert/main.go
echo "insert $1 end" 1>&2 
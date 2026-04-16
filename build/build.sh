#!/bin/bash

Os=""
Od=""
Path=""

while [[ $# -gt 0 ]]; do
    case "$1" in
        -Os) Os="$2"; shift 2 ;;
        -Od) Od="$2"; shift 2 ;;
        -Path) Path="$2"; shift 2 ;;
        *) echo "Unknown argument: $1"; exit 1 ;;
    esac
done

GOOS="$Os" go build -o "$Od" "$Path"
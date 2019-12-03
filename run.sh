#!/bin/sh

SCRIPT_DIR="$( dirname "$0" )"
SOURCE_DIR="$SCRIPT_DIR/aoc${1?Which problem number?}"
cat "$SOURCE_DIR/input.txt" | go run "$SOURCE_DIR/main.go"

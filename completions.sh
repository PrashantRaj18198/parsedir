#!/bin/sh
set -e
rm -rf completions
mkdir completions
for sh in bash zsh fish; do
	go run cmd/parsedir/main.go completion "$sh" >"completions/art.$sh"
done
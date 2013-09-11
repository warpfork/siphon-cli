#!/bin/bash
git submodule update --init
GOPATH="$PWD/.gopath/" go build -o siphon

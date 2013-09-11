#!/bin/bash
git submodule update --init
GOPATH=$PWD go build -o siphon

#!/bin/bash
GORACE="log_path=siphon-race.log" GOPATH="$PWD/.gopath/" go test -race

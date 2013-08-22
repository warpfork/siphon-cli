#!/bin/bash
GORACE="log_path=siphon-race.log" GOPATH=$PWD go test -race

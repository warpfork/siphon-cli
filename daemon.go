package main

import (
	"fmt"
)

type daemonOpts struct {

}

func init() {
	parser.AddCommand("daemon", "daemon", "Run a daemon to spawn hosts", &daemonOpts{})
}

func (opts *daemonOpts) Execute(args []string) error {
	return fmt.Errorf("Daemon mode is not implemented yet.")
}

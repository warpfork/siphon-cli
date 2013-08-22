package main

import (
	"fmt"
	"os"
	"polydawn.net/siphon"
)

type attachOpts struct {
	Address string `short:"L" long:"addr" default:"defaults to unix://siphon.sock" description:"Address of host to dial, of the form unix://path/to/socket"`
}

func init() {
	parser.AddCommand("attach", "attach", "Attach to a host", &attachOpts{
		Address: "unix://siphon.sock",
	})
}

func (opts *attachOpts) Execute(args []string) error {
	addr, err := ParseNewAddr(opts.Address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "siphon: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Attaching to %s\n", addr.Label())

	client := siphon.Connect(addr)

	client.Connect()
	client.Attach(os.Stdin, os.Stdout)

	return nil
}

package main

import (
	"fmt"
	"os"
	"polydawn.net/siphon"
)

func attach(opts attachOpts_t) {
	addr, err := ParseNewAddr(opts.Address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "siphon: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Attaching to %s\n", addr.Label())

	client := siphon.NewClient(addr)

	client.Connect()
	client.Attach(os.Stdin, os.Stdout)
}

package main

import (
	"fmt"
	"polydawn.net/siphon"
)

func attach(options *attachOpts_t, socket string) {
	fmt.Printf("Attaching to %s\n", socket)

	addr := siphon.NewAddr("Siphon-Attach", "unix", socket)
	client := siphon.NewClient(addr)

	client.Connect()
	client.Attach()
}

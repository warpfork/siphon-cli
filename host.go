package main

import (
	"fmt"
	"os"
	"os/exec"
	"polydawn.net/siphon"
)

func host(opts hostOpts_t) {
	addr, err := ParseNewAddr(opts.Address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "siphon: %s\n", err)
		os.Exit(1)
	}
	cmd := exec.Command(opts.Command)

	fmt.Printf("Hosting %s at %s\n", opts.Command, addr.Label())

	host := siphon.NewHost(cmd, addr)

	host.Serve(); defer host.UnServe()
	host.Start()
	exitCode := host.Wait()
	fmt.Printf("siphon: %s exited %d\r\n", opts.Command, exitCode)
}

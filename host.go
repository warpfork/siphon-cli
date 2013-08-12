package main

import (
	"fmt"
	"os/exec"
	"polydawn.net/siphon"
)

func host(options *attachOpts_t, socket, command string) {
	fmt.Printf("Hosting %s at %s\n", command, socket)

	cmd := exec.Command(command)
	addr := siphon.NewAddr("Siphon-Host", "unix", socket)
	host := siphon.NewHost(cmd, addr)

	host.Serve(); defer host.UnServe()
	host.Start()
}

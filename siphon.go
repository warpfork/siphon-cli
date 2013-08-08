package main

import (
	"polydawn.net/siphon"
	"os/exec"
)

func main() {
	cmd := exec.Command("bash")
	addr := siphon.NewAddr("test", "unix", "demo.sock")
	host := siphon.NewHost(cmd, addr)
	host.Serve(); defer host.UnServe()

	client := siphon.NewClient(addr)
	client.Connect()

	host.Start()

	client.Attach()
}

package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"polydawn.net/siphon"
	"syscall"
)

type hostOpts struct {
	Address string `short:"L" long:"addr" optional:"true" default:"defaults to unix://siphon.sock" description:"Address to bind to and await client attachings, of the form unix://path/to/socket" `
	Command string `short:"c" long:"command" optional:"true" default:"defaults to /bin/sh" description:"Command to execute inside the new psuedoterminal" `
}

func init() {
	parser.AddCommand("host", "host", "Host a process", &hostOpts{
		Address: "unix://siphon.sock",
		Command: "/bin/sh",
	})
}

func (opts *hostOpts) Execute(args []string) error {
	addr, err := ParseNewAddr(opts.Address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "siphon: %s\n", err)
		os.Exit(EXIT_BADARGS)
	}

	cmd := exec.Command(opts.Command)
	shutdownCh := HandleShutdown() //handle control-c gracefully

	host := siphon.NewHost(cmd, addr)

	//Give the shutdown handler a callback to close the host
	shutdownCh <- func() {
		host.UnServe()
	}

	fmt.Printf("Hosting %s at %s\n", opts.Command, addr.Label)

	serveErr := host.Serve()
	defer host.UnServe()
	if serveErr != nil {
		if serveOpError, ok := serveErr.(*net.OpError); ok && serveOpError.Err == syscall.EADDRINUSE {
			fmt.Fprintf(os.Stderr, "%s\n", serveErr)
			os.Exit(EXIT_BIND_IN_USE)
		} else {
			panic(serveErr)
		}
	}

	host.Start()
	exitCode := host.Wait()
	fmt.Printf("siphon: %s exited %d\r\n", opts.Command, exitCode)

	return nil
}

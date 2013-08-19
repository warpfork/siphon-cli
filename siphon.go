
package main

import (
	"os"
	"fmt"
	"github.com/jessevdk/go-flags"
	"polydawn.net/siphon"
	"strings"
)

//Options for attach
type attachOpts_t struct {
	Address string `short:"L" long:"addr" default:"defaults to unix://siphon.sock" description:"Address of host to dial, of the form unix://path/to/socket"`
}

//Options for host
type hostOpts_t struct {
	Address string `short:"L" long:"addr" optional:"true" default:"defaults to unix://siphon.sock" description:"Address to bind to and await client attachings, of the form unix://path/to/socket" `
	Command string `short:"c" long:"command" optional:"true" default:"defaults to /bin/sh" description:"Command to execute inside the new psuedoterminal" `
}

//Options for daemon
type daemonOpts_t struct {

}


func main() {
	defer panicHandler()

	//Create command options
	attachOpts := attachOpts_t{
		Address: "unix://siphon.sock",
	}
	hostOpts := hostOpts_t{
		Address: "unix://siphon.sock",
		Command: "/bin/sh",
	}
	daemonOpts := daemonOpts_t{}

	//Construct parser with commands
	parser := flags.NewNamedParser("siphon", flags.HelpFlag)
	parser.AddCommand("host",   "host", "Host a process", &hostOpts)
	parser.AddCommand("attach", "attach", "Attach to a host", &attachOpts)
	parser.AddCommand("daemon", "daemon", "Run a daemon to spawn hosts", &daemonOpts)

	//Parse arguments, handle errors/help
	_, err := parser.Parse()
	if err != nil {
		if err.(*flags.Error).Type == flags.ErrHelp { //Usage help
			fmt.Printf("%s", err)
		} else {
			fmt.Printf("Error parsing: %s\n", err)
		}
		os.Exit(1)
	}

	//Switch on command... sloppy, should integrate with go-flags
	switch os.Args[1] {

	case "attach":
		attach(attachOpts)

	case "host":
		host(hostOpts)

	case "daemon":
		fmt.Errorf("Daemon mode is not implemented yet.")

	default:
		fmt.Errorf("Please specify attach, host, or daemon.")

	}
}

func panicHandler() {
	// print only the error message (don't dump stacks).
	// unless any debug mode is on; then don't recover, because we want to dump stacks.
	if len(os.Getenv("DEBUG")) == 0 {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(2)
		}
	}
}

func ParseNewAddr(addr string) (siphon.Addr, error) {
	addrParts := strings.SplitN(addr, "://", 2)
	switch addrParts[0] {
	case "unix":
		return siphon.NewAddr(addr, "unix", addrParts[1]), nil
	// case "tcp":
	//	// lib siphon supports this, but it's so very likely to be a bad idea that i'm making you compile a program for it yourself if you want this.
	//	return siphon.NewAddr(addr, "tcp", addrParts[1]), nil
	default:
		return siphon.Addr{}, fmt.Errorf("invalid protocol format.  \"%s\"", addr)
	}
}

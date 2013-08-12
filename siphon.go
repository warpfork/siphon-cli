
package main

import (
	"os"
	"fmt"
	"github.com/jessevdk/go-flags"
)

//Options for attach
type attachOpts_t struct {

}

//Options for host
type hostOpts_t struct {

}

//Options for daemon
type daemonOpts_t struct {

}


func main() {

	//Create command options
	attachOpts := new(attachOpts_t)
	hostOpts   := new(hostOpts_t)
	daemonOpts := new(daemonOpts_t)

	//Construct parser with commands
	parser := flags.NewNamedParser("siphon", flags.HelpFlag)
	parser.AddCommand("host",   "host", "Host a process", attachOpts)
	parser.AddCommand("attach", "attach", "Attach to a host", hostOpts)
	parser.AddCommand("daemon", "daemon", "Run a daemon to spawn hosts", daemonOpts)

	//Parse arguments, handle errors/help
	args, err := parser.Parse() //whelp
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
		socket := ""
		if len(os.Args) < 3 {
			socket = "./siphon.sock" //default
		} else {
			socket = os.Args[2]
		}

		attach(attachOpts, socket)

	case "host":
		socket, command := "", ""
		if len(args) < 1 {
			socket = "./siphon.sock" //default
			command = "/bin/bash"
		} else if len(args) < 2 {
			socket = args[0]
			command = "/bin/bash"
		} else {
			socket = args[0]
			command = args[1]
		}

		host(attachOpts, socket, command)

	case "daemon":
		fmt.Errorf("Daemon mode is not implemented yet.")

	default:
		fmt.Errorf("Please specify attach, host, or daemon.")

	}
}

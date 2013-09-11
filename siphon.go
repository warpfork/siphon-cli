
package main

import (
	"os"
	"os/signal"
	"fmt"
	"github.com/jessevdk/go-flags"
	"polydawn.net/siphon"
	"strings"
)

var parser = flags.NewNamedParser("siphon", flags.Default)

var EXIT_BADARGS = 1
var EXIT_PANIC = 2
var EXIT_BIND_IN_USE = 14

func main() {
	defer panicHandler()

	//Parse arguments, handle errors/help
	_, err := parser.Parse()
	if err != nil {
		if err.(*flags.Error).Type == flags.ErrHelp { //Usage help
			fmt.Printf("%s", err)
		} else {
			fmt.Printf("Error parsing: %s\n", err)
		}
		os.Exit(EXIT_BADARGS)
	}
}

func panicHandler() {
	// print only the error message (don't dump stacks).
	// unless any debug mode is on; then don't recover, because we want to dump stacks.
	if len(os.Getenv("DEBUG")) == 0 {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(EXIT_PANIC)
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

//Handle interrupt signals gracefully.
//This returns a channel through which you can pass a callback. This lets you define your own shutdown behavior.
func HandleShutdown() chan func() {
	shutdown := make(chan os.Signal)    // gets interrupt signal from os/signal
	listenCh := make(chan func()) // siphon hands us a listener to close when shutting down
	var callback func()

	//Tell go to inform us of interrupts
	signal.Notify(shutdown, os.Interrupt)

	//Store the listener when siphons hands it off, and handle shutdown signal
	go func(listenerCh <- chan func()) {
		for {
			select {
				case <- shutdown:
					if callback != nil {
						callback()
					}
					fmt.Printf("Caught Ctrl-C\n")
					os.Exit(1)
				case callback = <- listenerCh:
			}
		}
	}(listenCh)

	return listenCh
}


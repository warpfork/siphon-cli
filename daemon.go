package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"polydawn.net/siphon"
	"time"
	"strings"
	"syscall"
)

type daemonOpts struct {
	Address string `short:"L" long:"addr" optional:"true" default:"defaults to unix://siphon.daemon.sock" description:"Address to bind to and await client attachings, of the form unix://path/to/socket.  Each new client connection will result in the spawning of a new siphon host process, which the client will be transparently redirected to."`
	HostAddress string `short:"H" long:"host-addr" optional:"true" default:"defaults to unix://siphon.#####.sock" description:"Pattern of addresses to create new hosts at, of the form unix://path/to/socket.  '#' characters will be replaced by a random [1-9] digit (the number of hosts this daemon can spawn is implicitly limited to 9^{the number of '#' chars})." `
	HostCommand string `short:"c" long:"host-command" optional:"true" default:"defaults to /bin/sh" description:"Command to execute inside the new psuedoterminal" `
}

func init() {
	parser.AddCommand("daemon", "daemon", "Run a daemon to spawn hosts", &daemonOpts{
		Address: "unix://siphon.daemon.sock",
		HostAddress: "unix://siphon.#####.sock",
		HostCommand: "/bin/sh",
	})
}

// honestly, just for function grouping
type Daemon struct {
	opts *daemonOpts
}

func (opts *daemonOpts) Execute(args []string) error {
	addr, err := ParseNewAddr(opts.Address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "siphon: %s\n", err)
		os.Exit(EXIT_BADARGS)
	}

	listenCh := HandleShutdown()
	daemon := &Daemon{ opts: opts }

	listener, err := net.Listen(addr.Proto, addr.Addr)
	if err != nil {
		panic(err)
	}
	listenCh <- listener

	fmt.Printf("Serving daemon at %s\n", addr.Label)

	for listener != nil {
		conn, err := listener.Accept();
		if err != nil {
			if err.(*net.OpError).Err.Error() == "use of closed network connection" {
				break
			}
			panic(err)
		}
		fmt.Printf("siphon: daemon: accepted new client connection %p\r\n", conn)
		go daemon.handleRemoteClient(siphon.NewNetConn(conn))
	}

	return nil
}

func (daemon *Daemon) handleRemoteClient(conn *siphon.Conn) {
	defer conn.Close()

	// do startup handshake
	var hai siphon.Hello
	if err := conn.Decode(&hai); err != nil {
		fmt.Printf("siphon: daemon: %s, dropping client %s\n", err, conn.Label())
		return
	}
	if hai.Siphon != "siphon" {
		fmt.Printf("siphon: daemon: Encountered a non-siphon protocol, dropping client %s\n", conn.Label())
		return
	}
	if hai.Hello != "client" {
		fmt.Printf("siphon: daemon: Unexpected hello from not a client protocol, dropping client %s\n", conn.Label())
		return
	}
	if err := conn.Encode(siphon.HelloAck{
		Siphon: "siphon",
		Hello: "daemon",
	}); err != nil {
		fmt.Printf("siphon: daemon: %s, dropping client %s\n", err, conn.Label())
		return
	}

	addr := daemon.launchHost()
	fmt.Printf("siphon: daemon: launched host at %s for client %s\n", addr.Label, conn.Label())

	if err := conn.Encode(siphon.Redirect{
		Addr: addr,
	}); err != nil {
		fmt.Printf("siphon: daemon: %s, dropping client %s\n", err, conn.Label())
		return
	}
}

func (daemon *Daemon) launchHost() siphon.Addr {
	for {
		siphonHostAddrStr := strings.Map(
			func(r rune) rune {
				switch r {
				case '#':
					return []rune(fmt.Sprintf("%d", rand.Int31n(9)+1))[0]
				default:
					return r
				}
			},
			daemon.opts.HostAddress,
		)

		addr, addrErr := ParseNewAddr(siphonHostAddrStr)
		if addrErr != nil {
			panic(addrErr)
		}

		fmt.Printf("siphon: daemon: attempting to launch host with address %s\n", addr.Label)
		err := daemon.attemptLaunchHost(siphonHostAddrStr)
		if err == errDaemonRetryHost {
			continue
		} else if err != nil {
			panic(err)
		} else {
			// returning a siphon.Addr struct is a bit of a clusterfuck, since we also have to be honest
			//  with ourselves that we just exec'd and passed the string form across a shell.  But would
			//  spitting that string across the network be any better?  No.  At least this way, you're
			//  only assuming that the siphon processes on the daemon side are both parsing the string
			//  the same way, and that seems ever so slightly less reckless.
			return addr
		}
	}
}

var errDaemonRetryHost = fmt.Errorf("siphon: daemon: requested host bind address already in use")

/**
 * Try to make a host at a chosen address.
 *
 * This may be rejected because the bind attempt finds the address already consumed, in which case it's
 * appropriate to simply choose a new random address and try this function again.
 */
func (daemon *Daemon) attemptLaunchHost(addr string) error {
	args := []string{
		"host",
		"-L", addr,
		"-c", daemon.opts.HostCommand,
	}
	cmd := exec.Command("siphon", args...)

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	exitCode := -1
	exitCh := make(chan bool)
	go func() {
		defer close(exitCh)
		err := cmd.Wait()
		if exitError, ok := err.(*exec.ExitError); ok {
			if waitStatus, ok := exitError.Sys().(syscall.WaitStatus); ok {
				exitCode = waitStatus.ExitStatus()
			} else { panic(exitError); }
		}
	}()

	// TODO: this entire wait is nonsense, and what we should really do is have the host send a message back to this daemon that parented it over stdout.
	//  But if we're going to start using stdout for business, we're going to tear a lot of shoddy logging statements out, add a non-for-humans mode, and take output format a great deal more seriously.
	select {
	case <- exitCh :
		if exitCode == 0 {
			// we probably don't actually expect the hosted process to return so quickly, so maybe this should later be some kind fo error
			return nil
		} else if exitCode == EXIT_BIND_IN_USE {
			return errDaemonRetryHost
		} else {
			return fmt.Errorf("siphon: daemon: unrecognized problem in spawning host (exit code %d)", exitCode)
		}
	case <- time.After(100 * time.Millisecond) :
		return nil
	}
}

package main

import (
	"polydawn.net/siphon"
	"github.com/kr/pty"
	"bytes"
	"io"
	"os/exec"
	"testing"
	"github.com/coocood/assrt"
)

/**
Expect:
 - all of the input to come back out, because terminals default to echo mode.
 - then the grep'd string should come out, because the command recieved it and matched.
 - the grep'd string should come out surrounded by the escape codes for color, since grep's auto mode should detect that we're in a terminal.
*/
func TestPtySanity(t *testing.T) {
	assert := assrt.NewAssert(t)

	c := exec.Command("grep", "--color=auto", "bar")
	f, err := pty.Start(c)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		f.Write([]byte("foo\nbar\nbaz\n"))
		/*
		All of the input must be written in a single call to prevent this test from occationally failing nondeterministically, because:
		grep operates stream-wise and will start printing output before it has all of its input,
		and 3/4ths of the output lines are actually from the terminal operating in echo mode on the same channel.
		So there's actually a race between the terminal itself (somewhere down in kernel land I suppose?) and the output of grep.
		*/
		f.Write([]byte{4}) // EOT
	}()

	outBuffer := new(bytes.Buffer)
	io.Copy(outBuffer, f)
	out := string(outBuffer.Bytes())

	expected := // I have no idea where the CR characters come from.
		"foo\r\n"+
		"bar\r\n"+
		"baz\r\n"+
		"[01;31m[Kbar[m[K\r\n";

	assert.Equal(
		expected,
		out,
	)
}

func TestGrep(t *testing.T) {
	assert := assrt.NewAssert(t)

	cmd := exec.Command("grep", "--color=auto", "bar")
	host := siphon.NewHost(cmd, siphon.NewInternalAddr())
	host.Start()

	go func() {
		stdin := host.StdinPipe()
		stdin.Write([]byte("foo\nbar\nbaz\n"))
		stdin.Write([]byte{4}) // EOT
	}()

	outBuffer := new(bytes.Buffer)
	io.Copy(outBuffer, host.StdoutPipe())
	out := string(outBuffer.Bytes())

	expected := // I have no idea where the CR characters come from.
		"foo\r\n"+
		"bar\r\n"+
		"baz\r\n"+
		"[01;31m[Kbar[m[K\r\n";

	assert.Equal(
		expected,
		out,
	)
}

/**
Getting any answer from the tty command at all is pretty good news, since `go test` doesn't let the shell's tty come through.
*/
func TestActuallyNewTty(t *testing.T) {
	assert := assrt.NewAssert(t)

	cmd := exec.Command("tty")
	host := siphon.NewHost(cmd, siphon.NewInternalAddr())
	outBuffer := new(bytes.Buffer)
	hostOut := host.StdoutPipe()
	host.Start()
	io.Copy(outBuffer, hostOut)

	innerTty := string(outBuffer.Bytes())

	assert.NotEqual(
		"", // what exactly you get varies on the current state of your machine, but something like /dev/pts/12 is reasonable.
		innerTty,
	)
}

func TestUnixSocket(t *testing.T) {
	assert := assrt.NewAssert(t)

	addr := siphon.NewAddr("test", "unix", "test.sock")

	cmd := exec.Command("cat", "-")
	host := siphon.NewHost(cmd, addr)
	host.Serve(); defer host.UnServe()

	client := siphon.NewClient(addr)
	client.Connect()
	//client.Attach() // you can't do this in a test.  there's no tty.

	//FIXME: there's really no guarantee that the host finished processing the client connect request by now
	host.Start()

	go func() {
		stdin := host.StdinPipe()
		stdin.Write([]byte("foo\nbar\nbaz\n"))
		stdin.Write([]byte{4}) // EOT
	}()

	outBuffer := new(bytes.Buffer)
	io.Copy(outBuffer, client.Stdout())
	out := string(outBuffer.Bytes())

	expected :=
		"foo\r\n"+
		"bar\r\n"+
		"baz\r\n";
	expected = expected + expected
	assert.Equal(
		expected,
		out,
	)
}

func Template(t *testing.T) {
	assert := assrt.NewAssert(t)
	assert.Equal(
		"yes",
		"no",
	)
}

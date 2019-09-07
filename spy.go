package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/creack/pty"
	"golang.org/x/crypto/ssh/terminal"
)

type Spy struct{}

func (spy *Spy) Write(data []byte) (int, error) {
	os.Stderr.Write(data)
	return len(data), nil
}

func main() {
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(0, oldState)

	c := exec.Command("/bin/sh")
	f, err := pty.Start(c)
	if err != nil {
		panic(err)
	}

    // when the terminal is resized we receive a SIGWINCH
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			//fmt.Fprintf(os.Stderr, "resize\n")
			if err := pty.InheritSize(os.Stdin, f); err != nil {
				fmt.Fprintf(os.Stderr, "resize error: %s\n", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH

	spy := new(Spy)
	mw := io.MultiWriter(os.Stdout, spy)

	go io.Copy(mw, f)
	go io.Copy(f, os.Stdin)

	c.Wait()
}

package main

import (
	"io"
	"os"
	"os/exec"

    "github.com/creack/pty"
	"golang.org/x/crypto/ssh/terminal"
)

type spy struct {}

func (spy *spy) Write(data []byte) (int, error) {
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

	spy2 := new(spy)
	mw := io.MultiWriter(os.Stdout, spy2)

	go io.Copy(mw, f)
	go io.Copy(f, os.Stdin)

	c.Wait()
}

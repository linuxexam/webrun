package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/creack/pty"
)

func TestXxx(t *testing.T) {
	c := exec.Command("dir")
	f, err := pty.Start(c)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		f.Write([]byte("foo\n"))
		f.Write([]byte("bar\n"))
		f.Write([]byte("baz\n"))
		f.Write([]byte{4})
	}()

	io.Copy(os.Stdout, f)

}

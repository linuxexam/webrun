package main

import (
	"os"
	"runtime"

	"github.com/coder/websocket"
)

var ptmx *os.File

func RunCommand(conn *websocket.Conn, name string, args ...string) error {
	if runtime.GOOS == "windows" {
		return RunCommand_Pipe(conn, name, args...)
	} else {
		return RunCommand_Pty(conn, name, args...)
	}
}

package main

import (
	"runtime"

	"github.com/coder/websocket"
)

func RunCommand(conn *websocket.Conn, name string, args ...string) error {
	if runtime.GOOS == "windows" {
		return RunCommand_Pipe(conn, args[0], args[1:]...)
	} else {
		return RunCommand_Pty(conn, args[0], args[1:]...)
	}
}

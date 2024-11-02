package run

import (
	"os"
	"os/exec"
	"runtime"

	"github.com/coder/websocket"
)

func StartCommand(conn *websocket.Conn, cmd *exec.Cmd) (ptmx *os.File, err error) {
	if runtime.GOOS == "windows" {
		return StartCommand_Pipe(conn, cmd)
	} else {
		return StartCommand_Pty(conn, cmd)
	}
}

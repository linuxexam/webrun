package run

import (
	"bytes"
	"context"
	"log"
	"os"
	"os/exec"

	"github.com/coder/websocket"
)

// ui --- websocket --- pipes --- stdin/stdout/stderr --- command process
// the command process has pipes as its stdin, stdout and stderr, while the
// control terminal is still the caller's tty
// This works on Windows as well, as there is no PTY involved
// Without PTY, some simulation of line descipline function built into TTY driver
// has to be somewhat impleted explicitly here. e.g. CTR+C
func StartCommand_Pipe(conn *websocket.Conn, cmd *exec.Cmd) (*os.File, error) {
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	// cmd stdout ---> conn
	go func() {
		buf := make([]byte, 1)
		for {
			_, err := stdout.Read(buf)
			if err != nil {
				log.Print(err)
				break
			}
			// \n --> \r\n
			msg := bytes.ReplaceAll(buf, []byte("\n"), []byte("\r\n"))
			if err := conn.Write(context.Background(), websocket.MessageBinary, msg); err != nil {
				log.Print(err)
				break
			}
		}
	}()

	// cmd stderr ---> conn
	go func() {
		buf := make([]byte, 1)
		for {
			_, err := stderr.Read(buf)
			if err != nil {
				log.Print(err)
				break
			}
			// \n --> \r\n
			msg := bytes.ReplaceAll(buf, []byte("\n"), []byte("\r\n"))
			if err := conn.Write(context.Background(), websocket.MessageBinary, msg); err != nil {
				log.Print(err)
				break
			}
		}
	}()

	// conn ---> cmd stdin
	go func() {
		for {
			_, msg, err := conn.Read(context.Background())
			if err != nil {
				log.Println(err)
				break
			}
			log.Printf("received from WebSocket: %v\n", msg)

			// CTRL + C = "0x03"
			if bytes.Contains(msg, []byte("\x03")) {
				break
			}
			// ENTER(\r) --- > \n
			msg = bytes.ReplaceAll(msg, []byte("\r"), []byte("\n"))

			// echo mode
			conn.Write(context.Background(), websocket.MessageBinary, msg)

			if _, err := stdin.Write(msg); err != nil {
				log.Println(err)
				break
			}
			log.Printf("write to cmd'stdin: %v", msg)
		}
	}()

	return nil, nil
}

package run

import (
	"context"
	"log"
	"os"
	"os/exec"

	"github.com/coder/websocket"
	"github.com/creack/pty"
)

// ui --- websocket --- PTY master --- PTY slave(stdin,stdout,stderr and tty) --- command process
func StartCommand_Pty(conn *websocket.Conn, cmd *exec.Cmd) (*os.File, error) {
	var err error
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return ptmx, err
	}

	// Read from ptmx, send to conn
	go func() {
		buf := make([]byte, 1)
		for {
			_, err := ptmx.Read(buf)
			if err != nil {
				log.Printf("ptmx.Read: %v", err)
				return
			}
			log.Printf("read from ptmx: %v", buf)
			if err := conn.Write(context.Background(), websocket.MessageBinary, buf); err != nil {
				log.Printf("conn.Write: %v", err)
				return
			}
		}
	}()

	// Read from conn, send to ptmx
	go func() {
		for {
			_, buf, err := conn.Read(context.Background())
			if err != nil {
				log.Printf("conn.Write: %v", err)
				break
			}
			log.Printf("read from ws-term: %v", buf)

			if _, err := ptmx.Write(buf); err != nil {
				log.Printf("ptmx.Write: %v", err)
				break
			}
		}
	}()

	return ptmx, nil
}

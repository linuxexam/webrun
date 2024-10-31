package main

import (
	"context"
	"log"
	"os/exec"

	"github.com/coder/websocket"
	"github.com/creack/pty"
)

// ui --- websocket --- PTY master --- PTY slave(stdin,stdout,stderr and tty) --- command process
func RunCommand_Pty(conn *websocket.Conn, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	var err error
	ptmx, err = pty.Start(cmd)
	if err != nil {
		return err
	}
	defer ptmx.Close()

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
	for {
		_, buf, err := conn.Read(context.Background())
		if err != nil {
			log.Printf("conn.Write: %v", err)
			break
		}
		log.Printf("read from websocket term: %v", buf)

		if _, err := ptmx.Write(buf); err != nil {
			log.Printf("ptmx.Write: %v", err)
			break
		}
	}

	// wait cmd
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

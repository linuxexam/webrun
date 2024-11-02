package main

import (
	"context"
	"log"
	"net/http"
	"os/exec"

	"github.com/coder/websocket"
	"github.com/linuxexam/webrun/run"
)

func (c *Client) ServeTerm(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if c.ptmx != nil {
			c.ptmx.Close()
		}
		DelClient(c.id)
	}()

	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	c.termConn = conn
	log.Println("ws-term websocket connected.")
	defer func() {
		conn.CloseNow()
		log.Println("ws-term websocket disconnected.")
	}()

	// Run the command
	ctx, _ := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, c.commandName, c.commandArgs...)

	ptmx, err := run.StartCommand(conn, cmd)
	if err != nil {
		log.Printf("start command error: %v\n", err)
	}
	c.ptmx = ptmx
	log.Printf("%v\n", c)

	if err := cmd.Wait(); err != nil {
		log.Printf("cmd %s failed: %v\n", c.commandName, err)
	}
	log.Printf("%s:%s finished.\n", c.id, c.commandName)
}

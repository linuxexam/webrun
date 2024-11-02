package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/coder/websocket"
	"github.com/creack/pty"
)

func (c *Client) ServeAdmin(w http.ResponseWriter, r *http.Request) {
	defer DelClient(c.id)

	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Printf("ws-admin websocket.Accept: %v", err)
	}
	log.Println("ws-admin websocket connected.")
	c.adminConn = conn

	defer func() {
		conn.CloseNow()
		log.Println("ws-admin websocket disconnected.")
	}()

	for {
		_, buf, err := conn.Read(context.Background())
		if err != nil {
			log.Printf("c.Read: %v", err)
			break
		}
		msg := string(buf)
		log.Printf("ws-admin: %s\n", msg)
		var data map[string]interface{}
		err = json.Unmarshal(buf, &data)
		if err != nil {
			log.Println("Error unmarshalling JSON:", err)
			break
		}
		var cols uint16
		var rows uint16
		if data["type"] == "winsize" {
			cols = uint16(data["cols"].(float64))
			rows = uint16(data["rows"].(float64))
		}
		if c.ptmx != nil {
			err := pty.Setsize(c.ptmx, &pty.Winsize{
				Rows: rows,
				Cols: cols,
			})
			if err != nil {
				log.Printf("Set ptmx size failed")
			}
			log.Printf("Set ptmx size to %d x %d\n", cols, rows)
		}
	}
}

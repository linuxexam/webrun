package main

import (
	"context"
	"embed"
	"encoding/json"
	"flag"
	"io/fs"
	"log"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/coder/websocket"
	"github.com/creack/pty"
)

//go:embed ui
var UI embed.FS

const dev = false

func main() {
	// parse args
	var listen = flag.String("listen", ":8080", "listen port")
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		if runtime.GOOS == "windows" {
			args = append(args, "cmd")
		} else {
			args = append(args, "bash")
		}
	}

	path, err := exec.LookPath(args[0])
	if err != nil {
		log.Fatalf("command %s doesn't exist!", args[0])
	}
	args[0] = path

	// web UI
	if dev {
		http.Handle("/", http.FileServer(http.Dir("ui")))
	} else {
		sub, err := fs.Sub(UI, "ui")
		if err != nil {
			panic(err)
		}
		http.Handle("/", http.FileServer(http.FS(sub)))
	}

	// websocket handler for terminal proto
	http.HandleFunc("/ws-term", func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("ws-term websocket connected.")
		defer func() {
			c.CloseNow()
			log.Println("ws-term websocket disconnected.")
		}()

		// Run the command
		if err := RunCommand(c, args[0], args[1:]...); err != nil {
			log.Printf("command exit: %v", err)
		}
	})

	// websocket handler for admin
	// e.g. winsize change event
	http.HandleFunc("/ws-admin", func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("ws-admin websocket connected.")
		defer func() {
			c.CloseNow()
			log.Println("ws-admin websocket disconnected.")
		}()

		for {
			_, buf, err := c.Read(context.Background())
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
			if ptmx != nil {
				err := pty.Setsize(ptmx, &pty.Winsize{
					Rows: rows,
					Cols: cols,
				})
				if err != nil {
					log.Printf("Set ptmx size failed")
				}
			}
		}
	})

	log.Fatal(http.ListenAndServe(*listen, nil))
}

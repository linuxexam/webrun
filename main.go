package main

import (
	"embed"
	"flag"
	"io/fs"
	"log"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/coder/websocket"
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

	// websocket handler
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer c.CloseNow()

		// Run the command
		if err := RunCommand(c, args[0], args[1:]...); err != nil {
			log.Printf("command exit: %v", err)
		}
	})

	log.Fatal(http.ListenAndServe(*listen, nil))
}

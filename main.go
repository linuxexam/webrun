package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/coder/websocket"
)

//go:embed ui
var UI embed.FS

const dev = false

func main() {

	if dev {
		http.Handle("/", http.FileServer(http.Dir("ui")))
	} else {

		sub, err := fs.Sub(UI, "ui")
		if err != nil {
			panic(err)
		}
		http.Handle("/", http.FileServer(http.FS(sub)))
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer c.CloseNow()

		if runtime.GOOS == "windows" {
			err = RunWithPipe(c, os.Args[1], os.Args[2:]...)
		} else {
			err = RunWithPty(c, os.Args[1], os.Args[2:]...)
		}
		if err != nil {
			log.Printf("command exit: %v", err)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
)

//go:embed ui
var UI embed.FS

func main() {
	//http.Handle("/", http.FileServer(http.Dir("ui")))

	sub, err := fs.Sub(UI, "ui")
	if err != nil {
		panic(err)
	}
	http.Handle("/", http.FileServer(http.FS(sub)))

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer c.CloseNow()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		for {
			_, buf, err := c.Read(ctx)
			if err != nil {
				log.Println(err)
				break
			}
			log.Printf("received: %v", buf)
		}
		c.Close(websocket.StatusNormalClosure, "")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

package main

import (
	"embed"
	"flag"
	"io/fs"
	"log"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/linuxexam/webrun/util"
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
			args = append(args, "sh")
		}
	}

	path, err := exec.LookPath(args[0])
	if err != nil {
		log.Fatalf("command %s doesn't exist!", args[0])
	}
	args[0] = path

	// web UI
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// TODO: login
		clientID := util.GenerateUUID()
		cookie := &http.Cookie{
			Name:     "client_id",
			Value:    clientID,
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
			Secure:   false,
		}
		http.SetCookie(w, cookie)
		NewClient(clientID, args[0], args[1:]...)

		if dev {
			http.FileServer(http.Dir("ui")).ServeHTTP(w, r)
		} else {
			sub, err := fs.Sub(UI, "ui")
			if err != nil {
				panic(err)
			}
			http.FileServer(http.FS(sub)).ServeHTTP(w, r)
		}
	})

	// websocket handler for admin
	// e.g. winsize change event
	http.HandleFunc("/ws-admin", func(w http.ResponseWriter, r *http.Request) {
		idCookie, err := r.Cookie("client_id")
		if err != nil {
			http.NotFound(w, r)
			return
		}
		client := FindClient(idCookie.Value)
		if client == nil {
			http.NotFound(w, r)
			return
		}
		client.ServeAdmin(w, r)
	})

	// websocket handler for terminal proto
	http.HandleFunc("/ws-term", func(w http.ResponseWriter, r *http.Request) {
		idCookie, err := r.Cookie("client_id")
		if err != nil {
			http.NotFound(w, r)
			return
		}
		client := FindClient(idCookie.Value)
		if client == nil {
			http.NotFound(w, r)
			return
		}
		client.ServeTerm(w, r)
	})

	log.Fatal(http.ListenAndServe(*listen, nil))
}

package main

import (
	"os"
	"sync"

	"github.com/coder/websocket"
)

type Client struct {
	id          string
	commandName string
	commandArgs []string
	adminConn   *websocket.Conn
	termConn    *websocket.Conn
	ptmx        *os.File
}

var clients map[string]*Client = make(map[string]*Client)
var clients_lock sync.Mutex

func AddClient(c *Client) {
	clients_lock.Lock()
	defer clients_lock.Unlock()
	clients[c.id] = c
}

func FindClient(id string) *Client {
	clients_lock.Lock()
	defer clients_lock.Unlock()
	return clients[id]
}

func DelClient(id string) *Client {
	clients_lock.Lock()
	defer clients_lock.Unlock()
	c := clients[id]
	delete(clients, id)
	return c
}

func NewClient(id string, cmdName string, args ...string) *Client {
	c := &Client{
		id:          id,
		commandName: cmdName,
		commandArgs: args,
	}
	AddClient(c)
	return c
}

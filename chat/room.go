package main

import (
	"golang.org/x/net/websocket"
	"log"
)

const (
	messageBufferSize = 256
)

type room struct {
	forward chan []byte
	join    chan *client
	leave   chan *client
	clients map[*client]bool
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
			log.Println("New client joined")

		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
			log.Println("Client left")

		case msg := <-r.forward:
			log.Println("Message received: ", string(msg))

			for client := range r.clients {
				client.send <- msg
				log.Println("Message sent to client: ", string(msg))
			}

		}
	}
}

func (r *room) socket(socket *websocket.Conn) {
	client := newClient(socket, r)

	r.join <- client

	defer func() {
		r.leave <- client
	}()

	go client.write()

	client.read()

}

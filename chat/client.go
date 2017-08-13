package main

import (
	"golang.org/x/net/websocket"
	"log"
)

type client struct {
	socket *websocket.Conn
	send   chan []byte
	room   *room
}

type message struct {
	Message string `json:"message"`
}

func newClient(socket *websocket.Conn, room *room) *client {
	return &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   room,
	}
}

func (c *client) read() {
	defer c.socket.Close()

	for {
		var m message

		if err := websocket.JSON.Receive(c.socket, &m); err != nil {
			log.Println(err)
			break
		}

		c.room.forward <- []byte(m.Message)
	}
}

func (c *client) write() {
	defer c.socket.Close()

	for messageBytes := range c.send {
		msg := message{string(messageBytes)}
		if err := websocket.JSON.Send(c.socket, msg); err != nil {
			log.Println(err)
			break
		}
	}
}

package main

import (
	server "github.com/chatty/appServer"
	"github.com/chatty/handlers"
	"golang.org/x/net/websocket"
)

func main() {
	handlers.Register("/", &handlers.TemplateHandler{Filename: "chat.html",
		Parser: &handlers.AppTemplateParser{PathPrefix: "handlers/templates"}})

	room := newRoom()
	handlers.Register("/room", websocket.Handler(room.socket))

	go room.run()

	server.Start()
}

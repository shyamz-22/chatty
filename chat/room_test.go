package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http/httptest"
	"golang.org/x/net/websocket"
	"log"
	"strings"
	"time"
)

var testServerUrl string
var webSocketUrl string
var testServer *httptest.Server
var r *room

var _ = Describe("Room", func() {

	Context("client management", func() {
		BeforeEach(func() {
			testServerUrl = ""
			webSocketUrl = ""
			r = newRoom()
			testServer = givenAppIsRunning(r)
		})

		AfterEach(func() {
			testServer.Close()
		})

		It("adds client to the room when a client joins", func() {
			go r.run()

			whenClientConnectsToSocket()

			Expect(r.clients).To(HaveLen(1))
			Expect(r.leave).To(BeEmpty())
			for _, value := range r.clients {
				Expect(value).To(BeTrue())
			}
		})

		It("can add multiple clients to room", func() {
			go r.run()

			whenClientConnectsToSocket()
			whenClientConnectsToSocket()
			whenClientConnectsToSocket()

			Expect(r.clients).To(HaveLen(3))
			Expect(r.leave).To(BeEmpty())
			for _, value := range r.clients {
				Expect(value).To(BeTrue())
			}
		})

		It("deletes client from the room when a client leaves", func() {
			go r.run()

			whenClientTerminatesSocket(whenClientConnectsToSocket)

			time.Sleep(1 * time.Millisecond)

			Expect(r.clients).To(HaveLen(0))
		})


		It("can receive and send message to connected client", func() {
			go r.run()

			ws := whenClientConnectsToSocket()
			ws.Write([]byte(`{"message":"Hi"}`))

			time.Sleep(1 * time.Millisecond)
			var m message
			err := websocket.JSON.Receive(ws, &m)
			if err != nil {
				Panic()
			}

			Expect(m.Message).To(Equal("Hi"))
		})

		It("can broadcast message to all connected clients", func() {
			go r.run()

			ws := whenClientConnectsToSocket()
			ws_another := whenClientConnectsToSocket()

			ws.Write([]byte(`{"message":"Hi"}`))

			var m message
			var m_another message

			websocket.JSON.Receive(ws, &m)
			websocket.JSON.Receive(ws_another, &m_another)

			Expect(m.Message).To(Equal("Hi"))
			Expect(m_another.Message).To(Equal("Hi"))
		})
	})
})

func whenClientConnectsToSocket() (ws *websocket.Conn) {
	ws, err := websocket.Dial(strings.Join([]string{webSocketUrl,"/"}, ""), "", testServerUrl)
	if err != nil {
		log.Fatal("cannot connect to websocket----------------", err)
	}
	return ws
}

func whenClientTerminatesSocket(ws func()(ws *websocket.Conn)) {
	ws().Close()
}

func givenAppIsRunning(r *room) (server *httptest.Server) {
	testServer := httptest.NewServer(websocket.Handler(r.socket))
	testServerUrl = testServer.URL
	webSocketUrl = strings.Join([]string{"ws://", strings.Replace(testServerUrl,"http://","", 1)}, "")
	return testServer
}
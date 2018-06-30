package main

import (
	"testing"
	"net/http/httptest"
	"strings"

	"golang.org/x/net/websocket"
	"log"
	"time"
	"bytes"
	"github.com/chatty/trace"
)

var (
	testServerUrl = ""
	webSocketUrl  = ""
)

func Test_ClientJoinsRoom(t *testing.T) {
	room := newRoom()

	app := givenAppIsRunning(room)
	defer app.Close()

	go room.run()

	whenClientConnectsToSocket()

	actualNumberOfClients := len(room.clients)
	expectedNumberOfClients := 1

	if actualNumberOfClients != expectedNumberOfClients {
		t.Fatalf("Expected No of client is '%d' but was '%d'", expectedNumberOfClients, actualNumberOfClients)
	}

	if len(room.leave) > 0 {
		t.Fatalf("Expected to find no client in left room but was found")
	}

	for _, available := range room.clients {
		if !available {
			t.Fatalf("Expected to client to be available for chat")
		}
	}
}

func Test_MultipleClientsJoinsRoom(t *testing.T) {
	room := newRoom()

	app := givenAppIsRunning(room)
	defer app.Close()

	go room.run()

	whenClientConnectsToSocket()
	whenClientConnectsToSocket()
	whenClientConnectsToSocket()

	actualNumberOfClients := len(room.clients)
	expectedNumberOfClients := 3

	if actualNumberOfClients != expectedNumberOfClients {
		t.Fatalf("Expected No of client is '%d' but was '%d'", expectedNumberOfClients, actualNumberOfClients)
	}

	if len(room.leave) > 0 {
		t.Fatalf("Expected to find no client in left room but was found")
	}

	for _, available := range room.clients {
		if !available {
			t.Fatalf("Expected to client to be available for chat")
		}
	}
}

func Test_ClientLeavesRoom(t *testing.T) {
	room := newRoom()

	app := givenAppIsRunning(room)
	defer app.Close()

	go room.run()

	ws := whenClientConnectsToSocket()
	ws.Close()

	time.Sleep(1 * time.Millisecond)

	actualNumberOfClients := len(room.clients)
	expectedNumberOfClients := 0

	if actualNumberOfClients != expectedNumberOfClients {
		t.Fatalf("Expected No of client is '%d' but was '%d'", expectedNumberOfClients, actualNumberOfClients)
	}
}

func Test_ClientCanSendAndReceive(t *testing.T) {
	room := newRoom()
	var buf bytes.Buffer
	tracer := trace.New(&buf)
	room.tracer = tracer

	app := givenAppIsRunning(room)
	defer app.Close()

	go room.run()

	ws := whenClientConnectsToSocket()
	ws.Write([]byte(`{"message":"Hi"}`))

	time.Sleep(1 * time.Millisecond)

	expectedMessage := `{"message":"Hi"}`
	actualMessage := buf.String()

	if !strings.Contains(actualMessage, expectedMessage) {
		t.Fatalf("Expected to receive message containing '%s' but was '%s'", expectedMessage, actualMessage)
	}
}

func Test_BroadcastMessageWorks(t *testing.T) {
	room := newRoom()
	var buf bytes.Buffer
	tracer := trace.New(&buf)
	room.tracer = tracer

	app := givenAppIsRunning(room)
	defer app.Close()

	go room.run()

	ws := whenClientConnectsToSocket()
	wsAnother := whenClientConnectsToSocket()

	expected := "Hi"
	ws.Write([]byte(`{"message":"Hi"}`))

	time.Sleep(1 * time.Millisecond)

	var m = struct {
		Message string `json:"Message"`
	}{}

	websocket.JSON.Receive(wsAnother, &m)

	actual := m.Message
	if expected != actual {
		t.Fatalf("Expected to receive %s but received %s", expected, actual)
	}
}

func givenAppIsRunning(r *room) (server *httptest.Server) {
	testServer := httptest.NewServer(r)
	testServerUrl = testServer.URL
	webSocketUrl = strings.Join([]string{"ws://", strings.Replace(testServerUrl, "http://", "", 1)}, "")
	return testServer
}

func whenClientConnectsToSocket() (ws *websocket.Conn) {
	ws, err := websocket.Dial(strings.Join([]string{webSocketUrl, "/"}, ""), "", testServerUrl)
	if err != nil {
		log.Fatal("cannot connect to websocket----------------", err)
	}
	return ws
}

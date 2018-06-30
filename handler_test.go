package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"strings"
)

func Test_ChatHandler(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/", nil )
	response := httptest.NewRecorder()

	th := &templateHandler{filename: "chat.html"}
	th.ServeHTTP(response, request)

	exp := "Let us Chat"
	act := response.Body.String()

	if !strings.Contains(act, exp){
		t.Fatalf("Expected %s to contain %s", exp, act)
	}
}

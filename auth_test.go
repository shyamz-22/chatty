package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"log"
	"github.com/chatty/config"
	"strings"
)

func Test_ChatRequiresAuthenticationWorks(t *testing.T) {

	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	th := RequiresAuth(&templateHandler{filename: "chat.html"})
	th.ServeHTTP(response, request)

	actualStatus := response.Result().StatusCode
	expectedStatus := http.StatusTemporaryRedirect

	if expectedStatus != actualStatus {
		log.Fatalf("Expected %d != %d", expectedStatus, actualStatus)
	}

	actualRedirect := response.Header().Get("Location")
	expectedRedirect := "/login"

	if expectedRedirect != actualRedirect {
		log.Fatalf("Expected %s != %s", expectedRedirect, actualRedirect)
	}
}

func Test_LoginRedirectsToOIDCProvider(t *testing.T) {
	// GIVEN
	configuration := config.NewJsonConfigLoader("configuration.json").Load()
	setUpIdProvider(configuration.Auth, false)

	request, _ := http.NewRequest(http.MethodGet, "/login", nil)
	response := httptest.NewRecorder()

	// WHEN
	loginHandler(response, request)

	actualStatus := response.Result().StatusCode
	expectedStatus := http.StatusTemporaryRedirect

	if expectedStatus != actualStatus {
		log.Fatal(response.Body)
		log.Fatalf("Expected %d != %d", expectedStatus, actualStatus)
	}

	actualRedirect := response.Header().Get("Location")
	expectedRedirect := "https://identity-sandbox.vwgroup.io/oidc/v1/authorize?client_id=ac2911d1-f69a-436e-a654-e4f35840a073%40apps_vw-dilab_com&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Flogin%2Fcallback&response_type=code&scope=openid&state="

	if !strings.Contains(actualRedirect, expectedRedirect) {
		log.Fatalf("Expected %s != %s", expectedRedirect, actualRedirect)
	}
}

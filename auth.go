package main

import (
	"net/http"
	"github.com/markbates/goth/gothic"
	"log"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	oidc "github.com/markbates/goth/providers/openidConnect"
	"github.com/chatty/config"
)

type authHandler struct {
	next http.Handler
}

func setUpIdProvider(auth config.AuthConfig, isSecure bool) {

	idP, err := oidc.New(auth.ClientId, auth.ClientSecret, auth.CallbackUrl, auth.DiscoveryUrl)

	if idP == nil || err != nil {
		panic("Please provide a  valid setup for OpenId Provider")
		return
	}

	store := sessions.NewCookieStore([]byte(auth.SecretKey))
	store.MaxAge(auth.CookieAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isSecure

	gothic.Store = store
	gothic.GetProviderName = func(req *http.Request) (string, error) {
		return "openid-connect", nil
	}

	goth.UseProviders(idP)
}


func RequiresAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("_gothic_session")

	if err == http.ErrNoCookie {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.next.ServeHTTP(w, r)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if gothUser, err := gothic.CompleteUserAuth(w, r); err == nil {
		log.Println(gothUser)
	} else {
		gothic.BeginAuthHandler(w, r)
	}
}

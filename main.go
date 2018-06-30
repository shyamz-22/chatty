package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/chatty/trace"
	"github.com/chatty/config"
)

// templ represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP handles the HTTP request.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	var isCloud = flag.Bool("isCloud", false, "Is cloud deployment")

	flag.Parse()

	var configuration *config.Configuration
	var isSecureCookie bool

	if !*isCloud {
		configuration = config.NewJsonConfigLoader("configuration.json").Load()
		isSecureCookie = false
	}

	//setup identity provider
	setUpIdProvider(configuration.Auth, isSecureCookie)


	r := newRoom()
	r.tracer = trace.New(os.Stdout)

	mux := http.NewServeMux()
	mux.HandleFunc("/login", loginHandler)
	mux.Handle("/", RequiresAuth(&templateHandler{filename: "chat.html"}))
	mux.Handle("/room", r)

	// get the room going
	go r.run()

	// start the web server
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}

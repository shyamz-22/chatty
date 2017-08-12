package handlers

import (
	"net/http"
	"html/template"
	"path/filepath"
	"sync"
	"fmt"
)

func Register(path string, handle http.Handler) {
	http.Handle(path, handle)
}

type TemplateParser interface {
	parse(fileName string) *template.Template
}

type TemplateHandler struct {
	once     sync.Once
	Filename string
	aTemplate *template.Template
	Parser TemplateParser
}

type AppTemplateParser struct {
	PathPrefix string
}

func (parser *AppTemplateParser) parse(fileName string) *template.Template {
	return ParseTemplate(parser.PathPrefix, fileName)
}

func (t *TemplateHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	t.once.Do(func() {
		t.aTemplate = t.Parser.parse(t.Filename)
	})

	t.aTemplate.Execute(writer, request)
}

func ParseTemplate(pathPrefix string, fileName string) *template.Template{
	fmt.Println(filepath.Join("templates"))
	return template.Must(template.ParseFiles(filepath.Join(pathPrefix, fileName)))
}

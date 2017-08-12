package handlers

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"html/template"

	"net/http"
	"net/http/httptest"
	"log"
	"io"
)

var timesInvoked = 0

var _ = Describe("chat application", func() {

	BeforeEach(func() {
		timesInvoked = 0
	})

	It("loads the chat page", func() {
		request := makeRequest(http.MethodGet, "/", nil)
		responseRecorder:= performRequest(&TemplateHandler{Filename: "chat.html", Parser: &AppTemplateParser{PathPrefix: "templates"}}, request)

		Expect(responseRecorder.Code).To(Equal(http.StatusOK))
		Expect(responseRecorder.Body.String()).To(ContainSubstring("Let's chat!!!"))
	})

	It("compiles template only once", func() {
		parser := &MockTemplateParser{PathPrefix: "templates"}
		handler := &TemplateHandler{Filename: "chat.html", Parser: parser}

		firstRequest := makeRequest(http.MethodGet, "/", nil)
		performRequest(handler, firstRequest)

		Expect(timesInvoked).To(Equal(1))

		secondRequest := makeRequest(http.MethodGet, "/", nil)
		performRequest(handler, secondRequest)

		Expect(timesInvoked).To(Equal(1))
	})
})

func performRequest(handler *TemplateHandler, r *http.Request) *httptest.ResponseRecorder {
	responseRecorder := httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, r)
	return responseRecorder
}

func makeRequest(method string, path string, body io.Reader) (*http.Request) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		log.Fatal(err)
	}
	return req
}

type MockTemplateParser struct {
	PathPrefix string
}

func (t *MockTemplateParser) parse(fileName string) *template.Template {
	timesInvoked++
	return ParseTemplate(t.PathPrefix, fileName)
}

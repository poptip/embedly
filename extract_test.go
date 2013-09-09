package embedly

import (
	"github.com/coocood/assrt"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
)

var (
	originalHost string
	statusCode   int
	results      []byte
	server       *http.Server
	assert       *assrt.Assert
	once         sync.Once
)

func setUp() {
	originalHost = Host
	Host = "http://localhost:12345"
	// Custom handler used for mocking request results.
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(statusCode)
		w.Write(results)
	})
	server = &http.Server{
		Addr:    ":12345",
		Handler: handler,
	}
	go func() {
		assert.MustNil(server.ListenAndServe())
	}()
}

func tearDown() {
	Host = originalHost
}

func mockRequest(status int, filename string) {
	statusCode = status
	if len(filename) > 0 {
		r, err := ioutil.ReadFile(filename)
		assert.MustNil(err)
		results = r
	}
}

func TestExtract(t *testing.T) {
	once.Do(setUp)
	defer once.Do(tearDown)
	assert = assrt.NewAssert(t)
	c := NewClient("")

	mockRequest(200, "responses.json")
	response, err := c.ExtractOne("http://www.theonion.com/articles/fasttalking-computer-hacker-just-has-to-break-thro,32000/", Options{})
	assert.MustNil(err)
	assert.Equal("Fast-Talking Computer Hacker Just Has To Break Through Encryption Shield Before Uploading Nano-Virus", response.Title)

	urls := []string{
		"http://google.com",
		"http://yahoo.com",
		"http://bing.com",
		"http://cnn.com",
		"http://bbc.com",
	}
	links, err := c.Extract(urls, Options{})
	assert.MustNil(err)
	assert.Equal(5, len(links))

	urls = []string{
		"http://google.com",
		"http://yahoo.com",
		"http://bing.com",
		"http://cnn.com",
		"http://bbc.com",
		"http://google.com",
		"http://yahoo.com",
		"http://bing.com",
		"http://cnn.com",
		"http://bbc.com",
		"http://google.com",
		"http://yahoo.com",
		"http://bing.com",
		"http://cnn.com",
		"http://bbc.com",
	}
	links, err = c.Extract(urls, Options{})
	assert.MustNil(err)
	assert.Equal(15, len(links))

	mockRequest(500, "error_response.json")
	_, err = c.ExtractOne("nope", Options{})
	assert.MustNotNil(err)
}

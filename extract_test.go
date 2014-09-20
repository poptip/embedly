package embedly

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/coocood/assrt"
)

var (
	originalHost string
	statusCode   int
	results      []byte
	server       *http.Server
	assert       *assrt.Assert
	setUpOnce    sync.Once
	tearDownOnce sync.Once
)

func setUp() {
	log.Println("Setup Fake API")
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
	log.Println("Teardown Fake API")
	Host = originalHost
}

func mockRequest(status int, filename string) {
	statusCode = status
	if len(filename) > 0 {
		r, err := ioutil.ReadFile("mocks/" + filename + ".json")
		assert.MustNil(err)
		results = r
	}
}

func TestExtract(t *testing.T) {
	setUpOnce.Do(setUp)
	defer tearDownOnce.Do(tearDown)
	assert = assrt.NewAssert(t)
	c := NewClient("")

	mockRequest(200, "response")
	response, err := c.ExtractOne("http://www.theonion.com/articles/fasttalking-computer-hacker-just-has-to-break-thro,32000/", Options{})
	assert.MustNil(err)
	assert.Equal("Fast-Talking Computer Hacker Just Has To Break Through Encryption Shield Before Uploading Nano-Virus", response.Title)
	assert.Equal(TypeHTML, response.Type)

	mockRequest(200, "giphy")
	response, err = c.ExtractOne("http://giphy.com/gifs/XYyT3ZRNzaflK", Options{})
	assert.MustNil(err)
	assert.Equal("Jim Carrey Animated GIF", response.Title)
	assert.Equal(TypeHTML, response.Type)

	mockRequest(200, "responses5")
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

	mockRequest(200, "responses10")
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
	for i, link := range links {
		assert.True(len(link.Title) > 0, strconv.Itoa(i)+"th link with empty title")
	}

	mockRequest(500, "error_response")
	_, err = c.ExtractOne("nope", Options{})
	assert.MustNotNil(err)
}

func TestRealAPI(t *testing.T) {
	apiKey := os.Getenv("EMBEDLY_API_KEY")
	if len(apiKey) == 0 {
		return
	}
	log.Println("Testing with the Embedly API")
	client := NewClient(apiKey)
	response, err := client.Extract([]string{"http://giphy.com/gifs/XYyT3ZRNzaflK"}, Options{})
	assert.MustNil(err)
	log.Printf("%+v", response[0])
	assert.Equal(t, "Jim Carrey Animated GIF", response[0].Title)
}

package embedly

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var (
	Host = "http://api.embed.ly"
)

type Client struct {
	key string
}

func NewClient(key string) *Client {
	return &Client{key}
}

// The main exported function that will extract urls.
func (c *Client) Extract(urls []string, options Options) ([]Response, error) {
	responses := make([]Response, len(urls))
	for i := 0; i < len(urls); i += 10 {
		to := len(urls)
		if to > i+10 {
			to = i + 10
		}
		r, err := c.extract(urls[i:to], options)
		if err != nil {
			return nil, err
		}
		for j, res := range r {
			responses[i+j] = res
		}
	}
	return responses, nil
}

// Shortcut for extracting one url.
func (c *Client) ExtractOne(url string, options Options) (*Response, error) {
	responses, err := c.extract([]string{url}, options)
	if err != nil {
		return nil, err
	}
	return &responses[0], nil
}

// extract will call extract 10 urls at max.
func (c *Client) extract(urls []string, options Options) ([]Response, error) {
	v := url.Values{}
	v.Add("key", c.key)
	if len(urls) == 0 {
		return nil, errors.New("At least one url is required")
	} else if len(urls) == 1 {
		v.Add("url", urls[0])
	} else {
		v.Add("urls", strings.Join(urls, ","))
	}

	addInt(&v, "maxwidth", options.MaxWidth)
	addInt(&v, "maxheight", options.MaxHeight)
	addInt(&v, "words", options.Words)
	addInt(&v, "chars", options.Chars)
	addBool(&v, "wmode", options.WMode)
	addBool(&v, "allowscripts", options.AllowScripts)
	addBool(&v, "nostyle", options.NoStyle)
	addBool(&v, "autoplay", options.Autoplay)
	addBool(&v, "videosrc", options.VideoSrc)
	addBool(&v, "frame", options.Frame)
	addBool(&v, "secure", options.Secure)

	// Make the request.
	resp, err := http.Get(Host + "/1/extract?" + v.Encode())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Got non 200 status code: %s %q", resp.Status, body)
	}

	// Read the JSON message from the body.
	defer resp.Body.Close()
	response := []Response{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

// addInt adds an int if non-zero.
func addInt(v *url.Values, name string, value int) {
	if value > 0 {
		v.Add(name, string(value))
	}
}

// addBool adds a boolean value if set to true.
func addBool(v *url.Values, name string, value bool) {
	if value {
		v.Add(name, "true")
	}
}

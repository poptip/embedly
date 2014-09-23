package embedly

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
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
		res, err := c.extract(urls[i:to], options)
		if err != nil {
			return nil, err
		}

		reslen := to - i
		if reslen > len(res) {
			reslen = len(res)
		}
		for j := 0; j < reslen; j++ {
			responses[i+j] = res[j]
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
	addr := Host + "/1/extract?"
	for i, u := range urls {
		urls[i] = url.QueryEscape(u)
	}
	if len(urls) == 0 {
		return nil, errors.New("At least one URL is required")
	} else {
		for _, url := range urls {
			if len(url) == 0 {
				return nil, errors.New("A URL cannot be empty")
			}
		}
		addr += "urls=" + strings.Join(urls, ",")
	}

	v := url.Values{}
	v.Add("key", c.key)
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
	addr += "&" + v.Encode()
	resp, err := http.Get(addr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 500 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Got non 200 status code: %s %q", resp.Status, body)
	}

	// Read the JSON message from the body.
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
		v.Add(name, strconv.Itoa(value))
	}
}

// addBool adds a boolean value if set to true.
func addBool(v *url.Values, name string, value bool) {
	if value {
		v.Add(name, "true")
	}
}

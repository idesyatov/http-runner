package httpclient

import (
	"net/http"
	"time"
)

type Client struct {
	http.Client
}

func NewClient(timeout time.Duration) *Client {
	return &Client{
		Client: http.Client{
			Timeout: timeout,
		},
	}
}

// SendRequest sends an HTTP request with the specified method, URL, and headers.
func (c *Client) SendRequest(method, url string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	// Set headers if provided
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return c.Do(req)
}

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

func (c *Client) SendRequest(method, url string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}
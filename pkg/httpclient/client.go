package httpclient

import (
	"bytes"
	"encoding/json"
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

// SendRequest sends an HTTP request with the specified method, URL, headers, and data.
func (c *Client) SendRequest(method, url string, headers map[string]string, data map[string]string) (*http.Response, error) {
	var body *bytes.Buffer

	// Serialize data to JSON if it's not nil
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonData)
	} else {
		body = bytes.NewBuffer([]byte{}) // Empty body for GET requests
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// Set headers if provided
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Set Content-Type header for JSON data if not already set
	if data != nil && headers["Content-Type"] == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.Do(req)
}

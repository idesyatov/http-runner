package httpclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"time"
)

type Client struct {
	http.Client
}

// NewClient builds an HTTP client with the given per-request timeout. When
// insecure is true, TLS certificate verification is skipped. When
// followRedirects is false, redirects are not followed (the last response is
// returned as-is).
func NewClient(timeout time.Duration, insecure, followRedirects bool) *Client {
	c := http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		},
	}
	if !followRedirects {
		c.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
	return &Client{Client: c}
}

// SendRequest sends an HTTP request with the specified method, URL, headers, and data.
func (c *Client) SendRequest(method, url string, headers map[string]string, data interface{}) (*http.Response, error) {
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

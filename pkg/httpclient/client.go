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
//
// maxIdleConns sizes the idle connection pool (both total and per-host) so that
// keep-alive actually reuses connections under load: net/http defaults to only
// 2 idle connections per host, which would churn TCP/TLS handshakes at higher
// concurrency and skew latency/throughput. Values <= 0 fall back to that
// default.
func NewClient(timeout time.Duration, insecure, followRedirects bool, maxIdleConns int) *Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
	}
	if maxIdleConns > 0 {
		transport.MaxIdleConns = maxIdleConns
		transport.MaxIdleConnsPerHost = maxIdleConns
	}
	c := http.Client{
		Timeout:   timeout,
		Transport: transport,
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

	// Set Content-Type header for JSON data if not already set. Use the request
	// header (canonicalised) rather than the raw map so a user-supplied header
	// in any case (e.g. "content-type") is honoured and not overwritten.
	if data != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.Do(req)
}

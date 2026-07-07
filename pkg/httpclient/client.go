package httpclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/http/httptrace"
	"time"
)

type Client struct {
	http.Client
}

// Trace holds per-request connection phase timings captured via httptrace.
// Phases that did not occur (for example DNS, connect and TLS on a reused
// keep-alive connection) stay at zero; Reused reports whether the underlying
// connection came from the idle pool.
type Trace struct {
	DNS     time.Duration // DNS resolution
	Connect time.Duration // TCP connection establishment
	TLS     time.Duration // TLS handshake
	TTFB    time.Duration // request start to first response byte
	Reused  bool          // connection was reused from the pool
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
		// Setting TLSClientConfig conservatively disables automatic HTTP/2, so
		// opt back in explicitly — otherwise HTTPS/2 servers would be measured
		// over HTTP/1.1 and misrepresent real-world performance.
		ForceAttemptHTTP2: true,
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

// SendRequest sends an HTTP request with the specified method, URL, headers and
// data. Alongside the response it returns a Trace with the connection phase
// timings for that request (zero-valued phases mean the step did not happen,
// e.g. a reused keep-alive connection). The Trace is non-nil whenever err is
// nil; httptrace hooks only fire for the standard transport, so a custom
// RoundTripper yields a zero Trace.
func (c *Client) SendRequest(method, url string, headers map[string]string, data interface{}) (*http.Response, *Trace, error) {
	var body *bytes.Buffer

	// Serialize data to JSON if it's not nil
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, nil, err
		}
		body = bytes.NewBuffer(jsonData)
	} else {
		body = bytes.NewBuffer([]byte{}) // Empty body for GET requests
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, nil, err
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

	tr := &Trace{}
	var start, dnsStart, connectStart, tlsStart time.Time
	trace := &httptrace.ClientTrace{
		DNSStart:             func(httptrace.DNSStartInfo) { dnsStart = time.Now() },
		DNSDone:              func(httptrace.DNSDoneInfo) { tr.DNS = time.Since(dnsStart) },
		ConnectStart:         func(_, _ string) { connectStart = time.Now() },
		ConnectDone:          func(_, _ string, _ error) { tr.Connect = time.Since(connectStart) },
		TLSHandshakeStart:    func() { tlsStart = time.Now() },
		TLSHandshakeDone:     func(tls.ConnectionState, error) { tr.TLS = time.Since(tlsStart) },
		GotConn:              func(info httptrace.GotConnInfo) { tr.Reused = info.Reused },
		GotFirstResponseByte: func() { tr.TTFB = time.Since(start) },
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	start = time.Now()
	resp, err := c.Do(req)
	if err != nil {
		return nil, nil, err
	}
	return resp, tr, nil
}

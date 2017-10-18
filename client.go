package serialized

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client for the Serialized.io API.
type Client struct {
	baseURL   *url.URL
	userAgent string

	pollInterval time.Duration

	accessKey       string
	secretAccessKey string

	httpClient *http.Client
}

// NewClient return a new Serialized.io Client.
func NewClient(opts ...func(*Client)) *Client {
	c := &Client{
		baseURL: &url.URL{
			Scheme: "https",
			Host:   "api.serialized.io",
		},
		userAgent:    "serialized-go/0.1.0",
		pollInterval: 2 * time.Second,
		httpClient:   &http.Client{},
	}

	for _, f := range opts {
		f(c)
	}

	return c
}

// WithBaseURL sets the Client base URL.
func WithBaseURL(rawurl string) func(*Client) {
	return func(c *Client) {
		if u, err := url.Parse(rawurl); err == nil {
			c.baseURL = u
		}
	}
}

// WithAccessKey sets the Client access key for authentication.
func WithAccessKey(key string) func(*Client) {
	return func(c *Client) {
		c.accessKey = key
	}
}

// WithSecretAccessKey sets the Client secret access key for authentication.
func WithSecretAccessKey(key string) func(*Client) {
	return func(c *Client) {
		c.secretAccessKey = key
	}
}

// WithPollInterval sets the interval used for polling the API.
func WithPollInterval(d time.Duration) func(*Client) {
	return func(c *Client) {
		c.pollInterval = d
	}
}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	u, err := c.baseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Serialized-Access-Key", c.accessKey)
	req.Header.Set("Serialized-Secret-Access-Key", c.secretAccessKey)

	return req, nil
}

func (c *Client) do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	req = req.WithContext(ctx)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return nil, err
	}
	defer resp.Body.Close()

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
		if err == io.EOF {
			err = nil
		}
	}

	return resp, err
}

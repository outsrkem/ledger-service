package uias

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	xRequestIdKey   = "X-Request-Id"
	xAuthTokenKey   = "X-Auth-Token"
	checkActionPath = "/v1/uias/action/check" // Permission verification path
)

// NewClientBuilder creates a new client builder
func NewClientBuilder() *ClientBuilder {
	// Set default configuration
	return &ClientBuilder{
		config: Config{
			Timeout:       5 * time.Second,
			SkipTlsVerify: false,
		},
	}
}

// WithEndpoint sets the service endpoint
func (b *ClientBuilder) WithEndpoint(endpoint string) *ClientBuilder {
	b.config.Endpoint = endpoint
	return b
}

// WithUrlPath sets the request path
func (b *ClientBuilder) WithUrlPath(path string) *ClientBuilder {
	b.config.UrlPath = path
	return b
}

// WithSkipTlsVerify sets whether to skip TLS verification
func (b *ClientBuilder) WithSkipTlsVerify(skip bool) *ClientBuilder {
	b.config.SkipTlsVerify = skip
	return b
}

// WithCACertPath sets the CA certificate path
func (b *ClientBuilder) WithCACertPath(path string) *ClientBuilder {
	b.config.CACertPath = path
	return b
}

// WithTimeout sets the request timeout duration
func (b *ClientBuilder) WithTimeout(timeout time.Duration) *ClientBuilder {
	b.config.Timeout = timeout
	return b
}

// Build constructs the client instance
func (b *ClientBuilder) Build() (*Client, error) {
	// Validate required configurations
	if b.config.Endpoint == "" {
		return nil, errors.New("endpoint is required")
	}
	if b.config.UrlPath == "" {
		b.config.UrlPath = checkActionPath
	}
	// Create HTTP client
	httpClient, err := createHttpClient(b.config.SkipTlsVerify, b.config.CACertPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create http client: %w", err)
	}
	// Set timeout
	httpClient.Timeout = b.config.Timeout
	return &Client{config: b.config, httpClient: httpClient}, nil
}

// VerifyAction verifies the action permission
func (c *Client) VerifyAction(ctx context.Context, requestId, token string, rawBody []byte) (*RespData, error) {
	// Build request headers
	headers := map[string]string{
		"Content-Type": "application/json; charset=utf-8",
		xRequestIdKey:  requestId,
		xAuthTokenKey:  token,
	}
	// Build request URL
	url := c.config.Endpoint + c.config.UrlPath
	// Send request
	response, err := sendHttpRequest(ctx, c.httpClient, "POST", url, rawBody, headers, c.config.Timeout)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	// Parse response body
	var resp *RespData
	err = json.Unmarshal(response.Body, &resp)
	if err != nil {
		panic(err)
	}
	// Check response status code
	if response.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}
	return resp, nil
}

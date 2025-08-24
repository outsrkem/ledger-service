package uias

import (
	"net/http"
	"time"
)

// RespData contains the returned data after successful token verification
type RespData struct {
	Metadata struct {
		Message string `json:"message"`
		Time    int64  `json:"time"`
		Ecode   string `json:"ecode"`
	} `json:"metadata"`
	Payload struct {
		Authentication int    `json:"authentication"`
		Error          string `json:"error"`
		Msg            struct {
			Action    string `json:"action"`
			Statement struct {
				Action string `json:"Action"`
				Effect string `json:"Effect"`
			} `json:"Statement"`
		} `json:"msg"`
		User struct {
			Domain struct {
				Id   string `json:"id"`
				Name string `json:"name"`
			} `json:"domain"`
			Id   string `json:"id"`
			Name struct {
				Account string `json:"account"`
			} `json:"name"`
		} `json:"user"`
	} `json:"payload"`
}

// Config contains client configuration parameters
type Config struct {
	Endpoint      string        // Service endpoint (e.g., "https://api.example.com")
	UrlPath       string        // API path to use for requests
	SkipTlsVerify bool          // Whether to skip TLS verification (default: false)
	CACertPath    string        // Path to CA certificate file
	Timeout       time.Duration // Request timeout duration (default: 5 seconds)
}

// Client represents a UIAS service client
type Client struct {
	config     Config
	httpClient *http.Client
}

// ClientBuilder helps construct a Client with specific configurations
type ClientBuilder struct {
	config Config
}

// Package datetimeclient provides a client for interacting with a datetime server.
package datetimeclient

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/codescalersinternships/datetime-client-eyadhussein/pkg/backoff"
)

// Client interface defines the method for getting the current date and time.
type Client interface {
	GetCurrentDateTime() ([]byte, error)
}

// RealClient implements Client and uses a http client for interacting with the datetime server.
type RealClient struct {
	baseUrl string
	port    string
	client  *http.Client
}

// NewRealClient creates and returns a new RealClient instance.
// It uses environment variables for baseUrl and port if not provided.
func NewRealClient(baseUrl, port string, timeout time.Duration) *RealClient {
	if baseUrl == "" {
		baseUrl = os.Getenv("SERVER_URL")
	}
	if port == "" {
		port = os.Getenv("PORT")
	}

	port = ":" + port

	return &RealClient{
		baseUrl: baseUrl,
		port:    port,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// GetCurrentDateTime sends a request to the datetime server and returns the current date and time.
// It uses a backoff strategy for retrying the request in case of failures.
// Returns the response body as a byte slice and any error encountered.
func (c *RealClient) GetCurrentDateTime() ([]byte, error) {
	backoff := backoff.NewRealBackOff(1, 3)
	req, err := http.NewRequest(http.MethodGet, c.baseUrl+c.port+"/datetime", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "text/plain;charset=UTF-8, application/json")

	resp, err := backoff.Retry(func() (*http.Response, error) {
		resp, err := c.client.Do(req)
		return resp, err
	})

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil

}

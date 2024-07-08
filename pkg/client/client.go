package datetimeclient

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/codescalersinternships/datetime-client-eyadhussein/pkg/backoff"
)

type HttpClient interface {
	struct {
		time.Duration
	}
	Do(req *http.Request) (*http.Response, error)
}

type Client interface {
	string
	HttpClient
}

type RealClient struct {
	baseUrl string
	port    string
	client  *http.Client
}

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

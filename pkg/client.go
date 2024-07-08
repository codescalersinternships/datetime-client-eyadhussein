package datetimeclient

import (
	"io"
	"net/http"
	"os"
	"time"
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
	host   string
	port   string
	client *http.Client
}

func NewRealClient(host, port string, timeout time.Duration) *RealClient {
	return &RealClient{
		host: host,
		port: port,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *RealClient) GetCurrentDateTime() ([]byte, error) {
	var host = c.host
	var port = ":" + c.port

	if serverUrlEnv := os.Getenv("SERVER_URL"); serverUrlEnv != "" {
		host = serverUrlEnv
	}
	if portEnv := os.Getenv("PORT"); portEnv != "" {
		port = ":" + portEnv
	}

	url := host + port

	req, err := http.NewRequest(http.MethodGet, url+"/datetime", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "text/plain;charset=UTF-8, application/json")

	resp, err := c.client.Do(req)

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

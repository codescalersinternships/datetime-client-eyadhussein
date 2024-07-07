package datetimeclient

import (
	"io"
	"net/http"
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
	url    string
	client *http.Client
}

func NewRealClient(url string, timeout time.Duration) *RealClient {
	return &RealClient{
		url: url,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *RealClient) GetCurrentDateTime() ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, c.url+"/datetime", nil)
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

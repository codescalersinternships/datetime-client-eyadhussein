package datetimeclient

import (
	"io"
	"net/http"
	"time"
)

type Client struct {
	url    string
	client *http.Client
}

func NewClient(url string, timeout time.Duration) *Client {
	return &Client{
		url: url,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) GetCurrentDateTime() ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, c.url, nil)
	req.Header.Add("Accept", "text/plain;charset=UTF-8, application/json")

	if err != nil {
		return nil, err
	}

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

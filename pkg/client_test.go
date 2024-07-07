package datetimeclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	testTimeout = 3 * time.Second

	validTimeout   = testTimeout - 1*time.Second
	inValidTimeout = testTimeout + 1*time.Second
)

type MockHttpClient struct {
	Client struct {
		Timeout time.Duration
	}
	SleepDuration time.Duration
}

func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	if m.SleepDuration > m.Client.Timeout {
		return &http.Response{StatusCode: http.StatusRequestTimeout}, context.DeadlineExceeded
	}
	return &http.Response{StatusCode: http.StatusOK}, nil
}

type MockClient struct {
	url    string
	client *MockHttpClient
}

func (m *MockClient) GetCurrentDateTime() ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, m.url+"/datetime", nil)
	if err != nil {
		return nil, err
	}
	_, err = m.client.Do(req)
	if err != nil {
		return nil, err
	}
	return []byte("2024-07-04 15:11:44"), nil
}

func NewMockClient(url string, timeout time.Duration, sleepDuration time.Duration) *MockClient {
	return &MockClient{
		url: url,
		client: &MockHttpClient{
			Client: struct{ Timeout time.Duration }{
				Timeout: timeout,
			},
			SleepDuration: sleepDuration,
		},
	}
}
func TestClient_GetCurrentDateTime(t *testing.T) {
	t.Run("handle json format", func(t *testing.T) {
		expected := "2024-07-07 02:39:09"
		jsonExpected, _ := json.Marshal(expected)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonExpected)
		}))

		defer mockServer.Close()

		client := NewRealClient(
			mockServer.URL,
			1*time.Second,
		)

		data, err := client.GetCurrentDateTime()

		if err != nil {
			t.Errorf("expected err to be nil got %v", err)
		}

		var resParsed string
		err = json.Unmarshal(data, &resParsed)

		if err != nil {
			t.Errorf("expected err to be nil got %v", err)
		}

		if resParsed != expected {
			t.Errorf("expected %s but got %s", expected, resParsed)
		}
	})

	t.Run("handle text/plain format", func(t *testing.T) {
		expected := "2024-07-07 02:39:09"
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprint(w, expected)
		}))

		defer mockServer.Close()

		client := NewRealClient(
			mockServer.URL,
			time.Duration(1)*time.Second,
		)

		data, err := client.GetCurrentDateTime()

		if err != nil {
			t.Errorf("expected err to be nil got %v", err)
		}

		if string(data) != expected {
			t.Errorf("expected %s but got %s", expected, data)
		}
	})

	t.Run("valid timeout", func(t *testing.T) {
		client := NewMockClient(
			"", testTimeout, validTimeout,
		)
		_, err := client.GetCurrentDateTime()

		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}
	})

	t.Run("invalid timeout", func(t *testing.T) {

		client := NewMockClient(
			"", testTimeout, inValidTimeout,
		)
		_, err := client.GetCurrentDateTime()

		if err == nil {
			t.Error("expected error but got nil")
		}
	})
}

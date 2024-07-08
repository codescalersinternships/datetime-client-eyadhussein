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
			"",
			1*time.Second,
		)

		data, err := client.GetCurrentDateTime()
		assertNoError(t, err)

		var resParsed string
		err = json.Unmarshal(data, &resParsed)

		assertNoError(t, err)
		assertEqual(t, resParsed, expected)
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
			"",
			time.Duration(1)*time.Second,
		)

		data, err := client.GetCurrentDateTime()

		assertNoError(t, err)
		assertEqual(t, string(data), expected)
	})

	t.Run("valid timeout", func(t *testing.T) {
		client := NewMockClient(
			"", testTimeout, validTimeout,
		)
		_, err := client.GetCurrentDateTime()

		assertNoError(t, err)
	})

	t.Run("invalid timeout", func(t *testing.T) {

		client := NewMockClient(
			"", testTimeout, inValidTimeout,
		)
		_, err := client.GetCurrentDateTime()

		assertError(t, err)
	})

	t.Run("correct endpoint", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/datetime" {
				t.Errorf("expected path /datetime, got %s", r.URL.Path)
			}
			fmt.Fprint(w, "2024-07-07 02:39:09")
		}))
		defer mockServer.Close()

		client := NewRealClient(mockServer.URL, "", 1*time.Second)
		_, err := client.GetCurrentDateTime()

		assertNoError(t, err)
	})

	t.Run("correct Accept header", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accept := r.Header.Get("Accept")
			expected := "text/plain;charset=UTF-8, application/json"
			if accept != expected {
				t.Errorf("expected Accept header %s, got %s", expected, accept)
			}
			fmt.Fprint(w, "2024-07-07 02:39:09")
		}))
		defer mockServer.Close()

		client := NewRealClient(mockServer.URL, "", 1*time.Second)
		_, err := client.GetCurrentDateTime()

		assertNoError(t, err)

	})
}

func assertEqual(t *testing.T, got, want any) {
	t.Helper()
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("expected nil but got %v", err)
	}
}

func assertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Error("expected error but got nil")
	}
}

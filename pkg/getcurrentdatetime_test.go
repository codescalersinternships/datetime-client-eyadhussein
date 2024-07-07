package datetimeclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_GetCurrentDateTime(t *testing.T) {
	t.Run("handle json format", func(t *testing.T) {
		expected := "2024-07-07 02:39:09"
		jsonExpected, _ := json.Marshal(expected)
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonExpected)
		}))

		defer mockServer.Close()

		client := NewClient(
			mockServer.URL,
			time.Duration(1)*time.Second,
		)

		data, err := client.GetCurrentDateTime()

		if err != nil {
			t.Error(err)
		}
		var resParsed string
		err = json.Unmarshal(data, &resParsed)
		if err != nil {
			t.Error(err)
		}

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

		client := NewClient(
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

	t.Run("verify timeout constraint", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(2 * time.Second)
			w.WriteHeader(http.StatusOK)
		}))

		defer mockServer.Close()

		client := NewClient(
			mockServer.URL,
			time.Duration(1)*time.Second,
		)

		_, err := client.GetCurrentDateTime()

		if err == nil {
			t.Errorf("expected timeout error but got nil")
		}
	})
}

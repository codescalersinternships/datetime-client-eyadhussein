package backoff

import (
	"errors"
	"io"
	"log"
	"net/http"
	"testing"
	"time"
)

const maxRetries = 3
const validRetries = 2

type SpySleeper struct {
	Calls int
}

func (s *SpySleeper) Sleep(d time.Duration) {
	s.Calls++
}

type MockBackOff struct {
	Duration time.Duration
	MaxRetry int
	sleeper  *SpySleeper
}

func NewMockBackOff(duration time.Duration, maxRetry int) *MockBackOff {
	return &MockBackOff{
		Duration: duration,
		MaxRetry: maxRetry,
		sleeper:  &SpySleeper{},
	}
}

func (b *MockBackOff) Retry(operation Operation) (*http.Response, error) {
	for i := 0; i < b.MaxRetry; i++ {
		resp, err := operation()

		if err == nil {
			return resp, nil
		}
		b.sleeper.Sleep(b.Duration)
	}
	return &http.Response{}, errors.New("reached maximum retries with no established connection")
}

func TestBackOff_Retry(t *testing.T) {
	t.Run("successful on first try", func(t *testing.T) {
		b := NewMockBackOff(1*time.Second, maxRetries)
		operation := func() (*http.Response, error) {
			return &http.Response{StatusCode: http.StatusOK}, nil
		}

		resp, err := b.Retry(operation)

		assertNoError(t, err)

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected success code but got %d", resp.StatusCode)
		}

		if b.sleeper.Calls != 0 {
			t.Errorf("expected no sleep calls but got %d", b.sleeper.Calls)
		}
	})

	t.Run("successful after retries", func(t *testing.T) {
		b := NewMockBackOff(1*time.Second, maxRetries)
		callCount := 0
		operation := func() (*http.Response, error) {
			callCount++
			if callCount < validRetries {
				return nil, errors.New("temporary error")
			}
			return &http.Response{StatusCode: http.StatusOK}, nil
		}

		resp, err := b.Retry(operation)

		assertNoError(t, err)

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected success code but got %d", resp.StatusCode)
		}

		if b.sleeper.Calls != validRetries-1 {
			t.Errorf("expected sleep to be called 1 time but was called %d times", b.sleeper.Calls)
		}
	})

	t.Run("failure after max retries", func(t *testing.T) {
		b := NewMockBackOff(1*time.Second, maxRetries)
		operation := func() (*http.Response, error) {
			return &http.Response{StatusCode: http.StatusServiceUnavailable}, errors.New("service unavailable try again later")
		}

		_, err := b.Retry(operation)

		if err == nil {
			t.Error("expected error but got nil")
		}
	})
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("expected nil but got %v", err)
	}
}

func ExampleBackOff_Retry() {
	backoff := NewRealBackOff(1*time.Second, 3)
	operation := func() (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK}, nil
	}

	resp, err := backoff.Retry(operation)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(body))
}

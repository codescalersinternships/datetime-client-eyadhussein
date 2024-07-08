// Package backoff provides functionality for retrying operations with a simple backoff strategy.
package backoff

import (
	"errors"
	"net/http"
	"time"
)

// Operation represents a function that returns an HTTP response and an error.
type Operation func() (*http.Response, error)

// BackOff interface defines the method for retrying a request for number of times before failure.
type BackOff interface {
	Retry(operation Operation) (*http.Response, error)
}

// RealBackOff implements a simple backoff strategy with a fixed duration between retries.
type RealBackOff struct {
	Duration time.Duration // time to wait between retry attempts.
	MaxRetry int           // maximum number of retry attempts
}

// NewRealBackOff creates and returns a new RealBackOff instance with the specified duration and maximum retries.
func NewRealBackOff(duration time.Duration, maxRetry int) *RealBackOff {
	return &RealBackOff{
		Duration: duration,
		MaxRetry: maxRetry,
	}
}

// Retry attempts to execute the given operation, retrying up to MaxRetry times with a fixed Duration between attempts.
// It returns the successful HTTP response or an error if all attempts fail.
func (b *RealBackOff) Retry(operation Operation) (*http.Response, error) {
	for i := 0; i < b.MaxRetry; i++ {
		resp, err := operation()

		if err == nil {
			return resp, nil
		}
		time.Sleep(b.Duration)
	}
	return &http.Response{}, errors.New("reached maximum retries with no established connection")
}

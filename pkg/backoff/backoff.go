package backoff

import (
	"errors"
	"net/http"
	"time"
)

type Operation func() (*http.Response, error)

type RealBackOff struct {
	Duration time.Duration
	MaxRetry int
}

func NewRealBackOff(duration time.Duration, maxRetry int) *RealBackOff {
	return &RealBackOff{
		Duration: duration,
		MaxRetry: maxRetry,
	}
}

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

package http

import (
	"net/http"
	"time"

	"github.com/avast/retry-go"
)

type RetryClient struct {
	Attempts int
	Timeout  time.Duration
}

func NewDefaultRetryClient() *RetryClient {
	return &RetryClient{
		Attempts: 5,
		Timeout:  10 * time.Second,
	}
}

func (r RetryClient) Do(req *http.Request) (*http.Response, error) {
	var resp *http.Response

	err := retry.Do(
		func() error {
			client := http.Client{Timeout: r.Timeout}
			var err error
			resp, err = client.Do(req)
			if err != nil {
				return err
			}

			return nil
		},
		retry.LastErrorOnly(true),
		retry.Attempts(uint(r.Attempts)),
		retry.Delay(1*time.Second),
		retry.DelayType(retry.BackOffDelay),
	)

	return resp, err
}

func (r RetryClient) Get(url string) (*http.Response, error) {
	var resp *http.Response

	err := retry.Do(
		func() error {
			client := http.Client{Timeout: r.Timeout}
			var err error
			resp, err = client.Get(url)
			if err != nil {
				return err
			}

			return nil
		},
		retry.LastErrorOnly(true),
		retry.Attempts(uint(r.Attempts)),
		retry.Delay(1*time.Second),
		retry.DelayType(retry.BackOffDelay),
	)

	return resp, err
}

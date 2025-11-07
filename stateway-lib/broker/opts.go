package broker

import "time"

type RequestOptions struct {
	Timeout time.Duration
}

type RequestOption func(*RequestOptions)

func WithTimeout(timeout time.Duration) RequestOption {
	return func(o *RequestOptions) {
		o.Timeout = timeout
	}
}

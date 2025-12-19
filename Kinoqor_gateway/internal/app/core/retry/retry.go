package retry

import (
	"context"
	"errors"
	myHttp "gateway/internal/app/core/http"
	"math"
	"net/http"
	"time"
)

type Retrier struct {
	Base    *myHttp.ClientBase
	Retries int
	Backoff time.Duration
}

func (r *Retrier) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	var lastErr error

	for attempt := 0; attempt <= r.Retries; attempt++ {
		resp, err := r.Base.Do(ctx, req)
		if err == nil && resp.StatusCode < http.StatusInternalServerError {
			return resp, nil
		}

		if err != nil {
			lastErr = err
		} else {
			if resp.Body != nil {
				_ = resp.Body.Close()
			}
			lastErr = errors.New(resp.Status)
		}

		if attempt < r.Retries {
			sleep := r.Backoff * time.Duration(math.Pow(2, float64(attempt)))
			select {
			case <-time.After(sleep):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	}

	return nil, lastErr
}

package retry

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"
	myhttp "users/internal/app/core/http"

	"github.com/google/go-cmp/cmp"
)

type roundTripperFunc func(r *http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func newMockClient(rt roundTripperFunc) *myhttp.ClientBase {
	return &myhttp.ClientBase{
		Client: &http.Client{Transport: rt},
	}
}

func newReq(t *testing.T) *http.Request {
	req, err := http.NewRequest(http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	return req
}

func TestRetrier_Do(t *testing.T) {
	tests := map[string]struct {
		retries   int
		mockCalls []func() (*http.Response, error)
		expected  int
		shouldErr bool
	}{
		"success: success on first try": {
			retries: 0,
			mockCalls: []func() (*http.Response, error){
				func() (*http.Response, error) {
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(bytes.NewBufferString("ok")),
					}, nil
				},
			},
			expected:  1,
			shouldErr: false,
		},

		"success: retry once then success": {
			retries: 1,
			mockCalls: []func() (*http.Response, error){
				func() (*http.Response, error) {
					return nil, errors.New("net")
				},
				func() (*http.Response, error) {
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(bytes.NewBufferString("ok")),
					}, nil
				},
			},
			expected:  2,
			shouldErr: false,
		},

		"invalid: no retry for 4xx": {
			retries: 5,
			mockCalls: []func() (*http.Response, error){
				func() (*http.Response, error) {
					return &http.Response{
						StatusCode: 400,
						Status:     "400 Bad Request",
						Body:       io.NopCloser(bytes.NewBufferString("bad")),
					}, nil
				},
			},
			expected:  1,
			shouldErr: false,
		},

		"error: all retries fail": {
			retries: 2,
			mockCalls: []func() (*http.Response, error){
				func() (*http.Response, error) { return nil, errors.New("x1") },
				func() (*http.Response, error) { return nil, errors.New("x2") },
				func() (*http.Response, error) { return nil, errors.New("x3") },
			},
			expected:  3,
			shouldErr: true,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			callIndex := 0

			rt := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
				fn := test.mockCalls[callIndex]
				callIndex++
				return fn()
			})

			r := &Retrier{
				Base:    newMockClient(rt),
				Retries: test.retries,
				Backoff: 1 * time.Millisecond,
			}

			resp, err := r.Do(context.Background(), newReq(t))

			if test.shouldErr && err == nil {
				t.Fatalf("want error, got nil")
			}
			if !test.shouldErr && err != nil {
				t.Fatalf("did not expect error, got: %v", err)
			}

			if diff := cmp.Diff(test.expected, callIndex); diff != "" {
				t.Errorf("call count mismatch (-expected +got): %v", diff)
			}

			if resp != nil {
				_ = resp.Body.Close()
			}
		})
	}
}

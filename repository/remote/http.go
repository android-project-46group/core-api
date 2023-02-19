package remote

import (
	"errors"
	"net"
	"net/http"
	"time"
)

const (
	attempts = 3
	waitSec  = 2
)

// Customized Transport that is a instance of [retryableRoundTripper] struct.
//
// By default
// - base is [http.DefaultTransport]
// - attempts is 3
// - waitTime is 2 seconds.
func newTransport() http.RoundTripper {
	return &retryableRoundTripper{
		base:     http.DefaultTransport,
		attempts: attempts,
		waitTime: waitSec * time.Second,
	}
}

// retryableRoundTripper implements [http.RoundTripper].
//
// It retries HTTP accesses if needed.
type retryableRoundTripper struct {
	// This is a underlying RoundTripper and used.
	base http.RoundTripper

	// Number of retry attempts.
	attempts int

	// Waiting time between retries.
	//
	// Now this is a fixed value rather than exponential backoff, etc.
	waitTime time.Duration
}

// shouldRetry determines if a retry should be made based on the response.
//
// true if network error or response code is due server-side.
func (r *retryableRoundTripper) shouldRetry(resp *http.Response, err error) bool {
	// network error
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return true
		}
	}

	// retry by response code (err is nil because the communication itself is successful)
	if resp != nil {
		if resp.StatusCode == 429 || (500 <= resp.StatusCode && resp.StatusCode <= 504) {
			return true
		}
	}

	return false
}

// retryableRoundTripper implements RoundTrip.
// RoundTrip is defined in [http.RoundTripper] interface.
//
// Retry HTTP accesses as necessary for the number of retries
// defined in the retryableRoundTripper structure.
func (r *retryableRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var (
		resp *http.Response
		err  error
	)

	for count := 0; count < r.attempts; count++ {
		resp, err = r.base.RoundTrip(req)

		if !r.shouldRetry(resp, err) {
			//nolint:wrapcheck
			return resp, err
		}

		select {
		case <-req.Context().Done():
			//nolint:wrapcheck
			return nil, req.Context().Err()

		case <-time.After(r.waitTime):
			// wait for retry
			continue
		}
	}

	//nolint:wrapcheck
	return resp, err
}

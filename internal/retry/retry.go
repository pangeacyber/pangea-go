package retry

import (
	"context"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

var (
	// A regular expression to match the error returned by net/http when the
	// scheme specified in the URL is invalid. This error isn't typed
	// specifically so we resort to matching on the error string.
	schemeErrorRe = regexp.MustCompile(`unsupported protocol scheme`)

	// A regular expression to match the error returned by net/http when the
	// TLS certificate is not trusted. This error isn't typed
	// specifically so we resort to matching on the error string.
	notTrustedErrorRe = regexp.MustCompile(`certificate is not trusted`)
)

type Manager struct {
	RetryWaitMin time.Duration
	RetryWaitMax time.Duration
	RetryMax     int
}

// CheckPolicy provides a policy for handling retries.
// It is called following each request with the response and error values returned by the http.Client.
// If CheckRetry returns false, the Client stops retrying and returns the response to the caller.
// If CheckRetry returns an error, that error value is returned in lieu of the error from the request.
func (r *Manager) CheckPolicy(ctx context.Context, resp *http.Response, err error) (bool, error) {
	// do not retry on context.Canceled or context.DeadlineExceeded
	if ctx.Err() != nil {
		return false, ctx.Err()
	}
	return baseRetryPolicy(resp, err), nil
}

func baseRetryPolicy(resp *http.Response, err error) bool {
	if err != nil {
		if v, ok := err.(*url.Error); ok {
			// Don't retry if the error was due to an invalid protocol scheme.
			if schemeErrorRe.MatchString(v.Error()) {
				return false
			}

			// Don't retry if the error was due to TLS cert verification failure.
			if notTrustedErrorRe.MatchString(v.Error()) {
				return false
			}
		}
		// The error is likely recoverable so retry.
		return true
	}

	// 429 Too Many Requests is recoverable. Sometimes the server puts
	// available to start processing request from client.
	if resp.StatusCode == http.StatusTooManyRequests {
		return true
	}

	// Check the response code. We retry on 500-range responses to allow the server time to recover
	if resp.StatusCode == 0 || (resp.StatusCode >= 500 && resp.StatusCode != http.StatusNotImplemented) {
		return true
	}
	return false
}

// Backoff will perform exponential backoff based on the attempt number and limited
// by the provided minimum and maximum durations.
func (r *Manager) Backoff(attemptNum int, resp *http.Response) time.Duration {
	mult := math.Pow(2, float64(attemptNum)) * float64(r.RetryWaitMin)
	sleep := time.Duration(mult)
	if float64(sleep) != mult || sleep > r.RetryWaitMax {
		sleep = r.RetryWaitMax
	}
	return sleep
}

package defaults

import (
	"net"
	"net/http"
	"runtime"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

// HTTPTransport returns a new http.Transport with similar default values to
// http.DefaultTransport, but with idle connections and keepalives disabled.
func HTTPTransport() *http.Transport {
	transport := HTTPPooledTransport()
	transport.DisableKeepAlives = true
	transport.MaxIdleConnsPerHost = -1
	return transport
}

// HTTPPooledTransport returns a new http.Transport with similar default values to http.DefaultTransport.
func HTTPPooledTransport() *http.Transport {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
		MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
	}
	return transport
}

// HTTPClient returns a new http.Client with similar default values to
// http.Client, but with a non-shared Transport, idle connections disabled, and
// keepalives disabled.
func HTTPClient() *http.Client {
	return &http.Client{
		Transport: HTTPTransport(),
	}
}

// HTTPPooledClient returns a new http.Client with similar default values to
// http.Client, but with a shared Transport. Do not use this function for
// transient clients as it can leak file descriptors over time. Only use this
// for clients that will be re-used for the same host(s).
func HTTPPooledClient() *http.Client {
	return &http.Client{
		Transport: HTTPPooledTransport(),
	}
}

// HTTPClientWithRetries returns a new http.Client with similar default values to
// http.Client, but with a non-shared Transport, idle connections disabled, and
// keepalives disabled and retries.
func HTTPClientWithRetries() *http.Client {
	cli := &retryablehttp.Client{
		HTTPClient:   HTTPPooledClient(),
		RetryWaitMin: 1 * time.Second,
		RetryWaitMax: 30 * time.Second,
		RetryMax:     4,
		CheckRetry:   retryablehttp.DefaultRetryPolicy,
		Backoff:      retryablehttp.DefaultBackoff,
	}
	return cli.StandardClient()
}

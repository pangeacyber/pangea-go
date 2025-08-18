package option

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/rs/zerolog"
)

type ConfigOption = pangea.ConfigOption

// WithAdditionalHeaders returns a ConfigOption that sets the AdditionalHeaders
// for the client.
func WithAdditionalHeaders(headers map[string]string) ConfigOption {
	return pangea.ConfigOptionFunc(func(c *pangea.Config) error {
		c.AdditionalHeaders = headers
		return nil
	})
}

// WithBaseURLTemplate returns a ConfigOption that sets the BaseURLTemplate for
// the client.
func WithBaseURLTemplate(baseURLTemplate string) ConfigOption {
	return pangea.ConfigOptionFunc(func(c *pangea.Config) error {
		c.BaseURLTemplate = baseURLTemplate
		return nil
	})
}

// WithDomain returns a ConfigOption that sets the API domain for the client.
func WithDomain(domain string) ConfigOption {
	return pangea.ConfigOptionFunc(func(c *pangea.Config) error {
		c.Domain = domain
		return nil
	})
}

// WithHTTPClient returns a ConfigOption that changes the underlying HTTP
// client used, which by default is [http.DefaultClient].
//
// For custom uses cases, it is recommended to provide an [*http.Client] with a
// custom [http.RoundTripper] as its transport, rather than directly
// implementing [HTTPClient].
func WithHTTPClient(client *http.Client) ConfigOption {
	return pangea.ConfigOptionFunc(func(c *pangea.Config) error {
		if client == nil {
			return fmt.Errorf("configoption: custom HTTP client cannot be nil")
		}

		c.HTTPClient = client
		return nil
	})
}

// WithLogger returns a ConfigOption that sets the logger for the client.
func WithLogger(logger *zerolog.Logger) ConfigOption {
	return pangea.ConfigOptionFunc(func(c *pangea.Config) error {
		c.Logger = logger
		return nil
	})
}

// WithMaxRetries returns a ConfigOption that sets the maximum number of retries
// that the client attempts to make. When given 0, the client only makes one
// request. By default, the client retries two times.
//
// WithMaxRetries panics when retries is negative.
func WithMaxRetries(retries int) ConfigOption {
	if retries < 0 {
		panic("configoption: cannot have fewer than 0 retries")
	}
	return pangea.ConfigOptionFunc(func(r *pangea.Config) error {
		r.MaxRetries = retries
		return nil
	})
}

// WithPollResultTimeout returns a ConfigOption that sets the timeout for
// polling results after a HTTP/202.
func WithPollResultTimeout(pollResultTimeout time.Duration) ConfigOption {
	return pangea.ConfigOptionFunc(func(c *pangea.Config) error {
		c.PollResultTimeout = pollResultTimeout
		return nil
	})
}

// WithQueuedRetryEnabled returns a ConfigOption that sets whether or not the
// client should retry queued requests.
func WithQueuedRetryEnabled(queuedRetryEnabled bool) ConfigOption {
	return pangea.ConfigOptionFunc(func(c *pangea.Config) error {
		c.QueuedRetryEnabled = queuedRetryEnabled
		return nil
	})
}

// WithToken returns a ConfigOption that sets the API token for the client.
func WithToken(token string) ConfigOption {
	return pangea.ConfigOptionFunc(func(c *pangea.Config) error {
		c.Token = token
		return nil
	})
}

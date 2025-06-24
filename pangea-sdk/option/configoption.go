package option

import (
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

// WithLogger returns a ConfigOption that sets the logger for the client.
func WithLogger(logger *zerolog.Logger) ConfigOption {
	return pangea.ConfigOptionFunc(func(c *pangea.Config) error {
		c.Logger = logger
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

// WithRetry returns a ConfigOption that sets whether or not the client should
// retry failed requests.
func WithRetry(retry bool) ConfigOption {
	return pangea.ConfigOptionFunc(func(c *pangea.Config) error {
		c.Retry = retry
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

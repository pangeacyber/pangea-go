// go:build integration
package pangeatesting

import (
	"testing"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/internal/defaults"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
)

func IntegrationConfig(t *testing.T, env TestEnvironment) *pangea.Config {
	return &pangea.Config{
		HTTPClient:         defaults.HTTPClient(),
		Domain:             GetTestDomain(t, env),
		Token:              GetTestToken(t, env),
		QueuedRetryEnabled: true,
		PollResultTimeout:  60 * time.Second,
		Retry:              true,
		RetryConfig: &pangea.RetryConfig{
			RetryMax: 4,
		},
	}
}

func IntegrationAuditVaultConfig(t *testing.T, env TestEnvironment) *pangea.Config {
	return &pangea.Config{
		HTTPClient:         defaults.HTTPClient(),
		Domain:             GetTestDomain(t, env),
		Token:              GetVaultSignatureTestToken(t, env),
		QueuedRetryEnabled: true,
		PollResultTimeout:  60 * time.Second,
		Retry:              true,
		RetryConfig: &pangea.RetryConfig{
			RetryMax: 4,
		},
	}
}

func IntegrationCustomSchemaConfig(t *testing.T, env TestEnvironment) *pangea.Config {
	return &pangea.Config{
		HTTPClient: defaults.HTTPClient(),
		Domain:     GetTestDomain(t, env),
		Token:      GetCustomSchemaTestToken(t, env),
	}
}

func IntegrationMultiConfigConfig(t *testing.T, env TestEnvironment) *pangea.Config {
	return &pangea.Config{
		HTTPClient: defaults.HTTPClient(),
		Domain:     GetTestDomain(t, env),
		Token:      GetMultiConfigTestToken(t, env),
	}
}

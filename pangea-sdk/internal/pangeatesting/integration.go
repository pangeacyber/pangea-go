//go:build integration

package pangeatesting

import (
	"testing"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/internal/defaults"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/internal/pangeautil"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
)

func IntegrationConfig(t *testing.T, env TestEnvironment) *pangea.Config {
	return &pangea.Config{
		HTTPClient:         defaults.HTTPClient(),
		Domain:             GetTestDomain(t, env),
		Token:              GetTestToken(t, env),
		QueuedRetryEnabled: true,
		PollResultTimeout:  60 * time.Second,
		Retry:              true,
		Logger:             pangeautil.GetDebugLogger(),
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
		Logger:             pangeautil.GetDebugLogger(),
		RetryConfig: &pangea.RetryConfig{
			RetryMax: 4,
		},
	}
}

func IntegrationCustomSchemaConfig(t *testing.T, env TestEnvironment) *pangea.Config {
	return &pangea.Config{
		HTTPClient:         defaults.HTTPClient(),
		Domain:             GetTestDomain(t, env),
		Token:              GetCustomSchemaTestToken(t, env),
		PollResultTimeout:  60 * time.Second,
		QueuedRetryEnabled: true,
		Logger:             pangeautil.GetDebugLogger(),
	}
}

func IntegrationMultiConfigConfig(t *testing.T, env TestEnvironment) *pangea.Config {
	return &pangea.Config{
		HTTPClient: defaults.HTTPClient(),
		Domain:     GetTestDomain(t, env),
		Token:      GetMultiConfigTestToken(t, env),
		Logger:     pangeautil.GetDebugLogger(),
	}
}

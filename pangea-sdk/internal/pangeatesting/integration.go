// go:build integration
package pangeatesting

import (
	"testing"

	"github.com/pangeacyber/pangea-go/pangea-sdk/internal/defaults"
	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
)

func IntegrationConfig(t *testing.T, env TestEnvironment) *pangea.Config {
	return &pangea.Config{
		HTTPClient: defaults.HTTPClient(),
		Domain:     GetTestDomain(t, env),
		Token:      GetTestToken(t, env),
		Retry:      true,
		RetryConfig: &pangea.RetryConfig{
			BackOff:  1.0,
			RetryMax: 10,
		},
	}
}

func IntegrationAuditVaultConfig(t *testing.T, env TestEnvironment) *pangea.Config {
	return &pangea.Config{
		HTTPClient: defaults.HTTPClient(),
		Domain:     GetTestDomain(t, env),
		Token:      GetVaultSignatureTestToken(t, env),
		Retry:      true,
		RetryConfig: &pangea.RetryConfig{
			BackOff:  1.0,
			RetryMax: 10,
		},
	}
}

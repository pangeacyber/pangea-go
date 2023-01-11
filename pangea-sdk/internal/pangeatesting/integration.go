// go:build integration
package pangeatesting

import (
	"testing"

	"github.com/pangeacyber/pangea-go/pangea-sdk/internal/defaults"
	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
)

var (
	integrationConfig = &pangea.Config{
		HTTPClient: defaults.HTTPClient(),
	}
)

func IntegrationConfig(t *testing.T, env TestEnvironment) *pangea.Config {
	integrationConfig.Domain = GetTestDomain(t, env)
	integrationConfig.Token = GetTestToken(t, env)
	return integrationConfig
}

func IntegrationAuditVaultConfig(t *testing.T, env TestEnvironment) *pangea.Config {
	integrationConfig.Domain = GetTestDomain(t, env)
	integrationConfig.Token = GetVaultSignatureTestToken(t, env)
	return integrationConfig
}

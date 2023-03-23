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
	}
}

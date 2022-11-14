// go:build integration
package pangeatesting

import (
	"os"
	"sync"
	"testing"

	"github.com/pangeacyber/go-pangea/internal/defaults"
	"github.com/pangeacyber/go-pangea/pangea"
)

const (
	integrationEnpointEnvVar = "PANGEA_INTEGRATION_DOMAIN"
)

var (
	onceInit          = sync.Once{}
	integrationConfig = &pangea.Config{
		HTTPClient: defaults.HTTPClient(),
	}
)

func initConfig() {
	integrationConfig.Domain = os.Getenv(integrationEnpointEnvVar)
}

func IntegrationConfig(t *testing.T) *pangea.Config {
	onceInit.Do(initConfig)
	if integrationConfig.Domain == "" {
		t.Skip("set " + integrationEnpointEnvVar + " env variables to run this test")
	}
	return integrationConfig
}

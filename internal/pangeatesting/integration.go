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
	integrationTokenEnvVar   = "PANGEA_TEST_INTEGRATION_TOKEN"
	integrationEnpointEnvVar = "PANGEA_TEST_INTEGRATION_ENDPOINT"
)

var (
	onceInit          = sync.Once{}
	integrationConfig = &pangea.Config{
		HTTPClient: defaults.HTTPClient(),
	}
)

func initConfig() {
	integrationConfig.Token = os.Getenv("PANGEA_TEST_INTEGRATION_TOKEN")
	integrationConfig.Domain = os.Getenv("PANGEA_TEST_INTEGRATION_ENDPOINT")
}

func IntegrationConfig(t *testing.T) *pangea.Config {
	onceInit.Do(initConfig)
	if integrationConfig.Token == "" || integrationConfig.Domain == "" {
		t.Skip("set PANGEA_TEST_INTEGRATION_TOKEN and PANGEA_TEST_INTEGRATION_ENDPOINT env variables to run this test")
	}
	return integrationConfig
}

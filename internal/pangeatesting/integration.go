//go:build integration
// +build integration

package pangeatesting

import (
	"os"

	"github.com/pangeacyber/go-pangea/internal/defaults"
	"github.com/pangeacyber/go-pangea/pangea"
)

var IntegrationConfig = &pangea.Config{
	HTTPClient: defaults.HTTPClient(),
}

func init() {
	IntegrationConfig.Token = os.Getenv("PANGEA_TEST_INTEGRAGION_TOKEN")
	IntegrationConfig.Endpoint = os.Getenv("PANGEA_TEST_INTEGRATION_ENDPOINT")
}

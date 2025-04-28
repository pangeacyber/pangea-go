//go:build integration

package pangeatesting

import (
	"testing"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/internal/pangeautil"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/option"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
)

func IntegrationConfig(t *testing.T, env TestEnvironment) *pangea.Config {
	config, err := pangea.NewConfig(
		option.WithDomain(GetTestDomain(t, env)),
		option.WithLogger(pangeautil.GetDebugLogger()),
		option.WithPollResultTimeout(60*time.Second),
		option.WithQueuedRetryEnabled(true),
		option.WithToken(GetTestToken(t, env)),
	)
	if err != nil {
		panic(err)
	}
	return config
}

func IntegrationAuditVaultConfig(t *testing.T, env TestEnvironment) *pangea.Config {
	config := IntegrationConfig(t, env)
	if err := config.Apply(option.WithToken(GetVaultSignatureTestToken(t, env))); err != nil {
		panic(err)
	}
	return config
}

func IntegrationCustomSchemaConfig(t *testing.T, env TestEnvironment) *pangea.Config {
	config := IntegrationConfig(t, env)
	if err := config.Apply(option.WithToken(GetCustomSchemaTestToken(t, env))); err != nil {
		panic(err)
	}
	return config
}

func IntegrationMultiConfigConfig(t *testing.T, env TestEnvironment) *pangea.Config {
	config := IntegrationConfig(t, env)
	if err := config.Apply(option.WithToken(GetMultiConfigTestToken(t, env))); err != nil {
		panic(err)
	}
	return config
}

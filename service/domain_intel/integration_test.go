// go:build integration
package domain_intel_test

import (
	"context"
	"testing"
	"time"

	"github.com/pangeacyber/go-pangea/internal/pangeatesting"
	"github.com/pangeacyber/go-pangea/pangea"
	"github.com/pangeacyber/go-pangea/service/domain_intel"
	"github.com/stretchr/testify/assert"
)

func Test_Integration_DomainLookup(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfgToken := pangeatesting.GetEnvVarOrSkip(t, "DOMAIN_INTEL_INTEGRATION_CONFIG_TOKEN")
	cfg := &pangea.Config{
		CfgToken: cfgToken,
	}
	cfg = cfg.Copy(pangeatesting.IntegrationConfig(t))
	domainintel, _ := domain_intel.New(cfg)

	input := &domain_intel.DomainLookupInput{
		Domain:   "teoghehofuuxo.su",
		Raw:      true,
		Verbose:  true,
		Provider: "crowdstrike",
	}
	out, err := domainintel.Lookup(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, out)
	assert.NotNil(t, out.Result.Data)
	assert.Equal(t, out.Result.Data.Verdict, "malicious")
}

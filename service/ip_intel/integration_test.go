// go:build integration
package ip_intel_test

import (
	"context"
	"testing"
	"time"

	"github.com/pangeacyber/go-pangea/internal/pangeatesting"
	"github.com/pangeacyber/go-pangea/pangea"
	"github.com/pangeacyber/go-pangea/service/ip_intel"
	"github.com/stretchr/testify/assert"
)

func Test_Integration_IpLookup(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfgToken := pangeatesting.GetEnvVarOrSkip(t, "IP_INTEL_INTEGRATION_CONFIG_TOKEN")
	cfg := &pangea.Config{
		CfgToken: cfgToken,
	}
	cfg = cfg.Copy(pangeatesting.IntegrationConfig(t))
	ipintel, _ := ip_intel.New(cfg)

	input := &ip_intel.IpLookupInput{
		Ip:       "93.231.182.110",
		Raw:      true,
		Verbose:  true,
		Provider: "crowdstrike",
	}
	out, err := ipintel.Lookup(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, out)
	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.Data)
	assert.Equal(t, out.Result.Data.Verdict, "malicious")
}

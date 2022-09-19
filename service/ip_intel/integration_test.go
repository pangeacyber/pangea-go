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
		ConfigID: cfgToken,
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

// Unknown IP
func Test_Integration_IpLookup_2(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfgToken := pangeatesting.GetEnvVarOrSkip(t, "IP_INTEL_INTEGRATION_CONFIG_TOKEN")
	cfg := &pangea.Config{
		ConfigID: cfgToken,
	}
	cfg = cfg.Copy(pangeatesting.IntegrationConfig(t))
	ipintel, _ := ip_intel.New(cfg)

	input := &ip_intel.IpLookupInput{
		Ip:       "4.4.8.8",
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

func Test_Integration_IpLookup_Error_BadIPFormat_1(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfgToken := pangeatesting.GetEnvVarOrSkip(t, "IP_INTEL_INTEGRATION_CONFIG_TOKEN")
	cfg := &pangea.Config{
		ConfigID: cfgToken,
	}
	cfg = cfg.Copy(pangeatesting.IntegrationConfig(t))
	ipintel, _ := ip_intel.New(cfg)

	input := &ip_intel.IpLookupInput{
		Ip:       "93.231.182.300",
		Raw:      true,
		Verbose:  true,
		Provider: "crowdstrike",
	}
	out, err := ipintel.Lookup(ctx, input)

	// FIXME: This should crash, uncomment and fix once enpoint is fixed
	assert.NoError(t, err)
	assert.NotNil(t, out)
	// apiErr := err.(*pangea.APIError)
	// assert.Equal(t, len(apiErr.PangeaErrors.Errors), 1)
	// assert.Equal(t, apiErr.PangeaErrors.Errors[0].Code, "BelowMinLength")
	// assert.Equal(t, apiErr.PangeaErrors.Errors[0].Detail, "'message' cannot have less than 1 characters")
	// assert.Equal(t, apiErr.PangeaErrors.Errors[0].Source, "/event/message")
}

func Test_Integration_IpLookup_Error_BadIPFormat_2(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfgToken := pangeatesting.GetEnvVarOrSkip(t, "IP_INTEL_INTEGRATION_CONFIG_TOKEN")
	cfg := &pangea.Config{
		ConfigID: cfgToken,
	}
	cfg = cfg.Copy(pangeatesting.IntegrationConfig(t))
	ipintel, _ := ip_intel.New(cfg)

	input := &ip_intel.IpLookupInput{
		Ip:       "notanip",
		Raw:      true,
		Verbose:  true,
		Provider: "crowdstrike",
	}
	out, err := ipintel.Lookup(ctx, input)

	// FIXME: This should crash, uncomment once service is fixed
	assert.NoError(t, err)
	assert.NotNil(t, out)

	// apiErr := err.(*pangea.APIError)
	// assert.Equal(t, len(apiErr.PangeaErrors.Errors), 1)
	// assert.Equal(t, apiErr.PangeaErrors.Errors[0].Code, "BelowMinLength")
	// assert.Equal(t, apiErr.PangeaErrors.Errors[0].Detail, "'message' cannot have less than 1 characters")
	// assert.Equal(t, apiErr.PangeaErrors.Errors[0].Source, "/event/message")
}

func Test_Integration_IpLookup_Error_BadToken(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfgToken := pangeatesting.GetEnvVarOrSkip(t, "IP_INTEL_INTEGRATION_CONFIG_TOKEN")
	cfg := &pangea.Config{
		ConfigID: cfgToken,
	}
	cfg = cfg.Copy(pangeatesting.IntegrationConfig(t))
	cfg.Token = "notarealtoken"
	ipintel, _ := ip_intel.New(cfg)

	input := &ip_intel.IpLookupInput{
		Ip:       "93.231.182.110",
		Raw:      true,
		Verbose:  true,
		Provider: "crowdstrike",
	}

	out, err := ipintel.Lookup(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, out)
	apiErr := err.(*pangea.APIError)
	assert.Equal(t, apiErr.Err.Error(), "API error: Not authorized to access this resource.")
}

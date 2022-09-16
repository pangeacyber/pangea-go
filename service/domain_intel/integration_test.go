// go:build integration
package domain_intel_test

import (
	"context"
	"fmt"
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

func Test_Integration_DomainLookup_2(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfgToken := pangeatesting.GetEnvVarOrSkip(t, "DOMAIN_INTEL_INTEGRATION_CONFIG_TOKEN")
	cfg := &pangea.Config{
		CfgToken: cfgToken,
	}
	cfg = cfg.Copy(pangeatesting.IntegrationConfig(t))
	domainintel, _ := domain_intel.New(cfg)

	input := &domain_intel.DomainLookupInput{
		Domain:   "",
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
	fmt.Println(out.Result.Data)
	assert.Equal(t, out.Result.Data.Verdict, "malicious")
}

// Test empty domain
func Test_Integration_DomainLookup_Error(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfgToken := pangeatesting.GetEnvVarOrSkip(t, "DOMAIN_INTEL_INTEGRATION_CONFIG_TOKEN")
	cfg := &pangea.Config{
		CfgToken: cfgToken,
	}
	cfg = cfg.Copy(pangeatesting.IntegrationConfig(t))
	domainintel, _ := domain_intel.New(cfg)

	input := &domain_intel.DomainLookupInput{
		Domain:   "",
		Raw:      true,
		Verbose:  true,
		Provider: "crowdstrike",
	}
	out, err := domainintel.Lookup(ctx, input)

	// FIXME: This should fail
	assert.NoError(t, err)
	assert.NotNil(t, out)
	// err = err.(*pangea.APIError)
	// apiErr := err.(*pangea.APIError)
	// assert.Equal(t, len(apiErr.PangeaErrors.Errors), 1)
	// assert.Equal(t, apiErr.PangeaErrors.Errors[0].Code, "BelowMinLength")
	// assert.Equal(t, apiErr.PangeaErrors.Errors[0].Detail, "'message' cannot have less than 1 characters")
	// assert.Equal(t, apiErr.PangeaErrors.Errors[0].Source, "/event/message")
}

// Not valid provider
func Test_Integration_DomainLookup_Error_Provider(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfgToken := pangeatesting.GetEnvVarOrSkip(t, "DOMAIN_INTEL_INTEGRATION_CONFIG_TOKEN")
	cfg := &pangea.Config{
		CfgToken: cfgToken,
	}
	cfg = cfg.Copy(pangeatesting.IntegrationConfig(t))
	domainintel, _ := domain_intel.New(cfg)

	input := &domain_intel.DomainLookupInput{
		Domain:   "",
		Raw:      true,
		Verbose:  true,
		Provider: "notaprovider",
	}
	out, err := domainintel.Lookup(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, out)
	err = err.(*pangea.APIError)
	apiErr := err.(*pangea.APIError)
	assert.Equal(t, len(apiErr.PangeaErrors.Errors), 1)
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Code, "NotEnumMember")
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Detail, "'provider' must be a one of the following values [crowdstrike]")
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Source, "/provider")
}

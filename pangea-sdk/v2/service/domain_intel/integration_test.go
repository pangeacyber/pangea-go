// go:build integration
package domain_intel_test

import (
	"context"
	"testing"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/internal/pangeatesting"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/domain_intel"
	"github.com/stretchr/testify/assert"
)

const (
	testingEnvironment = pangeatesting.Live
)

func intelDomainIntegrationCfg(t *testing.T) *pangea.Config {
	t.Helper()
	return pangeatesting.IntegrationConfig(t, testingEnvironment)
}

// Reputation domain malicious
func Test_Integration_DomainReputation(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := intelDomainIntegrationCfg(t)
	intelcli := domain_intel.New(cfg)

	input := &domain_intel.DomainReputationRequest{
		Domain:   "737updatesboeing.com",
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "domaintools",
	}
	resp, err := intelcli.Reputation(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result.Data)
	assert.Equal(t, resp.Result.Data.Verdict, "malicious")
}

// Reputation domain unknown
func Test_Integration_DomainReputation_2(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := intelDomainIntegrationCfg(t)
	intelcli := domain_intel.New(cfg)

	input := &domain_intel.DomainReputationRequest{
		Domain:   "google.com",
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "domaintools",
	}

	resp, err := intelcli.Reputation(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result.Data)
	assert.Equal(t, resp.Result.Data.Verdict, "benign")
}

// Test empty domain
func Test_Integration_DomainReputation_Error(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := intelDomainIntegrationCfg(t)
	intelcli := domain_intel.New(cfg)

	input := &domain_intel.DomainReputationRequest{
		Domain:   "",
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "domaintools",
	}

	resp, err := intelcli.Reputation(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, resp)
	err = err.(*pangea.APIError)
	apiErr := err.(*pangea.APIError)
	assert.Equal(t, len(apiErr.PangeaErrors.Errors), 1)
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Code, "BadFormatHostname")
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Detail, "'domain' must be a valid RFC1123 hostname")
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Source, "/domain")
}

// Bad auth token
func Test_Integration_DomainReputation_Error_BadAuthToken(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := intelDomainIntegrationCfg(t)
	cfg.Token = "notavalidtoken"
	intelcli := domain_intel.New(cfg)

	input := &domain_intel.DomainReputationRequest{
		Domain:   "737updatesboeing.com",
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "domaintools",
	}
	resp, err := intelcli.Reputation(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, resp)
	apiErr := err.(*pangea.APIError)
	assert.Equal(t, apiErr.Err.Error(), "API error: Not authorized to access this resource.")
}

// Not valid provider
func Test_Integration_DomainReputation_Error_Provider(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := intelDomainIntegrationCfg(t)
	intelcli := domain_intel.New(cfg)

	input := &domain_intel.DomainReputationRequest{
		Domain:   "737updatesboeing.com",
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "notaprovider",
	}
	resp, err := intelcli.Reputation(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, resp)
	err = err.(*pangea.APIError)
	apiErr := err.(*pangea.APIError)
	assert.Equal(t, len(apiErr.PangeaErrors.Errors), 1)
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Code, "NotEnumMember")
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Source, "/provider")
}
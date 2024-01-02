// go:build integration
package domain_intel_test

import (
	"context"
	"testing"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/internal/pangeatesting"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/service/domain_intel"
	"github.com/stretchr/testify/assert"
)

const (
	testingEnvironment = pangeatesting.Develop
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
		Provider: "crowdstrike",
	}
	resp, err := intelcli.Reputation(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result.Data)
	assert.Equal(t, resp.Result.Data.Verdict, "malicious")
}

func Test_Integration_DomainReputationBulk(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := intelDomainIntegrationCfg(t)
	intelcli := domain_intel.New(cfg)

	input := &domain_intel.DomainReputationBulkRequest{
		Domains:  []string{"pemewizubidob.cafij.co.za", "redbomb.com.tr", "kmbk8.hicp.net"},
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "crowdstrike",
	}
	resp, err := intelcli.ReputationBulk(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result.Data)
	assert.Equal(t, len(resp.Result.Data), 3)
	for _, di := range resp.Result.Data {
		assert.NotEmpty(t, di.Category)
		assert.NotEmpty(t, di.Score)
		assert.NotEmpty(t, di.Verdict)
	}
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
		Provider: "crowdstrike",
	}

	resp, err := intelcli.Reputation(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result.Data)
	assert.NotEmpty(t, resp.Result.Data.Verdict)
}

// func Test_Integration_DomainWhoIs(t *testing.T) {
// 	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancelFn()

// 	cfg := intelDomainIntegrationCfg(t)
// 	intelcli := domain_intel.New(cfg)

// 	input := &domain_intel.DomainWhoIsRequest{
// 		Domain:   "737updatesboeing.com",
// 		Raw:      pangea.Bool(true),
// 		Verbose:  pangea.Bool(true),
// 		Provider: "whoisxml",
// 	}
// 	resp, err := intelcli.WhoIs(ctx, input)
// 	if err != nil {
// 		t.Fatalf("expected no error got: %v", err)
// 	}

// 	assert.NotNil(t, resp)
// 	assert.NotNil(t, resp.Result.Data)
// 	assert.NotEmpty(t, resp.Result.Data.DomainName)
// 	assert.NotEmpty(t, resp.Result.Data.DomainAvailability)
// }

// Reputation domain unknown
func Test_Integration_DomainReputation_NotFound(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := intelDomainIntegrationCfg(t)
	intelcli := domain_intel.New(cfg)

	input := &domain_intel.DomainReputationRequest{
		Domain:   "thisshouldbeafakedomain123asd.com",
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "crowdstrike",
	}

	resp, err := intelcli.Reputation(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result.Data)
	assert.NotEmpty(t, resp.Result.Data.Verdict)
	assert.NotNil(t, resp.Result.Data.Category)
	assert.NotEmpty(t, resp.Result.Data.Score)
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
		Provider: "crowdstrike",
	}

	resp, err := intelcli.Reputation(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, resp)
	err = err.(*pangea.APIError)
	_, ok := err.(*pangea.APIError)
	assert.True(t, ok)
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
		Provider: "crowdstrike",
	}
	resp, err := intelcli.Reputation(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, resp)
	_, ok := err.(*pangea.APIError)
	assert.True(t, ok)
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
	_, ok := err.(*pangea.APIError)
	assert.True(t, ok)
}

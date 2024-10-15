//go:build integration

package url_intel_test

import (
	"context"
	"testing"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/internal/pangeatesting"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/service/url_intel"
	"github.com/stretchr/testify/assert"
)

var testingEnvironment = pangeatesting.LoadTestEnvironment("url-intel", pangeatesting.Live)

func Test_Integration_UrlReputation(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	intelcli := url_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &url_intel.UrlReputationRequest{
		Url:      "http://113.235.101.11:54384",
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "crowdstrike",
	}
	resp, err := intelcli.Reputation(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	assert.NotNil(t, resp.Result.Data)
	assert.Equal(t, resp.Result.Data.Verdict, "malicious")
}

func Test_Integration_UrlReputationBulk(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	intelcli := url_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &url_intel.UrlReputationBulkRequest{
		Urls:     []string{"http://113.235.101.11:54384", "http://45.14.49.109:54819"},
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "crowdstrike",
	}
	resp, err := intelcli.ReputationBulk(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	assert.NotNil(t, resp.Result.Data)
	assert.Equal(t, len(resp.Result.Data), 2)
}

func Test_Integration_UrlReputation_NotFound(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	intelcli := url_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &url_intel.UrlReputationRequest{
		Url:      "http://113.235.101.11:54384",
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "crowdstrike",
	}
	resp, err := intelcli.Reputation(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	assert.NotNil(t, resp.Result.Data)
	assert.NotEmpty(t, resp.Result.Data.Verdict)
	assert.NotNil(t, resp.Result.Data.Category)
	assert.NotEmpty(t, resp.Result.Data.Score)
}

func Test_Integration_UrlReputation_2(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	intelcli := url_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &url_intel.UrlReputationRequest{
		Url:      "http://google.com",
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "crowdstrike",
	}
	resp, err := intelcli.Reputation(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	assert.NotNil(t, resp.Result.Data)
	assert.Equal(t, resp.Result.Data.Verdict, "unknown")
}

func Test_Integration_UrlReputation_Error_BadToken(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := pangeatesting.IntegrationConfig(t, testingEnvironment)
	cfg.Token = "notarealtoken"
	intelcli := url_intel.New(cfg)

	input := &url_intel.UrlReputationRequest{
		Url:      "http://113.235.101.11:54384",
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "crowdstrike",
	}

	resp, err := intelcli.Reputation(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, resp)
	apiErr := err.(*pangea.APIError)
	assert.Equal(t, apiErr.Err.Error(), "API error: Not authorized to access this resource.")
}

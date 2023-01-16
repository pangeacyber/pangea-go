// go:build integration
package url_intel_test

import (
	"context"
	"testing"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/internal/pangeatesting"
	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/service/url_intel"
	"github.com/stretchr/testify/assert"
)

const (
	testingEnvironment = pangeatesting.Live
)

func Test_Integration_UrlLookup(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	urlintel := url_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &url_intel.UrlLookupRequest{
		Url:      "http://113.235.101.11:54384",
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "crowdstrike",
	}
	out, err := urlintel.Lookup(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, out)
	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.Data)
	assert.Equal(t, out.Result.Data.Verdict, "malicious")
}

func Test_Integration_UrlLookup_2(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	urlintel := url_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &url_intel.UrlLookupRequest{
		Url:      "http://google.com",
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "crowdstrike",
	}
	out, err := urlintel.Lookup(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, out)
	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.Data)
	assert.Equal(t, out.Result.Data.Verdict, "unknown")
}

func Test_Integration_UrlLookup_Error_BadToken(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := pangeatesting.IntegrationConfig(t, testingEnvironment)
	cfg.Token = "notarealtoken"
	urlintel := url_intel.New(cfg)

	input := &url_intel.UrlLookupRequest{
		Url:      "http://113.235.101.11:54384",
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "crowdstrike",
	}

	out, err := urlintel.Lookup(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, out)
	apiErr := err.(*pangea.APIError)
	assert.Equal(t, apiErr.Err.Error(), "API error: Not authorized to access this resource.")
}

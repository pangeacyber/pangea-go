// go:build integration
package ip_intel_test

import (
	"context"
	"testing"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/internal/pangeatesting"
	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/service/ip_intel"
	"github.com/stretchr/testify/assert"
)

func Test_Integration_IpLookup(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	ipintel := ip_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

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

	ipintel := ip_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &ip_intel.IpLookupInput{
		Ip:       "8.8.4.4",
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
	assert.Equal(t, out.Result.Data.Verdict, "unknown")
}

func Test_Integration_IpLookup_Error_BadIPFormat_1(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	ipintel := ip_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &ip_intel.IpLookupInput{
		Ip:       "93.231.182.300",
		Raw:      true,
		Verbose:  true,
		Provider: "crowdstrike",
	}
	out, err := ipintel.Lookup(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, out)
	apiErr := err.(*pangea.APIError)
	assert.Equal(t, len(apiErr.PangeaErrors.Errors), 1)
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Code, "BadFormatIPAddress")
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Detail, "'ip' must be a valid IPv4 or IPv6 address")
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Source, "/ip")
}

func Test_Integration_IpLookup_Error_BadIPFormat_2(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	ipintel := ip_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &ip_intel.IpLookupInput{
		Ip:       "notanip",
		Raw:      true,
		Verbose:  true,
		Provider: "crowdstrike",
	}
	out, err := ipintel.Lookup(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, out)
	apiErr := err.(*pangea.APIError)
	assert.Equal(t, len(apiErr.PangeaErrors.Errors), 1)
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Code, "BadFormatIPAddress")
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Detail, "'ip' must be a valid IPv4 or IPv6 address")
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Source, "/ip")
}

func Test_Integration_IpLookup_Error_BadToken(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := pangeatesting.IntegrationConfig(t, testingEnvironment)
	cfg.Token = "notarealtoken"
	ipintel := ip_intel.New(cfg)

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

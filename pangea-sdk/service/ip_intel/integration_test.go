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

const (
	testingEnvironment = pangeatesting.Live
)

func Test_Integration_IpLookup(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	ipintel := ip_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &ip_intel.IpLookupRequest{
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

	input := &ip_intel.IpLookupRequest{
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

func Test_Integration_IpLookup_DefaultProvider(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	ipintel := ip_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &ip_intel.IpLookupRequest{
		Ip: "93.231.182.110",
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

func Test_Integration_IpGeolocate(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	ipintel := ip_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &ip_intel.IpGeolocateRequest{
		Ip:       "93.231.182.110",
		Raw:      true,
		Verbose:  true,
		Provider: "digitalenvoy",
	}
	out, err := ipintel.Geolocate(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, out)
	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.Data)
	assert.Equal(t, out.Result.Data.Country, "Federal Republic Of Germany")
	assert.Equal(t, out.Result.Data.City, "unna")
	assert.Equal(t, out.Result.Data.PostalCode, "59425")
}

func Test_Integration_IpGeolocate_DefaultProvider(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	ipintel := ip_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &ip_intel.IpGeolocateRequest{
		Ip: "93.231.182.110",
	}
	out, err := ipintel.Geolocate(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, out)
	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.Data)
	assert.Equal(t, out.Result.Data.Country, "Federal Republic Of Germany")
	assert.Equal(t, out.Result.Data.City, "unna")
	assert.Equal(t, out.Result.Data.PostalCode, "59425")
}

func Test_Integration_IpGetDomain(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	ipintel := ip_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &ip_intel.IpDomainRequest{
		Ip:       "24.235.114.61",
		Raw:      true,
		Verbose:  true,
		Provider: "digitalenvoy",
	}
	out, err := ipintel.GetDomain(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, out)
	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.Data)
	assert.True(t, out.Result.Data.DomainFound)
	assert.Equal(t, out.Result.Data.Domain, "rogers.com")
}

func Test_Integration_IpGetDomain_DefaultProvider(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	ipintel := ip_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &ip_intel.IpDomainRequest{
		Ip: "24.235.114.61",
	}
	out, err := ipintel.GetDomain(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, out)
	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.Data)
	assert.True(t, out.Result.Data.DomainFound)
	assert.Equal(t, out.Result.Data.Domain, "rogers.com")
}

func Test_Integration_IpIsVPN(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	ipintel := ip_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &ip_intel.IpVPNRequest{
		Ip:       "2.56.189.74",
		Raw:      true,
		Verbose:  true,
		Provider: "digitalenvoy",
	}
	out, err := ipintel.IsVPN(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, out)
	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.Data)
	assert.True(t, out.Result.Data.IsVPN)
}

func Test_Integration_IpIsVPN_DefaultProvider(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	ipintel := ip_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &ip_intel.IpVPNRequest{
		Ip: "2.56.189.74",
	}
	out, err := ipintel.IsVPN(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, out)
	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.Data)
	assert.True(t, out.Result.Data.IsVPN)
}

func Test_Integration_IpIsProxy(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	ipintel := ip_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &ip_intel.IpProxyRequest{
		Ip:       "1.0.136.28",
		Raw:      true,
		Verbose:  true,
		Provider: "digitalenvoy",
	}
	out, err := ipintel.IsProxy(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, out)
	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.Data)
	assert.True(t, out.Result.Data.IsProxy)
}

func Test_Integration_IpIsProxy_DefaultProvider(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	ipintel := ip_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &ip_intel.IpProxyRequest{
		Ip: "1.0.136.28",
	}
	out, err := ipintel.IsProxy(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, out)
	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.Data)
	assert.True(t, out.Result.Data.IsProxy)
}

func Test_Integration_IpLookup_Error_BadIPFormat_1(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	ipintel := ip_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &ip_intel.IpLookupRequest{
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

	input := &ip_intel.IpLookupRequest{
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

	input := &ip_intel.IpLookupRequest{
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

func Test_Integration_IpReputation(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	ipintel := ip_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &ip_intel.IpReputationRequest{
		Ip:       "93.231.182.110",
		Raw:      true,
		Verbose:  true,
		Provider: "crowdstrike",
	}
	out, err := ipintel.Reputation(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, out)
	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.Data)
	assert.Equal(t, out.Result.Data.Verdict, "malicious")
}

// Unknown IP
func Test_Integration_IpReputation_2(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	ipintel := ip_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &ip_intel.IpReputationRequest{
		Ip:       "8.8.4.4",
		Raw:      true,
		Verbose:  true,
		Provider: "crowdstrike",
	}
	out, err := ipintel.Reputation(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, out)
	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.Data)
	assert.Equal(t, out.Result.Data.Verdict, "unknown")
}

func Test_Integration_IpReputation_Error_BadIPFormat_1(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	ipintel := ip_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &ip_intel.IpReputationRequest{
		Ip:       "93.231.182.300",
		Raw:      true,
		Verbose:  true,
		Provider: "crowdstrike",
	}
	out, err := ipintel.Reputation(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, out)
	apiErr := err.(*pangea.APIError)
	assert.Equal(t, len(apiErr.PangeaErrors.Errors), 1)
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Code, "BadFormatIPAddress")
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Detail, "'ip' must be a valid IPv4 or IPv6 address")
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Source, "/ip")
}

func Test_Integration_IpReputation_Error_BadIPFormat_2(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	ipintel := ip_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &ip_intel.IpReputationRequest{
		Ip:       "notanip",
		Raw:      true,
		Verbose:  true,
		Provider: "crowdstrike",
	}
	out, err := ipintel.Reputation(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, out)
	apiErr := err.(*pangea.APIError)
	assert.Equal(t, len(apiErr.PangeaErrors.Errors), 1)
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Code, "BadFormatIPAddress")
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Detail, "'ip' must be a valid IPv4 or IPv6 address")
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Source, "/ip")
}

func Test_Integration_IpReputation_Error_BadToken(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := pangeatesting.IntegrationConfig(t, testingEnvironment)
	cfg.Token = "notarealtoken"
	ipintel := ip_intel.New(cfg)

	input := &ip_intel.IpReputationRequest{
		Ip:       "93.231.182.110",
		Raw:      true,
		Verbose:  true,
		Provider: "crowdstrike",
	}

	out, err := ipintel.Reputation(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, out)
	apiErr := err.(*pangea.APIError)
	assert.Equal(t, apiErr.Err.Error(), "API error: Not authorized to access this resource.")
}

// go:build integration
package ip_intel_test

import (
	"context"
	"testing"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/internal/pangeatesting"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/ip_intel"
	"github.com/stretchr/testify/assert"
)

const (
	testingEnvironment = pangeatesting.Live
)

func Test_Integration_IpGeolocate(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	intelcli := ip_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &ip_intel.IpGeolocateRequest{
		Ip:       "93.231.182.110",
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "digitalelement",
	}
	resp, err := intelcli.Geolocate(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	assert.NotNil(t, resp.Result.Data)
	assert.Equal(t, resp.Result.Data.Country, "Federal Republic Of Germany")
	assert.Equal(t, resp.Result.Data.City, "unna")
	assert.Equal(t, resp.Result.Data.PostalCode, "59425")
}

func Test_Integration_IpGeolocate_DefaultProvider(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	intelcli := ip_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &ip_intel.IpGeolocateRequest{
		Ip: "93.231.182.110",
	}
	resp, err := intelcli.Geolocate(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	assert.NotNil(t, resp.Result.Data)
	assert.Equal(t, resp.Result.Data.Country, "Federal Republic Of Germany")
	assert.Equal(t, resp.Result.Data.City, "unna")
	assert.Equal(t, resp.Result.Data.PostalCode, "59425")
}

func Test_Integration_IpGetDomain(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	intelcli := ip_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &ip_intel.IpDomainRequest{
		Ip:       "24.235.114.61",
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "digitalelement",
	}
	resp, err := intelcli.GetDomain(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	assert.NotNil(t, resp.Result.Data)
	assert.True(t, resp.Result.Data.DomainFound)
	assert.Equal(t, resp.Result.Data.Domain, "rogers.com")
}

func Test_Integration_IpGetDomain_DefaultProvider(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	intelcli := ip_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &ip_intel.IpDomainRequest{
		Ip: "24.235.114.61",
	}
	resp, err := intelcli.GetDomain(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	assert.NotNil(t, resp.Result.Data)
	assert.True(t, resp.Result.Data.DomainFound)
	assert.Equal(t, resp.Result.Data.Domain, "rogers.com")
}

func Test_Integration_IpIsVPN(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	intelcli := ip_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &ip_intel.IpVPNRequest{
		Ip:       "2.56.189.74",
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "digitalelement",
	}
	resp, err := intelcli.IsVPN(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	assert.NotNil(t, resp.Result.Data)
	assert.True(t, resp.Result.Data.IsVPN)
}

func Test_Integration_IpIsVPN_DefaultProvider(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	intelcli := ip_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &ip_intel.IpVPNRequest{
		Ip: "2.56.189.74",
	}
	resp, err := intelcli.IsVPN(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	assert.NotNil(t, resp.Result.Data)
	assert.True(t, resp.Result.Data.IsVPN)
}

func Test_Integration_IpIsProxy(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	intelcli := ip_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &ip_intel.IpProxyRequest{
		Ip:       "34.201.32.172",
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "digitalelement",
	}
	resp, err := intelcli.IsProxy(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	assert.NotNil(t, resp.Result.Data)
	assert.True(t, resp.Result.Data.IsProxy)
}

func Test_Integration_IpIsProxy_DefaultProvider(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	intelcli := ip_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &ip_intel.IpProxyRequest{
		Ip: "34.201.32.172",
	}
	resp, err := intelcli.IsProxy(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	assert.NotNil(t, resp.Result.Data)
	assert.True(t, resp.Result.Data.IsProxy)
}

func Test_Integration_IpReputationCrowdstrikeProvider(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	intelcli := ip_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &ip_intel.IpReputationRequest{
		Ip:       "93.231.182.110",
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

func Test_Integration_IpReputation_CymruProvider(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	intelcli := ip_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &ip_intel.IpReputationRequest{
		Ip:       "93.231.182.110",
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "cymru",
	}
	resp, err := intelcli.Reputation(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	assert.NotNil(t, resp.Result.Data)
}

// Unknown IP
func Test_Integration_IpReputation_2(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	intelcli := ip_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &ip_intel.IpReputationRequest{
		Ip:       "8.8.4.4",
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

func Test_Integration_IpReputation_Error_BadIPFormat_1(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	intelcli := ip_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &ip_intel.IpReputationRequest{
		Ip:       "93.231.182.300",
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "crowdstrike",
	}
	resp, err := intelcli.Reputation(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, resp)
	apiErr := err.(*pangea.APIError)
	assert.Equal(t, len(apiErr.PangeaErrors.Errors), 1)
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Code, "BadFormatIPAddress")
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Detail, "'ip' must be a valid IPv4 or IPv6 address")
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Source, "/ip")
}

func Test_Integration_IpReputation_Error_BadIPFormat_2(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	intelcli := ip_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &ip_intel.IpReputationRequest{
		Ip:       "notanip",
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "crowdstrike",
	}
	resp, err := intelcli.Reputation(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, resp)
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
	intelcli := ip_intel.New(cfg)

	input := &ip_intel.IpReputationRequest{
		Ip:       "93.231.182.110",
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

package ip_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
)

// @summary Geolocate
//
// @description Retrieve information about the location of an IP address.
//
// @operationId ip_intel_post_v1_geolocate
//
// @example
//
//	input := &ip_intel.IpGeolocateRequest{
//		Ip:       "93.231.182.110",
//		Raw:      true,
//		Verbose:  true,
//		Provider: "digitalelement",
//	}
//
//	checkOutput, _, err := ipintel.geolocate(ctx, input)
func (e *ipIntel) Geolocate(ctx context.Context, input *IpGeolocateRequest) (*pangea.PangeaResponse[IpGeolocateResult], error) {
	req, err := e.Client.NewRequest("POST", "v1/geolocate", input)
	if err != nil {
		return nil, err
	}
	out := IpGeolocateResult{}
	resp, err := e.Client.Do(ctx, req, &out)

	if resp == nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[IpGeolocateResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, err
}

// @summary Reputation
//
// @description Retrieve a reputation score for an IP address from a provider,
// including an optional detailed report.
//
// @operationId ip_intel_post_v1_reputation
//
// @example
//
//	input := &ip_intel.IpReputationRequest{
//		Ip:       "93.231.182.110",
//		Raw:      true,
//		Verbose:  true,
//		Provider: "crowdstrike",
//	}
//
//	checkOutput, _, err := ipintel.Reputation(ctx, input)
func (e *ipIntel) Reputation(ctx context.Context, input *IpReputationRequest) (*pangea.PangeaResponse[IpReputationResult], error) {
	req, err := e.Client.NewRequest("POST", "v1/reputation", input)
	if err != nil {
		return nil, err
	}
	out := IpReputationResult{}
	resp, err := e.Client.Do(ctx, req, &out)

	if resp == nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[IpReputationResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, err
}

// @summary Domain
//
// @description Retrieve the domain name associated with an IP address.
//
// @operationId ip_intel_post_v1_domain
//
// @example
//
//	input := &ip_intel.IpDomainRequest{
//		Ip:       "93.231.182.110",
//		Raw:      true,
//		Verbose:  true,
//		Provider: "digitalelement",
//	}
//
//	checkOutput, _, err := ipintel.GetDomain(ctx, input)
func (e *ipIntel) GetDomain(ctx context.Context, input *IpDomainRequest) (*pangea.PangeaResponse[IpDomainResult], error) {
	req, err := e.Client.NewRequest("POST", "v1/domain", input)
	if err != nil {
		return nil, err
	}
	out := IpDomainResult{}
	resp, err := e.Client.Do(ctx, req, &out)

	if resp == nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[IpDomainResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, err
}

// @summary VPN
//
// @description Determine if an IP address is provided by a VPN service.
//
// @operationId ip_intel_post_v1_vpn
//
// @example
//
//	input := &ip_intel.IpVPNRequest{
//		Ip:       "93.231.182.110",
//		Raw:      true,
//		Verbose:  true,
//		Provider: "digitalelement",
//	}
//
//	checkOutput, _, err := ipintel.IsVPN(ctx, input)
func (e *ipIntel) IsVPN(ctx context.Context, input *IpVPNRequest) (*pangea.PangeaResponse[IpVPNResult], error) {
	req, err := e.Client.NewRequest("POST", "v1/vpn", input)
	if err != nil {
		return nil, err
	}
	out := IpVPNResult{}
	resp, err := e.Client.Do(ctx, req, &out)

	if resp == nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[IpVPNResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, err
}

// @summary Proxy
//
// @description Determine if an IP address is provided by a proxy service.
//
// @operationId ip_intel_post_v1_proxy
//
// @example
//
//	input := &ip_intel.IpProxyRequest{
//		Ip:       "93.231.182.110",
//		Raw:      true,
//		Verbose:  true,
//		Provider: "digitalelement",
//	}
//
//	checkOutput, _, err := ipintel.IsProxy(ctx, input)
func (e *ipIntel) IsProxy(ctx context.Context, input *IpProxyRequest) (*pangea.PangeaResponse[IpProxyResult], error) {
	req, err := e.Client.NewRequest("POST", "v1/proxy", input)
	if err != nil {
		return nil, err
	}
	out := IpProxyResult{}
	resp, err := e.Client.Do(ctx, req, &out)

	if resp == nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[IpProxyResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, err
}

type IpGeolocateRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Ip       string `json:"ip"`
	Verbose  *bool  `json:"verbose,omitempty"`
	Raw      *bool  `json:"raw,omitempty"`
	Provider string `json:"provider,omitempty"`
}

type IpReputationRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Ip       string `json:"ip"`
	Verbose  *bool  `json:"verbose,omitempty"`
	Raw      *bool  `json:"raw,omitempty"`
	Provider string `json:"provider,omitempty"`
}

type ReputationData struct {
	Category []string `json:"category"`
	Score    int      `json:"score"`
	Verdict  string   `json:"verdict"`
}

type IpReputationResult struct {
	Data       ReputationData `json:"data"`
	Parameters interface{}    `json:"parameters,omitempty"`
	RawData    interface{}    `json:"raw_data,omitempty"`
}

type IpGeolocateResult struct {
	Data       GeolocateData `json:"data"`
	Parameters interface{}   `json:"parameters,omitempty"`
	RawData    interface{}   `json:"raw_data,omitempty"`
}

type GeolocateData struct {
	Country     string  `json:"country"`
	City        string  `json:"city"`
	Latitude    float32 `json:"latitude"`
	Longitude   float32 `json:"longitude"`
	PostalCode  string  `json:"postal_code"`
	CountryCode string  `json:"country_code"`
}

type IpDomainRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Ip       string `json:"ip"`
	Verbose  *bool  `json:"verbose,omitempty"`
	Raw      *bool  `json:"raw,omitempty"`
	Provider string `json:"provider,omitempty"`
}

type IpDomainResult struct {
	Data       DomainData  `json:"data"`
	Parameters interface{} `json:"parameters,omitempty"`
	RawData    interface{} `json:"raw_data,omitempty"`
}

type DomainData struct {
	DomainFound bool   `json:"domain_found"`
	Domain      string `json:"domain,omitempty"`
}

type IpVPNRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Ip       string `json:"ip"`
	Verbose  *bool  `json:"verbose,omitempty"`
	Raw      *bool  `json:"raw,omitempty"`
	Provider string `json:"provider,omitempty"`
}

type IpVPNResult struct {
	Data       VPNData     `json:"data"`
	Parameters interface{} `json:"parameters,omitempty"`
	RawData    interface{} `json:"raw_data,omitempty"`
}

type VPNData struct {
	IsVPN bool `json:"is_vpn"`
}

type IpProxyRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Ip       string `json:"ip"`
	Verbose  *bool  `json:"verbose,omitempty"`
	Raw      *bool  `json:"raw,omitempty"`
	Provider string `json:"provider,omitempty"`
}

type IpProxyResult struct {
	Data       ProxyData   `json:"data"`
	Parameters interface{} `json:"parameters,omitempty"`
	RawData    interface{} `json:"raw_data,omitempty"`
}

type ProxyData struct {
	IsProxy bool `json:"is_proxy"`
}
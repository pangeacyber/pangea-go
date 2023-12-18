package ip_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/internal/request"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
)

// @summary Geolocate
//
// @description Retrieve geolocation information for an IP address from a provider,
// including an optional detailed report.
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
//	checkOutput, err := ipintel.Geolocate(ctx, input)
func (e *ipIntel) Geolocate(ctx context.Context, input *IpGeolocateRequest) (*pangea.PangeaResponse[IpGeolocateResult], error) {
	return request.DoPost(ctx, e.Client, "v1/geolocate", input, &IpGeolocateResult{})
}

// @summary Geolocate V2
//
// @description Retrieve geolocation information for a list of IP addresses, from a provider,
// including an optional detailed report.
//
// @operationId ip_intel_post_v2_geolocate
//
// @example
//
//	ips := [...]string{"93.231.182.110"}
//
//	input := &ip_intel.IpGeolocateBulkRequest{
//		Ips:      ips,
//		Raw:      true,
//		Verbose:  true,
//		Provider: "digitalelement",
//	}
//
//	checkOutput, err := ipintel.GeolocateBulk(ctx, input)
func (e *ipIntel) GeolocateBulk(ctx context.Context, input *IpGeolocateBulkRequest) (*pangea.PangeaResponse[IpGeolocateBulkResult], error) {
	return request.DoPost(ctx, e.Client, "v2/geolocate", input, &IpGeolocateBulkResult{})
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
//	checkOutput, err := ipintel.Reputation(ctx, input)
func (e *ipIntel) Reputation(ctx context.Context, input *IpReputationRequest) (*pangea.PangeaResponse[IpReputationResult], error) {
	return request.DoPost(ctx, e.Client, "v1/reputation", input, &IpReputationResult{})
}

// @summary Reputation V2
//
// @description Retrieve a reputation scores for a list of IP addresses, from a provider,
// including an optional detailed report.
//
// @operationId ip_intel_post_v2_reputation
//
// @example
//
//	ips := [...]string{"93.231.182.110"}
//
//	input := &ip_intel.IpReputationBulkRequest{
//		Ip:       ips,
//		Raw:      true,
//		Verbose:  true,
//		Provider: "crowdstrike",
//	}
//
//	checkOutput, err := ipintel.ReputationBulk(ctx, input)
func (e *ipIntel) ReputationBulk(ctx context.Context, input *IpReputationBulkRequest) (*pangea.PangeaResponse[IpReputationBulkResult], error) {
	return request.DoPost(ctx, e.Client, "v2/reputation", input, &IpReputationBulkResult{})
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
//	checkOutput, err := ipintel.GetDomain(ctx, input)
func (e *ipIntel) GetDomain(ctx context.Context, input *IpDomainRequest) (*pangea.PangeaResponse[IpDomainResult], error) {
	return request.DoPost(ctx, e.Client, "v1/domain", input, &IpDomainResult{})
}

// @summary Domain V2
//
// @description Retrieve the domain names associated with a list of IP addresses.
//
// @operationId ip_intel_post_v2_domain
//
// @example
//
//	ips := [...]string{"93.231.182.110"}
//
//	input := &ip_intel.IpDomainBulkRequest{
//		Ip:       ips,
//		Raw:      true,
//		Verbose:  true,
//		Provider: "digitalelement",
//	}
//
//	checkOutput, err := ipintel.GetDomainBulk(ctx, input)
func (e *ipIntel) GetDomainBulk(ctx context.Context, input *IpDomainBulkRequest) (*pangea.PangeaResponse[IpDomainBulkResult], error) {
	return request.DoPost(ctx, e.Client, "v2/domain", input, &IpDomainBulkResult{})
}

// @summary VPN
//
// @description Determine if an IP address originates from a VPN.
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
	return request.DoPost(ctx, e.Client, "v1/vpn", input, &IpVPNResult{})
}

// @summary VPN V2
//
// @description Determine which IP addresses originate from a VPN.
//
// @operationId ip_intel_post_v2_vpn
//
// @example
//
//	ips := [...]string{"93.231.182.110"}
//
//	input := &ip_intel.IpVPNBulkRequest{
//		Ip:       ips,
//		Raw:      true,
//		Verbose:  true,
//		Provider: "digitalelement",
//	}
//
//	checkOutput, err := ipintel.IsVPNBulk(ctx, input)
func (e *ipIntel) IsVPNBulk(ctx context.Context, input *IpVPNBulkRequest) (*pangea.PangeaResponse[IpVPNBulkResult], error) {
	return request.DoPost(ctx, e.Client, "v2/vpn", input, &IpVPNBulkResult{})
}

// @summary Proxy
//
// @description Determine if an IP address originates from a proxy.
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
//	checkOutput, err := ipintel.IsProxy(ctx, input)
func (e *ipIntel) IsProxy(ctx context.Context, input *IpProxyRequest) (*pangea.PangeaResponse[IpProxyResult], error) {
	return request.DoPost(ctx, e.Client, "v1/proxy", input, &IpProxyResult{})
}

// @summary Proxy V2
//
// @description Determine which IP addresses originate from a proxy.
//
// @operationId ip_intel_post_v1_proxy
//
// @example
//
//	ips := [...]string{"93.231.182.110"}
//
//	input := &ip_intel.IpProxyBulkRequest{
//		Ip:       ips,
//		Raw:      true,
//		Verbose:  true,
//		Provider: "digitalelement",
//	}
//
//	checkOutput, err := ipintel.IsProxyBulk(ctx, input)
func (e *ipIntel) IsProxyBulk(ctx context.Context, input *IpProxyBulkRequest) (*pangea.PangeaResponse[IpProxyBulkResult], error) {
	return request.DoPost(ctx, e.Client, "v2/proxy", input, &IpProxyBulkResult{})
}

type IpGeolocateRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Ip       string `json:"ip"`
	Verbose  *bool  `json:"verbose,omitempty"`
	Raw      *bool  `json:"raw,omitempty"`
	Provider string `json:"provider,omitempty"`
}

type IpGeolocateBulkRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Ips      []string `json:"ips"`
	Verbose  *bool    `json:"verbose,omitempty"`
	Raw      *bool    `json:"raw,omitempty"`
	Provider string   `json:"provider,omitempty"`
}

type IpReputationRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Ip       string `json:"ip"`
	Verbose  *bool  `json:"verbose,omitempty"`
	Raw      *bool  `json:"raw,omitempty"`
	Provider string `json:"provider,omitempty"`
}

type IpReputationBulkRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Ips      []string `json:"ips"`
	Verbose  *bool    `json:"verbose,omitempty"`
	Raw      *bool    `json:"raw,omitempty"`
	Provider string   `json:"provider,omitempty"`
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

type IpReputationBulkResult struct {
	Data       map[string]ReputationData `json:"data"`
	Parameters interface{}               `json:"parameters,omitempty"`
	RawData    interface{}               `json:"raw_data,omitempty"`
}

type IpGeolocateResult struct {
	Data       GeolocateData `json:"data"`
	Parameters interface{}   `json:"parameters,omitempty"`
	RawData    interface{}   `json:"raw_data,omitempty"`
}

type IpGeolocateBulkResult struct {
	Data       map[string]GeolocateData `json:"data"`
	Parameters interface{}              `json:"parameters,omitempty"`
	RawData    interface{}              `json:"raw_data,omitempty"`
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

type IpDomainBulkRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Ips      []string `json:"ips"`
	Verbose  *bool    `json:"verbose,omitempty"`
	Raw      *bool    `json:"raw,omitempty"`
	Provider string   `json:"provider,omitempty"`
}

type IpDomainResult struct {
	Data       DomainData  `json:"data"`
	Parameters interface{} `json:"parameters,omitempty"`
	RawData    interface{} `json:"raw_data,omitempty"`
}

type IpDomainBulkResult struct {
	Data       map[string]DomainData `json:"data"`
	Parameters interface{}           `json:"parameters,omitempty"`
	RawData    interface{}           `json:"raw_data,omitempty"`
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

type IpVPNBulkRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Ips      []string `json:"ips"`
	Verbose  *bool    `json:"verbose,omitempty"`
	Raw      *bool    `json:"raw,omitempty"`
	Provider string   `json:"provider,omitempty"`
}

type IpVPNResult struct {
	Data       VPNData     `json:"data"`
	Parameters interface{} `json:"parameters,omitempty"`
	RawData    interface{} `json:"raw_data,omitempty"`
}

type IpVPNBulkResult struct {
	Data       map[string]VPNData `json:"data"`
	Parameters interface{}        `json:"parameters,omitempty"`
	RawData    interface{}        `json:"raw_data,omitempty"`
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

type IpProxyBulkRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Ips      []string `json:"ips"`
	Verbose  *bool    `json:"verbose,omitempty"`
	Raw      *bool    `json:"raw,omitempty"`
	Provider string   `json:"provider,omitempty"`
}

type IpProxyResult struct {
	Data       ProxyData   `json:"data"`
	Parameters interface{} `json:"parameters,omitempty"`
	RawData    interface{} `json:"raw_data,omitempty"`
}

type IpProxyBulkResult struct {
	Data       map[string]ProxyData `json:"data"`
	Parameters interface{}          `json:"parameters,omitempty"`
	RawData    interface{}          `json:"raw_data,omitempty"`
}

type ProxyData struct {
	IsProxy bool `json:"is_proxy"`
}

package ip_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
)

// IP address Lookup
//
// Lookup an IP address to retrieve reputation data.
//
// Example:
//
//	 input := &ip_intel.IpLookupInput{
//	     Ip: "93.231.182.110",
//	     Raw: true,
//	     Verbose: true,
//	     Provider: "crowdstrike",
//	 }
//
//		checkOutput, _, err := ipintel.Lookup(ctx, input)
func (e *IpIntel) Lookup(ctx context.Context, input *IpLookupRequest) (*pangea.PangeaResponse[IpLookupResult], error) {
	req, err := e.Client.NewRequest("POST", "v1/lookup", input)
	if err != nil {
		return nil, err
	}
	out := IpLookupResult{}
	resp, err := e.Client.Do(ctx, req, &out)

	if resp == nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[IpLookupResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, err
}

// Geolocate an IP
//
// Retrieve information about the location of an IP address.
//
// Example:
//
//	 input := &ip_intel.IpGeolocateRequest{
//	     Ip: "93.231.182.110",
//	     Raw: true,
//	     Verbose: true,
//	     Provider: "digitalenvoy",
//	 }
//
//	checkOutput, _, err := ipintel.geolocate(ctx, input)
func (e *IpIntel) Geolocate(ctx context.Context, input *IpGeolocateRequest) (*pangea.PangeaResponse[IpGeolocateResult], error) {
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

// Look up domain for an IP
//
// Retrieve the domain name associated with an IP address.
//
// Example:
//
//	 input := &ip_intel.IpDomainRequest{
//	     Ip: "93.231.182.110",
//	     Raw: true,
//	     Verbose: true,
//	     Provider: "digitalenvoy",
//	 }
//
//	checkOutput, _, err := ipintel.GetDomain(ctx, input)
func (e *IpIntel) GetDomain(ctx context.Context, input *IpDomainRequest) (*pangea.PangeaResponse[IpDomainResult], error) {
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

// Check if an IP is a VPN
//
// Determine if an IP address is provided by a VPN service.
//
// Example:
//
//	 input := &ip_intel.IpVPNRequest{
//	     Ip: "93.231.182.110",
//	     Raw: true,
//	     Verbose: true,
//	     Provider: "digitalenvoy",
//	 }
//
//	checkOutput, _, err := ipintel.IsVPN(ctx, input)
func (e *IpIntel) IsVPN(ctx context.Context, input *IpVPNRequest) (*pangea.PangeaResponse[IpVPNResult], error) {
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

// Check for proxied IPs
//
// Determine if an IP address is provided by a proxy service.
//
// Example:
//
//	 input := &ip_intel.IpProxyRequest{
//	     Ip: "93.231.182.110",
//	     Raw: true,
//	     Verbose: true,
//	     Provider: "digitalenvoy",
//	 }
//
//	checkOutput, _, err := ipintel.IsProxy(ctx, input)
func (e *IpIntel) IsProxy(ctx context.Context, input *IpProxyRequest) (*pangea.PangeaResponse[IpProxyResult], error) {
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

type IpLookupRequest struct {
	Ip       string `json:"ip"`
	Verbose  bool   `json:"verbose,omitempty"`
	Raw      bool   `json:"raw,omitempty"`
	Provider string `json:"provider,omitempty"`
}

type LookupData struct {
	Category []string `json:"category"`
	Score    int      `json:"score"`
	Verdict  string   `json:"verdict"`
}

type IpLookupResult struct {
	Data       LookupData  `json:"data"`
	Parameters interface{} `json:"parameters,omitempty"`
	RawData    interface{} `json:"raw_data,omitempty"`
}

type IpGeolocateRequest struct {
	Ip       string `json:"ip"`
	Verbose  bool   `json:"verbose,omitempty"`
	Raw      bool   `json:"raw,omitempty"`
	Provider string `json:"provider,omitempty"`
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
	Ip       string `json:"ip"`
	Verbose  bool   `json:"verbose,omitempty"`
	Raw      bool   `json:"raw,omitempty"`
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
	Ip       string `json:"ip"`
	Verbose  bool   `json:"verbose,omitempty"`
	Raw      bool   `json:"raw,omitempty"`
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
	Ip       string `json:"ip"`
	Verbose  bool   `json:"verbose,omitempty"`
	Raw      bool   `json:"raw,omitempty"`
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

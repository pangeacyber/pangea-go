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

// IP address geolocate
//
// Geolocate an IP address.
//
// Example:
//
//	 input := &ip_intel.IpGeolocateRequest{
//	     Ip: "93.231.182.110",
//	     Raw: true,
//	     Verbose: true,
//	     Provider: "crowdstrike",
//	 }
//
//		checkOutput, _, err := ipintel.geolocate(ctx, input)
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

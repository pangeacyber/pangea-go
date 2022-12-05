package embargo

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
)

// Check IP
//
// Check this IP against known sanction and trade embargo lists.
//
// Example:
//
//	input := &embargo.IPCheckInput{
//		IP: pangea.String("213.24.238.26"),
//	}
//
//	checkResponse, err := embargocli.IPCheck(ctx, input)
func (e *Embargo) IPCheck(ctx context.Context, input *IPCheckInput) (*pangea.PangeaResponse[CheckOutput], error) {
	req, err := e.Client.NewRequest("POST", "v1/ip/check", input)
	if err != nil {
		return nil, err
	}
	out := CheckOutput{}
	resp, err := e.Client.Do(ctx, req, &out)
	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[CheckOutput]{
		Response: *resp,
		Result:   &out,
	}
	return &panresp, nil
}

// ISO Code Check
//
// Check this country against known sanction and trade embargo lists.
//
// Example:
//
//	input := &embargo.ISOCheckInput{
//		ISOCode: pangea.String("CU"),
//	}
//
//	checkResponse, err := embargocli.ISOCheck(ctx, input)
func (e *Embargo) ISOCheck(ctx context.Context, input *ISOCheckInput) (*pangea.PangeaResponse[CheckOutput], error) {
	req, err := e.Client.NewRequest("POST", "v1/iso/check", input)
	if err != nil {
		return nil, err
	}
	out := CheckOutput{}
	resp, err := e.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[CheckOutput]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

type IPCheckInput struct {
	// Check this IP against the enabled embargo lists.
	// Accepts both IPV4 and IPV6 strings.
	IP *string `json:"ip,omitempty"`
}

type ISOCheckInput struct {
	// Check this two character country ISO-code against the enabled embargo lists.
	ISOCode *string `json:"iso_code,omitempty"`
}

type Sanction struct {
	EmbargoedCountryISOCode string                 `json:"embargoed_country_iso_code"`
	IssuingCountry          string                 `json:"issuing_country"`
	ListName                string                 `json:"list_name"`
	EmbargoedCountryName    string                 `json:"embargoed_country_name"`
	Annotations             map[string]interface{} `json:"annotations"`
}

type CheckOutput struct {
	Count     int        `json:"count"`
	Sanctions []Sanction `json:"sanctions"`
}

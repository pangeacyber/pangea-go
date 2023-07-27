package embargo

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
)

// @summary Check IP
//
// @description Check this IP against known sanction and trade embargo lists.
//
// @operationId embargo_post_v1_ip_check
//
// @example
//
//	input := &embargo.IPCheckInput{
//		IP: pangea.String("190.6.64.94"),
//	}
//
//	checkResponse, err := embargocli.IPCheck(ctx, input)
func (e *embargo) IPCheck(ctx context.Context, input *IPCheckRequest) (*pangea.PangeaResponse[CheckResult], error) {
	req, err := e.Client.NewRequest("POST", "v1/ip/check", input)
	if err != nil {
		return nil, err
	}
	out := CheckResult{}
	resp, err := e.Client.Do(ctx, req, &out)
	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[CheckResult]{
		Response: *resp,
		Result:   &out,
	}
	return &panresp, nil
}

// @summary ISO Code Check
//
// @description Check this country against known sanction and trade embargo lists.
//
// @operationId embargo_post_v1_iso_check
//
// @example
//
//	input := &embargo.ISOCheckInput{
//		ISOCode: pangea.String("CU"),
//	}
//
//	checkResponse, err := embargocli.ISOCheck(ctx, input)
func (e *embargo) ISOCheck(ctx context.Context, input *ISOCheckRequest) (*pangea.PangeaResponse[CheckResult], error) {
	req, err := e.Client.NewRequest("POST", "v1/iso/check", input)
	if err != nil {
		return nil, err
	}
	out := CheckResult{}
	resp, err := e.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[CheckResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

type IPCheckRequest struct {
	pangea.BaseRequest

	// Check this IP against the enabled embargo lists.
	// Accepts both IPV4 and IPV6 strings.
	IP string `json:"ip,omitempty"`
}

type ISOCheckRequest struct {
	pangea.BaseRequest

	// Check this two character country ISO-code against the enabled embargo lists.
	ISOCode string `json:"iso_code,omitempty"`
}

type Sanction struct {
	EmbargoedCountryISOCode string                 `json:"embargoed_country_iso_code"`
	IssuingCountry          string                 `json:"issuing_country"`
	ListName                string                 `json:"list_name"`
	EmbargoedCountryName    string                 `json:"embargoed_country_name"`
	Annotations             map[string]interface{} `json:"annotations"`
}

type CheckResult struct {
	Count     int        `json:"count"`
	Sanctions []Sanction `json:"sanctions"`
}

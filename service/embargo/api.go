package embargo

import (
	"context"

	"github.com/pangeacyber/go-pangea/pangea"
)

func (e *Embargo) Check(ctx context.Context, input *CheckInput) (*CheckOutput, *pangea.Response, error) {
	req, err := e.Client.NewRequest("POST", "embargo", "v1/check", input)
	if err != nil {
		return nil, nil, err
	}
	out := CheckOutput{}
	resp, err := e.Client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

type CheckInput struct {
	IP      *string `json:"ip,omitempty"`
	ISOCode *string `json:"iso_code,omitempty"`
}

type Sanction struct {
	EmbargoedCountryISOCode *string                `json:"embargoed_country_iso_code"`
	IssuingCountry          *string                `json:"issuing_country"`
	ListName                *string                `json:"list_name"`
	EmbargoedCountryName    *string                `json:"embargoed_country_name"`
	Annotations             map[string]interface{} `json:"annotations"`
}

type CheckOutput struct {
	Count     *int        `json:"count"`
	Sanctions []*Sanction `json:"sanctions"`
}

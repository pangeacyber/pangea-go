package embargo

import (
	"context"

	"go-pangea/pangea"
)

type Embargo struct {
	*pangea.Client
}

func New(cfg pangea.Config) *Embargo {
	return &Embargo{
		Client: pangea.NewClient(cfg),
	}
}

type CheckInput struct {
	IP      *string `json:"ip,omitempty"`
	ISOCode *string `json:"iso_code"`
}

type SanctionOutput struct {
	ListName                *string `json:"list_name,omitempty"`
	EmbargoedCountryName    *string `json:"embargoed_country_name"`
	EmbargoedCountryISOCode *string `json:"embargoed_country_iso_code"`
	IssuingCountry          *string `json:"issuing_country"`
}

type CheckOutput struct {
	Count     *int              `json:"count"`
	Sanctions []*SanctionOutput `json:"sanctions"`
}

func (e *Embargo) Check(ctx context.Context, input *CheckInput) (*CheckOutput, *pangea.Response, error) {
	if input == nil {
		input = &CheckInput{}
	}
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

package embargo

import (
	"context"

	"github.com/pangeacyber/go-pangea/internal/pangeautil"
	"github.com/pangeacyber/go-pangea/pangea"
)

type CheckInput struct {
	IP      *string `json:"ip,omitempty"`
	ISOCode *string `json:"iso_code,omitempty"`
}

func (c CheckInput) String() string {
	return pangeautil.Stringify(c)
}

type Saction struct {
	ListName                *string `json:"list_name,omitempty"`
	EmbargoedCountryName    *string `json:"embargoed_country_name"`
	EmbargoedCountryISOCode *string `json:"embargoed_country_iso_code"`
	IssuingCountry          *string `json:"issuing_country"`
}

func (s Saction) String() string {
	return pangeautil.Stringify(s)
}

type CheckOutput struct {
	Count     *int       `json:"count"`
	Sanctions []*Saction `json:"sanctions"`
}

func (c CheckOutput) String() string {
	return pangeautil.Stringify(c)
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

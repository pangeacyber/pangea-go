package domain_intel

import (
	"context"

	"github.com/pangeacyber/go-pangea/pangea"
)

// Domain Lookup
//
// Lookup an internet domain to retrieve reputation data.
//
// Example:
//
//	 input := &domain_intel.DomainLookupInput{
//	     	Domain: "737updatesboeing.com",
//	     	Raw: true,
//	     	Verbose: true,
//	     Provider: "domaintools",
//	 }
//
//		checkResponse, err := domainintel.Lookup(ctx, input)
func (e *DomainIntel) Lookup(ctx context.Context, input *DomainLookupInput) (*pangea.PangeaResponse[DomainLookupOutput], error) {
	req, err := e.Client.NewRequest("POST", "v1/lookup", input)
	if err != nil {
		return nil, err
	}
	out := DomainLookupOutput{}
	resp, err := e.Client.Do(ctx, req, &out)

	if resp == nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[DomainLookupOutput]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, err
}

type DomainLookupInput struct {
	Domain   string `json:"domain"`
	Verbose  bool   `json:"verbose,omitempty"`
	Raw      bool   `json:"raw,omitempty"`
	Provider string `json:"provider,omitempty"`
}

type LookupData struct {
	Category []string `json:"category"`
	Score    int      `json:"score"`
	Verdict  string   `json:"verdict"`
}

type DomainLookupOutput struct {
	Data       LookupData  `json:"data"`
	Parameters interface{} `json:"parameters,omitempty"`
	RawData    interface{} `json:"raw_data,omitempty"`
}

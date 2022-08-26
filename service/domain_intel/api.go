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
//  input := &domain_intel.DomainLookupInput{
//      Parameters: DomainLookupParameters {
//      	Domain: "teoghehofuuxo.su",
//      	Raw: true,
//      	Verbose: true,
//      },
//      Provider: "crowdstrike",
//  }
//
//  checkOutput, _, err := domainintel.Lookup(ctx, input)
//
func (e *DomainIntel) Lookup(ctx context.Context, input *DomainLookupInput) (*DomainLookupOutput, *pangea.Response, error) {
	req, err := e.Client.NewRequest("POST", "v1/lookup", input)
	if err != nil {
		return nil, nil, err
	}
	out := DomainLookupOutput{}
	resp, err := e.Client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

type DomainLookupParameters struct {
	Domain  string `json:"domain"`
	Verbose bool   `json:"verbose,omitempty"`
	Raw     bool   `json:"raw,omitempty"`
}

type DomainLookupInput struct {
	Parameters DomainLookupParameters `json:"parameters"`
	Provider   string                 `json:"provider,omitempty"`
}

type LookupData struct {
	Category []string `json:"category"`
	Score    int      `json:"score"`
	Verdict  string   `json:"verdict"`
}

type DomainLookupOutput struct {
	Data      LookupData  `json:"data"`
	Parameter interface{} `json:"parameter,omitempty"`
	RawData   interface{} `json:"raw_data,omitempty"`
}

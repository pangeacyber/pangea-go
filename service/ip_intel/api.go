package ip_intel

import (
	"context"

	"github.com/pangeacyber/go-pangea/pangea"
)

// IP address Lookup
//
// Lookup an IP address to retrieve reputation data.
//
// Example:
//
//  input := &ip_intel.IpLookupInput{
//      Ip: "93.231.182.110",
//      Raw: true,
//      Verbose: true,
//      Provider: "crowdstrike",
//  }
//
//  checkOutput, _, err := ipintel.Lookup(ctx, input)
//
func (e *IpIntel) Lookup(ctx context.Context, input *IpLookupInput) (*IpLookupOutput, *pangea.Response, error) {
	req, err := e.Client.NewRequest("POST", "v1/lookup", input)
	if err != nil {
		return nil, nil, err
	}
	out := IpLookupOutput{}
	resp, err := e.Client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

type IpLookupInput struct {
	Ip      string `json:"ip"`
	Verbose bool   `json:"verbose,omitempty"`
	Raw     bool   `json:"raw,omitempty"`
	Provider   string             `json:"provider,omitempty"`
}

type LookupData struct {
	Category []string `json:"category"`
	Score    int      `json:"score"`
	Verdict  string   `json:"verdict"`
}

type IpLookupOutput struct {
	Data      LookupData  `json:"data"`
	Parameters interface{} `json:"parameters,omitempty"`
	RawData   interface{} `json:"raw_data,omitempty"`
}

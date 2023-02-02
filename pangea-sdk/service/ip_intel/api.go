package ip_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
)

// @summary Look up an IP
//
// @deprecated Use Reputation instead.
//
// @description Retrieve a reputation score for an IP address from a provider,
// including an optional detailed report.
//
// @example
//	input := &ip_intel.IpLookupRequest{
//		Ip: "93.231.182.110",
//		Raw: true,
//		Verbose: true,
//		Provider: "crowdstrike",
//	}
//
//	checkOutput, _, err := ipintel.Lookup(ctx, input)
func (e *IpIntel) Lookup(ctx context.Context, input *IpLookupRequest) (*pangea.PangeaResponse[IpLookupResult], error) {
	req, err := e.Client.NewRequest("POST", "v1/reputation", input)
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

// @summary Look up an IP reputation
//
// @description Retrieve a reputation score for an IP address from a provider,
// including an optional detailed report.
//
// @example
//
//	input := &ip_intel.IpReputationRequest{
//		Ip: "93.231.182.110",
//		Raw: true,
//		Verbose: true,
//		Provider: "crowdstrike",
//	}
//
//	checkOutput, _, err := ipintel.Reputation(ctx, input)
func (e *IpIntel) Reputation(ctx context.Context, input *IpReputationRequest) (*pangea.PangeaResponse[IpReputationResult], error) {
	req, err := e.Client.NewRequest("POST", "v1/reputation", input)
	if err != nil {
		return nil, err
	}
	out := IpReputationResult{}
	resp, err := e.Client.Do(ctx, req, &out)

	if resp == nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[IpReputationResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, err
}

// @deprecated Use IPReputationRequest
type IpLookupRequest struct {
	Ip       string `json:"ip"`
	Verbose  bool   `json:"verbose,omitempty"`
	Raw      bool   `json:"raw,omitempty"`
	Provider string `json:"provider,omitempty"`
}

// @deprecated Use ReputationData
type LookupData struct {
	Category []string `json:"category"`
	Score    int      `json:"score"`
	Verdict  string   `json:"verdict"`
}

// @deprecated Use IpReputationResult
type IpLookupResult struct {
	Data       LookupData  `json:"data"`
	Parameters interface{} `json:"parameters,omitempty"`
	RawData    interface{} `json:"raw_data,omitempty"`
}

type IpReputationRequest struct {
	Ip       string `json:"ip"`
	Verbose  bool   `json:"verbose,omitempty"`
	Raw      bool   `json:"raw,omitempty"`
	Provider string `json:"provider,omitempty"`
}

type ReputationData struct {
	Category []string `json:"category"`
	Score    int      `json:"score"`
	Verdict  string   `json:"verdict"`
}

type IpReputationResult struct {
	Data       ReputationData `json:"data"`
	Parameters interface{}    `json:"parameters,omitempty"`
	RawData    interface{}    `json:"raw_data,omitempty"`
}

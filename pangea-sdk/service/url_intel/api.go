package url_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
)

// @summary Reputation check
//
// @description Retrieve a reputation score for a URL from a provider, including an optional detailed report.
//
// @deprecated Use Reputation instead.
//
// @example
//
//	input := &url_intel.UrlLookupRequest{
//		Url: "http://113.235.101.11:54384",
//		Raw: true,
//		Verbose: true,
//		Provider: "crowdstrike",
//	}
//
//	resp, err := urlintel.Lookup(ctx, input)
func (e *UrlIntel) Lookup(ctx context.Context, input *UrlLookupRequest) (*pangea.PangeaResponse[UrlLookupResult], error) {
	req, err := e.Client.NewRequest("POST", "v1/reputation", input)
	if err != nil {
		return nil, err
	}
	out := UrlLookupResult{}
	resp, err := e.Client.Do(ctx, req, &out)

	if resp == nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[UrlLookupResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, err
}

// @deprecated Use UrlReputationRequest instead.
type UrlLookupRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Url      string `json:"url"`
	Verbose  bool   `json:"verbose,omitempty"`
	Raw      bool   `json:"raw,omitempty"`
	Provider string `json:"provider,omitempty"`
}

// @deprecated Use ReputationData instead.
type LookupData struct {
	Category []string `json:"category"`
	Score    int      `json:"score"`
	Verdict  string   `json:"verdict"`
}

// @deprecated Use UrlReputationResult instead.
type UrlLookupResult struct {
	Data       LookupData  `json:"data"`
	Parameters interface{} `json:"parameters,omitempty"`
	RawData    interface{} `json:"raw_data,omitempty"`
}

// @summary Reputation
//
// @description Retrieve a reputation score for a URL from a provider, including an optional detailed report.
//
// @operationId url_intel_post_v1_reputation
//
// @example
//
//	input := &url_intel.UrlReputationRequest{
//		Url:      "http://113.235.101.11:54384",
//		Raw:      true,
//		Verbose:  true,
//		Provider: "crowdstrike",
//	}
//
//	resp, err := urlintel.Reputation(ctx, input)
func (e *UrlIntel) Reputation(ctx context.Context, input *UrlReputationRequest) (*pangea.PangeaResponse[UrlReputationResult], error) {
	req, err := e.Client.NewRequest("POST", "v1/reputation", input)
	if err != nil {
		return nil, err
	}
	out := UrlReputationResult{}
	resp, err := e.Client.Do(ctx, req, &out)

	if resp == nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[UrlReputationResult]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, err
}

type UrlReputationRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Url      string `json:"url"`
	Verbose  bool   `json:"verbose,omitempty"`
	Raw      bool   `json:"raw,omitempty"`
	Provider string `json:"provider,omitempty"`
}

type ReputationData struct {
	Category []string `json:"category"`
	Score    int      `json:"score"`
	Verdict  string   `json:"verdict"`
}

type UrlReputationResult struct {
	Data       ReputationData `json:"data"`
	Parameters interface{}    `json:"parameters,omitempty"`
	RawData    interface{}    `json:"raw_data,omitempty"`
}

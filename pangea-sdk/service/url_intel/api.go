package url_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/internal/request"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
)

// @summary Reputation check
//
// @description Retrieve a reputation score for a URL, from a provider,
// including an optional detailed report.
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
//	checkOutput, err := urlintel.Reputation(ctx, input)
func (e *urlIntel) Reputation(ctx context.Context, input *UrlReputationRequest) (*pangea.PangeaResponse[UrlReputationResult], error) {
	return request.DoPost(ctx, e.Client, "v1/reputation", input, &UrlReputationResult{})
}

// @summary Reputation check V2
//
// @description Retrieve reputation scores for URLs, from a provider,
// including an optional detailed report.
//
// @operationId url_intel_post_v2_reputation
//
// @example
//
//	urls := [...]string{"http://113.235.101.11:54384"}
//
//	input := &url_intel.UrlReputationBulkRequest{
//		Urls:     urls,
//		Raw:      true,
//		Verbose:  true,
//		Provider: "crowdstrike",
//	}
//
//	checkOutput, err := urlintel.ReputationBulk(ctx, input)
func (e *urlIntel) ReputationBulk(ctx context.Context, input *UrlReputationBulkRequest) (*pangea.PangeaResponse[UrlReputationBulkResult], error) {
	return request.DoPost(ctx, e.Client, "v2/reputation", input, &UrlReputationBulkResult{})
}

type UrlReputationRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Url      string `json:"url"`
	Verbose  *bool  `json:"verbose,omitempty"`
	Raw      *bool  `json:"raw,omitempty"`
	Provider string `json:"provider,omitempty"`
}

type UrlReputationBulkRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	// The URL to be looked up
	Urls []string `json:"urls"`

	// Echo the API parameters in the response.
	Verbose *bool `json:"verbose,omitempty"`

	// Include raw data from this provider.
	Raw *bool `json:"raw,omitempty"`

	// Use reputation data from this provider.
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

type UrlReputationBulkResult struct {
	// High-level normalized results sent
	// by the Pangea service
	Data map[string]ReputationData `json:"data"`

	// The parameters, which were passed in
	// the request, echoed back
	Parameters map[string]any `json:"parameters,omitempty"`

	// The raw data from the provider.
	// Each provider's data will have its own format
	RawData map[string]any `json:"raw_data,omitempty"`
}

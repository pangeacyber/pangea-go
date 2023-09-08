package url_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/internal/request"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
)

// @summary Reputation check
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
//	checkOutput, _, err := urlintel.Reputation(ctx, input)
func (e *urlIntel) Reputation(ctx context.Context, input *UrlReputationRequest) (*pangea.PangeaResponse[UrlReputationResult], error) {
	return request.DoPost(ctx, e.Client, "v1/reputation", input, &UrlReputationResult{})
}

type UrlReputationRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	// The URL to be looked up
	Url string `json:"url,omitempty"`

	// URL list to be looked up.
	UrlList []string `json:"url_list,omitemtpy"`

	// Echo the API parameters in the response.
	Verbose *bool `json:"verbose,omitempty"`

	// Include raw data from this provider.
	Raw *bool `json:"raw,omitempty"`

	// Use reputation data from this provider.
	Provider string `json:"provider,omitempty"`
}

type ReputationData struct {
	// The categories that apply to this
	// indicator as determined by the provider
	Category []string `json:"category"`

	// The score, given by the Pangea service,
	// for the indicator
	Score int `json:"score"`

	// The verdict, given by the Pangea service,
	// for the indicator
	Verdict string `json:"verdict"`
}
type ReputationDataItem struct {
	ReputationData

	Indicator string `json:"indicator"`
}

type UrlReputationResult struct {
	// High-level normalized results sent
	// by the Pangea service
	Data ReputationData `json:"data"`

	// The parameters, which were passed in
	// the request, echoed back
	Parameters map[string]any `json:"parameters,omitempty"`

	// The raw data from the provider.
	// Each provider's data will have its own format
	RawData map[string]any `json:"raw_data,omitempty"`

	// High-level normalized list results sent
	// by the Pangea service
	DataDetails map[string]ReputationDataItem `json:"data_details"`
}

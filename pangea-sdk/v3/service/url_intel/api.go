package url_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/internal/request"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
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

	Url      string `json:"url"`
	Verbose  *bool  `json:"verbose,omitempty"`
	Raw      *bool  `json:"raw,omitempty"`
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

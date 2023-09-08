package domain_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/internal/request"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
)

// @summary Reputation check
//
// @description Lookup an internet domain to retrieve reputation data.
//
// @example
//
//	input := &domain_intel.DomainReputationInput{
//		Domain: "737updatesboeing.com",
//		Raw: true,
//		Verbose: true,
//		Provider: "domaintools",
//	}
//
//	checkResponse, err := domainintel.Reputation(ctx, input)
func (e *domainIntel) Reputation(ctx context.Context, input *DomainReputationRequest) (*pangea.PangeaResponse[DomainReputationResult], error) {
	return request.DoPost(ctx, e.Client, "v1/reputation", input, &DomainReputationResult{})
}

type DomainReputationRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	// The domain to be looked up.
	Domain string `json:"domain,omitempty"`

	// The domain list to be looked up.
	DomainList []string `json:"domain_list,omitempty"`

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

type DomainReputationResult struct {
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
	DataList map[string]ReputationDataItem `json:"data_list"`
}

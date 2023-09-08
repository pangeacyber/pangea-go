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
//	input := &domain_intel.DomainReputationRequest{
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

// @summary WhoIs check
//
// @description Retrieve who is for a domain from a provider, including an optional detailed report.
//
// @example
//
//	input := &domain_intel.DomainWhoIsRequest{
//		Domain: "737updatesboeing.com",
//		Raw: true,
//		Verbose: true,
//		Provider: "domaintools",
//	}
//
//	checkResponse, err := domainintel.Reputation(ctx, input)
func (e *domainIntel) WhoIs(ctx context.Context, input *DomainWhoIsRequest) (*pangea.PangeaResponse[DomainWhoIsResult], error) {
	return request.DoPost(ctx, e.Client, "v1/whois", input, &DomainWhoIsResult{})
}

type DomainReputationRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	// The domain to be looked up.
	Domain string `json:"domain"`

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

type DomainReputationResult struct {
	// High-level normalized results sent
	// by the Pangea service
	Data ReputationData `json:"data"`

	// The parameters, which were passed in
	// the request, echoed back
	Parameters interface{} `json:"parameters,omitempty"`

	// The raw data from the provider.
	// Each provider's data will have its own format
	RawData interface{} `json:"raw_data,omitempty"`
}

type DomainWhoIsRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	// The domain to be looked up.
	Domain string `json:"domain"`

	// Echo the API parameters in the response.
	Verbose *bool `json:"verbose,omitempty"`

	// Include raw data from this provider.
	Raw *bool `json:"raw,omitempty"`

	// Use reputation data from this provider.
	Provider string `json:"provider,omitempty"`
}

type WhoIsData struct {
	DomainName             string   `json:"domain_name"`
	DomainAvailability     string   `json:"domain_availability"`
	CreatedDate            string   `json:"created_date,omitempty"`
	UpdatedDate            string   `json:"updated_date,omitempty"`
	ExpiresDate            string   `json:"expires_date,omitempty"`
	HostNames              []string `json:"host_names,omitempty"`
	IPs                    []string `json:"ips,omitempty"`
	RegistrarName          string   `json:"registrar_name,omitempty"`
	ContactEmail           string   `json:"contact_email,omitempty"`
	EstimatedDomainAge     *int     `json:"estimated_domain_age,omitempty"`
	RegistrantOrganization string   `json:"registrant_organization,omitempty"`
	RegistrantCountry      string   `json:"registrant_country,omitempty"`
}

type DomainWhoIsResult struct {
	Data WhoIsData `json:"data"`

	// The parameters, which were passed in
	// the request, echoed back
	Parameters map[string]any `json:"parameters,omitempty"`

	// The raw data from the provider.
	// Each provider's data will have its own format
	RawData map[string]any `json:"raw_data,omitempty"`
}

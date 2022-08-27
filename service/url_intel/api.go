package url_intel

import (
	"context"

	"github.com/pangeacyber/go-pangea/pangea"
)

// URL Lookup
//
// Lookup a web URL to retrieve reputation data.
//
// Example:
//
//  input := &url_intel.UrlLookupInput{
//      Parameters: UrlLookupParameters {
//      	Url: "http://113.235.101.11:54384",
//      	Raw: true,
//      	Verbose: true,
//      },
//      Provider: "crowdstrike",
//  }
//
//  checkOutput, _, err := urlintel.Lookup(ctx, input)
//
func (e *UrlIntel) Lookup(ctx context.Context, input *UrlLookupInput) (*UrlLookupOutput, *pangea.Response, error) {
	req, err := e.Client.NewRequest("POST", "v1/lookup", input)
	if err != nil {
		return nil, nil, err
	}
	out := UrlLookupOutput{}
	resp, err := e.Client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

type UrlLookupParameters struct {
	Url     string `json:"url"`
	Verbose bool   `json:"verbose,omitempty"`
	Raw     bool   `json:"raw,omitempty"`
}

type UrlLookupInput struct {
	Parameters UrlLookupParameters `json:"parameters"`
	Provider   string              `json:"provider,omitempty"`
}

type LookupData struct {
	Category []string `json:"category"`
	Score    int      `json:"score"`
	Verdict  string   `json:"verdict"`
}

type UrlLookupOutput struct {
	Data      LookupData  `json:"data"`
	Parameter interface{} `json:"parameter,omitempty"`
	RawData   interface{} `json:"raw_data,omitempty"`
}

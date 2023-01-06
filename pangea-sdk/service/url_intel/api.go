package url_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
)

// Look up a URL
//
// Retrieve a reputation score for a URL from a provider, including an optional detailed report.
//
// Example:
//
//	 input := &url_intel.UrlLookupInput{
//	     Url: "http://113.235.101.11:54384",
//	     Raw: true,
//	     Verbose: true,
//	     Provider: "crowdstrike",
//	 }
//
//		checkOutput, _, err := urlintel.Lookup(ctx, input)
func (e *UrlIntel) Lookup(ctx context.Context, input *UrlLookupRequest) (*pangea.PangeaResponse[UrlLookupResult], error) {
	req, err := e.Client.NewRequest("POST", "v1/lookup", input)
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

type UrlLookupRequest struct {
	Url      string `json:"url"`
	Verbose  bool   `json:"verbose,omitempty"`
	Raw      bool   `json:"raw,omitempty"`
	Provider string `json:"provider,omitempty"`
}

type LookupData struct {
	Category []string `json:"category"`
	Score    int      `json:"score"`
	Verdict  string   `json:"verdict"`
}

type UrlLookupResult struct {
	Data       LookupData  `json:"data"`
	Parameters interface{} `json:"parameters,omitempty"`
	RawData    interface{} `json:"raw_data,omitempty"`
}

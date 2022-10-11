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
//	 input := &url_intel.UrlLookupInput{
//	     Url: "http://113.235.101.11:54384",
//	     Raw: true,
//	     Verbose: true,
//	     Provider: "crowdstrike",
//	 }
//
//		checkOutput, _, err := urlintel.Lookup(ctx, input)
func (e *UrlIntel) Lookup(ctx context.Context, input *UrlLookupInput) (*pangea.PangeaResponse[UrlLookupOutput], error) {
	req, err := e.Client.NewRequest("POST", "v1/lookup", input)
	if err != nil {
		return nil, err
	}
	out := UrlLookupOutput{}
	resp, err := e.Client.Do(ctx, req, &out)

	if resp == nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[UrlLookupOutput]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, err
}

type UrlLookupInput struct {
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

type UrlLookupOutput struct {
	Data       LookupData  `json:"data"`
	Parameters interface{} `json:"parameters,omitempty"`
	RawData    interface{} `json:"raw_data,omitempty"`
}

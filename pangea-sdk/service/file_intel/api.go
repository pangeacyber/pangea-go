package file_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
)

type FileLookupInput struct {
	Hash string `json:"hash"`

	// One of "sha256", "sha", "md5".
	HashType string `json:"hash_type"`

	// Echo the API parameters in the response.
	Verbose bool `json:"verbose,omitempty"`

	// Include raw data from this provider.
	Raw bool `json:"raw,omitempty"`

	// Use reputation data from this provider.
	Provider string `json:"provider,omitempty"`
}

type LookupData struct {
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

type FileLookupOutput struct {
	// High-level normalized results sent
	// by the Pangea service
	Data LookupData `json:"data"`

	// The parameters, which were passed in
	// the request, echoed back
	Parameters interface{} `json:"parameters,omitempty"`

	// The raw data from the provider.
	// Each provider's data will have its own format
	RawData interface{} `json:"raw_data,omitempty"`
}

// Look up a file
//
// Lookup a file's hash to retrieve reputation data.
//
// Example:
//
//	input := &file_intel.FileLookupInput{
//	    Hash: "322ccbd42b7e4fd3a9d0167ca2fa9f6483d9691364c431625f1df54270647ca8",
//	    HashType: "sha256",
//	    Raw: true,
//	    Verbose: true,
//	    Provider: "reversinglabs",
//	}
//
//	checkOutput, _, err := fileintel.Lookup(ctx, input)
func (e *FileIntel) Lookup(ctx context.Context, input *FileLookupInput) (*pangea.PangeaResponse[FileLookupOutput], error) {
	req, err := e.Client.NewRequest("POST", "v1/lookup", input)
	if err != nil {
		return nil, err
	}
	out := FileLookupOutput{}
	resp, err := e.Client.Do(ctx, req, &out)

	if resp == nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[FileLookupOutput]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, err
}

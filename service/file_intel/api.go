package file_intel

import (
	"context"

	"github.com/pangeacyber/go-pangea/pangea"
)

// File Lookup
//
// Lookup a file's hash to retrieve reputation data.
//
// Example:
//
//  input := &file_intel.FileLookupInput{
//      Hash: "322ccbd42b7e4fd3a9d0167ca2fa9f6483d9691364c431625f1df54270647ca8",
//      HashType: "sha256",
//      Raw: true,
//      Verbose: true,
//      Provider: "reversinglabs",
//  }
//
//  checkOutput, _, err := fileintel.Lookup(ctx, input)
//

type FileLookupInput struct {
	Hash     string `json:"hash"`
	HashType string `json:"hash_type"`
	Verbose  bool   `json:"verbose,omitempty"`
	Raw      bool   `json:"raw,omitempty"`
	Provider string `json:"provider,omitempty"`
}

type LookupData struct {
	Category []string `json:"category"`
	Score    int      `json:"score"`
	Verdict  string   `json:"verdict"`
}

type FileLookupOutput struct {
	Data       LookupData  `json:"data"`
	Parameters interface{} `json:"parameters,omitempty"`
	RawData    interface{} `json:"raw_data,omitempty"`
}

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

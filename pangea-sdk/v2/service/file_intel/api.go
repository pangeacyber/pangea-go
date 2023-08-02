package file_intel

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/internal/request"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
)

type FileReputationRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Hash string `json:"hash"`

	// One of "sha256", "sha", "md5".
	HashType string `json:"hash_type"`

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

type FileReputationResult struct {
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

// @summary Reputation check
//
// @description Lookup a file's hash to retrieve reputation data.
//
// @operationId file_intel_post_v1_reputation
//
// @example
//
//	input := &file_intel.FileReputationRequest{
//		Hash: "322ccbd42b7e4fd3a9d0167ca2fa9f6483d9691364c431625f1df54270647ca8",
//		HashType: "sha256",
//		Raw: true,
//		Verbose: true,
//		Provider: "reversinglabs",
//	}
//
//	checkOutput, _, err := fileintel.Reputation(ctx, input)
func (e *fileIntel) Reputation(ctx context.Context, input *FileReputationRequest) (*pangea.PangeaResponse[FileReputationResult], error) {
	return request.DoPost(ctx, e.Client, "v1/reputation", input, &FileReputationResult{})
}

// Create a FileReputationRequest from path file
func NewFileReputationRequestFromFilepath(fp string) (*FileReputationRequest, error) {
	f, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}

	return &FileReputationRequest{
		Hash:     hex.EncodeToString(h.Sum(nil)),
		HashType: "sha256",
	}, nil
}

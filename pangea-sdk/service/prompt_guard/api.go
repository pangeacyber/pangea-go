package prompt_guard

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/internal/request"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/pangea"
)

// @summary Guard (Beta)
//
// @description Guard messages.
//
// @operationId prompt_guard_post_v1beta_guard
//
// @example
//
//	input := &prompt_guard.GuardRequest{Messages: []prompt_guard.Message{{Role: "user", Content: "how are you?"}}}
//	response, err := client.Guard(ctx, input)
func (e *promptGuard) Guard(ctx context.Context, input *GuardRequest) (*pangea.PangeaResponse[GuardResult], error) {
	return request.DoPost(ctx, e.Client, "v1beta/guard", input, &GuardResult{})
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GuardRequest struct {
	pangea.BaseRequest

	Messages []Message `json:"messages"`
}

type GuardResult struct {
	Detected   bool   `json:"detected"`
	Type       string `json:"type,omitempty"`
	Detector   string `json:"detector,omitempty"`
	Confidence int    `json:"confidence"`
}

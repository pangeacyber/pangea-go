package prompt_guard

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/internal/request"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
)

// @summary Guard
//
// @description Guard messages.
//
// @operationId prompt_guard_post_v1_guard
//
// @example
//
//	input := &prompt_guard.GuardRequest{Messages: []prompt_guard.Message{{Role: "user", Content: "how are you?"}}}
//	response, err := client.Guard(ctx, input)
func (e *promptGuard) Guard(ctx context.Context, input *GuardRequest) (*pangea.PangeaResponse[GuardResult], error) {
	return request.DoPost(ctx, e.Client, "v1/guard", input, &GuardResult{})
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Classification struct {
	Category   string  `json:"category"`   // Classification category
	Detected   bool    `json:"detected"`   // Classification detection result
	Confidence float32 `json:"confidence"` // Confidence score for the classification
}

type GuardRequest struct {
	pangea.BaseRequest

	Messages  []Message `json:"messages"`            // Prompt content and role array. The content is the text that will be analyzed for redaction.
	Analyzers []string  `json:"analyzers,omitempty"` // Specific analyzers to be used in the call
	Classify  bool      `json:"classify,omitempty"`  // Boolean to enable classification of the content
}

type GuardResult struct {
	Detected        bool             `json:"detected"`           // Boolean response for if the prompt was considered malicious or not
	Type            string           `json:"type,omitempty"`     // Type of analysis, either direct or indirect
	Analyzer        string           `json:"analyzer,omitempty"` // Prompt Analyzers for identifying and rejecting properties of prompts
	Confidence      float32          `json:"confidence"`         // Percent of confidence in the detection result, ranging from 0 to 1
	Info            string           `json:"info,omitempty"`     // Extra information about the detection result
	Classifications []Classification `json:"classifications"`    // List of classification results with labels and confidence scores
}

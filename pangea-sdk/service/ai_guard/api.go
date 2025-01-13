package ai_guard

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/internal/request"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/pangea"
)

// @summary Text guard (Beta)
//
// @description Guard text.
//
// @operationId ai_guard_post_v1beta_text_guard
//
// @example
//
//	input := &ai_guard.TextGuardRequest{Text: "hello world"}
//	response, err := client.GuardText(ctx, input)
func (e *aiGuard) GuardText(ctx context.Context, input *TextGuardRequest) (*pangea.PangeaResponse[TextGuardResult], error) {
	return request.DoPost(ctx, e.Client, "v1beta/text/guard", input, &TextGuardResult{})
}

type AnalyzerResponse struct {
	Analyzer   string  `json:"analyzer"`
	Confidence float64 `json:"confidence"`
}

type PromptInjectionResult struct {
	AnalyzerResponses []AnalyzerResponse `json:"analyzer_responses"`
}

type PiiEntity struct {
	Type     string   `json:"type"`
	Value    string   `json:"value"`
	Redacted bool     `json:"redacted"`
	StartPos *float64 `json:"start_pos,omitempty"`
}

type PiiEntityResult struct {
	Entities []PiiEntity `json:"entities"`
}

type MaliciousEntity struct {
	Type     string                 `json:"type"`
	Value    string                 `json:"value"`
	Redacted *bool                  `json:"redacted,omitempty"`
	StartPos *float64               `json:"start_pos,omitempty"`
	Raw      map[string]interface{} `json:"raw,omitempty"`
}

type MaliciousEntityResult struct {
	Entities []MaliciousEntity `json:"entities"`
}

type TextGuardDetector[T any] struct {
	Detected bool `json:"detected"`
	Data     *T   `json:"data,omitempty"`
}

type TextGuardDetectors struct {
	PromptInjection *TextGuardDetector[PromptInjectionResult] `json:"prompt_injection,omitempty"`
	PiiEntity       *TextGuardDetector[PiiEntityResult]       `json:"pii_entity,omitempty"`
	MaliciousEntity *TextGuardDetector[MaliciousEntityResult] `json:"malicious_entity,omitempty"`
}

type TextGuardRequest struct {
	pangea.BaseRequest

	Text   string `json:"text"`
	Recipe string `json:"recipe,omitempty"`
	Debug  bool   `json:"debug,omitempty"`
}

type TextGuardResult struct {
	Detectors TextGuardDetectors `json:"detectors"`
	Prompt    string             `json:"prompt"`
}

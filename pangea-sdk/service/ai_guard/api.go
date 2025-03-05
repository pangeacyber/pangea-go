package ai_guard

import (
	"context"
	"errors"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/internal/request"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/pangea"
)

// @summary Text guard
//
// @description Guard text.
//
// @operationId ai_guard_post_v1_text_guard
//
// @example
//
//	input := &ai_guard.TextGuardRequest{Text: "hello world"}
//	response, err := client.GuardText(ctx, input)
func (e *aiGuard) GuardText(ctx context.Context, input *TextGuardRequest) (*pangea.PangeaResponse[TextGuardResult], error) {
	if input.Text == "" && input.Messages == nil {
		return nil, errors.New("one of `Text` or `Messages` must be defined")
	}

	return request.DoPost(ctx, e.Client, "v1/text/guard", input, &TextGuardResult{})
}

type AnalyzerResponse struct {
	Analyzer   string  `json:"analyzer"`
	Confidence float64 `json:"confidence"`
}

type PromptInjectionResult struct {
	Action            string             `json:"action"`
	AnalyzerResponses []AnalyzerResponse `json:"analyzer_responses"`
}

type PiiEntity struct {
	Type     string `json:"type"`
	Value    string `json:"value"`
	Action   string `json:"action"`
	StartPos *int   `json:"start_pos,omitempty"`
}

type PiiEntityResult struct {
	Entities []PiiEntity `json:"entities"`
}

type MaliciousEntity struct {
	Type     string                 `json:"type"`
	Value    string                 `json:"value"`
	Action   string                 `json:"action"`
	StartPos *int                   `json:"start_pos,omitempty"`
	Raw      map[string]interface{} `json:"raw,omitempty"`
}

type MaliciousEntityResult struct {
	Entities []MaliciousEntity `json:"entities"`
}

type SecretsEntity struct {
	Type          string `json:"type"`
	Value         string `json:"value"`
	Action        string `json:"action"`
	StartPos      *int   `json:"start_pos,omitempty"`
	RedactedValue string `json:"redacted_value,omitempty"`
}

type SecretsEntityResult struct {
	Entities []SecretsEntity `json:"entities"`
}

type LanguageDetectionResult struct {
	Language string `json:"language"`
	Action   string `json:"action"`
}

type CodeDetectionResult struct {
	Language string `json:"language"`
	Action   string `json:"action"`
}
type TextGuardDetector[T any] struct {
	Detected bool `json:"detected"`
	Data     *T   `json:"data,omitempty"`
}

type TextGuardDetectors struct {
	PromptInjection      *TextGuardDetector[PromptInjectionResult]   `json:"prompt_injection,omitempty"`
	PiiEntity            *TextGuardDetector[PiiEntityResult]         `json:"pii_entity,omitempty"`
	MaliciousEntity      *TextGuardDetector[MaliciousEntityResult]   `json:"malicious_entity,omitempty"`
	SecretsDetection     *TextGuardDetector[SecretsEntityResult]     `json:"secrets_detection,omitempty"`
	ProfanityAndToxicity *TextGuardDetector[any]                     `json:"profanity_and_toxicity,omitempty"`
	CustomEntity         *TextGuardDetector[any]                     `json:"custom_entity,omitempty"`
	LanguageDetection    *TextGuardDetector[LanguageDetectionResult] `json:"language_detection,omitempty"`
	CodeDetection        *TextGuardDetector[CodeDetectionResult]     `json:"code_detection,omitempty"`
}

// LogFields are additional fields to include in activity log
type LogFields struct {
	Citations string `json:"citations,omitempty"`  // Origin or source application of the event
	ExtraInfo string `json:"extra_info,omitempty"` // Stores supplementary details related to the event
	Model     string `json:"model,omitempty"`      // Model used to perform the event
	Source    string `json:"source,omitempty"`     // IP address of user or app or agent
	Tools     string `json:"tools,omitempty"`      // Tools used to perform the event
}

type TextGuardRequest struct {
	pangea.BaseRequest

	Text      string    `json:"text,omitempty"`       // Text to be scanned by AI Guard for PII, sensitive data, malicious content, and other data types defined by the configuration. Supports processing up to 10KB of text.
	Messages  any       `json:"messages,omitempty"`   // Structured messages data to be scanned by AI Guard for PII, sensitive data, malicious content, and other data types defined by the configuration. Supports processing up to 10KB of JSON text.
	Recipe    string    `json:"recipe,omitempty"`     // Recipe key of a configuration of data types and settings defined in the Pangea User Console. It specifies the rules that are to be applied to the text, such as defang malicious URLs.
	Debug     bool      `json:"debug,omitempty"`      // Setting this value to true will provide a detailed analysis of the text data
	LogFields LogFields `json:"log_fields,omitempty"` // Additional fields to include in activity log
}

type TextGuardResult struct {
	Detectors      TextGuardDetectors `json:"detectors"`       // Result of the recipe analyzing and input prompt.
	PromptText     string             `json:"prompt_text"`     // Updated prompt text, if applicable.
	PromptMessages any                `json:"prompt_messages"` // Updated prompt messages, if applicable.
	Blocked        bool               `json:"blocked"`
}

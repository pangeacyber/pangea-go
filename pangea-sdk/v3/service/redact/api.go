package redact

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/internal/request"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
)

// @summary Redact
//
// @description Redacts the content of a single text string.
//
// @operationId redact_post_v1_redact
//
// @example
//
//	input := &redact.TextInput{
//		Text: pangea.String("my phone number is 123-456-7890"),
//	}
//
//	redactResponse, err := redactcli.Redact(ctx, input)
func (r *redact) Redact(ctx context.Context, input *TextRequest) (*pangea.PangeaResponse[TextResult], error) {
	return request.DoPost(ctx, r.Client, "v1/redact", input, &TextResult{})
}

// @summary Redact structured
//
// @description Redacts text within a structured object.
//
// @operationId redact_post_v1_redact_structured
//
// @example
//
//	data := yourCustomDataStruct{
//		Secret: "My social security number is 0303456",
//	}
//
//	input := &redact.StructuredInput{
//		Data: data,
//	}
//
//	redactResponse, err := redactcli.RedactStructured(ctx, input)
func (r *redact) RedactStructured(ctx context.Context, input *StructuredRequest) (*pangea.PangeaResponse[StructuredResult], error) {
	return request.DoPost(ctx, r.Client, "v1/redact_structured", input, &StructuredResult{})
}

// @summary Unredact
//
// @description Decrypt or unredact fpe redactions
//
// @operationId redact_post_v1_unredact
func (r *redact) Unredact(ctx context.Context, input *UnredactRequest) (*pangea.PangeaResponse[UnredactResult], error) {
	return request.DoPost(ctx, r.Client, "v1/unredact", input, &UnredactResult{})
}

type TextRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	// The text to be redacted.
	// Text is a required field.
	Text *string `json:"text"`

	// If the response should include some debug Info.
	Debug *bool `json:"debug,omitempty"`

	// An array of redact rule short names
	Rules []string `json:"rules,omitempty"`

	// An array of redact rulesets short names
	Rulesets []string `json:"rulesets,omitempty"`

	// Setting this value to false will omit the redacted result only returning count
	ReturnResult *bool `json:"return_result,omitempty"`

	// A set of redaction method overrides for any enabled rule. These methods override the config declared methods
	RedactionMethodOverrides *RedactionMethodOverrides `json:"redaction_method_overrides,omitempty"`
}

type TextResult struct {
	// The redacted text.
	RedactedText *string `json:"redacted_text"`

	// Number of redactions present in the response
	Count int `json:"count"`

	Report *DebugReport `json:"report"`
}

type DebugReport struct {
	SummaryCounts     map[string]int     `json:"summary_counts"`
	RecognizerResults []RecognizerResult `json:"recognizer_results"`
}

type RecognizerResult struct {
	// FieldType is always populated on a successful response.
	FieldType string `json:"field_type"`

	// Score is always populated on a successful response.
	Score *float64 `json:"score"`

	// Text is always populated on a successful response.
	Text string `json:"text"`

	// Start is always populated on a successful response.
	Start int `json:"start"`

	// End is always populated on a successful response.
	End int `json:"end"`

	// Redacted is always populated on a successful response.
	Redacted bool `json:"redacted"`

	// DataKey is always populated on a successful response.
	DataKey string `json:"data_key"`
}

type StructuredRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	// Structured data to redact
	// Data is a required field.
	Data map[string]any `json:"data"`

	// JSON path(s) used to identify the specific JSON fields to redact in the structured data.
	// Note: If jsonp parameter is used, the data parameter must be in JSON format.
	JSONP []*string `json:"jsonp,omitempty"`

	// The format of the structured data to redact.
	Format *string `json:"format,omitempty"`

	// Setting this value to true will provide a detailed analysis of the redacted data and the rules that caused redaction.
	Debug *bool `json:"debug,omitempty"`

	// An array of redact rule short names
	Rules []string `json:"rules,omitempty"`

	// An array of redact rulesets short names
	Rulesets []string `json:"rulesets,omitempty"`

	// Setting this value to false will omit the redacted result only returning count
	ReturnResult *bool `json:"return_result,omitempty"`

	// A set of redaction method overrides for any enabled rule. These methods override the config declared methods
	RedactionMethodOverrides *RedactionMethodOverrides `json:"redaction_method_overrides,omitempty"`
}

type StructuredResult struct {
	// RedactedData is always populated on a successful response.
	RedactedData map[string]any `json:"redacted_data"`

	// Number of redactions present in the response
	Count int `json:"count"`

	Report *DebugReport `json:"report"`
}

type UnredactRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	// Data to unredact
	RedactedData any `json:"redacted_data"`

	// FPE context used to decrypt and unredact data (in base64)
	FPEContext string `json:"fpe_context"`
}

type UnredactResult struct {
	Data any `json:"data"`
}

type MaskingType string
type RedactType string
type FPEAlphabet string

const (
	RTmask           RedactType = "mask"
	RTpartialMasking RedactType = "partial_masking"
	RTreplacement    RedactType = "replacement"
	RTdetectOnly     RedactType = "detect_only"
	RThash           RedactType = "hash"
	RTfpe            RedactType = "fpe"
)

const (
	FPEAnumeric          FPEAlphabet = "numeric"
	FPEAaphaNumericLower FPEAlphabet = "alphanumericlower"
	FPEAalphanumeric     FPEAlphabet = "alphanumeric"
)

const (
	MTmask   MaskingType = "mask"
	MTunmask MaskingType = "unmask"
)

type PartialMasking struct {
	MaskingType       *MaskingType `json:"masking_type,omitempty"`
	UnmaskedFromLeft  *int         `json:"unmasked_from_left,omitempty"`
	UnmaskedFromRight *int         `json:"unmasked_from_right,omitempty"`
	MaskedFromLeft    *int         `json:"masked_from_left,omitempty"`
	MaskedFromRight   *int         `json:"masked_from_right,omitempty"`
	CharsToIgnore     []string     `json:"chars_to_ignore,omitempty"`
	MaskingChar       []string     `json:"masking_char,omitempty"`
}

type RedactionMethodOverrides struct {
	RedactionType  RedactType             `json:"redaction_type"`
	Hash           map[string]interface{} `json:"hash,omitempty"`
	FPEAlphabet    *FPEAlphabet           `json:"fpe_alphabet,omitempty"`
	PartialMasking *PartialMasking        `json:"partial_masking,omitempty"`
	RedactionValue *string                `json:"redaction_value,omitempty"`
}

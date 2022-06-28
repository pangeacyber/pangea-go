package redact

import (
	"context"

	"github.com/pangeacyber/go-pangea/pangea"
)

type TextInput struct {
	// The text to be redacted.
	// Text is a required field.
	Text *string `json:"text"`

	// The language of the text.
	// eg: "en"
	Lang *string `json:"lang,omitempty"`

	// If the response should include some debug Info.
	Debug *bool `json:"debug,omitempty"`
}

type TextOutput struct {
	// The redacted text.
	Text *string `json:"text"`

	Report *bool `json:"report"`
}

func (r *Redact) Redact(ctx context.Context, input *TextInput) (*TextOutput, *pangea.Response, error) {
	if input == nil {
		input = &TextInput{}
	}
	req, err := r.Client.NewRequest("POST", "redact", "v1/redact", input)
	if err != nil {
		return nil, nil, err
	}
	out := TextOutput{}
	resp, err := r.Client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}

	return &out, resp, nil
}

type StructuredInput struct {
	Data   *string `json:"data"`
	Format *string `json:"format"`
	Debug  *bool   `json:"debug"`
}

type RecognizerResult struct {
	FieldType *string `json:"field_type"`
	Score     *int    `json:"score"`
	Text      *string `json:"text"`
	Start     *int    `json:"start"`
	End       *int    `json:"end"`
	Redacted  *bool   `json:"redacted"`
	DataKey   *string `json:"data_key"`
}

type StructuredReportOutput struct {
	SummaryCounts     *string             `json:"summary_counts"`
	RecognizerResults []*RecognizerResult `json:"recognizer_results"`
}

type StructuredOutput struct {
	RedactedData *string                 `json:"redacted_data"`
	Report       *StructuredReportOutput `json:"report"`
}

func (r *Redact) RedactStructured(ctx context.Context, input *StructuredInput) (*StructuredOutput, *pangea.Response, error) {
	if input == nil {
		input = &StructuredInput{}
	}
	req, err := r.Client.NewRequest("POST", "redact", "v1/redact_structured", input)
	if err != nil {
		return nil, nil, err
	}
	out := StructuredOutput{}
	resp, err := r.Client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

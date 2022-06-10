package redact

import (
	"context"

	"go-pangea/pangea"
)

type Redact struct {
	Client *pangea.Client
}

type TextInput struct {
	Text  *string `json:"text"`
	Debug *bool   `json:"debug"`
}

type TextOutput struct {
	Text  *string `json:"text"`
	Debug *bool   `json:"debug"`
}

func (r *Redact) Text(ctx context.Context, input *TextInput) (*TextOutput, *pangea.Response, error) {
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

func (r *Redact) Structured(ctx context.Context, input *StructuredInput) (*StructuredOutput, *pangea.Response, error) {
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

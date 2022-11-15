package redact

import (
	"context"
	"encoding/json"

	"github.com/pangeacyber/pangea-go/packages/pangea-sdk/pangea"
)

// Redact
//
// Redacts the content of a single text string.
//
// Example:
//
//	input := &redact.TextInput{
//		Text: pangea.String("my phone number is 123-456-7890"),
//	}
//
//	redactResponse, err := redactcli.Redact(ctx, input)
func (r *Redact) Redact(ctx context.Context, input *TextInput) (*pangea.PangeaResponse[TextOutput], error) {
	req, err := r.Client.NewRequest("POST", "v1/redact", input)
	if err != nil {
		return nil, err
	}
	out := TextOutput{}
	resp, err := r.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}
	panresp := pangea.PangeaResponse[TextOutput]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

// Redact structured
//
// Redacts text within a structured object.
//
// Example:
//
//	type yourCustomDataStruct struct {
//		Secret string `json:"secret"`
//	}
//
//	input := &redact.StructuredInput{
//		JSONP: []*string{
//			pangea.String("$.secret"),
//		},
//	}
//	rawData := yourCustomDataStruct{
//		Secret: "My social security number is 0303456",
//	}
//	input.SetData(rawData)
//
//	redactResponse, err := redactcli.RedactStructured(ctx, input)
func (r *Redact) RedactStructured(ctx context.Context, input *StructuredInput) (*pangea.PangeaResponse[StructuredOutput], error) {
	if input == nil {
		input = &StructuredInput{}
	}
	req, err := r.Client.NewRequest("POST", "v1/redact_structured", input)
	if err != nil {
		return nil, err
	}
	out := StructuredOutput{}
	resp, err := r.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[StructuredOutput]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

type TextInput struct {
	// The text to be redacted.
	// Text is a required field.
	Text *string `json:"text"`

	// If the response should include some debug Info.
	Debug *bool `json:"debug,omitempty"`
}

type TextOutput struct {
	// The redacted text.
	RedactedText *string `json:"redacted_text"`

	Report *DebugReport `json:"report"`
}

type DebugReport struct {
	SummaryCounts     map[string]int      `json:"summary_counts"`
	RecognizerResults []*RecognizerResult `json:"recognizer_results"`
}

type RecognizerResult struct {
	// FieldType is always populated on a successful response.
	FieldType *string `json:"field_type"`

	// Score is always populated on a successful response.
	Score *int `json:"score"`

	// Text is always populated on a successful response.
	Text *string `json:"text"`

	// Start is always populated on a successful response.
	Start *int `json:"start"`

	// End is always populated on a successful response.
	End *int `json:"end"`

	// Redacted is always populated on a successful response.
	Redacted *bool `json:"redacted"` // FieldType is always populated on a successful response.

	DataKey *string `json:"data_key"`
}

type StructuredInput struct {
	// Structured data to redact
	// Data is a required field.
	Data json.RawMessage `json:"data"`

	// JSON path(s) used to identify the specific JSON fields to redact in the structured data.
	// Note: If jsonp parameter is used, the data parameter must be in JSON format.
	JSONP []*string `json:"jsonp,omitempty"`

	// The format of the structured data to redact.
	Format *string `json:"format,omitempty"`

	// Setting this value to true will provide a detailed analysis of the redacted data and the rules that caused redaction.
	Debug *bool `json:"debug,omitempty"`
}

// SetData marshal and sets the JSON encoding of obj into Data.
func (i *StructuredInput) SetData(obj interface{}) error {
	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	i.Data = b
	return nil
}

type StructuredOutput struct {
	// RedactedData is always populated on a successful response.
	RedactedData json.RawMessage `json:"redacted_data"`
	Report       *DebugReport    `json:"report"`
}

// GetRedactedData a parses the JSON-encoded RedactedData and stores the result in the value pointed to by obj.
// If v is nil or not a pointer, Unmarshal returns an InvalidUnmarshalError.
func (i *StructuredOutput) GetRedactedData(obj interface{}) error {
	return json.Unmarshal(i.RedactedData, obj)
}

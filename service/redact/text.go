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

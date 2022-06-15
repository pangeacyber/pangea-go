package redact

import (
	"context"

	"github.com/pangeacyber/go-pangea/pangea"
)

type Client interface {
	RedactText(ctx context.Context, input *TextInput) (*TextOutput, *pangea.Response, error)
	RedactStructured(ctx context.Context, input *StructuredInput) (*StructuredOutput, *pangea.Response, error)
}

type Redact struct {
	*pangea.Client
}

func New(cfg pangea.Config) *Redact {
	return &Redact{
		Client: pangea.NewClient(cfg),
	}
}

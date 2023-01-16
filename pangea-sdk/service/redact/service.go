package redact

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
)

type Client interface {
	RedactText(ctx context.Context, input *TextRequest) (*pangea.PangeaResponse[TextResult], error)
	RedactStructured(ctx context.Context, input *StructuredRequest) (*pangea.PangeaResponse[StructuredResult], error)
}

type Redact struct {
	*pangea.Client
}

func New(cfg *pangea.Config) *Redact {
	cli := &Redact{
		Client: pangea.NewClient("redact", cfg),
	}
	return cli
}

package redact

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
)

type Client interface {
	Redact(ctx context.Context, input *TextRequest) (*pangea.PangeaResponse[TextResult], error)
	RedactStructured(ctx context.Context, input *StructuredRequest) (*pangea.PangeaResponse[StructuredResult], error)
}

type redact struct {
	*pangea.Client
}

func New(cfg *pangea.Config) Client {
	cli := &redact{
		Client: pangea.NewClient("redact", false, cfg),
	}
	return cli
}

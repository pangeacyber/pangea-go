package redact

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
)

type Client interface {
	Redact(ctx context.Context, req *TextRequest) (*pangea.PangeaResponse[TextResult], error)
	RedactStructured(ctx context.Context, req *StructuredRequest) (*pangea.PangeaResponse[StructuredResult], error)

	// Base service methods
	pangea.BaseServicer
}

type redact struct {
	pangea.BaseService
}

func New(cfg *pangea.Config) Client {
	cli := &redact{
		BaseService: pangea.NewBaseService("redact", cfg),
	}
	return cli
}

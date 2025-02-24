package redact

import (
	"context"
	"fmt"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/pangea"
)

type Client interface {
	Redact(ctx context.Context, req *TextRequest) (*pangea.PangeaResponse[TextResult], error)
	RedactStructured(ctx context.Context, req *StructuredRequest) (*pangea.PangeaResponse[StructuredResult], error)
	Unredact(ctx context.Context, input *UnredactRequest) (*pangea.PangeaResponse[UnredactResult], error)

	// Base service methods
	pangea.BaseServicer
}

type redact struct {
	pangea.BaseService
}

func New(cfg *pangea.Config, opts ...Option) Client {
	cli := &redact{
		BaseService: pangea.NewBaseService("redact", cfg),
	}

	for _, opt := range opts {
		err := opt(cli)
		if err != nil {
			fmt.Println("Error applying options to redact service")
		}
	}
	return cli
}

type Option func(*redact) error

func WithConfigID(cid string) Option {
	return func(a *redact) error {
		return pangea.WithConfigID(cid)(&a.BaseService)
	}
}

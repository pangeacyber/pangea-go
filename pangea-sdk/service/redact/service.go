package redact

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
)

type Client interface {
	Redact(ctx context.Context, input *TextRequest) (*pangea.PangeaResponse[TextResult], error)
	RedactStructured(ctx context.Context, input *StructuredRequest) (*pangea.PangeaResponse[StructuredResult], error)
	// Base service methods
	GetPendingRequestID() []string
	PollResultByException(ctx context.Context, e pangea.AcceptedError) (*pangea.PangeaResponse[any], error)
	PollResultByID(ctx context.Context, rid string, v any) (*pangea.PangeaResponse[any], error)
	PollResultRaw(ctx context.Context, requestID string) (*pangea.PangeaResponse[map[string]any], error)
}

type redact struct {
	pangea.BaseService
}

func New(cfg *pangea.Config) Client {
	cli := &redact{
		BaseService: pangea.NewBaseService("redact", false, cfg),
	}
	return cli
}

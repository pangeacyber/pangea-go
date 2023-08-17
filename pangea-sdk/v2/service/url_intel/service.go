package url_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
)

type Client interface {
	Reputation(ctx context.Context, req *UrlReputationRequest) (*pangea.PangeaResponse[UrlReputationResult], error)

	// Base service methods
	GetPendingRequestID() []string
	PollResultByError(ctx context.Context, e pangea.AcceptedError) (*pangea.PangeaResponse[any], error)
	PollResultByID(ctx context.Context, rid string, v any) (*pangea.PangeaResponse[any], error)
	PollResultRaw(ctx context.Context, requestID string) (*pangea.PangeaResponse[map[string]any], error)
}

type urlIntel struct {
	pangea.BaseService
}

func New(cfg *pangea.Config) Client {
	cli := &urlIntel{
		BaseService: pangea.NewBaseService("url-intel", false, cfg),
	}
	return cli
}

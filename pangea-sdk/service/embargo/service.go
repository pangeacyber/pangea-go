package embargo

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
)

type Client interface {
	IPCheck(ctx context.Context, input *IPCheckRequest) (*pangea.PangeaResponse[CheckResult], error)
	ISOCheck(ctx context.Context, input *ISOCheckRequest) (*pangea.PangeaResponse[CheckResult], error)
	// Base service methods
	GetPendingRequestID() []string
	PollResultByError(ctx context.Context, e pangea.AcceptedError) (*pangea.PangeaResponse[any], error)
	PollResultByID(ctx context.Context, rid string, v any) (*pangea.PangeaResponse[any], error)
	PollResultRaw(ctx context.Context, requestID string) (*pangea.PangeaResponse[map[string]any], error)
}

type embargo struct {
	pangea.BaseService
}

func New(cfg *pangea.Config) Client {
	cli := &embargo{
		BaseService: pangea.NewBaseService("embargo", false, cfg),
	}
	return cli
}

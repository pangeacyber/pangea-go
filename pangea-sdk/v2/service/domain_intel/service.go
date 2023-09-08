package domain_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
)

type Client interface {
	Reputation(ctx context.Context, req *DomainReputationRequest) (*pangea.PangeaResponse[DomainReputationResult], error)
	WhoIs(ctx context.Context, req *DomainWhoIsRequest) (*pangea.PangeaResponse[DomainWhoIsResult], error)

	// Base service methods
	GetPendingRequestID() []string
	PollResultByError(ctx context.Context, e pangea.AcceptedError) (*pangea.PangeaResponse[any], error)
	PollResultByID(ctx context.Context, rid string, v any) (*pangea.PangeaResponse[any], error)
	PollResultRaw(ctx context.Context, requestID string) (*pangea.PangeaResponse[map[string]any], error)
}

type domainIntel struct {
	pangea.BaseService
}

func New(cfg *pangea.Config) Client {
	cli := &domainIntel{
		BaseService: pangea.NewBaseService("domain-intel", false, cfg),
	}

	return cli
}

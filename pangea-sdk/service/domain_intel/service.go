package domain_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
)

type Client interface {
	Reputation(ctx context.Context, req *DomainReputationRequest) (*pangea.PangeaResponse[DomainReputationResult], error)
	ReputationBulk(ctx context.Context, req *DomainReputationBulkRequest) (*pangea.PangeaResponse[DomainReputationBulkResult], error)
	WhoIs(ctx context.Context, req *DomainWhoIsRequest) (*pangea.PangeaResponse[DomainWhoIsResult], error)

	// Base service methods
	pangea.BaseServicer
}

type domainIntel struct {
	pangea.BaseService
}

func New(cfg *pangea.Config) Client {
	cli := &domainIntel{
		BaseService: pangea.NewBaseService("domain-intel", cfg),
	}

	return cli
}

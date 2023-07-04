package domain_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
)

type Client interface {
	Reputation(ctx context.Context, input *DomainReputationRequest) (*pangea.PangeaResponse[DomainReputationResult], error)
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

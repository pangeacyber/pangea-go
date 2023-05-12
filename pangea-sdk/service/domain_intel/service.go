package domain_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
)

type Client interface {
	Reputation(ctx context.Context, input *DomainReputationRequest) (*pangea.PangeaResponse[DomainReputationResult], error)
}

type domainIntel struct {
	*pangea.Client
}

func New(cfg *pangea.Config) Client {
	cli := &domainIntel{
		Client: pangea.NewClient("domain-intel", cfg),
	}

	return cli
}

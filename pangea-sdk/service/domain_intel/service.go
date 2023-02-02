package domain_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
)

type Client interface {
	Lookup(ctx context.Context, input *DomainLookupInput) (*pangea.PangeaResponse[DomainLookupOutput], error)
	Reputation(ctx context.Context, input *DomainReputationRequest) (*pangea.PangeaResponse[DomainReputationResult], error)
}

type DomainIntel struct {
	*pangea.Client
}

func New(cfg *pangea.Config) Client {
	cli := &DomainIntel{
		Client: pangea.NewClient("domain-intel", cfg),
	}

	return cli
}

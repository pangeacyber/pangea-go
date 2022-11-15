package domain_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/packages/pangea-sdk/pangea"
)

type Client interface {
	Lookup(ctx context.Context, input *DomainLookupInput) (*pangea.PangeaResponse[DomainLookupOutput], error)
}

type DomainIntel struct {
	*pangea.Client
}

func New(cfg *pangea.Config) *DomainIntel {
	cli := &DomainIntel{
		Client: pangea.NewClient("domain-intel", cfg),
	}

	return cli
}

package domain_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
)

type Client interface {
	Lookup(ctx context.Context, input *DomainLookupRequest) (*pangea.PangeaResponse[DomainLookupResult], error)
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

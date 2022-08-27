package domain_intel

import (
	"context"

	"github.com/pangeacyber/go-pangea/pangea"
)

type Client interface {
	Lookup(ctx context.Context, input *DomainLookupInput) (*DomainLookupOutput, *pangea.Response, error)
}

type DomainIntel struct {
	*pangea.Client
}

func New(cfg *pangea.Config, opts ...Option) (*DomainIntel, error) {
	cli := &DomainIntel{
		Client: pangea.NewClient("domain-intel", cfg),
	}
	for _, opt := range opts {
		err := opt(cli)
		if err != nil {
			return nil, err
		}
	}
	return cli, nil
}

type Option func(*DomainIntel) error

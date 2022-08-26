package ip_intel

import (
	"context"

	"github.com/pangeacyber/go-pangea/pangea"
)

type Client interface {
	Lookup(ctx context.Context, input *IpLookupInput) (*IpLookupOutput, *pangea.Response, error)
}

type IpIntel struct {
	*pangea.Client
}

func New(cfg *pangea.Config, opts ...Option) (*IpIntel, error) {
	cli := &IpIntel{
		Client: pangea.NewClient("ip-intel", cfg),
	}
	for _, opt := range opts {
		err := opt(cli)
		if err != nil {
			return nil, err
		}
	}
	return cli, nil
}

type Option func(*IpIntel) error

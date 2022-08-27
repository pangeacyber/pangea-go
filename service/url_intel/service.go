package url_intel

import (
	"context"

	"github.com/pangeacyber/go-pangea/pangea"
)

type Client interface {
	Lookup(ctx context.Context, input *UrlLookupInput) (*UrlLookupOutput, *pangea.Response, error)
}

type UrlIntel struct {
	*pangea.Client
}

func New(cfg *pangea.Config, opts ...Option) (*UrlIntel, error) {
	cli := &UrlIntel{
		Client: pangea.NewClient("url-intel", cfg),
	}
	for _, opt := range opts {
		err := opt(cli)
		if err != nil {
			return nil, err
		}
	}
	return cli, nil
}

type Option func(*UrlIntel) error

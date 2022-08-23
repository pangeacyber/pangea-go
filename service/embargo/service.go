package embargo

import (
	"context"

	"github.com/pangeacyber/go-pangea/pangea"
)

type Client interface {
	IPCheck(ctx context.Context, input *IPCheckInput) (*pangea.PangeaResponse[CheckOutput], error)
	ISOCheck(ctx context.Context, input *ISOCheckInput) (*pangea.PangeaResponse[CheckOutput], error)
}

type Embargo struct {
	*pangea.Client
}

func New(cfg *pangea.Config, opts ...Option) (*Embargo, error) {
	cli := &Embargo{
		Client: pangea.NewClient("embargo", cfg),
	}
	for _, opt := range opts {
		err := opt(cli)
		if err != nil {
			return nil, err
		}
	}
	return cli, nil
}

type Option func(*Embargo) error

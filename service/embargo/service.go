package embargo

import (
	"context"

	"github.com/pangeacyber/go-pangea/pangea"
)

type Client interface {
	Check(ctx context.Context, input *CheckInput) (*CheckOutput, *pangea.Response, error)
}

type Embargo struct {
	*pangea.Client
}

func New(cfg *pangea.Config, opts ...Option) (*Embargo, error) {
	cli := &Embargo{
		Client: pangea.NewClient(cfg),
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

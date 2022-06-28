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

func New(cfg *pangea.Config, optionalCfg ...*pangea.Config) *Embargo {
	return &Embargo{
		Client: pangea.NewClient(cfg),
	}
}

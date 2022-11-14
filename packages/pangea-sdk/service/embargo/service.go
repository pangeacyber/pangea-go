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

func New(cfg *pangea.Config) *Embargo {
	cli := &Embargo{
		Client: pangea.NewClient("embargo", cfg),
	}
	return cli
}

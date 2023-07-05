package embargo

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
)

type Client interface {
	IPCheck(ctx context.Context, input *IPCheckInput) (*pangea.PangeaResponse[CheckOutput], error)
	ISOCheck(ctx context.Context, input *ISOCheckInput) (*pangea.PangeaResponse[CheckOutput], error)
}

type Embargo struct {
	pangea.BaseService
}

func New(cfg *pangea.Config) *Embargo {
	cli := &Embargo{
		BaseService: pangea.NewBaseService("embargo", false, cfg),
	}
	return cli
}

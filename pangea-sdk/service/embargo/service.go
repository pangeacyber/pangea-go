package embargo

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
)

type Client interface {
	IPCheck(ctx context.Context, input *IPCheckRequest) (*pangea.PangeaResponse[CheckResult], error)
	ISOCheck(ctx context.Context, input *ISOCheckRequest) (*pangea.PangeaResponse[CheckResult], error)
}

type embargo struct {
	*pangea.Client
}

func New(cfg *pangea.Config) Client {
	cli := &embargo{
		Client: pangea.NewClient("embargo", false, cfg),
	}
	return cli
}

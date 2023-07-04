package embargo

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
)

type Client interface {
	IPCheck(ctx context.Context, input *IPCheckRequest) (*pangea.PangeaResponse[CheckResult], error)
	ISOCheck(ctx context.Context, input *ISOCheckRequest) (*pangea.PangeaResponse[CheckResult], error)
}

type embargo struct {
	pangea.BaseService
}

func New(cfg *pangea.Config) Client {
	cli := &embargo{
		BaseService: pangea.NewBaseService("embargo", false, cfg),
	}
	return cli
}

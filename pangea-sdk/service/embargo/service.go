package embargo

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/pangea"
)

type Client interface {
	IPCheck(ctx context.Context, req *IPCheckRequest) (*pangea.PangeaResponse[CheckResult], error)
	ISOCheck(ctx context.Context, req *ISOCheckRequest) (*pangea.PangeaResponse[CheckResult], error)

	// Base service methods
	pangea.BaseServicer
}

type embargo struct {
	pangea.BaseService
}

func New(cfg *pangea.Config) Client {
	cli := &embargo{
		BaseService: pangea.NewBaseService("embargo", cfg),
	}
	return cli
}

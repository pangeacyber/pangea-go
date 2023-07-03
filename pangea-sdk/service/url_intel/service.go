package url_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
)

type Client interface {
	Lookup(ctx context.Context, input *UrlLookupRequest) (*pangea.PangeaResponse[UrlLookupResult], error)
	Reputation(ctx context.Context, input *UrlReputationRequest) (*pangea.PangeaResponse[UrlReputationResult], error)
}

type UrlIntel struct {
	pangea.BaseService
}

func New(cfg *pangea.Config) Client {
	cli := &UrlIntel{
		BaseService: pangea.NewBaseService("url-intel", false, cfg),
	}
	return cli
}

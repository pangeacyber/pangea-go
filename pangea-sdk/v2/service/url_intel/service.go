package url_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
)

type Client interface {
	Reputation(ctx context.Context, req *UrlReputationRequest) (*pangea.PangeaResponse[UrlReputationResult], error)

	// Base service methods
	pangea.BaseServicer
}

type urlIntel struct {
	pangea.BaseService
}

func New(cfg *pangea.Config) Client {
	cli := &urlIntel{
		BaseService: pangea.NewBaseService("url-intel", cfg),
	}
	return cli
}

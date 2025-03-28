package url_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
)

type Client interface {
	Reputation(ctx context.Context, req *UrlReputationRequest) (*pangea.PangeaResponse[UrlReputationResult], error)
	ReputationBulk(ctx context.Context, req *UrlReputationBulkRequest) (*pangea.PangeaResponse[UrlReputationBulkResult], error)

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

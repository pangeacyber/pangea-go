package url_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
)

type Client interface {
	Reputation(ctx context.Context, input *UrlReputationRequest) (*pangea.PangeaResponse[UrlReputationResult], error)
}

type urlIntel struct {
	*pangea.Client
}

func New(cfg *pangea.Config) Client {
	cli := &urlIntel{
		Client: pangea.NewClient("url-intel", false, cfg),
	}
	return cli
}

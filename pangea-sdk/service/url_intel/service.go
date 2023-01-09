package url_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
)

type Client interface {
	Lookup(ctx context.Context, input *UrlLookupRequest) (*pangea.PangeaResponse[UrlLookupResult], error)
}

type UrlIntel struct {
	*pangea.Client
}

func New(cfg *pangea.Config) *UrlIntel {
	cli := &UrlIntel{
		Client: pangea.NewClient("url-intel", cfg),
	}
	return cli
}

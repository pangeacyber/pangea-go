package ip_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
)

type Client interface {
	Lookup(ctx context.Context, input *IpLookupRequest) (*pangea.PangeaResponse[IpLookupResult], error)
}

type IpIntel struct {
	*pangea.Client
}

func New(cfg *pangea.Config) *IpIntel {
	cli := &IpIntel{
		Client: pangea.NewClient("ip-intel", cfg),
	}
	return cli
}

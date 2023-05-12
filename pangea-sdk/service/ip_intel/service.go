package ip_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
)

type Client interface {
	Reputation(ctx context.Context, input *IpReputationRequest) (*pangea.PangeaResponse[IpReputationResult], error)
	Geolocate(ctx context.Context, input *IpGeolocateRequest) (*pangea.PangeaResponse[IpGeolocateResult], error)
	GetDomain(ctx context.Context, input *IpDomainRequest) (*pangea.PangeaResponse[IpDomainResult], error)
	IsVPN(ctx context.Context, input *IpVPNRequest) (*pangea.PangeaResponse[IpVPNResult], error)
	IsProxy(ctx context.Context, input *IpProxyRequest) (*pangea.PangeaResponse[IpProxyResult], error)
}

type ipIntel struct {
	*pangea.Client
}

func New(cfg *pangea.Config) Client {
	cli := &ipIntel{
		Client: pangea.NewClient("ip-intel", cfg),
	}
	return cli
}

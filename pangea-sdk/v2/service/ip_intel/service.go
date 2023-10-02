package ip_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
)

type Client interface {
	Reputation(ctx context.Context, req *IpReputationRequest) (*pangea.PangeaResponse[IpReputationResult], error)
	ReputationBulk(ctx context.Context, req *IpReputationBulkRequest) (*pangea.PangeaResponse[IpReputationBulkResult], error)
	Geolocate(ctx context.Context, req *IpGeolocateRequest) (*pangea.PangeaResponse[IpGeolocateResult], error)
	GetDomain(ctx context.Context, req *IpDomainRequest) (*pangea.PangeaResponse[IpDomainResult], error)
	IsVPN(ctx context.Context, req *IpVPNRequest) (*pangea.PangeaResponse[IpVPNResult], error)
	IsProxy(ctx context.Context, req *IpProxyRequest) (*pangea.PangeaResponse[IpProxyResult], error)

	// Base service methods
	pangea.BaseServicer
}

type ipIntel struct {
	pangea.BaseService
}

func New(cfg *pangea.Config) Client {
	cli := &ipIntel{
		BaseService: pangea.NewBaseService("ip-intel", cfg),
	}
	return cli
}

package ip_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
)

type Client interface {
	Reputation(ctx context.Context, req *IpReputationRequest) (*pangea.PangeaResponse[IpReputationResult], error)
	Geolocate(ctx context.Context, req *IpGeolocateRequest) (*pangea.PangeaResponse[IpGeolocateResult], error)
	GetDomain(ctx context.Context, req *IpDomainRequest) (*pangea.PangeaResponse[IpDomainResult], error)
	IsVPN(ctx context.Context, req *IpVPNRequest) (*pangea.PangeaResponse[IpVPNResult], error)
	IsProxy(ctx context.Context, req *IpProxyRequest) (*pangea.PangeaResponse[IpProxyResult], error)

	ReputationBulk(ctx context.Context, req *IpReputationBulkRequest) (*pangea.PangeaResponse[IpReputationBulkResult], error)
	GeolocateBulk(ctx context.Context, req *IpGeolocateBulkRequest) (*pangea.PangeaResponse[IpGeolocateBulkResult], error)
	GetDomainBulk(ctx context.Context, req *IpDomainBulkRequest) (*pangea.PangeaResponse[IpDomainBulkResult], error)
	IsVPNBulk(ctx context.Context, req *IpVPNBulkRequest) (*pangea.PangeaResponse[IpVPNBulkResult], error)
	IsProxyBulk(ctx context.Context, req *IpProxyBulkRequest) (*pangea.PangeaResponse[IpProxyBulkResult], error)

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

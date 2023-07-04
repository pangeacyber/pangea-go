package ip_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
)

type Client interface {
	Reputation(ctx context.Context, input *IpReputationRequest) (*pangea.PangeaResponse[IpReputationResult], error)
	Geolocate(ctx context.Context, input *IpGeolocateRequest) (*pangea.PangeaResponse[IpGeolocateResult], error)
	GetDomain(ctx context.Context, input *IpDomainRequest) (*pangea.PangeaResponse[IpDomainResult], error)
	IsVPN(ctx context.Context, input *IpVPNRequest) (*pangea.PangeaResponse[IpVPNResult], error)
	IsProxy(ctx context.Context, input *IpProxyRequest) (*pangea.PangeaResponse[IpProxyResult], error)

	// Base service methods
	GetPendingRequestID() []string
	PollResultByException(ctx context.Context, e pangea.AcceptedError) (*pangea.PangeaResponse[any], error)
	PollResultByID(ctx context.Context, rid string, v any) (*pangea.PangeaResponse[any], error)
	PollResultRaw(ctx context.Context, requestID string) (*pangea.PangeaResponse[map[string]any], error)
}

type ipIntel struct {
	pangea.BaseService
}

func New(cfg *pangea.Config) Client {
	cli := &ipIntel{
		BaseService: pangea.NewBaseService("ip-intel", false, cfg),
	}
	return cli
}

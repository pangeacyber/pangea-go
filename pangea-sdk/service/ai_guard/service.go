package ai_guard

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
)

// AI Guard API client.
type Client interface {
	GuardText(ctx context.Context, input *TextGuardRequest) (*pangea.PangeaResponse[TextGuardResult], error)
	GetServiceConfig(ctx context.Context, body GetServiceConfigParams) (*pangea.PangeaResponse[ServiceConfig], error)
	CreateServiceConfig(ctx context.Context, body CreateServiceConfigParams) (*pangea.PangeaResponse[ServiceConfig], error)
	UpdateServiceConfig(ctx context.Context, body UpdateServiceConfigParams) (*pangea.PangeaResponse[ServiceConfig], error)
	DeleteServiceConfig(ctx context.Context, body DeleteServiceConfigParams) (*pangea.PangeaResponse[ServiceConfig], error)
	ListServiceConfigs(ctx context.Context, body ListServiceConfigsParams) (*pangea.PangeaResponse[ListServiceConfigsResult], error)

	// Base service methods.
	pangea.BaseServicer
}

type aiGuard struct {
	pangea.BaseService
}

func New(cfg *pangea.Config) Client {
	cli := &aiGuard{
		BaseService: pangea.NewBaseService("ai-guard", cfg),
	}

	return cli
}

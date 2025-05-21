package prompt_guard

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
)

// Prompt Guard API client.
type Client interface {
	Guard(ctx context.Context, input *GuardRequest) (*pangea.PangeaResponse[GuardResult], error)
	GetServiceConfig(ctx context.Context, body GetServiceConfigParams) (*pangea.PangeaResponse[struct{}], error)
	CreateServiceConfig(ctx context.Context, body CreateServiceConfigParams) (*pangea.PangeaResponse[struct{}], error)
	UpdateServiceConfig(ctx context.Context, body UpdateServiceConfigParams) (*pangea.PangeaResponse[struct{}], error)
	DeleteServiceConfig(ctx context.Context, body DeleteServiceConfigParams) (*pangea.PangeaResponse[struct{}], error)
	ListServiceConfigs(ctx context.Context, body ListServiceConfigsParams) (*pangea.PangeaResponse[ListServiceConfigsResult], error)

	// Base service methods.
	pangea.BaseServicer
}

type promptGuard struct {
	pangea.BaseService
}

func New(cfg *pangea.Config) Client {
	return &promptGuard{BaseService: pangea.NewBaseService("prompt-guard", cfg)}
}

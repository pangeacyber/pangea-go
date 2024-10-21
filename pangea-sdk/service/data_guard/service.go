package data_guard

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/pangea"
)

// Data Guard API client.
type Client interface {
	GuardText(ctx context.Context, input *TextGuardRequest) (*pangea.PangeaResponse[TextGuardResult], error)
	GuardFile(ctx context.Context, input *FileGuardRequest) (*pangea.PangeaResponse[struct{}], error)

	// Base service methods.
	pangea.BaseServicer
}

type dataGuard struct {
	pangea.BaseService
}

func New(cfg *pangea.Config) Client {
	cli := &dataGuard{
		BaseService: pangea.NewBaseService("data-guard", cfg),
	}

	return cli
}

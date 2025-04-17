package ai_guard

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
)

// AI Guard API client.
type Client interface {
	GuardText(ctx context.Context, input *TextGuardRequest) (*pangea.PangeaResponse[TextGuardResult], error)

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

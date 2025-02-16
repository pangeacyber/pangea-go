package prompt_guard

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/pangea"
)

// Prompt Guard API client.
type Client interface {
	Guard(ctx context.Context, input *GuardRequest) (*pangea.PangeaResponse[GuardResult], error)

	// Base service methods.
	pangea.BaseServicer
}

type promptGuard struct {
	pangea.BaseService
}

func New(cfg *pangea.Config) Client {
	return &promptGuard{BaseService: pangea.NewBaseService("prompt-guard", cfg)}
}

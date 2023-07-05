package user_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
)

type Client interface {
	UserBreached(ctx context.Context, input *UserBreachedRequest) (*pangea.PangeaResponse[UserBreachedResult], error)
	PasswordBreached(ctx context.Context, input *UserPasswordBreachedRequest) (*pangea.PangeaResponse[UserPasswordBreachedResult], error)
	PollResult(ctx context.Context, e pangea.AcceptedError) (*pangea.PangeaResponse[any], error)
}

type userIntel struct {
	pangea.BaseService
}

func New(cfg *pangea.Config) Client {
	cli := &userIntel{
		BaseService: pangea.NewBaseService("user-intel", false, cfg),
	}
	return cli
}

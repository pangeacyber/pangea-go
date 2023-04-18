package user_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
)

type Client interface {
	UserBreached(ctx context.Context, input *UserBreachedRequest) (*pangea.PangeaResponse[UserBreachedResult], error)
	PasswordBreached(ctx context.Context, input *UserPasswordBreachedRequest) (*pangea.PangeaResponse[UserPasswordBreachedResult], error)
}

type userIntel struct {
	*pangea.Client
}

func New(cfg *pangea.Config) Client {
	cli := &userIntel{
		Client: pangea.NewClient("user-intel", cfg),
	}
	return cli
}

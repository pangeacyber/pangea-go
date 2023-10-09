package user_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
)

type Client interface {
	UserBreached(ctx context.Context, req *UserBreachedRequest) (*pangea.PangeaResponse[UserBreachedResult], error)
	PasswordBreached(ctx context.Context, req *UserPasswordBreachedRequest) (*pangea.PangeaResponse[UserPasswordBreachedResult], error)

	// Base service methods
	pangea.BaseServicer
}

type userIntel struct {
	pangea.BaseService
}

func New(cfg *pangea.Config) Client {
	cli := &userIntel{
		BaseService: pangea.NewBaseService("user-intel", cfg),
	}
	return cli
}

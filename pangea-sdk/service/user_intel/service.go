package user_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
)

type Client interface {
	UserBreached(ctx context.Context, req *UserBreachedRequest) (*pangea.PangeaResponse[UserBreachedResult], error)
	UserBreachedBulk(ctx context.Context, req *UserBreachedBulkRequest) (*pangea.PangeaResponse[UserBreachedBulkResult], error)
	PasswordBreached(ctx context.Context, req *UserPasswordBreachedRequest) (*pangea.PangeaResponse[UserPasswordBreachedResult], error)
	PasswordBreachedBulk(ctx context.Context, req *UserPasswordBreachedBulkRequest) (*pangea.PangeaResponse[UserPasswordBreachedBulkResult], error)

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

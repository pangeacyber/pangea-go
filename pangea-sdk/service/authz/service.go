package authz

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
)

type Client interface {
	TupleCreate(ctx context.Context, input *TupleCreateRequest) (*pangea.PangeaResponse[TupleCreateResult], error)
	TupleList(ctx context.Context, input *TupleListRequest) (*pangea.PangeaResponse[TupleListResult], error)
	TupleDelete(ctx context.Context, input *TupleDeleteRequest) (*pangea.PangeaResponse[TupleDeleteResult], error)
	Check(ctx context.Context, input *CheckRequest) (*pangea.PangeaResponse[CheckResult], error)
	ListResources(ctx context.Context, input *ListResourcesRequest) (*pangea.PangeaResponse[ListResourcesResult], error)
	ListSubjects(ctx context.Context, input *ListSubjectsRequest) (*pangea.PangeaResponse[ListSubjectsResult], error)

	// Base service methods
	pangea.BaseServicer
}

type authz struct {
	pangea.BaseService
}

func New(cfg *pangea.Config) Client {
	cli := &authz{
		BaseService: pangea.NewBaseService("authz", cfg),
	}
	return cli
}

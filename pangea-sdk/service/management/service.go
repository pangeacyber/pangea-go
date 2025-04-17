package management

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
)

// Management API client.
type Client interface {
	GetOrg(ctx context.Context, orgId string) (*pangea.PangeaResponse[Organization], error)
	UpdateOrg(ctx context.Context, input *UpdateOrgRequest) (*pangea.PangeaResponse[Organization], error)
	GetProject(ctx context.Context, projectId string) (*pangea.PangeaResponse[Project], error)
	ListProjects(ctx context.Context, input *ListProjectsRequest) (*pangea.PangeaResponse[ListProjectsResult], error)
	CreateProject(ctx context.Context, input *CreateProjectRequest) (*pangea.PangeaResponse[Project], error)
	UpdateProject(ctx context.Context, input *UpdateProjectRequest) (*pangea.PangeaResponse[Project], error)
	DeleteProject(ctx context.Context, projectId string) (*pangea.PangeaResponse[struct{}], error)
	CreateClient(ctx context.Context, input *CreatePlatformClientRequest) (*AccessClientCreateInfo, error)
	ListClients(ctx context.Context, input *ListPlatformClientsRequest) (*ListPlatformClientsResult, error)
	GetClient(ctx context.Context, clientId string) (*AccessClientInfo, error)
	UpdateClient(ctx context.Context, clientId string, input *UpdatePlatformClientRequest) (*AccessClientCreateInfo, error)
	DeleteClient(ctx context.Context, clientId string) error
	CreateClientSecret(ctx context.Context, clientId string, input *CreateClientSecretRequest) (*AccessClientSecretInfo, error)
	ListClientSecretMetadata(ctx context.Context, clientId string, input *ListClientSecretMetadataRequest) (*AccessClientSecretInfoListResult, error)
	RevokeClientSecret(ctx context.Context, clientId string, clientSecretId string) error
	UpdateClientSecret(ctx context.Context, clientId string, clientSecretId string, input *UpdateClientSecretRequest) (*AccessClientSecretInfo, error)
	ListClientRoles(ctx context.Context, clientId string, input *ListClientRolesRequest) (*AccessRolesListResult, error)
	GrantClientAccess(ctx context.Context, clientId string, input *GrantClientAccessRequest) error
	RevokeClientAccess(ctx context.Context, clientId string, input *RevokeClientAccessRequest) error

	// Base service methods.
	pangea.BaseServicer
}

type management struct {
	pangea.BaseService

	authorization pangea.BaseService
	console       pangea.BaseService
}

func New(cfg *pangea.Config) Client {
	return &management{
		authorization: pangea.NewBaseService("authorization.access", cfg),
		console:       pangea.NewBaseService("api.console", cfg),
	}
}

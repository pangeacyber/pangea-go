package management

import (
	"context"
	"fmt"
	"net/url"

	"github.com/google/go-querystring/query"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/internal/request"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
)

type Organization struct {
	Id         string  `json:"id"`
	Name       string  `json:"name"`
	Owner      string  `json:"owner"`
	OwnerEmail *string `json:"owner_email,omitempty"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  *string `json:"updated_at,omitempty"`
	Csp        string  `json:"csp"`
}

type Project struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	Org       string  `json:"org"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt *string `json:"updated_at,omitempty"`
	// The geographical region for the project.
	Geo string `json:"geo"`
	// The region for the project.
	Region string `json:"region"`
}

type GetOrgRequest struct {
	pangea.BaseRequest

	// An Organization Pangea ID
	Id string `json:"id" validate:"regexp=^poi_[a-z2-7]{32}$"`
}

type UpdateOrgRequest struct {
	pangea.BaseRequest

	// An Organization Pangea ID
	Id   string `json:"id" validate:"regexp=^poi_[a-z2-7]{32}$"`
	Name string `json:"name"`
}

type GetProjectRequest struct {
	pangea.BaseRequest

	// A Project Pangea ID
	Id string `json:"id" validate:"regexp=^ppi_[a-z2-7]{32}$"`
}

type ListProjectsFilter struct {
	Search *string `json:"search,omitempty"`
	Geo    *string `json:"geo,omitempty"`
	Region *string `json:"region,omitempty"`
}

type ListProjectsRequest struct {
	pangea.BaseRequest

	// An Organization Pangea ID
	OrgId  string              `json:"org_id" validate:"regexp=^poi_[a-z2-7]{32}$"`
	Filter *ListProjectsFilter `json:"filter,omitempty"`
	Offset *int32              `json:"offset,omitempty"`
	Limit  *int32              `json:"limit,omitempty"`
}

type ListProjectsResult struct {
	// A list of projects
	Results []Project `json:"results"`
	Count   int32     `json:"count"`
	Offset  *int32    `json:"offset,omitempty"`
}

type CreateProjectRequest struct {
	pangea.BaseRequest

	// An Organization Pangea ID
	OrgId string `json:"org_id" validate:"regexp=^poi_[a-z2-7]{32}$"`
	Name  string `json:"name"`
	// The geographical region for the project.
	Geo string `json:"geo"`
	// The region for the project.
	Region string `json:"region"`
}

type UpdateProjectRequest struct {
	pangea.BaseRequest

	// A Project Pangea ID
	Id   string `json:"id" validate:"regexp=^ppi_[a-z2-7]{32}$"`
	Name string `json:"name"`
}

type DeleteProjectRequest struct {
	pangea.BaseRequest

	// A Project Pangea ID
	Id string `json:"id" validate:"regexp=^ppi_[a-z2-7]{32}$"`
}

// AccessClientTokenAuth The authentication method for the token endpoint.
type AccessClientTokenAuth string

const (
	AccessClientTokenAuthBasic AccessClientTokenAuth = "client_secret_basic"
	AccessClientTokenAuthPost  AccessClientTokenAuth = "client_secret_post"
)

// AccessRole Service token information
type AccessRole struct {
	// The specific role being assigned to a client. Examples include 'manager' for service configurations or 'admin' for projects.
	Role string `json:"role"`
	// The role resource type. Examples include 'organization', 'project', or 'service_{snake_case(service)}_config'.
	Type    string  `json:"type"`
	Id      string  `json:"id"`
	Service *string `json:"service,omitempty"`
	// An ID for a service config
	ServiceConfigId *string `json:"service_config_id,omitempty" validate:"regexp=^pci_[a-z2-7]{32}$"`
}

type CreatePlatformClientRequest struct {
	pangea.BaseRequest

	ClientName string `json:"client_name"`
	// A list of space separated scopes. Examples include \"scope\": \"pangea:service:ai-guard:read pangea:service:redact:read\" for granting AI Guard & Redact API access. The actual service configurations the client has access to for those services is dictated through roles.
	Scope                   string                 `json:"scope" validate:"regexp=^[a-zA-Z0-9:*_\\/_-]+(?:\\\\s+[a-zA-Z0-9:*_\\/_-]+)*$"`
	TokenEndpointAuthMethod *AccessClientTokenAuth `json:"token_endpoint_auth_method,omitempty"`
	// A list of allowed redirect URIs for the client.
	RedirectUris []string `json:"redirect_uris,omitempty"`
	// A list of OAuth grant types that the client can use.
	GrantTypes []string `json:"grant_types,omitempty"`
	// A list of OAuth response types that the client can use.
	ResponseTypes []*string `json:"response_types,omitempty"`
	// A positive time duration in seconds or null
	ClientSecretExpiresIn *int32 `json:"client_secret_expires_in,omitempty"`
	// A positive time duration in seconds or null
	ClientTokenExpiresIn    *int32  `json:"client_token_expires_in,omitempty"`
	ClientSecretName        *string `json:"client_secret_name,omitempty"`
	ClientSecretDescription *string `json:"client_secret_description,omitempty"`
	// A list of roles. Roles are required to grant object access to clients, while client scopes dictate which API routes the clients may access. An example role: { \"type\": \"service_ai_guard_config\", \"id\": \"pci_xxx\", \"role\": \"manager\" }.
	Roles []AccessRole `json:"roles,omitempty"`
}

// AccessClientInfo API Client information
type AccessClientInfo struct {
	// An ID for a service account
	ClientId string `json:"client_id" validate:"regexp=^psa_[a-z2-7]{32}$"`
	// A time in ISO-8601 format
	CreatedAt string `json:"created_at"`
	// A time in ISO-8601 format
	UpdatedAt  string `json:"updated_at"`
	ClientName string `json:"client_name"`
	// A list of space separated scopes. Examples include \"scope\": \"pangea:service:ai-guard:read pangea:service:redact:read\" for granting AI Guard & Redact API access. The actual service configurations the client has access to for those services is dictated through roles.
	Scope                   *string                `json:"scope,omitempty" validate:"regexp=^[a-zA-Z0-9:*_\\/_-]+(?:\\\\s+[a-zA-Z0-9:*_\\/_-]+)*$"`
	TokenEndpointAuthMethod *AccessClientTokenAuth `json:"token_endpoint_auth_method,omitempty"`
	// A list of allowed redirect URIs for the client.
	RedirectUris []string `json:"redirect_uris,omitempty"`
	// A list of OAuth grant types that the client can use.
	GrantTypes []string `json:"grant_types,omitempty"`
	// A list of OAuth response types that the client can use.
	ResponseTypes []*string `json:"response_types,omitempty"`
	// A positive time duration in seconds or null
	ClientTokenExpiresIn *int32  `json:"client_token_expires_in,omitempty"`
	OwnerId              *string `json:"owner_id,omitempty"`
	OwnerUsername        *string `json:"owner_username,omitempty"`
	CreatorId            *string `json:"creator_id,omitempty"`
	ClientClass          *string `json:"client_class,omitempty"`
}

// AccessClientCreateInfo API Client information with initial secret
type AccessClientCreateInfo struct {
	// An ID for a service account
	ClientId string `json:"client_id" validate:"regexp=^psa_[a-z2-7]{32}$"`
	// A time in ISO-8601 format
	CreatedAt string `json:"created_at"`
	// A time in ISO-8601 format
	UpdatedAt  string `json:"updated_at"`
	ClientName string `json:"client_name"`
	// A list of space separated scopes. Examples include \"scope\": \"pangea:service:ai-guard:read pangea:service:redact:read\" for granting AI Guard & Redact API access. The actual service configurations the client has access to for those services is dictated through roles.
	Scope                   *string                `json:"scope,omitempty" validate:"regexp=^[a-zA-Z0-9:*_\\/_-]+(?:\\\\s+[a-zA-Z0-9:*_\\/_-]+)*$"`
	TokenEndpointAuthMethod *AccessClientTokenAuth `json:"token_endpoint_auth_method,omitempty"`
	// A list of allowed redirect URIs for the client.
	RedirectUris []string `json:"redirect_uris,omitempty"`
	// A list of OAuth grant types that the client can use.
	GrantTypes []string `json:"grant_types,omitempty"`
	// A list of OAuth response types that the client can use.
	ResponseTypes []string `json:"response_types,omitempty"`
	// A positive time duration in seconds or null
	ClientTokenExpiresIn *int32  `json:"client_token_expires_in,omitempty"`
	OwnerId              string  `json:"owner_id,omitempty"`
	OwnerUsername        *string `json:"owner_username,omitempty"`
	CreatorId            string  `json:"creator_id,omitempty"`
	ClientClass          *string `json:"client_class,omitempty"`
	ClientSecret         string  `json:"client_secret"`
	// A time in ISO-8601 format
	ClientSecretExpiresAt   string  `json:"client_secret_expires_at"`
	ClientSecretName        *string `json:"client_secret_name,omitempty"`
	ClientSecretDescription *string `json:"client_secret_description,omitempty"`
}

type ListPlatformClientsRequest struct {
	CreatedAt          *string  `url:"created_at,omitempty"`
	CreatedAtGt        *string  `url:"created_at__gt,omitempty"`
	CreatedAtGte       *string  `url:"created_at__gte,omitempty"`
	CreatedAtLt        *string  `url:"created_at__lt,omitempty"`
	CreatedAtLte       *string  `url:"created_at__lte,omitempty"`
	ClientId           *string  `url:"client_id,omitempty"`
	ClientIdContains   []string `url:"client_id__contains,omitempty"`
	ClientIdIn         []string `url:"client_id__in,omitempty"`
	ClientName         *string  `url:"client_name,omitempty"`
	ClientNameContains []string `url:"client_name__contains,omitempty"`
	ClientNameIn       []string `url:"client_name__in,omitempty"`
	Scopes             []string `url:"scopes,omitempty"`
	UpdatedAt          *string  `url:"updated_at,omitempty"`
	UpdatedAtGt        *string  `url:"updated_at__gt,omitempty"`
	UpdatedAtGte       *string  `url:"updated_at__gte,omitempty"`
	UpdatedAtLt        *string  `url:"updated_at__lt,omitempty"`
	UpdatedAtLte       *string  `url:"updated_at__lte,omitempty"`
	Last               *string  `url:"last,omitempty"`
	Order              *string  `url:"order,omitempty"`
	OrderBy            *string  `url:"order_by,omitempty"`
	Size               *int32   `url:"size,omitempty"`
}

type ListPlatformClientsResult struct {
	Clients []AccessClientInfo `json:"clients"`
	Count   int32              `json:"count"`
	Last    *string            `json:"last,omitempty"`
}

// URLQuery serializes [ListPlatformClientsRequest]'s query parameters as
// `url.Values`.
func (r ListPlatformClientsRequest) URLQuery() (values url.Values, err error) {
	values, err = query.Values(r)
	if err != nil {
		return nil, err
	}
	return values, nil
}

type UpdatePlatformClientRequest struct {
	pangea.BaseRequest

	Scope                   *string                `json:"scope,omitempty"`
	ClientName              *string                `json:"client_name,omitempty"`
	TokenEndpointAuthMethod *AccessClientTokenAuth `json:"token_endpoint_auth_method,omitempty"`
	RedirectUris            []string               `json:"redirect_uris,omitempty"`
	GrantTypes              []string               `json:"grant_types,omitempty"`
	ResponseTypes           []string               `json:"response_types,omitempty"`
}

type CreateClientSecretRequest struct {
	pangea.BaseRequest

	// A positive time duration in seconds
	ClientSecretExpiresIn   *int32  `json:"client_secret_expires_in,omitempty"`
	ClientSecretName        *string `json:"client_secret_name,omitempty"`
	ClientSecretDescription *string `json:"client_secret_description,omitempty"`
}

type UpdateClientSecretRequest struct {
	pangea.BaseRequest

	// A positive time duration in seconds
	ClientSecretExpiresIn   *int32  `json:"client_secret_expires_in,omitempty"`
	ClientSecretName        *string `json:"client_secret_name,omitempty"`
	ClientSecretDescription *string `json:"client_secret_description,omitempty"`
}

type AccessClientSecretInfo struct {
	// An ID for a service account
	ClientId string `json:"client_id" validate:"regexp=^psa_[a-z2-7]{32}$"`
	// An ID for an API Client secret
	ClientSecretId string `json:"client_secret_id" validate:"regexp=^pce_[a-z2-7]{32}$"`
	ClientSecret   string `json:"client_secret"`
	// A time in ISO-8601 format
	ClientSecretExpiresAt   string  `json:"client_secret_expires_at"`
	ClientSecretName        *string `json:"client_secret_name,omitempty"`
	ClientSecretDescription *string `json:"client_secret_description,omitempty"`
}

type ListClientSecretMetadataRequest struct {
	pangea.BaseRequest

	CreatedAt                *string  `url:"created_at,omitempty"` // Only records where created_at equals this value.
	CreatedAtGt              *string  `url:"created_at__gt,omitempty"`
	CreatedAtGte             *string  `url:"created_at__gte,omitempty"`
	CreatedAtLt              *string  `url:"created_at__lt,omitempty"`
	CreatedAtLte             *string  `url:"created_at__lte,omitempty"`
	ClientSecretName         *string  `url:"client_secret_name,omitempty"`
	ClientSecretNameContains []string `url:"client_secret_name__contains,omitempty"`
	ClientSecretNameIn       []string `url:"client_secret_name__in,omitempty"`
	Last                     *string  `url:"last,omitempty"`     // Reflected value from a previous response to obtain the next page of results.
	Order                    *string  `url:"order,omitempty"`    // Order results asc(ending) or desc(ending).
	OrderBy                  *string  `url:"order_by,omitempty"` // Which field to order results by.
	Size                     *int32   `url:"size,omitempty"`     // Maximum results to include in the response.
}

// URLQuery serializes [ListClientSecretMetadataRequest]'s query parameters as
// `url.Values`.
func (r ListClientSecretMetadataRequest) URLQuery() (values url.Values, err error) {
	values, err = query.Values(r)
	if err != nil {
		return nil, err
	}
	return values, nil
}

type AccessClientSecretMetadata struct {
	SourceIp    *string `json:"source_ip,omitempty"`
	UserAgent   *string `json:"user_agent,omitempty"`
	Creator     *string `json:"creator,omitempty"`
	CreatorId   *string `json:"creator_id,omitempty"`
	CreatorType *string `json:"creator_type,omitempty"`
}

// AccessClientSecretInfoWithMetadata Service account information
type AccessClientSecretInfoWithMetadata struct {
	// An ID for a service account
	ClientId *string `json:"client_id,omitempty" validate:"regexp=^psa_[a-z2-7]{32}$"`
	// An ID for an API Client secret
	ClientSecretId *string `json:"client_secret_id,omitempty" validate:"regexp=^pce_[a-z2-7]{32}$"`
	// A time in ISO-8601 format
	ClientSecretExpiresAt   *string `json:"client_secret_expires_at,omitempty"`
	ClientSecretName        *string `json:"client_secret_name,omitempty"`
	ClientSecretDescription *string `json:"client_secret_description,omitempty"`
	// A time in ISO-8601 format
	CreatedAt string `json:"created_at"`
	// A time in ISO-8601 format
	UpdatedAt            string                      `json:"updated_at"`
	ClientSecretMetadata *AccessClientSecretMetadata `json:"client_secret_metadata,omitempty"`
}

type AccessClientSecretInfoListResult struct {
	ClientSecrets []AccessClientSecretInfoWithMetadata `json:"client-secrets"`
	Count         int32                                `json:"count"`
	Last          *string                              `json:"last,omitempty"`
}

type ListClientRolesRequest struct {
	pangea.BaseRequest

	ResourceType *string `url:"resource_type,omitempty"`
	ResourceId   *string `url:"resource_id,omitempty"`
	Role         *string `url:"role,omitempty"`
}

// URLQuery serializes [ListClientRolesRequest]'s query parameters as
// `url.Values`.
func (r ListClientRolesRequest) URLQuery() (values url.Values, err error) {
	values, err = query.Values(r)
	if err != nil {
		return nil, err
	}
	return values, nil
}

type AccessRolesListResult struct {
	Roles []AccessRole `json:"roles"`
	Count int32        `json:"count"`
	Last  *string      `json:"last,omitempty"`
}

type GrantClientAccessRequest struct {
	pangea.BaseRequest

	// A list of roles. Roles are required to grant object access to clients, while client scopes dictate which API routes the clients may access. An example role: { \"type\": \"service_ai_guard_config\", \"id\": \"pci_xxx\", \"role\": \"manager\" }.
	Roles []AccessRole `json:"roles"`
	// A list of space separated scopes. Examples include \"scope\": \"pangea:service:ai-guard:read pangea:service:redact:read\" for granting AI Guard & Redact API access. The actual service configurations the client has access to for those services is dictated through roles.
	Scope string `json:"scope" validate:"regexp=^[a-zA-Z0-9:*_\\/_-]+(?:\\\\s+[a-zA-Z0-9:*_\\/_-]+)*$"`
}

type RevokeClientAccessRequest struct {
	pangea.BaseRequest

	// A list of roles. Roles are required to grant object access to clients, while client scopes dictate which API routes the clients may access. An example role: { \"type\": \"service_ai_guard_config\", \"id\": \"pci_xxx\", \"role\": \"manager\" }.
	Roles []AccessRole `json:"roles"`
	// A list of space separated scopes. Examples include \"scope\": \"pangea:service:ai-guard:read pangea:service:redact:read\" for granting AI Guard & Redact API access. The actual service configurations the client has access to for those services is dictated through roles.
	Scope string `json:"scope" validate:"regexp=^[a-zA-Z0-9:*_\\/_-]+(?:\\\\s+[a-zA-Z0-9:*_\\/_-]+)*$"`
}

// @summary Retrieve an organization
//
// @description Retrieve an organization
//
// @operationId api.console_post_v1beta_platform_org_get
func (e *management) GetOrg(ctx context.Context, orgId string) (*pangea.PangeaResponse[Organization], error) {
	return request.DoPost(ctx, e.console.Client, "v1beta/platform/org/get", &GetOrgRequest{Id: orgId}, &Organization{})
}

// @summary Update an organization
//
// @description Update an organization
//
// @operationId api.console_post_v1beta_platform_org_update
func (e *management) UpdateOrg(ctx context.Context, input *UpdateOrgRequest) (*pangea.PangeaResponse[Organization], error) {
	return request.DoPost(ctx, e.console.Client, "v1beta/platform/org/update", input, &Organization{})
}

// @summary Retrieve a project
//
// @description Retrieve a project
//
// @operationId api.console_post_v1beta_platform_project_get
func (e *management) GetProject(ctx context.Context, projectId string) (*pangea.PangeaResponse[Project], error) {
	return request.DoPost(ctx, e.console.Client, "v1beta/platform/project/get", &GetProjectRequest{Id: projectId}, &Project{})
}

// @summary List projects
//
// @description List projects
//
// @operationId api.console_post_v1beta_platform_project_list
func (e *management) ListProjects(ctx context.Context, input *ListProjectsRequest) (*pangea.PangeaResponse[ListProjectsResult], error) {
	return request.DoPost(ctx, e.console.Client, "v1beta/platform/project/list", input, &ListProjectsResult{})
}

// @summary Create a project
//
// @description Create a project
//
// @operationId api.console_post_v1beta_platform_project_create
func (e *management) CreateProject(ctx context.Context, input *CreateProjectRequest) (*pangea.PangeaResponse[Project], error) {
	return request.DoPost(ctx, e.console.Client, "v1beta/platform/project/create", input, &Project{})
}

// @summary Update a project
//
// @description Update a project
//
// @operationId api.console_post_v1beta_platform_project_update
func (e *management) UpdateProject(ctx context.Context, input *UpdateProjectRequest) (*pangea.PangeaResponse[Project], error) {
	return request.DoPost(ctx, e.console.Client, "v1beta/platform/project/update", input, &Project{})
}

// @summary Delete a project
//
// @description Delete a project
//
// @operationId api.console_post_v1beta_platform_project_delete
func (e *management) DeleteProject(ctx context.Context, projectId string) (*pangea.PangeaResponse[struct{}], error) {
	var result struct{}
	return request.DoPost(ctx, e.console.Client, "v1beta/platform/project/delete", &DeleteProjectRequest{Id: projectId}, &result)
}

// @summary Create platform client
//
// @description Create platform client
//
// @operationId createPlatformClient
func (e *management) CreateClient(ctx context.Context, input *CreatePlatformClientRequest) (*AccessClientCreateInfo, error) {
	return request.DoPostNonPangeaResponse(ctx, e.authorization.Client, "v1beta/oauth/clients/register", input, &AccessClientCreateInfo{})
}

// @summary List platform clients
//
// @description List platform clients
//
// @operationId listPlatformClients
func (e *management) ListClients(ctx context.Context, input *ListPlatformClientsRequest) (*ListPlatformClientsResult, error) {
	var result ListPlatformClientsResult
	err := request.Get(ctx, e.authorization.Client, "v1beta/oauth/clients", input, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// @summary Get a platform client
//
// @description Get a platform client
//
// @operationId getPlatformClient
func (e *management) GetClient(ctx context.Context, clientId string) (*AccessClientInfo, error) {
	var result AccessClientInfo
	err := request.Get(ctx, e.authorization.Client, fmt.Sprintf("v1beta/oauth/clients/%s", clientId), nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// @summary Update platform client
//
// @description Update platform client
//
// @operationId updatePlatformClient
func (e *management) UpdateClient(ctx context.Context, clientId string, input *UpdatePlatformClientRequest) (*AccessClientCreateInfo, error) {
	return request.DoPostNonPangeaResponse(
		ctx,
		e.authorization.Client,
		fmt.Sprintf("v1beta/oauth/clients/%s", clientId),
		input,
		&AccessClientCreateInfo{},
	)
}

// @summary Delete a platform client
//
// @description Delete a platform client
//
// @operationId deletePlatformClient
func (e *management) DeleteClient(ctx context.Context, clientId string) error {
	return request.Delete(ctx, e.authorization.Client, fmt.Sprintf("v1beta/oauth/clients/%s", clientId), nil, nil)
}

// @summary Create client secret
//
// @description Create client secret
//
// @operationId createClientSecret
func (e *management) CreateClientSecret(ctx context.Context, clientId string, input *CreateClientSecretRequest) (*AccessClientSecretInfo, error) {
	return request.DoPostNonPangeaResponse(
		ctx,
		e.authorization.Client,
		fmt.Sprintf("v1beta/oauth/clients/%s/secrets", clientId),
		input,
		&AccessClientSecretInfo{},
	)
}

// @summary List client secret metadata
//
// @description List client secret metadata
//
// @operationId listClientSecretMetadata
func (e *management) ListClientSecretMetadata(
	ctx context.Context,
	clientId string,
	input *ListClientSecretMetadataRequest,
) (*AccessClientSecretInfoListResult, error) {
	var result AccessClientSecretInfoListResult
	err := request.Get(
		ctx,
		e.authorization.Client,
		fmt.Sprintf("v1beta/oauth/clients/%s/secrets/metadata", clientId),
		input,
		&result,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// @summary Revoke client secret
//
// @description Revoke client secret
//
// @operationId revokeClientSecret
func (e *management) RevokeClientSecret(ctx context.Context, clientId string, clientSecretId string) error {
	return request.Delete(
		ctx,
		e.authorization.Client,
		fmt.Sprintf("v1beta/oauth/clients/%s/secrets/%s", clientId, clientSecretId),
		nil,
		nil,
	)
}

// @summary Update client secret
//
// @description Update client secret
//
// @operationId updateClientSecret
func (e *management) UpdateClientSecret(
	ctx context.Context,
	clientId string,
	clientSecretId string,
	input *UpdateClientSecretRequest,
) (*AccessClientSecretInfo, error) {
	return request.DoPostNonPangeaResponse(
		ctx,
		e.authorization.Client,
		fmt.Sprintf("v1beta/oauth/clients/%s/secrets/%s", clientId, clientSecretId),
		input,
		&AccessClientSecretInfo{},
	)
}

// @summary List client roles
//
// @description List client roles
//
// @operationId listClientRoles
func (e *management) ListClientRoles(ctx context.Context, clientId string, input *ListClientRolesRequest) (*AccessRolesListResult, error) {
	var result AccessRolesListResult
	err := request.Get(
		ctx,
		e.authorization.Client,
		fmt.Sprintf("v1beta/oauth/clients/%s/roles", clientId),
		input,
		&result,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// @summary Grant client access
//
// @description Grant client access
//
// @operationId grantClientRoles
func (e *management) GrantClientAccess(
	ctx context.Context,
	clientId string,
	input *GrantClientAccessRequest,
) error {
	var result struct{}
	_, err := request.DoPostNonPangeaResponse(
		ctx,
		e.authorization.Client,
		fmt.Sprintf("v1beta/oauth/clients/%s/grant", clientId),
		input,
		&result,
	)
	return err
}

// @summary Revoke client access
//
// @description Revoke client access
//
// @operationId revokeClientRoles
func (e *management) RevokeClientAccess(
	ctx context.Context,
	clientId string,
	input *RevokeClientAccessRequest,
) error {
	var result struct{}
	_, err := request.DoPostNonPangeaResponse(
		ctx,
		e.authorization.Client,
		fmt.Sprintf("v1beta/oauth/clients/%s/revoke", clientId),
		input,
		&result,
	)
	return err
}

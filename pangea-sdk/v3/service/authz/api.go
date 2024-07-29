package authz

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/internal/request"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
)

type Resource struct {
	Type string `json:"type"`
	ID   string `json:"id,omitempty"`
}

type Subject struct {
	Type   string `json:"type"`
	ID     string `json:"id,omitempty"`
	Action string `json:"action,omitempty"`
}

type Tuple struct {
	Resource Resource `json:"resource"`
	Relation string   `json:"relation"`
	Subject  Subject  `json:"subject"`
}

type TupleCreateRequest struct {
	pangea.BaseRequest

	Tuples []Tuple `json:"tuples"`
}

type TupleCreateResult struct {
}

// @summary Write tuples.
//
// @description Write tuples. The request will fail if tuples do not validate against the schema defined resource types.
//
// @operationId authz_post_v1_tuple_create
//
// @example
//
//	rCreate, err := cli.TupleCreate(ctx, &authz.TupleCreateRequest{
//		Tuples: []authz.Tuple{
//			authz.Tuple{
//				Resource: authz.Resource{
//					Type: "folder",
//					ID:        "folder_id",
//				},
//				Relation: "reader",
//				Subject: authz.Subject{
//					Type: "user",
//					ID:        "user_id",
//				},
//			},
//		},
//	})
func (e *authz) TupleCreate(ctx context.Context, input *TupleCreateRequest) (*pangea.PangeaResponse[TupleCreateResult], error) {
	return request.DoPost(ctx, e.Client, "v1/tuple/create", input, &TupleCreateResult{})
}

type TupleListFilter struct {
	pangea.FilterBase
	resourceType  *pangea.FilterMatch[string]
	resourceID    *pangea.FilterMatch[string]
	relation      *pangea.FilterMatch[string]
	subjectType   *pangea.FilterMatch[string]
	subjectID     *pangea.FilterMatch[string]
	subjectAction *pangea.FilterMatch[string]
}

func NewFilterUserList() *TupleListFilter {
	filter := make(pangea.Filter)
	return &TupleListFilter{
		FilterBase:    *pangea.NewFilterBase(filter),
		resourceType:  pangea.NewFilterMatch[string]("resource_type", &filter),
		resourceID:    pangea.NewFilterMatch[string]("resource_id", &filter),
		subjectType:   pangea.NewFilterMatch[string]("subject_type", &filter),
		subjectID:     pangea.NewFilterMatch[string]("subject_id", &filter),
		subjectAction: pangea.NewFilterMatch[string]("subject_action", &filter),
		relation:      pangea.NewFilterMatch[string]("relation", &filter),
	}
}

func (fu *TupleListFilter) ResourceType() *pangea.FilterMatch[string] {
	return fu.resourceType
}

func (fu *TupleListFilter) ResourceID() *pangea.FilterMatch[string] {
	return fu.resourceID
}

func (fu *TupleListFilter) SubjectType() *pangea.FilterMatch[string] {
	return fu.subjectType
}

func (fu *TupleListFilter) SubjectID() *pangea.FilterMatch[string] {
	return fu.subjectID
}

func (fu *TupleListFilter) SubjectAction() *pangea.FilterMatch[string] {
	return fu.subjectAction
}

func (fu *TupleListFilter) Relation() *pangea.FilterMatch[string] {
	return fu.relation
}

type ItemOrder string

const (
	IOasc  ItemOrder = "asc"
	IOdesc ItemOrder = "desc"
)

type TupleOrderBy string

const (
	TOBresourceType  TupleOrderBy = "resource_type"
	TOBresourceID    TupleOrderBy = "resource_id"
	TOBrelation      TupleOrderBy = "relation"
	TOBsubjectType   TupleOrderBy = "subject_type"
	TOBsubjectID     TupleOrderBy = "subject_id"
	TOBsubjectAction TupleOrderBy = "subject_action"
)

type TupleListRequest struct {
	pangea.BaseRequest

	Filter  pangea.Filter `json:"filter"`
	Size    int           `json:"size,omitempty"`
	Last    string        `json:"last,omitempty"`
	Order   ItemOrder     `json:"order,omitempty"`
	OrderBy TupleOrderBy  `json:"order_by,omitempty"`
}

type TupleListResult struct {
	Tuples []Tuple `json:"tuples"`
	Last   string  `json:"last,omitempty"`
	Count  int     `json:"count"`
}

// @summary Get tuples.
//
// @description Return a paginated list of filtered tuples. The filter is given in terms of a tuple. Fill out the fields that you want to filter. If the filter is empty it will return all the tuples.
//
// @operationId authz_post_v1_tuple_list
//
// @example
//
// filter := authz.NewFilterUserList()
// filter.ResourceType().Set(pangea.String("folder"))
// filter.ResourceID().Set(pangea.String("folder_id"))

//	rListWithResource, err := cli.TupleList(ctx, &authz.TupleListRequest{
//		Filter: filter.Filter(),
//	})
func (e *authz) TupleList(ctx context.Context, input *TupleListRequest) (*pangea.PangeaResponse[TupleListResult], error) {
	return request.DoPost(ctx, e.Client, "v1/tuple/list", input, &TupleListResult{})
}

type TupleDeleteRequest struct {
	pangea.BaseRequest

	Tuples []Tuple `json:"tuples"`
}

type TupleDeleteResult struct {
}

// @summary Delete tuples.
//
// @description Delete tuples.
//
// @operationId authz_post_v1_tuple_delete
//
// @example
//
//	rDelete, err := cli.TupleDelete(ctx, &authz.TupleDeleteRequest{
//		Tuples: []authz.Tuple{
//			authz.Tuple{
//				Resource: authz.Resource{
//					Type: "folder",
//					ID:        "folder_id",
//				},
//				Relation: "editor",
//				Subject: authz.Subject{
//					Type: "user",
//					ID:        "user_id",
//				},
//			},
//		},
//	})
func (e *authz) TupleDelete(ctx context.Context, input *TupleDeleteRequest) (*pangea.PangeaResponse[TupleDeleteResult], error) {
	return request.DoPost(ctx, e.Client, "v1/tuple/delete", input, &TupleDeleteResult{})
}

type CheckRequest struct {
	pangea.BaseRequest

	Resource   Resource       `json:"resource"`
	Action     string         `json:"action"`
	Subject    Subject        `json:"subject"`
	Debug      *bool          `json:"debug,omitempty"`
	Attributes map[string]any `json:"attributes,omitempty"`
}

type DebugPath struct {
	Type   string `json:"type"`
	ID     string `json:"id"`
	Action string `json:"action"`
}

type Debug struct {
	Path []DebugPath `json:"path"`
}

type CheckResult struct {
	SchemaID      string `json:"schema_id"`
	SchemaVersion int    `json:"schema_version"`
	Allowed       bool   `json:"allowed"`
	Depth         int    `json:"depth"`
	Debug         *Debug `json:"debug,omitempty"`
}

// @summary Perform a check request.
//
// @description Check if a subject has permission to do action on the resource.
//
// @operationId authz_post_v1_check
//
// @example
//
//	rCheck, err = cli.Check(ctx, &authz.CheckRequest{
//		Resource: authz.Resource{
//			Type: "folder",
//			ID:        "folder_id",
//		},
//		Action: "editor",
//		Subject: authz.Subject{
//			Type: "user",
//			ID:        "user_id",
//		},
//		Debug: pangea.Bool(true),
//	})
func (e *authz) Check(ctx context.Context, input *CheckRequest) (*pangea.PangeaResponse[CheckResult], error) {
	return request.DoPost(ctx, e.Client, "v1/check", input, &CheckResult{})
}

type ListResourcesRequest struct {
	pangea.BaseRequest

	Type       string         `json:"type"`
	Action     string         `json:"action"`
	Subject    Subject        `json:"subject"`
	Attributes map[string]any `json:"attributes,omitempty"`
}

type ListResourcesResult struct {
	IDs []string `json:"ids"`
}

// @summary List resources.
//
// @description Given a type, action, and subject, list all the resources in the type that the subject has permission to the action with.
//
// @operationId authz_post_v1_list-resources
//
// @example
//
//	rListResources, err := cli.ListResources(ctx, &authz.ListResourcesRequest{
//		Type: "folder",
//		Action:    "editor",
//		Subject: authz.Subject{
//			Type: "user",
//			ID:        "user_id",
//		},
//	})
func (e *authz) ListResources(ctx context.Context, input *ListResourcesRequest) (*pangea.PangeaResponse[ListResourcesResult], error) {
	return request.DoPost(ctx, e.Client, "v1/list-resources", input, &ListResourcesResult{})
}

type ListSubjectsRequest struct {
	pangea.BaseRequest

	Resource   Resource       `json:"resource"`
	Action     string         `json:"action"`
	Attributes map[string]any `json:"attributes,omitempty"` // A JSON object of attribute data.
}

type ListSubjectsResult struct {
	Subjects []Subject `json:"subjects"`
}

// @summary List subjects.
//
// @description Given a resource and an action, return the list of subjects who have the given action to the given resource.
//
// @operationId authz_post_v1_list-subjects
//
// @example
//
//	rListSubjects, err := cli.ListSubjects(ctx, &authz.ListSubjectsRequest{
//		Resource: authz.Resource{
//			Type: "folder",
//			ID:        "folder_id",
//		},
//		Action: "editor",
//	})
func (e *authz) ListSubjects(ctx context.Context, input *ListSubjectsRequest) (*pangea.PangeaResponse[ListSubjectsResult], error) {
	return request.DoPost(ctx, e.Client, "v1/list-subjects", input, &ListSubjectsResult{})
}

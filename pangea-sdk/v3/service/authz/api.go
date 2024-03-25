package authz

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/internal/request"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
)

type Resource struct {
	Namespace string `json:"namespace"`
	ID        string `json:"id,omitempty"`
}

type Subject struct {
	Namespace string `json:"namespace"`
	ID        string `json:"id,omitempty"`
	Action    string `json:"action,omitempty"`
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

// @summary Write tuples (Beta).
//
// @description Write tuples. The request will fail if tuples do not validate against the schema defined resource types.
//
// @operationId authz_post_v1beta_tuple_create
//
// @example
//
//	rCreate, err := cli.TupleCreate(ctx, &authz.TupleCreateRequest{
//		Tuples: []authz.Tuple{
//			authz.Tuple{
//				Resource: authz.Resource{
//					Namespace: "folder",
//					ID:        "folder_id",
//				},
//				Relation: "reader",
//				Subject: authz.Subject{
//					Namespace: "user",
//					ID:        "user_id",
//				},
//			},
//		},
//	})
func (e *authz) TupleCreate(ctx context.Context, input *TupleCreateRequest) (*pangea.PangeaResponse[TupleCreateResult], error) {
	return request.DoPost(ctx, e.Client, "v1beta/tuple/create", input, &TupleCreateResult{})
}

// type TupleListFilter struct {
// 	Resource *Resource `json:"resource,omitempty"`
// 	Relation string    `json:"relation,omitempty"`
// 	Subject  *Subject  `json:"subject,omitempty"`
// }

type TupleListFilter struct {
	pangea.FilterBase
	resourceNamespace *pangea.FilterMatch[string]
	resourceID        *pangea.FilterMatch[string]
	relation          *pangea.FilterMatch[string]
	subjectNamespace  *pangea.FilterMatch[string]
	subjectID         *pangea.FilterMatch[string]
	subjectAction     *pangea.FilterMatch[string]
}

func NewFilterUserList() *TupleListFilter {
	filter := make(pangea.Filter)
	return &TupleListFilter{
		FilterBase:        *pangea.NewFilterBase(filter),
		resourceNamespace: pangea.NewFilterMatch[string]("resource_namespace", &filter),
		resourceID:        pangea.NewFilterMatch[string]("resource_id", &filter),
		subjectNamespace:  pangea.NewFilterMatch[string]("subject_namespace", &filter),
		subjectID:         pangea.NewFilterMatch[string]("subject_id", &filter),
		subjectAction:     pangea.NewFilterMatch[string]("subject_action", &filter),
		relation:          pangea.NewFilterMatch[string]("relation", &filter),
	}
}

func (fu *TupleListFilter) ResourceNamespace() *pangea.FilterMatch[string] {
	return fu.resourceNamespace
}

func (fu *TupleListFilter) ResourceID() *pangea.FilterMatch[string] {
	return fu.resourceID
}

func (fu *TupleListFilter) SubjectNamespace() *pangea.FilterMatch[string] {
	return fu.subjectNamespace
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
	TOBresourceNamespace TupleOrderBy = "resource_namespace"
	TOBresourceID        TupleOrderBy = "resource_id"
	TOBrelation          TupleOrderBy = "relation"
	TOBsubjectNamespace  TupleOrderBy = "subject_namespace"
	TOBsubjectID         TupleOrderBy = "subject_id"
	TOBsubjectAction     TupleOrderBy = "subject_action"
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

// @summary Get tuples (Beta).
//
// @description Return a paginated list of filtered tuples. The filter is given in terms of a tuple. Fill out the fields that you want to filter. If the filter is empty it will return all the tuples.
//
// @operationId authz_post_v1beta_tuple_list
//
// @example
//
// filter := authz.NewFilterUserList()
// filter.ResourceNamespace().Set(pangea.String("folder"))
// filter.ResourceID().Set(pangea.String("folder_id"))

//	rListWithResource, err := cli.TupleList(ctx, &authz.TupleListRequest{
//		Filter: filter.Filter(),
//	})
func (e *authz) TupleList(ctx context.Context, input *TupleListRequest) (*pangea.PangeaResponse[TupleListResult], error) {
	return request.DoPost(ctx, e.Client, "v1beta/tuple/list", input, &TupleListResult{})
}

type TupleDeleteRequest struct {
	pangea.BaseRequest

	Tuples []Tuple `json:"tuples"`
}

type TupleDeleteResult struct {
}

// @summary Delete tuples (Beta).
//
// @description Delete tuples.
//
// @operationId authz_post_v1beta_tuple_delete
//
// @example
//
//	rDelete, err := cli.TupleDelete(ctx, &authz.TupleDeleteRequest{
//		Tuples: []authz.Tuple{
//			authz.Tuple{
//				Resource: authz.Resource{
//					Namespace: "folder",
//					ID:        "folder_id",
//				},
//				Relation: "editor",
//				Subject: authz.Subject{
//					Namespace: "user",
//					ID:        "user_id",
//				},
//			},
//		},
//	})
func (e *authz) TupleDelete(ctx context.Context, input *TupleDeleteRequest) (*pangea.PangeaResponse[TupleDeleteResult], error) {
	return request.DoPost(ctx, e.Client, "v1beta/tuple/delete", input, &TupleDeleteResult{})
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
	Namespace string `json:"namespace"`
	ID        string `json:"id"`
	Action    string `json:"action"`
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

// @summary Perform a check request (Beta).
//
// @description Check if a subject has permission to do action on the resource.
//
// @operationId authz_post_v1beta_check
//
// @example
//
//	rCheck, err = cli.Check(ctx, &authz.CheckRequest{
//		Resource: authz.Resource{
//			Namespace: "folder",
//			ID:        "folder_id",
//		},
//		Action: "editor",
//		Subject: authz.Subject{
//			Namespace: "user",
//			ID:        "user_id",
//		},
//		Debug: pangea.Bool(true),
//	})
func (e *authz) Check(ctx context.Context, input *CheckRequest) (*pangea.PangeaResponse[CheckResult], error) {
	return request.DoPost(ctx, e.Client, "v1beta/check", input, &CheckResult{})
}

type ListResourcesRequest struct {
	pangea.BaseRequest

	Namespace string  `json:"namespace"`
	Action    string  `json:"action"`
	Subject   Subject `json:"subject"`
}

type ListResourcesResult struct {
	IDs []string `json:"ids"`
}

// @summary List resources (Beta).
//
// @description Given a namespace, action, and subject, list all the resources in the namespace that the subject has permission to the action with.
//
// @operationId authz_post_v1beta_list-resources
//
// @example
//
//	rListResources, err := cli.ListResources(ctx, &authz.ListResourcesRequest{
//		Namespace: "folder",
//		Action:    "editor",
//		Subject: authz.Subject{
//			Namespace: "user",
//			ID:        "user_id",
//		},
//	})
func (e *authz) ListResources(ctx context.Context, input *ListResourcesRequest) (*pangea.PangeaResponse[ListResourcesResult], error) {
	return request.DoPost(ctx, e.Client, "v1beta/list-resources", input, &ListResourcesResult{})
}

type ListSubjectsRequest struct {
	pangea.BaseRequest

	Resource Resource `json:"resource"`
	Action   string   `json:"action"`
}

type ListSubjectsResult struct {
	Subjects []Subject `json:"subjects"`
}

// @summary List subjects (Beta).
//
// @description Given a resource and an action, return the list of subjects who have the given action to the given resource.
//
// @operationId authz_post_v1beta_list-subjects
//
// @example
//
//	rListSubjects, err := cli.ListSubjects(ctx, &authz.ListSubjectsRequest{
//		Resource: authz.Resource{
//			Namespace: "folder",
//			ID:        "folder_id",
//		},
//		Action: "editor",
//	})
func (e *authz) ListSubjects(ctx context.Context, input *ListSubjectsRequest) (*pangea.PangeaResponse[ListSubjectsResult], error) {
	return request.DoPost(ctx, e.Client, "v1beta/list-subjects", input, &ListSubjectsResult{})
}

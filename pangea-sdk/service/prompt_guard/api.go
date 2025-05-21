package prompt_guard

import (
	"context"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/internal/request"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
)

// @summary Guard
//
// @description Guard messages.
//
// @operationId prompt_guard_post_v1_guard
//
// @example
//
//	input := &prompt_guard.GuardRequest{Messages: []prompt_guard.Message{{Role: "user", Content: "how are you?"}}}
//	response, err := client.Guard(ctx, input)
func (e *promptGuard) Guard(ctx context.Context, input *GuardRequest) (*pangea.PangeaResponse[GuardResult], error) {
	return request.DoPost(ctx, e.Client, "v1/guard", input, &GuardResult{})
}

// @operationId prompt_guard_post_v1beta_config
func (e *promptGuard) GetServiceConfig(ctx context.Context, body GetServiceConfigParams) (*pangea.PangeaResponse[struct{}], error) {
	var result struct{}
	return request.DoPost(ctx, e.Client, "v1beta/config", &body, &result)
}

// @operationId prompt_guard_post_v1beta_config_create
func (e *promptGuard) CreateServiceConfig(ctx context.Context, body CreateServiceConfigParams) (*pangea.PangeaResponse[struct{}], error) {
	var result struct{}
	return request.DoPost(ctx, e.Client, "v1beta/config/create", &body, &result)
}

// @operationId prompt_guard_post_v1beta_config_update
func (e *promptGuard) UpdateServiceConfig(ctx context.Context, body UpdateServiceConfigParams) (*pangea.PangeaResponse[struct{}], error) {
	var result struct{}
	return request.DoPost(ctx, e.Client, "v1beta/config/update", &body, &result)
}

// @operationId prompt_guard_post_v1beta_config_delete
func (e *promptGuard) DeleteServiceConfig(ctx context.Context, body DeleteServiceConfigParams) (*pangea.PangeaResponse[struct{}], error) {
	var result struct{}
	return request.DoPost(ctx, e.Client, "v1beta/config/delete", &body, &result)
}

// @operationId prompt_guard_post_v1beta_config_list
func (e *promptGuard) ListServiceConfigs(ctx context.Context, body ListServiceConfigsParams) (*pangea.PangeaResponse[ListServiceConfigsResult], error) {
	return request.DoPost(ctx, e.Client, "v1beta/config/list", &body, &ListServiceConfigsResult{})
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Classification struct {
	Category   string  `json:"category"`   // Classification category
	Detected   bool    `json:"detected"`   // Classification detection result
	Confidence float32 `json:"confidence"` // Confidence score for the classification
}

type GuardRequest struct {
	pangea.BaseRequest

	Messages  []Message `json:"messages"`            // Prompt content and role array. The content is the text that will be analyzed for redaction.
	Analyzers []string  `json:"analyzers,omitempty"` // Specific analyzers to be used in the call
	Classify  bool      `json:"classify,omitempty"`  // Boolean to enable classification of the content
}

type GuardResult struct {
	Detected        bool             `json:"detected"`           // Boolean response for if the prompt was considered malicious or not
	Type            string           `json:"type,omitempty"`     // Type of analysis, either direct or indirect
	Analyzer        string           `json:"analyzer,omitempty"` // Prompt Analyzers for identifying and rejecting properties of prompts
	Confidence      float32          `json:"confidence"`         // Percent of confidence in the detection result, ranging from 0 to 1
	Info            string           `json:"info,omitempty"`     // Extra information about the detection result
	Classifications []Classification `json:"classifications"`    // List of classification results with labels and confidence scores
}

type AuditDataActivityConfigAreas struct {
	TextGuard bool `json:"text_guard"`
}

type AuditDataActivityConfig struct {
	Enabled              bool                         `json:"enabled"`
	AuditServiceConfigId string                       `json:"audit_service_config_id"`
	Areas                AuditDataActivityConfigAreas `json:"areas"`
}

type ServiceConfig struct {
	Id                          *string                  `json:"id,omitempty"`
	Version                     *string                  `json:"version,omitempty"`
	Analyzers                   map[string]bool          `json:"analyzers,omitempty"`
	MaliciousDetectionThreshold *float32                 `json:"malicious_detection_threshold,omitempty"`
	BenignDetectionThreshold    *float32                 `json:"benign_detection_threshold,omitempty"`
	AuditDataActivity           *AuditDataActivityConfig `json:"audit_data_activity,omitempty"`
}

type GetServiceConfigParams struct {
	pangea.BaseRequest
	ServiceConfig
}

type CreateServiceConfigParams struct {
	pangea.BaseRequest
	ServiceConfig
}

type UpdateServiceConfigParams struct {
	pangea.BaseRequest
	ServiceConfig
}

type DeleteServiceConfigParams struct {
	pangea.BaseRequest

	Id string `json:"id"`
}

type ServiceConfigListFilter struct {
	Id           *string    `json:"id,omitempty"`              // Only records where id equals this value.
	IdContains   []string   `json:"id__contains,omitempty"`    // Only records where id includes each substring.
	IdIn         []string   `json:"id__in,omitempty"`          // Only records where id equals one of the provided substrings.
	CreatedAt    *time.Time `json:"created_at,omitempty"`      // Only records where created_at equals this value.
	CreatedAtGt  *time.Time `json:"created_at__gt,omitempty"`  // Only records where created_at is greater than this value.
	CreatedAtGte *time.Time `json:"created_at__gte,omitempty"` // Only records where created_at is greater than or equal to this value.
	CreatedAtLt  *time.Time `json:"created_at__lt,omitempty"`  // Only records where created_at is less than this value.
	CreatedAtLte *time.Time `json:"created_at__lte,omitempty"` // Only records where created_at is less than or equal to this value.
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`      // Only records where updated_at equals this value.
	UpdatedAtGt  *time.Time `json:"updated_at__gt,omitempty"`  // Only records where updated_at is greater than this value.
	UpdatedAtGte *time.Time `json:"updated_at__gte,omitempty"` // Only records where updated_at is greater than or equal to this value.
	UpdatedAtLt  *time.Time `json:"updated_at__lt,omitempty"`  // Only records where updated_at is less than this value.
	UpdatedAtLte *time.Time `json:"updated_at__lte,omitempty"` // Only records where updated_at is less than or equal to this value.
}

type ListServiceConfigsParams struct {
	pangea.BaseRequest

	Filter  *ServiceConfigListFilter `json:"filter,omitempty"`
	Last    *string                  `json:"last,omitempty"`     // Reflected value from a previous response to obtain the next page of results.
	Order   *string                  `json:"order,omitempty"`    // Order results asc(ending) or desc(ending).
	OrderBy *string                  `json:"order_by,omitempty"` // Which field to order results by.
	Size    *int                     `json:"size,omitempty"`     // Maximum results to include in the response.
}

type ListServiceConfigsResult struct {
	Count *int            `json:"count,omitempty"` // The total number of service configs matched by the list request.
	Last  *string         `json:"last,omitempty"`  // Used to fetch the next page of the current listing when provided in a repeated request's last parameter.
	Items []ServiceConfig `json:"items,omitempty"`
}

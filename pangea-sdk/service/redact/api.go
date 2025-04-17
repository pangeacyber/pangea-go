package redact

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/internal/request"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
)

// @summary Redact
//
// @description Redacts the content of a single text string.
//
// @operationId redact_post_v1_redact
//
// @example
//
//	input := &redact.TextInput{
//		Text: pangea.String("my phone number is 123-456-7890"),
//	}
//
//	redactResponse, err := redactcli.Redact(ctx, input)
func (r *redact) Redact(ctx context.Context, input *TextRequest) (*pangea.PangeaResponse[TextResult], error) {
	return request.DoPost(ctx, r.Client, "v1/redact", input, &TextResult{})
}

// @summary Redact structured
//
// @description Redacts text within a structured object.
//
// @operationId redact_post_v1_redact_structured
//
// @example
//
//	data := yourCustomDataStruct{
//		Secret: "My social security number is 0303456",
//	}
//
//	input := &redact.StructuredInput{
//		Data: data,
//	}
//
//	redactResponse, err := redactcli.RedactStructured(ctx, input)
func (r *redact) RedactStructured(ctx context.Context, input *StructuredRequest) (*pangea.PangeaResponse[StructuredResult], error) {
	return request.DoPost(ctx, r.Client, "v1/redact_structured", input, &StructuredResult{})
}

// @summary Unredact
//
// @description Decrypt or unredact fpe redactions
//
// @operationId redact_post_v1_unredact
//
// @example
//
//	redactcli.Unredact(ctx, &UnredactRequest{
//		RedactedData: "redacted data",
//		FPEContext:   "gAyHpblmIoUXKTiYY8xKiQ==",
//	})
func (r *redact) Unredact(ctx context.Context, input *UnredactRequest) (*pangea.PangeaResponse[UnredactResult], error) {
	return request.DoPost(ctx, r.Client, "v1/unredact", input, &UnredactResult{})
}

type GetServiceConfigConfigRequest struct {
	pangea.BaseRequest

	Id string `json:"id"`
}

// @summary Get a service config.
//
// @description Get a service config.
//
// @operationId redact_post_v1beta_config
func (a *redact) GetServiceConfig(ctx context.Context, configId string) (*pangea.PangeaResponse[ServiceConfig], error) {
	return request.DoPost(ctx, a.Client, "v1beta/config", &GetServiceConfigConfigRequest{Id: configId}, &ServiceConfig{})
}

// @summary Create a service config.
//
// @description Create a service config.
//
// @operationId redact_post_v1beta_config_create
func (a *redact) CreateServiceConfig(ctx context.Context, input *CreateServiceConfigRequest) (*pangea.PangeaResponse[ServiceConfig], error) {
	return request.DoPost(ctx, a.Client, "v1beta/config/create", input, &ServiceConfig{})
}

// @summary Update a service config.
//
// @description Update a service config.
//
// @operationId redact_post_v1beta_config_update
func (a *redact) UpdateServiceConfig(ctx context.Context, input *UpdateServiceConfigRequest) (*pangea.PangeaResponse[ServiceConfig], error) {
	return request.DoPost(ctx, a.Client, "v1beta/config/update", input, &ServiceConfig{})
}

// @summary Delete a service config.
//
// @description Delete a service config.
//
// @operationId redact_post_v1beta_config_delete
func (a *redact) DeleteServiceConfig(ctx context.Context, configId string) (*pangea.PangeaResponse[ServiceConfig], error) {
	return request.DoPost(ctx, a.Client, "v1beta/config/delete", &DeleteServiceConfigRequest{Id: configId}, &ServiceConfig{})
}

// @summary List service configs.
//
// @description List service configs.
//
// @operationId redact_post_v1beta_config_list
func (a *redact) ListServiceConfigs(ctx context.Context, input *ListServiceConfigsRequest) (*pangea.PangeaResponse[ServiceConfig], error) {
	return request.DoPost(ctx, a.Client, "v1beta/config/list", input, &ServiceConfig{})
}

type TextRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	// The text to be redacted.
	// Text is a required field.
	Text *string `json:"text"`

	// If the response should include some debug Info.
	Debug *bool `json:"debug,omitempty"`

	// An array of redact rule short names
	Rules []string `json:"rules,omitempty"`

	// An array of redact rulesets short names
	Rulesets []string `json:"rulesets,omitempty"`

	// Setting this value to false will omit the redacted result only returning count
	ReturnResult *bool `json:"return_result,omitempty"`

	// A set of redaction method overrides for any enabled rule. These methods override the config declared methods
	RedactionMethodOverrides map[string]Redaction `json:"redaction_method_overrides,omitempty"`

	VaultParameters *VaultParameters `json:"vault_parameters,omitempty"`

	// Is this redact call going to be used in an LLM request?
	LLMrequest *bool `json:"llm_request,omitempty"`
}

type VaultParameters struct {
	// A vault key ID of an exportable key used to redact with FPE instead of using the service config default.
	FPEkeyID string `json:"fpe_key_id,omitempty"`

	// A vault secret ID of a secret used to salt a hash instead of using the service config default.
	SaltSecretID string `json:"salt_secret_id,omitempty"`
}

type TextResult struct {
	// The redacted text.
	RedactedText *string `json:"redacted_text"`

	// Number of redactions present in the response
	Count int `json:"count"`

	// Describes the decision process for redactions
	Report *DebugReport `json:"report"`

	// FPE context used to encrypt and redact data
	FPEContext *string `json:"fpe_context"`
}

type DebugReport struct {
	SummaryCounts     map[string]int     `json:"summary_counts"`
	RecognizerResults []RecognizerResult `json:"recognizer_results"`
}

type RecognizerResult struct {
	// FieldType is always populated on a successful response.
	FieldType string `json:"field_type"`

	// Score is always populated on a successful response.
	Score *float64 `json:"score"`

	// Text is always populated on a successful response.
	Text string `json:"text"`

	// Start is always populated on a successful response.
	Start int `json:"start"`

	// End is always populated on a successful response.
	End int `json:"end"`

	// Redacted is always populated on a successful response.
	Redacted bool `json:"redacted"`

	// DataKey is always populated on a successful response.
	DataKey string `json:"data_key"`
}

type StructuredRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	// Structured data to redact
	// Data is a required field.
	Data map[string]any `json:"data"`

	// JSON path(s) used to identify the specific JSON fields to redact in the structured data.
	// Note: If jsonp parameter is used, the data parameter must be in JSON format.
	JSONP []*string `json:"jsonp,omitempty"`

	// The format of the structured data to redact.
	Format *string `json:"format,omitempty"`

	// Setting this value to true will provide a detailed analysis of the redacted data and the rules that caused redaction.
	Debug *bool `json:"debug,omitempty"`

	// An array of redact rule short names
	Rules []string `json:"rules,omitempty"`

	// An array of redact rulesets short names
	Rulesets []string `json:"rulesets,omitempty"`

	// Setting this value to false will omit the redacted result only returning count
	ReturnResult *bool `json:"return_result,omitempty"`

	// A set of redaction method overrides for any enabled rule. These methods override the config declared methods
	RedactionMethodOverrides map[string]Redaction `json:"redaction_method_overrides,omitempty"`

	VaultParameters *VaultParameters `json:"vault_parameters,omitempty"`

	// Is this redact call going to be used in an LLM request?
	LLMrequest *bool `json:"llm_request,omitempty"`
}

type StructuredResult struct {
	// RedactedData is always populated on a successful response.
	RedactedData map[string]any `json:"redacted_data"`

	// Number of redactions present in the response
	Count int `json:"count"`

	// Describes the decision process for redactions
	Report *DebugReport `json:"report"`

	// FPE context used to encrypt and redact data
	FPEContext *string `json:"fpe_context"`
}

type UnredactRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	// Data to unredact
	RedactedData any `json:"redacted_data"`

	// FPE context used to decrypt and unredact data (in base64)
	FPEContext string `json:"fpe_context"`
}

type UnredactResult struct {
	Data any `json:"data"`
}

type MaskingType string
type RedactType string
type FPEAlphabet string

const (
	RTmask           RedactType = "mask"
	RTpartialMasking RedactType = "partial_masking"
	RTreplacement    RedactType = "replacement"
	RTdetectOnly     RedactType = "detect_only"
	RThash           RedactType = "hash"
	RTfpe            RedactType = "fpe"
)

const (
	FPEAnumeric           FPEAlphabet = "numeric"
	FPEAalphaNumericLower FPEAlphabet = "alphanumericlower"
	FPEAalphanumeric      FPEAlphabet = "alphanumeric"
)

const (
	MTmask   MaskingType = "mask"
	MTunmask MaskingType = "unmask"
)

type PartialMasking struct {
	MaskingType       *MaskingType `json:"masking_type,omitempty"`
	UnmaskedFromLeft  *int         `json:"unmasked_from_left,omitempty"`
	UnmaskedFromRight *int         `json:"unmasked_from_right,omitempty"`
	MaskedFromLeft    *int         `json:"masked_from_left,omitempty"`
	MaskedFromRight   *int         `json:"masked_from_right,omitempty"`
	CharsToIgnore     []string     `json:"chars_to_ignore,omitempty"`
	MaskingChar       []string     `json:"masking_char,omitempty"`
}

type Redaction struct {
	RedactionType  RedactType             `json:"redaction_type"`
	Hash           map[string]interface{} `json:"hash,omitempty"`
	FPEAlphabet    *FPEAlphabet           `json:"fpe_alphabet,omitempty"`
	PartialMasking *PartialMasking        `json:"partial_masking,omitempty"`
	RedactionValue *string                `json:"redaction_value,omitempty"`
}

type Rule struct {
	EntityName            string   `json:"entity_name"`
	MatchThreshold        *float32 `json:"match_threshold,omitempty"`
	ContextValues         []string `json:"context_values,omitempty"`
	NegativeContextValues []string `json:"negative_context_values,omitempty"`
	Name                  *string  `json:"name,omitempty"`
	Description           *string  `json:"description,omitempty"`
}

type Ruleset struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Rules       []string `json:"rules,omitempty"`
}

type ServiceConfig struct {
	Version      *string  `json:"version,omitempty"`
	Id           *string  `json:"id,omitempty"`
	Name         *string  `json:"name,omitempty"`
	UpdatedAt    *string  `json:"updated_at,omitempty"`
	EnabledRules []string `json:"enabled_rules,omitempty"`
	// Always run service config enabled rules across all redact calls regardless of flags?
	EnforceEnabledRules *bool                `json:"enforce_enabled_rules,omitempty"`
	Redactions          map[string]Redaction `json:"redactions,omitempty"`
	// Service config used to create the secret
	VaultServiceConfigId *string `json:"vault_service_config_id,omitempty"`
	// Pangea only allows hashing to be done using a salt value to prevent brute-force attacks.
	SaltVaultSecretId *string `json:"salt_vault_secret_id,omitempty"`
	// The ID of the key used by FF3 Encryption algorithms for FPE.
	FpeVaultSecretId   *string            `json:"fpe_vault_secret_id,omitempty"`
	Rules              map[string]Rule    `json:"rules,omitempty"`
	Rulesets           map[string]Ruleset `json:"rulesets,omitempty"`
	SupportedLanguages []string           `json:"supported_languages,omitempty"`
}

type CreateServiceConfigRequest struct {
	pangea.BaseRequest

	// An ID for a service config
	Id string `json:"id" validate:"regexp=^pci_[a-z2-7]{32}$"`
}

type UpdateServiceConfigRequest struct {
	pangea.BaseRequest
	ServiceConfig
}

type DeleteServiceConfigRequest struct {
	pangea.BaseRequest

	// An ID for a service config
	Id string `json:"id" validate:"regexp=^pci_[a-z2-7]{32}$"`
}

type ServiceConfigListFilter struct {
	// Only records where id equals this value.
	Id *string `json:"id,omitempty"`
	// Only records where id includes each substring.
	IdContains []string `json:"id__contains,omitempty"`
	// Only records where id equals one of the provided substrings.
	IdIn []string `json:"id__in,omitempty"`
	// Only records where created_at equals this value.
	CreatedAt *string `json:"created_at,omitempty"`
	// Only records where created_at is greater than this value.
	CreatedAtGt *string `json:"created_at__gt,omitempty"`
	// Only records where created_at is greater than or equal to this value.
	CreatedAtGte *string `json:"created_at__gte,omitempty"`
	// Only records where created_at is less than this value.
	CreatedAtLt *string `json:"created_at__lt,omitempty"`
	// Only records where created_at is less than or equal to this value.
	CreatedAtLte *string `json:"created_at__lte,omitempty"`
	// Only records where updated_at equals this value.
	UpdatedAt *string `json:"updated_at,omitempty"`
	// Only records where updated_at is greater than this value.
	UpdatedAtGt *string `json:"updated_at__gt,omitempty"`
	// Only records where updated_at is greater than or equal to this value.
	UpdatedAtGte *string `json:"updated_at__gte,omitempty"`
	// Only records where updated_at is less than this value.
	UpdatedAtLt *string `json:"updated_at__lt,omitempty"`
	// Only records where updated_at is less than or equal to this value.
	UpdatedAtLte *string `json:"updated_at__lte,omitempty"`
}

// ListServiceConfigsRequest List or filter/search records.
type ListServiceConfigsRequest struct {
	pangea.BaseRequest

	Filter *ServiceConfigListFilter `json:"filter,omitempty"`
	// Reflected value from a previous response to obtain the next page of results.
	Last *string `json:"last,omitempty"`
	// Order results asc(ending) or desc(ending).
	Order *string `json:"order,omitempty"`
	// Which field to order results by.
	OrderBy *string `json:"order_by,omitempty"`
	// Maximum results to include in the response.
	Size *int32 `json:"size,omitempty"`
}

package ai_guard

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/internal/request"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
)

// @summary Text guard
//
// @description Guard text.
//
// @operationId ai_guard_post_v1_text_guard
//
// @example
//
//	input := &ai_guard.TextGuardRequest{Text: "hello world"}
//	response, err := client.GuardText(ctx, input)
func (e *aiGuard) GuardText(ctx context.Context, input *TextGuardRequest) (*pangea.PangeaResponse[TextGuardResult], error) {
	if input.Text == "" && input.Messages == nil {
		return nil, errors.New("one of `Text` or `Messages` must be defined")
	}

	return request.DoPost(ctx, e.Client, "v1/text/guard", input, &TextGuardResult{})
}

// @summary Guard LLM input and output
//
// @description Analyze and redact content to avoid manipulation of the model, addition of malicious content, and other undesirable data transfers.
//
// @operationId ai_guard_post_v1beta_guard
func (e *aiGuard) Guard(ctx context.Context, body GuardRequest) (*pangea.PangeaResponse[GuardResult], error) {
	return request.DoPost(ctx, e.Client, "v1beta/guard", &body, &GuardResult{})
}

// @operationId ai_guard_post_v1beta_config
func (e *aiGuard) GetServiceConfig(ctx context.Context, body GetServiceConfigParams) (*pangea.PangeaResponse[ServiceConfig], error) {
	return request.DoPost(ctx, e.Client, "v1beta/config", &body, &ServiceConfig{})
}

// @operationId ai_guard_post_v1beta_config_create
func (e *aiGuard) CreateServiceConfig(ctx context.Context, body CreateServiceConfigParams) (*pangea.PangeaResponse[ServiceConfig], error) {
	return request.DoPost(ctx, e.Client, "v1beta/config/create", &body, &ServiceConfig{})
}

// @operationId ai_guard_post_v1beta_config_update
func (e *aiGuard) UpdateServiceConfig(ctx context.Context, body UpdateServiceConfigParams) (*pangea.PangeaResponse[ServiceConfig], error) {
	return request.DoPost(ctx, e.Client, "v1beta/config/update", &body, &ServiceConfig{})
}

// @operationId ai_guard_post_v1beta_config_delete
func (e *aiGuard) DeleteServiceConfig(ctx context.Context, body DeleteServiceConfigParams) (*pangea.PangeaResponse[ServiceConfig], error) {
	return request.DoPost(ctx, e.Client, "v1beta/config/delete", &body, &ServiceConfig{})
}

// @operationId ai_guard_post_v1beta_config_list
func (e *aiGuard) ListServiceConfigs(ctx context.Context, body ListServiceConfigsParams) (*pangea.PangeaResponse[ListServiceConfigsResult], error) {
	return request.DoPost(ctx, e.Client, "v1beta/config/list", &body, &ListServiceConfigsResult{})
}

type TopicDetectionOverride struct {
	Disabled  *bool    `json:"disabled,omitempty"`
	Action    *string  `json:"action,omitempty"`
	Topics    []string `json:"topics,omitempty"`
	Threshold *float64 `json:"threshold,omitempty"`
}

// This is named "prompt injection" in the API spec even though it is also used
// for many other detectors.
type PromptInjectionAction string

const (
	PromptInjectionActionReport PromptInjectionAction = "report"
	PromptInjectionActionBlock  PromptInjectionAction = "block"
)

type MaliciousEntityAction string

const (
	MaliciousEntityActionReport   MaliciousEntityAction = "report"
	MaliciousEntityActionDefang   MaliciousEntityAction = "defang"
	MaliciousEntityActionDisabled MaliciousEntityAction = "disabled"
	MaliciousEntityActionBlock    MaliciousEntityAction = "block"
)

type PiiEntityAction string

const (
	PiiEntityActionDisabled       PiiEntityAction = "disabled"
	PiiEntityActionReport         PiiEntityAction = "report"
	PiiEntityActionBlock          PiiEntityAction = "block"
	PiiEntityActionMask           PiiEntityAction = "mask"
	PiiEntityActionPartialMasking PiiEntityAction = "partial_masking"
	PiiEntityActionReplacement    PiiEntityAction = "replacement"
	PiiEntityActionHash           PiiEntityAction = "hash"
	PiiEntityActionFPE            PiiEntityAction = "fpe"
)

// Override models
type CodeDetectionOverride struct {
	Disabled *bool   `json:"disabled,omitempty"`
	Action   *string `json:"action,omitempty"`
}

type LanguageDetectionOverride struct {
	Disabled *bool    `json:"disabled,omitempty"`
	Allow    []string `json:"allow,omitempty"`
	Block    []string `json:"block,omitempty"`
	Report   []string `json:"report,omitempty"`
}

type PromptInjectionOverride struct {
	Disabled *bool                  `json:"disabled,omitempty"`
	Action   *PromptInjectionAction `json:"action,omitempty"`
}

type SelfHarmOverride struct {
	Disabled  *bool                  `json:"disabled,omitempty"`
	Action    *PromptInjectionAction `json:"action,omitempty"`
	Threshold *float64               `json:"threshold,omitempty"`
}

type GibberishOverride struct {
	Disabled *bool                  `json:"disabled,omitempty"`
	Action   *PromptInjectionAction `json:"action,omitempty"`
}

type RoleplayOverride struct {
	Disabled *bool                  `json:"disabled,omitempty"`
	Action   *PromptInjectionAction `json:"action,omitempty"`
}

type SentimentOverride struct {
	Disabled  *bool                  `json:"disabled,omitempty"`
	Action    *PromptInjectionAction `json:"action,omitempty"`
	Threshold *float64               `json:"threshold,omitempty"`
}

type MaliciousEntityOverride struct {
	Disabled  *bool                  `json:"disabled,omitempty"`
	IPAddress *MaliciousEntityAction `json:"ip_address,omitempty"`
	URL       *MaliciousEntityAction `json:"url,omitempty"`
	Domain    *MaliciousEntityAction `json:"domain,omitempty"`
}

type CompetitorsOverride struct {
	Disabled *bool                  `json:"disabled,omitempty"`
	Action   *PromptInjectionAction `json:"action,omitempty"`
}

type PiiEntityOverride struct {
	Disabled         *bool            `json:"disabled,omitempty"`
	EmailAddress     *PiiEntityAction `json:"email_address,omitempty"`
	NRP              *PiiEntityAction `json:"nrp,omitempty"`
	Location         *PiiEntityAction `json:"location,omitempty"`
	Person           *PiiEntityAction `json:"person,omitempty"`
	PhoneNumber      *PiiEntityAction `json:"phone_number,omitempty"`
	DateTime         *PiiEntityAction `json:"date_time,omitempty"`
	IPAddress        *PiiEntityAction `json:"ip_address,omitempty"`
	URL              *PiiEntityAction `json:"url,omitempty"`
	Money            *PiiEntityAction `json:"money,omitempty"`
	CreditCard       *PiiEntityAction `json:"credit_card,omitempty"`
	Crypto           *PiiEntityAction `json:"crypto,omitempty"`
	IBANCode         *PiiEntityAction `json:"iban_code,omitempty"`
	USBankNumber     *PiiEntityAction `json:"us_bank_number,omitempty"`
	NIF              *PiiEntityAction `json:"nif,omitempty"`
	AUABN            *PiiEntityAction `json:"au_abn,omitempty"`
	AUACN            *PiiEntityAction `json:"au_acn,omitempty"`
	AUTFN            *PiiEntityAction `json:"au_tfn,omitempty"`
	MedicalLicense   *PiiEntityAction `json:"medical_license,omitempty"`
	UKNHS            *PiiEntityAction `json:"uk_nhs,omitempty"`
	AUMedicare       *PiiEntityAction `json:"au_medicare,omitempty"`
	USDriversLicense *PiiEntityAction `json:"us_drivers_license,omitempty"`
	USITIN           *PiiEntityAction `json:"us_itin,omitempty"`
	USPassport       *PiiEntityAction `json:"us_passport,omitempty"`
	USSSN            *PiiEntityAction `json:"us_ssn,omitempty"`
}

type SecretsDetectionOverride struct {
	Disabled                          *bool            `json:"disabled,omitempty"`
	SlackToken                        *PiiEntityAction `json:"slack_token,omitempty"`
	RSAPrivateKey                     *PiiEntityAction `json:"rsa_private_key,omitempty"`
	SSHDSAPrivateKey                  *PiiEntityAction `json:"ssh_dsa_private_key,omitempty"`
	SSHECPrivateKey                   *PiiEntityAction `json:"ssh_ec_private_key,omitempty"`
	PGPPrivateKeyBlock                *PiiEntityAction `json:"pgp_private_key_block,omitempty"`
	AmazonAWSAccessKeyID              *PiiEntityAction `json:"amazon_aws_access_key_id,omitempty"`
	AmazonAWSSecretAccessKey          *PiiEntityAction `json:"amazon_aws_secret_access_key,omitempty"`
	AmazonMWSAuthToken                *PiiEntityAction `json:"amazon_mws_auth_token,omitempty"`
	FacebookAccessToken               *PiiEntityAction `json:"facebook_access_token,omitempty"`
	GitHubAccessToken                 *PiiEntityAction `json:"github_access_token,omitempty"`
	JWTToken                          *PiiEntityAction `json:"jwt_token,omitempty"`
	GoogleAPIKey                      *PiiEntityAction `json:"google_api_key,omitempty"`
	GoogleCloudPlatformAPIKey         *PiiEntityAction `json:"google_cloud_platform_api_key,omitempty"`
	GoogleDriveAPIKey                 *PiiEntityAction `json:"google_drive_api_key,omitempty"`
	GoogleCloudPlatformServiceAccount *PiiEntityAction `json:"google_cloud_platform_service_account,omitempty"`
	GoogleGmailAPIKey                 *PiiEntityAction `json:"google_gmail_api_key,omitempty"`
	YouTubeAPIKey                     *PiiEntityAction `json:"youtube_api_key,omitempty"`
	MailchimpAPIKey                   *PiiEntityAction `json:"mailchimp_api_key,omitempty"`
	MailgunAPIKey                     *PiiEntityAction `json:"mailgun_api_key,omitempty"`
	BasicAuth                         *PiiEntityAction `json:"basic_auth,omitempty"`
	PicaticAPIKey                     *PiiEntityAction `json:"picatic_api_key,omitempty"`
	SlackWebhook                      *PiiEntityAction `json:"slack_webhook,omitempty"`
	StripeAPIKey                      *PiiEntityAction `json:"stripe_api_key,omitempty"`
	StripeRestrictedAPIKey            *PiiEntityAction `json:"stripe_restricted_api_key,omitempty"`
	SquareAccessToken                 *PiiEntityAction `json:"square_access_token,omitempty"`
	SquareOAuthSecret                 *PiiEntityAction `json:"square_oauth_secret,omitempty"`
	TwilioAPIKey                      *PiiEntityAction `json:"twilio_api_key,omitempty"`
	PangeaToken                       *PiiEntityAction `json:"pangea_token,omitempty"`
}

type Overrides struct {
	IgnoreRecipe      *bool                      `json:"ignore_recipe,omitempty"` // Bypass existing Recipe content and create an on-the-fly Recipe.
	CodeDetection     *CodeDetectionOverride     `json:"code_detection,omitempty"`
	Competitors       *CompetitorsOverride       `json:"competitors,omitempty"`
	Gibberish         *GibberishOverride         `json:"gibberish,omitempty"`
	LanguageDetection *LanguageDetectionOverride `json:"language_detection,omitempty"`
	MaliciousEntity   *MaliciousEntityOverride   `json:"malicious_entity,omitempty"`
	PiiEntity         *PiiEntityOverride         `json:"pii_entity,omitempty"`
	PromptInjection   *PromptInjectionOverride   `json:"prompt_injection,omitempty"`
	Roleplay          *RoleplayOverride          `json:"roleplay,omitempty"`
	SecretsDetection  *SecretsDetectionOverride  `json:"secrets_detection,omitempty"`
	SelfHarm          *SelfHarmOverride          `json:"selfharm,omitempty"`
	Sentiment         *SentimentOverride         `json:"sentiment,omitempty"`
	TopicDetection    *TopicDetectionOverride    `json:"topic,omitempty"`
}

type AnalyzerResponse struct {
	Analyzer   string  `json:"analyzer"`
	Confidence float64 `json:"confidence"`
}

type PromptInjectionResult struct {
	Action            string             `json:"action"`
	AnalyzerResponses []AnalyzerResponse `json:"analyzer_responses"`
}

// TODO: remove in favor of RedactEntity
type PiiEntity = RedactEntity

// TODO: remove in favor of RedactEntityResult
type PiiEntityResult struct {
	Entities []PiiEntity `json:"entities"` // Detected redaction rules.
}

type RedactEntity struct {
	Type     string `json:"type"`
	Value    string `json:"value"`
	Action   string `json:"action"` // The action taken on this Entity
	StartPos *int   `json:"start_pos,omitempty"`
}

type RedactEntityResult struct {
	Entities []RedactEntity `json:"entities"` // Detected redaction rules.
}

type MaliciousEntity struct {
	Type     string                 `json:"type"`
	Value    string                 `json:"value"`
	Action   string                 `json:"action"`
	StartPos *int                   `json:"start_pos,omitempty"`
	Raw      map[string]interface{} `json:"raw,omitempty"`
}

type MaliciousEntityResult struct {
	Entities []MaliciousEntity `json:"entities"` // Detected harmful items.
}

type SecretsEntity struct {
	Type          string `json:"type"`
	Value         string `json:"value"`
	Action        string `json:"action"` // The action taken on this Entity
	StartPos      *int   `json:"start_pos,omitempty"`
	RedactedValue string `json:"redacted_value,omitempty"`
}

type SecretsEntityResult struct {
	Entities []SecretsEntity `json:"entities"`
}

type LanguageDetectionResult struct {
	Language string `json:"language"`
	Action   string `json:"action"`
}

type Topic struct {
	Topic      string  `json:"topic"`
	Confidence float64 `json:"confidence"`
}

type TopicDetectionResult struct {
	Action string  `json:"action"` // The action taken by this Detector
	Topics []Topic `json:"topics"` // List of topics detected
}

type CodeDetectionResult struct {
	Language string `json:"language"`
	Action   string `json:"action"`
}

type SingleEntityResult struct {
	Action   string   `json:"action"`   // The action taken by this Detector
	Entities []string `json:"entities"` // Detected entities
}

type Classification struct {
	Category   string  `json:"category"`
	Confidence float64 `json:"confidence"`
}

type ClassificationResult struct {
	Action          string           `json:"action"` // The action taken by this Detector
	Classifications []Classification `json:"classifications"`
}

type TextGuardDetector[T any] struct {
	Detected bool `json:"detected"`
	Data     *T   `json:"data,omitempty"`
}

type TextGuardDetectors struct {
	PromptInjection      *TextGuardDetector[PromptInjectionResult]   `json:"prompt_injection,omitempty"`
	Gibberish            *TextGuardDetector[ClassificationResult]    `json:"gibberish,omitempty"`
	Sentiment            *TextGuardDetector[ClassificationResult]    `json:"sentiment,omitempty"`
	SelfHarm             *TextGuardDetector[ClassificationResult]    `json:"selfharm,omitempty"`
	PiiEntity            *TextGuardDetector[PiiEntityResult]         `json:"pii_entity,omitempty"`
	MaliciousEntity      *TextGuardDetector[MaliciousEntityResult]   `json:"malicious_entity,omitempty"`
	CustomEntity         *TextGuardDetector[RedactEntityResult]      `json:"custom_entity,omitempty"`
	SecretsDetection     *TextGuardDetector[SecretsEntityResult]     `json:"secrets_detection,omitempty"`
	Competitors          *TextGuardDetector[SingleEntityResult]      `json:"competitors,omitempty"`
	ProfanityAndToxicity *TextGuardDetector[ClassificationResult]    `json:"profanity_and_toxicity,omitempty"`
	LanguageDetection    *TextGuardDetector[LanguageDetectionResult] `json:"language_detection,omitempty"`
	TopicDetection       *TextGuardDetector[TopicDetectionResult]    `json:"topic,omitempty"`
	CodeDetection        *TextGuardDetector[CodeDetectionResult]     `json:"code_detection,omitempty"`
}

// LogFields are additional fields to include in activity log
type LogFields struct {
	Citations string `json:"citations,omitempty"`  // Origin or source application of the event
	ExtraInfo string `json:"extra_info,omitempty"` // Stores supplementary details related to the event
	Model     string `json:"model,omitempty"`      // Model used to perform the event
	Source    string `json:"source,omitempty"`     // IP address of user or app or agent
	Tools     string `json:"tools,omitempty"`      // Tools used to perform the event
}

type TextGuardRequest struct {
	pangea.BaseRequest

	Text      string     `json:"text,omitempty"`     // Text to be scanned by AI Guard for PII, sensitive data, malicious content, and other data types defined by the configuration. Supports processing up to 10KB of text.
	Messages  any        `json:"messages,omitempty"` // Structured messages data to be scanned by AI Guard for PII, sensitive data, malicious content, and other data types defined by the configuration. Supports processing up to 10KB of JSON text.
	Recipe    string     `json:"recipe,omitempty"`   // Recipe key of a configuration of data types and settings defined in the Pangea User Console. It specifies the rules that are to be applied to the text, such as defang malicious URLs.
	Debug     bool       `json:"debug,omitempty"`    // Setting this value to true will provide a detailed analysis of the text data
	Overrides *Overrides `json:"overrides,omitempty"`
	LogFields LogFields  `json:"log_fields,omitempty"` // Additional fields to include in activity log
}

type TextGuardResult struct {
	Detectors      TextGuardDetectors `json:"detectors"`       // Result of the recipe analyzing and input prompt.
	PromptText     string             `json:"prompt_text"`     // Updated prompt text, if applicable.
	PromptMessages any                `json:"prompt_messages"` // Updated prompt messages, if applicable.
	Blocked        bool               `json:"blocked"`         // Whether or not the prompt triggered a block detection.
	Recipe         string             `json:"recipe"`          // The Recipe that was used.
}

type TextContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ImageContent struct {
	Type     string `json:"type"`
	ImageSrc string `json:"image_src" validate:"regexp=^(data:image\\/(jpeg|png|gif);base64,|https?:\\/\\/).*"`
}

type MultimodalContentInner struct {
	TextContent  *TextContent
	ImageContent *ImageContent
}

func (src MultimodalContentInner) MarshalJSON() ([]byte, error) {
	if src.TextContent != nil {
		return json.Marshal(&src.TextContent)
	}

	if src.ImageContent != nil {
		return json.Marshal(&src.ImageContent)
	}

	return nil, nil
}

type MultimodalContent struct {
	OfString         *string
	OfArrayOfContent []MultimodalContentInner
}

func (src MultimodalContent) MarshalJSON() ([]byte, error) {
	if src.OfString != nil {
		return json.Marshal(&src.OfString)
	}

	if src.OfArrayOfContent != nil {
		return json.Marshal(&src.OfArrayOfContent)
	}

	return nil, nil
}

type MultimodalMessage struct {
	Role    string            `json:"role"`
	Content MultimodalContent `json:"content"`
}

type GuardRequest struct {
	pangea.BaseRequest

	Messages           []MultimodalMessage `json:"messages,omitempty"` // Prompt content and role array in JSON format. The `content` is the multimodal text or image input that will be analyzed.
	Recipe             string              `json:"recipe,omitempty"`   // Recipe key of a configuration of data types and settings defined in the Pangea User Console. It specifies the rules that are to be applied to the text, such as defang malicious URLs.
	Debug              bool                `json:"debug,omitempty"`    // Setting this value to true will provide a detailed analysis of the text data
	Overrides          Overrides           `json:"overrides,omitempty"`
	AppName            string              `json:"app_name,omitempty"`             // Name of source application.
	LlmProvider        string              `json:"llm_provider,omitempty"`         // Underlying LLM.  Example: 'OpenAI'.
	Model              string              `json:"model,omitempty"`                // Model used to perform the event. Example: 'gpt'.
	ModelVersion       string              `json:"model_version,omitempty"`        // Model version used to perform the event. Example: '3.5'.
	RequestTokenCount  int32               `json:"request_token_count,omitempty"`  // Number of tokens in the request.
	ResponseTokenCount int32               `json:"response_token_count,omitempty"` // Number of tokens in the response.
	SourceIp           string              `json:"source_ip,omitempty"`            // IP address of user or app or agent.
	SourceLocation     string              `json:"source_location,omitempty"`      // Location of user or app or agent.
	TenantId           string              `json:"tenant_id,omitempty"`            // For gateway-like integrations with multi-tenant support.
	SensorMode         string              `json:"sensor_mode,omitempty"`          // (AIDR) sensor mode.
	Context            map[string]any      `json:"context,omitempty"`              // (AIDR) Logging schema.
}

type GuardResult struct {
	PromptMessages map[string]any     `json:"prompt_messages,omitempty"` // Updated structured prompt.
	Blocked        bool               `json:"blocked,omitempty"`         // Whether or not the prompt triggered a block detection.
	Recipe         string             `json:"recipe,omitempty"`          // The Recipe that was used.
	Detectors      TextGuardDetectors `json:"detectors"`
	FpeContext     string             `json:"fpe_context,omitempty"` // If an FPE redaction method returned results, this will be the context passed to unredact.
}

type GetServiceConfigParams struct {
	pangea.BaseRequest

	Id string `json:"id"`
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

type AuditDataActivityConfigAreas struct {
	TextGuard bool `json:"text_guard"`
}

type AuditDataActivityConfig struct {
	Enabled              bool                         `json:"enabled"`
	AuditServiceConfigId string                       `json:"audit_service_config_id"`
	Areas                AuditDataActivityConfigAreas `json:"areas"`
}

type ConnectionsConfigPromptGuard struct {
	Enabled             *bool    `json:"enabled,omitempty"`
	ConfigId            *string  `json:"config_id,omitempty"`
	ConfidenceThreshold *float32 `json:"confidence_threshold,omitempty"`
}

type ConnectionsConfigIpIntel struct {
	Enabled            *bool    `json:"enabled,omitempty"`
	ConfigId           *string  `json:"config_id,omitempty"`
	ReputationProvider *string  `json:"reputation_provider,omitempty"`
	RiskThreshold      *float32 `json:"risk_threshold,omitempty"`
}

type ConnectionsConfigUserIntel struct {
	Enabled        *bool   `json:"enabled,omitempty"`
	ConfigId       *string `json:"config_id,omitempty"`
	BreachProvider *string `json:"breach_provider,omitempty"`
}
type ConnectionsConfigUrlIntel struct {
	Enabled            *bool    `json:"enabled,omitempty"`
	ConfigId           *string  `json:"config_id,omitempty"`
	ReputationProvider *string  `json:"reputation_provider,omitempty"`
	RiskThreshold      *float32 `json:"risk_threshold,omitempty"`
}

type ConnectionsConfigFileScan struct {
	Enabled       *bool    `json:"enabled,omitempty"`
	ConfigId      *string  `json:"config_id,omitempty"`
	ScanProvider  *string  `json:"scan_provider,omitempty"`
	RiskThreshold *float32 `json:"risk_threshold,omitempty"`
}

type ConnectionsConfigRedact struct {
	Enabled  *bool   `json:"enabled,omitempty"`
	ConfigId *string `json:"config_id,omitempty"`
}

type ConnectionsConfigVault struct {
	ConfigId *string `json:"config_id,omitempty"`
}

type ConnectionsConfigLingua struct {
	Enabled *bool `json:"enabled,omitempty"`
}

type ConnectionsConfigCode struct {
	Enabled *bool `json:"enabled,omitempty"`
}

type ConnectionsConfig struct {
	PromptGuard *ConnectionsConfigPromptGuard `json:"prompt_guard,omitempty"`
	IpIntel     *ConnectionsConfigIpIntel     `json:"ip_intel,omitempty"`
	UserIntel   *ConnectionsConfigUserIntel   `json:"user_intel,omitempty"`
	UrlIntel    *ConnectionsConfigUrlIntel    `json:"url_intel,omitempty"`
	DomainIntel *ConnectionsConfigUrlIntel    `json:"domain_intel,omitempty"`
	FileScan    *ConnectionsConfigFileScan    `json:"file_scan,omitempty"`
	Redact      *ConnectionsConfigRedact      `json:"redact,omitempty"`
	Vault       *ConnectionsConfigVault       `json:"vault,omitempty"`
	Lingua      *ConnectionsConfigLingua      `json:"lingua,omitempty"`
	Code        *ConnectionsConfigCode        `json:"code,omitempty"`
}

type ServiceConfig struct {
	Id                *string                  `json:"id,omitempty"`
	Name              *string                  `json:"name,omitempty"`
	AuditDataActivity *AuditDataActivityConfig `json:"audit_data_activity,omitempty"`
	Connections       *ConnectionsConfig       `json:"connections,omitempty"`
	Recipes           map[string]any           `json:"recipes,omitempty"`
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

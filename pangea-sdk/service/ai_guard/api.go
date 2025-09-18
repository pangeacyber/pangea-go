package ai_guard

import (
	"context"
	"errors"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/internal/request"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
)

// @summary Guard LLM input and output
//
// @description Analyze and redact content to avoid manipulation of the model, addition of malicious content, and other undesirable data transfers.
//
// @operationId ai_guard_post_v1_guard
//
// @example
//
//	input := &ai_guard.GuardRequest{Text: "hello world"}
//	response, err := client.Guard(ctx, input)
func (e *aiGuard) Guard(ctx context.Context, input GuardRequest) (*pangea.PangeaResponse[GuardResult], error) {
	return request.DoPost(ctx, e.Client, "v1/guard", &input, &GuardResult{})
}

// @summary Guard LLM input and output
//
// @description Analyze and redact content to avoid manipulation of the model, addition of malicious content, and other undesirable data transfers.
//
// @operationId ai_guard_post_v1_guard_async
//
// @example
//
//	input := &ai_guard.GuardRequest{Text: "hello world"}
//	response, err := client.Guard(ctx, input)
func (e *aiGuard) GuardAsync(ctx context.Context, input GuardRequest) (*pangea.PangeaResponse[GuardResult], error) {
	response, err := request.DoPostNoQueue(ctx, e.Client, "v1/guard_async", &input, &GuardResult{})
	if err != nil {
		acceptedErr, ok := err.(*pangea.AcceptedError)
		if ok {
			return &pangea.PangeaResponse[GuardResult]{
				AcceptedResult: &acceptedErr.AcceptedResult,
				Response:       acceptedErr.Response,
				Result:         nil,
			}, nil
		}

		return nil, err
	}
	return response, nil
}

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

// GuardTextWithRelevantContent sends only relevant messages to the AI Guard
// service. It implements the logic to filter messages, call the API, and then
// patch the results back.
func (e *aiGuard) GuardTextWithRelevantContent(ctx context.Context, input *TextGuardRequest) (*pangea.PangeaResponse[TextGuardResult], error) {
	if input.Text != "" || input.Messages == nil {
		return e.GuardText(ctx, input)
	}

	originalMessages := input.Messages
	relevantMessages, originalIndices := getRelevantContent(originalMessages)

	relevantRequest := *input
	relevantRequest.Messages = relevantMessages

	resp, err := e.GuardText(ctx, &relevantRequest)
	if err != nil {
		return nil, err
	}

	if resp != nil && resp.Result != nil && resp.Result.PromptMessages != nil {
		transformedMessages := resp.Result.PromptMessages
		resp.Result.PromptMessages = patchMessages(originalMessages, originalIndices, transformedMessages)
	}

	return resp, nil
}

func getRelevantContent(messages []PromptMessage) ([]PromptMessage, []int) {
	if len(messages) == 0 {
		return []PromptMessage{}, []int{}
	}

	var systemMessages []PromptMessage
	var systemIndices []int
	for i, msg := range messages {
		if msg.Role == "system" {
			systemMessages = append(systemMessages, msg)
			systemIndices = append(systemIndices, i)
		}
	}

	if messages[len(messages)-1].Role == "assistant" {
		relevantMessages := append(systemMessages, messages[len(messages)-1])
		relevantIndices := append(systemIndices, len(messages)-1)
		return relevantMessages, relevantIndices
	}

	lastAssistantIndex := -1
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == "assistant" {
			lastAssistantIndex = i
			break
		}
	}

	var relevantMessages []PromptMessage
	var indices []int
	for i, msg := range messages {
		if msg.Role == "system" || i > lastAssistantIndex {
			relevantMessages = append(relevantMessages, msg)
			indices = append(indices, i)
		}
	}

	return relevantMessages, indices
}

func patchMessages(original []PromptMessage, originalIndices []int, transformed []PromptMessage) []PromptMessage {
	if len(original) == len(transformed) {
		return transformed
	}

	result := make([]PromptMessage, len(original))
	copy(result, original)

	for i, transformedMsg := range transformed {
		originalIndex := originalIndices[i]
		result[originalIndex] = transformedMsg
	}

	return result
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

type PromptMessage struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type TextGuardRequest struct {
	pangea.BaseRequest

	Text      string          `json:"text,omitzero"`       // Text to be scanned by AI Guard for PII, sensitive data, malicious content, and other data types defined by the configuration. Supports processing up to 10KB of text.
	Messages  []PromptMessage `json:"messages,omitzero"`   // Structured messages data to be scanned by AI Guard for PII, sensitive data, malicious content, and other data types defined by the configuration. Supports processing up to 10KB of JSON text.
	Recipe    string          `json:"recipe,omitzero"`     // Recipe key of a configuration of data types and settings defined in the Pangea User Console. It specifies the rules that are to be applied to the text, such as defang malicious URLs.
	Debug     bool            `json:"debug,omitzero"`      // Setting this value to true will provide a detailed analysis of the text data
	Overrides Overrides       `json:"overrides,omitzero"`  // Overrides flags.
	LogFields LogFields       `json:"log_fields,omitzero"` // Additional fields to include in activity log
}

type TextGuardResult struct {
	Detectors      TextGuardDetectors `json:"detectors"`                   // Result of the recipe analyzing and input prompt.
	AccessRules    any                `json:"access_rules"`                // Result of the recipe evaluating configured rules
	Blocked        bool               `json:"blocked"`                     // Whether or not the prompt triggered a block detection.
	FpeContext     string             `json:"fpe_context" format:"base64"` // If an FPE redaction method returned results, this will be the context passed to unredact.
	PromptMessages []PromptMessage    `json:"prompt_messages"`             // Updated structured prompt, if applicable.
	PromptText     string             `json:"prompt_text"`                 // Updated prompt text, if applicable.
	Recipe         string             `json:"recipe"`                      // The Recipe that was used.
	Transformed    bool               `json:"transformed"`                 // Whether or not the original input was transformed.
}

// (AIDR) Event Type.
type GuardEventType string

const (
	GuardEventTypeInput       GuardEventType = "input"
	GuardEventTypeOutput      GuardEventType = "output"
	GuardEventTypeToolInput   GuardEventType = "tool_input"
	GuardEventTypeToolOutput  GuardEventType = "tool_output"
	GuardEventTypeToolListing GuardEventType = "tool_listing"
)

// The properties ServerName, Tools are required.
type GuardExtraInfoMcpToolParam struct {
	// MCP server name
	ServerName string   `json:"server_name"`
	Tools      []string `json:"tools,omitzero"`
}

// (AIDR) Logging schema.
type GuardExtraInfoParam struct {
	// The group of subject actor.
	ActorGroup string `json:"actor_group,omitzero"`
	// Name of subject actor/service account.
	ActorName string `json:"actor_name,omitzero"`
	// The group of source application/agent.
	AppGroup string `json:"app_group,omitzero"`
	// Name of source application/agent.
	AppName string `json:"app_name,omitzero"`
	// Version of the source application/agent.
	AppVersion string `json:"app_version,omitzero"`
	// Geographic region or data center.
	SourceRegion string `json:"source_region,omitzero"`
	// Sub tenant of the user or organization
	SubTenant string `json:"sub_tenant,omitzero"`
	// Each item groups tools for a given MCP server.
	McpTools    []GuardExtraInfoMcpToolParam `json:"mcp_tools,omitzero"`
	ExtraFields map[string]any               `json:"-"`
}

type CodeDetectionAction string

const (
	CodeDetectionActionReport CodeDetectionAction = "report"
	CodeDetectionActionBlock  CodeDetectionAction = "block"
)

type CompetitorsAction string

const (
	CompetitorsActionReport CompetitorsAction = "report"
	CompetitorsActionBlock  CompetitorsAction = "block"
)

type LanguageDetectionItemsAction string

const (
	LanguageDetectionItemsActionEmpty  LanguageDetectionItemsAction = ""
	LanguageDetectionItemsActionReport LanguageDetectionItemsAction = "report"
	LanguageDetectionItemsActionAllow  LanguageDetectionItemsAction = "allow"
	LanguageDetectionItemsActionBlock  LanguageDetectionItemsAction = "block"
)

type TopicDetectionItemsAction string

const (
	TopicDetectionItemsActionEmpty  TopicDetectionItemsAction = ""
	TopicDetectionItemsActionReport TopicDetectionItemsAction = "report"
	TopicDetectionItemsActionBlock  TopicDetectionItemsAction = "block"
)

type GuardOverridesCodeParam struct {
	Disabled  bool    `json:"disabled,omitzero"`
	Threshold float64 `json:"threshold,omitzero"`
	// Any of "report", "block".
	Action CodeDetectionAction `json:"action,omitzero"`
}

type GuardOverridesCompetitorsParam struct {
	Disabled bool `json:"disabled,omitzero"`
	// Any of "report", "block".
	Action CompetitorsAction `json:"action,omitzero"`
}

type GuardOverridesConfidentialAndPiiEntityParam struct {
	Disabled bool `json:"disabled,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	AuAbn PiiEntityAction `json:"au_abn,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	AuAcn PiiEntityAction `json:"au_acn,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	AuMedicare PiiEntityAction `json:"au_medicare,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	AuTfn PiiEntityAction `json:"au_tfn,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	CreditCard PiiEntityAction `json:"credit_card,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	Crypto PiiEntityAction `json:"crypto,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	DateTime PiiEntityAction `json:"date_time,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	EmailAddress PiiEntityAction `json:"email_address,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	FinNric PiiEntityAction `json:"fin/nric,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	IbanCode PiiEntityAction `json:"iban_code,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	IPAddress PiiEntityAction `json:"ip_address,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	Location PiiEntityAction `json:"location,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	MedicalLicense PiiEntityAction `json:"medical_license,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	Money PiiEntityAction `json:"money,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	Nif PiiEntityAction `json:"nif,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	Nrp PiiEntityAction `json:"nrp,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	Person PiiEntityAction `json:"person,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	PhoneNumber PiiEntityAction `json:"phone_number,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	UkNhs PiiEntityAction `json:"uk_nhs,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	URL PiiEntityAction `json:"url,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	UsBankNumber PiiEntityAction `json:"us_bank_number,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	UsDriversLicense PiiEntityAction `json:"us_drivers_license,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	UsItin PiiEntityAction `json:"us_itin,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	UsPassport PiiEntityAction `json:"us_passport,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	UsSsn PiiEntityAction `json:"us_ssn,omitzero"`
}

type GuardOverridesImageParam struct {
	Disabled  bool    `json:"disabled,omitzero"`
	Threshold float64 `json:"threshold,omitzero"`
	// Any of "", "report", "block".
	Action string   `json:"action,omitzero"`
	Topics []string `json:"topics,omitzero"`
}

type LanguageDetectionItemsParam struct {
	Disabled  bool    `json:"disabled,omitzero"`
	Threshold float64 `json:"threshold,omitzero"`
	// Any of "", "report", "allow", "block".
	Action    LanguageDetectionItemsAction `json:"action,omitzero"`
	Languages []string                     `json:"languages,omitzero"`
}

type GuardOverridesMaliciousEntityParam struct {
	Disabled bool `json:"disabled,omitzero"`
	// Any of "report", "defang", "disabled", "block".
	Domain MaliciousEntityAction `json:"domain,omitzero"`
	// Any of "report", "defang", "disabled", "block".
	IPAddress MaliciousEntityAction `json:"ip_address,omitzero"`
	// Any of "report", "defang", "disabled", "block".
	URL MaliciousEntityAction `json:"url,omitzero"`
}

type GuardOverridesMaliciousPromptParam struct {
	Disabled bool `json:"disabled,omitzero"`
	// Any of "report", "block".
	Action PromptInjectionAction `json:"action,omitzero"`
}

type GuardOverridesSecretAndKeyEntityParam struct {
	Disabled bool `json:"disabled,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	AmazonAwsAccessKeyID PiiEntityAction `json:"amazon_aws_access_key_id,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	AmazonAwsSecretAccessKey PiiEntityAction `json:"amazon_aws_secret_access_key,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	AmazonMwsAuthToken PiiEntityAction `json:"amazon_mws_auth_token,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	BasicAuth PiiEntityAction `json:"basic_auth,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	FacebookAccessToken PiiEntityAction `json:"facebook_access_token,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	GitHubAccessToken PiiEntityAction `json:"github_access_token,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	GoogleAPIKey PiiEntityAction `json:"google_api_key,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	GoogleCloudPlatformAPIKey PiiEntityAction `json:"google_cloud_platform_api_key,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	GoogleCloudPlatformServiceAccount PiiEntityAction `json:"google_cloud_platform_service_account,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	GoogleDriveAPIKey PiiEntityAction `json:"google_drive_api_key,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	GoogleGmailAPIKey PiiEntityAction `json:"google_gmail_api_key,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	JwtToken PiiEntityAction `json:"jwt_token,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	MailchimpAPIKey PiiEntityAction `json:"mailchimp_api_key,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	MailgunAPIKey PiiEntityAction `json:"mailgun_api_key,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	PangeaToken PiiEntityAction `json:"pangea_token,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	PgpPrivateKeyBlock PiiEntityAction `json:"pgp_private_key_block,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	PicaticAPIKey PiiEntityAction `json:"picatic_api_key,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	RsaPrivateKey PiiEntityAction `json:"rsa_private_key,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	SlackToken PiiEntityAction `json:"slack_token,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	SlackWebhook PiiEntityAction `json:"slack_webhook,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	SquareAccessToken PiiEntityAction `json:"square_access_token,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	SquareOAuthSecret PiiEntityAction `json:"square_oauth_secret,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	SSHDsaPrivateKey PiiEntityAction `json:"ssh_dsa_private_key,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	SSHEcPrivateKey PiiEntityAction `json:"ssh_ec_private_key,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	StripeAPIKey PiiEntityAction `json:"stripe_api_key,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	StripeRestrictedAPIKey PiiEntityAction `json:"stripe_restricted_api_key,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	TwilioAPIKey PiiEntityAction `json:"twilio_api_key,omitzero"`
	// Any of "disabled", "report", "block", "mask", "partial_masking", "replacement",
	// "hash", "fpe".
	YoutubeAPIKey PiiEntityAction `json:"youtube_api_key,omitzero"`
}

type TopicDetectionItemsParam struct {
	Disabled  bool    `json:"disabled,omitzero"`
	Threshold float64 `json:"threshold,omitzero"`
	// Any of "", "report", "block".
	Action TopicDetectionItemsAction `json:"action,omitzero"`
	Topics []string                  `json:"topics,omitzero"`
}

// Overrides flags. Note: This parameter has no effect when the request is made by
// AIDR
type GuardOverridesParam struct {
	// Bypass existing Recipe content and create an on-the-fly Recipe.
	IgnoreRecipe             bool                                        `json:"ignore_recipe,omitzero"`
	Code                     GuardOverridesCodeParam                     `json:"code,omitzero"`
	Competitors              GuardOverridesCompetitorsParam              `json:"competitors,omitzero"`
	ConfidentialAndPiiEntity GuardOverridesConfidentialAndPiiEntityParam `json:"confidential_and_pii_entity,omitzero"`
	Image                    GuardOverridesImageParam                    `json:"image,omitzero"`
	Language                 LanguageDetectionItemsParam                 `json:"language,omitzero"`
	MaliciousEntity          GuardOverridesMaliciousEntityParam          `json:"malicious_entity,omitzero"`
	MaliciousPrompt          GuardOverridesMaliciousPromptParam          `json:"malicious_prompt,omitzero"`
	SecretAndKeyEntity       GuardOverridesSecretAndKeyEntityParam       `json:"secret_and_key_entity,omitzero"`
	Topic                    TopicDetectionItemsParam                    `json:"topic,omitzero"`
}

// The property Input is required.
type GuardRequest struct {
	pangea.BaseRequest

	// 'messages' (required) contains Prompt content and role array in JSON format. The
	// `content` is the multimodal text or image input that will be analyzed.
	// Additional properties such as 'tools' may be provided for analysis.
	Input any `json:"input,omitzero"`
	// User/Service account id/service account
	ActorID string `json:"actor_id,omitzero"`
	// Id of source application/agent
	AppID string `json:"app_id,omitzero"`
	// (AIDR) collector instance id.
	CollectorInstanceID string `json:"collector_instance_id,omitzero"`
	// Provide input and output token count.
	CountTokens bool `json:"count_tokens,omitzero"`
	// Setting this value to true will provide a detailed analysis of the text data
	Debug bool `json:"debug,omitzero"`
	// Underlying LLM. Example: 'OpenAI'.
	LlmProvider string `json:"llm_provider,omitzero"`
	// Model used to perform the event. Example: 'gpt'.
	Model string `json:"model,omitzero"`
	// Model version used to perform the event. Example: '3.5'.
	ModelVersion string `json:"model_version,omitzero"`
	// Recipe key of a configuration of data types and settings defined in the Pangea
	// User Console. It specifies the rules that are to be applied to the text, such as
	// defang malicious URLs. Note: This parameter has no effect when the request is
	// made by AIDR
	Recipe string `json:"recipe,omitzero"`
	// Number of tokens in the request.
	RequestTokenCount int64 `json:"request_token_count,omitzero"`
	// Number of tokens in the response.
	ResponseTokenCount int64 `json:"response_token_count,omitzero"`
	// IP address of user or app or agent.
	SourceIP string `json:"source_ip,omitzero"`
	// Location of user or app or agent.
	SourceLocation string `json:"source_location,omitzero"`
	// For gateway-like integrations with multi-tenant support.
	TenantID string `json:"tenant_id,omitzero"`
	// (AIDR) Event Type.
	//
	// Any of "input", "output", "tool_input", "tool_output", "tool_listing".
	EventType GuardEventType `json:"event_type,omitzero"`
	// (AIDR) Logging schema.
	ExtraInfo GuardExtraInfoParam `json:"extra_info,omitzero"`
	// Overrides flags. Note: This parameter has no effect when the request is made by
	// AIDR
	Overrides GuardOverridesParam `json:"overrides,omitzero"`
}

type HardeningResult struct {
	// The action taken by this Detector
	Action string `json:"action"`
	// Descriptive information about the hardening detector execution
	Message string `json:"message"`
	// Number of tokens counted in the last user prompt
	TokenCount float64 `json:"token_count"`
}

type LanguageResult struct {
	// The action taken by this Detector
	Action   string `json:"action"`
	Language string `json:"language"`
}

type TopicResultTopic struct {
	Confidence float64 `json:"confidence"`
	Topic      string  `json:"topic"`
}

type TopicResult struct {
	// The action taken by this Detector
	Action string `json:"action"`
	// List of topics detected
	Topics []TopicResultTopic `json:"topics"`
}

type GuardResultDetectorsCode struct {
	// Details about the detected code.
	Data LanguageResult `json:"data"`
	// Whether or not the Code was detected.
	Detected bool `json:"detected"`
}

type GuardResultDetectorsCompetitors struct {
	// Details about the detected entities.
	Data SingleEntityResult `json:"data"`
	// Whether or not the Competitors were detected.
	Detected bool `json:"detected"`
}

type GuardResultDetectorsConfidentialAndPiiEntity struct {
	// Details about the detected entities.
	Data RedactEntityResult `json:"data"`
	// Whether or not the PII Entities were detected.
	Detected bool `json:"detected"`
}

type GuardResultDetectorsCustomEntity struct {
	// Details about the detected entities.
	Data RedactEntityResult `json:"data"`
	// Whether or not the Custom Entities were detected.
	Detected bool `json:"detected"`
}

type GuardResultDetectorsLanguage struct {
	// Details about the detected languages.
	Data LanguageResult `json:"data"`
	// Whether or not the Languages were detected.
	Detected bool `json:"detected"`
}

type GuardResultDetectorsMaliciousEntity struct {
	// Details about the detected entities.
	Data MaliciousEntityResult `json:"data"`
	// Whether or not the Malicious Entities were detected.
	Detected bool `json:"detected"`
}

type GuardResultDetectorsMaliciousPrompt struct {
	// Details about the analyzers.
	Data PromptInjectionResult `json:"data"`
	// Whether or not the Malicious Prompt was detected.
	Detected bool `json:"detected"`
}

type GuardResultDetectorsPromptHardening struct {
	// Details about the detected languages.
	Data HardeningResult `json:"data"`
}

type GuardResultDetectorsSecretAndKeyEntity struct {
	// Details about the detected entities.
	Data RedactEntityResult `json:"data"`
	// Whether or not the Secret Entities were detected.
	Detected bool `json:"detected"`
}

type GuardResultDetectorsTopic struct {
	// Details about the detected topics.
	Data TopicResult `json:"data"`
	// Whether or not the Topics were detected.
	Detected bool `json:"detected"`
}

// Result of the recipe analyzing and input prompt.
type GuardResultDetectors struct {
	Code                     GuardResultDetectorsCode                     `json:"code"`
	Competitors              GuardResultDetectorsCompetitors              `json:"competitors"`
	ConfidentialAndPiiEntity GuardResultDetectorsConfidentialAndPiiEntity `json:"confidential_and_pii_entity"`
	CustomEntity             GuardResultDetectorsCustomEntity             `json:"custom_entity"`
	Language                 GuardResultDetectorsLanguage                 `json:"language"`
	MaliciousEntity          GuardResultDetectorsMaliciousEntity          `json:"malicious_entity"`
	MaliciousPrompt          GuardResultDetectorsMaliciousPrompt          `json:"malicious_prompt"`
	PromptHardening          GuardResultDetectorsPromptHardening          `json:"prompt_hardening"`
	SecretAndKeyEntity       GuardResultDetectorsSecretAndKeyEntity       `json:"secret_and_key_entity"`
	Topic                    GuardResultDetectorsTopic                    `json:"topic"`
}

type GuardResult struct {
	// Result of the recipe analyzing and input prompt.
	Detectors GuardResultDetectors `json:"detectors"`
	// Result of the recipe evaluating configured rules
	AccessRules any `json:"access_rules"`
	// Whether or not the prompt triggered a block detection.
	Blocked bool `json:"blocked"`
	// If an FPE redaction method returned results, this will be the context passed to
	// unredact.
	FpeContext string `json:"fpe_context" format:"base64"`
	// Number of tokens counted in the input
	InputTokenCount float64 `json:"input_token_count"`
	// Updated structured prompt.
	Output any `json:"output"`
	// Number of tokens counted in the output
	OutputTokenCount float64 `json:"output_token_count"`
	// The Recipe that was used.
	Recipe string `json:"recipe"`
	// Whether or not the original input was transformed.
	Transformed bool `json:"transformed"`
}

package ai_guard

import (
	"context"
	"errors"

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

type TopicDetectionOverride struct {
	Disabled *bool    `json:"disabled,omitempty"`
	Block    []string `json:"block,omitempty"`
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
	TopicDetection    *TopicDetectionOverride    `json:"topic_detection,omitempty"`
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
	TopicDetection       *TextGuardDetector[TopicDetectionResult]    `json:"topic_detection,omitempty"`
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

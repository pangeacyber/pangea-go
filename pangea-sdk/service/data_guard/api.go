package data_guard

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/internal/request"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/service/user_intel"
)

// @summary Text guard (Beta)
//
// @description Guard text.
//
// @operationId data_guard_post_v1_text_guard
//
// @example
//
//	input := &data_guard.TextGuardRequest{Text: "hello world"}
//	response, err := client.GuardText(ctx, input)
func (e *dataGuard) GuardText(ctx context.Context, input *TextGuardRequest) (*pangea.PangeaResponse[TextGuardResult], error) {
	return request.DoPost(ctx, e.Client, "v1/text/guard", input, &TextGuardResult{})
}

// @summary File guard (Beta)
//
// @description Guard file URLs.
//
// @operationId data_guard_post_v1_file_guard
//
// @example
//
//	input := &data_guard.FileGuardRequest{FileUrl: "https://example.org/file.txt"}
//	response, err := client.GuardFile(ctx, input)
func (e *dataGuard) GuardFile(ctx context.Context, input *FileGuardRequest) (*pangea.PangeaResponse[struct{}], error) {
	var result struct{}
	return request.DoPost(ctx, e.Client, "v1/file/guard", input, &result)
}

type TextGuardArtifact struct {
	Defanged bool   `json:"defanged"`
	End      int    `json:"end"`
	Start    int    `json:"start"`
	Type     string `json:"type"`
	Value    string `json:"value"`
	Verdict  string `json:"verdict,omitempty"`
}

type TextGuardSecurityIssues struct {
	CompromisedEmailAddresses int `json:"compromised_email_addresses"`
	MaliciousDomainCount      int `json:"malicious_domain_count"`
	MaliciousIPCount          int `json:"malicious_ip_count"`
	MaliciousURLCount         int `json:"malicious_url_count"`
	RedactedItemCount         int `json:"redacted_item_count"`
}

type TextGuardFindings struct {
	ArtifactCount  int                     `json:"artifact_count"`
	MaliciousCount int                     `json:"malicious_count"`
	SecurityIssues TextGuardSecurityIssues `json:"security_issues"`
}

type RedactRecognizerResult struct {
	FieldType string `json:"field_type"` // The entity name.
	Score     int    `json:"score"`      // The certainty score that the entity matches this specific snippet.
	Text      string `json:"text"`       // The text snippet that matched.
	Start     int    `json:"start"`      // The starting index of a snippet.
	End       int    `json:"end"`        // The ending index of a snippet.
	Redacted  bool   `json:"redacted"`   // Indicates if this rule was used to anonymize a text snippet.
}

type RedactReport struct {
	Count             int                      `json:"count"`
	RecognizerResults []RedactRecognizerResult `json:"recognizer_results"`
}

type IntelResults struct {
	Category []string `json:"category"` // The categories that apply to this indicator as determined by the provider.
	Score    int      `json:"score"`    // The score, given by the Pangea service, for the indicator.
	Verdict  string   `json:"verdict"`  // The verdict for the indicator.
}

type TextGuardReport struct {
	DomainIntel *IntelResults                `json:"domain_intel,omitempty"`
	IPIntel     *IntelResults                `json:"ip_intel,omitempty"`
	Redact      RedactReport                 `json:"redact"`
	URLIntel    *IntelResults                `json:"url_intel,omitempty"`
	UserIntel   *user_intel.UserBreachedData `json:"user_intel,omitempty"`
}

type TextGuardRequest struct {
	pangea.BaseRequest

	Text   string `json:"text"`
	Recipe string `json:"recipe,omitempty"`
	Debug  bool   `json:"debug,omitempty"`
}

type TextGuardResult struct {
	Artifacts      []TextGuardArtifact `json:"artifacts,omitempty"`
	Findings       TextGuardFindings   `json:"findings"`
	RedactedPrompt string              `json:"redacted_prompt"`

	// `debug=true` only.
	Report *TextGuardReport `json:"report,omitempty"`
}

type FileGuardRequest struct {
	pangea.BaseRequest

	FileUrl string `json:"file_url"`
}

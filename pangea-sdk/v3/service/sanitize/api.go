package sanitize

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/internal/request"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
)

func (e *sanitize) Sanitize(ctx context.Context, input *SanitizeRequest, file io.ReadSeeker) (*pangea.PangeaResponse[SanitizeResult], error) {
	if input == nil {
		return nil, errors.New("nil input")
	}

	if input.TransferMethod == pangea.TMpostURL {
		params, err := pangea.GetUploadFileParams(file)
		if err != nil {
			return nil, err
		}
		input.SHA256 = params.SHA256
		input.CRC32C = params.CRC32C
		input.Size = pangea.Int(params.Size)
	}

	name := "file"
	if input.TransferMethod == pangea.TMmultipart {
		name = "upload"
	}

	fd := pangea.FileData{
		File: file,
		Name: name,
	}

	return request.DoPostWithFile(ctx, e.Client, "v1beta/sanitize", input, &SanitizeResult{}, fd)
}

func (e *sanitize) RequestUploadURL(ctx context.Context, input *SanitizeRequest) (*pangea.PangeaResponse[SanitizeResult], error) {
	if input.TransferMethod == pangea.TMmultipart || input.TransferMethod == pangea.TMdestURL || input.TransferMethod == pangea.TMsourceURL {
		return nil, fmt.Errorf("transfer method [%s] is not supported in RequestUploadURL. Use Sanitize() method instead.", input.TransferMethod)
	}

	if input.TransferMethod == pangea.TMpostURL && (input.SHA256 == "" || input.CRC32C == "" || input.Size == nil) {
		return nil, errors.New("Need to set SHA256, CRC32C and Size in order to use TMpostURL")
	}

	return request.GetUploadURL(ctx, e.Client, "v1beta/sanitize", input, &SanitizeResult{})
}

// SanitizeFile represents the SanitizeFile API request model.
type SanitizeFile struct {
	ScanProvider string `json:"scan_provider,omitempty"`
	CDRProvider  string `json:"cdr_provider,omitempty"`
}

// SanitizeContent represents the SanitizeContent API request model.
type SanitizeContent struct {
	URLIntel            *bool  `json:"url_intel,omitempty"`
	URLIntelProvider    string `json:"url_intel_provider,omitempty"`
	DomainIntel         *bool  `json:"domain_intel,omitempty"`
	DomainIntelProvider string `json:"domain_intel_provider,omitempty"`
	Defang              *bool  `json:"defang,omitempty"`
	DefangThreshold     *int   `json:"defang_threshold,omitempty"`
	Redact              *bool  `json:"redact,omitempty"`
	RemoveAttachments   *bool  `json:"remove_attachments,omitempty"`
	RemoveInteractive   *bool  `json:"remove_interactive,omitempty"`
}

// SanitizeShareOutput represents the SanitizeShareOutput API request model.
type SanitizeShareOutput struct {
	Enabled      *bool  `json:"enabled,omitempty"`
	OutputFolder string `json:"output_folder,omitempty"`
}

// SanitizeRequest represents the SanitizeRequest API request model.
type SanitizeRequest struct {
	pangea.BaseRequest
	pangea.TransferRequest

	SourceURL        string               `json:"source_url,omitempty"`
	ShareID          string               `json:"share_id,omitempty"`
	File             *SanitizeFile        `json:"file,omitempty"`
	Content          *SanitizeContent     `json:"content,omitempty"`
	ShareOutput      *SanitizeShareOutput `json:"share_output,omitempty"`
	Size             *int                 `json:"size,omitempty"`
	CRC32C           string               `json:"crc32c,omitempty"`
	SHA256           string               `json:"sha256,omitempty"`
	UploadedFileName string               `json:"uploaded_file_name,omitempty"`
}

// DefangData represents the DefangData PangeaResponseResult.
type DefangData struct {
	ExternalURLsCount    int    `json:"external_urls_count"`
	ExternalDomainsCount int    `json:"external_domains_count"`
	DefangedCount        int    `json:"defanged_count"`
	URLIntelSummary      string `json:"url_intel_summary"`
	DomainIntelSummary   string `json:"domain_intel_summary"`
}

// RedactData represents the RedactData PangeaResponseResult.
type RedactData struct {
	RedactionCount int            `json:"redaction_count"`
	SummaryCounts  map[string]int `json:"summary_counts"`
}

// CDR represents the CDR PangeaResponseResult.
type CDR struct {
	FileAttachmentsRemoved     int `json:"file_attachments_removed"`
	InteractiveContentsRemoved int `json:"interactive_contents_removed"`
}

// SanitizeData represents the SanitizeData PangeaResponseResult.
type SanitizeData struct {
	Defang        *DefangData `json:"defang,omitempty"`
	Redact        *RedactData `json:"redact,omitempty"`
	MaliciousFile bool        `json:"malicious_file"`
	CDR           *CDR        `json:"cdr,omitempty"`
}

// SanitizeResult represents the SanitizeResult PangeaResponseResult.
type SanitizeResult struct {
	DestURL     *string                `json:"dest_url,omitempty"`
	DestShareID *string                `json:"dest_share_id,omitempty"`
	Data        SanitizeData           `json:"data"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
}

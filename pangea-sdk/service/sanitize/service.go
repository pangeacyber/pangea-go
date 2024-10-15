package sanitize

import (
	"context"
	"io"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
)

type Client interface {
	// Sanitize.
	Sanitize(ctx context.Context, input *SanitizeRequest, file io.ReadSeeker) (*pangea.PangeaResponse[SanitizeResult], error)

	// Sanitize via presigned URL.
	RequestUploadURL(ctx context.Context, input *SanitizeRequest) (*pangea.PangeaResponse[SanitizeResult], error)

	// Base service methods
	pangea.BaseServicer
}

type sanitize struct {
	pangea.BaseService
}

func New(cfg *pangea.Config) Client {
	cli := &sanitize{
		BaseService: pangea.NewBaseService("sanitize", cfg),
	}
	return cli
}

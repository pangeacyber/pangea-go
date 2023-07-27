package file_scan

import (
	"context"
	"io"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
)

type Client interface {
	Scan(ctx context.Context, input *FileScanRequest, file io.Reader) (*pangea.PangeaResponse[FileScanResult], error)

	// Base service methods
	GetPendingRequestID() []string
	PollResultByError(ctx context.Context, e pangea.AcceptedError) (*pangea.PangeaResponse[any], error)
	PollResultByID(ctx context.Context, rid string, v any) (*pangea.PangeaResponse[any], error)
	PollResultRaw(ctx context.Context, requestID string) (*pangea.PangeaResponse[map[string]any], error)
}

type FileScan struct {
	pangea.BaseService
}

func New(cfg *pangea.Config) Client {
	cli := &FileScan{
		BaseService: pangea.NewBaseService("file-scan", false, cfg),
	}
	return cli
}

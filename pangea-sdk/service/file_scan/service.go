package file_scan

import (
	"context"
	"io"

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
)

type Client interface {
	Scan(ctx context.Context, input *FileScanRequest, file io.Reader) (*pangea.PangeaResponse[FileScanResult], error)
	PollResult(ctx context.Context, e pangea.AcceptedError) (*pangea.PangeaResponse[any], error)
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

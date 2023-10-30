package file_scan

import (
	"context"
	"io"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
)

type Client interface {
	Scan(ctx context.Context, input *FileScanRequest, file io.Reader) (*pangea.PangeaResponse[FileScanResult], error)

	// Base service methods
	pangea.BaseServicer
}

type FileScan struct {
	pangea.BaseService
}

func New(cfg *pangea.Config) Client {
	cli := &FileScan{
		BaseService: pangea.NewBaseService("file-scan", cfg),
	}
	return cli
}

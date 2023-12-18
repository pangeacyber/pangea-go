package file_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
)

type Client interface {
	Reputation(ctx context.Context, req *FileReputationRequest) (*pangea.PangeaResponse[FileReputationResult], error)
	ReputationBulk(ctx context.Context, req *FileReputationBulkRequest) (*pangea.PangeaResponse[FileReputationBulkResult], error)

	// Base service methods
	pangea.BaseServicer
}

type fileIntel struct {
	pangea.BaseService
}

func New(cfg *pangea.Config) Client {
	cli := &fileIntel{
		BaseService: pangea.NewBaseService("file-intel", cfg),
	}
	return cli
}

package file_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
)

type Client interface {
	Reputation(ctx context.Context, req *FileReputationRequest) (*pangea.PangeaResponse[FileReputationResult], error)

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
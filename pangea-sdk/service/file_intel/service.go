package file_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
)

type Client interface {
	Reputation(ctx context.Context, input *FileReputationRequest) (*pangea.PangeaResponse[FileReputationResult], error)
}

type fileIntel struct {
	pangea.BaseService
}

func New(cfg *pangea.Config) Client {
	cli := &fileIntel{
		BaseService: pangea.NewBaseService("file-intel", false, cfg),
	}
	return cli
}

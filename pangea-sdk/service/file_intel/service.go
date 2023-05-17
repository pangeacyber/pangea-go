package file_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
)

type Client interface {
	Lookup(ctx context.Context, input *FileLookupInput) (*pangea.PangeaResponse[FileLookupOutput], error)
	Reputation(ctx context.Context, input *FileReputationRequest) (*pangea.PangeaResponse[FileReputationResult], error)
}

type FileIntel struct {
	pangea.BaseService
}

func New(cfg *pangea.Config) Client {
	cli := &FileIntel{
		BaseService: pangea.NewBaseService("file-intel", cfg),
	}
	return cli
}

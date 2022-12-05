package file_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
)

type Client interface {
	Lookup(ctx context.Context, input *FileLookupInput) (*pangea.PangeaResponse[FileLookupOutput], error)
}

type FileIntel struct {
	*pangea.Client
}

func New(cfg *pangea.Config) *FileIntel {
	cli := &FileIntel{
		Client: pangea.NewClient("file-intel", cfg),
	}
	return cli
}

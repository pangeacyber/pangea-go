package file_intel

import (
	"context"

	"github.com/pangeacyber/go-pangea/pangea"
)

type Client interface {
	Lookup(ctx context.Context, input *FileLookupInput) (*pangea.PangeaResponse[FileLookupOutput], error)
}

type FileIntel struct {
	*pangea.Client
}

func New(cfg *pangea.Config, opts ...Option) (*FileIntel, error) {
	cli := &FileIntel{
		Client: pangea.NewClient("file-intel", cfg),
	}
	for _, opt := range opts {
		err := opt(cli)
		if err != nil {
			return nil, err
		}
	}
	return cli, nil
}

type Option func(*FileIntel) error

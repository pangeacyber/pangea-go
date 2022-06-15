package audit

import (
	"context"

	"github.com/pangeacyber/go-pangea/pangea"
)

type Client interface {
	Log(ctx context.Context, input *LogInput) (*LogOutput, *pangea.Response, error)
	Search(ctx context.Context, input *SerarchInput) (*SearchOutput, *pangea.Response, error)
}

type Audit struct {
	*pangea.Client
}

func New(cfg pangea.Config) *Audit {
	return &Audit{
		Client: pangea.NewClient(cfg),
	}
}

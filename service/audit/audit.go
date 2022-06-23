package audit

import (
	"context"

	"github.com/pangeacyber/go-pangea/pangea"
)

type Client interface {
	Log(context.Context, *LogInput) (*LogOutput, *pangea.Response, error)
	Search(context.Context, *SerarchInput) (*SearchOutput, *pangea.Response, error)
	Root(context.Context, *RootInput) (*RootOutput, *pangea.Response, error)
}

type Audit struct {
	*pangea.Client
}

func New(cfg pangea.Config) *Audit {
	return &Audit{
		Client: pangea.NewClient(cfg),
	}
}

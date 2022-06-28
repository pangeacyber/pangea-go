package audit

import (
	"context"

	"github.com/pangeacyber/go-pangea/pangea"
)

type Client interface {
	Log(context.Context, *LogInput) (*LogOutput, *pangea.Response, error)
	Search(context.Context, *SearchInput) (*SearchOutput, *pangea.Response, error)
	SearchResults(context.Context, *SeachResultInput) (*SeachResultOutput, *pangea.Response, error)
	Root(context.Context, *RootInput) (*Root, *pangea.Response, error)
}

type Audit struct {
	*pangea.Client
}

func New(cfg pangea.Config) *Audit {
	return &Audit{
		Client: pangea.NewClient(cfg),
	}
}

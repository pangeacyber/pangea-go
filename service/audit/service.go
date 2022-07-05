package audit

import (
	"context"
	"fmt"

	"github.com/pangeacyber/go-pangea/internal/signer"
	"github.com/pangeacyber/go-pangea/pangea"
)

type Client interface {
	Log(context.Context, *LogInput) (*LogOutput, *pangea.Response, error)
	Search(context.Context, *SearchInput) (*SearchOutput, *pangea.Response, error)
	SearchResults(context.Context, *SeachResultInput) (*SeachResultOutput, *pangea.Response, error)
	Root(context.Context, *RootInput) (*RootOutput, *pangea.Response, error)
}

type Audit struct {
	*pangea.Client

	SignLogs bool
	Signer   signer.Signer

	VerifyRecords bool
	Verifier      signer.Verifier
}

func New(cfg *pangea.Config, opts ...Option) (*Audit, error) {
	cli := &Audit{
		Client: pangea.NewClient(cfg),
	}
	for _, opt := range opts {
		err := opt(cli)
		if err != nil {
			return nil, err
		}
	}
	return cli, nil
}

type Option func(*Audit) error

func WithLogSignatureVerificationEnabled(filename string) Option {
	return func(a *Audit) error {
		a.VerifyRecords = true
		v, err := signer.NewPrivateKeyFromFile(filename)
		if err != nil {
			return fmt.Errorf("audit: failed verifier creation: %w", err)
		}
		a.Verifier = v
		return nil
	}
}

func WithLogSigningEnabled(filename string) Option {
	return func(a *Audit) error {
		a.SignLogs = true
		s, err := signer.NewPrivateKeyFromFile(filename)
		if err != nil {
			return fmt.Errorf("audit: failed signer creation: %w", err)
		}
		a.Signer = s
		return nil
	}
}

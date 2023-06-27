package audit

import (
	"context"
	"fmt"

	"github.com/pangeacyber/pangea-go/pangea-sdk/internal/signer"
	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
)

type Client interface {
	Log(context.Context, IEvent, bool) (*pangea.PangeaResponse[LogResult], error)
	Search(context.Context, *SearchInput, IEvent) (*pangea.PangeaResponse[SearchOutput], error)
	SearchResults(context.Context, *SearchResultInput, IEvent) (*pangea.PangeaResponse[SearchResultOutput], error)
	Root(context.Context, *RootInput) (*pangea.PangeaResponse[RootOutput], error)
}

type IEvent interface {
	GetTenantID() string
	SetTenantID(string)
	NewFromJSON([]byte) (IEvent, error)
}

type LogSigningMode int

const (
	Unsigned  LogSigningMode = 0
	LocalSign                = 1
)

type audit struct {
	*pangea.Client

	signer                *signer.Signer
	verifyProofs          bool
	skipEventVerification bool
	publicKeyInfo         map[string]string
	rp                    RootsProvider
	lastUnpRootHash       *string
	tenantID              string
}

func New(cfg *pangea.Config, opts ...Option) (Client, error) {
	cli := &audit{
		Client:                pangea.NewClient("audit", true, cfg),
		skipEventVerification: false,
		rp:                    nil,
		lastUnpRootHash:       nil,
		signer:                nil,
		publicKeyInfo:         nil,
		tenantID:              "",
	}
	for _, opt := range opts {
		err := opt(cli)
		if err != nil {
			return nil, err
		}
	}
	return cli, nil
}

type Option func(*audit) error

func WithLogProofVerificationEnabled() Option {
	return func(a *audit) error {
		a.verifyProofs = true
		return nil
	}
}

func WithLogLocalSigning(filename string) Option {
	return func(a *audit) error {
		s, err := signer.NewSignerFromPrivateKeyFile(filename)
		if err != nil {
			return fmt.Errorf("audit: failed signer creation: %w", err)
		}
		a.signer = &s
		return nil
	}
}

func WithTenantID(tenantID string) Option {
	return func(a *audit) error {
		a.tenantID = tenantID
		return nil
	}
}

func DisableEventVerification() Option {
	return func(a *audit) error {
		a.skipEventVerification = true
		return nil
	}
}

func SetPublicKeyInfo(pkinfo map[string]string) Option {
	return func(a *audit) error {
		a.publicKeyInfo = pkinfo
		return nil
	}
}

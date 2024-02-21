package audit

import (
	"context"
	"fmt"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/internal/signer"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
)

type Client interface {
	Log(ctx context.Context, event any, verbose bool) (*pangea.PangeaResponse[LogResult], error)
	LogBulk(ctx context.Context, event []any, verbose bool) (*pangea.PangeaResponse[LogBulkResult], error)
	LogBulkAsync(ctx context.Context, event []any, verbose bool) (*pangea.PangeaResponse[LogBulkResult], error)
	Search(ctx context.Context, req *SearchInput) (*pangea.PangeaResponse[SearchOutput], error)
	SearchResults(ctx context.Context, req *SearchResultsInput) (*pangea.PangeaResponse[SearchResultsOutput], error)
	Root(ctx context.Context, req *RootInput) (*pangea.PangeaResponse[RootOutput], error)

	// Base service methods
	pangea.BaseServicer
}

type Tenanter interface {
	Tenant() string
	SetTenant(string)
}

type LogSigningMode int

const (
	Unsigned  LogSigningMode = 0
	LocalSign LogSigningMode = 1
)

type audit struct {
	pangea.BaseService

	signer                *signer.Signer
	verifyProofs          bool
	skipEventVerification bool
	publicKeyInfo         map[string]string
	rp                    RootsProvider
	lastUnpRootHash       *string
	tenantID              string
	schema                any
}

func New(cfg *pangea.Config, opts ...Option) (Client, error) {
	cli := &audit{
		BaseService:           pangea.NewBaseService("audit", cfg),
		skipEventVerification: false,
		rp:                    nil,
		lastUnpRootHash:       nil,
		signer:                nil,
		publicKeyInfo:         nil,
		tenantID:              "",
		schema:                StandardEvent{},
	}

	// FIXME: Just to still support ConfigID in PangeaConfig. Remove when deprecated
	if cfg.ConfigID != "" {
		err := WithConfigID(cfg.ConfigID)(cli)
		if err != nil {
			return nil, err
		}
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

func WithConfigID(cid string) Option {
	return func(a *audit) error {
		return pangea.WithConfigID(cid)(&a.BaseService)
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

func WithCustomSchema(schema any) Option {
	return func(a *audit) error {
		a.schema = schema
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

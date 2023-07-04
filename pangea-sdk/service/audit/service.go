package audit

import (
	"context"
	"fmt"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/internal/signer"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
)

type Client interface {
	Log(context.Context, any, bool) (*pangea.PangeaResponse[LogResult], error)
	Search(context.Context, *SearchInput) (*pangea.PangeaResponse[SearchOutput], error)
	SearchResults(context.Context, *SearchResultsInput) (*pangea.PangeaResponse[SearchResultsOutput], error)
	Root(context.Context, *RootInput) (*pangea.PangeaResponse[RootOutput], error)

	// Base service methods
	GetPendingRequestID() []string
	PollResultByError(ctx context.Context, e pangea.AcceptedError) (*pangea.PangeaResponse[any], error)
	PollResultByID(ctx context.Context, rid string, v any) (*pangea.PangeaResponse[any], error)
	PollResultRaw(ctx context.Context, requestID string) (*pangea.PangeaResponse[map[string]any], error)
}

type Tenanter interface {
	Tenant() string
	SetTenant(string)
}

type LogSigningMode int

const (
	Unsigned  LogSigningMode = 0
	LocalSign                = 1
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
		BaseService:           pangea.NewBaseService("audit", true, cfg),
		skipEventVerification: false,
		rp:                    nil,
		lastUnpRootHash:       nil,
		signer:                nil,
		publicKeyInfo:         nil,
		tenantID:              "",
		schema:                StandardEvent{},
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

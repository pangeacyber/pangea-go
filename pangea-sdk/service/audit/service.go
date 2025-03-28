package audit

import (
	"context"
	"fmt"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/internal/signer"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
)

type Client interface {
	Log(ctx context.Context, event any, verbose bool) (*pangea.PangeaResponse[LogResult], error)
	LogBulk(ctx context.Context, event []any, verbose bool) (*pangea.PangeaResponse[LogBulkResult], error)
	LogBulkAsync(ctx context.Context, event []any, verbose bool) (*pangea.PangeaResponse[LogBulkResult], error)
	Search(ctx context.Context, req *SearchInput) (*pangea.PangeaResponse[SearchOutput], error)
	SearchResults(ctx context.Context, req *SearchResultsInput) (*pangea.PangeaResponse[SearchResultsOutput], error)
	Root(ctx context.Context, req *RootInput) (*pangea.PangeaResponse[RootOutput], error)

	// Get all search results as a compressed (gzip) CSV file.
	DownloadResults(ctx context.Context, input *DownloadRequest) (*pangea.PangeaResponse[DownloadResult], error)

	// This API allows 3rd party vendors (like Auth0) to stream events to this
	// endpoint where the structure of the payload varies across different
	// vendors.
	LogStream(ctx context.Context, input pangea.ConfigIDer) (*pangea.PangeaResponse[struct{}], error)

	// Bulk export of data from the Secure Audit Log, with optional filtering.
	Export(ctx context.Context, input *ExportRequest) (*pangea.PangeaResponse[struct{}], error)

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

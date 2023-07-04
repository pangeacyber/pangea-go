package file_intel

import (
	"context"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
)

type Client interface {
	Reputation(ctx context.Context, input *FileReputationRequest) (*pangea.PangeaResponse[FileReputationResult], error)

	// Base service methods
	GetPendingRequestID() []string
	PollResultByException(ctx context.Context, e pangea.AcceptedError) (*pangea.PangeaResponse[any], error)
	PollResultByID(ctx context.Context, rid string, v any) (*pangea.PangeaResponse[any], error)
	PollResultRaw(ctx context.Context, requestID string) (*pangea.PangeaResponse[map[string]any], error)
}

type fileIntel struct {
	pangea.BaseService
}

func New(cfg *pangea.Config) Client {
	cli := &fileIntel{
		BaseService: pangea.NewBaseService("file-intel", false, cfg),
	}
	return cli
}

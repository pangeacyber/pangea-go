package pangea

import (
	"context"
	"errors"
)

type BaseService struct {
	Client *Client
}

func NewBaseService(name string, baseCfg *Config) BaseService {
	cfg := baseCfg.Copy()
	if cfg.Logger == nil {
		cfg.Logger = GetDefaultPangeaLogger()
	}
	bs := BaseService{
		Client: NewClient(name, cfg),
	}
	return bs
}

func (bs *BaseService) PollResultByError(ctx context.Context, e AcceptedError) (*PangeaResponse[any], error) {
	if e.RequestID == nil {
		return nil, errors.New("Request ID is empty")
	}

	resp, err := bs.PollResultByID(ctx, *e.RequestID, e.ResultField)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (bs *BaseService) PollResultByID(ctx context.Context, rid string, v any) (*PangeaResponse[any], error) {
	resp, err := bs.Client.FetchAcceptedResponse(ctx, rid, v)
	if err != nil {
		return nil, err
	}

	return &PangeaResponse[any]{
		Response: *resp,
		Result:   &v,
	}, nil
}

func (bs *BaseService) PollResultRaw(ctx context.Context, rid string) (*PangeaResponse[map[string]any], error) {
	r := make(map[string]any)
	resp, err := bs.Client.FetchAcceptedResponse(ctx, rid, &r)
	if err != nil {
		return nil, err
	}

	return &PangeaResponse[map[string]any]{
		Response: *resp,
		Result:   &r,
	}, nil
}

func (bs *BaseService) GetNumRequestsInProgress() int {
	return bs.Client.GetNumRequestsInProgress()
}

func (bs *BaseService) WaitGroup() {
	bs.Client.WaitGroup()
}

func (bs *BaseService) Close() {
	// Shall we do anything else here?
	bs.Client.KillRequestsInProgress()
}

func (bs *BaseService) KillRequestsInProgress() {
	bs.Client.KillRequestsInProgress()
	bs.Client.WaitGroup()
}

func (bs *BaseService) GetPendingRequestID() []string {
	return bs.Client.GetPendingRequestID()
}

type BaseServicer interface {
	GetPendingRequestID() []string
	PollResultByError(ctx context.Context, e AcceptedError) (*PangeaResponse[any], error)
	PollResultByID(ctx context.Context, rid string, v any) (*PangeaResponse[any], error)
	PollResultRaw(ctx context.Context, requestID string) (*PangeaResponse[map[string]any], error)
}

type Option func(*BaseService) error

func WithConfigID(cid string) Option {
	return func(b *BaseService) error {
		return ClientWithConfigID(cid)(b.Client)
	}
}

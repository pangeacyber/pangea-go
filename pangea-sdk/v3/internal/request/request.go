package request

import (
	"context"
	"errors"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
)

func DoPost[T any](ctx context.Context, c *pangea.Client, urlStr string, input pangea.ConfigIDer, out *T) (*pangea.PangeaResponse[T], error) {
	if input == nil || out == nil {
		return nil, errors.New("nil pointer to struct")
	}

	if c == nil {
		return nil, errors.New("nil client")
	}

	req, err := c.NewRequest("POST", urlStr, input)
	if err != nil {
		return nil, err
	}
	resp, err := c.Do(ctx, req, out)
	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[T]{
		Response: *resp,
		Result:   out,
	}
	return &panresp, nil
}

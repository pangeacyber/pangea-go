package request

import (
	"context"
	"errors"
	"io"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
)

func DoPost[T any](ctx context.Context, c *pangea.Client, path string, input pangea.ConfigIDer, out *T) (*pangea.PangeaResponse[T], error) {
	if input == nil || out == nil {
		return nil, errors.New("nil pointer to struct")
	}

	if c == nil {
		return nil, errors.New("nil client")
	}

	url, err := c.GetURL(path)
	if err != nil {
		c.Logger.Error().
			Str("service", c.ServiceName()).
			Str("method", "DoPost").
			Err(err)
		return nil, err
	}

	req, err := c.NewRequest("POST", url, input)
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

func DoPostWithFile[T any](ctx context.Context, c *pangea.Client, path string, input any, out *T, file io.Reader) (*pangea.PangeaResponse[T], error) {
	url, err := c.GetURL(path)
	if err != nil {
		c.Logger.Error().
			Str("service", c.ServiceName()).
			Str("method", "DoPost").
			Err(err)
		return nil, err
	}

	var resp *pangea.Response
	err = nil
	v, ok := input.(pangea.TransferRequester)

	if ok && v.GetTransferMethod() == pangea.TMdirect { // Check TrasnferMethod
		resp, err = c.PostPresignedURL(ctx, url, input, out, file)
	} else {
		resp, err = c.PostMultipart(ctx, url, input, out, file)
	}

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[T]{
		Response: *resp,
		Result:   out,
	}
	return &panresp, nil
}

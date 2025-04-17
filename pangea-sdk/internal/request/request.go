package request

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
)

type Queryer interface {
	URLQuery() (url.Values, error)
}

func getPostRequest(c *pangea.Client, path string, input pangea.ConfigIDer) (*http.Request, error) {
	if input == nil {
		return nil, errors.New("nil pointer to struct")
	}

	if c == nil {
		return nil, errors.New("nil client")
	}

	url, err := c.GetURL(path)
	if err != nil {
		c.Logger.Error().
			Str("service", c.ServiceName()).
			Str("method", "GetURL").
			Err(err)
		return nil, err
	}

	req, err := c.NewRequest("POST", url, input)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// Delete makes a DELETE request with the given URL, params, and optionally
// deserializes to a response.
func Delete(ctx context.Context, client *pangea.Client, path string, params any, res any) error {
	url, err := client.GetURL(path)
	if err != nil {
		return err
	}

	if params != nil {
		if queryer, ok := params.(Queryer); ok {
			query, err := queryer.URLQuery()
			if err != nil {
				return err
			}
			url.RawQuery = query.Encode()
		}
	}

	req, err := client.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	_, err = client.Do(ctx, req, res, false)
	if err != nil {
		return err
	}

	return nil
}

// Get makes a GET request with the given URL, params, and optionally
// deserializes to a response.
func Get(ctx context.Context, client *pangea.Client, path string, params any, res any) error {
	url, err := client.GetURL(path)
	if err != nil {
		return err
	}

	if params != nil {
		if queryer, ok := params.(Queryer); ok {
			query, err := queryer.URLQuery()
			if err != nil {
				return err
			}
			url.RawQuery = query.Encode()
		}
	}

	req, err := client.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	response, err := client.BareDo(ctx, req)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, res)
}

func DoPost[T any](ctx context.Context, c *pangea.Client, path string, input pangea.ConfigIDer, out *T) (*pangea.PangeaResponse[T], error) {
	if out == nil {
		return nil, errors.New("nil pointer to struct")
	}

	req, err := getPostRequest(c, path, input)
	if err != nil {
		return nil, err
	}

	// Pass true to HANDLE 202 in queue
	resp, err := c.Do(ctx, req, out, true)
	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[T]{
		Response: *resp,
		Result:   out,
	}
	return &panresp, nil
}

func DoPostNonPangeaResponse[T any](ctx context.Context, c *pangea.Client, path string, input pangea.ConfigIDer, out *T) (*T, error) {
	if out == nil {
		return nil, errors.New("nil pointer to struct")
	}

	req, err := getPostRequest(c, path, input)
	if err != nil {
		return nil, err
	}

	// Pass true to HANDLE 202 in queue
	_, err = c.Do(ctx, req, out, true)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func DoPostNoQueue[T any](ctx context.Context, c *pangea.Client, path string, input pangea.ConfigIDer, out *T) (*pangea.PangeaResponse[T], error) {
	if out == nil {
		return nil, errors.New("nil pointer to struct")
	}

	req, err := getPostRequest(c, path, input)
	if err != nil {
		return nil, err
	}

	// Pass false to NOT handle 202 in queue
	resp, err := c.Do(ctx, req, out, false)
	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[T]{
		Response: *resp,
		Result:   out,
	}
	return &panresp, nil
}

func GetUploadURL[T any](ctx context.Context, c *pangea.Client, path string, input pangea.ConfigIDer, out *T) (*pangea.PangeaResponse[T], error) {
	url, err := c.GetURL(path)
	if err != nil {
		c.Logger.Error().
			Str("service", c.ServiceName()).
			Str("method", "DoPost").
			Err(err)
		return nil, err
	}

	pr, ar, err := c.GetPresignedURL(ctx, url, input)
	if err != nil {
		c.Logger.Error().
			Str("service", c.ServiceName()).
			Str("method", "GetPresignedURL").
			Err(err)
		return nil, err
	}

	if pr == nil {
		err = errors.New("GetPresignedURL return nil response pointer")
		c.Logger.Error().
			Str("service", c.ServiceName()).
			Str("method", "GetPresignedURL").
			Err(err)
		return nil, err
	}

	return &pangea.PangeaResponse[T]{
		Response:       *pr,
		AcceptedResult: ar,
		Result:         out,
	}, nil
}

func DoPostWithFile[T any](ctx context.Context, c *pangea.Client, path string, input pangea.ConfigIDer, out *T, fd pangea.FileData) (*pangea.PangeaResponse[T], error) {
	url, err := c.GetURL(path)
	if err != nil {
		c.Logger.Error().
			Str("service", c.ServiceName()).
			Str("method", "DoPost").
			Err(err)
		return nil, err
	}

	var resp *pangea.Response
	v, ok := input.(pangea.TransferRequester)

	if ok && v.GetTransferMethod() == pangea.TMpostURL { // Check TransferMethod
		resp, err = c.FullPostPresignedURL(ctx, url, input, out, fd)
	} else {
		resp, err = c.PostMultipart(ctx, url, input, out, fd)
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

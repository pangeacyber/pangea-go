package pangea

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	userAgent = "go-pangea"
)

var errNonNilContext = errors.New("context must be non-nil")

type EndpointConfig struct {
	Scheme string
	Region string
	CSP    string
}

type Config struct {
	// The Berarer token used to authenticate requests.
	Token string

	// The HTTP client to be used by the client.
	//  It defaults to http.DefaultClient
	HTTPClient *http.Client

	// EndpointConfig is the configuration for the endpoint.
	//  It overrides the Endpoint field if non-nil.
	EndpointConfig *EndpointConfig

	// Base URL for API requests.
	// 	BaseURL should always be specified with a trailing slash.
	//  Used for testing.
	Endpoint string

	// AdditionalHeaders is a map of additional headers to be sent with the request.
	AdditionalHeaders map[string]string
}

// A Client manages communication with the Pangea API.
type Client struct {
	// The auth token of the user.
	Token string

	// The client's config.
	Config Config

	// User agent used when communicating with the Pangea API.
	UserAgent string
}

func NewClient(cfg Config) *Client {
	c := &Client{
		Token: cfg.Token,
	}
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = http.DefaultClient
	}
	c.Config = cfg
	c.UserAgent = userAgent
	return c
}

func mergeHeaders(req *http.Request, additionalHeaders map[string]string) {
	for k, v := range additionalHeaders {
		req.Header.Add(k, v)
	}
}

// Path should be absolute and start with a slash.
func (c *Client) serviceUrl(service, path string) (string, error) {
	cfg := c.Config
	if cfg.EndpointConfig != nil {
		if cfg.EndpointConfig.Region != "" {
			return fmt.Sprintf("%s://%s.%s.%s.pangea.cloud/%s", cfg.EndpointConfig.Scheme, service, cfg.EndpointConfig.Region, cfg.EndpointConfig.CSP, path), nil
		}
		return fmt.Sprintf("%s://%s.%s.pangea.cloud/%s", cfg.EndpointConfig.Scheme, service, cfg.EndpointConfig.CSP, path), nil
	}
	u, err := url.Parse(cfg.Endpoint + path)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the Endpoint or EndpointConfig of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, service, urlStr string, body interface{}) (*http.Request, error) {
	u, err := c.serviceUrl(service, urlStr)
	if err != nil {
		return nil, err
	}
	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u, buf)
	if err != nil {
		return nil, err
	}
	if c.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	mergeHeaders(req, c.Config.AdditionalHeaders)
	return req, nil
}

type Response struct {
	HTTPResponse *http.Response
	ResponseMetadata
	Result json.RawMessage `json:"result"`
}

func (r *Response) UnMarshalResult(target interface{}) error {
	return json.Unmarshal(r.Result, target)
}

// newResponse takes a http.Response and tries to parse the body into a base pangea API response.
func newResponse(r *http.Response) (*Response, error) {
	response := &Response{HTTPResponse: r}
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, NewAPIError(err, r, nil)
	}
	if err := json.Unmarshal(data, response); err != nil {
		return nil, NewUnMarshalError(err, data, r, nil)
	}
	return response, nil
}

// BareDo sends an API request and lets you handle the api response.
//	If an error or API Error occurs, the error will contain more information. Otherwise you
// 	are supposed to read and close the response's Body.
//
// 	The provided ctx must be non-nil, if it is nil an error is returned. If it is
// 	canceled or times out, ctx.Err() will be returned.
func (c *Client) BareDo(ctx context.Context, req *http.Request) (*Response, error) {
	if ctx == nil {
		return nil, errNonNilContext
	}

	resp, err := c.Config.HTTPClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, NewAPIError(err, resp, nil)
	}

	response, err := newResponse(resp)
	if err != nil {
		return nil, err
	}

	err = CheckResponse(response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v is nil, and no error hapens, the response is returned as is.
// The provided ctx must be non-nil, if it is nil an error is returned. If it
// is canceled or times out, ctx.Err() will be returned.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.BareDo(ctx, req)
	if err != nil {
		return nil, err
	}
	switch v := v.(type) {
	case nil:
	default:
		err = resp.UnMarshalResult(v)
		if err != nil {
			return nil, NewUnMarshalError(err, resp.Result, resp.HTTPResponse, &resp.ResponseMetadata)
		}
	}
	return resp, nil
}

// compareHTTPResponse returns whether two http.Response objects are equal or not.
// Currently, only StatusCode is checked. This function is used when implementing the
// Is(error) bool interface for the custom error types in this package.
func compareHTTPResponse(r1, r2 *http.Response) bool {
	if r1 == nil && r2 == nil {
		return true
	}

	if r1 != nil && r2 != nil {
		return r1.StatusCode == r2.StatusCode
	}
	return false
}

func CheckResponse(r *Response) error {
	if r.HTTPResponse.StatusCode == http.StatusAccepted {
		return &AcceptedError{ResponseMetadata: r.ResponseMetadata}
	}
	if r.HTTPResponse.StatusCode <= http.StatusOK {
		return nil
	}
	return &APIError{
		HTTPResponse:     r.HTTPResponse,
		ResponseMetadata: &r.ResponseMetadata,
	}
}

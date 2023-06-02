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
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/pangeacyber/pangea-go/pangea-sdk/internal/defaults"
)

const (
	version         = "1.8.0"
	pangeaUserAgent = "pangea-go/" + version
)

var errNonNilContext = errors.New("context must be non-nil")

type IBaseRequest interface {
	SetConfigID(configID string)
	GetConfigID() string
}

type BaseRequest struct {
	ConfigID string `json:"config_id,omitempty"`
}

func (br *BaseRequest) GetConfigID() string {
	return br.ConfigID
}

func (br *BaseRequest) SetConfigID(c string) {
	br.ConfigID = c
}

type RetryConfig struct {
	RetryWaitMin time.Duration // Minimum time to wait
	RetryWaitMax time.Duration // Maximum time to wait
	RetryMax     int           // Maximum number of retries
}

type Config struct {
	// The Bearer token used to authenticate requests.
	Token string

	// Config ID for multi-config projects
	ConfigID string

	// The HTTP client to be used by the client.
	//  It defaults to defaults.HTTPClient
	HTTPClient *http.Client

	// Base domain for API requests.
	Domain string

	// Set to true to use plain http
	Insecure bool

	// Set to "local" for testing locally
	Enviroment string

	// AdditionalHeaders is a map of additional headers to be sent with the request.
	AdditionalHeaders map[string]string

	// Custom user agent is a string to be added to pangea sdk user agent header and identify app
	CustomUserAgent string

	// if it should retry request
	// if HTTPClient is set in the config this value won't take effect
	Retry bool

	// Retry config defaults to a base retry option
	RetryConfig *RetryConfig
}

// A Client manages communication with the Pangea API.
type Client struct {
	// The auth token of the user.
	token string

	// The client's config.
	config *Config

	// User agent used when communicating with the Pangea API.
	userAgent string

	// The identifier for the service
	serviceName string

	// Flag to check config ID on request
	checkConfigID bool
}

func NewClient(service string, checkConfigID bool, baseCfg *Config, additionalConfigs ...*Config) *Client {
	cfg := baseCfg.Copy()
	cfg.MergeIn(additionalConfigs...)
	cfg.HTTPClient = chooseHTTPClient(cfg)
	var userAgent string
	if len(baseCfg.CustomUserAgent) > 0 {
		userAgent = pangeaUserAgent + " " + baseCfg.CustomUserAgent
	} else {
		userAgent = pangeaUserAgent
	}
	return &Client{
		serviceName:   service,
		token:         cfg.Token,
		config:        cfg,
		userAgent:     userAgent,
		checkConfigID: checkConfigID,
	}
}

func chooseHTTPClient(cfg *Config) *http.Client {
	if cfg.HTTPClient != nil {
		return cfg.HTTPClient
	}
	if cfg.Retry {
		if cfg.RetryConfig != nil {
			cli := retryablehttp.NewClient()
			cli.RetryMax = cfg.RetryConfig.RetryMax
			cli.RetryWaitMin = cfg.RetryConfig.RetryWaitMin
			cli.RetryWaitMax = cfg.RetryConfig.RetryWaitMax
			return cli.StandardClient()
		}
		return defaults.HTTPClientWithRetries()
	}
	return defaults.HTTPClient()
}

func mergeHeaders(req *http.Request, additionalHeaders map[string]string) {
	for k, v := range additionalHeaders {
		// We don't want to overwrite pangea headers with user additional headers. Ignore them.
		_, ok := req.Header[k]
		if !ok {
			req.Header.Add(k, v)
		}
	}
}

func (c *Client) serviceUrl(service, path string) (string, error) {
	cfg := c.config
	endpoint := ""
	// Remove slashes, just in case
	path = strings.TrimPrefix(path, "/")
	domain := strings.TrimSuffix(cfg.Domain, "/")

	if strings.HasPrefix(cfg.Domain, "http://") || strings.HasPrefix(cfg.Domain, "https://") {
		// URL
		endpoint = fmt.Sprintf("%s/%s", domain, path)
	} else {
		scheme := "https://"
		if cfg.Insecure == true {
			scheme = "http://"
		}
		if cfg.Enviroment == "local" {
			// If we are testing locally do not use service
			endpoint = fmt.Sprintf("%s%s/%s", scheme, domain, path)
		} else {
			endpoint = fmt.Sprintf("%s%s.%s/%s", scheme, service, domain, path)
		}
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the Domain of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, urlStr string, body IBaseRequest) (*http.Request, error) {
	u, err := c.serviceUrl(c.serviceName, urlStr)
	if err != nil {
		return nil, err
	}

	if c.checkConfigID && c.config.ConfigID != "" && body.GetConfigID() == "" {
		body.SetConfigID(c.config.ConfigID)
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
	if c.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
	mergeHeaders(req, c.config.AdditionalHeaders)
	return req, nil
}

type PangeaResponse[T any] struct {
	Response
	Result *T
}

func (r *Response) UnmarshalResult(target interface{}) error {
	return json.Unmarshal(r.RawResult, target)
}

// newResponse takes a http.Response and tries to parse the body into a base pangea API response.
func newResponse(r *http.Response) (*Response, error) {
	response := &Response{HTTPResponse: r}
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, NewUnmarshalError(err, []byte{}, r)
	}
	if err := json.Unmarshal(data, response); err != nil {
		return nil, NewUnmarshalError(err, data, r)
	}
	return response, nil
}

// BareDo sends an API request and lets you handle the api response.
//
//	If an error or API Error occurs, the error will contain more information. Otherwise you
//	are supposed to read and close the response's Body.
func (c *Client) BareDo(ctx context.Context, req *http.Request) (*http.Response, error) {
	resp, err := c.config.HTTPClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v is nil, and no error hapens, the response is returned as is.
// The provided ctx must be non-nil, if it is nil an error is returned. If it
// is canceled or times out, ctx.Err() will be returned.
//
// The provided ctx must be non-nil, if it is nil an error is returned. If it is
// canceled or times out, ctx.Err() will be returned.
func (c *Client) Do(ctx context.Context, req *http.Request, v any) (*Response, error) {
	if ctx == nil {
		return nil, errNonNilContext
	}

	resp, err := c.BareDo(ctx, req)
	if err != nil {
		return nil, err
	}

	response, err := newResponse(resp)
	if err != nil {
		return nil, err
	}

	err = CheckResponse(response)
	if err != nil {
		// Return APIError
		return nil, err
	}

	switch v.(type) {
	case nil:
		// This should never be fired to user because Client is to internal use
		return response, fmt.Errorf("Not initialized struct. Can't unmarshal result from response")
	default:
		err = response.UnmarshalResult(v)
		if err != nil {
			return nil, NewAPIError(err, response)
		}
	}

	return response, nil
}

func CheckResponse(r *Response) error {
	if r.HTTPResponse.StatusCode == http.StatusAccepted {
		return &AcceptedError{ResponseHeader: r.ResponseHeader}
	}

	if r.HTTPResponse.StatusCode == http.StatusOK && *r.ResponseHeader.Status == "Success" {
		return nil
	}

	var apiError error

	var pa PangeaErrors
	err := r.UnmarshalResult(&pa)
	if err != nil {
		pa = PangeaErrors{}
		apiError = fmt.Errorf("API error: %s. Unmarshall Error: %s.", *r.ResponseHeader.Summary, err.Error())
	} else {
		apiError = fmt.Errorf("API error: %s.", *r.ResponseHeader.Summary)
	}

	return &APIError{
		BaseError: BaseError{
			Err:          apiError,
			HTTPResponse: r.HTTPResponse,
		},
		ResponseHeader: &r.ResponseHeader,
		RawResult:      r.RawResult,
		PangeaErrors:   pa,
	}
}

func configHeaderName(key string) string {
	return fmt.Sprintf("x-pangea-%v-config-id", key)
}

// MergeIn merges the passed in configs into the existing config object.
func (c *Config) MergeIn(cfgs ...*Config) {
	for _, other := range cfgs {
		mergeInConfig(c, other)
	}
}

func mergeInConfig(dst *Config, other *Config) {
	if other == nil {
		return
	}

	if other.Token != "" {
		dst.Token = other.Token
	}

	if other.Domain != "" {
		dst.Domain = other.Domain
	}

	if other.Enviroment != "" {
		dst.Enviroment = other.Enviroment
	}

	dst.Insecure = other.Insecure

	if other.AdditionalHeaders != nil {
		dst.AdditionalHeaders = other.AdditionalHeaders
	}

	if other.Retry {
		dst.Retry = other.Retry
	}

	if other.RetryConfig != nil {
		dst.RetryConfig = other.RetryConfig
	}

	dst.ConfigID = other.ConfigID
}

// Copy will return a shallow copy of the Config object. If any additional
// configurations are provided they will be merged into the new config returned.
func (c *Config) Copy(cfgs ...*Config) *Config {
	dst := &Config{}
	dst.MergeIn(c)

	for _, cfg := range cfgs {
		dst.MergeIn(cfg)
	}

	return dst
}

// FetchAcceptedResponse retries the
func (c *Client) FetchAcceptedResponse(ctx context.Context, reqID string, v interface{}) (*Response, error) {
	req, err := c.NewRequest("GET", fmt.Sprintf("request/%v", reqID), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Do(ctx, req, v)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

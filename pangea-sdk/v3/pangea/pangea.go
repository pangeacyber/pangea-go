package pangea

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/internal/defaults"
	pu "github.com/pangeacyber/pangea-go/pangea-sdk/v3/internal/pangeautil"
	"github.com/rs/zerolog"
)

const (
	version         = "3.0.0"
	pangeaUserAgent = "pangea-go/" + version
)

var (
	initFileWriterOnce sync.Once
	file               *os.File
)

var errNonNilContext = errors.New("context must be non-nil")

type ConfigIDer interface {
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
	BackOff      float32       //Exponential back of factor
}

type Config struct {
	// The Bearer token used to authenticate requests.
	Token string

	// Config ID for multi-config projects
	// @deprecated: Set config_id with options in service initialization if supported
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

	// Enable queued request retry support
	QueuedRetryEnabled bool

	// Timeout used to poll results after 202 (in secs)
	PollResultTimeout time.Duration

	// Retry config defaults to a base retry option
	RetryConfig *RetryConfig

	// Logger
	Logger *zerolog.Logger
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

	// Map to save pending requests id
	pendingRequestID map[string]bool

	// Client logger
	Logger zerolog.Logger

	// config ID on request
	configID string
}

func NewClient(service string, baseCfg *Config, additionalConfigs ...*Config) *Client {
	cfg := baseCfg.Copy()
	cfg.MergeIn(additionalConfigs...)

	if cfg.Logger == nil {
		cfg.Logger = GetDefaultPangeaLogger()
	}

	cfg.HTTPClient = chooseHTTPClient(cfg)

	var userAgent string
	if len(baseCfg.CustomUserAgent) > 0 {
		userAgent = pangeaUserAgent + " " + baseCfg.CustomUserAgent
	} else {
		userAgent = pangeaUserAgent
	}

	return &Client{
		serviceName:      service,
		token:            cfg.Token,
		config:           cfg,
		userAgent:        userAgent,
		configID:         "",
		pendingRequestID: make(map[string]bool),
		Logger:           *cfg.Logger,
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
			cli.Logger = nil
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
func (c *Client) NewRequest(method, urlStr string, body ConfigIDer) (*http.Request, error) {
	u, err := c.serviceUrl(c.serviceName, urlStr)
	if err != nil {
		c.Logger.Error().
			Str("service", c.serviceName).
			Str("method", "NewRequest").
			Err(err)
		return nil, err
	}

	c.Logger.Info().
		Str("service", c.serviceName).
		Str("method", "NewRequest").
		Str("action", method).
		Str("url", u).
		Send()

	if c.configID != "" && body.GetConfigID() == "" {
		body.SetConfigID(c.configID)
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	c.Logger.Debug().
		Str("service", c.serviceName).
		Str("method", "NewRequest").
		Str("action", method).
		Str("url", u).
		Interface("data", body).
		Send()

	req, err := http.NewRequest(method, u, buf)
	if err != nil {
		c.Logger.Error().
			Str("service", c.serviceName).
			Str("method", "NewRequest.http").
			Err(err)
		return nil, err
	}

	c.SetHeaders(req)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (c *Client) NewRequestMultPart(method, urlStr string, body interface{}, file io.Reader) (*http.Request, error) {
	u, err := c.serviceUrl(c.serviceName, urlStr)
	if err != nil {
		c.Logger.Error().
			Str("service", c.serviceName).
			Str("method", "NewRequestMultPart").
			Err(err)
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

	c.Logger.Debug().
		Str("service", c.serviceName).
		Str("method", "NewRequestMultPart").
		Str("action", method).
		Str("url", u).
		Interface("data", body).
		Send()

	// Prepare multi part form
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	var fw io.Writer

	// Write request body
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name=request;`))
	h.Set("Content-Type", "application/json")
	if fw, err = w.CreatePart(h); err != nil {
		return nil, err
	}
	if _, err = io.Copy(fw, buf); err != nil {
		return nil, err
	}

	// Write file
	if fw, err = w.CreateFormFile("upload", "filename.exe"); err != nil {
		return nil, err
	}
	if _, err = io.Copy(fw, file); err != nil {
		return nil, err
	}

	// close the multipart writer.
	w.Close()
	req, err := http.NewRequest(method, u, &b)
	if err != nil {
		return nil, err
	}

	c.SetHeaders(req)
	req.Header.Set("Content-Type", w.FormDataContentType())

	return req, nil
}

func (c *Client) SetHeaders(req *http.Request) {
	if c.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
	mergeHeaders(req, c.config.AdditionalHeaders)
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
		c.Logger.Error().
			Str("service", c.serviceName).
			Str("method", "Do").
			Err(errNonNilContext)
		return nil, errNonNilContext
	}

	resp, err := c.BareDo(ctx, req)
	if err != nil {
		c.Logger.Error().
			Str("service", c.serviceName).
			Str("method", "Do.BareDo").
			Err(err)
		return nil, err
	}

	r, err := newResponse(resp)
	if err != nil {
		c.Logger.Error().
			Str("service", c.serviceName).
			Str("method", "Do.newResponse").
			Err(err)
		return nil, err
	}

	r, err = c.handledQueued(ctx, r)
	if err != nil {
		return nil, err
	}

	err = c.CheckResponse(r, v)
	if err != nil {
		c.Logger.Error().
			Str("service", c.serviceName).
			Str("method", "Do.handleQueued").
			Err(err)
		return nil, err
	}

	u := ""
	if r != nil && r.HTTPResponse != nil && r.HTTPResponse.Request != nil && r.HTTPResponse.Request.URL != nil {
		u = r.HTTPResponse.Request.URL.String()
	}

	err = c.CheckResponse(r, v)
	if err != nil {
		c.Logger.Error().
			Str("service", c.serviceName).
			Str("method", "Do.CheckResponse").
			Str("url", u).
			Err(err)
		// Return APIError
		return nil, err
	}

	switch v.(type) {
	case nil:
		// This should never be fired to user because Client is to internal use
		err := fmt.Errorf("Not initialized struct. Can't unmarshal result from response")
		c.Logger.Error().
			Str("service", c.serviceName).
			Str("method", "Do").
			Err(err)
		return r, err
	default:
		err = r.UnmarshalResult(v)
		if err != nil {
			c.Logger.Error().
				Str("service", c.serviceName).
				Str("method", "Do.UnmarshalResult").
				Str("url", u).
				Err(err)
			return nil, NewAPIError(err, r)
		}
	}

	return r, nil
}

func (c *Client) getDelay(retry_count int, start time.Time) time.Duration {
	delay := time.Duration(retry_count*retry_count) * time.Second
	elapsed := time.Since(start)
	//  if with this delay exceed timeout, reduce delay
	if elapsed+delay > c.config.PollResultTimeout {
		delay = c.config.PollResultTimeout - elapsed
	}
	return delay
}

func (c *Client) reachTimeout(start time.Time) bool {
	return time.Since(start) >= c.config.PollResultTimeout
}

func (c *Client) handledQueued(ctx context.Context, r *Response) (*Response, error) {
	if r.HTTPResponse.StatusCode == http.StatusAccepted && (r != nil && r.RequestID != nil) {
		c.addPendingRequestID(*r.RequestID)
	} else {
		return r, nil
	}

	if c.config.QueuedRetryEnabled == false || r == nil || r.HTTPResponse.StatusCode != http.StatusAccepted {
		return r, nil
	}

	start := time.Now()
	var retry = 1
	u := fmt.Sprintf("request/%v", *r.RequestID)

	c.Logger.Info().
		Str("service", c.serviceName).
		Str("method", "handledQueued.Start").
		Str("url", u).
		Send()

	for r.HTTPResponse.StatusCode == http.StatusAccepted && !c.reachTimeout(start) {
		delay := c.getDelay(retry, start)
		if pu.Sleep(delay, ctx) == false {
			// If context closed, return inmediatly
			return r, nil
		}

		req, err := c.NewRequest("GET", fmt.Sprintf("request/%v", *r.RequestID), nil)
		if err != nil {
			return nil, err
		}

		if ctx == nil {
			return nil, errNonNilContext
		}

		resp, err := c.BareDo(ctx, req)
		if err != nil {
			return nil, err
		}

		r, err = newResponse(resp)
		if err != nil {
			return nil, err
		}

		retry++
	}

	if r.HTTPResponse.StatusCode != http.StatusAccepted {
		c.removePendingRequestID(*r.RequestID)
	}

	c.Logger.Info().
		Str("service", c.serviceName).
		Str("method", "handledQueued.Exit").
		Str("url", u).
		Send()

	return r, nil
}

func (c *Client) CheckResponse(r *Response, v any) error {
	if r.HTTPResponse.StatusCode == http.StatusAccepted {
		return &AcceptedError{
			ResponseHeader: r.ResponseHeader,
			ResultField:    v,
		}
	}

	if r.HTTPResponse.StatusCode == http.StatusOK && *r.ResponseHeader.Status == "Success" {
		return nil
	}

	var apiError error

	var pa PangeaErrors
	err := r.UnmarshalResult(&pa)
	em := ""
	if err != nil {
		pa = PangeaErrors{}
		em = fmt.Sprintf("API error: %s. Unmarshall Error: %s.", *r.ResponseHeader.Summary, err.Error())
		apiError = fmt.Errorf(em)
	} else {
		em = fmt.Sprintf("API error: %s.", *r.ResponseHeader.Summary)
		apiError = fmt.Errorf(em)
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

	dst.QueuedRetryEnabled = other.QueuedRetryEnabled
	dst.PollResultTimeout = other.PollResultTimeout
	dst.Logger = other.Logger
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
	u := fmt.Sprintf("request/%v", reqID)
	c.Logger.Info().
		Str("service", c.serviceName).
		Str("method", "FetchAcceptedResponse").
		Str("url", u).
		Send()

	req, err := c.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Do(ctx, req, v)
	if err != nil {
		return nil, err
	}

	c.removePendingRequestID(reqID)

	return resp, nil
}

type BaseService struct {
	Client *Client
}

type BaseServicer interface {
	GetPendingRequestID() []string
	PollResultByError(ctx context.Context, e AcceptedError) (*PangeaResponse[any], error)
	PollResultByID(ctx context.Context, rid string, v any) (*PangeaResponse[any], error)
	PollResultRaw(ctx context.Context, requestID string) (*PangeaResponse[map[string]any], error)
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

type Option func(*BaseService) error

func WithConfigID(cid string) Option {
	return func(b *BaseService) error {
		return ClientWithConfigID(cid)(b.Client)
	}
}

type ClientOption func(*Client) error

func ClientWithConfigID(cid string) ClientOption {
	return func(b *Client) error {
		b.configID = cid
		return nil
	}
}

func (c *Client) GetPendingRequestID() []string {
	keys := make([]string, len(c.pendingRequestID))
	i := 0
	for k := range c.pendingRequestID {
		keys[i] = k
		i++
	}
	return keys
}

func (bs *BaseService) GetPendingRequestID() []string {
	return bs.Client.GetPendingRequestID()
}

func (c *Client) addPendingRequestID(rid string) {
	c.pendingRequestID[rid] = true
}

func (c *Client) removePendingRequestID(rid string) {
	delete(c.pendingRequestID, rid)
}

func initFileWriter() {
	// Open the output file
	filename := "pangea_sdk_log.json"
	var err error
	file, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// Where should we close this file?
	if err != nil {
		fmt.Printf("Failed to open log file: %s. Logger will go to stdout", filename)
		file = os.Stdout
	}
}

func GetDefaultPangeaLogger() *zerolog.Logger {
	// Set up the logger
	initFileWriterOnce.Do(initFileWriter)

	zerolog.TimestampFieldName = "time"
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "message"

	// Set up the JSON file writer as the output
	logger := zerolog.New(file).With().Timestamp().Logger()
	return &logger
}

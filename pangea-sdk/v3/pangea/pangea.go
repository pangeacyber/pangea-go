package pangea

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	version         = "3.8.0-beta.2"
	pangeaUserAgent = "pangea-go/" + version
)

var (
	initFileWriterOnce sync.Once
	file               *os.File
)

var errNonNilContext = errors.New("context must be non-nil")

type TransferMethod string

const (
	TMmultipart TransferMethod = "multipart"
	TMpostURL   TransferMethod = "post-url"
	TMputURL    TransferMethod = "put-url"
	TMsourceURL TransferMethod = "source-url"
	TMdestURL   TransferMethod = "dest-url"
)

type ConfigIDer interface {
	SetConfigID(configID string)
	GetConfigID() string
}

type BaseRequest struct {
	ConfigID string `json:"config_id,omitempty"`
}

type TransferRequester interface {
	GetTransferMethod() TransferMethod
}

type TransferRequest struct {
	TransferMethod TransferMethod `json:"transfer_method,omitempty"`
}

func (tr TransferRequest) GetTransferMethod() TransferMethod {
	return tr.TransferMethod
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

type FileData struct {
	File    io.Reader
	Name    string
	Details map[string]string
}

func (c *Client) ServiceName() string {
	return c.serviceName
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

func (c *Client) GetRequestIDURL(rid string) (string, error) {
	return c.GetURL(fmt.Sprintf("request/%v", rid))
}

func (c *Client) GetURL(path string) (string, error) {
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
		if cfg.Insecure {
			scheme = "http://"
		}
		if cfg.Enviroment == "local" {
			// If we are testing locally do not use service
			endpoint = fmt.Sprintf("%s%s/%s", scheme, domain, path)
		} else {
			endpoint = fmt.Sprintf("%s%s.%s/%s", scheme, c.serviceName, domain, path)
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
func (c *Client) NewRequest(method, url string, body any) (*http.Request, error) {
	c.Logger.Info().
		Str("service", c.serviceName).
		Str("method", "NewRequest").
		Str("action", method).
		Str("url", url).
		Send()

	if c.configID != "" {
		v, ok := body.(ConfigIDer)
		if ok && v.GetConfigID() == "" {
			v.SetConfigID(c.configID)
		}
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
		Str("url", url).
		Interface("data", body).
		Send()

	req, err := http.NewRequest(method, url, buf)
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

func (c *Client) GetPresignedURL(ctx context.Context, url string, input any) (*Response, *AcceptedResult, error) {
	req, err := c.NewRequest("POST", url, input)
	if err != nil {
		return nil, nil, err
	}

	pr, err := c.simplePost(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	c.Logger.Debug().
		Str("service", c.serviceName).
		Str("method", "GetPresignedURL").
		Str("url", url).
		Interface("header", pr.ResponseHeader).
		Interface("result", pr.RawResult).
		Send()

	if *pr.Status == "Success" { // If this request already success, just return. Not need to poll presigned URL
		return pr, nil, nil
	}

	err = c.CheckResponse(pr, &AcceptedResult{})
	var ae *AcceptedError
	var ok bool
	if err != nil {
		// This should return AcceptedError
		ae, ok = err.(*AcceptedError)
		if !ok {
			c.Logger.Error().
				Str("service", c.serviceName).
				Str("method", "GetPresignedURL").
				Str("url", url).
				Err(err)
			// Return APIError
			return nil, nil, err
		}
	}

	ar, err := c.pollPresignedURL(ctx, ae)
	if err != nil {
		return nil, nil, err
	}

	return pr, ar, nil
}

func (c *Client) DownloadFile(ctx context.Context, url string) (*AttachedFile, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.Logger.Error().
			Str("service", c.serviceName).
			Str("method", "DownloadFile.NewRequest").
			Err(err)
		return nil, err
	}

	resp, err := c.BareDo(ctx, req)
	if err != nil {
		c.Logger.Error().
			Str("service", c.serviceName).
			Str("method", "DownloadFile.BareDo").
			Err(err)
		return nil, err
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	filename, _ := pu.GetFilenameFromContentDisposition(resp.Header.Get("Content-Disposition"))
	if filename == "" {
		filename = pu.GetFileNameFromURL(url)
		if filename == "" {
			filename = "default_filename"
		}
	}

	return &AttachedFile{
		Filename:    filename,
		File:        data,
		ContentType: resp.Header.Get("Content-Type"),
	}, nil

}

func (c *Client) UploadFile(ctx context.Context, url string, tm TransferMethod, fd FileData) error {
	if tm == TMputURL && fd.Details != nil {
		return fmt.Errorf("data param should be nil in order to use TransferMethod %s", TMputURL)
	}

	var req *http.Request
	var err error

	if tm == TMputURL {
		var buffer bytes.Buffer
		_, err = io.Copy(&buffer, fd.File)
		if err != nil {
			return err
		}
		req, err = http.NewRequest("PUT", url, bytes.NewReader(buffer.Bytes()))
	} else {
		req, err = c.NewRequestForm("POST", url, fd, false)
	}

	if err != nil {
		return err
	}

	psURLr, err := c.BareDo(ctx, req)
	if err != nil {
		c.Logger.Error().
			Str("service", c.serviceName).
			Str("method", "UploadFile.BareDo").
			Err(err)
		return err
	}

	if psURLr.StatusCode < 200 || psURLr.StatusCode >= 300 {
		defer psURLr.Body.Close()
		return errors.New("presigned url upload failure")
	}
	return nil
}

func (c *Client) FullPostPresignedURL(ctx context.Context, url string, input ConfigIDer, out any, fd FileData) (*Response, error) {
	pr, ar, err := c.GetPresignedURL(ctx, url, input)
	if err != nil {
		c.Logger.Error().
			Str("service", c.serviceName).
			Str("method", "PostPresignedURL.GetPresignedURL").
			Err(err)
		return nil, err
	}

	if ar != nil { // This is the case that GetPresignedURL return an already success response
		fds := FileData{
			File:    fd.File,
			Name:    fd.Name,
			Details: ar.PostFormData,
		}

		err = c.UploadFile(ctx, ar.PostURL, TMpostURL, fds)
		if err != nil {
			c.Logger.Error().
				Str("service", c.serviceName).
				Str("method", "PostPresignedURL.PostPresignedURL").
				Err(err)
			return nil, err
		}

		pr, err = c.handledQueued(ctx, pr)
		if err != nil {
			c.Logger.Error().
				Str("service", c.serviceName).
				Str("method", "PostPresignedURL.handleQueued").
				Err(err)
			return nil, err
		}
	}

	err = c.CheckResponse(pr, out)
	if err != nil {
		c.Logger.Error().
			Str("service", c.serviceName).
			Str("method", "PostPresignedURL.CheckResponse").
			Str("url", url).
			Err(err)
		// Return APIError
		return nil, err
	}

	return pr, nil
}

func (c *Client) PostMultipart(ctx context.Context, url string, input any, out any, fd FileData) (*Response, error) {
	req, err := c.NewRequestMultipart("POST", url, input, fd)
	if err != nil {
		return nil, err
	}
	return c.Do(ctx, req, out, true)
}

func (c *Client) NewRequestMultipart(method, url string, body any, fd FileData) (*http.Request, error) {
	if c.configID != "" {
		v, ok := body.(ConfigIDer)
		if ok && v.GetConfigID() == "" {
			v.SetConfigID(c.configID)
		}
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
		Str("method", "NewRequestMultipart").
		Str("action", method).
		Str("url", url).
		Interface("data", body).
		Send()

	// Prepare multi part form
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	var fw io.Writer

	var err error

	// Write request body
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name=request;`)
	h.Set("Content-Type", "application/json")
	if fw, err = w.CreatePart(h); err != nil {
		return nil, err
	}
	if _, err = io.Copy(fw, buf); err != nil {
		return nil, err
	}

	// Write file
	if fw, err = w.CreateFormFile(fd.Name, fd.Name); err != nil {
		return nil, err
	}
	if _, err = io.Copy(fw, fd.File); err != nil {
		return nil, err
	}

	// close the multipart writer.
	w.Close()
	req, err := http.NewRequest(method, url, &b)
	if err != nil {
		return nil, err
	}

	c.SetHeaders(req)
	req.Header.Set("Content-Type", w.FormDataContentType())

	return req, nil
}

func (c *Client) NewRequestForm(method, url string, fd FileData, setHeaders bool) (*http.Request, error) {
	c.Logger.Debug().
		Str("service", c.serviceName).
		Str("method", "NewRequestForm").
		Str("action", method).
		Str("url", url).
		Interface("data", fd.Details).
		Send()

	// Prepare multi part form
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Write request body fields
	if fd.Details != nil {
		for key, value := range fd.Details {
			if err := w.WriteField(key, value); err != nil {
				return nil, err
			}
		}
	}

	// Write file
	var err error
	part, err := w.CreateFormFile(fd.Name, fd.Name)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, fd.File)
	if err != nil {
		return nil, err
	}

	// close the multipart writer.
	w.Close()
	req, err := http.NewRequest(method, url, &b)
	if err != nil {
		return nil, err
	}

	if setHeaders {
		c.SetHeaders(req)
	}
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

func (c *Client) simplePost(ctx context.Context, req *http.Request) (*Response, error) {
	if ctx == nil {
		c.Logger.Error().
			Str("service", c.serviceName).
			Str("method", "simplePost").
			Err(errNonNilContext)
		return nil, errNonNilContext
	}

	resp, err := c.BareDo(ctx, req)
	if err != nil {
		c.Logger.Error().
			Str("service", c.serviceName).
			Str("method", "simplePost.BareDo").
			Err(err)
		return nil, err
	}

	r, err := newResponse(resp)
	if err != nil {
		c.Logger.Error().
			Str("service", c.serviceName).
			Str("method", "simplePost.newResponse").
			Err(err)
		return nil, err
	}

	return r, nil
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v is nil, and no error happens, the response is returned as is.
// The provided ctx must be non-nil, if it is nil an error is returned. If it
// is canceled or times out, ctx.Err() will be returned.
//
// The provided ctx must be non-nil, if it is nil an error is returned. If it is
// canceled or times out, ctx.Err() will be returned.
func (c *Client) Do(ctx context.Context, req *http.Request, v any, handleQueue bool) (*Response, error) {
	r, err := c.simplePost(ctx, req)
	if err != nil {
		return nil, err
	}

	u := ""
	if r != nil && r.HTTPResponse != nil && r.HTTPResponse.Request != nil && r.HTTPResponse.Request.URL != nil {
		u = r.HTTPResponse.Request.URL.String()
	}

	c.Logger.Debug().
		Str("service", c.serviceName).
		Str("method", "Do").
		Str("url", u).
		Interface("header", r.ResponseHeader).
		Interface("result", r.RawResult).
		Send()

	if handleQueue {
		r, err = c.handledQueued(ctx, r)
		if err != nil {
			c.Logger.Error().
				Str("service", c.serviceName).
				Str("method", "Do.handleQueued").
				Err(err)
			return nil, err
		}
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

func (c *Client) pollPresignedURL(ctx context.Context, ae *AcceptedError) (*AcceptedResult, error) {
	if ae.AcceptedResult.HasUploadURL() {
		return &ae.AcceptedResult, nil
	}

	u, err := c.GetRequestIDURL(*ae.RequestID)
	if err != nil {
		c.Logger.Error().
			Str("service", c.ServiceName()).
			Str("method", "PollPresignedURL").
			Err(err)
		return nil, err
	}

	var aeLoop = ae
	var ok bool
	start := time.Now()
	var retry = 1

	for !aeLoop.AcceptedResult.HasUploadURL() && !c.reachTimeout(start) {
		delay := c.getDelay(retry, start)
		if !pu.Sleep(delay, ctx) {
			// If context closed, return inmediatly
			return nil, errors.New("context closed")
		}

		req, err := c.NewRequest("GET", u, nil)
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

		r, err := newResponse(resp)
		if err != nil {
			return nil, err
		}

		err = c.CheckResponse(r, ae.ResultField)
		aeLoop, ok = err.(*AcceptedError)
		if !ok {
			return nil, err
		}
		retry++
	}

	if c.reachTimeout(start) {
		return nil, aeLoop
	}

	return &aeLoop.AcceptedResult, nil
}

func (c *Client) handledQueued(ctx context.Context, r *Response) (*Response, error) {
	if r.HTTPResponse.StatusCode == http.StatusAccepted && (r != nil && r.RequestID != nil) {
		c.addPendingRequestID(*r.RequestID)
	} else {
		return r, nil
	}

	if !c.config.QueuedRetryEnabled || r == nil || r.HTTPResponse.StatusCode != http.StatusAccepted {
		return r, nil
	}

	start := time.Now()
	var retry = 1
	u, err := c.GetRequestIDURL(*r.RequestID)
	if err != nil {
		c.Logger.Error().
			Str("service", c.ServiceName()).
			Str("method", "handledQueued").
			Err(err)
		return nil, err
	}

	c.Logger.Info().
		Str("service", c.serviceName).
		Str("method", "handledQueued.Start").
		Str("url", u).
		Send()

	for r.HTTPResponse.StatusCode == http.StatusAccepted && !c.reachTimeout(start) {
		delay := c.getDelay(retry, start)
		if !pu.Sleep(delay, ctx) {
			// If context closed, return inmediatly
			return r, nil
		}

		req, err := c.NewRequest("GET", u, nil)
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

	c.Logger.Debug().
		Str("service", c.serviceName).
		Str("method", "handleQueued").
		Str("url", u).
		Interface("header", r.ResponseHeader).
		Interface("result", r.RawResult).
		Send()

	return r, nil
}

func (c *Client) CheckResponse(r *Response, v any) error {
	if r.HTTPResponse.StatusCode == http.StatusAccepted {
		var ar AcceptedResult
		err := r.UnmarshalResult(&ar)
		if err != nil {
			ar = AcceptedResult{}
		}

		return &AcceptedError{
			ResponseHeader: r.ResponseHeader,
			ResultField:    v,
			AcceptedResult: ar,
			Response:       *r,
		}
	}

	if r.HTTPResponse.StatusCode == http.StatusOK && *r.ResponseHeader.Status == "Success" {
		switch v.(type) {
		case nil:
			// This should never be fired to user because Client is to internal use
			return fmt.Errorf("not initialized struct. Can't unmarshal result from response")
		default:
			err := r.UnmarshalResult(v)
			if err != nil {
				return NewAPIError(err, r)
			}
		}
		return nil
	}

	var apiError error
	var pa PangeaErrors
	em := ""

	err := r.UnmarshalResult(&pa)
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
	u, err := c.GetRequestIDURL(reqID)
	if err != nil {
		c.Logger.Error().
			Str("service", c.ServiceName()).
			Str("method", "FetchAcceptedResponse").
			Err(err)
		return nil, err
	}

	c.Logger.Info().
		Str("service", c.serviceName).
		Str("method", "FetchAcceptedResponse").
		Str("url", u).
		Send()

	req, err := c.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Do(ctx, req, v, false)
	if err != nil {
		return nil, err
	}

	c.removePendingRequestID(reqID)

	return resp, nil
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

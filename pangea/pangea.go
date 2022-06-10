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

var defaultConfig = &Config{
	HTTPClient: http.DefaultClient,
	EndpointConfig: &EndpointConfig{
		Scheme: "https",
		Region: "us-east-1",
		CSP:    "aws",
	},
}

// A Client manages communication with the Pangea API.
type Client struct {
	// The auth token of the user.
	Token string

	// The client's config.
	Config *Config

	// User agent used when communicating with the Pangea API.
	UserAgent string
}

func NewClient(token string, cfg *Config) *Client {
	c := &Client{
		Token: token,
	}
	if cfg == nil {
		cfg = defaultConfig
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

type ResponseMetadata struct {
	RequestID    *string `json:"request_id"`    // The request ID
	RequestTime  *string `json:"request_time"`  // The time the request was issued, ISO8601
	ResponseTime *string `json:"response_time"` // The time the response was issued, ISO8601
	Status       *string `json:"status"`        //
	StatusCode   *int    `json:"status_code"`
	Summary      *string `json:"summary"`
}

type Response struct {
	HTTPResponse *http.Response
	ResponseMetadata
	Result json.RawMessage `json:"result"`
}

func (r *Response) UnMarshalResult(target interface{}) error {
	return json.Unmarshal(r.Result, target)
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

	req = req.WithContext(ctx)
	resp, err := c.Config.HTTPClient.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return nil, err
	}
	response := &Response{HTTPResponse: resp}

	err = CheckResponse(resp)
	if err != nil {
		defer resp.Body.Close()
		aerr, ok := err.(*AcceptedError)
		if ok {
			b, readErr := ioutil.ReadAll(resp.Body)
			if readErr != nil {
				return response, readErr
			}

			aerr.Raw = b
			err = aerr
		}
	}
	return response, err
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v is nil, and no error hapens, the response is returned as is.
// The provided ctx must be non-nil, if it is nil an error is returned. If it
// is canceled or times out, ctx.Err() will be returned.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.BareDo(ctx, req)
	if err != nil {
		return resp, err
	}
	defer resp.HTTPResponse.Body.Close()
	decErr := json.NewDecoder(resp.HTTPResponse.Body).Decode(resp)
	if decErr == io.EOF {
		decErr = nil // ignore EOF errors caused by empty response body
	}
	if decErr != nil {
		err = decErr
	}
	switch v := v.(type) {
	case nil:
	default:
		resultErr := resp.UnMarshalResult(v)
		if resultErr == io.EOF {
			resultErr = nil // ignore EOF errors caused by empty response body
		}
		if resultErr != nil {
			err = decErr
		}
	}
	return resp, err
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

type ErrorResponse struct {
	Response *http.Response // HTTP response that caused this error
	ResponseMetadata
	Result json.RawMessage `json:"result"`
	// the error if any ocurred durin sending or encoding the request
	Err error
}

func (r *ErrorResponse) Error() string {
	s := fmt.Sprintf("%v %v: %d ", r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode)
	if r.Summary != nil {
		s += fmt.Sprintf("%v", r.Summary)
	}
	if len(r.Result) > 0 {
		s += fmt.Sprintf("%+v", r.Result)
	}
	if r.Err != nil {
		s += fmt.Sprintf("; %v", r.Err)
	}
	return s
}

// Is returns whether the provided error equals this error.
func (r *ErrorResponse) Is(target error) bool {
	v, ok := target.(*ErrorResponse)
	if !ok {
		return false
	}

	if r.Summary != v.Summary ||
		!compareHTTPResponse(r.Response, v.Response) {
		return false
	}

	if len(r.Result) != len(v.Result) {
		return false
	}

	for idx := range r.Result {
		if r.Result[idx] != v.Result[idx] {
			return false
		}
	}
	return true
}

// AcceptedError occurs when Pangea returns 202 Accepted response
// which means the request was send to Queue to process asynchronously.
type AcceptedError struct {
	// Raw contains the response body.
	Raw []byte
}

func (*AcceptedError) Error() string {
	return "request scheduled on Pangea side: please check the status of the request later"
}

// Is returns whether the provided error equals this error.
func (ae *AcceptedError) Is(target error) bool {
	v, ok := target.(*AcceptedError)
	if !ok {
		return false
	}
	return bytes.Compare(ae.Raw, v.Raw) == 0
}

func CheckResponse(r *http.Response) error {
	if r.StatusCode == http.StatusAccepted {
		return &AcceptedError{}
	}

	if r.StatusCode <= 200 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		if err := json.Unmarshal(data, errorResponse); err != nil {
			errorResponse.Err = fmt.Errorf("error decoding error response: %s", err)
		}
	}
	return errorResponse
}

// Bool is a helper routine that allocates a new bool value
// to store v and returns a pointer to it.
func Bool(v bool) *bool { return &v }

// Int is a helper routine that allocates a new int value
// to store v and returns a pointer to it.
func Int(v int) *int { return &v }

// String is a helper routine that allocates a new string value
// to store v and returns a pointer to it.
func String(v string) *string { return &v }

// StringValue is a helper routine that returns the value of a string pointer or a default value if nil
func StringValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

// IntValue is a helper routine that returns the value of a int pointer or a default value if nil
func IntValue(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}

// BoolValue is a helper routine that returns the value of a bool pointer or a default value if nil
func BoolValue(v *bool) bool {
	if v == nil {
		return false
	}
	return *v
}

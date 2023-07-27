package pangea

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type PangeaErrors struct {
	Errors []ErrorField `json:"errors,omitempty"`
}

type ErrorField struct {
	Code   string `json:"code"`
	Detail string `json:"detail"`
	Source string `json:"source"`
	Path   string `json:"path,omitempty"`
}

type BaseError struct {
	// The underlying error that triggered this one, if any.
	Err error

	// the HTTP response
	HTTPResponse *http.Response
}

type APIError struct {
	BaseError

	// the reponse header of the request if any
	ResponseHeader *ResponseHeader

	// the result of the request
	RawResult json.RawMessage

	// error details
	PangeaErrors PangeaErrors
}

func NewAPIError(err error, r *Response) *APIError {
	var errRes error
	var pa PangeaErrors
	errUnm := r.UnmarshalResult(&pa)
	if errUnm != nil {
		pa = PangeaErrors{}
		errRes = fmt.Errorf("Error: %s. Unmarshall Error: %s.", err.Error(), errUnm.Error())
	} else {
		errRes = fmt.Errorf("Error: %s.", err.Error())
	}

	return &APIError{
		BaseError: BaseError{
			Err:          errRes,
			HTTPResponse: r.HTTPResponse,
		},
		RawResult:      r.RawResult,
		ResponseHeader: &r.ResponseHeader,
		PangeaErrors:   pa,
	}
}

func (e *BaseError) Error() string {
	return e.Err.Error()
}

func (e *APIError) Error() string {
	b := new(bytes.Buffer)
	if e.HTTPResponse != nil {
		b.WriteString(fmt.Sprintf("pangea: %v %v", e.HTTPResponse.Request.Method, e.HTTPResponse.Request.URL))
		pad(b, ": ")
		if e.ResponseHeader != nil {
			b.WriteString(fmt.Sprintf("%v", e.ResponseHeader.String()))
		} else {
			b.WriteString(fmt.Sprintf("%v", e.HTTPResponse.StatusCode))
		}
	}
	if e.Err != nil {
		pad(b, ": ")
		b.WriteString(e.Err.Error())
	}
	return b.String()
}

type UnmarshalError struct {
	BaseError

	Bytes []byte
}

func NewUnmarshalError(err error, bytes []byte, r *http.Response) *UnmarshalError {
	return &UnmarshalError{
		BaseError: BaseError{
			Err:          err,
			HTTPResponse: r,
		},
		Bytes: bytes,
	}
}

func (e *UnmarshalError) Error() string {
	b := new(bytes.Buffer)
	b.WriteString("Failed to unmarshall body")
	if len(e.Bytes) > 0 {
		pad(b, ": ")
		b.WriteString(fmt.Sprintf("body: %s", string(e.Bytes)))
	}
	return b.String()
}

// pad appends str to the buffer if the buffer already has some data.
func pad(b *bytes.Buffer, str string) {
	if b.Len() == 0 {
		return
	}
	b.WriteString(str)
}

type AcceptedError struct {
	ResponseHeader
	ResultField any
}

func (e *AcceptedError) Error() string {
	return fmt.Sprintf("request scheduled on Pangea side: please check the status of the request %v later", e.ReqID())
}

// Is returns whether the provided error equals this error.
func (e *AcceptedError) Is(target error) bool {
	v, ok := target.(*AcceptedError)
	if !ok {
		return false
	}
	return v.RequestID == e.RequestID
}

// ReqID is a helper function to get the request ID
func (e *AcceptedError) ReqID() string {
	return StringValue(e.RequestID)
}

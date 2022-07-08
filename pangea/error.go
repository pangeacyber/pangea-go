package pangea

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type APIError struct {
	// the HTTP response
	HTTPResponse *http.Response

	// the reponse metadata of the request if any
	ResponseMetadata *ResponseMetadata

	// the result of the request
	Result json.RawMessage

	// The underlying error that triggered this one, if any.
	Err error
}

func NewAPIError(err error, r *http.Response, metadata *ResponseMetadata) *APIError {
	return &APIError{
		Err:              err,
		HTTPResponse:     r,
		ResponseMetadata: metadata,
	}
}

func (e *APIError) Error() string {
	b := new(bytes.Buffer)
	if e.HTTPResponse != nil {
		b.WriteString(fmt.Sprintf("pangea: %v %v", e.HTTPResponse.Request.Method, e.HTTPResponse.Request.URL))
		pad(b, ": ")
		if e.ResponseMetadata != nil {
			b.WriteString(fmt.Sprintf("%v", e.ResponseMetadata.String()))
		} else {
			b.WriteString(fmt.Sprintf("%v", e.HTTPResponse.StatusCode))
		}
		if len(e.Result) > 0 {
			pad(b, ": ")
			b.WriteString(fmt.Sprintf("body: %s", e.Result))
		}
	}
	if e.Err != nil {
		pad(b, ": ")
		b.WriteString(e.Err.Error())
	}
	return b.String()
}

type UnMarshalError struct {
	APIError
	Bytes []byte
}

func NewUnMarshalError(err error, bytes []byte, r *http.Response, metadata *ResponseMetadata) *UnMarshalError {
	return &UnMarshalError{
		APIError: APIError{
			Err:              err,
			HTTPResponse:     r,
			ResponseMetadata: metadata,
		},
		Bytes: bytes,
	}
}

func (e *UnMarshalError) Error() string {
	b := new(bytes.Buffer)
	b.WriteString("pangea: failed to unmarshall body")
	errMsg := e.APIError.Error()
	if errMsg != "" {
		pad(b, ": ")
		b.WriteString(errMsg)
	}
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
	ResponseMetadata
}

func (e *AcceptedError) Error() string {
	return fmt.Sprintf("request scheduled on Pangea side: please check the status of the request %v later", e.ResponseMetadata.RequestID)
}

// Is returns whether the provided error equals this error.
func (e *AcceptedError) Is(target error) bool {
	v, ok := target.(*AcceptedError)
	if !ok {
		return false
	}
	return v.RequestID == e.RequestID
}

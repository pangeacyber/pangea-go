package pangea

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type ResponseHeader struct {
	// The request ID
	RequestID *string `json:"request_id"`

	// The time the request was issued, ISO8601
	RequestTime *string `json:"request_time"`

	// The time the response was issued, ISO8601
	ResponseTime *string `json:"response_time"`

	// The HTTP status code msg
	Status *string `json:"status"`

	// The summary of the response
	Summary *string `json:"summary"`
}

type Response struct {
	ResponseHeader
	HTTPResponse *http.Response
	// Query raw result
	RawResult   json.RawMessage `json:"result"`
	rawResponse []byte
}

func (r *ResponseHeader) String() string {
	if r == nil {
		return ""
	}
	return fmt.Sprintf("request_id: %v, request_time: %v, response_time: %v, status: %v, summary: %v",
		StringValue(r.RequestID), StringValue(r.RequestTime), StringValue(r.ResponseTime),
		StringValue(r.Status), StringValue(r.Summary))
}

type PangeaResponse[T any] struct {
	Response
	Result         *T
	AcceptedResult *AcceptedResult
}

func (r *Response) UnmarshalResult(target interface{}) error {
	return json.Unmarshal(r.RawResult, target)
}

// newResponse takes a http.Response and tries to parse the body into a base pangea API response.
func newResponse(r *http.Response) (*Response, error) {
	response := &Response{HTTPResponse: r}
	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, NewUnmarshalError(err, []byte{}, r)
	}
	response.rawResponse = data
	if err := json.Unmarshal(data, response); err != nil {
		return nil, NewUnmarshalError(err, data, r)
	}
	return response, nil
}

// MarshalJSON implements the json.Marshaler interface for CustomType.
func (r Response) MarshalJSON() ([]byte, error) {
	if r.rawResponse == nil {
		return nil, errors.New("Unable to read response body")
	}
	b := make([]byte, len(r.rawResponse))
	// Using the copy function to copy the contents.
	nc := copy(b, r.rawResponse)
	if nc != len(r.rawResponse) {
		return nil, errors.New("Unable to copy raw response")
	}
	return b, nil
}

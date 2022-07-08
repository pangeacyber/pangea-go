package pangea

import (
	"fmt"
)

type ResponseMetadata struct {
	// The request ID
	RequestID *string `json:"request_id"`

	// The time the request was issued, ISO8601
	RequestTime *string `json:"request_time"`

	// The time the response was issued, ISO8601
	ResponseTime *string `json:"response_time"`

	// The HTTP status code msg
	Status *string `json:"status"`

	// The HTTP status code
	StatusCode *int `json:"status_code"`

	// The summary of the response
	Summary *string `json:"summary"`
}

func (r *ResponseMetadata) String() string {
	if r == nil {
		return ""
	}
	return fmt.Sprintf("request_id: %v, request_time: %v, response_time: %v, status_code: %v, status: %v, summary: %v",
		StringValue(r.RequestID), StringValue(r.RequestTime), StringValue(r.ResponseTime), IntValue(r.StatusCode),
		StringValue(r.Status), StringValue(r.Summary))
}

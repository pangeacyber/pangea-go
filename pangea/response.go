package pangea

import (
	"bytes"
	"fmt"
)

type ResponseMetadata struct {
	RequestID    *string `json:"request_id"`    // The request ID
	RequestTime  *string `json:"request_time"`  // The time the request was issued, ISO8601
	ResponseTime *string `json:"response_time"` // The time the response was issued, ISO8601
	Status       *string `json:"status"`        //
	StatusCode   *int    `json:"status_code"`
	Summary      *string `json:"summary"`
}

func (r *ResponseMetadata) String() string {
	b := new(bytes.Buffer)
	if r == nil {
		return ""
	}
	if r.RequestID != nil {
		b.WriteString(fmt.Sprintf("request_id: %v", *r.RequestID))
	}
	if r.RequestTime != nil {
		pad(b, ", ")
		b.WriteString(fmt.Sprintf("request_time: %v", *r.RequestTime))
	}
	if r.ResponseTime != nil {
		pad(b, ", ")
		b.WriteString(fmt.Sprintf("response_time: %v", *r.ResponseTime))
	}
	if r.Status != nil {
		pad(b, ", ")
		b.WriteString(fmt.Sprintf("status: %v", *r.Status))
	}
	if r.StatusCode != nil {
		pad(b, ", ")
		b.WriteString(fmt.Sprintf("status_code: %v", *r.StatusCode))
	}
	if r.Summary != nil {
		pad(b, ", ")
		b.WriteString(fmt.Sprintf("summary: %v", *r.Summary))
	}
	return b.String()
}

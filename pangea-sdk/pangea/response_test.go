//go:build unit

package pangea_test

import (
	"testing"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/pangea"
)

func TestResponseHeader_String(t *testing.T) {
	want := "request_id: some-id, " +
		"request_time: 1970-01-01T00:00:00Z, " +
		"response_time: 1970-01-01T00:00:10Z, " +
		"status: I'm a teapot, " +
		"summary: I'm a teapot"

	metadata := &pangea.ResponseHeader{
		RequestID:    pangea.String("some-id"),
		RequestTime:  pangea.String("1970-01-01T00:00:00Z"),
		ResponseTime: pangea.String("1970-01-01T00:00:10Z"),
		Status:       pangea.String("I'm a teapot"),
		Summary:      pangea.String("I'm a teapot"),
	}

	if got := metadata.String(); got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestResponseHeader_String_When_ResponseHeader_Is_Nil_Return_EmptyString(t *testing.T) {
	var header *pangea.ResponseHeader
	if got := header.String(); got != "" {
		t.Errorf("got %v, want empty string", got)
	}
}

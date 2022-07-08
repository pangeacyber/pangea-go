package pangea_test

import (
	"testing"

	"github.com/pangeacyber/go-pangea/pangea"
)

func TestResponseMetadata_String(t *testing.T) {
	want := "request_id: some-id, " +
		"request_time: 1970-01-01T00:00:00Z, " +
		"response_time: 1970-01-01T00:00:10Z, " +
		"status_code: 418, " +
		"status: I'm a teapot, " +
		"summary: I'm a teapot"

	metadata := &pangea.ResponseMetadata{
		RequestID:    pangea.String("some-id"),
		RequestTime:  pangea.String("1970-01-01T00:00:00Z"),
		ResponseTime: pangea.String("1970-01-01T00:00:10Z"),
		StatusCode:   pangea.Int(418),
		Status:       pangea.String("I'm a teapot"),
		Summary:      pangea.String("I'm a teapot"),
	}

	if got := metadata.String(); got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestResponseMetadata_String_When_ResponseMetadata_Is_Nil_Return_EmptyString(t *testing.T) {
	var metadata *pangea.ResponseMetadata
	if got := metadata.String(); got != "" {
		t.Errorf("got %v, want empty string", got)
	}
}

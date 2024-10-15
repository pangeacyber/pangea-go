//go:build unit

package redact_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/internal/pangeatesting"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/service/redact"
	"github.com/stretchr/testify/assert"
)

func TestRedact(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	mux.HandleFunc("/v1/redact", func(w http.ResponseWriter, r *http.Request) {
		pangeatesting.TestMethod(t, r, "POST")
		pangeatesting.TestBody(t, r, `{"text":"My phone number is: 110303456"}`)
		fmt.Fprint(w,
			`{
				"request_id": "some-id",
				"request_time": "1970-01-01T00:00:00Z",
				"response_time": "1970-01-01T00:00:10Z",
				"status": "Success",
				"result":{
					"redacted_text": "My phone number is: <PHONE_NUMBER>",
					"count": 1
				},
				"summary": "success"
			}`)
	})

	client := redact.New(pangeatesting.TestConfig(url))
	input := &redact.TextRequest{
		Text: pangea.String("My phone number is: 110303456"),
	}
	ctx := context.Background()
	got, err := client.Redact(ctx, input)

	assert.NoError(t, err)

	want := &redact.TextResult{
		RedactedText: pangea.String("My phone number is: <PHONE_NUMBER>"),
		Count:        1,
	}
	assert.Equal(t, want, got.Result)
}

func TestRedactError(t *testing.T) {
	f := func(cfg *pangea.Config) error {
		client := redact.New(cfg)
		_, err := client.Redact(context.Background(), nil)
		return err
	}
	pangeatesting.TestNewRequestAndDoFailure(t, "Redact.Redact", f)
}

func TestRedactStructuredError(t *testing.T) {
	f := func(cfg *pangea.Config) error {
		client := redact.New(cfg)
		_, err := client.RedactStructured(context.Background(), nil)
		return err
	}
	pangeatesting.TestNewRequestAndDoFailure(t, "Redact.RedactStructured", f)
}

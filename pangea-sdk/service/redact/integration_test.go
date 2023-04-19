// go:build unit
package redact_test

import (
	"context"
	"testing"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/internal/pangeatesting"
	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/service/redact"
	"github.com/stretchr/testify/assert"
)

const (
	testingEnvironment = pangeatesting.Live
)

func redactIntegrationCfg(t *testing.T) *pangea.Config {
	t.Helper()
	return pangeatesting.IntegrationConfig(t, testingEnvironment)
}

func Test_Integration_Redact(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelFn()

	cfg := redactIntegrationCfg(t)
	client := redact.New(cfg)

	redacted := "My Phone number is <PHONE_NUMBER>"

	input := &redact.TextInput{
		Text: pangea.String("My Phone number is 415-867-5309"),
	}
	out, err := client.Redact(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NoError(t, err)
	assert.NotNil(t, out.Result)
	assert.Equal(t, redacted, *out.Result.RedactedText)
	assert.Equal(t, 1, out.Result.Count)
}

func Test_Integration_Redact_NoResult(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelFn()

	cfg := redactIntegrationCfg(t)
	client := redact.New(cfg)

	input := &redact.TextInput{
		Text:         pangea.String("My Phone number is 415-867-5309"),
		ReturnResult: pangea.Bool(false),
	}
	out, err := client.Redact(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NoError(t, err)
	assert.NotNil(t, out.Result)
	assert.Nil(t, out.Result.RedactedText)
	assert.Equal(t, 1, out.Result.Count)

}

func Test_Integration_Redact_Structured(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelFn()

	cfg := redactIntegrationCfg(t)
	cfg.Retry = true
	client := redact.New(cfg)

	data := map[string]any{"phone": "415-867-5309"}
	redacted := map[string]any{"phone": "<PHONE_NUMBER>"}

	input := &redact.StructuredInput{
		Data: data,
	}
	out, err := client.RedactStructured(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NoError(t, err)
	assert.NotNil(t, out.Result)
	assert.Equal(t, redacted, out.Result.RedactedData)
	assert.Equal(t, 1, out.Result.Count)
}

func Test_Integration_Redact_Structured_NoResult(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelFn()

	cfg := redactIntegrationCfg(t)
	cfg.Retry = true
	client := redact.New(cfg)

	data := map[string]any{"phone": "415-867-5309"}

	input := &redact.StructuredInput{
		Data:         data,
		ReturnResult: pangea.Bool(false),
	}
	out, err := client.RedactStructured(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NoError(t, err)
	assert.NotNil(t, out.Result)
	assert.Nil(t, out.Result.RedactedData)
	assert.Equal(t, 1, out.Result.Count)
}

func Test_Integration_Redact_Error_BadToken(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelFn()

	cfg := redactIntegrationCfg(t)
	cfg.Retry = true
	cfg.Token = "notarealtoken"
	client := redact.New(cfg)

	input := &redact.TextInput{
		Text: pangea.String(""),
	}
	out, err := client.Redact(ctx, input)
	assert.Error(t, err)
	assert.Nil(t, out)
	apiErr := err.(*pangea.APIError)
	assert.Equal(t, apiErr.Err.Error(), "API error: Not authorized to access this resource.")
}

//go:build integration

package redact_test

import (
	"context"
	"testing"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/internal/pangeatesting"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/service/redact"
	"github.com/stretchr/testify/assert"
)

var testingEnvironment = pangeatesting.LoadTestEnvironment("redact", pangeatesting.Live)

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

	input := &redact.TextRequest{
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

func Test_Integration_Redact_DebugTrue(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelFn()

	cfg := redactIntegrationCfg(t)
	client := redact.New(cfg)

	redacted := "My Phone number is <PHONE_NUMBER>"

	input := &redact.TextRequest{
		Text:  pangea.String("My Phone number is 415-867-5309"),
		Debug: pangea.Bool(true),
	}
	out, err := client.Redact(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NoError(t, err)
	assert.NotNil(t, out.Result)
	assert.Equal(t, redacted, *out.Result.RedactedText)
	assert.Equal(t, 1, out.Result.Count)
	assert.NotNil(t, out.Result.Report.RecognizerResults)
	assert.NotEmpty(t, out.Result.Report.RecognizerResults)
	assert.NotNil(t, out.Result.Report.RecognizerResults[0].Score)
	assert.NotNil(t, out.Result.Report.RecognizerResults[0].Text)

}

func Test_Integration_Redact_NoResult(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelFn()

	cfg := redactIntegrationCfg(t)
	client := redact.New(cfg)

	input := &redact.TextRequest{
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

	input := &redact.StructuredRequest{
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

func Test_Integration_Redact_Structured_DebugTrue(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelFn()

	cfg := redactIntegrationCfg(t)
	cfg.Retry = true
	client := redact.New(cfg)

	data := map[string]any{"phone": "415-867-5309"}
	redacted := map[string]any{"phone": "<PHONE_NUMBER>"}

	input := &redact.StructuredRequest{
		Data:  data,
		Debug: pangea.Bool(true),
	}
	out, err := client.RedactStructured(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NoError(t, err)
	assert.NotNil(t, out.Result)
	assert.Equal(t, redacted, out.Result.RedactedData)
	assert.Equal(t, 1, out.Result.Count)
	assert.NotNil(t, out.Result.Report.RecognizerResults)
	assert.NotEmpty(t, out.Result.Report.RecognizerResults)
	assert.NotNil(t, out.Result.Report.RecognizerResults[0].Score)
	assert.NotNil(t, out.Result.Report.RecognizerResults[0].Text)
}

func Test_Integration_Redact_Structured_NoResult(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelFn()

	cfg := redactIntegrationCfg(t)
	cfg.Retry = true
	client := redact.New(cfg)

	data := map[string]any{"phone": "415-867-5309"}

	input := &redact.StructuredRequest{
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

	input := &redact.TextRequest{
		Text: pangea.String(""),
	}
	out, err := client.Redact(ctx, input)
	assert.Error(t, err)
	assert.Nil(t, out)
	apiErr := err.(*pangea.APIError)
	assert.Equal(t, apiErr.Err.Error(), "API error: Not authorized to access this resource.")
}

func Test_Integration_Multi_Config_1_Redact(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := pangeatesting.IntegrationMultiConfigConfig(t, testingEnvironment)
	ConfigID := pangeatesting.GetConfigID(t, testingEnvironment, "redact", 1)
	client := redact.New(cfg, redact.WithConfigID(ConfigID))

	redacted := "My Phone number is <PHONE_NUMBER>"

	input := &redact.TextRequest{
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

func Test_Integration_Multi_Config_2_Log(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := pangeatesting.IntegrationMultiConfigConfig(t, testingEnvironment)
	ConfigID := pangeatesting.GetConfigID(t, testingEnvironment, "redact", 2)
	client := redact.New(cfg, redact.WithConfigID(ConfigID))

	redacted := "My Phone number is <PHONE_NUMBER>"

	input := &redact.TextRequest{
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

func Test_Integration_Multi_Config_No_ConfigID(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := pangeatesting.IntegrationMultiConfigConfig(t, testingEnvironment)
	client := redact.New(cfg)

	input := &redact.TextRequest{
		Text: pangea.String("My Phone number is 415-867-5309"),
	}
	out, err := client.Redact(ctx, input)
	assert.Error(t, err)
	assert.Nil(t, out)
}

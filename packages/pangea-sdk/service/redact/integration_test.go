// go:build unit
package redact_test

import (
	"context"
	"testing"
	"time"

	"github.com/pangeacyber/pangea-go/packages/pangea-sdk/internal/pangeatesting"
	"github.com/pangeacyber/pangea-go/packages/pangea-sdk/pangea"
	"github.com/pangeacyber/pangea-go/packages/pangea-sdk/service/redact"
	"github.com/stretchr/testify/assert"
)

func redactIntegrationCfg(t *testing.T) *pangea.Config {
	t.Helper()
	token := pangeatesting.GetEnvVarOrSkip(t, "PANGEA_INTEGRATION_REDACT_TOKEN")
	if token == "" {
		t.Skip("set PANGEA_INTEGRATION_REDACT_TOKEN env variables to run this test")
	}
	cfg := &pangea.Config{
		Token: token,
	}
	return cfg.Copy(pangeatesting.IntegrationConfig(t))
}

func Test_Integration_Redact(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelFn()

	cfg := redactIntegrationCfg(t)
	client := redact.New(cfg)

	input := &redact.TextInput{
		Text: pangea.String("My Phone number is 110045638"),
	}
	out, err := client.Redact(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, out.Response)
}

func Test_Integration_Redact_Structured(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelFn()

	cfg := redactIntegrationCfg(t)
	cfg.Retry = true
	client := redact.New(cfg)

	input := &redact.TextInput{
		Text: pangea.String("My Phone number is 110045638"),
	}
	out, err := client.Redact(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, out.Response)
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
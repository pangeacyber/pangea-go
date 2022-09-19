// go:build unit
package redact_test

import (
	"context"
	"testing"
	"time"

	"github.com/pangeacyber/go-pangea/internal/pangeatesting"
	"github.com/pangeacyber/go-pangea/pangea"
	"github.com/pangeacyber/go-pangea/service/redact"
	"github.com/stretchr/testify/assert"
)

func Test_Integration_Redact(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelFn()

	cfgToken := pangeatesting.GetEnvVarOrSkip(t, "REDACT_INTEGRATION_CONFIG_TOKEN")
	cfg := &pangea.Config{
		ConfigID: cfgToken,
		Retry:    true,
	}
	cfg = cfg.Copy(pangeatesting.IntegrationConfig(t))
	client, _ := redact.New(cfg)

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

	cfgToken := pangeatesting.GetEnvVarOrSkip(t, "REDACT_INTEGRATION_CONFIG_TOKEN")
	cfg := &pangea.Config{
		ConfigID: cfgToken,
		Retry:    true,
	}
	cfg = cfg.Copy(pangeatesting.IntegrationConfig(t))
	client, _ := redact.New(cfg)

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

	cfgToken := pangeatesting.GetEnvVarOrSkip(t, "REDACT_INTEGRATION_CONFIG_TOKEN")
	cfg := &pangea.Config{
		ConfigID: cfgToken,
		Retry:    true,
	}
	cfg = cfg.Copy(pangeatesting.IntegrationConfig(t))
	cfg.Token = "notarealtoken"
	client, _ := redact.New(cfg)

	input := &redact.TextInput{
		Text: pangea.String(""),
	}
	out, err := client.Redact(ctx, input)
	assert.Error(t, err)
	assert.Nil(t, out)
	apiErr := err.(*pangea.APIError)
	assert.Equal(t, apiErr.Err.Error(), "API error: Not authorized to access this resource.")
}

func Test_Integration_Redact_Error_BadConfigID(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelFn()

	cfgToken := pangeatesting.GetEnvVarOrSkip(t, "REDACT_INTEGRATION_CONFIG_TOKEN")
	cfg := &pangea.Config{
		ConfigID: cfgToken,
		Retry:    true,
	}
	cfg = cfg.Copy(pangeatesting.IntegrationConfig(t))
	cfg.ConfigID = "notarealconfigid"
	client, _ := redact.New(cfg)

	input := &redact.TextInput{
		Text: pangea.String(""),
	}

	out, err := client.Redact(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, out)
	apiErr := err.(*pangea.APIError)
	assert.Equal(t, apiErr.Err.Error(), "API error: Missing Config ID, you can provide using the config header X-Pangea-redact-Config-Id or adding a token scope `service:redact:*:config:r`.")
}

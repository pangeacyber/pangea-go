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
		CfgToken: cfgToken,
		Retry:    true,
	}
	cfg = cfg.Copy(pangeatesting.IntegrationConfig(t))
	client, _ := redact.New(cfg)

	input := &redact.TextInput{
		Text: pangea.String("My Phone number is 110045638"),
	}
	out, _, err := client.Redact(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, out)
}

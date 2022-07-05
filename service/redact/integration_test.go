// go:build integration
package redact_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/pangeacyber/go-pangea/internal/pangeatesting"
	"github.com/pangeacyber/go-pangea/pangea"
	"github.com/pangeacyber/go-pangea/service/redact"
	"github.com/stretchr/testify/assert"
)

func Test_Integration_Check(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := &pangea.Config{
		CfgToken: os.Getenv("REDACT_INTEGRATION_CONFIG_TOKEN"),
	}
	cfg = cfg.Copy(pangeatesting.IntegrationConfig)
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

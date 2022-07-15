// go:build integration
package embargo_test

import (
	"context"
	"testing"
	"time"

	"github.com/pangeacyber/go-pangea/internal/pangeatesting"
	"github.com/pangeacyber/go-pangea/pangea"
	"github.com/pangeacyber/go-pangea/service/embargo"
	"github.com/stretchr/testify/assert"
)

func Test_Integration_Check(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfgToken := pangeatesting.GetEnvVarOrSkip(t, "EMBARGO_INTEGRATION_CONFIG_TOKEN")
	cfg := &pangea.Config{
		CfgToken: cfgToken,
	}
	cfg = cfg.Copy(pangeatesting.IntegrationConfig(t))
	client, _ := embargo.New(cfg)

	input := &embargo.CheckInput{
		ISOCode: pangea.String("CU"),
	}
	out, _, err := client.Check(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, out)
	assert.NotNil(t, out.Count)
	assert.NotZero(t, *out.Count)
}

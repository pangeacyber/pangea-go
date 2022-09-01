// go:build integration
package file_intel_test

import (
	"context"
	"testing"
	"time"

	"github.com/pangeacyber/go-pangea/internal/pangeatesting"
	"github.com/pangeacyber/go-pangea/pangea"
	"github.com/pangeacyber/go-pangea/service/file_intel"
	"github.com/stretchr/testify/assert"
)

func Test_Integration_FileLookup(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfgToken := pangeatesting.GetEnvVarOrSkip(t, "FILE_INTEL_INTEGRATION_CONFIG_TOKEN")
	cfg := &pangea.Config{
		CfgToken: cfgToken,
	}
	cfg = cfg.Copy(pangeatesting.IntegrationConfig(t))
	fileintel, _ := file_intel.New(cfg)

	input := &file_intel.FileLookupInput{
		Hash:     "142b638c6a60b60c7f9928da4fb85a5a8e1422a9ffdc9ee49e17e56ccca9cf6e",
		HashType: "sha256",
		Raw:      true,
		Verbose:  true,
		Provider: "reversinglabs",
	}
	out, err := fileintel.Lookup(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, out)
	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.Data)
	assert.Equal(t, "malicious", out.Result.Data.Verdict)
}

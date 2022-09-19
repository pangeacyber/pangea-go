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
		ConfigID: cfgToken,
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

func Test_Integration_FileLookup_ErrorBadHash(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfgToken := pangeatesting.GetEnvVarOrSkip(t, "FILE_INTEL_INTEGRATION_CONFIG_TOKEN")
	cfg := &pangea.Config{
		ConfigID: cfgToken,
	}
	cfg = cfg.Copy(pangeatesting.IntegrationConfig(t))
	fileintel, _ := file_intel.New(cfg)

	input := &file_intel.FileLookupInput{
		Hash:     "notarealhash",
		HashType: "sha256",
		Raw:      true,
		Verbose:  true,
		Provider: "reversinglabs",
	}
	out, err := fileintel.Lookup(ctx, input)

	// FIXME: This should crash, uncomment and fix once enpoint is fixed
	assert.NoError(t, err)
	assert.NotNil(t, out)
	// apiErr := err.(*pangea.APIError)
	// assert.Equal(t, len(apiErr.PangeaErrors.Errors), 1)
	// assert.Equal(t, apiErr.PangeaErrors.Errors[0].Code, "BelowMinLength")
	// assert.Equal(t, apiErr.PangeaErrors.Errors[0].Detail, "'message' cannot have less than 1 characters")
	// assert.Equal(t, apiErr.PangeaErrors.Errors[0].Source, "/event/message")
}

func Test_Integration_FileLookup_ErrorBadToken(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfgToken := pangeatesting.GetEnvVarOrSkip(t, "FILE_INTEL_INTEGRATION_CONFIG_TOKEN")
	cfg := &pangea.Config{
		ConfigID: cfgToken,
	}
	cfg = cfg.Copy(pangeatesting.IntegrationConfig(t))
	cfg.Token = "142b638c6a60b60c7f9928da4fb85a5a8e1422a9ffdc9ee49e17e56ccca9cf6e"
	fileintel, _ := file_intel.New(cfg)

	input := &file_intel.FileLookupInput{
		Hash:     "notarealhash",
		HashType: "sha256",
		Raw:      true,
		Verbose:  true,
		Provider: "reversinglabs",
	}
	out, err := fileintel.Lookup(ctx, input)
	assert.Error(t, err)
	assert.Nil(t, out)
	apiErr := err.(*pangea.APIError)
	assert.Equal(t, apiErr.Err.Error(), "API error: Not authorized to access this resource.")
}

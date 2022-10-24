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

func intelFileIntegrationCfg(t *testing.T) *pangea.Config {
	t.Helper()
	token := pangeatesting.GetEnvVarOrSkip(t, "PANGEA_INTEGRATION_FILE_INTEL_TOKEN")
	if token == "" {
		t.Skip("set PANGEA_INTEGRATION_FILE_INTEL_TOKEN env variables to run this test")
	}
	cfg := &pangea.Config{
		Token: token,
	}
	return cfg.Copy(pangeatesting.IntegrationConfig(t))
}

func Test_Integration_FileLookup(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := intelFileIntegrationCfg(t)
	fileintel := file_intel.New(cfg)

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
	assert.Equal(t, "Trojan", out.Result.Data.Category[0])
}

func Test_Integration_FileLookup_2(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := intelFileIntegrationCfg(t)
	fileintel := file_intel.New(cfg)

	input := &file_intel.FileLookupInput{
		Hash:     "322ccbd42b7e4fd3a9d0167ca2fa9f6483d9691364c431625f1df542706",
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
	assert.Equal(t, "", out.Result.Data.Verdict)
	assert.Equal(t, "Not Provided", out.Result.Data.Category[0])
}

func Test_Integration_FileLookup_ErrorBadHash(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := intelFileIntegrationCfg(t)
	fileintel := file_intel.New(cfg)

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
	assert.Equal(t, len(apiErr.PangeaErrors.Errors), 1)
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Code, "BelowMinLength")
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Detail, "'hash' cannot have less than 32 characters")
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Source, "/hash")
}

func Test_Integration_FileLookup_ErrorBadHashType(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := intelFileIntegrationCfg(t)
	fileintel := file_intel.New(cfg)

	input := &file_intel.FileLookupInput{
		Hash:     "142b638c6a60b60c7f9928da4fb85a5a8e1422a9ffdc9ee49e17e56ccca9cf6e",
		HashType: "notarealhashtype",
		Raw:      true,
		Verbose:  true,
		Provider: "reversinglabs",
	}
	out, err := fileintel.Lookup(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, out)
	apiErr := err.(*pangea.APIError)
	assert.Equal(t, len(apiErr.PangeaErrors.Errors), 1)
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Code, "NotEnumMember")
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Detail, "'hash_type' must be a one of the following values [sha256 sha1 md5]")
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Source, "/hash_type")
}

func Test_Integration_FileLookup_ErrorBadToken(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := intelFileIntegrationCfg(t)
	cfg.Token = "142b638c6a60b60c7f9928da4fb85a5a8e1422a9ffdc9ee49e17e56ccca9cf6e"
	fileintel := file_intel.New(cfg)

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

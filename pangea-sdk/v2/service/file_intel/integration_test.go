// go:build integration
package file_intel_test

import (
	"context"
	"testing"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/internal/pangeatesting"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/file_intel"
	"github.com/stretchr/testify/assert"
)

const (
	testingEnvironment = pangeatesting.Live
)

func intelFileIntegrationCfg(t *testing.T) *pangea.Config {
	t.Helper()
	return pangeatesting.IntegrationConfig(t, testingEnvironment)
}

func Test_Integration_FileReputation(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := intelFileIntegrationCfg(t)
	intelcli := file_intel.New(cfg)

	input := &file_intel.FileReputationRequest{
		Hash:     "142b638c6a60b60c7f9928da4fb85a5a8e1422a9ffdc9ee49e17e56ccca9cf6e",
		HashType: "sha256",
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "reversinglabs",
	}
	resp, err := intelcli.Reputation(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	assert.NotNil(t, resp.Result.Data)
	assert.Equal(t, "malicious", resp.Result.Data.Verdict)
	assert.Equal(t, "Trojan", resp.Result.Data.Category[0])
}

func Test_Integration_FileReputation_2(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := intelFileIntegrationCfg(t)
	intelcli := file_intel.New(cfg)

	input := &file_intel.FileReputationRequest{
		Hash:     "322ccbd42b7e4fd3a9d0167ca2fa9f6483d9691364c431625f1df542706",
		HashType: "sha256",
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "reversinglabs",
	}
	resp, err := intelcli.Reputation(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	assert.Equal(t, "unknown", resp.Result.Data.Verdict)
	assert.Equal(t, "Not Provided", resp.Result.Data.Category[0])
}

func Test_Integration_FileReputation_ErrorBadHash(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := intelFileIntegrationCfg(t)
	intelcli := file_intel.New(cfg)

	input := &file_intel.FileReputationRequest{
		Hash:     "notarealhash",
		HashType: "sha256",
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "reversinglabs",
	}
	resp, err := intelcli.Reputation(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, resp)
	apiErr := err.(*pangea.APIError)
	assert.Equal(t, len(apiErr.PangeaErrors.Errors), 1)
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Code, "BelowMinLength")
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Detail, "'hash' cannot have less than 32 characters")
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Source, "/hash")
}

func Test_Integration_FileReputation_ErrorBadHashType(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := intelFileIntegrationCfg(t)
	intelcli := file_intel.New(cfg)

	input := &file_intel.FileReputationRequest{
		Hash:     "142b638c6a60b60c7f9928da4fb85a5a8e1422a9ffdc9ee49e17e56ccca9cf6e",
		HashType: "notarealhashtype",
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "reversinglabs",
	}
	resp, err := intelcli.Reputation(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, resp)
	apiErr := err.(*pangea.APIError)
	assert.Equal(t, len(apiErr.PangeaErrors.Errors), 1)
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Code, "NotEnumMember")
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Detail, "'hash_type' must be a one of the following values [sha256 sha1 md5]")
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Source, "/hash_type")
}

func Test_Integration_FileReputation_ErrorBadToken(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := intelFileIntegrationCfg(t)
	cfg.Token = "142b638c6a60b60c7f9928da4fb85a5a8e1422a9ffdc9ee49e17e56ccca9cf6e"
	intelcli := file_intel.New(cfg)

	input := &file_intel.FileReputationRequest{
		Hash:     "notarealhash",
		HashType: "sha256",
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "reversinglabs",
	}
	resp, err := intelcli.Reputation(ctx, input)
	assert.Error(t, err)
	assert.Nil(t, resp)
	apiErr := err.(*pangea.APIError)
	assert.Equal(t, apiErr.Err.Error(), "API error: Not authorized to access this resource.")
}

func Test_Integration_NewFileReputationRequestFromFilepath(t *testing.T) {
	input, err := file_intel.NewFileReputationRequestFromFilepath("./api.go")

	assert.NoError(t, err)
	assert.NotEmpty(t, input.Hash)
	assert.Equal(t, "sha256", input.HashType)
}

func Test_Integration_NewFileReputationRequestFromFilepath_WrongFile(t *testing.T) {
	input, err := file_intel.NewFileReputationRequestFromFilepath("./not/a/real/path/file.exe")

	assert.Error(t, err)
	assert.Nil(t, input)
}
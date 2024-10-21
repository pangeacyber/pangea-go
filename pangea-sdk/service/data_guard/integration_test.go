//go:build integration

package data_guard_test

import (
	"context"
	"testing"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/internal/pangeatesting"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/service/data_guard"
	"github.com/stretchr/testify/assert"
)

var testingEnvironment = pangeatesting.LoadTestEnvironment("data-guard", pangeatesting.Live)

func TestTextGuard(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	client := data_guard.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &data_guard.TextGuardRequest{Text: "hello world"}
	out, err := client.GuardText(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.RedactedPrompt)
	assert.Zero(t, out.Result.Findings.ArtifactCount)
	assert.Zero(t, out.Result.Findings.MaliciousCount)

	input = &data_guard.TextGuardRequest{Text: "security@pangea.cloud"}
	out, err = client.GuardText(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.RedactedPrompt)
	assert.NotZero(t, out.Result.Findings.ArtifactCount)
	assert.Zero(t, out.Result.Findings.MaliciousCount)
}

func TestFileGuard(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	client := data_guard.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &data_guard.FileGuardRequest{FileUrl: "https://pangea.cloud/robots.txt"}
	out, err := client.GuardFile(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, out.Result)
}

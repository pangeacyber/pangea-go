//go:build integration

package ai_guard_test

import (
	"context"
	"testing"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/internal/pangeatesting"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/service/ai_guard"
	"github.com/stretchr/testify/assert"
)

var testingEnvironment = pangeatesting.LoadTestEnvironment("ai-guard", pangeatesting.Live)

func TestTextGuard(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	client := ai_guard.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &ai_guard.TextGuardRequest{Text: "hello world", Recipe: "pangea_prompt_guard"}
	out, err := client.GuardText(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.Prompt)
	assert.False(t, out.Result.Detectors.PromptInjection.Detected)
	assert.False(t, out.Result.Detectors.PiiEntity.Detected)
	assert.False(t, out.Result.Detectors.MaliciousEntity.Detected)

	input = &ai_guard.TextGuardRequest{Text: "security@pangea.cloud", Recipe: "pangea_prompt_guard"}
	out, err = client.GuardText(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.Prompt)
	assert.True(t, out.Result.Detectors.PiiEntity.Detected)
}

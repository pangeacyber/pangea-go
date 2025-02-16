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

	input := &ai_guard.TextGuardRequest{Text: "what was pangea?"}
	out, err := client.GuardText(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.PromptText)
	assert.False(t, out.Result.Detectors.PromptInjection.Detected)
	if out.Result.Detectors.PiiEntity != nil {
		assert.False(t, out.Result.Detectors.PiiEntity.Detected)
	}
	if out.Result.Detectors.MaliciousEntity != nil {
		assert.False(t, out.Result.Detectors.MaliciousEntity.Detected)
	}
}

func TestTextGuard_Messages(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	client := ai_guard.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &ai_guard.TextGuardRequest{Messages: []map[string]interface{}{
		{
			"role":    "user",
			"content": "what was pangea?",
		},
	}}
	out, err := client.GuardText(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.PromptMessages)
}

func TestTextGuard_Llm_Input(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	client := ai_guard.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &ai_guard.TextGuardRequest{LlmInput: map[string]interface{}{
		"model": "gpt-4o",
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": "what was pangea?",
			},
		},
	}}
	out, err := client.GuardText(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.PromptMessages)
}

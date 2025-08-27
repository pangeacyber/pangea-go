//go:build mock

package ai_guard_test

import (
	"context"
	"os"
	"testing"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/option"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/ai_guard"
	"github.com/stretchr/testify/assert"
)

func TestGuardText_Text(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}

	config, err := pangea.NewConfig(
		option.WithBaseURLTemplate(baseURL),
		option.WithToken("my API token"),
	)
	assert.NoError(t, err)
	client := ai_guard.New(config)

	_, err = client.GuardText(context.TODO(), &ai_guard.TextGuardRequest{Text: "what was pangea?"})
	assert.NoError(t, err)
}

func TestGuardText_Messages(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}

	config, err := pangea.NewConfig(
		option.WithBaseURLTemplate(baseURL),
		option.WithToken("my API token"),
	)
	assert.NoError(t, err)
	client := ai_guard.New(config)

	_, err = client.GuardText(context.TODO(), &ai_guard.TextGuardRequest{Messages: []ai_guard.PromptMessage{{
		Role:    "user",
		Content: "what was pangea?",
	}}})
	assert.NoError(t, err)
}

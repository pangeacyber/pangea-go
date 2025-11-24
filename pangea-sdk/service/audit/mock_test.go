//go:build mock

package audit_test

import (
	"context"
	"os"
	"testing"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/option"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/audit"
	"github.com/stretchr/testify/assert"
)

func TestLog_StandardEvent(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}

	config, err := pangea.NewConfig(
		option.WithBaseURLTemplate(baseURL),
		option.WithToken("my API token"),
	)
	assert.NoError(t, err)

	client, err := audit.New(config, audit.DisableEventVerification())
	assert.NoError(t, err)

	_, err = client.Log(context.TODO(), &audit.StandardEvent{Message: "test"}, true)
	assert.NoError(t, err)
}

func TestLog_AdditionalHeaders(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}

	config, err := pangea.NewConfig(
		option.WithBaseURLTemplate(baseURL),
		option.WithToken("my API token"),
	)
	assert.NoError(t, err)

	client, err := audit.New(config, audit.DisableEventVerification())
	assert.NoError(t, err)

	_, err = client.LogRequest(context.TODO(), audit.LogRequest{
		BaseRequest: pangea.BaseRequest{
			AdditionalHeaders: map[string]string{
				"Test-Header": "test-value",
			},
		},
		LogEvent: audit.LogEvent{
			Event: &audit.StandardEvent{Message: "test"},
		},
		Verbose: false,
	})
	assert.NoError(t, err)
}

func TestSearch_AdditionalHeaders(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}

	config, err := pangea.NewConfig(
		option.WithBaseURLTemplate(baseURL),
		option.WithToken("my API token"),
	)
	assert.NoError(t, err)

	client, err := audit.New(config, audit.DisableEventVerification())
	assert.NoError(t, err)

	_, err = client.Search(context.TODO(), &audit.SearchInput{
		BaseRequest: pangea.BaseRequest{
			AdditionalHeaders: map[string]string{
				"Test-Header": "test-value",
			},
		},
		Query: "message:test",
	})
	assert.NoError(t, err)
}

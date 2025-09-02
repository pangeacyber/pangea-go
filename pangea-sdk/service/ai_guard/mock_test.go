//go:build mock

package ai_guard_test

import (
	"context"
	"os"
	"testing"
	"time"

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

	_, err = client.GuardText(context.TODO(), &ai_guard.TextGuardRequest{Messages: []map[string]interface{}{
		{
			"role":    "user",
			"content": "what was pangea?",
		},
	}})
	assert.NoError(t, err)
}

func TestGuard(t *testing.T) {
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

	response, err := client.Guard(context.TODO(), ai_guard.GuardRequest{Input: map[string]any{
		"messages": []ai_guard.MultimodalMessage{
			{
				Role: "user",
				Content: ai_guard.MultimodalContent{
					OfString: pangea.P("what was pangea?"),
				},
			},
		}}})
	assert.NoError(t, err)
	assert.NotNil(t, response.Result.Detectors)

	response, err = client.Guard(context.TODO(), ai_guard.GuardRequest{
		Input: map[string]any{
			"messages": []ai_guard.MultimodalMessage{
				{
					Role: "user",
					Content: ai_guard.MultimodalContent{
						OfArrayOfContent: []ai_guard.MultimodalContentInner{
							{
								TextContent: &ai_guard.TextContent{
									Type: "text",
									Text: "what was pangea?",
								},
							},
							{
								ImageContent: &ai_guard.ImageContent{
									Type:     "image",
									ImageSrc: "https://example.org/favicon.ico",
								},
							},
						},
					},
				},
			},
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, response.Result.Detectors)
}

func TestGuardAsync(t *testing.T) {
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

	response, err := client.Guard(context.TODO(), ai_guard.GuardRequest{
		Input: map[string]any{
			"messages": []ai_guard.MultimodalMessage{
				{
					Role: "user",
					Content: ai_guard.MultimodalContent{
						OfString: pangea.P("what was pangea?"),
					},
				}},
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, response.Result.Detectors)

	response, err = client.Guard(context.TODO(), ai_guard.GuardRequest{
		Input: map[string]any{
			"messages": []ai_guard.MultimodalMessage{
				{
					Role: "user",
					Content: ai_guard.MultimodalContent{
						OfArrayOfContent: []ai_guard.MultimodalContentInner{
							{
								TextContent: &ai_guard.TextContent{
									Type: "text",
									Text: "what was pangea?",
								},
							},
							{
								ImageContent: &ai_guard.ImageContent{
									Type:     "image",
									ImageSrc: "https://example.org/favicon.ico",
								},
							},
						},
					},
				}},
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, response.Result.Detectors)
}

func TestGetServiceConfig(t *testing.T) {
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

	_, err = client.GetServiceConfig(context.TODO(), ai_guard.GetServiceConfigParams{Id: "pci_hgxhjm6w2l72ujzoxtfugaf6her7d27w"})
	assert.NoError(t, err)
}

func TestCreateServiceConfig(t *testing.T) {
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

	_, err = client.CreateServiceConfig(context.TODO(), ai_guard.CreateServiceConfigParams{
		ServiceConfig: ai_guard.ServiceConfig{
			Id:   pangea.String("pci_hgxhjm6w2l72ujzoxtfugaf6her7d27w"),
			Name: pangea.String("my-service-config-name"),
		},
	})
	assert.NoError(t, err)
}

func TestUpdateServiceConfig(t *testing.T) {
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

	_, err = client.UpdateServiceConfig(context.TODO(), ai_guard.UpdateServiceConfigParams{
		ServiceConfig: ai_guard.ServiceConfig{
			Id:   pangea.String("pci_hgxhjm6w2l72ujzoxtfugaf6her7d27w"),
			Name: pangea.String("my-service-config-name"),
		},
	})
	assert.NoError(t, err)
}

func TestDeleteServiceConfig(t *testing.T) {
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

	_, err = client.DeleteServiceConfig(context.TODO(), ai_guard.DeleteServiceConfigParams{Id: "pci_hgxhjm6w2l72ujzoxtfugaf6her7d27w"})
	assert.NoError(t, err)
}

func TestListServiceConfigs(t *testing.T) {
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

	_, err = client.ListServiceConfigs(context.TODO(), ai_guard.ListServiceConfigsParams{
		Filter: &ai_guard.ServiceConfigListFilter{
			Id:        pangea.String("pci_hgxhjm6w2l72ujzoxtfugaf6her7d27w"),
			CreatedAt: pangea.P(time.Now()),
		},
		Last:    pangea.P("last"),
		Order:   pangea.P("asc"),
		OrderBy: pangea.P("created_at"),
		Size:    pangea.P(1),
	})
	assert.NoError(t, err)
}

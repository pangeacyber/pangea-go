//go:build mock

package prompt_guard_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/option"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/prompt_guard"
	"github.com/stretchr/testify/assert"
)

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
	client := prompt_guard.New(config)

	_, err = client.GetServiceConfig(context.TODO(), prompt_guard.GetServiceConfigParams{
		ServiceConfig: prompt_guard.ServiceConfig{
			Id:      pangea.String("id"),
			Version: pangea.String("version"),
		},
	})
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
	client := prompt_guard.New(config)

	_, err = client.CreateServiceConfig(context.TODO(), prompt_guard.CreateServiceConfigParams{
		ServiceConfig: prompt_guard.ServiceConfig{
			Id: pangea.String("id"),
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
	client := prompt_guard.New(config)

	_, err = client.UpdateServiceConfig(context.TODO(), prompt_guard.UpdateServiceConfigParams{
		ServiceConfig: prompt_guard.ServiceConfig{
			Id:      pangea.String("id"),
			Version: pangea.String("version"),
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
	client := prompt_guard.New(config)

	_, err = client.DeleteServiceConfig(context.TODO(), prompt_guard.DeleteServiceConfigParams{Id: "id"})
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
	client := prompt_guard.New(config)

	response, err := client.ListServiceConfigs(context.TODO(), prompt_guard.ListServiceConfigsParams{
		Filter: &prompt_guard.ServiceConfigListFilter{
			Id:        pangea.String("id"),
			CreatedAt: pangea.P(time.Now()),
		},
		Last:    pangea.P("last"),
		Order:   pangea.P("asc"),
		OrderBy: pangea.P("created_at"),
		Size:    pangea.P(1),
	})
	assert.NoError(t, err)
	assert.NotNil(t, response.Result)
	assert.NotNil(t, response.Result.Count)
	assert.NotNil(t, response.Result.Items)
}

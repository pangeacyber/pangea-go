// go:build integration
package audit_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/pangeacyber/go-pangea/internal/pangeatesting"
	"github.com/pangeacyber/go-pangea/pangea"
	"github.com/pangeacyber/go-pangea/service/audit"
	"github.com/stretchr/testify/assert"
)

func Test_Integration_Root(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := &pangea.Config{
		CfgToken: os.Getenv("AUDIT_INTEGRATION_CONFIG_TOKEN"),
	}
	cfg = cfg.Copy(pangeatesting.IntegrationConfig)
	client, _ := audit.New(cfg)

	input := &audit.RootInput{}
	out, _, err := client.Root(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, out)
	assert.NotNil(t, out.Data)
	assert.NotNil(t, out.Data.RootHash)
	assert.NotEmpty(t, *out.Data.RootHash)
	assert.NotNil(t, out.Data.TreeName)
	assert.NotEmpty(t, *out.Data.TreeName)
	assert.NotNil(t, out.Data.Size)
	assert.NotNil(t, out.Data.ConsistencyProof)
}

func Test_Integration_Search(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := &pangea.Config{
		CfgToken: os.Getenv("AUDIT_INTEGRATION_CONFIG_TOKEN"),
	}
	cfg = cfg.Copy(pangeatesting.IntegrationConfig)
	client, _ := audit.New(cfg)

	input := &audit.SearchInput{
		Query: pangea.String("message:test"),
	}
	out, _, err := client.Search(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, out)
	assert.NotNil(t, out.ID)
	assert.NotEmpty(t, out.ID)
	assert.NotNil(t, out.ExpiresAt)
	assert.NotNil(t, out.Count)
	assert.NotNil(t, out.Events)
}

func Test_Integration_SearchResults(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancelFn()

	cfg := &pangea.Config{
		CfgToken: os.Getenv("AUDIT_INTEGRATION_CONFIG_TOKEN"),
	}
	cfg = cfg.Copy(pangeatesting.IntegrationConfig)
	client, _ := audit.New(cfg)
	searchInput := &audit.SearchInput{
		Query: pangea.String("message:test"),
	}
	searchOut, _, err := client.Search(ctx, searchInput)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, searchOut)
	assert.NotNil(t, searchOut.ID)

	searchResultInput := &audit.SeachResultInput{
		ID: searchOut.ID,
	}
	out, _, err := client.SearchResults(ctx, searchResultInput)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, out.Events)
	assert.NotNil(t, out.Count)
}

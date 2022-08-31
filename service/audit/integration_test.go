// go:build integration && !unit
package audit_test

import (
	"context"
	"testing"
	"time"

	"github.com/pangeacyber/go-pangea/internal/pangeatesting"
	"github.com/pangeacyber/go-pangea/pangea"
	"github.com/pangeacyber/go-pangea/service/audit"
	"github.com/stretchr/testify/assert"
)

func auditIntegrationCfg(t *testing.T) *pangea.Config {
	t.Helper()
	cfgToken := pangeatesting.GetEnvVarOrSkip(t, "AUDIT_INTEGRATION_CONFIG_TOKEN")
	cfg := &pangea.Config{
		CfgToken: cfgToken,
	}
	return cfg.Copy(pangeatesting.IntegrationConfig(t))
}

func Test_Integration_Log(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg)

	input := &audit.LogInput{
		Event: &audit.Event{
			Message: pangea.String("Integration test msg"),
		},
		ReturnHash: pangea.Bool(true),
		Verbose:    pangea.Bool(true),
	}

	out, err := client.Log(ctx, input)

	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.Hash)
	assert.NotNil(t, out.Result.Event)
	assert.NotNil(t, out.Result.CanonicalEventBase64)
	assert.NotEmpty(t, *out.Result.Hash)
	// assert.NotEmpty(t, *out.Result.Event)
	assert.NotEmpty(t, *out.Result.CanonicalEventBase64)
}

func Test_Integration_LogWithSignature(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg)

	input := &audit.LogInput{
		Event: &audit.Event{
			Message: pangea.String("Integration test msg"),
		},
		ReturnHash: pangea.Bool(true),
		Verbose:    pangea.Bool(true),
	}

	out, err := client.Log(ctx, input)

	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.Hash)
	assert.NotNil(t, out.Result.Event)
	assert.NotNil(t, out.Result.CanonicalEventBase64)
	assert.NotEmpty(t, *out.Result.Hash)
	// assert.NotEmpty(t, *out.Result.Event)
	assert.NotEmpty(t, *out.Result.CanonicalEventBase64)
}

func Test_Integration_Root(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg)

	input := &audit.RootInput{}
	out, err := client.Root(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.Data)
	assert.NotNil(t, out.Result.Data.RootHash)
	assert.NotEmpty(t, *out.Result.Data.RootHash)
	assert.NotNil(t, out.Result.Data.TreeName)
	assert.NotEmpty(t, *out.Result.Data.TreeName)
	assert.NotNil(t, out.Result.Data.Size)
	// TODO: Fix test
	// assert.NotNil(t, out.Data.ConsistencyProof)
}

func Test_Integration_Search(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg)

	input := &audit.SearchInput{
		Query: pangea.String("message:test"),
	}
	out, err := client.Search(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.ID)
	assert.NotEmpty(t, out.Result.ID)
	assert.NotNil(t, out.Result.ExpiresAt)
	assert.NotNil(t, out.Result.Count)
	assert.NotNil(t, out.Result.Events)
}

func Test_Integration_SearchResults(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg)
	searchInput := &audit.SearchInput{
		Query: pangea.String("message:test"),
	}
	searchOut, err := client.Search(ctx, searchInput)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, searchOut.Result)
	assert.NotNil(t, searchOut.Result.ID)

	searchResultInput := &audit.SearchResultInput{
		ID: searchOut.Result.ID,
	}
	out, err := client.SearchResults(ctx, searchResultInput)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, out.Result.Events)
	assert.NotNil(t, out.Result.Count)
}

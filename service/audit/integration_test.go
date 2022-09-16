// go:build integration && !unit
package audit_test

import (
	"context"
	"fmt"
	"math/rand"
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
	assert.NoError(t, err)
	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.Hash)
	assert.NotNil(t, out.Result.EventEnvelope)
	assert.NotNil(t, out.Result.EventEnvelope.Event)
	assert.NotNil(t, out.Result.EventEnvelope.Event.Message)
	assert.NotEmpty(t, *out.Result.Hash)
	assert.NotNil(t, out.Result.CanonicalEnvelopeBase64)
	assert.NotEmpty(t, *out.Result.CanonicalEnvelopeBase64)
}

func Test_Integration_Log_NoVerbose(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg)

	input := &audit.LogInput{
		Event: &audit.Event{
			Message: pangea.String("Integration test msg"),
		},
		ReturnHash: pangea.Bool(true),
		Verbose:    pangea.Bool(false),
	}

	out, err := client.Log(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.Hash)
	assert.NotEmpty(t, *out.Result.Hash)
	assert.Nil(t, out.Result.EventEnvelope)
	assert.Nil(t, out.Result.CanonicalEnvelopeBase64)
}

func Test_Integration_Log_Silent(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg)

	input := &audit.LogInput{
		Event: &audit.Event{
			Message: pangea.String("Integration test msg"),
		},
		ReturnHash: pangea.Bool(false),
		Verbose:    pangea.Bool(false),
	}

	out, err := client.Log(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, out.Result)
	assert.Nil(t, out.Result.Hash)
	assert.Nil(t, out.Result.EventEnvelope)
	assert.Nil(t, out.Result.CanonicalEnvelopeBase64)
}

// Fails because empty message
func Test_Integration_Log_Error(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg)

	input := &audit.LogInput{
		Event: &audit.Event{
			Message: pangea.String(""),
		},
		ReturnHash: pangea.Bool(true),
		Verbose:    pangea.Bool(true),
	}

	out, err := client.Log(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, out)
	err = err.(*pangea.APIError)
	apiErr := err.(*pangea.APIError)
	assert.Equal(t, len(apiErr.PangeaErrors.Errors), 1)
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Code, "BelowMinLength")
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Detail, "'message' cannot have less than 1 characters")
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Source, "/event/message")
}

func Test_Integration_Signatures(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg,
		audit.WithLogSigningEnabled("./testdata/privkey"),
		audit.WithLogSignatureVerificationEnabled(),
	)

	msg := "sigtest" + "100"
	logInput := &audit.LogInput{
		Event: &audit.Event{
			Message: pangea.String(msg),
			Source:  pangea.String("Source"),
			Status:  pangea.String("Status"),
			Target:  pangea.String("Target"),
			Actor:   pangea.String("Actor"),
			Action:  pangea.String("Action"),
			New:     pangea.String("New"),
			Old:     pangea.String("Old"),
		},
		ReturnHash: pangea.Bool(true),
		Verbose:    pangea.Bool(true),
	}

	outLog, err := client.Log(ctx, logInput)
	assert.NoError(t, err)

	fmt.Println("Event signature: ", *outLog.Result.EventEnvelope.Signature)
	fmt.Println("Encoded public key: ", *outLog.Result.EventEnvelope.PublicKey)

	searchInput := &audit.SearchInput{
		Query:      pangea.String(fmt.Sprintf("message:%s", msg)),
		MaxResults: pangea.Int(1),
	}
	// signature verification is done inside search
	out, err := client.Search(ctx, searchInput)

	assert.NoError(t, err)
	assert.NotNil(t, out)
	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.Count)
	assert.Equal(t, 1, pangea.IntValue(out.Result.Count))
	assert.Equal(t, len(out.Result.Events), 1)
	assert.NotNil(t, out.Result.Events[0].EventEnvelope.Signature)
	assert.NotNil(t, out.Result.Events[0].EventEnvelope.PublicKey)
}

func Test_Integration_Root(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg)

	input := &audit.RootInput{}
	out, err := client.Root(ctx, input)
	assert.NoError(t, err)

	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.Data)
	assert.NotNil(t, out.Result.Data.RootHash)
	assert.NotEmpty(t, *out.Result.Data.RootHash)
	assert.NotNil(t, out.Result.Data.TreeName)
	assert.NotEmpty(t, *out.Result.Data.TreeName)
	assert.NotNil(t, out.Result.Data.Size)
}

func Test_Integration_Proof(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg, audit.WithLogProofVerificationEnabled())

	maxResults := 4
	limit := 2

	input := &audit.SearchInput{
		IncludeHash:            pangea.Bool(true),
		IncludeMembershipProof: pangea.Bool(true),
		IncludeRoot:            pangea.Bool(true),
		MaxResults:             pangea.Int(maxResults),
		Limit:                  pangea.Int(limit),
		Query:                  pangea.String(""),
	}
	out, err := client.Search(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, out)
	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.ID)
	assert.NotEmpty(t, out.Result.ID)
	assert.NotNil(t, out.Result.ExpiresAt)
	assert.Equal(t, maxResults, pangea.IntValue(out.Result.Count))
	assert.NotNil(t, out.Result.Root)
	assert.Equal(t, limit, len(out.Result.Events))
	for _, event := range out.Result.Events {
		assert.NotNil(t, event.Hash)
	}
}

func Test_Integration_Search(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg)

	maxResults := 4
	limit := 2

	input := &audit.SearchInput{
		MaxResults: pangea.Int(maxResults),
		Limit:      pangea.Int(limit),
		Query:      pangea.String(""),
	}
	out, err := client.Search(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.ID)
	assert.NotEmpty(t, out.Result.ID)
	assert.NotNil(t, out.Result.ExpiresAt)
	assert.Equal(t, maxResults, pangea.IntValue(out.Result.Count))
	assert.Equal(t, limit, len(out.Result.Events))
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
	assert.NoError(t, err)
	assert.NotNil(t, searchOut.Result)
	assert.NotNil(t, searchOut.Result.ID)

	searchResultInput := &audit.SearchResultInput{
		ID: searchOut.Result.ID,
	}
	out, err := client.SearchResults(ctx, searchResultInput)
	assert.NoError(t, err)
	assert.NotNil(t, out.Result.Events)
	assert.NotNil(t, out.Result.Count)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

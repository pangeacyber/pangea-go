// go:build integration && !unit
package audit_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pangeacyber/go-pangea/internal/pangeatesting"
	pu "github.com/pangeacyber/go-pangea/internal/pangeautil"
	"github.com/pangeacyber/go-pangea/pangea"
	"github.com/pangeacyber/go-pangea/service/audit"
	"github.com/stretchr/testify/assert"
)

func auditIntegrationCfg(t *testing.T) *pangea.Config {
	t.Helper()
	cfgToken := pangeatesting.GetEnvVarOrSkip(t, "AUDIT_INTEGRATION_CONFIG_TOKEN")
	cfg := &pangea.Config{
		ConfigID: cfgToken,
	}
	return cfg.Copy(pangeatesting.IntegrationConfig(t))
}

func Test_Integration_Log(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg)

	event := audit.Event{
		Message: "Integration test msg",
	}

	out, err := client.Log(ctx, event, true, true)
	assert.NoError(t, err)
	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.Hash)
	assert.NotNil(t, out.Result.EventEnvelope)
	assert.NotNil(t, out.Result.EventEnvelope.Event)
	assert.NotNil(t, out.Result.EventEnvelope.Event.Message)
	assert.NotEmpty(t, out.Result.Hash)
	assert.NotNil(t, out.Result.CanonicalEnvelopeBase64)
	assert.NotEmpty(t, out.Result.CanonicalEnvelopeBase64)
}

func Test_Integration_Log_NoVerbose(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg)

	event := audit.Event{
		Message: "Integration test msg",
	}

	out, err := client.Log(ctx, event, false, true)
	assert.NoError(t, err)
	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.Hash)
	assert.NotEmpty(t, out.Result.Hash)
	assert.Nil(t, out.Result.EventEnvelope)
	assert.Empty(t, out.Result.CanonicalEnvelopeBase64)
}

func Test_Integration_Log_Silent(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg)

	event := audit.Event{
		Message: "Integration test msg",
	}

	out, err := client.Log(ctx, event, false, false)
	assert.NoError(t, err)
	assert.NotNil(t, out.Result)
	assert.Empty(t, out.Result.Hash)
	assert.Nil(t, out.Result.EventEnvelope)
	assert.Empty(t, out.Result.CanonicalEnvelopeBase64)
}

func Test_Integration_Log_Error_BadAuthToken(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	cfg.Token = "notavalidtoken"
	client, _ := audit.New(cfg)

	event := audit.Event{
		Message: "Integration test msg",
	}

	out, err := client.Log(ctx, event, true, true)
	assert.Error(t, err)
	assert.Nil(t, out)
	apiErr := err.(*pangea.APIError)
	assert.Equal(t, apiErr.Err.Error(), "API error: Not authorized to access this resource.")
}

// Not a valid config ID
func Test_Integration_Log_Error_BadConfidID(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	cfg.ConfigID = "notavalidid"
	client, _ := audit.New(cfg)

	event := audit.Event{
		Message: "Integration test msg",
	}

	out, err := client.Log(ctx, event, true, true)
	assert.Error(t, err)
	assert.Nil(t, out)
	apiErr := err.(*pangea.APIError)
	assert.Equal(t, apiErr.Err.Error(), "API error: Missing Config ID, you can provide using the config header X-Pangea-audit-Config-Id or adding a token scope `service:audit:*:config:r`.")
}

// Fails because empty message
func Test_Integration_Log_Error_EmptyMessage(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg)

	event := audit.Event{
		Message: "",
	}

	out, err := client.Log(ctx, event, true, true)

	assert.Error(t, err)
	assert.Nil(t, out)
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
		audit.WithLogProofVerificationEnabled(),
	)

	ts := pu.PangeaTimestamp(time.Date(2022, time.Month(11), 27, 12, 23, 37, 123456, time.UTC))

	msg := "sigtest" + "100"
	event := audit.Event{
		Message:   msg,
		Source:    "Source",
		Status:    "Status",
		Target:    "Target",
		Actor:     "Actor",
		Action:    "Action",
		New:       "New",
		Old:       "Old",
		Timestamp: pangea.PangeaTime(ts),
	}

	_, err := client.Log(ctx, event, true, true)
	assert.NoError(t, err)

	searchInput := &audit.SearchInput{
		Query:      fmt.Sprintf("message:%s", msg),
		MaxResults: 1,
	}
	// signature verification is done inside search
	out, err := client.Search(ctx, searchInput)

	assert.NoError(t, err)
	assert.NotNil(t, out)
	assert.NotNil(t, out.Result)
	assert.Equal(t, out.Result.Count, 1)
	assert.Equal(t, len(out.Result.Events), 1)
	assert.NotNil(t, out.Result.Events[0].EventEnvelope.Signature)
	assert.Equal(t, *out.Result.Events[0].EventEnvelope.Signature, "6qQGLhiIfWRqrPpcoVXFhtAxKr4iqzU5MT0iJm77ky3DWif6YS5PS4k/5CQ7rogwT09gdo8Dx2Ak8wrTYW1XBA==")
	assert.NotNil(t, out.Result.Events[0].EventEnvelope.PublicKey)
	assert.Equal(t, *out.Result.Events[0].EventEnvelope.PublicKey, "lvOyDMpK2DQ16NI8G41yINl01wMHzINBahtDPoh4+mE=")
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
	assert.NotEmpty(t, out.Result.Data.RootHash)
	assert.NotEmpty(t, out.Result.Data.TreeName)
	assert.NotEmpty(t, out.Result.Data.Size)
}

func Test_Integration_Proof(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg, audit.WithLogProofVerificationEnabled())

	maxResults := 4
	limit := 2

	input := &audit.SearchInput{
		IncludeHash:            true,
		IncludeMembershipProof: true,
		IncludeRoot:            true,
		MaxResults:             maxResults,
		Limit:                  limit,
		Query:                  "",
	}
	out, err := client.Search(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, out)
	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.ID)
	assert.NotEmpty(t, out.Result.ID)
	assert.NotNil(t, out.Result.ExpiresAt)
	assert.Equal(t, maxResults, out.Result.Count)
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
		MaxResults:             maxResults,
		Limit:                  limit,
		Query:                  "",
		IncludeHash:            true,
		IncludeMembershipProof: true,
		IncludeRoot:            true,
	}
	out, err := client.Search(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.ID)
	assert.NotEmpty(t, out.Result.ID)
	assert.NotNil(t, out.Result.ExpiresAt)
	assert.Equal(t, maxResults, out.Result.Count)
	assert.Equal(t, limit, len(out.Result.Events))
}

func Test_Integration_SearchResults(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg)
	searchInput := &audit.SearchInput{
		Query:      "message:test",
		MaxResults: 4,
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

func Test_Integration_SearchAll(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg)
	searchInput := &audit.SearchInput{
		Query:                  "message:Integration test msg",
		IncludeRoot:            true,
		IncludeMembershipProof: true,
		Limit:                  2,
	}
	root, se, err := audit.SearchAll(ctx, client, searchInput)

	assert.NoError(t, err)
	assert.NotNil(t, root)
	assert.NotEmpty(t, se)

	ve := se.VerifiableRecords()
	assert.NotEmpty(t, ve)

}

func Test_Integration_SearchAllAndValidate(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg)
	searchInput := &audit.SearchInput{
		Query:                  "message:Integration test msg",
		IncludeMembershipProof: true,
		IncludeRoot:            true,
	}

	root, se, err := audit.SearchAllAndValidate(ctx, client, searchInput)

	assert.NoError(t, err)
	assert.NotNil(t, root)
	assert.NotEmpty(t, se)

}

// go:build integration && !unit
package audit_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/internal/pangeatesting"
	pu "github.com/pangeacyber/pangea-go/pangea-sdk/internal/pangeautil"
	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/service/audit"
	"github.com/stretchr/testify/assert"
)

const (
	ACTOR              = "go-sdk"
	MSG_NO_SIGNED      = "test-message"
	MSG_JSON           = "JSON-message"
	STATUS_SIGNED      = "signed"
	MSG_SIGNED         = "sign-test"
	STATUS_NO_SIGNED   = "no-signed"
	testingEnvironment = pangeatesting.Live
)

func auditIntegrationCfg(t *testing.T) *pangea.Config {
	t.Helper()
	return pangeatesting.IntegrationConfig(t, testingEnvironment)
}

func Test_Integration_Log_NoVerbose(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg)

	event := audit.Event{
		Message: MSG_NO_SIGNED,
		Actor:   ACTOR,
		Status:  MSG_NO_SIGNED,
	}

	out, err := client.Log(ctx, event, false)
	assert.NoError(t, err)
	assert.NotNil(t, out.Result)
	assert.NotEmpty(t, out.Result.Hash)
	assert.Nil(t, out.Result.EventEnvelope)
	assert.Nil(t, out.Result.ConsistencyProof)
	assert.Nil(t, out.Result.MembershipProof)
	assert.Equal(t, out.Result.ConcistencyVerification, audit.NotVerified)
	assert.Equal(t, out.Result.MembershipVerification, audit.NotVerified)
	assert.Equal(t, out.Result.SignatureVerification, audit.NotVerified)
}

func Test_Integration_Log_VerboseNoVerify(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg)

	event := audit.Event{
		Message: MSG_NO_SIGNED,
		Actor:   ACTOR,
		Status:  MSG_NO_SIGNED,
	}

	out, err := client.Log(ctx, event, true)
	assert.NoError(t, err)
	assert.NotNil(t, out.Result)
	assert.NotEmpty(t, out.Result.Hash)
	assert.NotNil(t, out.Result.EventEnvelope)
	assert.NotNil(t, out.Result.EventEnvelope.Event)
	assert.NotNil(t, out.Result.EventEnvelope.Event.Message)
	assert.Nil(t, out.Result.ConsistencyProof)
	assert.NotNil(t, out.Result.MembershipProof)
	assert.Equal(t, out.Result.ConcistencyVerification, audit.NotVerified)
	assert.Equal(t, out.Result.MembershipVerification, audit.NotVerified)
	assert.Equal(t, out.Result.SignatureVerification, audit.NotVerified)
}

func Test_Integration_Log_TenantID(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg, audit.WithTenantID("mytenantid"))

	event := audit.Event{
		Message: MSG_NO_SIGNED,
		Actor:   ACTOR,
		Status:  MSG_NO_SIGNED,
	}

	out, err := client.Log(ctx, event, true)
	assert.NoError(t, err)
	assert.NotNil(t, out.Result)
	assert.NotEmpty(t, out.Result.Hash)
	assert.NotNil(t, out.Result.EventEnvelope)
	assert.NotNil(t, out.Result.EventEnvelope.Event)
	assert.NotNil(t, out.Result.EventEnvelope.Event.Message)
	assert.Nil(t, out.Result.ConsistencyProof)
	assert.NotNil(t, out.Result.MembershipProof)
	assert.Equal(t, out.Result.ConcistencyVerification, audit.NotVerified)
	assert.Equal(t, out.Result.MembershipVerification, audit.NotVerified)
	assert.Equal(t, out.Result.SignatureVerification, audit.NotVerified)
	assert.Equal(t, *out.Result.EventEnvelope.Event.TenantID, "mytenantid")
}

func Test_Integration_Log_VerboseAndVerify(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg, audit.WithLogProofVerificationEnabled())

	event := audit.Event{
		Message: MSG_NO_SIGNED,
		Actor:   ACTOR,
		Status:  MSG_NO_SIGNED,
	}

	out, err := client.Log(ctx, event, true)
	assert.NoError(t, err)
	assert.NotNil(t, out.Result)
	assert.NotEmpty(t, out.Result.Hash)
	assert.NotNil(t, out.Result.EventEnvelope)
	assert.NotNil(t, out.Result.EventEnvelope.Event)
	assert.NotNil(t, out.Result.EventEnvelope.Event.Message)
	assert.Nil(t, out.Result.ConsistencyProof)
	assert.NotNil(t, out.Result.MembershipProof)
	assert.Equal(t, out.Result.ConcistencyVerification, audit.NotVerified) // First log cant be consistency verified
	assert.Equal(t, out.Result.MembershipVerification, audit.Success)
	assert.Equal(t, out.Result.SignatureVerification, audit.NotVerified)

	out, err = client.Log(ctx, event, true)
	assert.NoError(t, err)
	assert.NotNil(t, out.Result)
	assert.NotEmpty(t, out.Result.Hash)
	assert.NotNil(t, out.Result.EventEnvelope)
	assert.NotNil(t, out.Result.EventEnvelope.Event)
	assert.NotNil(t, out.Result.EventEnvelope.Event.Message)
	assert.NotNil(t, out.Result.ConsistencyProof)
	assert.NotNil(t, out.Result.MembershipProof)
	assert.Equal(t, out.Result.ConcistencyVerification, audit.Success) // Second log can be verified
	assert.Equal(t, out.Result.MembershipVerification, audit.Success)
	assert.Equal(t, out.Result.SignatureVerification, audit.NotVerified)
}

func Test_Integration_Local_Signatures(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg,
		audit.WithLogLocalSigning("./testdata/privkey"),
		audit.WithLogProofVerificationEnabled(),
	)

	ts := pu.PangeaTimestamp(time.Date(2022, time.Month(11), 27, 12, 23, 37, 123456, time.UTC))

	event := audit.Event{
		Message:   MSG_SIGNED,
		Source:    "Source",
		Status:    STATUS_SIGNED,
		Target:    "Target",
		Actor:     ACTOR,
		Action:    "Action",
		New:       "New",
		Old:       "Old",
		Timestamp: pangea.PangeaTime(ts),
	}

	out, err := client.Log(ctx, event, true)
	assert.NoError(t, err)
	assert.NotNil(t, out)
	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.EventEnvelope.Signature)
	assert.NotNil(t, out.Result.EventEnvelope.PublicKey)
	assert.Equal(t, *out.Result.EventEnvelope.PublicKey, "lvOyDMpK2DQ16NI8G41yINl01wMHzINBahtDPoh4+mE=")
	assert.Equal(t, out.Result.SignatureVerification, audit.Success)
}

func Test_Integration_Local_Signatures_and_TenantID(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg,
		audit.WithLogLocalSigning("./testdata/privkey"),
		audit.WithLogProofVerificationEnabled(),
		audit.WithTenantID("mytenantid"),
	)

	ts := pu.PangeaTimestamp(time.Date(2022, time.Month(11), 27, 12, 23, 37, 123456, time.UTC))

	event := audit.Event{
		Message:   MSG_SIGNED,
		Source:    "Source",
		Status:    STATUS_SIGNED,
		Target:    "Target",
		Actor:     ACTOR,
		Action:    "Action",
		New:       "New",
		Old:       "Old",
		Timestamp: pangea.PangeaTime(ts),
	}

	out, err := client.Log(ctx, event, true)
	assert.NoError(t, err)
	assert.NotNil(t, out)
	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.EventEnvelope.Signature)
	assert.NotNil(t, out.Result.EventEnvelope.PublicKey)
	assert.Equal(t, *out.Result.EventEnvelope.PublicKey, "lvOyDMpK2DQ16NI8G41yINl01wMHzINBahtDPoh4+mE=")
	assert.Equal(t, out.Result.SignatureVerification, audit.Success)
	assert.Equal(t, *out.Result.EventEnvelope.Event.TenantID, "mytenantid")
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

func Test_Integration_Root_2(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg)
	TreeSize := 1

	input := &audit.RootInput{
		TreeSize: TreeSize,
	}
	out, err := client.Root(ctx, input)
	assert.NoError(t, err)

	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.Data)
	assert.NotEmpty(t, out.Result.Data.RootHash)
	assert.NotEmpty(t, out.Result.Data.TreeName)
	assert.Equal(t, TreeSize, out.Result.Data.Size)
}

func Test_Integration_Search_Results_NoVerify(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg)
	maxResults := 10
	limit := 2

	input := &audit.SearchInput{
		MaxResults: maxResults,
		Limit:      limit,
		Order:      "desc",
		Query:      "message:",
		Verbose:    pangea.Bool(false),
	}

	outSearch, err := client.Search(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, outSearch.Result)
	assert.NotNil(t, outSearch.Result.ID)
	assert.NotEmpty(t, outSearch.Result.ID)
	assert.NotNil(t, outSearch.Result.ExpiresAt)
	assert.LessOrEqual(t, outSearch.Result.Count, maxResults)
	assert.Greater(t, outSearch.Result.Count, 0)
	assert.Equal(t, limit, len(outSearch.Result.Events))

	for _, e := range outSearch.Result.Events {
		assert.Nil(t, e.MembershipProof)
		assert.Nil(t, e.Published)
		assert.Nil(t, e.LeafIndex)
		assert.Equal(t, audit.NotVerified, e.ConsistencyVerification)
		assert.Equal(t, audit.NotVerified, e.MembershipVerification)
	}

	// Test results
	resultsLimit := 2
	searchResultInput := &audit.SearchResultInput{
		ID:    outSearch.Result.ID,
		Limit: resultsLimit,
	}
	outResults, err := client.SearchResults(ctx, searchResultInput)
	assert.NoError(t, err)
	assert.LessOrEqual(t, outResults.Result.Count, maxResults)
	assert.Greater(t, outResults.Result.Count, 0)
	assert.Equal(t, resultsLimit, len(outResults.Result.Events))
	for _, e := range outResults.Result.Events {
		assert.Nil(t, e.MembershipProof)
		assert.Nil(t, e.Published)
		assert.Nil(t, e.LeafIndex)
		assert.Equal(t, audit.NotVerified, e.ConsistencyVerification)
		assert.Equal(t, audit.NotVerified, e.MembershipVerification)
		assert.Equal(t, "NotVerified", e.MembershipVerification.String())
	}

}

func Test_Integration_Search_Results_Verify(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg, audit.WithLogProofVerificationEnabled())
	maxResults := 10
	limit := 2

	input := &audit.SearchInput{
		Query:      fmt.Sprintf("message:%s status:%s actor: %s", MSG_SIGNED, STATUS_SIGNED, ACTOR),
		MaxResults: maxResults,
		Order:      "asc",
		Limit:      limit,
	}

	outSearch, err := client.Search(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, outSearch)
	assert.NotNil(t, outSearch.Result)
	assert.NotNil(t, outSearch.Result.ID)
	assert.NotEmpty(t, outSearch.Result.ID)
	assert.NotNil(t, outSearch.Result.ExpiresAt)
	assert.LessOrEqual(t, outSearch.Result.Count, maxResults)
	assert.Greater(t, outSearch.Result.Count, 0)
	assert.Equal(t, limit, len(outSearch.Result.Events))
	for _, e := range outSearch.Result.Events {
		assert.NotEmpty(t, e.MembershipProof)
		assert.Equal(t, audit.Success, e.MembershipVerification)
		assert.Equal(t, audit.Success, e.SignatureVerification)
		assert.Equal(t, "Success", e.SignatureVerification.String())
	}

	// Test results
	resultsLimit := 2
	searchResultInput := &audit.SearchResultInput{
		ID:    outSearch.Result.ID,
		Limit: resultsLimit,
	}
	outResults, err := client.SearchResults(ctx, searchResultInput)
	assert.NoError(t, err)
	assert.LessOrEqual(t, outResults.Result.Count, maxResults)
	assert.Greater(t, outResults.Result.Count, 0)
	assert.Equal(t, resultsLimit, len(outResults.Result.Events))
	for _, e := range outResults.Result.Events {
		assert.NotEmpty(t, e.MembershipProof)
		assert.Equal(t, audit.Success, e.MembershipVerification)
		assert.Equal(t, audit.Success, e.SignatureVerification)
	}

}

func Test_Integration_SearchAll(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg)
	searchInput := &audit.SearchInput{
		Query:   "message:test-message",
		Verbose: pangea.Bool(true),
		Limit:   10,
	}
	_, se, err := audit.SearchAll(ctx, client, searchInput)

	assert.NoError(t, err)
	assert.NotEmpty(t, se)
	ve := se.VerifiableRecords()
	assert.NotEmpty(t, ve)
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

	out, err := client.Log(ctx, event, true)
	assert.Error(t, err)
	assert.Nil(t, out)
	apiErr := err.(*pangea.APIError)
	assert.Equal(t, apiErr.Err.Error(), "API error: Not authorized to access this resource.")
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

	out, err := client.Log(ctx, event, true)

	assert.Error(t, err)
	assert.Nil(t, out)
	apiErr := err.(*pangea.APIError)
	assert.Equal(t, len(apiErr.PangeaErrors.Errors), 1)
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Code, "BelowMinLength")
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Detail, "'message' cannot have less than 1 characters")
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Source, "/event/message")
}

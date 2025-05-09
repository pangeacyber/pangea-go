//go:build integration

package audit_test

import (
	"context"
	"testing"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/internal/pangeatesting"
	pu "github.com/pangeacyber/pangea-go/pangea-sdk/v5/internal/pangeautil"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/audit"
	"github.com/stretchr/testify/assert"
)

const (
	ACTOR                          = "go-sdk"
	MSG_NO_SIGNED                  = "test-message"
	MSG_JSON                       = "JSON-message"
	STATUS_SIGNED                  = "signed"
	MSG_SIGNED                     = "sign-test"
	STATUS_NO_SIGNED               = "no-signed"
	ACTION_VAULT                   = "vault-sign"
	ACTION_LOCAL                   = "local-sign"
	MSG_CUSTOM_SCHEMA_NO_SIGNED    = "go-sdk-custom-schema-no-signed"
	MSG_CUSTOM_SCHEMA_SIGNED_LOCAL = "go-sdk-custom-schema-sign-local"
	MSG_CUSTOM_SCHEMA_SIGNED_VAULT = "go-sdk-custom-schema-sign-vault"
	LONG_FIELD                     = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed lacinia, orci eget commodo commodo non."
)

var testingEnvironment = pangeatesting.LoadTestEnvironment("audit", pangeatesting.Live)

var customSchemaEvent = pangeatesting.CustomSchemaEvent{
	Message:       MSG_CUSTOM_SCHEMA_NO_SIGNED,
	FieldInt:      1,
	FieldBool:     true,
	FieldStrShort: STATUS_NO_SIGNED,
	FieldStrLong:  LONG_FIELD,
	FieldTime:     pangea.PangeaTime(pu.PangeaTimestamp(time.Now())),
}

func auditIntegrationCfg(t *testing.T) *pangea.Config {
	t.Helper()
	return pangeatesting.IntegrationConfig(t, testingEnvironment)
}

func auditVaultIntegrationCfg(t *testing.T) *pangea.Config {
	t.Helper()
	return pangeatesting.IntegrationAuditVaultConfig(t, testingEnvironment)
}

func auditCustomSchemaCfg(t *testing.T) *pangea.Config {
	t.Helper()
	return pangeatesting.IntegrationCustomSchemaConfig(t, testingEnvironment)
}

func Test_Integration_Log_NoVerbose(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg)

	event := &audit.StandardEvent{
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

	event := &audit.StandardEvent{
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
	e := (out.Result.EventEnvelope.Event).(*audit.StandardEvent)
	assert.NotNil(t, e.Message)
	assert.Equal(t, e.Message, MSG_NO_SIGNED)
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

	event := &audit.StandardEvent{
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
	e := (out.Result.EventEnvelope.Event).(*audit.StandardEvent)
	assert.NotNil(t, e.Message)
	assert.Nil(t, out.Result.ConsistencyProof)
	assert.NotNil(t, out.Result.MembershipProof)
	assert.Equal(t, out.Result.ConcistencyVerification, audit.NotVerified)
	assert.Equal(t, out.Result.MembershipVerification, audit.NotVerified)
	assert.Equal(t, out.Result.SignatureVerification, audit.NotVerified)
	assert.Equal(t, e.TenantID, "mytenantid")
}

func Test_Integration_Log_VerboseAndVerify(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg, audit.WithLogProofVerificationEnabled())

	event := &audit.StandardEvent{
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
	e := (out.Result.EventEnvelope.Event).(*audit.StandardEvent)
	assert.NotNil(t, e.Message)
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
	e = (out.Result.EventEnvelope.Event).(*audit.StandardEvent)
	assert.NotNil(t, e.Message)
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

	event := &audit.StandardEvent{
		Message:   MSG_SIGNED,
		Source:    "Source",
		Status:    STATUS_SIGNED,
		Target:    "Target",
		Actor:     ACTOR,
		Action:    ACTION_LOCAL,
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
	assert.Equal(t, *out.Result.EventEnvelope.PublicKey, `{"algorithm":"ED25519","key":"-----BEGIN PUBLIC KEY-----\nMCowBQYDK2VwAyEAlvOyDMpK2DQ16NI8G41yINl01wMHzINBahtDPoh4+mE=\n-----END PUBLIC KEY-----\n"}`)
	assert.Equal(t, audit.Success, out.Result.SignatureVerification)
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

	event := &audit.StandardEvent{
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
	assert.Equal(t, *out.Result.EventEnvelope.PublicKey, `{"algorithm":"ED25519","key":"-----BEGIN PUBLIC KEY-----\nMCowBQYDK2VwAyEAlvOyDMpK2DQ16NI8G41yINl01wMHzINBahtDPoh4+mE=\n-----END PUBLIC KEY-----\n"}`)
	assert.Equal(t, out.Result.SignatureVerification, audit.Success)
	e := (out.Result.EventEnvelope.Event).(*audit.StandardEvent)
	assert.Equal(t, e.TenantID, "mytenantid")
}

// Custom schema tests
func Test_Integration_CustomSchema_Log_NoVerbose(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := auditCustomSchemaCfg(t)
	client, _ := audit.New(cfg)

	out, err := client.Log(ctx, &customSchemaEvent, false)
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

func Test_Integration_CustomSchema_Log_VerboseNoVerify(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := auditCustomSchemaCfg(t)
	client, _ := audit.New(cfg, audit.WithCustomSchema(pangeatesting.CustomSchemaEvent{}))

	out, err := client.Log(ctx, &customSchemaEvent, true)
	assert.NoError(t, err)
	assert.NotNil(t, out.Result)
	assert.NotEmpty(t, out.Result.Hash)
	assert.NotNil(t, out.Result.EventEnvelope)
	assert.NotNil(t, out.Result.EventEnvelope.Event)
	e := (out.Result.EventEnvelope.Event).(*pangeatesting.CustomSchemaEvent)
	assert.NotNil(t, e.Message)
	assert.Equal(t, MSG_CUSTOM_SCHEMA_NO_SIGNED, e.Message)
	assert.Nil(t, out.Result.ConsistencyProof)
	assert.NotNil(t, out.Result.MembershipProof)
	assert.Equal(t, out.Result.ConcistencyVerification, audit.NotVerified)
	assert.Equal(t, out.Result.MembershipVerification, audit.NotVerified)
	assert.Equal(t, out.Result.SignatureVerification, audit.NotVerified)
}

func Test_Integration_CustomSchema_Log_VerboseAndVerify(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := auditCustomSchemaCfg(t)
	client, _ := audit.New(cfg, audit.WithLogProofVerificationEnabled(), audit.WithCustomSchema(pangeatesting.CustomSchemaEvent{}))

	out, err := client.Log(ctx, &customSchemaEvent, true)
	assert.NoError(t, err)
	assert.NotNil(t, out.Result)
	assert.NotEmpty(t, out.Result.Hash)
	assert.NotNil(t, out.Result.EventEnvelope)
	assert.NotNil(t, out.Result.EventEnvelope.Event)
	e := (out.Result.EventEnvelope.Event).(*pangeatesting.CustomSchemaEvent)
	assert.NotNil(t, e.Message)
	assert.Equal(t, MSG_CUSTOM_SCHEMA_NO_SIGNED, e.Message)
	assert.Nil(t, out.Result.ConsistencyProof)
	assert.NotNil(t, out.Result.MembershipProof)
	assert.Equal(t, out.Result.ConcistencyVerification, audit.NotVerified) // First log cant be consistency verified
	assert.Equal(t, out.Result.MembershipVerification, audit.Success)
	assert.Equal(t, out.Result.SignatureVerification, audit.NotVerified)

	out, err = client.Log(ctx, &customSchemaEvent, true)
	assert.NoError(t, err)
	assert.NotNil(t, out.Result)
	assert.NotEmpty(t, out.Result.Hash)
	assert.NotNil(t, out.Result.EventEnvelope)
	assert.NotNil(t, out.Result.EventEnvelope.Event)
	e = (out.Result.EventEnvelope.Event).(*pangeatesting.CustomSchemaEvent)
	assert.NotNil(t, e.Message)
	assert.Equal(t, MSG_CUSTOM_SCHEMA_NO_SIGNED, e.Message)
	assert.NotNil(t, out.Result.ConsistencyProof)
	assert.NotNil(t, out.Result.MembershipProof)
	assert.Equal(t, out.Result.ConcistencyVerification, audit.Success) // Second log can be verified
	assert.Equal(t, out.Result.MembershipVerification, audit.Success)
	assert.Equal(t, out.Result.SignatureVerification, audit.NotVerified)
}

func Test_Integration_CustomSchema_Local_Signatures(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelFn()

	cfg := auditCustomSchemaCfg(t)
	client, _ := audit.New(cfg,
		audit.WithLogLocalSigning("./testdata/privkey"),
		audit.WithLogProofVerificationEnabled(),
		audit.WithCustomSchema(pangeatesting.CustomSchemaEvent{}),
	)

	var event = pangeatesting.CustomSchemaEvent{
		Message:       MSG_CUSTOM_SCHEMA_SIGNED_LOCAL,
		FieldInt:      1,
		FieldBool:     true,
		FieldStrShort: STATUS_NO_SIGNED,
		FieldStrLong:  LONG_FIELD,
		FieldTime:     pangea.PangeaTime(pu.PangeaTimestamp(time.Now())),
	}

	out, err := client.Log(ctx, &event, true)
	assert.NoError(t, err)
	assert.NotNil(t, out)
	assert.NotNil(t, out.Result)
	assert.NotNil(t, out.Result.EventEnvelope.Signature)
	assert.NotNil(t, out.Result.EventEnvelope.PublicKey)
	assert.Equal(t, *out.Result.EventEnvelope.PublicKey, `{"algorithm":"ED25519","key":"-----BEGIN PUBLIC KEY-----\nMCowBQYDK2VwAyEAlvOyDMpK2DQ16NI8G41yINl01wMHzINBahtDPoh4+mE=\n-----END PUBLIC KEY-----\n"}`)
	assert.Equal(t, audit.Success, out.Result.SignatureVerification)
	e := (out.Result.EventEnvelope.Event).(*pangeatesting.CustomSchemaEvent)
	assert.NotNil(t, e.Message)
	assert.Equal(t, MSG_CUSTOM_SCHEMA_SIGNED_LOCAL, e.Message)
}

func Test_Integration_Vault_Signatures(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelFn()

	cfg := auditVaultIntegrationCfg(t)
	client, _ := audit.New(cfg,
		audit.WithLogProofVerificationEnabled(),
	)

	ts := pu.PangeaTimestamp(time.Date(2022, time.Month(11), 27, 12, 23, 37, 123456, time.UTC))

	event := &audit.StandardEvent{
		Message:   MSG_SIGNED,
		Source:    "Source",
		Status:    STATUS_SIGNED,
		Target:    "Target",
		Actor:     ACTOR,
		Action:    ACTION_VAULT,
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
	assert.Equal(t, audit.Success, out.Result.SignatureVerification)
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
	client, _ := audit.New(cfg, audit.WithCustomSchema(pangeatesting.CustomSchemaEvent{}))
	maxResults := 5
	limit := 2

	input := &audit.SearchInput{
		MaxResults: maxResults,
		Limit:      limit,
		Order:      "desc",
		Query:      "message:\"\"",
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
	searchResultInput := &audit.SearchResultsInput{
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
	maxResults := 5
	limit := 2
	ct := time.Now().UTC()
	start := ct.Add(-3 * 24 * time.Hour).Truncate(1 * time.Microsecond)

	input := &audit.SearchInput{
		Query:      "message:" + MSG_SIGNED,
		MaxResults: maxResults,
		Order:      "asc",
		Limit:      limit,
		Start:      &start,
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
	searchResultInput := &audit.SearchResultsInput{
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
	ct := time.Now().UTC()
	start := ct.Add(-30 * 24 * time.Hour).Truncate(1 * time.Microsecond)

	searchInput := &audit.SearchInput{
		Query:   "message:" + MSG_NO_SIGNED,
		Verbose: pangea.Bool(true),
		Limit:   3,
		Order:   "asc",
		Start:   &start,
	}
	_, se, err := audit.SearchAll(ctx, client, searchInput)

	assert.NoError(t, err)
	assert.NotEmpty(t, se)
	ve := se.VerifiableRecords()
	assert.NotEmpty(t, ve)
}

// Custom schema tests
func Test_Integration_CustomSchema_Search_Results_NoVerify(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	cfg := auditCustomSchemaCfg(t)
	client, _ := audit.New(cfg, audit.WithCustomSchema(pangeatesting.CustomSchemaEvent{}))
	maxResults := 10
	limit := 2

	input := &audit.SearchInput{
		MaxResults: maxResults,
		Limit:      limit,
		Order:      "desc",
		Query:      "message:\"\"",
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
	searchResultInput := &audit.SearchResultsInput{
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

func Test_Integration_CustomSchema_Search_Results_Verify(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	cfg := auditCustomSchemaCfg(t)
	client, _ := audit.New(cfg, audit.WithLogProofVerificationEnabled(), audit.WithCustomSchema(pangeatesting.CustomSchemaEvent{}))
	maxResults := 10
	limit := 2

	input := &audit.SearchInput{
		Query:      "message:" + MSG_CUSTOM_SCHEMA_SIGNED_LOCAL,
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
	for _, se := range outSearch.Result.Events {
		assert.NotEmpty(t, se.MembershipProof)
		assert.Equal(t, audit.Success, se.MembershipVerification)
		assert.Equal(t, audit.Success, se.SignatureVerification)
		assert.Equal(t, "Success", se.SignatureVerification.String())
		e := (se.EventEnvelope.Event).(*pangeatesting.CustomSchemaEvent)
		assert.Equal(t, MSG_CUSTOM_SCHEMA_SIGNED_LOCAL, e.Message)
	}

	// Test results
	resultsLimit := 2
	searchResultInput := &audit.SearchResultsInput{
		ID:    outSearch.Result.ID,
		Limit: resultsLimit,
	}
	outResults, err := client.SearchResults(ctx, searchResultInput)
	assert.NoError(t, err)
	assert.LessOrEqual(t, outResults.Result.Count, maxResults)
	assert.Greater(t, outResults.Result.Count, 0)
	assert.Equal(t, resultsLimit, len(outResults.Result.Events))
	for _, se := range outResults.Result.Events {
		assert.NotEmpty(t, se.MembershipProof)
		assert.Equal(t, audit.Success, se.MembershipVerification)
		assert.Equal(t, audit.Success, se.SignatureVerification)
		e := (se.EventEnvelope.Event).(*pangeatesting.CustomSchemaEvent)
		assert.Equal(t, MSG_CUSTOM_SCHEMA_SIGNED_LOCAL, e.Message)
	}

}

func Test_Integration_CustomSchema_SearchAll(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg)
	searchInput := &audit.SearchInput{
		Query:   `message:""`,
		Verbose: pangea.Bool(true),
		Limit:   10,
		Order:   "asc",
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

	event := &audit.StandardEvent{
		Message: "Integration test msg",
	}

	out, err := client.Log(ctx, event, true)
	assert.Error(t, err)
	assert.Nil(t, out)
	apiErr := err.(*pangea.APIError)
	assert.Equal(t, apiErr.Err.Error(), "API error: Not authorized to access this resource.")
	assert.NotEmpty(t, apiErr.Error())
}

func Test_Integration_Multi_Config_1_Log(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := pangeatesting.IntegrationMultiConfigConfig(t, testingEnvironment)
	ConfigID := pangeatesting.GetConfigID(t, testingEnvironment, "audit", 1)
	client, _ := audit.New(cfg, audit.WithConfigID(ConfigID))

	event := &audit.StandardEvent{
		Message: MSG_NO_SIGNED,
		Actor:   ACTOR,
		Status:  STATUS_NO_SIGNED,
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

func Test_Integration_Multi_Config_2_Log(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := pangeatesting.IntegrationMultiConfigConfig(t, testingEnvironment)
	ConfigID := pangeatesting.GetConfigID(t, testingEnvironment, "audit", 2)
	client, _ := audit.New(cfg, audit.WithConfigID(ConfigID))

	event := &audit.StandardEvent{
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

func Test_Integration_Multi_Config_No_ConfigID(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := pangeatesting.IntegrationMultiConfigConfig(t, testingEnvironment)
	client, _ := audit.New(cfg, audit.WithConfigID(""))

	event := &audit.StandardEvent{
		Message: MSG_NO_SIGNED,
		Actor:   ACTOR,
		Status:  STATUS_NO_SIGNED,
	}

	out, err := client.Log(ctx, event, false)

	assert.Error(t, err)
	assert.Nil(t, out)
}

// Fails because empty message
// FIXME: Uncomment when fixed in backend
// func Test_Integration_Log_Error_EmptyMessage(t *testing.T) {
// 	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancelFn()

// 	cfg := auditIntegrationCfg(t)
// 	client, _ := audit.New(cfg)

// 	event := &audit.Event{
// 		Message: "",
// 	}

// 	out, err := client.Log(ctx, event, true)

// 	assert.Error(t, err)
// 	assert.Nil(t, out)
// 	apiErr := err.(*pangea.APIError)
// 	assert.Equal(t, len(apiErr.PangeaErrors.Errors), 1)
// 	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Code, "BelowMinLength")
// 	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Detail, "'message' cannot have less than 1 characters")
// 	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Source, "/event/message")
// }

func Test_Integration_Log_Bulk(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg)

	event := &audit.StandardEvent{
		Message: MSG_NO_SIGNED,
		Actor:   ACTOR,
		Status:  STATUS_NO_SIGNED,
	}

	out, err := client.LogBulk(ctx, []any{event, event}, true)
	assert.NoError(t, err)
	assert.NotNil(t, out.Result)
	for _, r := range out.Result.Results {
		assert.NotEmpty(t, r.Hash)
		assert.NotNil(t, r.EventEnvelope)
		assert.NotNil(t, r.EventEnvelope.Event)
		e := (r.EventEnvelope.Event).(*audit.StandardEvent)
		assert.NotNil(t, e.Message)
		assert.Equal(t, e.Message, MSG_NO_SIGNED)
		assert.Nil(t, r.ConsistencyProof)
		assert.Equal(t, r.ConcistencyVerification, audit.NotVerified)
		assert.Equal(t, r.MembershipVerification, audit.NotVerified)
		assert.Equal(t, r.SignatureVerification, audit.NotVerified)
	}
}

func Test_Integration_Log_Bulk_And_Sign(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg,
		audit.WithLogLocalSigning("./testdata/privkey"),
		audit.WithLogProofVerificationEnabled(),
	)

	event := &audit.StandardEvent{
		Message: MSG_SIGNED,
		Actor:   ACTOR,
		Status:  STATUS_SIGNED,
	}

	out, err := client.LogBulk(ctx, []any{event, event}, true)
	assert.NoError(t, err)
	assert.NotNil(t, out.Result)
	for _, r := range out.Result.Results {
		assert.NotEmpty(t, r.Hash)
		assert.NotNil(t, r.EventEnvelope)
		assert.NotNil(t, r.EventEnvelope.Event)
		e := (r.EventEnvelope.Event).(*audit.StandardEvent)
		assert.NotNil(t, e.Message)
		assert.Equal(t, e.Message, MSG_SIGNED)
		assert.Nil(t, r.ConsistencyProof)
		assert.Equal(t, r.ConcistencyVerification, audit.NotVerified)
		assert.Equal(t, r.MembershipVerification, audit.NotVerified)
		assert.Equal(t, r.SignatureVerification, audit.Success)
	}
}

func Test_Integration_Log_Bulk_Async(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg)

	event := &audit.StandardEvent{
		Message: MSG_NO_SIGNED,
		Actor:   ACTOR,
		Status:  MSG_NO_SIGNED,
	}

	out, err := client.LogBulkAsync(ctx, []any{event, event}, true)
	assert.NoError(t, err)
	assert.NotNil(t, out)
	assert.Nil(t, out.Result)
	assert.NotNil(t, out.AcceptedResult)
}

func Test_Integration_Download(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg, audit.WithCustomSchema(pangeatesting.CustomSchemaEvent{}))
	maxResults := 5
	limit := 2

	input := &audit.SearchInput{
		MaxResults: maxResults,
		Limit:      limit,
		Order:      "desc",
		Query:      "message:\"\"",
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

	outDownload, err := client.DownloadResults(ctx, &audit.DownloadRequest{
		ResultID: outSearch.Result.ID,
		Format:   audit.DFjson,
	})

	assert.NoError(t, err)
	assert.NotNil(t, outDownload)
	assert.NotNil(t, outDownload.Result)
	assert.NotEmpty(t, outDownload.Result.DestURL)

	af, err := client.DownloadFile(ctx, outDownload.Result.DestURL)

	af.Save(pangea.AttachedFileSaveInfo{})
}

type LogStreamEventData struct {
	ClientID     string  `json:"client_id"`
	Connection   *string `json:"connection,omitempty"`
	ConnectionID *string `json:"connection_id,omitempty"`
	Date         string  `json:"date"`
	Description  string  `json:"description"`
	IP           string  `json:"ip"`
	Strategy     *string `json:"strategy,omitempty"`
	StrategyType *string `json:"strategy_type,omitempty"`
	Type         string  `json:"type"`
	UserAgent    string  `json:"user_agent"`
	UserID       string  `json:"user_id"`
}

type LogStreamEvent struct {
	LogID string             `json:"log_id"`
	Data  LogStreamEventData `json:"data"`
}

type LogStreamRequest struct {
	pangea.BaseRequest

	Logs []LogStreamEvent `json:"logs"`
}

func Test_Integration_LogStream(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	cfg := pangeatesting.IntegrationMultiConfigConfig(t, testingEnvironment)
	cfgId := pangeatesting.GetConfigID(t, testingEnvironment, "audit", 3)
	client, _ := audit.New(cfg, audit.WithConfigID(cfgId))

	logStreamEvent := LogStreamEvent{
		LogID: "some log ID",
		Data: LogStreamEventData{
			ClientID:    "test client ID",
			Date:        "2024-03-29T17:26:50.193Z",
			Description: "Create a log stream",
			IP:          "127.0.0.1",
			Type:        "some_type",
			UserAgent:   "AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0",
			UserID:      "test user ID",
		},
	}
	input := LogStreamRequest{
		Logs: []LogStreamEvent{logStreamEvent},
	}
	response, err := client.LogStream(ctx, &input)

	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func Test_Integration_Export_Download(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	cfg := auditIntegrationCfg(t)
	client, _ := audit.New(cfg)

	exportRes, err := client.Export(ctx, &audit.ExportRequest{
		Start:   pangea.String("1d"),
		Verbose: pangea.Bool(false),
	})
	assert.NoError(t, err)
	assert.NotNil(t, exportRes)
	assert.Equal(t, "Accepted", pangea.StringValue(exportRes.Status))
	assert.NotEmpty(t, exportRes.RequestID)

	retry := 0
	for retry < 10 {
		_, err := client.PollResultByID(ctx, *exportRes.RequestID, &audit.DownloadResult{})
		if err == nil {
			break
		}

		// Wait until result should be ready
		time.Sleep(time.Duration(3 * time.Second))
		retry++
	}

	downloadRes, err := client.DownloadResults(ctx, &audit.DownloadRequest{
		RequestID: *exportRes.RequestID,
	})
	assert.NoError(t, err)
	assert.Equal(t, "Success", pangea.StringValue(downloadRes.Status))
	assert.NotEmpty(t, downloadRes.Result.DestURL)
}

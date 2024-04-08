package audit

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	pu "github.com/pangeacyber/pangea-go/pangea-sdk/v3/internal/pangeautil"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/internal/request"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/internal/signer"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
)

// @summary Log an entry
//
// @description Create a log entry in the Secure Audit Log.
//
// @operationId audit_post_v1_log
//
// @example
//
//	event := audit.Event{
//		Message: "hello world",
//	 }
//
//	logResponse, err := auditcli.Log(ctx, event, true)
func (a *audit) Log(ctx context.Context, event any, verbose bool) (*pangea.PangeaResponse[LogResult], error) {
	input, err := a.getLogRequest(event, verbose)
	if err != nil {
		return nil, err
	}

	resp, err := request.DoPost(ctx, a.Client, "v1/log", input, &LogResult{})
	if err != nil {
		return nil, err
	}

	err = a.processLogResult(ctx, resp.Result)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// @summary Log multiple entries
//
// @description Create multiple log entries in the Secure Audit Log.
//
// @operationId audit_post_v2_log
//
// @example
//
//	event := audit.Event{
//		Message: "hello world",
//	 }
//	events := []audit.Event{event}
//
//	logResponse, err := auditcli.LogBulk(ctx, events, true)
func (a *audit) LogBulk(ctx context.Context, events []any, verbose bool) (*pangea.PangeaResponse[LogBulkResult], error) {
	input, err := a.getLogBulkRequest(events, verbose)
	if err != nil {
		return nil, err
	}

	resp, err := request.DoPost(ctx, a.Client, "v2/log", input, &LogBulkResult{})
	if err != nil {
		return nil, err
	}

	err = a.processLogBulkResult(ctx, resp.Result)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// @summary Log multiple entries asynchronously
//
// @description Asynchronously create multiple log entries in the Secure Audit Log.
//
// @operationId audit_post_v2_log_async
//
// @example
//
//	event := audit.Event{
//		Message: "hello world",
//	 }
//	events := []audit.Event{event}
//
//	logResponse, err := auditcli.LogBulkAsync(ctx, events, true)
func (a *audit) LogBulkAsync(ctx context.Context, events []any, verbose bool) (*pangea.PangeaResponse[LogBulkResult], error) {
	input, err := a.getLogBulkRequest(events, verbose)
	if err != nil {
		return nil, err
	}

	resp, err := request.DoPostNoQueue(ctx, a.Client, "v2/log_async", input, &LogBulkResult{})
	if err != nil {
		ae, ok := err.(*pangea.AcceptedError)
		if ok {
			return &pangea.PangeaResponse[LogBulkResult]{
				Response:       ae.Response,
				AcceptedResult: &ae.AcceptedResult,
				Result:         nil,
			}, nil
		}

		return nil, err
	}

	err = a.processLogBulkResult(ctx, resp.Result)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (a *audit) getLogBulkRequest(events []any, verbose bool) (*LogBulkRequest, error) {
	var logEvents []*LogEvent

	for i := range events {
		logEvent, err := a.getLogEvent(events[i])
		if err != nil {
			return nil, err
		}
		logEvents = append(logEvents, logEvent)
	}

	input := &LogBulkRequest{
		Events: logEvents,
	}

	input.Verbose = verbose

	if a.verifyProofs {
		input.Verbose = true
	}

	return input, nil
}

func (a *audit) getLogRequest(event any, verbose bool) (*LogRequest, error) {
	// Overwrite tenant id if user set it on event
	if st, ok := event.(Tenanter); ok {
		if st.Tenant() == "" && a.tenantID != "" {
			st.SetTenant(a.tenantID)
		}
	}

	logEvent, err := a.getLogEvent(event)
	if err != nil {
		return nil, err
	}

	input := &LogRequest{
		LogEvent: *logEvent,
	}

	input.Verbose = verbose

	if a.verifyProofs {
		input.Verbose = true
		input.PrevRoot = a.lastUnpRootHash
	}

	return input, nil
}

func (a *audit) getLogEvent(event any) (*LogEvent, error) {
	// Overwrite tenant id if user set it on event
	if st, ok := event.(Tenanter); ok {
		if st.Tenant() == "" && a.tenantID != "" {
			st.SetTenant(a.tenantID)
		}
	}

	le := &LogEvent{
		Event: event,
	}

	if a.signer != nil {
		err := le.SignEvent(*a.signer, a.publicKeyInfo)
		if err != nil {
			return nil, err
		}
	}

	return le, nil
}

// @summary Search for events
//
// @description Search for events that match the provided search criteria.
//
// @operationId audit_post_v1_search
//
// @example
//
//	input := &audit.SearchInput{
//		Query:                  pangea.String("message:log-123"),
//		IncludeMembershipProof: pangea.Bool(true),
//	}
//
//	searchResponse, err := auditcli.Search(ctx, input)
func (a *audit) Search(ctx context.Context, input *SearchInput) (*pangea.PangeaResponse[SearchOutput], error) {
	if input == nil {
		return nil, errors.New("nil input")
	}

	if a.verifyProofs {
		// Need this info to verify
		input.Verbose = pangea.Bool(true)
	}

	resp, err := request.DoPost(ctx, a.Client, "v1/search", input, &SearchOutput{})
	if err != nil {
		return nil, err
	}

	err = a.processSearchEvents(ctx, resp.Result.Events, resp.Result.Root, resp.Result.UnpublishedRoot)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// @summary Search results
//
// @description Page through results from a previous search.
//
// @operationId audit_post_v1_results
//
// @example
//
//	input := &audit.SearchResultsInput{
//		ID: "pas_sqilrhruwu54uggihqj3aie24wrctakr",
//	}
//
//	res, err := auditcli.SearchResults(ctx, input)
func (a *audit) SearchResults(ctx context.Context, input *SearchResultsInput) (*pangea.PangeaResponse[SearchResultsOutput], error) {
	resp, err := request.DoPost(ctx, a.Client, "v1/results", input, &SearchResultsOutput{})
	if err != nil {
		return nil, err
	}

	err = a.processSearchEvents(ctx, resp.Result.Events, resp.Result.Root, resp.Result.UnpublishedRoot)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// @summary Tamperproof verification
//
// @description Returns current root hash and consistency proof.
//
// @operationId audit_post_v1_root
//
// @example
//
//	input := &audit.RootInput{
//		TreeSize: pangea.Int(10),
//	}
//
//	rootResponse, err := auditcli.Root(ctx, input)
func (a *audit) Root(ctx context.Context, input *RootInput) (*pangea.PangeaResponse[RootOutput], error) {
	return request.DoPost(ctx, a.Client, "v1/root", input, &RootOutput{})
}

// SearchAll is a helper function to return all the search results for a search with pages
func SearchAll(ctx context.Context, client Client, input *SearchInput) (*Root, SearchEvents, error) {
	resp, err := client.Search(ctx, input)
	if err != nil {
		return nil, nil, err
	}
	events := make(SearchEvents, 0, resp.Result.Count)
	events = append(events, resp.Result.Events...)
	for resp.Result.Count > len(events) {
		s := SearchResultsInput{
			ID: resp.Result.ID,
		}
		sOut, err := client.SearchResults(ctx, &s)
		if err != nil {
			return nil, nil, err
		}
		if len(sOut.Result.Events) == 0 {
			break
		}
		events = append(events, sOut.Result.Events...)
	}
	return resp.Result.Root, events, nil
}

func (a *audit) processLogBulkResult(ctx context.Context, br *LogBulkResult) error {
	for i := range br.Results {
		err := a.processLogResult(ctx, &br.Results[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *audit) processLogResult(ctx context.Context, log *LogResult) error {
	if log == nil {
		return nil
	}

	ee, err := a.newEventEnvelopeFromMap(log.RawEnvelope)
	if err != nil {
		return err
	}

	log.EventEnvelope = ee

	if !a.skipEventVerification {
		if VerifyHash(log.RawEnvelope, log.Hash) == Failed {
			return fmt.Errorf("audit: Failed hash verification of event. Hash: [%s]", log.Hash)
		}
	}

	nurh := log.UnpublishedRootHash
	if log.EventEnvelope != nil {
		log.SignatureVerification = log.EventEnvelope.VerifySignature()
	}

	if a.verifyProofs {
		if nurh != nil && log.MembershipProof != nil {
			res, _ := VerifyMembershipProof(*nurh, log.Hash, *log.MembershipProof)
			log.MembershipVerification = res

			if log.ConsistencyProof != nil && a.lastUnpRootHash != nil {
				b, _ := verifyConsistencyProof(*a.lastUnpRootHash, *nurh, *log.ConsistencyProof)
				if b {
					log.ConcistencyVerification = Success
				} else {
					log.ConcistencyVerification = Failed
				}
			}
		}
	}

	if nurh != nil {
		a.lastUnpRootHash = nurh
	}
	return nil
}

func (a *audit) processSearchEvents(ctx context.Context, events SearchEvents, root *Root, unpRoot *Root) error {
	var roots map[int]Root

	var err error
	for _, event := range events {
		event.EventEnvelope, err = a.newEventEnvelopeFromMap(event.RawEnvelope)
		if err != nil {
			return err
		}
	}

	if a.verifyProofs && root != nil {
		if a.rp == nil {
			a.rp = NewArweaveRootsProvider(root.TreeName)
		}
		treeSizes := treeSizes(root, events)
		roots = a.rp.UpdateRoots(ctx, treeSizes)
	}

	for _, event := range events {
		if !a.skipEventVerification {
			if VerifyHash(event.RawEnvelope, event.Hash) == Failed {
				return fmt.Errorf("audit: cannot verify hash of record. Hash: [%s]", event.Hash)
			}
			event.SignatureVerification = event.EventEnvelope.VerifySignature()
		}

		if a.verifyProofs {
			if event.Published != nil && *event.Published {
				event.VerifyMembershipProof(root)
				event.VerifyConsistencyProof(roots)
			} else {
				event.VerifyMembershipProof(unpRoot)
			}
		}
	}
	return nil
}

type LogEvent struct {
	// A structured event describing an auditable activity.
	Event any `json:"event"`

	// An optional client-side signature for forgery protection.
	// max len of 256 bytes
	Signature *string `json:"signature,omitempty"`

	// The base64-encoded ed25519 public key used for the signature, if one is provided
	PublicKey *string `json:"public_key,omitempty"`
}

type LogRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	LogEvent

	// If true, be verbose in the response; include root, membership and consistency proof, etc.
	// default: false
	Verbose bool `json:"verbose"`

	// Previous unpublished root
	PrevRoot *string `json:"prev_root,omitempty"`
}

type LogBulkRequest struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	Events []*LogEvent `json:"events"`

	// If true, be verbose in the response; include root, membership and consistency proof, etc.
	// default: false
	Verbose bool `json:"verbose"`
}

func (i *LogEvent) SignEvent(s signer.Signer, pki map[string]string) error {
	b, err := pu.CanonicalizeStruct(&i.Event)
	if err != nil {
		return err
	}

	signature, err := s.Sign(b)
	if err != nil {
		return err
	}

	if pki == nil {
		pki = make(map[string]string)
	}

	pk, err := s.PublicKey()
	if err != nil {
		return err
	}

	pki["key"] = pk
	pki["algorithm"] = s.GetAlgorithm()

	pkib, err := pu.CanonicalizeStruct(pki)
	if err != nil {
		return err
	}

	i.Signature = pangea.String(base64.StdEncoding.EncodeToString(signature))
	i.PublicKey = pangea.String(string(pkib))
	return nil
}

type StandardEvent struct {
	// Record who performed the auditable activity.
	// max len is 128 bytes
	// examples:
	// 	John Doe
	//  user-id
	//  DennisNedry@InGen.com
	Actor string `json:"actor,omitempty"`

	// The auditable action that occurred."
	// max len is 32 bytes
	// examples:
	// 	created
	//  deleted
	//  updated
	Action string `json:"action,omitempty"`

	// A message describing a detailed account of what happened.
	// This can be recorded as free-form text or as a JSON-formatted string.
	// Message is a required field.
	// max len of 65536 bytes
	Message string `json:"message"`

	// The value of a record after it was changed.
	// max len of 65536 bytes
	New string `json:"new,omitempty"`

	// The value of a record before it was changed.
	// max len of 65536 bytes
	Old string `json:"old,omitempty"`

	// Used to record the location from where an activity occurred.
	// max len of 128 bytes
	Source string `json:"source,omitempty"`

	// Record whether or not the activity was successful.
	// examples:
	//  failure
	//  success
	// max len of 32 bytes
	Status string `json:"status,omitempty"`

	// Used to record the specific record that was targeted by the auditable activity.
	// max len of 128 bytes
	Target string `json:"target,omitempty"`

	// An optional client-supplied timestamp.
	Timestamp *pu.PangeaTimestamp `json:"timestamp,omitempty"`

	// TenantID field
	TenantID string `json:"tenant_id,omitempty"`
}

func (e *StandardEvent) Tenant() string {
	return e.TenantID
}

func (e *StandardEvent) SetTenant(tid string) {
	e.TenantID = tid
}

type EventEnvelope struct {
	// A structured record describing that <actor> did <action> on <target>
	// changing it from <old> to <new> and the operation was <status>,
	// and/or a free-form <message>.
	Event any `json:"event"`

	// An optional client-side signature for forgery protection.
	// max len of 256 bytes
	Signature *string `json:"signature,omitempty"`

	// The base64-encoded ed25519 public key used for the signature, if one is provided
	PublicKey *string `json:"public_key,omitempty"`

	// A server-supplied timestamp.
	ReceivedAt *pu.PangeaTimestamp `json:"received_at,omitempty"`
}

func (a *audit) newEventEnvelopeFromMap(m map[string]any) (*EventEnvelope, error) {
	if m == nil {
		return nil, nil
	}

	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	var ee = EventEnvelope{}
	err = json.Unmarshal(b, &ee)
	if err != nil {
		return nil, err
	}

	if ee.Event != nil {
		b, err = json.Marshal(ee.Event)
		if err != nil {
			return nil, err
		}

		v := reflect.New(reflect.TypeOf(a.schema)).Interface()

		err = json.Unmarshal(b, &v)
		if err != nil {
			return nil, err
		}

		ee.Event = v
	}
	return &ee, nil
}

type EventVerification int

const (
	NotVerified EventVerification = iota
	Success
	Failed
)

type DownloadFormat string

const (
	DFjson DownloadFormat = "json" // JSON.
	DFcsv  DownloadFormat = "csv"  // CSV.
)

func (ev EventVerification) String() string {
	switch ev {
	case NotVerified:
		return "NotVerified"
	case Success:
		return "Success"
	case Failed:
		return "Failed"
	}
	return "unknown"
}

type LogResult struct {
	EventEnvelope *EventEnvelope

	RawEnvelope map[string]any `json:"envelope"`

	// The hash of the event data.
	// max len of 64 bytes
	Hash string `json:"hash"`

	UnpublishedRootHash     *string   `json:"unpublished_root,omitempty"`
	MembershipProof         *string   `json:"membership_proof,omitempty"`
	ConsistencyProof        *[]string `json:"consistency_proof,omitempty"`
	MembershipVerification  EventVerification
	ConcistencyVerification EventVerification
	SignatureVerification   EventVerification
}

type LogBulkResult struct {
	Results []LogResult `json:"results"`
}

type SearchInput struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	// Natural search string; list of keywords with optional `<option>:<value>` qualifiers.
	//
	// Query is a required field.
	//
	// The following optional qualifiers are supported:
	//	* action:
	//	* actor:
	//	* message:
	//	* new:
	//	* old:
	//	* status:
	//	* target:
	//
	// examples:
	//		actor:root target:/etc/shadow
	Query string `json:"query"`

	// Specify the sort order of the response. "asc" or "desc"
	Order string `json:"order,omitempty"`

	// Name of column to sort the results by.
	OrderBy string `json:"order_by,omitempty"`

	// The start of the time range to perform the search on.
	Start *time.Time `json:"start,omitempty"`

	// The end of the time range to perform the search on. All records up to the latest if left out.
	End *time.Time `json:"end,omitempty"`

	// Number of audit records to include from the first page of the results.
	Limit int `json:"limit,omitempty"`

	// Maximum number of results to return.
	// min 1 max 10000
	MaxResults int `json:"max_results,omitempty"`

	// If true include root, membership and consistency proof
	Verbose *bool `json:"verbose,omitempty"`

	// A list of keys to restrict the search results to. Useful for partitioning data available to the query string.
	SearchRestriction *SearchRestriction `json:"search_restriction,omitempty"`
}

type SearchRestriction struct {
	// A list of actors to restrict the search to.
	Actor []string `json:"actor,omitempty"`

	// A list of sources to restrict the search to.
	Source []string `json:"source,omitempty"`

	// A list of targets to restrict the search to.
	Target []string `json:"target,omitempty"`

	// A list of actions to restrict the search to.
	Action []string `json:"action,omitempty"`

	// A list of status to restrict the search to.
	Status []string `json:"status,omitempty"`
}

type SearchOutput struct {
	// Identifier to supply to search_results API to fetch/paginate through search results.
	// ID is always populated on a successful response.
	ID string `json:"id"`

	// The time when the results will no longer be available to page through via the results API.
	// ExpiresAt is always populated on a successful response.
	ExpiresAt *time.Time `json:"expires_at"`

	// The total number of results that were returned by the search.
	// Count is always populated on a successful response.
	Count int `json:"count"`

	// A list of matching audit records.
	// Events is always populated on a successful response.
	Events SearchEvents `json:"events"`

	// A root of a Merkle Tree
	Root *Root `json:"root,omitempty"`

	// A unpublished root of a Merkle Tree
	UnpublishedRoot *Root `json:"unpublished_root,omitempty"`
}

type SearchEvents []*SearchEvent

// VerifiableRecords returns a slice of records that can be verifiable by the published proof
func (events SearchEvents) VerifiableRecords() SearchEvents {
	evs := make(SearchEvents, 0)
	for _, event := range events {
		if event.IsVerifiable() {
			evs = append(evs, event)
		}
	}
	return evs
}

type SearchEvent struct {
	// Include Event data and security information
	EventEnvelope *EventEnvelope

	RawEnvelope map[string]any `json:"envelope"`

	// The record's hash
	// len of 64 bytes
	Hash string `json:"hash"`

	// The index of the leaf of the Merkle Tree where this record was inserted.
	LeafIndex *int `json:"leaf_index"`

	// A cryptographic proof that the record has been persisted in the log.
	MembershipProof *string `json:"membership_proof"`

	Published *bool `json:"published"`

	MembershipVerification  EventVerification
	ConsistencyVerification EventVerification
	SignatureVerification   EventVerification
}

// IsVerifiable checks if a record can be verifiable with the published proof
func (event *SearchEvent) IsVerifiable() bool {
	return event.LeafIndex != nil && *event.LeafIndex >= 0
}

func (ee *SearchEvent) VerifyMembershipProof(root *Root) {
	if root == nil || ee.MembershipProof == nil {
		ee.MembershipVerification = NotVerified
	} else {
		res, err := VerifyMembershipProof(root.RootHash, ee.Hash, *ee.MembershipProof)
		if err != nil {
			ee.MembershipVerification = Failed
		} else {
			ee.MembershipVerification = res
		}
	}
}

func (ee *SearchEvent) VerifyConsistencyProof(publishedRoots map[int]Root) {
	if ee.Published == nil || !*ee.Published || ee.LeafIndex == nil {
		ee.ConsistencyVerification = NotVerified
		return
	}
	idx := *ee.LeafIndex
	if idx < 0 {
		ee.ConsistencyVerification = Failed
		return
	}
	current, ok := publishedRoots[idx]
	if !ok || current.ConsistencyProof == nil {
		ee.ConsistencyVerification = NotVerified
		return
	}
	previous, ok := publishedRoots[idx-1]
	if !ok {
		ee.ConsistencyVerification = NotVerified
		return
	}
	verified, err := verifyConsistencyProof(previous.RootHash, current.RootHash, *current.ConsistencyProof)
	if err != nil {
		ee.ConsistencyVerification = Failed
		return
	}
	if verified {
		ee.ConsistencyVerification = Success
	} else {
		ee.ConsistencyVerification = Failed
	}
}

func (ee *EventEnvelope) VerifySignature() EventVerification {
	// Both nil, so NotVerified
	if ee.Signature == nil && ee.PublicKey == nil {
		return NotVerified
	}

	// If just one nil, it's an error so Failed
	if ee.Signature == nil || ee.PublicKey == nil {
		return Failed
	}

	b, err := pu.CanonicalizeStruct(ee.Event)
	if err != nil {
		return NotVerified
	}

	sig, err := base64.StdEncoding.DecodeString(*ee.Signature)
	if err != nil {
		return Failed
	}

	publicKey, err := ee.getPublicKey()
	if err != nil {
		return Failed
	}

	v, err := signer.NewVerifierFromPubKey(publicKey)
	if err != nil {
		return Failed
	}

	if v != nil {
		ver, err := v.Verify(b, sig)
		if err != nil {
			return NotVerified
		}
		if ver {
			return Success
		} else {
			return Failed
		}
	}
	return NotVerified
}

func (ee EventEnvelope) getPublicKey() (string, error) {
	// Should never enter this case
	if ee.PublicKey == nil {
		return "", errors.New("public key field nil pointer")
	}

	pkinfo := make(map[string]any)
	err := json.Unmarshal([]byte(*ee.PublicKey), &pkinfo)
	if err != nil {
		return *ee.PublicKey, nil
	}

	val, ok := pkinfo["key"]
	if !ok {
		return "", errors.New("'key' field not present in json")
	}

	ret, ok := val.(string)
	if !ok {
		return "", errors.New("value is not a string")
	}

	return ret, nil
}

type SearchResultsInput struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	// A search results identifier returned by the search call
	// ID is a required field
	ID string `json:"id"`

	// Number of audit records to include from the first page of the results.
	Limit int `json:"limit,omitempty"`

	// Offset from the start of the result set to start returning results from.
	Offset *int `json:"offset,omitempty"`
}

type SearchResultsOutput struct {
	// The total number of results that were returned by the search.
	// Count is always populated on a successful response.
	Count int `json:"count"`

	// A list of matching audit records.
	// Events is always populated on a successful response.
	Events SearchEvents `json:"events"`

	// A root of a Merkle Tree
	Root *Root `json:"root"`

	// A unpublished root of a Merkle Tree
	UnpublishedRoot *Root `json:"unpublished_root"`
}

type RootInput struct {
	// Base request has ConfigID for multi-config projects
	pangea.BaseRequest

	// The size of the tree (the number of records)
	TreeSize int `json:"tree_size,omitempty"`
}

type Root struct {
	// The name of the Merkle Tree
	TreeName string `json:"tree_name"`

	// The size of the tree (the number of records)
	Size int `json:"size"`

	// The root hash
	// max len of 64 bytes
	RootHash string `json:"root_hash"`

	// The URL where this root has been published
	URL *string `json:"url"`

	// The date/time when this root was published
	PublishedAt *time.Time `json:"published_at"`

	// Consistency proof to verify that this root is a continuation of the previous one
	ConsistencyProof *[]string `json:"consistency_proof"`
}

type RootOutput struct {
	Data Root `json:"data"`
}

type DownloadRequest struct {
	pangea.BaseRequest

	// ID returned by the search API.
	ResultID string `json:"result_id"`

	// Format for the records.
	Format DownloadFormat `json:"format,omitempty"`
}

type DownloadResult struct {
	// URL where search results can be downloaded.
	DestURL string `json:"dest_url"`
}

// @summary Download search results
//
// @description Get all search results as a compressed (gzip) CSV file.
//
// @operationId audit_post_v1_download_results
//
// @example
//
//	response, err := client.DownloadResults(ctx, &audit.DownloadRequest{
//		ResultID: "pas_[...]",
//		Format:   audit.DFcsv,
//	})
func (a *audit) DownloadResults(ctx context.Context, input *DownloadRequest) (*pangea.PangeaResponse[DownloadResult], error) {
	return request.DoPost(ctx, a.Client, "v1/download_results", input, &DownloadResult{})
}

// @summary Log streaming endpoint
//
// @description This API allows 3rd party vendors (like Auth0) to stream events
// to this endpoint where the structure of the payload varies across different
// vendors.
//
// @operationId audit_post_v1_log_stream
//
// @example
//
//	type LogStreamEventData struct {
//		ClientID     string  `json:"client_id"`
//		Connection   *string `json:"connection,omitempty"`
//		ConnectionID *string `json:"connection_id,omitempty"`
//		Date         string  `json:"date"`
//		Description  string  `json:"description"`
//		IP           string  `json:"ip"`
//		Strategy     *string `json:"strategy,omitempty"`
//		StrategyType *string `json:"strategy_type,omitempty"`
//		Type         string  `json:"type"`
//		UserAgent    string  `json:"user_agent"`
//		UserID       string  `json:"user_id"`
//	}
//
//	type LogStreamEvent struct {
//		LogID string             `json:"log_id"`
//		Data  LogStreamEventData `json:"data"`
//	}
//
//	type LogStreamInput struct {
//		pangea.BaseRequest
//
//		Logs []LogStreamEvent `json:"logs"`
//	}
//
//	logStreamEvent := LogStreamEvent{
//		LogID: "some log ID",
//		Data: LogStreamEventData{
//			ClientID:    "test client ID",
//			Date:        "2024-03-29T17:26:50.193Z",
//			Description: "Create a log stream",
//			IP:          "127.0.0.1",
//			Type:        "some_type",
//			UserAgent:   "AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0",
//			UserID:      "test user ID",
//		},
//	}
//
//	input := LogStreamInput{
//		Logs: []LogStreamEvent{logStreamEvent},
//	}
//
//	response, err := client.LogStream(ctx, &input)
func (a *audit) LogStream(ctx context.Context, input pangea.ConfigIDer) (*pangea.PangeaResponse[struct{}], error) {
	var result struct{}
	return request.DoPost(ctx, a.Client, "v1/log_stream", input, &result)
}

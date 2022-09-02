package audit

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/pangeacyber/go-pangea/internal/pangeautil"
	"github.com/pangeacyber/go-pangea/internal/signer"
	"github.com/pangeacyber/go-pangea/pangea"
	"github.com/pangeacyber/go-pangea/pangea/hash"
)

// Log an entry
//
// Create a log entry in the Secure Audit Log.
//
// Example:
//
//	input := &audit.LogInput{
//		Event: &audit.LogEventInput{
//			Message: pangea.String("some important message."),
//		},
//		ReturnHash: pangea.Bool(true),
//	}
//
//	logResponse, err := auditcli.Log(ctx, input)
func (a *Audit) Log(ctx context.Context, input *LogInput) (*pangea.PangeaResponse[LogOutput], error) {
	if a.SignLogs {
		err := input.Sign(a.Signer)
		if err != nil {
			return nil, err
		}
	}
	req, err := a.Client.NewRequest("POST", "v1/log", input)
	if err != nil {
		return nil, err
	}

	var out LogOutput
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[LogOutput]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

// Search for events
//
// Search for events that match the provided search criteria.
//
// Example:
//
//	input := &audit.SearchInput{
//		Query:                  pangea.String("message:log-123"),
//		IncludeMembershipProof: pangea.Bool(true),
//	}
//
//	searchResponse, err := auditcli.Search(ctx, input)
func (a *Audit) Search(ctx context.Context, input *SearchInput) (*pangea.PangeaResponse[SearchOutput], error) {
	if a.VerifyProofs && (!pangea.BoolValue(input.IncludeHash) || !pangea.BoolValue(input.IncludeMembershipProof) || !pangea.BoolValue(input.IncludeRoot)) {
		return nil, fmt.Errorf("audit: should include hash, membership_proof and root if VerifyProofs is true")
	}
	req, err := a.Client.NewRequest("POST", "v1/search", input)
	if err != nil {
		return nil, err
	}
	out := SearchOutput{}
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	err = a.verifyRecords(out.Events, out.Root)
	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[SearchOutput]{
		Response: *resp,
		Result:   &out,
	}
	return &panresp, nil
}

// SearchResults is used to page through results from a previous search.
func (a *Audit) SearchResults(ctx context.Context, input *SearchResultInput) (*pangea.PangeaResponse[SearchResultOutput], error) {
	req, err := a.Client.NewRequest("POST", "v1/results", input)
	if err != nil {
		return nil, err
	}
	out := SearchResultOutput{}
	resp, err := a.Client.Do(ctx, req, &out)
	if err != nil {
		return nil, err
	}
	err = a.verifyRecords(out.Events, out.Root)
	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[SearchResultOutput]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

// Retrieve tamperproof verification
//
// Root returns current root hash and consistency proof.
//
// Example:
//
//	input := &audit.RootInput{
//		TreeSize: pangea.Int(10),
//	}
//
//	rootResponse, err := auditcli.Root(ctx, input)
func (a *Audit) Root(ctx context.Context, input *RootInput) (*pangea.PangeaResponse[RootOutput], error) {
	req, err := a.Client.NewRequest("POST", "v1/root", input)
	if err != nil {
		return nil, err
	}
	var out RootOutput
	resp, err := a.Client.Do(ctx, req, &out)

	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[RootOutput]{
		Response: *resp,
		Result:   &out,
	}
	return &panresp, nil
}

// SearchAll is a helper function to return all the search results for a search with pages
func SearchAll(ctx context.Context, client Client, input *SearchInput) (*Root, SearchEvents, error) {
	resp, err := client.Search(ctx, input)
	if err != nil {
		return nil, nil, err
	}
	events := make(SearchEvents, 0, *resp.Result.Count)
	events = append(events, resp.Result.Events...)
	for pangea.IntValue(resp.Result.Count) > len(events) {
		s := SearchResultInput{
			ID:                     resp.Result.ID,
			IncludeMembershipProof: input.IncludeMembershipProof,
			IncludeHash:            input.IncludeHash,
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

// SearchAll is a helper function to return all the search results for a search with pages and valitade membership proof and consitency proof
func SearchAllAndValidate(ctx context.Context, client Client, input *SearchInput, r RootsProvider) (*Root, ValidateEvents, error) {
	root, events, err := SearchAll(ctx, client, input)
	if err != nil {
		return nil, nil, err
	}
	vEvents, err := VerifyAuditRecords(ctx, r, root, events, true)
	if err != nil {
		return nil, nil, err
	}
	return root, vEvents, nil
}

func (a *Audit) verifyRecords(events SearchEvents, root *Root) error {
	for idx, event := range events {
		if a.VerifyProofs && !event.VerifyHash() {
			return fmt.Errorf("audit: cannot verify hash of record [%v]", idx)
		}

		if a.VerifyProofs && !event.VerifyMembershipProof(root) {
			return fmt.Errorf("audit: cannot verify membership proof of record [%v]", idx)
		}

		if a.VerifySignature && !event.EventEnvelope.VerifySignature() {
			return fmt.Errorf("audit: cannot verify signature of record [%v]", idx)
		}
	}
	return nil
}

type LogInput struct {
	// A structured event describing an auditable activity.
	Event *Event `json:"event"`

	// Return the event's hash with response.
	ReturnHash *bool `json:"return_hash"`

	// If true, be verbose in the response; include canonical events, create time, hashes, etc.
	// default: false
	Verbose *bool `json:"verbose"`

	// An optional client-side signature for forgery protection.
	// max len of 256 bytes
	Signature *string `json:"signature,omitempty"`

	// The base64-encoded ed25519 public key used for the signature, if one is provided
	PublicKey *string `json:"public_key,omitempty"`
}

func (i *LogInput) Sign(s signer.Signer) error {
	b, err := newsSignedMessageFromRecord(i.Event.Actor, i.Event.Action, i.Event.Message, i.Event.New,
		i.Event.Old, i.Event.Source, i.Event.Status, i.Event.Target, i.Event.Timestamp)
	if err != nil {
		return err
	}
	signature, err := s.Sign(b)
	if err != nil {
		return err
	}
	i.Signature = pangea.String(base64.StdEncoding.EncodeToString(signature))
	i.PublicKey = pangea.String(s.PublicKey())
	return nil
}

type Event struct {
	// Record who performed the auditable activity.
	// max len is 128 bytes
	// examples:
	// 	John Doe
	//  user-id
	//  DennisNedry@InGen.com
	Actor *string `json:"actor,omitempty"`

	// The auditable action that occurred."
	// examples:
	// 	created
	//  deleted
	//  updated
	Action *string `json:"action,omitempty"`

	// A message describing a detailed account of what happened.
	// This can be recorded as free-form text or as a JSON-formatted string.
	// Message is a required field.
	// max len of 65536 bytes
	Message *string `json:"message"`

	// The value of a record after it was changed.
	// max len of 65536 bytes
	New *string `json:"new,omitempty"`

	// The value of a record before it was changed.
	// max len of 65536 bytes
	Old *string `json:"old,omitempty"`

	// Used to record the location from where an activity occurred.
	// max len of 128 bytes
	Source *string `json:"source,omitempty"`

	// Record whether or not the activity was successful.
	// examples:
	//  failure
	//  success
	// max len of 32 bytes
	Status *string `json:"status,omitempty"`

	// Used to record the specific record that was targeted by the auditable activity.
	// max len of 128 bytes
	Target *string `json:"target,omitempty"`

	// An optional client-supplied timestamp.
	Timestamp *string `json:"timestamp,omitempty"`

	// Timestamp set by the server
	ReceivedAt *string `json:"received_at,omitempty"`
}

type EventEnvelope struct {
	// A structured record describing that <actor> did <action> on <target>
	// changing it from <old> to <new> and the operation was <status>,
	// and/or a free-form <message>.
	Event *Event `json:"event"`

	// An optional client-side signature for forgery protection.
	// max len of 256 bytes
	Signature *string `json:"signature,omitempty"`

	// The base64-encoded ed25519 public key used for the signature, if one is provided
	PublicKey *string `json:"public_key,omitempty"`

	// A server-supplied timestamp.
	ReceivedAt *string `json:"received_at,omitempty"`
}

type LogOutput struct {
	EventEnvelope *EventEnvelope `json:"envelope"`

	// The hash of the event data.
	// max len of 64 bytes
	Hash *string `json:"hash"`

	// A base64 encoded canonical JSON form of the event, used for hashing.
	CanonicalEventBase64 *string `json:"canonical_event_base64"`
}

type SearchInput struct {
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
	Query *string `json:"query"`

	// Specify the sort order of the response. "asc" or "desc"
	Order *string `json:"order,omitempty"`

	// Name of column to sort the results by.
	OrderBy *string `json:"order_by,omitempty"`

	// If set, the last value from the response to fetch the next page from.
	Last *string `json:"last,omitempty"`

	// The start of the time range to perform the search on.
	Start *time.Time `json:"start,omitempty"`

	// The end of the time range to perform the search on. All records up to the latest if left out.
	End *time.Time `json:"end,omitempty"`

	// Number of audit records to include from the first page of the results.
	Limit *int `json:"limit,omitempty"`

	// Maximum number of results to return.
	// min 1 max 10000
	MaxResults *int `json:"max_results,omitempty"`

	// If true, include membership proofs for each record in the first page.
	IncludeMembershipProof *bool `json:"include_membership_proof,omitempty"`

	// If true, include hashes for each record in the first page.
	IncludeHash *bool `json:"include_hash,omitempty"`

	// If true, include the Merkle root hash of the tree in the first page.
	IncludeRoot *bool `json:"include_root,omitempty"`

	// A list of keys to restrict the search results to. Useful for partitioning data available to the query string.
	SearchRestriction *SearchRestriction `json:"search_restriction,omitempty"`
}

type SearchRestriction struct {
	// A list of actors to restrict the search to.
	Actor []*string `json:"actor,omitempty"`

	// A list of sources to restrict the search to.
	Source []*string `json:"source,omitempty"`

	// A list of targets to restrict the search to.
	Target []*string `json:"target,omitempty"`
}

type SearchOutput struct {
	// Identifier to supply to search_results API to fetch/paginate through search results.
	// ID is always populated on a successful response.
	ID *string `json:"id"`

	// The time when the results will no longer be available to page through via the results API.
	// ExpiresAt is always populated on a successful response.
	ExpiresAt *time.Time `json:"expires_at"`

	// The total number of results that were returned by the search.
	// Count is always populated on a successful response.
	Count *int `json:"count"`

	// A root of a Merkle Tree
	Root *Root `json:"root"`

	// A list of matching audit records.
	// Events is always populated on a successful response.
	Events SearchEvents `json:"events"`
}

type SearchEvents []*SearchEvent

// VerifiableRecords retuns a slice of records that can be verifiable by the published proof
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
	EventEnvelope EventEnvelope `json:"envelope"`

	// The record's hash
	// len of 64 bytes
	Hash *string `json:"hash"`

	// The index of the leaf of the Merkle Tree where this record was inserted.
	LeafIndex *int `json:"leaf_index"`

	// A cryptographic proof that the record has been persisted in the log.
	MembershipProof *string `json:"membership_proof"`
}

// IsVerifiable checks if a record can be verfiable with the published proof
func (event *SearchEvent) IsVerifiable() bool {
	return event.LeafIndex != nil
}

func (ee *SearchEvent) VerifyHash() bool {
	if ee.Hash == nil {
		// FIXME: Why?
		return true
	}

	eventCanon, err := pangeautil.CanonicalizeJSONMarshall((ee.EventEnvelope))
	if err != nil {
		return false
	}
	eventHash := hash.Encode(eventCanon)
	if err != nil {
		return false
	}

	return pangea.StringValue(ee.Hash) == eventHash.String()
}

func (ee *SearchEvent) VerifyMembershipProof(root *Root) bool {
	if root == nil || ee.MembershipProof == nil {
		return true
	}

	b, err := VerifyMembershipProof(*root, *ee, false)
	if err != nil {
		return false
	}
	return b
}

func (ee *EventEnvelope) VerifySignature() bool {
	if ee.Signature == nil {
		return true
	}

	b, err := newsSignedMessageFromRecord(ee.Event.Actor, ee.Event.Action, ee.Event.Message, ee.Event.New,
		ee.Event.Old, ee.Event.Source, ee.Event.Status, ee.Event.Target, ee.Event.Timestamp)
	if err != nil {
		return false
	}

	sig, err := base64.StdEncoding.DecodeString(pangea.StringValue(ee.Signature))
	if err != nil {
		return false
	}

	pubKey, err := base64.StdEncoding.DecodeString(pangea.StringValue(ee.PublicKey))
	if err != nil {
		return false
	}

	v := signer.NewVerifierFromPubKey(pubKey)
	return v.Verify(b, sig)
}

type SearchResultInput struct {
	// A search results identifier returned by the search call
	// ID is a required field
	ID *string `json:"id"`

	// If true, include membership proofs for each record in the first page.
	IncludeMembershipProof *bool `json:"include_membership_proof,omitempty"`

	// If true, include hashes for each record in the first page.
	IncludeHash *bool `json:"include_hash,omitempty"`

	// If true, include the Merkle root hash of the tree in the first page.
	IncludeRoot *bool `json:"include_root,omitempty"`

	// Number of audit records to include from the first page of the results.
	Limit *int `json:"limit,omitempty"`

	// Offset from the start of the result set to start returning results from.
	Offset *int `json:"offset,omitempty"`
}

type SearchResultOutput struct {
	// The total number of results that were returned by the search.
	// Count is always populated on a successful response.
	Count *int `json:"count"`

	// A list of matching audit records.
	// Events is always populated on a successful response.
	Events SearchEvents `json:"events"`

	// A root of a Merkle Tree
	Root *Root `json:"root"`
}

type RootInput struct {
	// The size of the tree (the number of records)
	TreeSize *int `json:"tree_size,omitempty"`
}

type Root struct {
	// The name of the Merkle Tree
	TreeName *string `json:"tree_name"`

	// The size of the tree (the number of records)
	Size *int `json:"size"`

	// The root hash
	// max len of 64 bytes
	RootHash *string `json:"root_hash"`

	// The URL where this root has been published
	URL *string `json:"url"`

	// The date/time when this root was published
	PublishedAt *time.Time `json:"published_at"`

	// Consistency proof to verify that this root is a continuation of the previous one
	ConsistencyProof []*string `json:"consistency_proof"`
}

type RootOutput struct {
	Data *Root `json:"data"`
}

type signedMessage struct {
	Actor     *string `json:"actor"`
	Action    *string `json:"action"`
	Message   *string `json:"message"`
	New       *string `json:"new"`
	Old       *string `json:"old"`
	Source    *string `json:"source"`
	Status    *string `json:"status"`
	Target    *string `json:"target"`
	Timestamp *string `json:"timestamp"`
}

func newsSignedMessageFromRecord(actor, action, message, new, old, source, status, target, timestamp *string) ([]byte, error) {
	return pangeautil.CanonicalizeJSONMarshall(
		signedMessage{
			Actor:     actor,
			Action:    action,
			Message:   message,
			New:       new,
			Old:       old,
			Source:    source,
			Status:    status,
			Target:    target,
			Timestamp: timestamp,
		},
	)
}

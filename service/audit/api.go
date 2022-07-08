package audit

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/pangeacyber/go-pangea/internal/pangeautil"
	"github.com/pangeacyber/go-pangea/internal/signer"
	"github.com/pangeacyber/go-pangea/pangea"
)

// Log creates a log entry in the Secure Audit Log.
func (a *Audit) Log(ctx context.Context, input *LogInput) (*LogOutput, *pangea.Response, error) {
	if a.SignLogs {
		err := input.Event.Sign(a.Signer)
		if err != nil {
			return nil, nil, err
		}
	}
	req, err := a.Client.NewRequest("POST", "v1/log", input)
	if err != nil {
		return nil, nil, err
	}

	var out LogOutput
	resp, err := a.Client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}

	return &out, resp, nil
}

// Search the Secure Audit Log
func (a *Audit) Search(ctx context.Context, input *SearchInput) (*SearchOutput, *pangea.Response, error) {
	req, err := a.Client.NewRequest("POST", "v1/search", input)
	if err != nil {
		return nil, nil, err
	}
	out := SearchOutput{}
	resp, err := a.Client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}
	err = a.verifyRecords(out.Events)
	if err != nil {
		return nil, nil, err
	}
	return &out, resp, nil
}

// SearchResults is used to page through results from a previous search.
func (a *Audit) SearchResults(ctx context.Context, input *SeachResultInput) (*SeachResultOutput, *pangea.Response, error) {
	req, err := a.Client.NewRequest("POST", "v1/results", input)
	if err != nil {
		return nil, nil, err
	}
	out := SeachResultOutput{}
	resp, err := a.Client.Do(ctx, req, &out)
	if err != nil {
		return nil, nil, err
	}
	err = a.verifyRecords(out.Events)
	if err != nil {
		return nil, nil, err
	}
	return &out, resp, nil
}

// Root returns current root hash and consistency proof.
func (a *Audit) Root(ctx context.Context, input *RootInput) (*RootOutput, *pangea.Response, error) {
	req, err := a.Client.NewRequest("POST", "v1/root", input)
	if err != nil {
		return nil, nil, err
	}
	var out RootOutput
	resp, err := a.Client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

func (a *Audit) verifyRecords(events Events) error {
	if a.VerifyRecords {
		for idx, event := range events {
			verified := event.Record.VerifySignature(a.Verifier)
			if !verified {
				return fmt.Errorf("audit: cannot verify signature of record [%v]", idx)
			}
		}
	}
	return nil
}

type LogInput struct {
	// A structured event describing an auditable activity.
	Event *LogEventInput `json:"event"`

	// Return the event's hash with response.
	ReturnHash *bool `json:"return_hash"`

	// If true, be verbose in the response; include canonical events, create time, hashes, etc.
	// default: false
	Verbose *bool `json:"verbose"`
}

type LogEventInput struct {
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

	// An optional client-side signature for forgery protection.
	// max len of 256 bytes
	Signature *string `json:"signature,omitempty"`

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
}

func (i *LogEventInput) Sign(s signer.Signer) error {
	b, err := newsSignedMessageFromRecord(i.Actor, i.Action, i.Message, i.New,
		i.Old, i.Source, i.Status, i.Target, i.Timestamp)
	if err != nil {
		return err
	}
	signature, err := s.Sign(b)
	if err != nil {
		return err
	}
	i.Signature = pangea.String(base64.StdEncoding.EncodeToString(signature))
	return nil
}

type LogOutput struct {
	// The hash of the event data.
	// max len of 64 bytes
	Hash *string `json:"hash"`

	// The event that was logged. Includes additional server-added timestamps.
	Event *LogEventOutput `json:"event"`

	// A base64 encoded canonical JSON form of the event, used for hashing.
	CanonicalEventBase64 *string `json:"canonical_event_base64"`
}

type LogEventOutput struct {
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

	// An optional client-side signature for forgery protection.
	// max len of 256 bytes
	Signature *string `json:"signature,omitempty"`

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
	Start *string `json:"start,omitempty"`

	// The end of the time range to perform the search on. All records up to the latest if left out.
	End *string `json:"end,omitempty"`

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
	Events Events `json:"events"`
}

type Events []*EventEnvelope

// VerifiableRecords retuns a slice of records that can be verifiable by the published proof
func (events Events) VerifiableRecords() Events {
	evs := make(Events, 0)
	for _, event := range events {
		if event.IsVerifiable() {
			evs = append(evs, event)
		}
	}
	return evs
}

type EventEnvelope struct {
	// A structured record describing that <actor> did <action> on <target>
	// changing it from <old> to <new> and the operation was <status>,
	// and/or a free-form <message>.
	Record *Record `json:"event"`

	// A cryptographic proof that the record has been persisted in the log.
	MembershipProof *string `json:"membership_proof"`

	// The record's hash
	// len of 64 bytes
	Hash *string `json:"hash"`

	// The index of the leaf of the Merkle Tree where this record was inserted.
	LeafIndex *int `json:"leaf_index"`
}

// IsVerifiable checks if a record can be verfiable with the published proof
func (event *EventEnvelope) IsVerifiable() bool {
	return event.LeafIndex != nil
}

type Record struct {
	// An identifier for _who_ the audit record is about.
	Actor *string `json:"actor,omitempty"`

	// What action was performed on a record.
	// examples:
	// 	created
	//  deleted
	//  updated
	Action *string `json:"action,omitempty"`

	// A free form text field describing the event.
	// Message is always populated on a successful response.
	Message *string `json:"message"`

	// The value of a record _after_ it was changed.
	New *string `json:"new,omitempty"`

	// The value of a record _before_ it was changed.
	Old *string `json:"old,omitempty"`

	// An optional client-side signature for forgery protection.
	// max len of 256 bytes
	Signature *string `json:"signature,omitempty"`

	// The source of a record. Can be used to hard-split logged and searched data.
	// max len of 128 bytes
	Source *string `json:"source,omitempty"`

	// The status or result of the event
	// examples:
	//  failure
	//  success
	// max len of 32 bytes
	Status *string `json:"status,omitempty"`

	// An identifier for what the audit record is about.
	// max len of 128 bytes
	Target *string `json:"target,omitempty"`

	Timestamp *string `json:"timestamp"`
}

func (r *Record) VerifySignature(verifier signer.Verifier) bool {
	b, err := newsSignedMessageFromRecord(r.Actor, r.Action, r.Message, r.New,
		r.Old, r.Source, r.Status, r.Target, r.Timestamp)
	if err != nil {
		return false
	}
	sig, err := base64.StdEncoding.DecodeString(pangea.StringValue(r.Signature))
	if err != nil {
		return false
	}
	return verifier.Verify(b, sig)
}

type SeachResultInput struct {
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

type SeachResultOutput struct {
	// The total number of results that were returned by the search.
	// Count is always populated on a successful response.
	Count *int `json:"count"`

	// A list of matching audit records.
	// Events is always populated on a successful response.
	Events Events `json:"events"`

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

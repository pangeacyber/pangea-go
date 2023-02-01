package audit

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	pu "github.com/pangeacyber/pangea-go/pangea-sdk/internal/pangeautil"
	"github.com/pangeacyber/pangea-go/pangea-sdk/internal/signer"
	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
)

// @summary Log an entry
//
// @description Create a log entry in the Secure Audit Log.
//
// @example
//
//	event := audit.Event{
//		Message: "Integration test msg",
//	 }
//
//	logResponse, err := auditcli.Log(ctx, event, true)
func (a *Audit) Log(ctx context.Context, event Event, verbose bool) (*pangea.PangeaResponse[LogOutput], error) {
	input := LogInput{
		Event:   event,
		Verbose: verbose,
	}

	if a.VerifyProofs {
		input.Verbose = true
		input.PrevRoot = a.lastUnpRootHash
	}

	if a.SignLogsMode == LocalSign && a.Signer != nil {
		err := input.SignEvent(*a.Signer, a.publicKeyInfo)
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

	out.EventEnvelope, err = newEventEnvelopeFromMap(out.RawEnvelope)
	if err != nil {
		return nil, err
	}
	err = a.processLogResponse(ctx, &out)
	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[LogOutput]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

// @summary Search for events
//
// @description Search for events that match the provided search criteria.
//
// @example
//
//	input := &audit.SearchInput{
//		Query:                  pangea.String("message:log-123"),
//		IncludeMembershipProof: pangea.Bool(true),
//	}
//
//	searchResponse, err := auditcli.Search(ctx, input)
func (a *Audit) Search(ctx context.Context, input *SearchInput) (*pangea.PangeaResponse[SearchOutput], error) {
	if input == nil {
		return nil, errors.New("nil input")
	}

	if a.VerifyProofs {
		// Need this info to verify
		input.Verbose = pangea.Bool(true)
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

	err = a.processSearchEvents(ctx, out.Events, out.Root, out.UnpublishedRoot)
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
	err = a.processSearchEvents(ctx, out.Events, out.Root, out.UnpublishedRoot)
	if err != nil {
		return nil, err
	}

	panresp := pangea.PangeaResponse[SearchResultOutput]{
		Response: *resp,
		Result:   &out,
	}

	return &panresp, nil
}

// @summary Retrieve tamperproof verification
//
// @description Root returns current root hash and consistency proof.
//
// @example
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
	events := make(SearchEvents, 0, resp.Result.Count)
	events = append(events, resp.Result.Events...)
	for resp.Result.Count > len(events) {
		s := SearchResultInput{
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

func (a *Audit) processLogResponse(ctx context.Context, log *LogOutput) error {
	if log == nil {
		return nil
	}

	nurh := log.UnpublishedRootHash
	if VerifyHash(log.RawEnvelope, log.Hash) == Failed {
		return fmt.Errorf("audit: Failed hash verification of event. Hash: [%s]", log.Hash)
	}
	if log.EventEnvelope != nil {
		log.SignatureVerification = log.EventEnvelope.VerifySignature()
	}

	if a.VerifyProofs {
		if VerifyHash(log.RawEnvelope, log.Hash) == Failed {
			return fmt.Errorf("audit: cannot verify hash of event. Hash: [%s]", log.Hash)
		}

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

func (a *Audit) processSearchEvents(ctx context.Context, events SearchEvents, root *Root, unpRoot *Root) error {
	var roots map[int]Root

	var err error
	for _, event := range events {
		event.EventEnvelope, err = newEventEnvelopeFromMap(event.RawEnvelope)
		if err != nil {
			return err
		}
	}

	if a.VerifyProofs && root != nil {
		if a.rp == nil {
			a.rp = NewArweaveRootsProvider(root.TreeName)
		}
		treeSizes := treeSizes(root, events)
		roots = a.rp.UpdateRoots(ctx, treeSizes)
	}

	for _, event := range events {
		if !a.SkipEventVerification {
			if VerifyHash(event.RawEnvelope, event.Hash) == Failed {
				return fmt.Errorf("audit: cannot verify hash of record. Hash: [%s]", event.Hash)
			}
			event.SignatureVerification = event.EventEnvelope.VerifySignature()
		}

		if a.VerifyProofs {
			if event.Published != nil && *event.Published == true {
				event.VerifyMembershipProof(root)
				event.VerifyConsistencyProof(roots)
			} else {
				event.VerifyMembershipProof(unpRoot)
			}
		}
	}
	return nil
}

type LogInput struct {
	// A structured event describing an auditable activity.
	Event Event `json:"event"`

	// If true, be verbose in the response; include root, membership and consistency proof, etc.
	// default: false
	Verbose bool `json:"verbose"`

	// An optional client-side signature for forgery protection.
	// max len of 256 bytes
	Signature *string `json:"signature,omitempty"`

	// The base64-encoded ed25519 public key used for the signature, if one is provided
	PublicKey *string `json:"public_key,omitempty"`

	// Previous unpublished root
	PrevRoot *string `json:"prev_root,omitempty"`
}

func (i *LogInput) SignEvent(s signer.Signer, pki map[string]string) error {
	b, err := pu.CanonicalizeStruct(&i.Event)
	if err != nil {
		return err
	}

	signature, err := s.Sign(b)
	if err != nil {
		return err
	}

	pk, err := s.PublicKey()
	if err != nil {
		return err
	}

	if pki != nil {
		pki["key"] = pk
	} else {
		pki = map[string]string{
			"key": pk,
		}
	}
	pkib, err := pu.CanonicalizeStruct(pki)
	if err != nil {
		return err
	}

	i.Signature = pangea.String(base64.StdEncoding.EncodeToString(signature))
	i.PublicKey = pangea.String(string(pkib))
	return nil
}

type Event struct {
	// Record who performed the auditable activity.
	// max len is 128 bytes
	// examples:
	// 	John Doe
	//  user-id
	//  DennisNedry@InGen.com
	Actor string `json:"actor,omitempty"`

	// The auditable action that occurred."
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
	ReceivedAt *pu.PangeaTimestamp `json:"received_at,omitempty"`
}

func newEventEnvelopeFromMap(m *map[string]any) (*EventEnvelope, error) {
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

	return &ee, nil
}

type EventVerification int

const (
	NotVerified EventVerification = iota
	Success
	Failed
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

type LogOutput struct {
	EventEnvelope *EventEnvelope

	RawEnvelope *map[string]any `json:"envelope"`

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
	EventEnvelope *EventEnvelope

	RawEnvelope *map[string]any `json:"envelope"`

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

// IsVerifiable checks if a record can be verfiable with the published proof
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
	if ee.Published == nil || *ee.Published != true || ee.LeafIndex == nil {
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
	return
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
	if v == nil {
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
		return "", errors.New("Public key field nil pointer")
	}

	pkinfo := make(map[string]any)
	err := json.Unmarshal([]byte(*ee.PublicKey), &pkinfo)
	if err == nil {
		val, ok := pkinfo["key"]
		if ok {
			ret := val.(string)
			return ret, nil
		} else {
			return "", errors.New("'key' field not present in json")
		}
	} else {
		return *ee.PublicKey, nil
	}
}

type SearchResultInput struct {
	// A search results identifier returned by the search call
	// ID is a required field
	ID string `json:"id"`

	// Number of audit records to include from the first page of the results.
	Limit int `json:"limit,omitempty"`

	// Offset from the start of the result set to start returning results from.
	Offset *int `json:"offset,omitempty"`
}

type SearchResultOutput struct {
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

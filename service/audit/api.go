package audit

import (
	"context"

	"github.com/pangeacyber/go-pangea/internal/pangeautil"
	"github.com/pangeacyber/go-pangea/pangea"
)

type LogInput struct {
	// A structured record describing that <actor> did <action> on <target>
	// changing it from <old> to <new> and the operation was <status>, and/or a free-form <message>.
	Data *Record `json:"data"`

	// If the response should include the hash.
	ReturnHash *bool `json:"return_hash"`
}

func (l LogInput) String() string {
	return pangeautil.Stringify(l)
}

type Record struct {
	// An identifier for _who_ the audit record is about.
	Actor *string `json:"actor,omitempty"`

	// What action was performed on a record.
	// eg: "created", "deleted", "updated"
	Action *string `json:"action,omitempty"`

	// A free form text field describing the event.
	// Message is a required field.
	Message *string `json:"message"`

	// The value of a record _after_ it was changed.
	New *string `json:"new,omitempty"`

	// The value of a record _before_ it was changed.
	Old *string `json:"old,omitempty"`

	// The source of a record. Can be used to hard-split logged and searched data.
	Source *string `json:"source,omitempty"`

	// The status or result of the event
	// eg: "failure", "success"
	Status *string `json:"status,omitempty"`

	// An identifier for what the audit record is about.
	Target *string `json:"target,omitempty"`
}

func (r Record) String() string {
	return pangeautil.Stringify(r)
}

type LogOutput struct {
	// The hash of the log.
	Hash *string `json:"hash"`
}

func (l LogOutput) String() string {
	return pangeautil.Stringify(l)
}

func (a *Audit) Log(ctx context.Context, input *LogInput) (*LogOutput, *pangea.Response, error) {
	if input == nil {
		input = &LogInput{}
	}
	req, err := a.Client.NewRequest("POST", "audit", "v1/audit/log", input)
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

type SerarchInput struct {
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
	// Some examples:
	//		actor:root target:/etc/shadow
	//		viewed x-ray
	Query *string `json:"query"`

	// Specify the sort order of the response. "asc" or "desc"
	Order *string `json:"order,omitempty"`

	// If set, the last value from the response to fetch the next page from.
	Last *string `json:"last,omitempty"`

	// The start of the time range to perform the search on.
	Start *string `json:"start,omitempty"`

	// The end of the time range to perform the search on. All records up to the latest if left out.
	End *string `json:"end,omitempty"`

	// A list of sources that the search can apply to. If empty or not provided, matches only the default source.
	Sources []*string `json:"sources,omitempty"`

	// If the response should include the memebership proof.
	IncludeMembershipProof *bool `json:"include_membership_proof,omitempty"`
}

func (s SerarchInput) String() string {
	return pangeautil.Stringify(s)
}

type SearchOutput struct {
	Root *Root `json:"root"`

	// A list of matching audit records.
	Audits AuditRecords `json:"audits"`

	// An opaque identifier that can be used to fetch the next page of the results.
	Last string `json:"last"`
}

func (s SearchOutput) String() string {
	return pangeautil.Stringify(s)
}

type AuditRecord struct {
	// A structured record describing that <actor> did <action> on <target>
	// changing it from <old> to <new> and the operation was <status>, and/or a free-form <message>.
	Data *Record `json:"data"`

	// A list of hashes that prove the membership of the log with the root hash.
	MembershipProof *string `json:"membership_proof"`

	// The hash of the log.
	Hash *string `json:"hash"`

	// The index of the leaf in the log.
	LeafIndex *int `json:"leaf_index"`
}

func (a AuditRecord) String() string {
	return pangeautil.Stringify(a)
}

// IsVerifiable checks if a record can be verfiable with the published proof
func (record *AuditRecord) IsVerifiable() bool {
	return record.LeafIndex != nil
}

type AuditRecords []*AuditRecord

func (a AuditRecords) String() string {
	return pangeautil.Stringify(a)
}

// VerifiableRecords retuns a slice of records that can be verifiable by the published proof
func (records AuditRecords) VerifiableRecords() AuditRecords {
	r := make(AuditRecords, 0)
	for _, record := range records {
		if record.IsVerifiable() {
			r = append(r, record)
		}
	}
	return r
}

func (a *Audit) Search(ctx context.Context, input *SerarchInput) (*SearchOutput, *pangea.Response, error) {
	if input == nil {
		input = &SerarchInput{}
	}
	req, err := a.Client.NewRequest("POST", "audit", "v1/search", input)
	if err != nil {
		return nil, nil, err
	}
	out := SearchOutput{}
	resp, err := a.Client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

type Root struct {
	TreeName         *string   `json:"tree_name"`
	Size             *int      `json:"size"`
	RootHash         *string   `json:"root_hash"`
	URL              *string   `json:"url"`
	PublishedAt      *string   `json:"published_at"`
	ConsistencyProof []*string `json:"consistency_proof"`
}

func (r Root) String() string {
	return pangeautil.Stringify(r)
}

type RootInput struct {
	TreeSize *int `json:"tree_size"`
}

func (r RootInput) String() string {
	return pangeautil.Stringify(r)
}

type RootOutput struct {
	Data *Root `json:"data"`
}

func (r RootOutput) String() string {
	return pangeautil.Stringify(r)
}

func (a *Audit) Root(ctx context.Context, input *RootInput) (*RootOutput, *pangea.Response, error) {
	if input == nil {
		input = &RootInput{}
	}
	req, err := a.Client.NewRequest("POST", "audit", "v1/root", input)
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

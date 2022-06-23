package audit

import (
	"context"

	"github.com/pangeacyber/go-pangea/pangea"
)

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

type SearchOutput struct {
	Root *Root `json:"root"`

	// A list of matching audit records.
	Audits AuditRecords `json:"audits"`

	// An opaque identifier that can be used to fetch the next page of the results.
	Last string `json:"last"`
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

// IsVerifiable checks if a record can be verfiable with the published proof
func (record *AuditRecord) IsVerifiable() bool {
	return record.LeafIndex != nil
}

type AuditRecords []*AuditRecord

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

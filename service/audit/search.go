package audit

import (
	"context"
	"fmt"
	"strings"

	"github.com/pangeacyber/go-pangea/pangea"
	"github.com/pangeacyber/go-pangea/pangea/hash"
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
	Audits []*AuditRecord `json:"audits"`

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

type ProofSide uint

const (
	Left ProofSide = iota
	Right
)

type ProofItem struct {
	Side ProofSide
	Hash hash.Hash
}

type Proof []ProofItem

func DecodeProof(s string) (Proof, error) {
	items := strings.Split(s, ",")
	if len(items) == 0 {
		return nil, nil
	}

	proof := make(Proof, 0, len(items))
	for _, item := range items {
		parts := strings.Split(item, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("audit: inconsistent proof item with no separation: %v", item)
		}
		h, err := hash.Decode(parts[1])
		if err != nil {
			return nil, fmt.Errorf("audit: invalid hash in proof item: %w", err)
		}
		proofItem := ProofItem{
			Hash: h,
		}
		switch parts[0] {
		case "l":
			proofItem.Side = Left
		case "r":
			proofItem.Side = Right
		default:
			return nil, fmt.Errorf("audit: inconsistent proof item with no side declaration: side: %v hash: %v", parts[0], parts[1])
		}
		proof = append(proof, proofItem)
	}
	return proof, nil
}

func VerifyMembershipProof(root *Root, auditOutput *AuditRecord, required bool) (bool, error) {
	membershipProof := pangea.StringValue(auditOutput.MembershipProof)
	if membershipProof == "" {
		return !required, nil
	}
	targetHash, err := hash.Decode(pangea.StringValue(auditOutput.Hash))
	if err != nil {
		return false, err
	}
	rootHash, err := hash.Decode(pangea.StringValue(root.RootHash))
	if err != nil {
		return false, err
	}
	proof, err := DecodeProof(membershipProof)
	if err != nil {
		return false, err
	}
	return verifyLogProof(targetHash, rootHash, proof), nil
}

func verifyLogProof(target, root hash.Hash, proof Proof) bool {
	h := target
	for _, proofItem := range proof {
		switch proofItem.Side {
		case Left:
			h = hash.Pair(proofItem.Hash).With(h)
		case Right:
			h = hash.Pair(h).With(proofItem.Hash)
		}
	}
	return root.Equal(h)
}

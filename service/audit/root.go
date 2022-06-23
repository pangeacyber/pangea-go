package audit

import (
	"context"

	"github.com/pangeacyber/go-pangea/pangea"
	"github.com/pangeacyber/go-pangea/pangea/hash"
)

type Root struct {
	TreeName         *string   `json:"tree_name"`
	Size             *int      `json:"size"`
	RootHash         *string   `json:"root_hash"`
	URL              *string   `json:"url"`
	PublishedAt      *string   `json:"published_at"`
	ConsistencyProof []*string `json:"consistency_proof"`
}

type RootInput struct {
	TreeSize *uint `json:"tree_size"`
}

type RootOutput struct {
	Data *Root `json:"data"`
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

type rootProofItem struct {
	Hash  hash.Hash
	Proof proof
}

type rootProof []rootProofItem

func decodeRootProof(rawProof []*string) (rootProof, error) {
	if len(rawProof) == 0 {
		return nil, nil
	}
	proof := make(rootProof, 0, len(rawProof))
	for _, rawItem := range rawProof {
		root, membershipProof, err := splitConsistencyProof(pangea.StringValue(rawItem))
		if err != nil {
			return nil, err
		}
		_, decodedRoot, err := hashFromHashProof(root)
		if err != nil {
			return nil, err
		}
		decodedMembershipProof, err := DecodeProof(membershipProof)
		if err != nil {
			return nil, err
		}
		proof = append(proof, rootProofItem{
			Hash:  decodedRoot,
			Proof: decodedMembershipProof,
		})
	}
	return proof, nil
}

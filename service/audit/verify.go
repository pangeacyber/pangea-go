package audit

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/pangeacyber/go-pangea/pangea"
	"github.com/pangeacyber/go-pangea/pangea/hash"
)

type proofSide uint

const (
	Left proofSide = iota
	Right
)

type proofItem struct {
	Side proofSide
	Hash hash.Hash
}

type proof []proofItem

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
		decodedMembershipProof, err := decodeProof(membershipProof)
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

func decodeProof(s string) (proof, error) {
	items := strings.Split(s, ",")
	if len(items) == 0 {
		return nil, nil
	}

	p := make(proof, 0, len(items))
	for _, item := range items {
		side, h, err := hashFromHashProof(item)
		if err != nil {
			return nil, err
		}
		pi := proofItem{
			Hash: h,
		}
		switch side {
		case "l":
			pi.Side = Left
		case "r":
			pi.Side = Right
		default:
			return nil, fmt.Errorf("audit: inconsistent proof item with no side declaration: %v", item)
		}
		p = append(p, pi)
	}
	return p, nil
}

func VerifyMembershipProof(root Root, event EventEnvelope, required bool) (bool, error) {
	membershipProof := pangea.StringValue(event.MembershipProof)
	if membershipProof == "" {
		return !required, nil
	}
	targetHash, err := hash.Decode(pangea.StringValue(event.Hash))
	if err != nil {
		return false, err
	}
	rootHash, err := hash.Decode(pangea.StringValue(root.RootHash))
	if err != nil {
		return false, err
	}
	proof, err := decodeProof(membershipProof)
	if err != nil {
		return false, err
	}
	return verifyLogProof(targetHash, rootHash, proof), nil
}

func verifyLogProof(target, root hash.Hash, p proof) bool {
	h := target
	for _, proofItem := range p {
		switch proofItem.Side {
		case Left:
			h = hash.Pair(proofItem.Hash).With(h)
		case Right:
			h = hash.Pair(h).With(proofItem.Hash)
		}
	}
	return root.Equal(h)
}

func VerifyConsistencyProof(publishedRoots map[int]Root, event EventEnvelope, required bool) bool {
	if event.LeafIndex == nil {
		return !required
	}
	idx := *event.LeafIndex
	if idx <= 1 {
		return !required
	}
	current, ok := publishedRoots[idx]
	if !ok {
		return false
	}
	previous, ok := publishedRoots[idx-1]
	if !ok {
		return false
	}
	verified, err := verifyConsistencyProof(previous, current)
	if err != nil {
		return false
	}
	return verified
}

func verifyConsistencyProof(old, new Root) (bool, error) {
	oldHash, err := hash.Decode(pangea.StringValue(old.RootHash))
	if err != nil {
		return false, err
	}
	newHash, err := hash.Decode(pangea.StringValue(new.RootHash))
	if err != nil {
		return false, err
	}
	consistencyProof, err := decodeRootProof(new.ConsistencyProof)
	if err != nil {
		return false, err
	}
	if len(consistencyProof) == 0 {
		return false, fmt.Errorf("audit: consistency proof is empty")
	}
	rootHash := consistencyProof[0].Hash
	for i := 1; i < len(consistencyProof); i++ {
		rootHash = hash.Pair(consistencyProof[i].Hash).With(rootHash)
	}
	if !rootHash.Equal(oldHash) {
		return false, nil
	}
	for _, item := range consistencyProof {
		if !verifyLogProof(item.Hash, newHash, item.Proof) {
			return false, nil
		}
	}
	return true, nil
}

func VerifyAuditRecords(ctx context.Context, rp RootsProvider, root *Root, events Events, required bool) (bool, error) {
	if root == nil || len(events) == 0 {
		return false, fmt.Errorf("audit: empty root or events")
	}

	treeSizes := treeSizes(root, events)
	roots, err := rp.Roots(ctx, treeSizes)
	if err != nil {
		return false, err
	}

	for _, event := range events {
		if !VerifyConsistencyProof(roots, *event, required) {
			return false, fmt.Errorf("audit: consistency proof failed for event %v", *event)
		}
	}
	return true, nil
}

func VerifyAuditRecordsWithArweave(ctx context.Context, root *Root, events Events, required bool) (bool, error) {
	arweavecli := NewArweaveRootsProvider(*root.TreeName)
	return VerifyAuditRecords(ctx, arweavecli, root, events, required)
}

func treeSizes(root *Root, events Events) []string {
	treeSizes := make(map[int]struct{}, 0)
	treeSizes[*root.Size] = struct{}{}
	for _, event := range events {
		leafIdx := *event.LeafIndex
		treeSizes[leafIdx] = struct{}{}
		if leafIdx > 1 {
			treeSizes[leafIdx-1] = struct{}{}
		}
	}
	sizes := make([]string, 0, len(treeSizes))
	for size := range treeSizes {
		sizes = append(sizes, strconv.Itoa(size))
	}
	return sizes
}

func hashFromHashProof(proof string) (string, hash.Hash, error) {
	label, rawHash, err := splitHashProof(proof)
	if err != nil {
		return "", nil, err
	}
	h, err := hash.Decode(rawHash)
	if err != nil {
		return "", nil, fmt.Errorf("audit: invalid hash in hash proof: %w", err)
	}
	return label, h, nil
}

// splitConsistencyProof splits a consistency proof into the root hash and the membership proof.
// the proof has the format "x:<root-hash>,l:<left-hash>,r:<right-hash>".
// the result should be
func splitConsistencyProof(proof string) (root, membershipProof string, err error) {
	root, membershipProof, ok := strings.Cut(proof, ",")
	if !ok {
		return "", "", fmt.Errorf("audit: invalid format: expected `,` in: %v", proof)
	}
	return root, membershipProof, nil
}

// splitHashProof splits a hash proof into it's label and the hash.
// the expected format is "<label>:<hash>".
func splitHashProof(proof string) (label, rawHash string, err error) {
	label, rawHash, ok := strings.Cut(proof, ":")
	if !ok {
		return "", "", fmt.Errorf("audit: invalid format: expected `:` in: %v", proof)
	}
	return label, rawHash, nil
}

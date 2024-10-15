package audit

import (
	"fmt"
	"strconv"
	"strings"

	pu "github.com/pangeacyber/pangea-go/pangea-sdk/v3/internal/pangeautil"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea/hash"
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

func VerifyHash(ee map[string]any, h string) EventVerification {
	if h == "" || ee == nil {
		return NotVerified
	}
	eventCanon, err := pu.CanonicalizeStruct((ee))
	if err != nil {
		return NotVerified
	}
	eventHash := hash.Encode(eventCanon)
	if h == eventHash.String() {
		return Success
	}
	return Failed
}

func decodeRootProof(rawProof []string) (rootProof, error) {
	proof := make(rootProof, 0, len(rawProof))
	for _, rawItem := range rawProof {
		root, membershipProof, err := splitConsistencyProof(rawItem)
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

func VerifyMembershipProof(rootHashEnc, h string, membershipProof string) (EventVerification, error) {
	targetHash, err := hash.Decode(h)
	if err != nil {
		return Failed, err
	}
	rootHash, err := hash.Decode(rootHashEnc)
	if err != nil {
		return Failed, err
	}
	proof, err := decodeProof(membershipProof)
	if err != nil {
		return Failed, err
	}
	if verifyLogProof(targetHash, rootHash, proof) {
		return Success, nil
	}
	return Failed, nil
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

func verifyConsistencyProof(oldRootHash, newRootHash string, consistencyProof []string) (bool, error) {
	oldHash, err := hash.Decode(oldRootHash)
	if err != nil {
		return false, err
	}
	newHash, err := hash.Decode(newRootHash)
	if err != nil {
		return false, err
	}
	proof, err := decodeRootProof(consistencyProof)
	if err != nil {
		return false, err
	}
	if len(proof) == 0 {
		return false, fmt.Errorf("audit: consistency proof is empty")
	}
	rootHash := proof[0].Hash
	for i := 1; i < len(proof); i++ {
		rootHash = hash.Pair(proof[i].Hash).With(rootHash)
	}
	if !rootHash.Equal(oldHash) {
		return false, nil
	}
	for _, item := range proof {
		if !verifyLogProof(item.Hash, newHash, item.Proof) {
			return false, nil
		}
	}
	return true, nil
}

func treeSizes(root *Root, events SearchEvents) []string {
	treeSizes := make(map[int]struct{}, 0)
	treeSizes[root.Size] = struct{}{}
	for _, event := range events {
		if event.LeafIndex != nil {
			leafIdx := *event.LeafIndex + 1
			treeSizes[leafIdx] = struct{}{}
			if leafIdx > 1 {
				treeSizes[leafIdx-1] = struct{}{}
			}
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

type ValidatedEvent struct {
	// the event that was validated
	Event *EventEnvelope

	// True if the event was successfully validated nil if there is no membership to validate
	MembershipProofStatus *bool

	// True if the event was successfully validated nil if there is no hash to validate
	ConsistencyProofStatus *bool
}

type ValidateEvents []*ValidatedEvent

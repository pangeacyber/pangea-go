package audit

import (
	"fmt"
	"strings"
)

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

package audit

import (
	"fmt"

	"github.com/pangeacyber/go-pangea/pangea/hash"
)

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

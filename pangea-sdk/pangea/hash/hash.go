package hash

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Hash []byte

func Encode(d []byte) Hash {
	b := sha256.Sum256(d)
	return Hash(b[:])
}

func Decode(hash string) (Hash, error) {
	decoded, err := hex.DecodeString(hash)
	if err != nil {
		return nil, fmt.Errorf("hash: invalid hash: %w", err)
	}
	return Hash(decoded), nil
}

func (h Hash) String() string {
	return hex.EncodeToString(h)
}

func (h Hash) Equal(other Hash) bool {
	return bytes.Equal(h, other)
}

// Pair is an intermidiate type to pair to hashes together from left to right.
type Pair []byte

// With returns a new Hash with the given hash appended to the right and rehashed.
func (p Pair) With(h Hash) Hash {
	b := sha256.Sum256(append(p, h...))
	return Hash(b[:])
}

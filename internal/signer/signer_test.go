package signer_test

import (
	"testing"

	"github.com/pangeacyber/go-pangea/internal/signer"
	"github.com/stretchr/testify/assert"
)

func TestSigner(t *testing.T) {
	_, err := signer.NewKeyPairFromFile("../../utils/privkey")
	assert.NoError(t, err)
}

func TestSigner_BadFile(t *testing.T) {
	_, err := signer.NewKeyPairFromFile("Not a file")
	assert.Error(t, err)
}

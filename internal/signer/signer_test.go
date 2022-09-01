package signer_test

import (
	"testing"

	"github.com/pangeacyber/go-pangea/internal/signer"
	"github.com/stretchr/testify/assert"
)

func TestSigner(t *testing.T) {
	_, err := signer.NewSignerFromPrivateKeyFile("./testdata/privkey")
	assert.NoError(t, err)
}

func TestSigner_BadFile(t *testing.T) {
	_, err := signer.NewSignerFromPrivateKeyFile("Not a file")
	assert.Error(t, err)
}

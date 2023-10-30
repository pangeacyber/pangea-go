package signer_test

import (
	"encoding/base64"
	"testing"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/internal/signer"
	"github.com/stretchr/testify/assert"
)

func TestSigner(t *testing.T) {
	s, err := signer.NewSignerFromPrivateKeyFile("./testdata/privkey")
	assert.NoError(t, err)
	msg := "Hello signed world"
	signature, err := s.Sign([]byte(msg))
	signBase64 := base64.StdEncoding.EncodeToString(signature)

	pk, err := s.PublicKey()
	assert.NoError(t, err)
	assert.Equal(t, pk, "-----BEGIN PUBLIC KEY-----\nMCowBQYDK2VwAyEAlvOyDMpK2DQ16NI8G41yINl01wMHzINBahtDPoh4+mE=\n-----END PUBLIC KEY-----\n")
	assert.Equal(t, signBase64, "IYmIUBKWu5yLHM1u3bAw7dvVg1MPc7FLDWSz6d9oqn4FoCu9Bk6ta/lXvvXZUpa7hCm6RhU0VdBzh53x3mKiDQ==")

	v, err := signer.NewVerifierFromPubKey(pk)
	assert.NoError(t, err)
	sDecoded, err := base64.StdEncoding.DecodeString(signBase64)
	res, err := v.Verify([]byte(msg), sDecoded)
	assert.NoError(t, err)
	assert.True(t, res)
	assert.Equal(t, "ED25519", s.GetAlgorithm())
}

func TestVerifier(t *testing.T) {
	pk := "lvOyDMpK2DQ16NI8G41yINl01wMHzINBahtDPoh4+mE="
	s := "IYmIUBKWu5yLHM1u3bAw7dvVg1MPc7FLDWSz6d9oqn4FoCu9Bk6ta/lXvvXZUpa7hCm6RhU0VdBzh53x3mKiDQ=="
	msg := []byte("Hello signed world")
	v, err := signer.NewVerifierFromPubKey(pk)
	assert.NoError(t, err)

	sDecoded, err := base64.StdEncoding.DecodeString(s)
	res, err := v.Verify([]byte(msg), sDecoded)
	assert.NoError(t, err)
	assert.True(t, res)

}

func TestSigner_BadFile(t *testing.T) {
	_, err := signer.NewSignerFromPrivateKeyFile("Not a file")
	assert.Error(t, err)
}

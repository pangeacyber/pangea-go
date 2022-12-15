package signer_test

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/pangeacyber/pangea-go/pangea-sdk/internal/signer"
	"github.com/stretchr/testify/assert"
)

func TestSigner(t *testing.T) {
	s, err := signer.NewSignerFromPrivateKeyFile("./testdata/privkey")
	assert.NoError(t, err)

	signature, err := s.Sign([]byte("Hello signed world"))
	signBase64 := base64.StdEncoding.EncodeToString(signature)

	fmt.Println("Signature base64 is: ", signBase64)
	fmt.Println("Publick key base64 is: ", s.PublicKey())
	assert.Equal(t, s.PublicKey(), "lvOyDMpK2DQ16NI8G41yINl01wMHzINBahtDPoh4+mE=")
	assert.Equal(t, signBase64, "IYmIUBKWu5yLHM1u3bAw7dvVg1MPc7FLDWSz6d9oqn4FoCu9Bk6ta/lXvvXZUpa7hCm6RhU0VdBzh53x3mKiDQ==")
}

func TestSigner_BadFile(t *testing.T) {
	_, err := signer.NewSignerFromPrivateKeyFile("Not a file")
	assert.Error(t, err)
}

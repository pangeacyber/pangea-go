package signer

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

type Signer interface {
	Sign(msg []byte) ([]byte, error)
	PublicKey() string
}

type Verifier interface {
	Verify(msg, sig []byte) bool
}

type signer ed25519.PrivateKey
type verifier ed25519.PublicKey

func NewSignerFromPrivateKeyFile(name string) (signer, error) {
	b, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("signer: cannot read file %v: %w", name, err)
	}

	rawPrivateKey, err := ssh.ParseRawPrivateKey(b)
	if err != nil {
		return nil, fmt.Errorf("signer: cannot parse private key: %w", err)
	}

	privateKey, ok := rawPrivateKey.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("signer: cannot convert to ED25519 key")
	}

	return (signer)(privateKey), nil
}

func (s signer) Sign(msg []byte) ([]byte, error) {
	return (ed25519.PrivateKey)(s).Sign(rand.Reader, msg, crypto.Hash(0))
}

func (s signer) PublicKey() string {
	return base64.StdEncoding.EncodeToString((ed25519.PrivateKey)(s).Public().(ed25519.PublicKey))
}

func NewVerifierFromPubKey(pubkey []byte) Verifier {
	return (verifier)(pubkey)
}

func (v verifier) Verify(msg, sig []byte) bool {
	return ed25519.Verify((ed25519.PublicKey)(v), msg, sig)
}

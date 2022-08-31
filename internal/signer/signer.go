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

type SignerVerifier interface {
	Signer
	Verifier
}

type KeyPair struct {
	priv ed25519.PrivateKey
	pub  ed25519.PublicKey
}

// NewKeyPairFromFile
func NewKeyPairFromFile(name string) (*KeyPair, error) {
	b, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("signer: cannot read file %v: %w", name, err)
	}

	rawPrivateKey, err := ssh.ParseRawPrivateKey(b)
	if err != nil {
		return nil, fmt.Errorf("signer: cannot parse private key: %w", err)
	}

	privateKey, ok := rawPrivateKey.(ed25519.PrivateKey)
	if ok != true {
		return nil, fmt.Errorf("signer: cannot convert to ED25519 key")
	}

	publicKey := privateKey.Public().(ed25519.PublicKey)
	return &KeyPair{
		priv: privateKey,
		pub:  publicKey,
	}, nil

}

func (k *KeyPair) Sign(msg []byte) ([]byte, error) {
	return k.priv.Sign(rand.Reader, msg, crypto.Hash(0))
}

func (k *KeyPair) Verify(msg, sig []byte) bool {
	return ed25519.Verify(k.pub, msg, sig)
}

func (k *KeyPair) PublicKey() string {
	return base64.StdEncoding.EncodeToString(k.pub)
}

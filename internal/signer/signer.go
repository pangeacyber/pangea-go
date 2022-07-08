package signer

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

type Signer interface {
	Sign(msg []byte) ([]byte, error)
}

type Verifier interface {
	Verify(msg, sig []byte) bool
}

type SignerVerifier interface {
	Signer
	Verifier
}

type PrivateKey struct {
	priv ed25519.PrivateKey
	pub  ed25519.PublicKey
}

// NewPrivateKeyFromFile
func NewPrivateKeyFromFile(name string) (*PrivateKey, error) {
	b, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("signer: cannot read file %v: %w", name, err)
	}
	pemBlock, _ := pem.Decode(b)
	if pemBlock == nil {
		return nil, fmt.Errorf("signer: cannot decode file as PEM encoding")
	}

	rawPrivateKey, err := x509.ParsePKCS8PrivateKey(pemBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("signer: cannot decode PKCS8 format in pem file: %w", err)
	}
	privateKey, ok := rawPrivateKey.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("signer: cannot use private key of type %T as ed25519", rawPrivateKey)
	}
	return &PrivateKey{
		priv: privateKey,
		pub:  privateKey.Public().(ed25519.PublicKey),
	}, nil
}

func (k *PrivateKey) Sign(msg []byte) ([]byte, error) {
	return k.priv.Sign(rand.Reader, msg, crypto.Hash(0))
}

func (k *PrivateKey) Verify(msg, sig []byte) bool {
	return ed25519.Verify(k.pub, msg, sig)
}

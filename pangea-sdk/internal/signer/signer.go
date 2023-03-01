package signer

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

type Signer interface {
	Sign(msg []byte) ([]byte, error)
	PublicKey() string
}

type Verifier interface {
	Verify(msg, sig []byte) (bool, error)
}

type signer ed25519.PrivateKey
type ed25519Verifier struct {
	pubkey ed25519.PublicKey
}

func NewSignerFromPrivateKeyFile(name string) (Signer, error) {
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

func NewVerifierFromPubKey(pkPem string) (Verifier, error) {
	if strings.HasPrefix(pkPem, "-----") {
		block, _ := pem.Decode([]byte(pkPem))
		if block == nil {
			return nil, errors.New("Failed to decode PEM block")
		}

		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}

		switch pub := pub.(type) {
		// TODO: Add support for more kind of signatures
		// case *rsa.PublicKey:
		// case *dsa.PublicKey:
		// case *ecdsa.PublicKey:
		case ed25519.PublicKey:
			return ed25519Verifier{
				pubkey: pub,
			}, nil
		default:
			return nil, errors.New("unknown type of public key")
		}
	} else {
		// Done to keep backward compatibility with old key format without header
		pubKey, err := base64.StdEncoding.DecodeString(pkPem)
		if err != nil {
			return nil, err
		}
		return ed25519Verifier{
			pubkey: pubKey,
		}, nil
	}

}

func (v ed25519Verifier) Verify(msg, sig []byte) (bool, error) {
	return ed25519.Verify(v.pubkey, msg, sig), nil
}

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

	"github.com/pangeacyber/pangea-go/pangea-sdk/service/vault"
	"golang.org/x/crypto/ssh"
)

type Signer interface {
	Sign(msg []byte) ([]byte, error)
	PublicKey() (string, error)
	GetAlgorithm() string
}

type Verifier interface {
	Verify(msg, sig []byte) (bool, error)
}

type signerEd25519 ed25519.PrivateKey
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
	if ok {
		return (signerEd25519)(privateKey), nil
	}

	return nil, fmt.Errorf("Not supported key type")
}

func (s signerEd25519) Sign(msg []byte) ([]byte, error) {
	return (ed25519.PrivateKey)(s).Sign(rand.Reader, msg, crypto.Hash(0))
}

func (s signerEd25519) PublicKey() (string, error) {
	// public key
	b, err := x509.MarshalPKIXPublicKey((ed25519.PrivateKey)(s).Public())
	if err != nil {
		return "", err
	}

	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: b,
	}

	pubPEM := string(pem.EncodeToMemory(block))
	return pubPEM, nil
}

func (s signerEd25519) GetAlgorithm() string {
	return string(vault.AAed25519)
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
		case ed25519.PublicKey:
			return ed25519Verifier{
				pubkey: pub,
			}, nil
		// TODO: Add support for more kind of signatures
		// case *rsa.PublicKey:
		// case *dsa.PublicKey:
		// case *ecdsa.PublicKey:
		default:
			return nil, errors.New("Not supported key type")
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

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
	pub  crypto.PublicKey
}

// NewKeyPairFromFile
func NewKeyPairFromFile(name string) (*KeyPair, error) {
	b, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("signer: cannot read file %v: %w", name, err)
	}
	fmt.Println(string(b))
	// pemBlock, _ := pem.Decode(b)
	// if pemBlock == nil {
	// 	return nil, fmt.Errorf("signer: cannot decode file as PEM encoding")
	// }

	rawPrivateKey, err := ssh.ParseRawPrivateKey(b)
	if err != nil {
		fmt.Println("cannot parse raw private key")
		return nil, fmt.Errorf("signer: cannot parse private key: %w", err)
	} else {
		fmt.Println("Parse OK")
	}

	fmt.Println(rawPrivateKey.(*ed25519.PrivateKey))

	privateKey, ok := rawPrivateKey.(*ed25519.PrivateKey)
	if ok != true {
		return nil, fmt.Errorf("signer: cannot convert to ED25519 key")
	}

	realPrivateKey := ed25519.NewKeyFromSeed([]byte(*privateKey))

	fmt.Println(realPrivateKey)

	fmt.Println(privateKey)
	return &KeyPair{
		priv: *privateKey,
		pub:  realPrivateKey.Public(),
	}, nil

	return nil, nil
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

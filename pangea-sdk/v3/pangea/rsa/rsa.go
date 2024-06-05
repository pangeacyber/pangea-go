package rsa

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"

	cryptorand "crypto/rand"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
	// "gitlab.com/pangeacyber/pangea-cli/plugins/vault/common"
)

// Generates rsa key pairs
func GenerateKeyPair(kbIn int) (pubKey *rsa.PublicKey, privKey *rsa.PrivateKey, err error) {
	var (
		keyBits  int
		seedSize = 4096
	)

	keySizes := []int{2048, 3072, 4096, 0}
	for _, keyBits = range keySizes {
		if kbIn == keyBits {
			break
		}
	}
	if keyBits == 0 {
		return nil, nil, fmt.Errorf("invalid key bits value: %d", kbIn)
	}

	rand := cryptorand.Reader
	seed := make([]byte, seedSize)
	if _, err := io.ReadFull(rand, seed); err != nil {
		return nil, nil, fmt.Errorf("generate asymmetric key failed: %w", err)
	}

	privKey, err = rsa.GenerateKey(rand, keyBits)
	if err != nil {
		return nil, nil, fmt.Errorf("generate asymmetric key failed: %w", err)
	}

	pubKey = &privKey.PublicKey

	return pubKey, privKey, nil
}

// Encode Private Key to PKCS1 format embedded in a PEM Block
func EncodePEMPrivateKey(privKey crypto.PrivateKey) ([]byte, error) {
	rsaPrivKey, ok := privKey.(*rsa.PrivateKey)
	if !ok {
		return nil, pangea.ErrInvalidPrivateKey
	}

	encodedKey := x509.MarshalPKCS1PrivateKey(rsaPrivKey)

	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: encodedKey,
	}
	return pem.EncodeToMemory(block), nil
}

// Encode Public Key to PKIX format embedded in a PEM Block
func EncodePEMPublicKey(pubKey crypto.PublicKey) ([]byte, error) {
	encodedKey, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return nil, pangea.ErrInvalidPrivateKey
	}

	block := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: encodedKey,
	}
	return pem.EncodeToMemory(block), nil
}

func DecryptSHA512(rsaKey *rsa.PrivateKey, msgEncr []byte) ([]byte, error) {
	decryptedBytes, err := rsaKey.Decrypt(nil, msgEncr, &rsa.OAEPOptions{Hash: crypto.SHA512})
	if errors.Is(err, rsa.ErrDecryption) {
		fmt.Println(err)
		return nil, pangea.ErrDecryptionFailed
	}
	if err != nil {
		return nil, pangea.ErrInvalidPrivateKey
	}
	return decryptedBytes, nil
}

package vault

import (
	cryptorsa "crypto/rsa"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/internal"
	"golang.org/x/crypto/pbkdf2"
)

type KEMhashAlgorithm string
type KEMKDF string

const (
	KEMHAsha512    KEMhashAlgorithm = "sha512"
	KEMKDFpbkdf2   KEMKDF           = "pbkdf2"
	iterationCount                  = 1000000
)

type KEMDecryptInput struct {
	Cipher               []byte
	EncryptedSalt        []byte
	AsymmetricAlgorithm  string
	AsymmetricPrivateKey cryptorsa.PrivateKey
	SymmetricAlgorithm   string
	Password             string
	HashAlgorithm        KEMhashAlgorithm
	IterationCount       int
	KDF                  KEMKDF
}

func KEMDecrypt(input KEMDecryptInput) ([]byte, error) {
	if input.AsymmetricAlgorithm != string(EEArsa4096_no_padding_kem) {
		return nil, errors.New("unsupported asymmetric algorithm")
	}

	if input.IterationCount < 1 {
		return nil, errors.New("invalid IterationCount value")
	}

	if input.KDF != KEMKDFpbkdf2 {
		return nil, errors.New("unsupported KDF")
	}

	if input.HashAlgorithm != KEMHAsha512 {
		return nil, errors.New("unsupported HashAlgorithm")
	}

	if input.EncryptedSalt == nil {
		return nil, errors.New("EncryptedSalt field should not be nil")
	}

	keyLength, err := GetSymmetricKeyLength(input.SymmetricAlgorithm)
	if err != nil {
		return nil, err
	}

	salt := internal.RSADecryptNoPadding(input.AsymmetricPrivateKey, input.EncryptedSalt)
	symmetricKey := pbkdf2.Key([]byte(input.Password), salt, input.IterationCount, keyLength, sha512.New)

	message, err := AESDecrypt(input.SymmetricAlgorithm, symmetricKey, input.Cipher, nil)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func NewKEMDecryptInput(r ExportResult, password string, rsaPrivKey cryptorsa.PrivateKey) (*KEMDecryptInput, error) {
	// Decide which field to decrypt, just once should be pressent
	cipherEncoded := r.PrivateKey
	if cipherEncoded == nil {
		cipherEncoded = r.Key
	}
	if cipherEncoded == nil {
		return nil, errors.New("no field to decrypt")
	}

	cipher := make([]byte, len(*cipherEncoded))
	n, err := base64.StdEncoding.Decode(cipher, []byte(*cipherEncoded))
	if err != nil {
		return nil, fmt.Errorf("failed to decode cipher. %v", err)
	}
	cipher = cipher[0:n]

	decodedEncryptedSalt := make([]byte, len(r.EncryptedSalt))
	n, err = base64.StdEncoding.Decode(decodedEncryptedSalt, []byte(r.EncryptedSalt))
	if err != nil {
		return nil, fmt.Errorf("failed to decode encrypted salt. %v", err)
	}
	decodedEncryptedSalt = decodedEncryptedSalt[0:n]

	return &KEMDecryptInput{
		Cipher:               cipher,
		EncryptedSalt:        decodedEncryptedSalt,
		AsymmetricAlgorithm:  r.AsymmetricAlgorithm,
		AsymmetricPrivateKey: rsaPrivKey,
		SymmetricAlgorithm:   r.SymmetricAlgorithm,
		Password:             password,
		HashAlgorithm:        KEMhashAlgorithm(r.HashAlgorithm),
		IterationCount:       r.IterationCount,
		KDF:                  KEMKDF(r.KDF),
	}, nil

}

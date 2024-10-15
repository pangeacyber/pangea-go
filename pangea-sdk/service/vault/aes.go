package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

const NONCE_SIZE = 12 // Standard nonce size for GCM

func AESDecrypt(algorithm string, key []byte, msg []byte, additionalData []byte) ([]byte, error) {
	if algorithm != string(SYAaes256_gcm) {
		return nil, fmt.Errorf("unsupported algorithm: '%s'", algorithm)
	}

	return aesDecryptGCM(key, msg, additionalData)
}

func aesDecryptGCM(key []byte, cipherText []byte, additionalData []byte) ([]byte, error) {
	var plaintext []byte

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	nonce := cipherText[:NONCE_SIZE]
	cipherText = cipherText[NONCE_SIZE:]

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err = aesgcm.Open(plaintext, nonce, cipherText, additionalData)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

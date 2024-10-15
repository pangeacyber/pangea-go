package internal

import (
	"crypto/rsa"
	"math/big"
)

func RSADecryptNoPadding(pk rsa.PrivateKey, cipherText []byte) []byte {
	c := new(big.Int).SetBytes(cipherText)
	return c.Exp(c, pk.D, pk.N).Bytes()
}

func RSAEncryptNoPadding(pub rsa.PublicKey, data []byte) []byte {
	encrypted := new(big.Int)
	e := big.NewInt(int64(pub.E))
	payload := new(big.Int).SetBytes(data)
	return encrypted.Exp(payload, e, pub.N).Bytes()
}

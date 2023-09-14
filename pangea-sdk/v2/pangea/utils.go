package pangea

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

func HashSHA256(i string) string {
	b := sha256.Sum256([]byte(i))
	return hex.EncodeToString(b[:])
}

func HashSHA1(i string) string {
	b := sha1.Sum([]byte(i))
	return hex.EncodeToString(b[:])
}

func HashSHA512(i string) string {
	b := sha512.Sum512([]byte(i))
	return hex.EncodeToString(b[:])
}

func GetHashPrefix(h string, len uint) string {
	return h[0:len]
}

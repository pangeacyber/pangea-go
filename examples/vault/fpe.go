// This example demonstrates how to use Vault's format-preserving encryption (FPE)
// to encrypt and decrypt text without changing its length.

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/service/vault"
)

func main() {
	// Set up a Pangea Vault client.
	token := os.Getenv("PANGEA_VAULT_TOKEN")
	if token == "" {
		log.Fatal("Missing `PANGEA_VAULT_TOKEN` environment variable.")
	}

	vaultClient := vault.New(&pangea.Config{
		Token:  token,
		Domain: os.Getenv("PANGEA_DOMAIN"),
	})

	ctx := context.Background()

	// Plain text that we'll encrypt.
	plainText := "123-4567-8901"

	// Optional tweak string.
	tweak := "MTIzMTIzMT=="

	// Generate an encryption key.
	generated, err := vaultClient.SymmetricGenerate(ctx, &vault.SymmetricGenerateRequest{
		Algorithm: vault.SYAaes_ff3_1_256,
		Purpose:   vault.KPfpe,
		CommonGenerateRequest: vault.CommonGenerateRequest{
			Name: "go-fpe-example-" + time.Now().Format(time.RFC3339),
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	keyId := generated.Result.ID

	// Encrypt the plain text.
	encrypted, err := vaultClient.EncryptTransform(ctx, &vault.EncryptTransformRequest{
		ID:        keyId,
		PlainText: plainText,
		Tweak:     &tweak,
		Alphabet:  vault.TAnumeric,
	})
	if err != nil {
		log.Fatal(err)
	}
	encryptedText := encrypted.Result.CipherText
	fmt.Printf("Plain text: %s. Encrypted text: %s.\n", plainText, encryptedText)

	// Decrypt the result to get back the text we started with.
	decrypted, err := vaultClient.DecryptTransform(ctx, &vault.DecryptTransformRequest{
		ID:         keyId,
		CipherText: encryptedText,
		Tweak:      tweak,
		Alphabet:   vault.TAnumeric,
	})
	if err != nil {
		log.Fatal(err)
	}
	decryptedText := decrypted.Result.PlainText
	fmt.Printf("Original text: %s. Decrypted text: %s.\n", plainText, decryptedText)
}

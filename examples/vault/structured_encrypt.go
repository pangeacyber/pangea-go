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
	token := os.Getenv("PANGEA_VAULT_TOKEN")
	if token == "" {
		log.Fatal("missing PANGEA_VAULT_TOKEN environment variable")
	}

	domain := os.Getenv("PANGEA_DOMAIN")
	if domain == "" {
		log.Fatal("missing PANGEA_DOMAIN environment variable")
	}

	vaultClient := vault.New(&pangea.Config{
		Token:  token,
		Domain: domain,
	})

	ctx := context.Background()

	// First create an encryption key, either from the Pangea Console or
	// programmatically as below.
	generateInput := &vault.SymmetricGenerateRequest{
		Algorithm: vault.SYAaes256_cfb,
		Purpose:   vault.KPencryption,
		CommonGenerateRequest: vault.CommonGenerateRequest{
			Name: "Go encrypt example " + time.Now().Format(time.RFC3339),
		},
	}
	generateResponse, err := vaultClient.SymmetricGenerate(ctx, generateInput)
	if err != nil {
		log.Fatal(err)
	}
	encryptionKeyId := generateResponse.Result.ID

	// Structured data that we'll encrypt.
	data := map[string]interface{}{
		"foo":  [4]interface{}{1, 2, "true", "false"},
		"some": "thing",
	}

	encryptInput := &vault.EncryptStructuredRequest{
		ID:             encryptionKeyId,
		StructuredData: data,
		Filter:         "$.foo[2:4]",
	}
	encryptResponse, err := vaultClient.EncryptStructured(ctx, encryptInput)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Encrypted result:", pangea.Stringify(encryptResponse.Result.StructuredData))
}

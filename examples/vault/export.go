// vault encrypt is an example of how to use the encrypt/decrypt methods
package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/pangea/rsa"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/service/vault"
)

func main() {
	token := os.Getenv("PANGEA_VAULT_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	vaultcli := vault.New(&pangea.Config{
		Token:  token,
		Domain: os.Getenv("PANGEA_DOMAIN"),
	})

	ctx := context.Background()

	fmt.Println("Generate key with exportable field on true...")
	name := "Go encrypt example " + time.Now().Format(time.RFC3339)
	generateInput := &vault.SymmetricGenerateRequest{
		Algorithm: vault.SYAaes256_cbc,
		Purpose:   vault.KPencryption,
		CommonGenerateRequest: vault.CommonGenerateRequest{
			Name: name,
		},
		Exportable: pangea.Bool(true),
	}
	generateResponse, err := vaultcli.SymmetricGenerate(ctx, generateInput)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pangea.Stringify(generateResponse.Result))

	id := generateResponse.Result.ID

	// Export with no encryption
	fmt.Println("Exporting key without encryption...")
	rExp, err := vaultcli.Export(ctx,
		&vault.ExportRequest{
			ID: id,
		})

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pangea.Stringify(rExp.Result))

	// Export with encryption
	fmt.Println("Export key encrypted...")

	// Generate a RSA key pair
	rsaPubKey, rsaPrivKey, err := rsa.GenerateKeyPair(4096)
	if err != nil {
		log.Fatal(err)
	}

	// Should send public key in PEM format to encrypt exported key
	rsaPubKeyPEM, err := rsa.EncodePEMPublicKey(rsaPubKey)
	if err != nil {
		log.Fatal(err)
	}

	ea := vault.EEArsa4096_oaep_sha512
	rExpEnc, err := vaultcli.Export(ctx,
		&vault.ExportRequest{
			ID:                  id,
			Version:             pangea.Int(1),
			AsymmetricPublicKey: pangea.String(string(rsaPubKeyPEM)),
			AsymmetricAlgorithm: &ea,
		})
	if err != nil {
		log.Fatal(err)
	}

	// Decode base64 key field
	expKeyDec, err := base64.StdEncoding.DecodeString(*rExpEnc.Result.Key)
	if err != nil {
		log.Fatal(err)
	}

	// Decrypt decoded field
	expKey, err := rsa.DecryptSHA512(rsaPrivKey, expKeyDec)
	if err != nil {
		log.Fatal(err)
	}

	// Use decrypted key
	fmt.Println("Decrypted key:")
	fmt.Println(string(expKey))
}

// vault encrypt is an example of how to use the encrypt/decrypt methods
package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/service/vault"
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

	fmt.Println("Generate key...")
	generateInput := &vault.SymmetricGenerateRequest{
		Algorithm: vault.SYAaes128_cfb,
		Purpose:   vault.KPencryption,
		CommonGenerateRequest: vault.CommonGenerateRequest{
			Name: "My key's name",
		},
	}
	generateResponse, err := vaultcli.SymmetricGenerate(ctx, generateInput)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pangea.Stringify(generateResponse.Result))

	fmt.Println("Encrypt...")
	message := "messagetoencrypt"
	data := base64.StdEncoding.EncodeToString([]byte(message))

	encryptInput := &vault.EncryptRequest{
		ID:        generateResponse.Result.ID,
		PlainText: data,
	}

	encryptResponse, err := vaultcli.Encrypt(ctx, encryptInput)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pangea.Stringify(encryptResponse.Result))

	fmt.Println("Decrypt...")
	decryptInput := &vault.DecryptRequest{
		ID:         generateResponse.Result.ID,
		CipherText: encryptResponse.Result.CipherText,
	}

	decryptResponse, err := vaultcli.Decrypt(ctx, decryptInput)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pangea.Stringify(decryptResponse.Result))

	if decryptResponse.Result.PlainText == data {
		fmt.Println("Encrypt/Decrypt success")
	} else {
		fmt.Println("Encrypt/Decrypt failed")
	}

}

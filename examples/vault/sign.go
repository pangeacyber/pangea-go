// vault sign is an example of how to use the sign/verify methods
package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/service/vault"
)

func main() {
	token := os.Getenv("PANGEA_VAULT_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	vaultcli := vault.New(&pangea.Config{
		Token:           token,
		BaseURLTemplate: os.Getenv("PANGEA_URL_TEMPLATE"),
	})

	ctx := context.Background()

	fmt.Println("Generate key...")
	name := "Go sign example " + time.Now().Format(time.RFC3339)
	generateInput := &vault.AsymmetricGenerateRequest{
		Algorithm: vault.AAed25519,
		Purpose:   vault.KPsigning,
		CommonGenerateRequest: vault.CommonGenerateRequest{
			Name: name,
		},
	}
	generateResponse, err := vaultcli.AsymmetricGenerate(ctx, generateInput)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pangea.Stringify(generateResponse.Result))

	fmt.Println("sign...")
	message := "messagetosign"
	data := base64.StdEncoding.EncodeToString([]byte(message))

	signInput := &vault.SignRequest{
		ID:      generateResponse.Result.ID,
		Message: data,
	}

	signResponse, err := vaultcli.Sign(ctx, signInput)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pangea.Stringify(signResponse.Result))

	fmt.Println("Verify...")
	verifyInput := &vault.VerifyRequest{
		ID:        generateResponse.Result.ID,
		Message:   data,
		Signature: signResponse.Result.Signature,
	}

	verifyResponse, err := vaultcli.Verify(ctx, verifyInput)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pangea.Stringify(verifyResponse.Result))

	if verifyResponse.Result.ValidSignature {
		fmt.Println("Verify success")
	} else {
		fmt.Println("Verify failed")
	}

}

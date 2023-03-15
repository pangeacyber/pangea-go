// embargo check is an example of how to use the check method
package main

import (
	"context"
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

	const (
		secretV1 = "my first secret"
		secretV2 = "my second secret"
	)

	ctx := context.Background()

	fmt.Println("Store secret...")
	storeInput := &vault.SecretStoreRequest{
		Secret: secretV1,
		CommonStoreRequest: vault.CommonStoreRequest{
			Name: "My secret's name",
		},
	}
	storeResponse, err := vaultcli.SecretStore(ctx, storeInput)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pangea.Stringify(storeResponse.Result))

	fmt.Println("Rotate secret...")
	rotateInput := &vault.SecretRotateRequest{
		CommonRotateRequest: vault.CommonRotateRequest{
			ID: storeResponse.Result.ID,
		},
		Secret: secretV1,
	}

	rotateResponse, err := vaultcli.SecretRotate(ctx, rotateInput)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pangea.Stringify(rotateResponse.Result))

	fmt.Println("Get last version")
	getInput := &vault.GetRequest{
		ID: storeResponse.Result.ID,
	}

	getResponse, err := vaultcli.Get(ctx, getInput)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pangea.Stringify(getResponse.Result))

	fmt.Println("Get version 1")
	getInput = &vault.GetRequest{
		ID:      storeResponse.Result.ID,
		Version: "1",
	}

	getResponse, err = vaultcli.Get(ctx, getInput)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pangea.Stringify(getResponse.Result))

}

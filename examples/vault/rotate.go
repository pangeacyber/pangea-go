// vault rotate is an example of how to use the rotate method
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
	name := "Go rotate example " + time.Now().Format(time.RFC3339)
	storeInput := &vault.SecretStoreRequest{
		Secret: secretV1,
		CommonStoreRequest: vault.CommonStoreRequest{
			Name: name,
			Type: vault.ITsecret,
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

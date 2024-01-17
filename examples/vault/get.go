// vault rotate is an example of how to use the rotate method
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/vault"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/audit"
)

func main() {
	token := os.Getenv("PANGEA_VAULT_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	auditTokenID := os.Getenv("PANGEA_AUDIT_TOKEN_ID")
	if auditTokenID == "" {
		log.Fatal("Unauthorized: No audit token id present")
	}
    vaultConfig := pangea.Config{
		Token:  token,
		Domain: os.Getenv("PANGEA_DOMAIN"),
	}
	vaultObj := vault.New(&vaultConfig)

	ctx := context.Background()

	fmt.Println("Fetch the audit token...")
	getRequest := &vault.GetRequest{
        ID: auditTokenID,
	}
	storeResponse, err := vaultObj.Get(ctx, getRequest)
	if err != nil {
		log.Fatal(err)
	}
    auditToken := storeResponse.Result.CurrentVersion.Secret

	fmt.Println("Initialize Log...")
    auditConfig := pangea.Config{
		Token:  *auditToken,
		Domain: os.Getenv("PANGEA_DOMAIN"),
	}
    auditObj, err := audit.New(&auditConfig)
	if err != nil {
		log.Fatal("failed to create audit client")
	}
    event := &audit.StandardEvent{
		Message: "Hello, World!",
	}
    lr, err := auditObj.Log(ctx, event, true)
	if err != nil {
		log.Fatal(err)
	}

	e := (lr.Result.EventEnvelope.Event).(*audit.StandardEvent)
	fmt.Printf("Logged event: %s", pangea.Stringify(e))
}

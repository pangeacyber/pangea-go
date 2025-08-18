// vault get is an example of how to get a token to access audit service
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/option"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/audit"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/vault"
)

func main() {
	token := os.Getenv("PANGEA_VAULT_TOKEN")
	if token == "" {
		log.Fatal("Error: No token present")
	}

	auditTokenID := os.Getenv("PANGEA_AUDIT_TOKEN_VAULT_ID")
	if auditTokenID == "" {
		log.Fatal("Error: No audit token id present")
	}
	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
	)
	if err != nil {
		log.Fatal("Failed to create config: %v", err)
	}
	vaultClient := vault.New(config)

	ctx := context.Background()

	fmt.Println("Fetch the audit token...")
	getRequest := &vault.GetRequest{
		ID: auditTokenID,
	}
	storeResponse, err := vaultClient.Get(ctx, getRequest)
	if err != nil {
		log.Fatal(err)
	}
	auditToken := storeResponse.Result.ItemVersions[0].Token
	if auditToken == nil {
		log.Fatal("Unexpected: token not present")
	}

	fmt.Println("Initialize Log...")
	auditConfig, err := pangea.NewConfig(
		option.WithToken(*auditToken),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
	)
	if err != nil {
		log.Fatal("Failed to create config: %v", err)
	}
	auditClient, err := audit.New(auditConfig)
	if err != nil {
		log.Fatal("failed to create audit client")
	}
	event := &audit.StandardEvent{
		Message: "Hello, World!",
	}
	lr, err := auditClient.Log(ctx, event, true)
	if err != nil {
		log.Fatal(err)
	}

	e := (lr.Result.EventEnvelope.Event).(*audit.StandardEvent)
	fmt.Printf("Logged event: %s", pangea.Stringify(e))
}

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/service/audit"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/service/vault"
)

func main() {
	vaultToken := os.Getenv("PANGEA_VAULT_TOKEN")
	domain := os.Getenv("PANGEA_DOMAIN")
    auditTokenVaultId := os.Getenv("PANGEA_AUDIT_TOKEN_VAULT_ID")
	if vaultToken == "" {
		log.Fatal("Unauthorized: No vault token present")
	}

    vaultcli := vault.New(&pangea.Config{
		Token:  vaultToken,
		Domain: domain,
	})
    getInput := &vault.GetRequest{
		ID: auditTokenVaultId,
	}
	ctx := context.Background()
    getResponse, err := vaultcli.Get(ctx, getInput)
	if err != nil {
		log.Fatal(err)
	}
    auditToken := getResponse.Result.CurrentVersion.Secret
    if auditToken == nil {
		log.Fatal("Unexpected nil auditToken")
    }

	auditcli, err := audit.New(&pangea.Config{
		Token:  *auditToken,
		Domain: domain,
	})
	if err != nil {
		log.Fatal("failed to create audit client")
	}

	event1 := &audit.StandardEvent{
		Message: "Sign up",
		Actor:   "go-sdk",
	}

	resp, err := auditcli.LogBulkAsync(ctx, []any{event1}, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Success. Request_id: %s\n", *resp.RequestID)

}

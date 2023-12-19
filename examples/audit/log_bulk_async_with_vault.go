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
	vault_token := os.Getenv("PANGEA_VAULT_TOKEN")
	domain := os.Getenv("PANGEA_DOMAIN")
    audit_token_vault_id := os.Getenv("PANGEA_AUDIT_TOKEN_VAULT_ID")
	if vault_token == "" {
		log.Fatal("Unauthorized: No vault token present")
	}

    vaultcli := vault.New(&pangea.Config{
		Token:  vault_token,
		Domain: domain,
	})
    getInput := &vault.GetRequest{
		ID: audit_token_vault_id,
	}
	ctx := context.Background()
    getResponse, err := vaultcli.Get(ctx, getInput)
	if err != nil {
		log.Fatal(err)
	}
    audit_token := getResponse.Result.CurrentVersion.Secret

	auditcli, err := audit.New(&pangea.Config{
		Token:  *audit_token,
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

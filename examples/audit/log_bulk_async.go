package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/service/audit"
)

func main() {
	token := os.Getenv("PANGEA_AUDIT_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	auditcli, err := audit.New(&pangea.Config{
		Token:  token,
		Domain: os.Getenv("PANGEA_DOMAIN"),
	})
	if err != nil {
		log.Fatal("failed to create audit client")
	}

	ctx := context.Background()
	event1 := &audit.StandardEvent{
		Message: "Sign up",
		Actor:   "go-sdk",
	}

	event2 := &audit.StandardEvent{
		Message: "Sign in",
		Actor:   "go-sdk",
	}

	fmt.Println("Sending multiple events...")

	resp, err := auditcli.LogBulkAsync(ctx, []any{event1, event2}, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Success. Request_id: %s\n", *resp.RequestID)

}

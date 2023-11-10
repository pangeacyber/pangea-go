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

	_, err = auditcli.LogBulkAsync(ctx, []any{event1, event2}, true)
	_, ok := err.(*pangea.AcceptedError)
	if ok {
		fmt.Println("\tAcceptedError as expected")
	} else {
		fmt.Println("\tUnexpected error")
		fmt.Println(err)
	}

}

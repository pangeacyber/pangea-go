package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/service/audit"
)

func main() {
	token := os.Getenv("PANGEA_TOKEN")
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
	event := audit.Event{
		Message: "Hello, World!",
	}

	fmt.Printf("Logging: %s\n", event.Message)

	logResponse, err := auditcli.Log(ctx, event, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Response: %s", pangea.Stringify(logResponse.Result))
}

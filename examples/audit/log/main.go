package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/go-pangea/pangea"
	"github.com/pangeacyber/go-pangea/service/audit"
)

func main() {

	configID := os.Getenv("AUDIT_CONFIG_ID")
	if configID == "" {
		log.Fatal("Configuration: No config ID present")
	}

	token := os.Getenv("AUDIT_AUTH_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	auditcli, err := audit.New(&pangea.Config{
		Token:    token,
		Domain:   os.Getenv("PANGEA_DOMAIN"),
		CfgToken: configID,
	})
	if err != nil {
		log.Fatal("failed to create audit client")
	}

	ctx := context.Background()
	input := &audit.LogInput{
		Event: &audit.LogEventInput{
			Message: pangea.String("Hello, World!"),
		},
		ReturnHash: pangea.Bool(true),
		Verbose:    pangea.Bool(false),
	}

	fmt.Printf("Logging: %s\n", *input.Event.Message)

	logResponse, err := auditcli.Log(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Response: %s", pangea.Stringify(logResponse.Result))
}

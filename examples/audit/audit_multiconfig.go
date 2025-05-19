package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/option"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/audit"
)

func main() {
	token := os.Getenv("PANGEA_AUDIT_MULTICONFIG_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	configId := os.Getenv("PANGEA_AUDIT_CONFIG_ID")
	if configId == "" {
		log.Fatal("Need to set PANGEA_AUDIT_CONFIG_ID env var")
	}

	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
	)
	if err != nil {
		log.Fatal("failed to create config")
	}

	// Set configId in service construction
	auditcli, err := audit.New(config, audit.WithConfigID(configId))
	if err != nil {
		log.Fatal("failed to create audit client")
	}

	ctx := context.Background()
	event := &audit.StandardEvent{
		Message: "Hello, World!",
	}

	fmt.Printf("Logging: %s\n", event.Message)

	lr, err := auditcli.Log(ctx, event, true)
	if err != nil {
		log.Fatal(err)
	}

	e := (lr.Result.EventEnvelope.Event).(*audit.StandardEvent)
	fmt.Printf("Logged event: %s", pangea.Stringify(e))
}

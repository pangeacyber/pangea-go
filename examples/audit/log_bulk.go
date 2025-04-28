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
	token := os.Getenv("PANGEA_AUDIT_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
	)
	if err != nil {
		log.Fatal("failed to create config")
	}

	auditcli, err := audit.New(config)
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

	fmt.Println("Logging multiple events...")

	lr, err := auditcli.LogBulk(ctx, []any{event1, event2}, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Results:")
	for _, r := range lr.Result.Results {
		e := (r.EventEnvelope.Event).(*audit.StandardEvent)
		fmt.Printf("\tLogged event: %s", pangea.Stringify(e))
	}

}

// Search for events with a custom schema.

package main

import (
	"context"
	"examples/audit/util"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/option"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/audit"
)

func main() {
	token := os.Getenv("PANGEA_AUDIT_CUSTOM_SCHEMA_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
		option.WithPollResultTimeout(60*time.Second),
		option.WithQueuedRetryEnabled(true),
	)
	if err != nil {
		log.Fatal("failed to create config")
	}

	auditcli, err := audit.New(
		config,
		audit.WithCustomSchema(util.CustomSchemaEvent{}),
	)
	if err != nil {
		log.Fatal("failed to create audit client")
	}

	ctx := context.Background()
	input := &audit.SearchInput{
		Query:   `message:""`,
		Limit:   3,
		Verbose: pangea.Bool(false),
	}

	sr, err := auditcli.Search(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	for i, se := range sr.Result.Events {
		ec := (se.EventEnvelope.Event).(*util.CustomSchemaEvent)
		fmt.Printf("Event %d:\n %s\n", i, pangea.Stringify(*ec))
	}
}

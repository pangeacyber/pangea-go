// Search for events with the standard schema.

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/service/audit"
)

func main() {
	token := os.Getenv("PANGEA_AUDIT_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	auditcli, err := audit.New(&pangea.Config{
		Token:              token,
		Domain:             os.Getenv("PANGEA_DOMAIN"),
		PollResultTimeout:  60 * time.Second,
		QueuedRetryEnabled: true,
	})
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
		ec := (se.EventEnvelope.Event).(*audit.StandardEvent)
		fmt.Printf("Event %d:\n %s\n", i, pangea.Stringify(*ec))
	}

}

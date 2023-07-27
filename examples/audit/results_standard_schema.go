package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/audit"
)

func main() {
	token := os.Getenv("PANGEA_AUDIT_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	auditcli, err := audit.New(
		&pangea.Config{
			Token:  token,
			Domain: os.Getenv("PANGEA_DOMAIN"),
		},
	)
	if err != nil {
		log.Fatal("failed to create audit client")
	}

	ctx := context.Background()
	si := &audit.SearchInput{
		Query:   `message:""`,
		Limit:   10,
		Verbose: pangea.Bool(false),
	}

	sr, err := auditcli.Search(ctx, si)
	if err != nil {
		log.Fatal(err)
	}

	ri := &audit.SearchResultsInput{
		ID:    sr.Result.ID,
		Limit: 2,
	}

	rr, err := auditcli.SearchResults(ctx, ri)
	if err != nil {
		log.Fatal(err)
	}

	for i, se := range rr.Result.Events {
		ec := (se.EventEnvelope.Event).(*audit.StandardEvent)
		fmt.Printf("Event %d:\n %s\n", i, pangea.Stringify(*ec))
	}

}

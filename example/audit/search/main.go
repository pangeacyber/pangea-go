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
	token := os.Getenv("PANGEA_AUTH_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	auditcli := audit.New(pangea.Config{
		Token: token,
		EndpointConfig: &pangea.EndpointConfig{
			Scheme: "https",
			CSP:    "aws",
		},
	})

	ctx := context.Background()
	input := &audit.SearchInput{
		Query:                  pangea.String("message:log-123"),
		IncludeMembershipProof: pangea.Bool(true),
	}

	searchOutput, _, err := auditcli.Search(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	verified, err := audit.VerifyAuditRecordsWithArweave(ctx, searchOutput.Root, searchOutput.Events.VerifiableRecords(), true)
	if err != nil {
		log.Fatal(err)
	}

	if !verified {
		log.Fatal("failed validation of audit records")
	}

	fmt.Println(searchOutput.String())
}

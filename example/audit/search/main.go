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

	configID := os.Getenv("AUDIT_CONFIG_ID")
	if token == "" {
		log.Fatal("Configuration: No config ID present")
	}

	auditcli, err := audit.New(&pangea.Config{
		Token:    token,
		Domain:   "dev.pangea.cloud/",
		Insecure: false,
		CfgToken: configID,
	})
	if err != nil {
		log.Fatal("failed to create audit client")
	}

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

	fmt.Println(pangea.Stringify(searchOutput))
}

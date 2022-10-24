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
	input := &audit.SearchInput{
		Query: "message: Hello, World!",
	}

	searchResponse, err := auditcli.Search(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	verified, err := audit.VerifyAuditRecordsWithArweave(ctx, &searchResponse.Result.Root, searchResponse.Result.Events.VerifiableRecords(), true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(pangea.Stringify(searchResponse.Result))

	if len(verified) == 0 {
		log.Fatal("failed validation of audit records")
	}
}

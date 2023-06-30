package main

import (
	"context"
	"examples/audit/util"
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

	auditcli, err := audit.New(&pangea.Config{
		Token:  token,
		Domain: os.Getenv("PANGEA_DOMAIN"),
	})
	if err != nil {
		log.Fatal("failed to create audit client")
	}

	ctx := context.Background()
	input := &audit.SearchInput{
		Query:   "message:\"\"",
		Limit:   2,
		Verbose: pangea.Bool(false),
	}

	searchResponse, err := auditcli.Search(ctx, input, &util.CustomSchemaEvent{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(pangea.Stringify(searchResponse.Result))

}

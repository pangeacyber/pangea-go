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
		Domain:   os.Getenv("PANGEA_DOMAIN"),
		CfgToken: configID,
	})
	if err != nil {
		log.Fatal("failed to create audit client")
	}

	ctx := context.Background()
	input := &audit.LogInput{
		Event: &audit.LogEventInput{
			Message: pangea.String("some important message."),
		},
		ReturnHash: pangea.Bool(true),
		Verbose:    pangea.Bool(true),
	}

	logResponse, err := auditcli.Log(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(pangea.Stringify(logResponse.Result))
}

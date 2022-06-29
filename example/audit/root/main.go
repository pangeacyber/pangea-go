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

	auditcli := audit.New(&pangea.Config{
		Token: token,
		EndpointConfig: &pangea.EndpointConfig{
			Scheme: "https",
			CSP:    "aws",
		},
		CfgToken: configID,
	})

	ctx := context.Background()
	input := &audit.RootInput{
		TreeSize: pangea.Int(10),
	}

	rootOutput, _, err := auditcli.Root(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(pangea.Stringify(rootOutput))
}

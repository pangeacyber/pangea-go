package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/go-pangea/pangea"
	"github.com/pangeacyber/go-pangea/service/redact"
)

func main() {
	token := os.Getenv("PANGEA_AUTH_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	configID := os.Getenv("REDACT_CONFIG_ID")
	if token == "" {
		log.Fatal("Configuration: No config ID present")
	}

	redactcli := redact.New(&pangea.Config{
		Token: token,
		EndpointConfig: &pangea.EndpointConfig{
			Scheme: "https",
			CSP:    "aws",
		},
		CfgToken: configID,
	})

	ctx := context.Background()
	input := &redact.TextInput{
		Text: pangea.String("my phone number is 123-456-7890"),
	}

	redactOutput, _, err := redactcli.Redact(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(pangea.Stringify(redactOutput))
}

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
	token := os.Getenv("REDACT_AUTH_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	configID := os.Getenv("REDACT_CONFIG_ID")
	if token == "" {
		log.Fatal("Configuration: No config ID present")
	}

	redactcli := redact.New(&pangea.Config{
		Token:    token,
		Domain:   os.Getenv("PANGEA_DOMAIN"),
		ConfigID: configID,
	})

	var text = "Hello, my phone number is 123-456-7890"

	fmt.Printf("Redacting PII from: '%s'\n", text)

	ctx := context.Background()
	input := &redact.TextInput{
		Text: pangea.String(text),
	}

	redactResponse, err := redactcli.Redact(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Response: %s", pangea.Stringify(redactResponse.Result))
}

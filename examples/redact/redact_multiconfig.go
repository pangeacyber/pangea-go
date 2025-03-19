package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/service/redact"
)

func main() {
	token := os.Getenv("PANGEA_REDACT_MULTICONFIG_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	configId := os.Getenv("PANGEA_REDACT_CONFIG_ID")
	if token == "" {
		log.Fatal("Need to set PANGEA_REDACT_CONFIG_ID env var")
	}

	// Set configId in service construction
	redactcli := redact.New(&pangea.Config{
		Token:           token,
		BaseURLTemplate: os.Getenv("PANGEA_URL_TEMPLATE"),
	}, redact.WithConfigID(configId))

	var text = "Hello, my phone number is 123-456-7890"

	fmt.Printf("Redacting PII from: '%s'\n", text)

	ctx := context.Background()
	input := &redact.TextRequest{
		Text: pangea.String(text),
	}

	redactResponse, err := redactcli.Redact(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Redacted text: %s", pangea.Stringify(redactResponse.Result.RedactedText))
}

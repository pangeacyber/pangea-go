package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/option"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/redact"
)

func main() {
	token := os.Getenv("PANGEA_REDACT_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
	)
	if err != nil {
		log.Fatalf("expected no error got: %v", err)
	}
	redactcli := redact.New(config)

	var text = "Visit our web is https://pangea.cloud"

	fmt.Printf("Redacting PII from: '%s'\n", text)

	ctx := context.Background()
	input := &redact.TextRequest{
		Text:       pangea.String(text),
		LLMrequest: pangea.Bool(true),
	}

	redactResponse, err := redactcli.Redact(ctx, input)
	if err != nil {
		log.Fatal(err)
	}
	if redactResponse.Result.FPEContext == nil {
		log.Fatal("FPEContext is nil")
	}

	fmt.Printf("Redacted text: %s", pangea.Stringify(redactResponse.Result.RedactedText))

	unredactResponse, err := redactcli.Unredact(ctx, &redact.UnredactRequest{
		RedactedData: redactResponse.Result.RedactedText,
		FPEContext:   *redactResponse.Result.FPEContext,
	})

	if err != nil {
		log.Fatal(err)
	}

	data := unredactResponse.Result.Data.(string)
	fmt.Println("Unredacted text: ", data)
}

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/redact"
)

func main() {
	token := os.Getenv("PANGEA_REDACT_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	redactcli := redact.New(&pangea.Config{
		Token:  token,
		Domain: os.Getenv("PANGEA_DOMAIN"),
	})

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

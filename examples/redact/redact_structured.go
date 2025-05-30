package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/redact"
)

type yourCustomDataStruct struct {
	Secret string `json:"secret"`
}

func main() {
	token := os.Getenv("PANGEA_REDACT_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	redactcli := redact.New(&pangea.Config{
		Token:  token,
		Domain: os.Getenv("PANGEA_DOMAIN"),
	})

	ctx := context.Background()

	data := map[string]any{
		"Secret": "My phone number is 415-867-5309",
	}

	input := &redact.StructuredRequest{
		Data: data,
	}

	redactResponse, err := redactcli.RedactStructured(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(pangea.Stringify(redactResponse.Result.RedactedData))
}

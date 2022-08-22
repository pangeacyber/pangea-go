package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/go-pangea/pangea"
	"github.com/pangeacyber/go-pangea/service/redact"
)

type yourCustomDataStruct struct {
	Secret string `json:"secret"`
}

func main() {
	token := os.Getenv("PANGEA_AUTH_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	configID := os.Getenv("REDACT_CONFIG_ID")
	if token == "" {
		log.Fatal("Configuration: No config ID present")
	}

	redactcli, err := redact.New(&pangea.Config{
		Token:    token,
		Domain:   os.Getenv("PANGEA_DOMAIN"),
		CfgToken: configID,
	})
	if err != nil {
		log.Fatal("failed to create redact client")
	}

	ctx := context.Background()

	input := &redact.StructuredInput{
		JSONP: []*string{
			pangea.String("$.secret"),
		},
	}
	rawData := yourCustomDataStruct{
		Secret: "My social security number is 0303456",
	}
	input.SetData(rawData)

	redactOutput, _, err := redactcli.RedactStructured(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	var redactedData yourCustomDataStruct
	redactOutput.GetRedactedData(&redactedData)

	fmt.Println(pangea.Stringify(redactedData))
}

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
	token := os.Getenv("PANGEA_REDACT_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	redactcli := redact.New(&pangea.Config{
		Token:  token,
		Domain: os.Getenv("PANGEA_DOMAIN"),
	})

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

	redactResponse, err := redactcli.RedactStructured(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	var redactedData yourCustomDataStruct
	redactResponse.Result.GetRedactedData(&redactedData)

	fmt.Println(pangea.Stringify(redactedData))
}

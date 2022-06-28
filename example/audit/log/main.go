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

	auditcli := audit.New(&pangea.Config{
		Token: token,
		EndpointConfig: &pangea.EndpointConfig{
			Scheme: "https",
			CSP:    "aws",
		},
	})

	ctx := context.Background()
	input := &audit.LogInput{
		Event: &audit.LogEventInput{
			Message: pangea.String("some important message."),
		},
		ReturnHash: pangea.Bool(true),
	}

	logOutput, _, err := auditcli.Log(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(pangea.Stringify(logOutput))
}

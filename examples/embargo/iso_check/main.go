// embargo check is an example of how to use the check method
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/go-pangea/pangea"
	"github.com/pangeacyber/go-pangea/service/embargo"
)

func main() {
	token := os.Getenv("EMBARGO_AUTH_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	configID := os.Getenv("EMBARGO_CONFIG_ID")
	if token == "" {
		log.Fatal("Configuration: No config ID present")
	}

	embargocli := embargo.New(&pangea.Config{
		Token:    token,
		Domain:   os.Getenv("PANGEA_DOMAIN"),
		ConfigID: configID,
	})

	ctx := context.Background()
	input := &embargo.ISOCheckInput{
		ISOCode: pangea.String("CU"),
	}

	fmt.Printf("Checking Embargo ISO code: '%s'\n", *input.ISOCode)

	checkResponse, err := embargocli.ISOCheck(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Response: %s", pangea.Stringify(checkResponse.Result))
}

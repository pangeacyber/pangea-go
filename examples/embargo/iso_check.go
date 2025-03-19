// embargo check is an example of how to use the check method
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/service/embargo"
)

func main() {
	token := os.Getenv("PANGEA_EMBARGO_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	embargocli := embargo.New(&pangea.Config{
		Token:           token,
		BaseURLTemplate: os.Getenv("PANGEA_URL_TEMPLATE"),
	})

	ctx := context.Background()
	input := &embargo.ISOCheckRequest{
		ISOCode: "CU",
	}

	fmt.Printf("Checking Embargo ISO code: '%s'\n", input.ISOCode)

	checkResponse, err := embargocli.ISOCheck(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Response: %s", pangea.Stringify(checkResponse.Result))
}

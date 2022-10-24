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
	token := os.Getenv("PANGEA_EMBARGO_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	embargocli := embargo.New(&pangea.Config{
		Token:  token,
		Domain: os.Getenv("PANGEA_DOMAIN"),
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

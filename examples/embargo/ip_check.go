// embargo check is an example of how to use the check method
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/service/embargo"
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
	input := &embargo.IPCheckInput{
		IP: pangea.String("213.24.238.26"),
	}

	checkResponse, err := embargocli.IPCheck(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(pangea.Stringify(checkResponse.Result))
}

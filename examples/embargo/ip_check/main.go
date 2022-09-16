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
	token := os.Getenv("PANGEA_AUTH_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	configID := os.Getenv("EMBARGO_CONFIG_ID")
	if token == "" {
		log.Fatal("Configuration: No config ID present")
	}

	embargocli, err := embargo.New(&pangea.Config{
		Token:    token,
		Domain:   os.Getenv("PANGEA_DOMAIN"),
		ConfigID: configID,
	})
	if err != nil {
		log.Fatal("failed to create embargo client")
	}

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

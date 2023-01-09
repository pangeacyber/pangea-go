// intel domain lookup is an example of how to use the lookup method
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/service/url_intel"
)

func main() {
	token := os.Getenv("PANGEA_URL_INTEL_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	intelcli := url_intel.New(&pangea.Config{
		Token:  token,
		Domain: os.Getenv("PANGEA_DOMAIN"),
	})

	ctx := context.Background()
	input := &url_intel.UrlLookupRequest{
		Url:      "http://113.235.101.11:54384",
		Raw:      true,
		Verbose:  true,
		Provider: "crowdstrike",
	}

	response, err := intelcli.Lookup(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(pangea.Stringify(response.Result))
}

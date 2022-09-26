// intel domain lookup is an example of how to use the lookup method
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/go-pangea/pangea"
	"github.com/pangeacyber/go-pangea/service/domain_intel"
)

func main() {
	token := os.Getenv("INTEL_AUTH_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	configID := os.Getenv("INTEL_DOMAIN_CONFIG_ID")
	if token == "" {
		log.Fatal("Configuration: No config ID present")
	}

	intelcli := domain_intel.New(&pangea.Config{
		Token:    token,
		Domain:   os.Getenv("PANGEA_DOMAIN"),
		ConfigID: configID,
	})

	ctx := context.Background()
	input := &domain_intel.DomainLookupInput{
		Domain:   "teoghehofuuxo.su",
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

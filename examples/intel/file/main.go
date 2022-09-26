// intel file lookup is an example of how to use the lookup method
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/go-pangea/pangea"
	"github.com/pangeacyber/go-pangea/service/file_intel"
)

func main() {
	token := os.Getenv("INTEL_AUTH_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	configID := os.Getenv("INTEL_FILE_CONFIG_ID")
	if token == "" {
		log.Fatal("Configuration: No config ID present")
	}

	intelcli := file_intel.New(&pangea.Config{
		Token:    token,
		Domain:   os.Getenv("PANGEA_DOMAIN"),
		ConfigID: configID,
	})

	ctx := context.Background()
	input := &file_intel.FileLookupInput{
		Hash:     "322ccbd42b7e4fd3a9d0167ca2fa9f6483d9691364c431625f1df542706",
		HashType: "sha256",
		Raw:      true,
		Verbose:  true,
		Provider: "reversinglabs",
	}
	response, err := intelcli.Lookup(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(pangea.Stringify(response.Result))
}

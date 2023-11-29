// Example of how to lookup a file's reputation by hash given a local file path
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/file_intel"
)

func main() {
	token := os.Getenv("PANGEA_INTEL_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	intelcli := file_intel.New(&pangea.Config{
		Token:  token,
		Domain: os.Getenv("PANGEA_DOMAIN"),
	})

	ctx := context.Background()
	input, err := file_intel.NewFileReputationRequestFromFilepath("./go.mod")
	if err != nil {
		log.Fatal(err)
		return
	}

	input.Raw = pangea.Bool(true)
	input.Verbose = pangea.Bool(true)

	resp, err := intelcli.Reputation(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(pangea.Stringify(resp.Result.Data))
}

// intel domain lookup is an example of how to use the lookup method
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/ip_intel"
)

func main() {
	fmt.Println("Checking if IP is a proxy...")
	token := os.Getenv("PANGEA_INTEL_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	intelcli := ip_intel.New(&pangea.Config{
		Token:  token,
		Domain: os.Getenv("PANGEA_DOMAIN"),
	})

	ctx := context.Background()
	input := &ip_intel.IpProxyRequest{
		Ip:       "34.201.32.172",
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "digitalelement",
	}

	resp, err := intelcli.IsProxy(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	if resp.Result.Data.IsProxy {
		fmt.Println("IP is a proxy")
	} else {
		fmt.Println("IP is not a proxy")
	}
}

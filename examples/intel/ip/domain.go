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
	fmt.Println("Checking IP's domain...")
	token := os.Getenv("PANGEA_INTEL_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	intelcli := ip_intel.New(&pangea.Config{
		Token:  token,
		Domain: os.Getenv("PANGEA_DOMAIN"),
	})

	ctx := context.Background()
	input := &ip_intel.IpDomainRequest{
		Ip:       "24.235.114.61",
		Raw:      true,
		Verbose:  true,
		Provider: "digitalelement",
	}

	resp, err := intelcli.GetDomain(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	if resp.Result.Data.DomainFound {
		fmt.Printf("IP's domain is %s\n", resp.Result.Data.Domain)
	} else {
		fmt.Println("IP's domain not found")
	}
}

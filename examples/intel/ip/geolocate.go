// intel domain lookup is an example of how to use the lookup method
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/service/ip_intel"
)

func main() {
	token := os.Getenv("PANGEA_IP_INTEL_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	intelcli := ip_intel.New(&pangea.Config{
		Token:  token,
		Domain: os.Getenv("PANGEA_DOMAIN"),
	})

	ctx := context.Background()
	input := &ip_intel.IpGeolocateRequest{
		Ip:       "93.231.182.110",
		Raw:      true,
		Verbose:  true,
		Provider: "digitalenvoy",
	}

	response, err := intelcli.Geolocate(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(pangea.Stringify(response.Result))
}

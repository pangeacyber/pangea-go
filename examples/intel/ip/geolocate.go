// Example of how to look up geolocation information for an IP address
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
	fmt.Println("Geolocating IP...")
	token := os.Getenv("PANGEA_INTEL_TOKEN")
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
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "digitalelement",
	}

	resp, err := intelcli.Geolocate(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(pangea.Stringify(resp.Result.Data))
}

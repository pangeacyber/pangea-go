// intel domain lookup is an example of how to use the lookup method
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/ip_intel"
)

func PrintData(ip string, data ip_intel.GeolocateData) {
	fmt.Printf("\t Indicator: %s\n", ip)
	fmt.Printf("\t\t Country: %s\n", data.Country)
	fmt.Printf("\t\t City: %s\n", data.City)
	fmt.Printf("\t\t Latitude: %f\n", data.Latitude)
	fmt.Printf("\t\t Longitude: %f\n", data.Longitude)
	fmt.Printf("\t\t PostalCode: %s\n", data.PostalCode)
	fmt.Printf("\t\t CountryCode: %s\n", data.CountryCode)
}

func PrintBulkData(data map[string]ip_intel.GeolocateData) {
	for k, v := range data {
		PrintData(k, v)
	}
}

func main() {
	fmt.Println("Geolocating IP...")
	token := os.Getenv("PANGEA_INTEL_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	intelcli := ip_intel.New(&pangea.Config{
		Token:              token,
		Domain:             os.Getenv("PANGEA_DOMAIN"),
		QueuedRetryEnabled: true,
		PollResultTimeout:  60 * time.Second,
	})

	ctx := context.Background()
	input := &ip_intel.IpGeolocateBulkRequest{
		Ips:     []string{"93.231.182.110", "190.28.74.251"},
		Raw:     pangea.Bool(true),
		Verbose: pangea.Bool(true),
	}

	resp, err := intelcli.GeolocateBulk(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Result:")
	PrintBulkData(resp.Result.Data)
}

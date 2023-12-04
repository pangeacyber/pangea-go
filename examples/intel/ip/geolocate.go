// Example of how to look up geolocation information for an IP address
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/service/ip_intel"
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
	ip := "93.231.182.110"
	input := &ip_intel.IpGeolocateRequest{
		Ip:       ip,
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "digitalelement",
	}

	resp, err := intelcli.Geolocate(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Result:")
	PrintData(ip, resp.Result.Data)
}

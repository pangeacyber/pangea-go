// Example of how to look up a domain for an IP address
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/service/ip_intel"
)

func PrintData(ip string, data ip_intel.DomainData) {
	if data.DomainFound {
		fmt.Printf("\t IP %s domain is %s\n", ip, data.Domain)
	} else {
		fmt.Println("\t IP %s domain not found\n", ip)
	}
}

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
	ip := "24.235.114.61"
	input := &ip_intel.IpDomainRequest{
		Ip:       ip,
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "digitalelement",
	}

	resp, err := intelcli.GetDomain(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Result:")
	PrintData(ip, resp.Result.Data)
}

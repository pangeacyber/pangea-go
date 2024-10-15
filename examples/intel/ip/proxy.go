// Example of how to look up if an IP address belongs to a proxy service
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/service/ip_intel"
)

func PrintData(ip string, data ip_intel.ProxyData) {
	if data.IsProxy {
		fmt.Printf("\t IP %s is a proxy\n", ip)
	} else {
		fmt.Printf("\t IP %s is not a proxy\n", ip)
	}
}

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
	ip := "34.201.32.172"
	input := &ip_intel.IpProxyRequest{
		Ip:       ip,
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "digitalelement",
	}

	resp, err := intelcli.IsProxy(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Result:")
	PrintData(ip, resp.Result.Data)
}

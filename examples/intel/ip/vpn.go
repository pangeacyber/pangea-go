// Example of how to look up if an IP address belongs to a VPN service
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
	fmt.Println("Checking if IP is a vpn...")
	token := os.Getenv("PANGEA_INTEL_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	intelcli := ip_intel.New(&pangea.Config{
		Token:  token,
		Domain: os.Getenv("PANGEA_DOMAIN"),
	})

	ctx := context.Background()
	input := &ip_intel.IpVPNRequest{
		Ip:       "2.56.189.74",
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "digitalelement",
	}

	resp, err := intelcli.IsVPN(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	if resp.Result.Data.IsVPN {
		fmt.Println("IP is a VPN")
	} else {
		fmt.Println("IP is not a VPN")
	}
}

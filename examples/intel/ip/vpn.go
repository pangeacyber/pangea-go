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
		Raw:      true,
		Verbose:  true,
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

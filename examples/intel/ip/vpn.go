// Example of how to look up if an IP address belongs to a VPN service
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/service/ip_intel"
)

func PrintData(ip string, data ip_intel.VPNData) {
	if data.IsVPN {
		fmt.Printf("\t IP %s is a VPN\n", ip)
	} else {
		fmt.Printf("\t IP %s is not a VPN\n", ip)
	}
}

func main() {
	fmt.Println("Checking if IP is a vpn...")
	token := os.Getenv("PANGEA_INTEL_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	intelcli := ip_intel.New(&pangea.Config{
		Token:           token,
		BaseURLTemplate: os.Getenv("PANGEA_URL_TEMPLATE"),
	})

	ctx := context.Background()
	ip := "2.56.189.74"
	input := &ip_intel.IpVPNRequest{
		Ip:       ip,
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "digitalelement",
	}

	resp, err := intelcli.IsVPN(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Result:")
	PrintData(ip, resp.Result.Data)
}

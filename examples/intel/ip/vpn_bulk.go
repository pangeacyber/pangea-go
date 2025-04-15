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

func PrintData(ip string, data ip_intel.VPNData) {
	if data.IsVPN {
		fmt.Printf("\t IP %s is a VPN\n", ip)
	} else {
		fmt.Printf("\t IP %s is not a VPN\n", ip)
	}
}

func PrintBulkData(data map[string]ip_intel.VPNData) {
	for k, v := range data {
		PrintData(k, v)
	}
}

func main() {
	fmt.Println("Checking if IP is a vpn...")
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
	input := &ip_intel.IpVPNBulkRequest{
		Ips:     []string{"2.56.189.74", "190.28.74.251"},
		Raw:     pangea.Bool(true),
		Verbose: pangea.Bool(true),
	}

	resp, err := intelcli.IsVPNBulk(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Result:")
	PrintBulkData(resp.Result.Data)
}

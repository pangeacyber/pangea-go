// intel domain lookup is an example of how to use the lookup method
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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

func PrintBulkData(data map[string]ip_intel.DomainData) {
	for k, v := range data {
		PrintData(k, v)
	}
}

func main() {
	fmt.Println("Checking IP's domain...")
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
	input := &ip_intel.IpDomainBulkRequest{
		Ips:     []string{"93.231.182.110", "190.28.74.251"},
		Raw:     pangea.Bool(true),
		Verbose: pangea.Bool(true),
	}

	resp, err := intelcli.GetDomainBulk(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Result:")
	PrintBulkData(resp.Result.Data)
}

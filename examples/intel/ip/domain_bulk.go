// intel domain lookup is an example of how to use the lookup method
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/option"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/ip_intel"
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

	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
		option.WithQueuedRetryEnabled(true),
		option.WithPollResultTimeout(60*time.Second),
	)
	if err != nil {
		log.Fatalf("expected no error got: %v", err)
	}
	intelcli := ip_intel.New(config)

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

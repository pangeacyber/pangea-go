// intel domain lookup is an example of how to use the lookup method
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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

func PrintBulkData(data map[string]ip_intel.ProxyData) {
	for k, v := range data {
		PrintData(k, v)
	}
}

func main() {
	fmt.Println("Checking if IP is a proxy...")
	token := os.Getenv("PANGEA_INTEL_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	intelcli := ip_intel.New(&pangea.Config{
		Token:              token,
		BaseURLTemplate:    os.Getenv("PANGEA_URL_TEMPLATE"),
		QueuedRetryEnabled: true,
		PollResultTimeout:  60 * time.Second,
	})

	ctx := context.Background()
	input := &ip_intel.IpProxyBulkRequest{
		Ips:     []string{"34.201.32.172", "190.28.74.251"},
		Raw:     pangea.Bool(true),
		Verbose: pangea.Bool(true),
	}

	resp, err := intelcli.IsProxyBulk(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Result:")
	PrintBulkData(resp.Result.Data)
}

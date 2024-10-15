// Example of how to look up an IP's reputation using reversinglabs
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/service/ip_intel"
)

func PrintData(indicator string, data ip_intel.ReputationData) {
	fmt.Printf("\t Indicator: %s\n", indicator)
	fmt.Printf("\t\t Verdict: %s\n", data.Verdict)
	fmt.Printf("\t\t Score: %d\n", data.Score)
	fmt.Printf("\t\t Category: %s\n", pangea.Stringify(data.Category))
}

func main() {
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
	input := &ip_intel.IpReputationRequest{
		Ip:       ip,
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "cymru",
	}

	resp, err := intelcli.Reputation(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Result:")
	PrintData(ip, resp.Result.Data)
}

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
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/url_intel"
)

func PrintData(indicator string, data url_intel.ReputationData) {
	fmt.Printf("\t Indicator: %s\n", indicator)
	fmt.Printf("\t\t Verdict: %s\n", data.Verdict)
	fmt.Printf("\t\t Score: %d\n", data.Score)
	fmt.Printf("\t\t Category: %s\n", pangea.Stringify(data.Category))
}

func PrintBulkData(data map[string]url_intel.ReputationData) {
	for k, v := range data {
		PrintData(k, v)
	}
}

func main() {
	fmt.Println("Checking URL...")
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
	intelcli := url_intel.New(config)

	ctx := context.Background()
	input := &url_intel.UrlReputationBulkRequest{
		Urls:     []string{"http://113.235.101.11:54384", "http://45.14.49.109:54819"},
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "crowdstrike",
	}

	resp, err := intelcli.ReputationBulk(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Result:")
	PrintBulkData(resp.Result.Data)
}

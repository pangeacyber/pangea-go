package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/option"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/url_intel"
	"github.com/spf13/cobra"
)

func init() {
	urlCmd := &cobra.Command{
		Use:   "url",
		Short: "URL intel examples",
	}
	intelCmd.AddCommand(urlCmd)

	urlReputationCmd := &cobra.Command{
		Use:  "reputation",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return urlReputation(cmd)
		},
	}
	urlCmd.AddCommand(urlReputationCmd)

	urlReputationBulkCmd := &cobra.Command{
		Use:  "reputation_bulk",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return urlReputationBulk(cmd)
		},
	}
	urlCmd.AddCommand(urlReputationBulkCmd)
}

func PrintURLData(indicator string, data url_intel.ReputationData) {
	fmt.Printf("\t Indicator: %s\n", indicator)
	fmt.Printf("\t\t Verdict: %s\n", data.Verdict)
	fmt.Printf("\t\t Score: %d\n", data.Score)
	fmt.Printf("\t\t Category: %s\n", pangea.Stringify(data.Category))
}

func urlReputation(cmd *cobra.Command) error {
	fmt.Println("Checking URL...")
	token := os.Getenv("PANGEA_INTEL_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
	)
	if err != nil {
		log.Fatalf("expected no error got: %v", err)
	}
	intelcli := url_intel.New(config)

	ctx := context.Background()
	indicator := "http://113.235.101.11:54384"
	input := &url_intel.UrlReputationRequest{
		Url:      indicator,
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "crowdstrike",
	}

	resp, err := intelcli.Reputation(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Result:")
	PrintURLData(indicator, resp.Result.Data)
	return nil
}

func PrintURLBulkData(data map[string]url_intel.ReputationData) {
	for k, v := range data {
		PrintURLData(k, v)
	}
}

func urlReputationBulk(cmd *cobra.Command) error {
	fmt.Println("Checking URL...")
	token := os.Getenv("PANGEA_INTEL_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
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
	PrintURLBulkData(resp.Result.Data)
	return nil
}

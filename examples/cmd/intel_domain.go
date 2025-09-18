package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/option"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/domain_intel"
	"github.com/spf13/cobra"
)

func init() {
	domainCmd := &cobra.Command{
		Use:   "domain",
		Short: "Domain intel examples",
	}
	intelCmd.AddCommand(domainCmd)

	domainReputationCmd := &cobra.Command{
		Use:  "reputation",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return domainReputation(cmd)
		},
	}
	domainCmd.AddCommand(domainReputationCmd)

	domainReputationBulkCmd := &cobra.Command{
		Use:  "reputation_bulk",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return domainReputationBulk(cmd)
		},
	}
	domainCmd.AddCommand(domainReputationBulkCmd)

	domainWhoisCmd := &cobra.Command{
		Use:  "whois",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return domainWhois(cmd)
		},
	}
	domainCmd.AddCommand(domainWhoisCmd)
}

func PrintData(indicator string, data domain_intel.ReputationData) {
	fmt.Printf("\t Indicator: %s\n", indicator)
	fmt.Printf("\t\t Verdict: %s\n", data.Verdict)
	fmt.Printf("\t\t Score: %d\n", data.Score)
	fmt.Printf("\t\t Category: %s\n", pangea.Stringify(data.Category))
}

func domainReputation(cmd *cobra.Command) error {
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
	intelcli := domain_intel.New(config)

	ctx := context.Background()
	indicator := "737updatesboeing.com"
	input := &domain_intel.DomainReputationRequest{
		Domain:   indicator,
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "domaintools",
	}

	resp, err := intelcli.Reputation(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Result:")
	PrintData(indicator, resp.Result.Data)
	return nil
}

func PrintBulkData(data map[string]domain_intel.ReputationData) {
	for k, v := range data {
		PrintData(k, v)
	}
}

func domainReputationBulk(cmd *cobra.Command) error {
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
	intelcli := domain_intel.New(config)

	ctx := context.Background()
	input := &domain_intel.DomainReputationBulkRequest{
		Domains:  []string{"pemewizubidob.cafij.co.za", "redbomb.com.tr", "kmbk8.hicp.net"},
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
	return nil
}

func domainWhois(cmd *cobra.Command) error {
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
	intelcli := domain_intel.New(config)

	ctx := context.Background()
	input := &domain_intel.DomainWhoIsRequest{
		Domain:   "737updatesboeing.com",
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "whoisxml",
	}

	resp, err := intelcli.WhoIs(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(pangea.Stringify(resp.Result.Data))
	return nil
}

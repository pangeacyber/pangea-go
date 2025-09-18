package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/option"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/embargo"
	"github.com/spf13/cobra"
)

func init() {
	embargoCmd := &cobra.Command{
		Use:   "embargo",
		Short: "Embargo examples",
	}

	ipCheckCmd := &cobra.Command{
		Use:  "ip_check",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return embargoIPCheck(cmd)
		},
	}
	embargoCmd.AddCommand(ipCheckCmd)

	isoCheckCmd := &cobra.Command{
		Use:  "iso_check",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return embargoISOCheck(cmd)
		},
	}
	embargoCmd.AddCommand(isoCheckCmd)

	ExamplesCmd.AddCommand(embargoCmd)
}

func embargoISOCheck(cmd *cobra.Command) error {
	token := os.Getenv("PANGEA_EMBARGO_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	config, err := pangea.NewConfig(option.WithToken(token), option.WithDomain(os.Getenv("PANGEA_DOMAIN")))
	if err != nil {
		log.Fatalf("expected no error got: %v", err)
	}
	embargocli := embargo.New(config)

	ctx := context.Background()
	input := &embargo.ISOCheckRequest{
		ISOCode: "CU",
	}

	fmt.Printf("Checking Embargo ISO code: '%s'\n", input.ISOCode)

	checkResponse, err := embargocli.ISOCheck(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Response: %s", pangea.Stringify(checkResponse.Result))
	return nil
}

func embargoIPCheck(cmd *cobra.Command) error {
	token := os.Getenv("PANGEA_EMBARGO_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	config, err := pangea.NewConfig(option.WithToken(token), option.WithDomain(os.Getenv("PANGEA_DOMAIN")))
	if err != nil {
		log.Fatalf("expected no error got: %v", err)
	}
	embargocli := embargo.New(config)

	ctx := context.Background()
	input := &embargo.IPCheckRequest{
		IP: "213.24.238.26",
	}

	checkResponse, err := embargocli.IPCheck(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(pangea.Stringify(checkResponse.Result))
	return nil
}

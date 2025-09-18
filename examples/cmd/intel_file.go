package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/option"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/file_intel"
	"github.com/spf13/cobra"
)

func init() {
	fileCmd := &cobra.Command{
		Use:   "file",
		Short: "File intel examples",
	}
	intelCmd.AddCommand(fileCmd)

	fileReputationCmd := &cobra.Command{
		Use:  "reputation",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return fileReputation(cmd)
		},
	}
	fileCmd.AddCommand(fileReputationCmd)

	fileReputationBulkCmd := &cobra.Command{
		Use:  "reputation_bulk",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return fileReputationBulk(cmd)
		},
	}
	fileCmd.AddCommand(fileReputationBulkCmd)

	fileFilepathReputationCmd := &cobra.Command{
		Use:  "filepath_reputation",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return fileFilepathReputation(cmd)
		},
	}
	fileCmd.AddCommand(fileFilepathReputationCmd)
}

func PrintFileData(indicator string, data file_intel.ReputationData) {
	fmt.Printf("\t Indicator: %s\n", indicator)
	fmt.Printf("\t\t Verdict: %s\n", data.Verdict)
	fmt.Printf("\t\t Score: %d\n", data.Score)
	fmt.Printf("\t\t Category: %s\n", pangea.Stringify(data.Category))
}

func fileReputation(cmd *cobra.Command) error {
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
	intelcli := file_intel.New(config)

	ctx := context.Background()
	indicator := "142b638c6a60b60c7f9928da4fb85a5a8e1422a9ffdc9ee49e17e56ccca9cf6e"
	input := &file_intel.FileReputationRequest{
		Hash:     indicator,
		HashType: "sha256",
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "reversinglabs",
	}
	resp, err := intelcli.Reputation(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Result:")
	PrintFileData(indicator, resp.Result.Data)
	return nil
}

func PrintFileBulkData(data map[string]file_intel.ReputationData) {
	for k, v := range data {
		PrintFileData(k, v)
	}
}

func fileReputationBulk(cmd *cobra.Command) error {
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
	intelcli := file_intel.New(config)

	ctx := context.Background()
	input := &file_intel.FileReputationBulkRequest{
		Hashes:   []string{"142b638c6a60b60c7f9928da4fb85a5a8e1422a9ffdc9ee49e17e56ccca9cf6e", "179e2b8a4162372cd9344b81793cbf74a9513a002eda3324e6331243f3137a63"},
		HashType: "sha256",
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "reversinglabs",
	}
	resp, err := intelcli.ReputationBulk(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Result:")
	PrintFileBulkData(resp.Result.Data)
	return nil
}

func fileFilepathReputation(cmd *cobra.Command) error {
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
	intelcli := file_intel.New(config)

	ctx := context.Background()
	input, err := file_intel.NewFileReputationRequestFromFilepath("./go.mod")
	if err != nil {
		log.Fatal(err)
		return err
	}

	input.Raw = pangea.Bool(true)
	input.Verbose = pangea.Bool(true)

	resp, err := intelcli.Reputation(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(pangea.Stringify(resp.Result.Data))
	return nil
}

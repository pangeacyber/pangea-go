package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/option"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/redact"
	"github.com/spf13/cobra"
)

func init() {
	redactCmd := &cobra.Command{
		Use:   "redact",
		Short: "Redact examples",
	}

	redactTextCmd := &cobra.Command{
		Use:  "redact",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return redactText(cmd)
		},
	}
	redactCmd.AddCommand(redactTextCmd)

	redactLLMRequestCmd := &cobra.Command{
		Use:  "redact_llm_request",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return redactLLMRequest(cmd)
		},
	}
	redactCmd.AddCommand(redactLLMRequestCmd)

	redactMulticonfigCmd := &cobra.Command{
		Use:  "redact_multiconfig",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return redactMulticonfig(cmd)
		},
	}
	redactCmd.AddCommand(redactMulticonfigCmd)

	redactStructuredCmd := &cobra.Command{
		Use:  "redact_structured",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return redactStructured(cmd)
		},
	}
	redactCmd.AddCommand(redactStructuredCmd)

	unredactCmd := &cobra.Command{
		Use:  "unredact",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return unredact(cmd)
		},
	}
	redactCmd.AddCommand(unredactCmd)

	ExamplesCmd.AddCommand(redactCmd)
}

func redactLLMRequest(cmd *cobra.Command) error {
	token := os.Getenv("PANGEA_REDACT_TOKEN")
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
	redactcli := redact.New(config)

	var text = "Visit our web is https://pangea.cloud"

	fmt.Printf("Redacting PII from: '%s'\n", text)

	ctx := context.Background()
	input := &redact.TextRequest{
		Text:       pangea.String(text),
		LLMrequest: pangea.Bool(true),
	}

	redactResponse, err := redactcli.Redact(ctx, input)
	if err != nil {
		log.Fatal(err)
	}
	if redactResponse.Result.FPEContext == nil {
		log.Fatal("FPEContext is nil")
	}

	fmt.Printf("Redacted text: %s", pangea.Stringify(redactResponse.Result.RedactedText))

	unredactResponse, err := redactcli.Unredact(ctx, &redact.UnredactRequest{
		RedactedData: redactResponse.Result.RedactedText,
		FPEContext:   *redactResponse.Result.FPEContext,
	})

	if err != nil {
		log.Fatal(err)
	}

	data := unredactResponse.Result.Data.(string)
	fmt.Println("Unredacted text: ", data)
	return nil
}

func redactMulticonfig(cmd *cobra.Command) error {
	token := os.Getenv("PANGEA_REDACT_MULTICONFIG_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	configId := os.Getenv("PANGEA_REDACT_CONFIG_ID")
	if token == "" {
		log.Fatal("Need to set PANGEA_REDACT_CONFIG_ID env var")
	}

	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
	)
	if err != nil {
		log.Fatalf("expected no error got: %v", err)
	}

	// Set configId in service construction
	redactcli := redact.New(config, redact.WithConfigID(configId))

	var text = "Hello, my phone number is 123-456-7890"

	fmt.Printf("Redacting PII from: '%s'\n", text)

	ctx := context.Background()
	input := &redact.TextRequest{
		Text: pangea.String(text),
	}

	redactResponse, err := redactcli.Redact(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Redacted text: %s", pangea.Stringify(redactResponse.Result.RedactedText))
	return nil
}

func redactStructured(cmd *cobra.Command) error {
	token := os.Getenv("PANGEA_REDACT_TOKEN")
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
	redactcli := redact.New(config)

	ctx := context.Background()

	data := map[string]any{
		"Secret": "My phone number is 415-867-5309",
	}

	input := &redact.StructuredRequest{
		Data: data,
	}

	redactResponse, err := redactcli.RedactStructured(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(pangea.Stringify(redactResponse.Result.RedactedData))
	return nil
}

func unredact(cmd *cobra.Command) error {
	token := os.Getenv("PANGEA_REDACT_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	keyID := os.Getenv("PANGEA_VAULT_FPE_KEY_ID")
	if keyID == "" {
		log.Fatal("Unauthorized: No key id present")
	}

	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
	)
	if err != nil {
		log.Fatalf("expected no error got: %v", err)
	}
	redactcli := redact.New(config)

	var text = "Visit our web is https://pangea.cloud"

	fmt.Printf("Redacting PII from: '%s'\n", text)

	ctx := context.Background()
	input := &redact.TextRequest{
		Text: pangea.String(text),
		VaultParameters: &redact.VaultParameters{
			FPEkeyID: keyID,
		},
	}

	redactResponse, err := redactcli.Redact(ctx, input)
	if err != nil {
		log.Fatal(err)
	}
	if redactResponse.Result.FPEContext == nil {
		log.Fatal("FPEContext is nil")
	}

	fmt.Printf("Redacted text: %s", pangea.Stringify(redactResponse.Result.RedactedText))

	unredactResponse, err := redactcli.Unredact(ctx, &redact.UnredactRequest{
		RedactedData: redactResponse.Result.RedactedText,
		FPEContext:   *redactResponse.Result.FPEContext,
	})

	if err != nil {
		log.Fatal(err)
	}

	data := unredactResponse.Result.Data.(string)
	fmt.Println("Unredacted text: ", data)
	return nil
}

func redactText(cmd *cobra.Command) error {
	token := os.Getenv("PANGEA_REDACT_TOKEN")
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
	redactcli := redact.New(config)

	var text = "Hello, my phone number is 123-456-7890"

	fmt.Printf("Redacting PII from: '%s'\n", text)

	ctx := context.Background()
	input := &redact.TextRequest{
		Text: pangea.String(text),
	}

	redactResponse, err := redactcli.Redact(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Redacted text: %s", pangea.Stringify(redactResponse.Result.RedactedText))
	return nil
}

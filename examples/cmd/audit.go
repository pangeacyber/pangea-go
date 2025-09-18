package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"time"

	"github.com/pangeacyber/pangea-go/examples/cmd/audit/util"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/option"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/audit"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/vault"
	"github.com/spf13/cobra"
)

func init() {
	auditCmd := &cobra.Command{
		Use:   "audit",
		Short: "Audit examples",
	}

	logBulkCmd := &cobra.Command{
		Use:  "log_bulk",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return auditLogBulk(cmd)
		},
	}
	auditCmd.AddCommand(logBulkCmd)

	logBulkAsyncCmd := &cobra.Command{
		Use:  "log_bulk_async",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return auditLogBulkAsync(cmd)
		},
	}
	auditCmd.AddCommand(logBulkAsyncCmd)

	logBulkAsyncWithVaultCmd := &cobra.Command{
		Use:  "log_bulk_async_with_vault",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return auditLogBulkAsyncWithVault(cmd)
		},
	}
	auditCmd.AddCommand(logBulkAsyncWithVaultCmd)

	logCustomSchemaCmd := &cobra.Command{
		Use:  "log_custom_schema",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return auditLogCustomSchema(cmd)
		},
	}
	auditCmd.AddCommand(logCustomSchemaCmd)

	logStandardSchemaCmd := &cobra.Command{
		Use:  "log_standard_schema",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return auditLogStandardSchema(cmd)
		},
	}
	auditCmd.AddCommand(logStandardSchemaCmd)

	resultsCustomSchemaCmd := &cobra.Command{
		Use:  "results_custom_schema",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return auditResultsCustomSchema(cmd)
		},
	}
	auditCmd.AddCommand(resultsCustomSchemaCmd)

	resultsStandardSchemaCmd := &cobra.Command{
		Use:  "results_standard_schema",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return auditResultsStandardSchema(cmd)
		},
	}
	auditCmd.AddCommand(resultsStandardSchemaCmd)

	rootCmd := &cobra.Command{
		Use:  "root",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return auditRoot(cmd)
		},
	}
	auditCmd.AddCommand(rootCmd)

	searchCustomSchemaCmd := &cobra.Command{
		Use:  "search_custom_schema",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return auditSearchCustomSchema(cmd)
		},
	}
	auditCmd.AddCommand(searchCustomSchemaCmd)

	searchStandardSchemaCmd := &cobra.Command{
		Use:  "search_standard_schema",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return auditSearchStandardSchema(cmd)
		},
	}
	auditCmd.AddCommand(searchStandardSchemaCmd)

	auditMulticonfigCmd := &cobra.Command{
		Use:  "audit_multiconfig",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return auditMulticonfig(cmd)
		},
	}
	auditCmd.AddCommand(auditMulticonfigCmd)

	ExamplesCmd.AddCommand(auditCmd)
}

func auditResultsCustomSchema(cmd *cobra.Command) error {
	token := os.Getenv("PANGEA_AUDIT_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
		option.WithPollResultTimeout(60*time.Second),
		option.WithQueuedRetryEnabled(true),
	)
	if err != nil {
		log.Fatal("failed to create config")
	}

	auditcli, err := audit.New(
		config,
		audit.WithCustomSchema(util.CustomSchemaEvent{}),
	)
	if err != nil {
		log.Fatal("failed to create audit client")
	}

	ctx := context.Background()
	si := &audit.SearchInput{
		Query:   `message:""`,
		Limit:   10,
		Verbose: pangea.Bool(false),
	}

	sr, err := auditcli.Search(ctx, si)
	if err != nil {
		log.Fatal(err)
	}

	ri := &audit.SearchResultsInput{
		ID:    sr.Result.ID,
		Limit: 3,
	}

	rr, err := auditcli.SearchResults(ctx, ri)
	if err != nil {
		log.Fatal(err)
	}

	for i, se := range rr.Result.Events {
		ec := (se.EventEnvelope.Event).(*util.CustomSchemaEvent)
		fmt.Printf("Event %d:\n %s\n", i, pangea.Stringify(*ec))
	}
	return nil
}

func auditResultsStandardSchema(cmd *cobra.Command) error {
	token := os.Getenv("PANGEA_AUDIT_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
		option.WithPollResultTimeout(60*time.Second),
		option.WithQueuedRetryEnabled(true),
	)
	if err != nil {
		log.Fatal("failed to create config")
	}

	auditcli, err := audit.New(config)
	if err != nil {
		log.Fatal("failed to create audit client")
	}

	ctx := context.Background()
	si := &audit.SearchInput{
		Query:   `message:""`,
		Limit:   10,
		Verbose: pangea.Bool(false),
	}

	sr, err := auditcli.Search(ctx, si)
	if err != nil {
		log.Fatal(err)
	}

	ri := &audit.SearchResultsInput{
		ID:    sr.Result.ID,
		Limit: 2,
	}

	rr, err := auditcli.SearchResults(ctx, ri)
	if err != nil {
		log.Fatal(err)
	}

	for i, se := range rr.Result.Events {
		ec := (se.EventEnvelope.Event).(*audit.StandardEvent)
		fmt.Printf("Event %d:\n %s\n", i, pangea.Stringify(*ec))
	}
	return nil
}

func auditRoot(cmd *cobra.Command) error {
	token := os.Getenv("PANGEA_AUDIT_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
	)
	if err != nil {
		log.Fatal("failed to create config")
	}

	auditcli, err := audit.New(config)
	if err != nil {
		log.Fatal("failed to create audit client")
	}

	ctx := context.Background()
	input := &audit.RootInput{
		TreeSize: 10,
	}

	rootResponse, err := auditcli.Root(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Result: ", pangea.Stringify(rootResponse.Result))
	return nil
}

func auditSearchCustomSchema(cmd *cobra.Command) error {
	token := os.Getenv("PANGEA_AUDIT_CUSTOM_SCHEMA_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
		option.WithPollResultTimeout(60*time.Second),
		option.WithQueuedRetryEnabled(true),
	)
	if err != nil {
		log.Fatal("failed to create config")
	}

	auditcli, err := audit.New(
		config,
		audit.WithCustomSchema(util.CustomSchemaEvent{}),
	)
	if err != nil {
		log.Fatal("failed to create audit client")
	}

	ctx := context.Background()
	input := &audit.SearchInput{
		Query:   `message:""`,
		Limit:   3,
		Verbose: pangea.Bool(false),
	}

	sr, err := auditcli.Search(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	for i, se := range sr.Result.Events {
		ec := (se.EventEnvelope.Event).(*util.CustomSchemaEvent)
		fmt.Printf("Event %d:\n %s\n", i, pangea.Stringify(*ec))
	}
	return nil
}

func auditSearchStandardSchema(cmd *cobra.Command) error {
	token := os.Getenv("PANGEA_AUDIT_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
		option.WithPollResultTimeout(60*time.Second),
		option.WithQueuedRetryEnabled(true),
	)
	if err != nil {
		log.Fatal("failed to create config")
	}

	auditcli, err := audit.New(config)
	if err != nil {
		log.Fatal("failed to create audit client")
	}

	ctx := context.Background()
	input := &audit.SearchInput{
		Query:   `message:""`,
		Limit:   3,
		Verbose: pangea.Bool(false),
	}

	sr, err := auditcli.Search(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	for i, se := range sr.Result.Events {
		ec := (se.EventEnvelope.Event).(*audit.StandardEvent)
		fmt.Printf("Event %d:\n %s\n", i, pangea.Stringify(*ec))
	}
	return nil
}

func auditMulticonfig(cmd *cobra.Command) error {
	token := os.Getenv("PANGEA_AUDIT_MULTICONFIG_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	configId := os.Getenv("PANGEA_AUDIT_CONFIG_ID")
	if configId == "" {
		log.Fatal("Need to set PANGEA_AUDIT_CONFIG_ID env var")
	}

	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
	)
	if err != nil {
		log.Fatal("failed to create config")
	}

	// Set configId in service construction
	auditcli, err := audit.New(config, audit.WithConfigID(configId))
	if err != nil {
		log.Fatal("failed to create audit client")
	}

	ctx := context.Background()
	event := &audit.StandardEvent{
		Message: "Hello, World!",
	}

	fmt.Printf("Logging: %s\n", event.Message)

	lr, err := auditcli.Log(ctx, event, true)
	if err != nil {
		log.Fatal(err)
	}

	e := (lr.Result.EventEnvelope.Event).(*audit.StandardEvent)
	fmt.Printf("Logged event: %s", pangea.Stringify(e))
	return nil
}

func auditLogStandardSchema(cmd *cobra.Command) error {
	token := os.Getenv("PANGEA_AUDIT_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
	)
	if err != nil {
		log.Fatal("failed to create config")
	}

	auditcli, err := audit.New(config)
	if err != nil {
		log.Fatal("failed to create audit client")
	}

	ctx := context.Background()
	event := &audit.StandardEvent{
		Message: "Hello, World!",
	}

	fmt.Printf("Logging: %s\n", event.Message)

	lr, err := auditcli.Log(ctx, event, true)
	if err != nil {
		log.Fatal(err)
	}

	e := (lr.Result.EventEnvelope.Event).(*audit.StandardEvent)
	fmt.Printf("Logged event: %s", pangea.Stringify(e))
	return nil
}

func auditLogCustomSchema(cmd *cobra.Command) error {
	token := os.Getenv("PANGEA_AUDIT_CUSTOM_SCHEMA_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
	)
	if err != nil {
		log.Fatal("failed to create config")
	}

	auditcli, err := audit.New(
		config,
		audit.WithCustomSchema(util.CustomSchemaEvent{}),
	)
	if err != nil {
		log.Fatal("failed to create audit client")
	}

	ctx := context.Background()
	msg := "go-sdk-custom-schema-no-signed"
	var event = &util.CustomSchemaEvent{
		Message:       msg,
		FieldInt:      1,
		FieldBool:     true,
		FieldStrShort: "no-signed",
		FieldStrLong:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed lacinia, orci eget commodo commodo non.",
		FieldTime:     pangea.Time(time.Now().Truncate(time.Microsecond)),
	}

	fmt.Printf("Logging: %s\n", event.Message)

	lr, err := auditcli.Log(ctx, event, true)
	if err != nil {
		log.Fatal(err)
	}

	e := (lr.Result.EventEnvelope.Event).(*util.CustomSchemaEvent)
	fmt.Printf("Logged event: %s", pangea.Stringify(e))
	return nil
}

func auditLogBulkAsyncWithVault(cmd *cobra.Command) error {
	vaultToken := os.Getenv("PANGEA_VAULT_TOKEN")
	domain := os.Getenv("PANGEA_DOMAIN")
	auditTokenVaultId := os.Getenv("PANGEA_AUDIT_TOKEN_VAULT_ID")
	if vaultToken == "" {
		log.Fatal("Unauthorized: No vault token present")
	}

	config, err := pangea.NewConfig(
		option.WithToken(vaultToken),
		option.WithDomain(domain),
	)
	if err != nil {
		log.Fatal("failed to create config")
	}

	vaultcli := vault.New(config)
	getInput := &vault.GetRequest{
		ID: auditTokenVaultId,
	}
	ctx := context.Background()
	getResponse, err := vaultcli.Get(ctx, getInput)
	if err != nil {
		log.Fatal(err)
	}
	auditToken := getResponse.Result.ItemVersions[0].Token
	if auditToken == nil {
		log.Fatal("Unexpected nil auditToken")
	}

	config, err = pangea.NewConfig(option.WithToken(*auditToken), option.WithDomain(domain))
	if err != nil {
		log.Fatalf("expected no error got: %v", err)
	}
	auditcli, err := audit.New(config)
	if err != nil {
		log.Fatal("failed to create audit client")
	}

	event1 := &audit.StandardEvent{
		Message: "Sign up",
		Actor:   "go-sdk",
	}

	resp, err := auditcli.LogBulkAsync(ctx, []any{event1}, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Success. Request_id: %s\n", *resp.RequestID)
	return nil
}

func auditLogBulkAsync(cmd *cobra.Command) error {
	token := os.Getenv("PANGEA_AUDIT_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
	)
	if err != nil {
		log.Fatal("failed to create config")
	}

	auditcli, err := audit.New(config)
	if err != nil {
		log.Fatal("failed to create audit client")
	}

	ctx := context.Background()
	event1 := &audit.StandardEvent{
		Message: "Sign up",
		Actor:   "go-sdk",
	}

	event2 := &audit.StandardEvent{
		Message: "Sign in",
		Actor:   "go-sdk",
	}

	fmt.Println("Sending multiple events...")

	resp, err := auditcli.LogBulkAsync(ctx, []any{event1, event2}, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Success. Request_id: %s\n", *resp.RequestID)
	return nil
}

func auditLogBulk(cmd *cobra.Command) error {
	token := os.Getenv("PANGEA_AUDIT_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
	)
	if err != nil {
		log.Fatal("failed to create config")
	}

	auditcli, err := audit.New(config)
	if err != nil {
		log.Fatal("failed to create audit client")
	}

	ctx := context.Background()
	event1 := &audit.StandardEvent{
		Message: "Sign up",
		Actor:   "go-sdk",
	}

	event2 := &audit.StandardEvent{
		Message: "Sign in",
		Actor:   "go-sdk",
	}

	fmt.Println("Logging multiple events...")

	lr, err := auditcli.LogBulk(ctx, []any{event1, event2}, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Results:")
	for _, r := range lr.Result.Results {
		e := (r.EventEnvelope.Event).(*audit.StandardEvent)
		fmt.Printf("\tLogged event: %s", pangea.Stringify(e))
	}
	return nil
}

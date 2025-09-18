package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/option"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/user_intel"
	"github.com/spf13/cobra"
)

func init() {
	userCmd := &cobra.Command{
		Use:   "user",
		Short: "User intel examples",
	}
	intelCmd.AddCommand(userCmd)

	userBreachedByEmailCmd := &cobra.Command{
		Use:  "user_breached_by_email",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return userBreachedByEmail(cmd)
		},
	}
	userCmd.AddCommand(userBreachedByEmailCmd)

	passwordBreachedCmd := &cobra.Command{
		Use:  "password_breached",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return passwordBreached(cmd)
		},
	}
	userCmd.AddCommand(passwordBreachedCmd)

	userBreachedByIPCmd := &cobra.Command{
		Use:  "user_breached_by_ip",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return userBreachedByIP(cmd)
		},
	}
	userCmd.AddCommand(userBreachedByIPCmd)

	userBreachedByPhoneCmd := &cobra.Command{
		Use:  "user_breached_by_phone",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return userBreachedByPhone(cmd)
		},
	}
	userCmd.AddCommand(userBreachedByPhoneCmd)

	userBreachedByUsernameCmd := &cobra.Command{
		Use:  "user_breached_by_username",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return userBreachedByUsername(cmd)
		},
	}
	userCmd.AddCommand(userBreachedByUsernameCmd)
}

func userBreachedByEmail(cmd *cobra.Command) error {
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
	intelcli := user_intel.New(config)

	ctx := context.Background()
	input := &user_intel.UserBreachedRequest{
		Email:    "test@example.com",
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "spycloud",
	}

	resp, err := intelcli.UserBreached(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(pangea.Stringify(resp.Result.Data))
	return nil
}

func passwordBreached(cmd *cobra.Command) error {
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
	intelcli := user_intel.New(config)

	// Set the password you would like to check
	// Observe proper safety with passwords, do not check them into source control etc.
	password := "mypassword"
	// Calculate its hash, it could be sha256, sha512 or sha1
	hash := pangea.HashSHA256(password)
	// get the hash prefix, just the first 5 characters
	hashPrefix := pangea.GetHashPrefix(hash, 5)

	ctx := context.Background()
	input := &user_intel.UserPasswordBreachedRequest{
		// set the right hash_type here, sha256, sha512 or sha1
		HashType:   user_intel.HTsha265,
		HashPrefix: hashPrefix,
		Raw:        pangea.Bool(true),
		Verbose:    pangea.Bool(true),
		Provider:   "spycloud",
	}

	r, err := intelcli.PasswordBreached(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	// IsPasswordBreached is a helper function that can simplify searching the response's raw data for the full hash
	s, err := user_intel.IsPasswordBreached(r, hash)
	if err != nil {
		log.Fatal(err)
	}

	if s == user_intel.PSbreached {
		fmt.Printf("Password '%s' has been breached.\n", password)
	} else if s == user_intel.PSunbreached {
		fmt.Printf("Password '%s' has not been breached.\n", password)
	} else if s == user_intel.PSinconclusive {
		fmt.Printf("Not enough information to confirm if password '%s' has been or has not been breached.\n", password)
	} else {
		fmt.Printf("Unknown status: %d.\n", s)
	}
	return nil
}

func userBreachedByIP(cmd *cobra.Command) error {
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
	intelcli := user_intel.New(config)

	ctx := context.Background()
	input := &user_intel.UserBreachedRequest{
		IP:       "192.168.140.37",
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "spycloud",
	}

	resp, err := intelcli.UserBreached(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(pangea.Stringify(resp.Result.Data))
	return nil
}

func userBreachedByPhone(cmd *cobra.Command) error {
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
	intelcli := user_intel.New(config)

	ctx := context.Background()
	input := &user_intel.UserBreachedRequest{
		PhoneNumber: "8005550123",
		Raw:         pangea.Bool(true),
		Verbose:     pangea.Bool(true),
		Provider:    "spycloud",
	}

	resp, err := intelcli.UserBreached(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(pangea.Stringify(resp.Result.Data))
	return nil
}

func userBreachedByUsername(cmd *cobra.Command) error {
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
	intelcli := user_intel.New(config)

	ctx := context.Background()
	input := &user_intel.UserBreachedRequest{
		Username: "shortpatrick",
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "spycloud",
	}

	resp, err := intelcli.UserBreached(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(pangea.Stringify(resp.Result.Data))
	return nil
}

package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/option"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/ip_intel"
	"github.com/spf13/cobra"
)

func init() {
	ipCmd := &cobra.Command{
		Use:   "ip",
		Short: "IP intel examples",
	}

	intelCmd.AddCommand(ipCmd)

	ipReputationCymruCmd := &cobra.Command{
		Use:  "reputation_cymru",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return ipReputationCymru(cmd)
		},
	}
	ipCmd.AddCommand(ipReputationCymruCmd)

	ipDomainBulkCmd := &cobra.Command{
		Use:  "domain_bulk",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return ipDomainBulk(cmd)
		},
	}
	ipCmd.AddCommand(ipDomainBulkCmd)

	ipDomainCmd := &cobra.Command{
		Use:  "domain",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return ipDomain(cmd)
		},
	}
	ipCmd.AddCommand(ipDomainCmd)

	ipGeolocateBulkCmd := &cobra.Command{
		Use:  "geolocate_bulk",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return ipGeolocateBulk(cmd)
		},
	}
	ipCmd.AddCommand(ipGeolocateBulkCmd)

	ipGeolocateCmd := &cobra.Command{
		Use:  "geolocate",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return ipGeolocate(cmd)
		},
	}
	ipCmd.AddCommand(ipGeolocateCmd)

	ipProxyBulkCmd := &cobra.Command{
		Use:  "proxy_bulk",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return ipProxyBulk(cmd)
		},
	}
	ipCmd.AddCommand(ipProxyBulkCmd)

	ipProxyCmd := &cobra.Command{
		Use:  "proxy",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return ipProxy(cmd)
		},
	}
	ipCmd.AddCommand(ipProxyCmd)

	ipReputationBulkCmd := &cobra.Command{
		Use:  "reputation_bulk",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return ipReputationBulk(cmd)
		},
	}
	ipCmd.AddCommand(ipReputationBulkCmd)

	ipReputationCrowdstrikeCmd := &cobra.Command{
		Use:  "reputation_crowdstrike",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return ipReputationCrowdstrike(cmd)
		},
	}
	ipCmd.AddCommand(ipReputationCrowdstrikeCmd)

	ipVPNBulkCmd := &cobra.Command{
		Use:  "vpn_bulk",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return ipVPNBulk(cmd)
		},
	}
	ipCmd.AddCommand(ipVPNBulkCmd)

	ipVPNCmd := &cobra.Command{
		Use:  "vpn",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return ipVPN(cmd)
		},
	}
	ipCmd.AddCommand(ipVPNCmd)
}

func PrintIPData(indicator string, data ip_intel.ReputationData) {
	fmt.Printf("\t Indicator: %s\n", indicator)
	fmt.Printf("\t\t Verdict: %s\n", data.Verdict)
	fmt.Printf("\t\t Score: %d\n", data.Score)
	fmt.Printf("\t\t Category: %s\n", pangea.Stringify(data.Category))
}

func ipReputationCymru(cmd *cobra.Command) error {
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
	intelcli := ip_intel.New(config)

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
	PrintIPData(ip, resp.Result.Data)
	return nil
}

func ipDomainBulk(cmd *cobra.Command) error {
	fmt.Println("Checking IP's domain...")
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
	intelcli := ip_intel.New(config)

	ctx := context.Background()
	input := &ip_intel.IpDomainBulkRequest{
		Ips:     []string{"93.231.182.110", "190.28.74.251"},
		Raw:     pangea.Bool(true),
		Verbose: pangea.Bool(true),
	}

	resp, err := intelcli.GetDomainBulk(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Result:")
	for k, v := range resp.Result.Data {
		if v.DomainFound {
			fmt.Printf("\t IP %s domain is %s\n", k, v.Domain)
		} else {
			fmt.Printf("\t IP %s domain not found\n", k)
		}
	}
	return nil
}

func ipDomain(cmd *cobra.Command) error {
	fmt.Println("Checking IP's domain...")
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
	intelcli := ip_intel.New(config)

	ctx := context.Background()
	ip := "24.235.114.61"
	input := &ip_intel.IpDomainRequest{
		Ip:       ip,
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "digitalelement",
	}

	resp, err := intelcli.GetDomain(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Result:")
	if resp.Result.Data.DomainFound {
		fmt.Printf("\t IP %s domain is %s\n", ip, resp.Result.Data.Domain)
	} else {
		fmt.Printf("\t IP %s domain not found\n", ip)
	}
	return nil
}

func ipGeolocateBulk(cmd *cobra.Command) error {
	fmt.Println("Geolocating IP...")
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
	intelcli := ip_intel.New(config)

	ctx := context.Background()
	input := &ip_intel.IpGeolocateBulkRequest{
		Ips:     []string{"93.231.182.110", "190.28.74.251"},
		Raw:     pangea.Bool(true),
		Verbose: pangea.Bool(true),
	}

	resp, err := intelcli.GeolocateBulk(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Result:")
	for k, v := range resp.Result.Data {
		fmt.Printf("\t Indicator: %s\n", k)
		fmt.Printf("\t\t Country: %s\n", v.Country)
		fmt.Printf("\t\t City: %s\n", v.City)
		fmt.Printf("\t\t Latitude: %f\n", v.Latitude)
		fmt.Printf("\t\t Longitude: %f\n", v.Longitude)
		fmt.Printf("\t\t PostalCode: %s\n", v.PostalCode)
		fmt.Printf("\t\t CountryCode: %s\n", v.CountryCode)
	}
	return nil
}

func ipGeolocate(cmd *cobra.Command) error {
	fmt.Println("Geolocating IP...")
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
	intelcli := ip_intel.New(config)

	ctx := context.Background()
	ip := "93.231.182.110"
	input := &ip_intel.IpGeolocateRequest{
		Ip:       ip,
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "digitalelement",
	}

	resp, err := intelcli.Geolocate(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Result:")
	fmt.Printf("\t Indicator: %s\n", ip)
	fmt.Printf("\t\t Country: %s\n", resp.Result.Data.Country)
	fmt.Printf("\t\t City: %s\n", resp.Result.Data.City)
	fmt.Printf("\t\t Latitude: %f\n", resp.Result.Data.Latitude)
	fmt.Printf("\t\t Longitude: %f\n", resp.Result.Data.Longitude)
	fmt.Printf("\t\t PostalCode: %s\n", resp.Result.Data.PostalCode)
	fmt.Printf("\t\t CountryCode: %s\n", resp.Result.Data.CountryCode)
	return nil
}

func ipProxyBulk(cmd *cobra.Command) error {
	fmt.Println("Checking if IP is a proxy...")
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
	intelcli := ip_intel.New(config)

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
	for k, v := range resp.Result.Data {
		if v.IsProxy {
			fmt.Printf("\t IP %s is a proxy\n", k)
		} else {
			fmt.Printf("\t IP %s is not a proxy\n", k)
		}
	}
	return nil
}

func ipProxy(cmd *cobra.Command) error {
	fmt.Println("Checking if IP is a proxy...")
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
	intelcli := ip_intel.New(config)

	ctx := context.Background()
	ip := "34.201.32.172"
	input := &ip_intel.IpProxyRequest{
		Ip:       ip,
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "digitalelement",
	}

	resp, err := intelcli.IsProxy(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Result:")
	if resp.Result.Data.IsProxy {
		fmt.Printf("\t IP %s is a proxy\n", ip)
	} else {
		fmt.Printf("\t IP %s is not a proxy\n", ip)
	}
	return nil
}

func ipReputationBulk(cmd *cobra.Command) error {
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
	intelcli := ip_intel.New(config)

	ctx := context.Background()
	input := &ip_intel.IpReputationBulkRequest{
		Ips:      []string{"93.231.182.110", "190.28.74.251"},
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "crowdstrike",
	}

	resp, err := intelcli.ReputationBulk(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Result:")
	for k, v := range resp.Result.Data {
		PrintIPData(k, v)
	}
	return nil
}

func ipReputationCrowdstrike(cmd *cobra.Command) error {
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
	intelcli := ip_intel.New(config)

	ctx := context.Background()
	ip := "93.231.182.110"
	input := &ip_intel.IpReputationRequest{
		Ip:       ip,
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "crowdstrike",
	}

	resp, err := intelcli.Reputation(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Result:")
	PrintIPData(ip, resp.Result.Data)
	return nil
}

func ipVPNBulk(cmd *cobra.Command) error {
	fmt.Println("Checking if IP is a vpn...")
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
	intelcli := ip_intel.New(config)

	ctx := context.Background()
	input := &ip_intel.IpVPNBulkRequest{
		Ips:     []string{"2.56.189.74", "190.28.74.251"},
		Raw:     pangea.Bool(true),
		Verbose: pangea.Bool(true),
	}

	resp, err := intelcli.IsVPNBulk(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Result:")
	for k, v := range resp.Result.Data {
		if v.IsVPN {
			fmt.Printf("\t IP %s is a VPN\n", k)
		} else {
			fmt.Printf("\t IP %s is not a VPN\n", k)
		}
	}
	return nil
}

func ipVPN(cmd *cobra.Command) error {
	fmt.Println("Checking if IP is a vpn...")
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
	intelcli := ip_intel.New(config)

	ctx := context.Background()
	ip := "2.56.189.74"
	input := &ip_intel.IpVPNRequest{
		Ip:       ip,
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "digitalelement",
	}

	resp, err := intelcli.IsVPN(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Result:")
	if resp.Result.Data.IsVPN {
		fmt.Printf("\t IP %s is a VPN\n", ip)
	} else {
		fmt.Printf("\t IP %s is not a VPN\n", ip)
	}
	return nil
}

package cmd

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/option"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/domain_intel"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/url_intel"
	"github.com/spf13/cobra"
)

func init() {
	defangCmd := &cobra.Command{
		Use:   "defang",
		Short: "Defang a URL",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return defang(cmd, args[0])
		},
	}
	intelCmd.AddCommand(defangCmd)
}

var defangedSchemes = map[string]string{
	"http":  "hxxp",
	"https": "hxxps",
}

func DefangURL(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	ds, ok := defangedSchemes[u.Scheme]
	if !ok {
		ds = "xxxx"
	}
	return strings.Replace(s, u.Scheme, ds, 1), nil
}

func GetDomain(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	p := strings.Split(u.Host, ":")
	return p[0], nil
}

func defang(cmd *cobra.Command, url string) error {
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
	urlc := url_intel.New(config)

	ctx := context.Background()
	urlReq := &url_intel.UrlReputationRequest{
		Url:      url,
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "crowdstrike",
	}

	urlResp, err := urlc.Reputation(ctx, urlReq)
	if err != nil {
		log.Fatal(err)
	}

	if urlResp.Result.Data.Verdict == "malicious" {
		du, err := DefangURL(url)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(fmt.Sprintf("Defanged URL: %s", du))
		return nil
	}

	domain, err := GetDomain(url)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(domain)

	config, err = pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
	)
	if err != nil {
		log.Fatalf("expected no error got: %v", err)
	}
	domainc := domain_intel.New(config)

	ctx = context.Background()
	domainReq := &domain_intel.DomainReputationRequest{
		Domain:   domain,
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "domaintools",
	}

	domainResp, err := domainc.Reputation(ctx, domainReq)
	if err != nil {
		log.Fatal(err)
	}

	if domainResp.Result.Data.Verdict == "malicious" {
		du, err := DefangURL(url)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(fmt.Sprintf("Defanged URL: %s", du))
		return nil
	}

	fmt.Println(fmt.Sprintf("URL: %s seems to be secure", url))
	return nil
}

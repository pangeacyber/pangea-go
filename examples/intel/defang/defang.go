// Domain amd URL intel can be used to defang malicious URLs
package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/domain_intel"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/url_intel"
)

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

func main() {
	token := os.Getenv("PANGEA_INTEL_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	urlc := url_intel.New(&pangea.Config{
		Token:  token,
		Domain: os.Getenv("PANGEA_DOMAIN"),
	})

	url := "http://113.235.101.11:54384"

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
		return
	}

	domain, err := GetDomain(url)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(domain)
	domainc := domain_intel.New(&pangea.Config{
		Token:  token,
		Domain: os.Getenv("PANGEA_DOMAIN"),
	})

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
		return
	}

	fmt.Println(fmt.Sprintf("URL: %s seems to be secure", url))
}

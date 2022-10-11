// go:build integration
package url_intel_test

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	"github.com/pangeacyber/go-pangea/internal/pangeatesting"
// 	"github.com/pangeacyber/go-pangea/pangea"
// 	"github.com/pangeacyber/go-pangea/service/url_intel"
// 	"github.com/stretchr/testify/assert"
// )

// func Test_Integration_UrlLookup(t *testing.T) {
// 	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancelFn()

// 	cfgToken := pangeatesting.GetEnvVarOrSkip(t, "URL_INTEL_INTEGRATION_CONFIG_TOKEN")
// 	cfg := &pangea.Config{
// 		ConfigID: cfgToken,
// 	}
// 	cfg = cfg.Copy(pangeatesting.IntegrationConfig(t))
// 	urlintel := url_intel.New(cfg)

// 	input := &url_intel.UrlLookupInput{
// 		Url:      "http://113.235.101.11:54384",
// 		Raw:      true,
// 		Verbose:  true,
// 		Provider: "crowdstrike",
// 	}
// 	out, err := urlintel.Lookup(ctx, input)
// 	if err != nil {
// 		t.Fatalf("expected no error got: %v", err)
// 	}

// 	assert.NotNil(t, out)
// 	assert.NotNil(t, out.Result)
// 	assert.NotNil(t, out.Result.Data)
// 	assert.Equal(t, out.Result.Data.Verdict, "malicious")
// }

// func Test_Integration_UrlLookup_2(t *testing.T) {
// 	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancelFn()

// 	cfgToken := pangeatesting.GetEnvVarOrSkip(t, "URL_INTEL_INTEGRATION_CONFIG_TOKEN")
// 	cfg := &pangea.Config{
// 		ConfigID: cfgToken,
// 	}
// 	cfg = cfg.Copy(pangeatesting.IntegrationConfig(t))
// 	urlintel := url_intel.New(cfg)

// 	input := &url_intel.UrlLookupInput{
// 		Url:      "http://google.com",
// 		Raw:      true,
// 		Verbose:  true,
// 		Provider: "crowdstrike",
// 	}
// 	out, err := urlintel.Lookup(ctx, input)
// 	if err != nil {
// 		t.Fatalf("expected no error got: %v", err)
// 	}

// 	assert.NotNil(t, out)
// 	assert.NotNil(t, out.Result)
// 	assert.NotNil(t, out.Result.Data)
// 	assert.Equal(t, out.Result.Data.Verdict, "unknown")
// }

// func Test_Integration_UrlLookup_Error_BadToken(t *testing.T) {
// 	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancelFn()

// 	cfgToken := pangeatesting.GetEnvVarOrSkip(t, "URL_INTEL_INTEGRATION_CONFIG_TOKEN")
// 	cfg := &pangea.Config{
// 		ConfigID: cfgToken,
// 	}
// 	cfg = cfg.Copy(pangeatesting.IntegrationConfig(t))
// 	cfg.Token = "notarealtoken"
// 	urlintel := url_intel.New(cfg)

// 	input := &url_intel.UrlLookupInput{
// 		Url:      "http://113.235.101.11:54384",
// 		Raw:      true,
// 		Verbose:  true,
// 		Provider: "crowdstrike",
// 	}

// 	out, err := urlintel.Lookup(ctx, input)

// 	assert.Error(t, err)
// 	assert.Nil(t, out)
// 	apiErr := err.(*pangea.APIError)
// 	assert.Equal(t, apiErr.Err.Error(), "API error: Not authorized to access this resource.")
// }

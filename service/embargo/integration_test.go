// go:build integration
package embargo_test

import (
	"context"
	"testing"
	"time"

	"github.com/pangeacyber/go-pangea/internal/pangeatesting"
	"github.com/pangeacyber/go-pangea/pangea"
	"github.com/pangeacyber/go-pangea/service/embargo"
	"github.com/stretchr/testify/assert"
)

// Check ISO with sanctions
func Test_Integration_Check(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfgToken := pangeatesting.GetEnvVarOrSkip(t, "EMBARGO_INTEGRATION_CONFIG_TOKEN")
	cfg := &pangea.Config{
		ConfigID: cfgToken,
	}
	cfg = cfg.Copy(pangeatesting.IntegrationConfig(t))
	client := embargo.New(cfg)

	input := &embargo.ISOCheckInput{
		ISOCode: pangea.String("CU"),
	}
	out, err := client.ISOCheck(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, out.Result)
	assert.NotZero(t, out.Result.Count)
	assert.GreaterOrEqual(t, len(out.Result.Sanctions), 1)
	assert.Equal(t, out.Result.Sanctions[0].EmbargoedCountryName, "Cuba")
}

// Check ISO without sanctions
func Test_Integration_Check_2(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfgToken := pangeatesting.GetEnvVarOrSkip(t, "EMBARGO_INTEGRATION_CONFIG_TOKEN")
	cfg := &pangea.Config{
		ConfigID: cfgToken,
	}
	cfg = cfg.Copy(pangeatesting.IntegrationConfig(t))
	client := embargo.New(cfg)

	input := &embargo.ISOCheckInput{
		ISOCode: pangea.String("AR"),
	}
	out, err := client.ISOCheck(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, out.Result)
	assert.Zero(t, out.Result.Count)
}

func Test_Integration_Check_Error_BadISO(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfgToken := pangeatesting.GetEnvVarOrSkip(t, "EMBARGO_INTEGRATION_CONFIG_TOKEN")
	cfg := &pangea.Config{
		ConfigID: cfgToken,
	}
	cfg = cfg.Copy(pangeatesting.IntegrationConfig(t))
	client := embargo.New(cfg)

	input := &embargo.ISOCheckInput{
		ISOCode: pangea.String("NotAnISOcode"),
	}
	out, err := client.ISOCheck(ctx, input)

	assert.NotNil(t, err)
	assert.Nil(t, out)
	apiErr := err.(*pangea.APIError)
	assert.Equal(t, len(apiErr.PangeaErrors.Errors), 1)
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Code, "DoesNotMatchPattern")
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Detail, "'iso_code' must match the given regex: ^[a-zA-Z]{2}$")
	assert.Equal(t, apiErr.PangeaErrors.Errors[0].Source, "/iso_code")
}

func Test_Integration_Check_Error_BadToken(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfgToken := pangeatesting.GetEnvVarOrSkip(t, "EMBARGO_INTEGRATION_CONFIG_TOKEN")
	cfg := &pangea.Config{
		ConfigID: cfgToken,
	}
	cfg = cfg.Copy(pangeatesting.IntegrationConfig(t))
	cfg.Token = "notatoken"
	client := embargo.New(cfg)

	input := &embargo.ISOCheckInput{
		ISOCode: pangea.String("CU"),
	}
	out, err := client.ISOCheck(ctx, input)

	assert.NotNil(t, err)
	assert.Nil(t, out)
	apiErr := err.(*pangea.APIError)
	assert.Equal(t, apiErr.Err.Error(), "API error: Not authorized to access this resource.")

}

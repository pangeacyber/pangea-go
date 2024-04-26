//go:build integration

package embargo_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/internal/pangeatesting"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/service/embargo"
	"github.com/stretchr/testify/assert"
)

var testingEnvironment = pangeatesting.LoadTestEnvironment("embargo", pangeatesting.Live)

func embargoIntegrationCfg(t *testing.T) *pangea.Config {
	t.Helper()
	return pangeatesting.IntegrationConfig(t, testingEnvironment)
}

// Check ISO with sanctions
func Test_Integration_Check(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := embargoIntegrationCfg(t)
	client := embargo.New(cfg)

	input := &embargo.ISOCheckRequest{
		ISOCode: "CU",
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

	cfg := embargoIntegrationCfg(t)
	client := embargo.New(cfg)

	input := &embargo.ISOCheckRequest{
		ISOCode: "AR",
	}
	out, err := client.ISOCheck(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, out.Result)
	assert.Zero(t, out.Result.Count)

	rr, err := json.Marshal(out)
	assert.NoError(t, err)
	assert.NotNil(t, rr)
	assert.True(t, len(rr) > 0)
	fmt.Println("Marshalled response:")
	fmt.Println(string(rr))
}

func Test_Integration_Check_Error_BadISO(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := embargoIntegrationCfg(t)
	client := embargo.New(cfg)

	input := &embargo.ISOCheckRequest{
		ISOCode: "NotAnISOcode",
	}
	out, err := client.ISOCheck(ctx, input)

	assert.NotNil(t, err)
	assert.Nil(t, out)
	fmt.Println(err.Error())
	apiErr := err.(*pangea.APIError)
	assert.Equal(t, 1, len(apiErr.PangeaErrors.Errors))
	assert.Equal(t, "DoesNotMatchPattern", apiErr.PangeaErrors.Errors[0].Code)
	assert.Equal(t, "/iso_code", apiErr.PangeaErrors.Errors[0].Source)
}

func Test_Integration_Check_Error_BadToken(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg := embargoIntegrationCfg(t)
	cfg.Token = "notatoken"
	client := embargo.New(cfg)

	input := &embargo.ISOCheckRequest{
		ISOCode: "CU",
	}
	out, err := client.ISOCheck(ctx, input)

	assert.NotNil(t, err)
	assert.Nil(t, out)
	fmt.Println(err.Error())
	apiErr := err.(*pangea.APIError)
	assert.Equal(t, apiErr.Err.Error(), "API error: Not authorized to access this resource.")
}

// go:build integration
package user_intel_test

import (
	"context"
	"testing"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/internal/pangeatesting"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/service/user_intel"
	"github.com/stretchr/testify/assert"
)

const (
	testingEnvironment = pangeatesting.Develop
)

func Test_Integration_UserBreachedByPhone(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	intelcli := user_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &user_intel.UserBreachedRequest{
		PhoneNumber: "8005550123",
		Raw:         pangea.Bool(true),
		Verbose:     pangea.Bool(true),
		Provider:    "spycloud",
	}
	resp, err := intelcli.UserBreached(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	assert.NotNil(t, resp.Result.Data)
	assert.True(t, resp.Result.Data.FoundInBreach)
	assert.Greater(t, resp.Result.Data.BreachCount, 0)
}

func Test_Integration_UserBreachedByPhoneBulk(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	intelcli := user_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &user_intel.UserBreachedBulkRequest{
		PhoneNumbers: []string{"8005550123", "8005550124"},
		Raw:          pangea.Bool(true),
		Verbose:      pangea.Bool(true),
		Provider:     "spycloud",
	}
	resp, err := intelcli.UserBreachedBulk(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	assert.NotNil(t, resp.Result.Data)
	assert.Equal(t, len(resp.Result.Data), 2)
}

func Test_Integration_UserBreachedByEmail(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	intelcli := user_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &user_intel.UserBreachedRequest{
		Email:    "test@example.com",
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "spycloud",
	}
	resp, err := intelcli.UserBreached(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	assert.NotNil(t, resp.Result.Data)
	assert.True(t, resp.Result.Data.FoundInBreach)
	assert.Greater(t, resp.Result.Data.BreachCount, 0)
}

func Test_Integration_UserBreachedByEmailBulk(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	intelcli := user_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &user_intel.UserBreachedBulkRequest{
		Emails:   []string{"test@example.com", "noreply@example.com"},
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "spycloud",
	}
	resp, err := intelcli.UserBreachedBulk(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	assert.NotNil(t, resp.Result.Data)
	assert.Equal(t, len(resp.Result.Data), 2)
}

func Test_Integration_UserBreachedByUsername(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	intelcli := user_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &user_intel.UserBreachedRequest{
		Username: "shortpatrick",
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "spycloud",
	}
	resp, err := intelcli.UserBreached(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	assert.NotNil(t, resp.Result.Data)
	assert.True(t, resp.Result.Data.FoundInBreach)
	assert.Greater(t, resp.Result.Data.BreachCount, 0)
}

func Test_Integration_UserBreachedByUsernameBulk(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	intelcli := user_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &user_intel.UserBreachedBulkRequest{
		Usernames: []string{"shortpatrick", "user1"},
		Raw:       pangea.Bool(true),
		Verbose:   pangea.Bool(true),
		Provider:  "spycloud",
	}
	resp, err := intelcli.UserBreachedBulk(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	assert.NotNil(t, resp.Result.Data)
	assert.Equal(t, len(resp.Result.Data), 2)
}

func Test_Integration_UserBreachedByIP(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	intelcli := user_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &user_intel.UserBreachedRequest{
		IP:       "192.168.140.37",
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "spycloud",
	}
	resp, err := intelcli.UserBreached(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	assert.NotNil(t, resp.Result.Data)
	assert.True(t, resp.Result.Data.FoundInBreach)
	assert.Greater(t, resp.Result.Data.BreachCount, 0)
}

func Test_Integration_UserBreachedByIPBulk(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	intelcli := user_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &user_intel.UserBreachedBulkRequest{
		IPs:      []string{"192.168.140.37", "1.1.1.1"},
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "spycloud",
	}
	resp, err := intelcli.UserBreachedBulk(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	assert.NotNil(t, resp.Result.Data)
	assert.Equal(t, len(resp.Result.Data), 2)
}

func Test_Integration_UserBreachedDefaultProvider(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	intelcli := user_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &user_intel.UserBreachedRequest{
		PhoneNumber: "8005550123",
		Raw:         pangea.Bool(true),
		Verbose:     pangea.Bool(true),
	}
	resp, err := intelcli.UserBreached(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	assert.NotNil(t, resp.Result.Data)
	assert.True(t, resp.Result.Data.FoundInBreach)
	assert.Greater(t, resp.Result.Data.BreachCount, 0)
}

func Test_Integration_UserBreached_Error_BadToken(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	cfg := pangeatesting.IntegrationConfig(t, testingEnvironment)
	cfg.Token = "notarealtoken"
	intelcli := user_intel.New(cfg)

	input := &user_intel.UserBreachedRequest{
		PhoneNumber: "8005550123",
		Raw:         pangea.Bool(true),
		Verbose:     pangea.Bool(true),
	}
	resp, err := intelcli.UserBreached(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, resp)
	apiErr := err.(*pangea.APIError)
	assert.Equal(t, apiErr.Err.Error(), "API error: Not authorized to access this resource.")
}

func Test_Integration_PasswordBreached(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	intelcli := user_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &user_intel.UserPasswordBreachedRequest{
		HashType:   user_intel.HTsha265,
		HashPrefix: "5baa6",
		Raw:        pangea.Bool(true),
		Verbose:    pangea.Bool(true),
		Provider:   "spycloud",
	}
	resp, err := intelcli.PasswordBreached(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	assert.NotNil(t, resp.Result.Data)
	assert.True(t, resp.Result.Data.FoundInBreach)
	assert.Greater(t, resp.Result.Data.BreachCount, 0)
}

func Test_Integration_PasswordBreachedBulk(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	intelcli := user_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &user_intel.UserPasswordBreachedBulkRequest{
		HashType:     user_intel.HTsha265,
		HashPrefixes: []string{"5baa6", "5baa7"},
		Raw:          pangea.Bool(true),
		Verbose:      pangea.Bool(true),
		Provider:     "spycloud",
	}
	resp, err := intelcli.PasswordBreachedBulk(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	assert.NotNil(t, resp.Result.Data)
	assert.Equal(t, len(resp.Result.Data), 2)
}

func Test_Integration_PasswordBreachedDefaultProvider(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	intelcli := user_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))

	input := &user_intel.UserPasswordBreachedRequest{
		HashType:   user_intel.HTsha265,
		HashPrefix: "5baa6",
		Raw:        pangea.Bool(true),
		Verbose:    pangea.Bool(true),
	}
	resp, err := intelcli.PasswordBreached(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	assert.NotNil(t, resp.Result.Data)
	assert.True(t, resp.Result.Data.FoundInBreach)
	assert.Greater(t, resp.Result.Data.BreachCount, 0)
}

func Test_Integration_PasswordBreachedFullWorkflow(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()

	intelcli := user_intel.New(pangeatesting.IntegrationConfig(t, testingEnvironment))
	password := "admin123"
	h := pangea.HashSHA256(password)
	hp := pangea.GetHashPrefix(h, 5)

	input := &user_intel.UserPasswordBreachedRequest{
		HashType:   user_intel.HTsha265,
		HashPrefix: hp,
		Raw:        pangea.Bool(true),
		Verbose:    pangea.Bool(true),
	}
	resp, err := intelcli.PasswordBreached(ctx, input)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	assert.NotNil(t, resp.Result.Data)
	assert.True(t, resp.Result.Data.FoundInBreach)
	assert.Greater(t, resp.Result.Data.BreachCount, 0)

	s, err := user_intel.IsPasswordBreached(resp, h)
	assert.NoError(t, err)
	assert.Equal(t, user_intel.PSbreached, s)
}

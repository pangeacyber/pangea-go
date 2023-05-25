// go:build integration
package file_scan_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/internal/pangeatesting"
	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/service/file_scan"
	"github.com/stretchr/testify/assert"
)

const (
	testingEnvironment = pangeatesting.Develop
	EICAR              = "X5O!P%@AP[4\\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*\n"
)

func Test_Integration_FileScan(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancelFn()

	cfg := pangeatesting.IntegrationConfig(t, testingEnvironment)
	cfg.PollResultTimeout = 60 * time.Second
	client := file_scan.New(cfg)

	input := &file_scan.FileScanRequest{
		Raw:      true,
		Verbose:  true,
		Provider: "reversinglabs",
	}
	file := strings.NewReader(EICAR)

	resp, err := client.Scan(ctx, input, file)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	assert.NotNil(t, resp.Result.Data)
	assert.Equal(t, resp.Result.Data.Verdict, "malicious")
}

func Test_Integration_FileScanCanceled(t *testing.T) {
	// In this case we'll setup 7 seconds timeout to context, so, once this timeout is reached, function should return AcceptedError inmediatly
	ctx, cancelFn := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancelFn()

	cfg := pangeatesting.IntegrationConfig(t, testingEnvironment)
	cfg.PollResultTimeout = 60 * time.Second
	client := file_scan.New(cfg)

	input := &file_scan.FileScanRequest{
		Raw:      true,
		Verbose:  true,
		Provider: "reversinglabs",
	}
	file := strings.NewReader(EICAR)

	resp, err := client.Scan(ctx, input, file)
	assert.Error(t, err)
	assert.Nil(t, resp)
	_, ok := err.(*pangea.AcceptedError)
	assert.True(t, ok)
}

func Test_Integration_FileScan_NoRetry(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancelFn()

	cfg := pangeatesting.IntegrationConfig(t, testingEnvironment)
	cfg.QueuedRetryEnabled = false
	client := file_scan.New(cfg)

	input := &file_scan.FileScanRequest{
		Raw:      true,
		Verbose:  true,
		Provider: "reversinglabs",
	}
	file := strings.NewReader(EICAR)

	resp, err := client.Scan(ctx, input, file)
	assert.Error(t, err)
	assert.Nil(t, resp)
	ae := err.(*pangea.AcceptedError)

	// Wait until result should be ready
	time.Sleep(time.Duration(60 * time.Second))

	pr, err := client.PollResult(ctx, *ae)
	assert.NoError(t, err)
	assert.NotNil(t, pr)
	assert.NotNil(t, pr.Result)

	r := (*pr.Result).(*file_scan.FileScanResult)
	assert.Equal(t, r.Data.Verdict, "malicious")
}

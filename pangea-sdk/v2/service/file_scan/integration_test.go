// go:build integration
package file_scan_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/internal/pangeatesting"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/file_scan"
	"github.com/stretchr/testify/assert"
)

const (
	testingEnvironment = pangeatesting.Live
	TESTFILE_PATH      = "./testdata/testfile.pdf"
)

func Test_Integration_FileScan_crowdstrike(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancelFn()

	cfg := pangeatesting.IntegrationConfig(t, testingEnvironment)
	cfg.PollResultTimeout = 60 * time.Second
	client := file_scan.New(cfg)

	input := &file_scan.FileScanRequest{
		Raw:      true,
		Verbose:  true,
		Provider: "crowdstrike",
	}

	file, err := os.Open(TESTFILE_PATH)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	resp, err := client.Scan(ctx, input, file)
	if err != nil {
		t.Fatalf("expected no error got: %v", err.Error())
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	assert.NotNil(t, resp.Result.Data)
	assert.Equal(t, resp.Result.Data.Verdict, "benign")
}

func Test_Integration_FileScan_NoRetry_crowdstrike(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancelFn()

	cfg := pangeatesting.IntegrationConfig(t, testingEnvironment)
	cfg.QueuedRetryEnabled = false
	client := file_scan.New(cfg)

	input := &file_scan.FileScanRequest{
		Raw:      true,
		Verbose:  true,
		Provider: "crowdstrike",
	}

	file, err := os.Open(TESTFILE_PATH)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	resp, err := client.Scan(ctx, input, file)
	assert.Error(t, err)
	assert.Nil(t, resp)
	ae := err.(*pangea.AcceptedError)

	// Wait until result should be ready
	time.Sleep(time.Duration(10 * time.Second))

	pr, err := client.PollResultByError(ctx, *ae)
	assert.NoError(t, err)
	assert.NotNil(t, pr)
	assert.NotNil(t, pr.Result)

	r := (*pr.Result).(*file_scan.FileScanResult)
	assert.Equal(t, r.Data.Verdict, "benign")
}

func Test_Integration_FileScan_reversinglabs(t *testing.T) {
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

	file, err := os.Open(TESTFILE_PATH)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	resp, err := client.Scan(ctx, input, file)
	if err != nil {
		t.Fatalf("expected no error got: %v", err.Error())
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	assert.NotNil(t, resp.Result.Data)
	assert.Equal(t, resp.Result.Data.Verdict, "benign")
}

func Test_Integration_FileScan_NoRetry_reversinglabs(t *testing.T) {
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

	file, err := os.Open(TESTFILE_PATH)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}

	resp, err := client.Scan(ctx, input, file)
	assert.Error(t, err)
	assert.Nil(t, resp)
	ae := err.(*pangea.AcceptedError)

	// Wait until result should be ready
	time.Sleep(time.Duration(40 * time.Second))

	pr, err := client.PollResultByError(ctx, *ae)
	assert.NoError(t, err)
	assert.NotNil(t, pr)
	assert.NotNil(t, pr.Result)

	r := (*pr.Result).(*file_scan.FileScanResult)
	assert.Equal(t, r.Data.Verdict, "benign")
}

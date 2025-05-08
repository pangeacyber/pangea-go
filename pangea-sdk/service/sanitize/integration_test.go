//go:build integration

package sanitize_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/internal/pangeatesting"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/sanitize"
	"github.com/stretchr/testify/assert"
)

const (
	TESTFILE_PATH = "./testdata/test-sanitize.txt"
)

var testingEnvironment = pangeatesting.LoadTestEnvironment("sanitize", pangeatesting.Live)

func Test_Integration_SanitizeAndShare(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	// The Sanitize config in the regular org was obsoleted by a breaking
	// change, so the custom schema org is used instead.
	cfg := pangeatesting.IntegrationCustomSchemaConfig(t, testingEnvironment)
	cfg.PollResultTimeout = 2 * time.Minute
	client := sanitize.New(cfg)

	file, err := os.Open(TESTFILE_PATH)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resp, err := client.Sanitize(ctx, &sanitize.SanitizeRequest{
		Content: &sanitize.SanitizeContent{
			URLIntel:            pangea.Bool(true),
			URLIntelProvider:    "crowdstrike",
			DomainIntel:         pangea.Bool(true),
			DomainIntelProvider: "crowdstrike",
			Defang:              pangea.Bool(true),
			DefangThreshold:     pangea.Int(20),
			Redact:              pangea.Bool(true),
		},
		ShareOutput: &sanitize.SanitizeShareOutput{
			Enabled:      pangea.Bool(true),
			OutputFolder: "sdk_test/sanitize/",
		},
		File: &sanitize.SanitizeFile{
			ScanProvider: "crowdstrike",
		},
		TransferRequest: pangea.TransferRequest{
			TransferMethod: pangea.TMpostURL,
		},
		UploadedFileName: "uploaded_file",
	}, file)
	if err != nil {
		acceptedError, isAcceptedError := err.(*pangea.AcceptedError)
		if isAcceptedError {
			t.Logf("result of request '%s' was not ready in time", *acceptedError.RequestID)
			return
		}

		t.Fatalf("unexpected error: %v", err.Error())
	}

	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	assert.Empty(t, resp.Result.DestURL)
	assert.NotEmpty(t, resp.Result.DestShareID)
	assert.Greater(t, resp.Result.Data.Redact.RedactionCount, 0)
	assert.NotEmpty(t, resp.Result.Data.Redact.SummaryCounts)
	assert.Greater(t, resp.Result.Data.Defang.ExternalURLsCount, 0)
	assert.Greater(t, resp.Result.Data.Defang.ExternalURLsCount, 0)
	assert.Equal(t, resp.Result.Data.Defang.DefangedCount, 0)
	assert.NotEmpty(t, resp.Result.Data.Defang.DomainIntelSummary)
	assert.False(t, resp.Result.Data.MaliciousFile)
}

func Test_Integration_SanitizeNoShare(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	// The Sanitize config in the regular org was obsoleted by a breaking
	// change, so the custom schema org is used instead.
	cfg := pangeatesting.IntegrationCustomSchemaConfig(t, testingEnvironment)
	cfg.PollResultTimeout = 2 * time.Minute
	client := sanitize.New(cfg)

	file, err := os.Open(TESTFILE_PATH)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resp, err := client.Sanitize(ctx, &sanitize.SanitizeRequest{
		Content: &sanitize.SanitizeContent{
			URLIntel:            pangea.Bool(true),
			URLIntelProvider:    "crowdstrike",
			DomainIntel:         pangea.Bool(true),
			DomainIntelProvider: "crowdstrike",
			Defang:              pangea.Bool(true),
			DefangThreshold:     pangea.Int(20),
			Redact:              pangea.Bool(true),
		},
		ShareOutput: &sanitize.SanitizeShareOutput{
			Enabled: pangea.Bool(false),
		},
		File: &sanitize.SanitizeFile{
			ScanProvider: "crowdstrike",
		},
		TransferRequest: pangea.TransferRequest{
			TransferMethod: pangea.TMpostURL,
		},
		UploadedFileName: "uploaded_file",
	}, file)

	if err != nil {
		acceptedError, isAcceptedError := err.(*pangea.AcceptedError)
		if isAcceptedError {
			t.Logf("Result of request '%s' was not ready in time", *acceptedError.RequestID)
			return
		}
	}

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	assert.NotEmpty(t, resp.Result.DestURL)
	assert.Empty(t, resp.Result.DestShareID)
	assert.Greater(t, resp.Result.Data.Redact.RedactionCount, 0)
	assert.NotEmpty(t, resp.Result.Data.Redact.SummaryCounts)
	assert.Greater(t, resp.Result.Data.Defang.ExternalURLsCount, 0)
	assert.Greater(t, resp.Result.Data.Defang.ExternalURLsCount, 0)
	assert.Equal(t, resp.Result.Data.Defang.DefangedCount, 0)
	assert.NotEmpty(t, resp.Result.Data.Defang.DomainIntelSummary)
	assert.False(t, resp.Result.Data.MaliciousFile)

	af, err := client.DownloadFile(ctx, *resp.Result.DestURL)
	assert.NoError(t, err)

	err = af.Save(pangea.AttachedFileSaveInfo{
		Folder: "./",
	})
	assert.NoError(t, err)
}

func Test_Integration_SanitizeAllDefaults(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	// The Sanitize config in the regular org was obsoleted by a breaking
	// change, so the custom schema org is used instead.
	cfg := pangeatesting.IntegrationCustomSchemaConfig(t, testingEnvironment)
	cfg.PollResultTimeout = 2 * time.Minute
	client := sanitize.New(cfg)

	file, err := os.Open(TESTFILE_PATH)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resp, err := client.Sanitize(ctx, &sanitize.SanitizeRequest{
		TransferRequest: pangea.TransferRequest{
			TransferMethod: pangea.TMpostURL,
		},
		UploadedFileName: "uploaded_file",
	}, file)

	if err != nil {
		acceptedError, isAcceptedError := err.(*pangea.AcceptedError)
		if isAcceptedError {
			t.Logf("Result of request '%s' was not ready in time", *acceptedError.RequestID)
			return
		}
	}

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	assert.NotEmpty(t, resp.Result.DestURL)
	assert.Empty(t, resp.Result.DestShareID)
	assert.Nil(t, resp.Result.Data.Redact)
	assert.Greater(t, resp.Result.Data.Defang.ExternalURLsCount, 0)
	assert.Greater(t, resp.Result.Data.Defang.ExternalURLsCount, 0)
	assert.Equal(t, resp.Result.Data.Defang.DefangedCount, 0)
	assert.NotEmpty(t, resp.Result.Data.Defang.DomainIntelSummary)
	assert.False(t, resp.Result.Data.MaliciousFile)

	af, err := client.DownloadFile(ctx, *resp.Result.DestURL)
	assert.NoError(t, err)

	err = af.Save(pangea.AttachedFileSaveInfo{
		Folder: "./",
	})
	assert.NoError(t, err)
}

func Test_Integration_SanitizeMultipart(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	// The Sanitize config in the regular org was obsoleted by a breaking
	// change, so the custom schema org is used instead.
	cfg := pangeatesting.IntegrationCustomSchemaConfig(t, testingEnvironment)
	cfg.PollResultTimeout = 2 * time.Minute
	client := sanitize.New(cfg)

	file, err := os.Open(TESTFILE_PATH)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resp, err := client.Sanitize(ctx, &sanitize.SanitizeRequest{
		Content: &sanitize.SanitizeContent{
			URLIntel:            pangea.Bool(true),
			URLIntelProvider:    "crowdstrike",
			DomainIntel:         pangea.Bool(true),
			DomainIntelProvider: "crowdstrike",
			Defang:              pangea.Bool(true),
			DefangThreshold:     pangea.Int(20),
			Redact:              pangea.Bool(true),
		},
		ShareOutput: &sanitize.SanitizeShareOutput{
			Enabled: pangea.Bool(false),
		},
		File: &sanitize.SanitizeFile{
			ScanProvider: "crowdstrike",
		},
		TransferRequest: pangea.TransferRequest{
			TransferMethod: pangea.TMmultipart,
		},
		UploadedFileName: "uploaded_file",
	}, file)

	if err != nil {
		acceptedError, isAcceptedError := err.(*pangea.AcceptedError)
		if isAcceptedError {
			t.Logf("Result of request '%s' was not ready in time", *acceptedError.RequestID)
			return
		}
	}

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Result)
	assert.NotEmpty(t, resp.Result.DestURL)
	assert.Empty(t, resp.Result.DestShareID)
	assert.Greater(t, resp.Result.Data.Redact.RedactionCount, 0)
	assert.NotEmpty(t, resp.Result.Data.Redact.SummaryCounts)
	assert.Greater(t, resp.Result.Data.Defang.ExternalURLsCount, 0)
	assert.Greater(t, resp.Result.Data.Defang.ExternalURLsCount, 0)
	assert.Equal(t, resp.Result.Data.Defang.DefangedCount, 0)
	assert.NotEmpty(t, resp.Result.Data.Defang.DomainIntelSummary)
	assert.False(t, resp.Result.Data.MaliciousFile)

	af, err := client.DownloadFile(ctx, *resp.Result.DestURL)
	assert.NoError(t, err)

	err = af.Save(pangea.AttachedFileSaveInfo{
		Folder: "./",
	})
	assert.NoError(t, err)
}

func Test_Integration_Sanitize_SplitUpload_Post(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	// The Sanitize config in the regular org was obsoleted by a breaking
	// change, so the custom schema org is used instead.
	cfg := pangeatesting.IntegrationCustomSchemaConfig(t, testingEnvironment)
	client := sanitize.New(cfg)

	file, err := os.Open(TESTFILE_PATH)
	assert.NoError(t, err)

	params, err := pangea.GetUploadFileParams(file)
	assert.NoError(t, err)

	resp, err := client.RequestUploadURL(ctx, &sanitize.SanitizeRequest{
		Content: &sanitize.SanitizeContent{
			URLIntel:            pangea.Bool(true),
			URLIntelProvider:    "crowdstrike",
			DomainIntel:         pangea.Bool(true),
			DomainIntelProvider: "crowdstrike",
			Defang:              pangea.Bool(true),
			DefangThreshold:     pangea.Int(20),
			Redact:              pangea.Bool(true),
		},
		ShareOutput: &sanitize.SanitizeShareOutput{
			Enabled: pangea.Bool(false),
		},
		File: &sanitize.SanitizeFile{
			ScanProvider: "crowdstrike",
		},
		TransferRequest: pangea.TransferRequest{
			TransferMethod: pangea.TMpostURL,
		},
		UploadedFileName: "uploaded_file",
		CRC32C:           params.CRC32C,
		SHA256:           params.SHA256,
		Size:             pangea.Int(params.Size),
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.AcceptedResult)
	assert.NotEmpty(t, resp.AcceptedResult.PostURL)
	assert.Empty(t, resp.AcceptedResult.PutURL)

	url := resp.AcceptedResult.PostURL
	ud := resp.AcceptedResult.PostFormData

	fd := pangea.FileData{
		File:    file,
		Name:    "someName",
		Details: ud,
	}

	uploader := pangea.NewFileUploader()
	err = uploader.UploadFile(ctx, url, pangea.TMpostURL, fd)
	assert.NoError(t, err)

	var pollResponse *pangea.PangeaResponse[any]
	attempts := 0

	for attempts < 120/3 {
		// Wait until result should be ready
		time.Sleep(time.Duration(3 * time.Second))

		pollResponse, err = client.PollResultByID(ctx, *resp.RequestID, &sanitize.SanitizeResult{})
		if err == nil {
			break
		}
		attempts++
	}

	if err != nil {
		t.Logf("result of request '%s' was not ready in time\n", *resp.RequestID)
		return
	}

	assert.NoError(t, err)
	assert.NotNil(t, pollResponse)
	assert.NotNil(t, pollResponse.Result)

	r := (*pollResponse.Result).(*sanitize.SanitizeResult)
	assert.NotNil(t, r)
	assert.NotEmpty(t, r.DestURL)
	assert.Empty(t, r.DestShareID)
	assert.Greater(t, r.Data.Redact.RedactionCount, 0)
	assert.NotEmpty(t, r.Data.Redact.SummaryCounts)
	assert.Greater(t, r.Data.Defang.ExternalURLsCount, 0)
	assert.Greater(t, r.Data.Defang.ExternalURLsCount, 0)
	assert.Equal(t, r.Data.Defang.DefangedCount, 0)
	assert.NotEmpty(t, r.Data.Defang.DomainIntelSummary)
	assert.False(t, r.Data.MaliciousFile)
}

func Test_Integration_Sanitize_SplitUpload_Put(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	// The Sanitize config in the regular org was obsoleted by a breaking
	// change, so the custom schema org is used instead.
	cfg := pangeatesting.IntegrationCustomSchemaConfig(t, testingEnvironment)
	client := sanitize.New(cfg)

	file, err := os.Open(TESTFILE_PATH)
	assert.NoError(t, err)

	resp, err := client.RequestUploadURL(ctx, &sanitize.SanitizeRequest{
		Content: &sanitize.SanitizeContent{
			URLIntel:            pangea.Bool(true),
			URLIntelProvider:    "crowdstrike",
			DomainIntel:         pangea.Bool(true),
			DomainIntelProvider: "crowdstrike",
			Defang:              pangea.Bool(true),
			DefangThreshold:     pangea.Int(20),
			Redact:              pangea.Bool(true),
		},
		ShareOutput: &sanitize.SanitizeShareOutput{
			Enabled: pangea.Bool(false),
		},
		File: &sanitize.SanitizeFile{
			ScanProvider: "crowdstrike",
		},
		TransferRequest: pangea.TransferRequest{
			TransferMethod: pangea.TMputURL,
		},
		UploadedFileName: "uploaded_file",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.AcceptedResult)
	assert.Empty(t, resp.AcceptedResult.PostURL)
	assert.NotEmpty(t, resp.AcceptedResult.PutURL)

	url := resp.AcceptedResult.PutURL

	fd := pangea.FileData{
		File: file,
		Name: "someName",
	}

	uploader := pangea.NewFileUploader()
	err = uploader.UploadFile(ctx, url, pangea.TMputURL, fd)
	assert.NoError(t, err)

	var pollResponse *pangea.PangeaResponse[any]
	attempts := 0

	for attempts < 120/3 {
		// Wait until result should be ready
		time.Sleep(time.Duration(3 * time.Second))

		pollResponse, err = client.PollResultByID(ctx, *resp.RequestID, &sanitize.SanitizeResult{})
		if err == nil {
			break
		}
		attempts++
	}

	if err != nil {
		t.Logf("result of request '%s' was not ready in time\n", *resp.RequestID)
		return
	}

	assert.NotNil(t, pollResponse)
	assert.NotNil(t, pollResponse.Result)

	r := (*pollResponse.Result).(*sanitize.SanitizeResult)
	assert.NotNil(t, r)
	assert.NotEmpty(t, r.DestURL)
	assert.Empty(t, r.DestShareID)
	assert.Greater(t, r.Data.Redact.RedactionCount, 0)
	assert.NotEmpty(t, r.Data.Redact.SummaryCounts)
	assert.Greater(t, r.Data.Defang.ExternalURLsCount, 0)
	assert.Greater(t, r.Data.Defang.ExternalURLsCount, 0)
	assert.Equal(t, r.Data.Defang.DefangedCount, 0)
	assert.NotEmpty(t, r.Data.Defang.DomainIntelSummary)
	assert.False(t, r.Data.MaliciousFile)
}

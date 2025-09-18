package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/option"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/sanitize"
	"github.com/spf13/cobra"
)

const SANITIZE_TESTFILE_PATH = "./examples/sanitize/testdata/test-sanitize.txt"

func init() {
	sanitizeCmd := &cobra.Command{
		Use:   "sanitize",
		Short: "Sanitize examples",
	}

	sanitizeNoShareCmd := &cobra.Command{
		Use:  "sanitize_no_share",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return sanitizeNoShare(cmd)
		},
	}
	sanitizeCmd.AddCommand(sanitizeNoShareCmd)

	sanitizeAndShareCmd := &cobra.Command{
		Use:  "sanitize_and_share",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return sanitizeAndShare(cmd)
		},
	}
	sanitizeCmd.AddCommand(sanitizeAndShareCmd)

	sanitizeMultipartCmd := &cobra.Command{
		Use:  "sanitize_multipart",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return sanitizeMultipart(cmd)
		},
	}
	sanitizeCmd.AddCommand(sanitizeMultipartCmd)

	sanitizeSplitUploadPostURLCmd := &cobra.Command{
		Use:  "sanitize_split_upload_post_url",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return sanitizeSplitUploadPostURL(cmd)
		},
	}
	sanitizeCmd.AddCommand(sanitizeSplitUploadPostURLCmd)

	sanitizeSplitUploadPutURLCmd := &cobra.Command{
		Use:  "sanitize_split_upload_put_url",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return sanitizeSplitUploadPutURL(cmd)
		},
	}
	sanitizeCmd.AddCommand(sanitizeSplitUploadPutURLCmd)

	ExamplesCmd.AddCommand(sanitizeCmd)
}

func sanitizeAndShare(cmd *cobra.Command) error {
	// Load Pangea token from environment variables
	token := os.Getenv("PANGEA_SANITIZE_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present.")
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancelFn()

	// Create a new Sanitize client with Pangea token and domain
	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
		option.WithQueuedRetryEnabled(true),
		option.WithPollResultTimeout(120*time.Second),
	)
	if err != nil {
		log.Fatalf("expected no error got: %v", err)
	}
	client := sanitize.New(config)

	file, err := os.Open(SANITIZE_TESTFILE_PATH)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	// Get file upload params
	params, err := pangea.GetUploadFileParams(file)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	// Create a SanitizeRequest.
	// In this case TransferMethod is set to TMpostURL
	// Set CRC32C, SHA256 and Size
	input := &sanitize.SanitizeRequest{
		Content: &sanitize.SanitizeContent{
			URLIntel:            pangea.Bool(true),
			URLIntelProvider:    "crowdstrike",
			DomainIntel:         pangea.Bool(true),
			DomainIntelProvider: "crowdstrike",
			Defang:              pangea.Bool(true),
			DefangThreshold:     pangea.Int(20),
			Redact:              pangea.Bool(true),
		},
		// Enable Secure Share output and set output folder
		ShareOutput: &sanitize.SanitizeShareOutput{
			Enabled:      pangea.Bool(true),
			OutputFolder: "sdk_test/sanitize/",
		},
		File: &sanitize.SanitizeFile{
			ScanProvider: "crowdstrike",
		},
		UploadedFileName: "uploaded_file",
		TransferRequest: pangea.TransferRequest{
			TransferMethod: pangea.TMpostURL,
		},
		SHA256: params.SHA256,
		CRC32C: params.CRC32C,
		Size:   pangea.Int(params.Size),
	}

	fmt.Println("Sending Sanitize request...")
	resp, err := client.Sanitize(ctx, input, file)
	if err != nil {
		acceptedError, isAcceptedError := err.(*pangea.AcceptedError)
		if isAcceptedError {
			log.Printf("Result of request '%s' was not ready in time.\n", *acceptedError.RequestID)
			return err
		}

		log.Fatalf("Failed to process Sanitize request. Unexpected error: %v", err.Error())
	}

	fmt.Println("File Sanitized:")
	if resp.Result.DestShareID != nil {
		fmt.Printf("\tShare ID: %s\n", *resp.Result.DestShareID)
	}
	if resp.Result.DestURL != nil {
		fmt.Printf("\tDest URL: %s\n", *resp.Result.DestURL)
	}

	fmt.Printf("\tRedact data: %s\n", pangea.Stringify(resp.Result.Data.Redact))
	fmt.Printf("\tDefang data: %s\n", pangea.Stringify(resp.Result.Data.Defang))

	if resp.Result.Data.MaliciousFile {
		fmt.Println("File IS malicious")
	} else {
		fmt.Println("File is NOT malicious")
	}
	return nil
}

func sanitizeMultipart(cmd *cobra.Command) error {
	// Load Pangea token from environment variables
	token := os.Getenv("PANGEA_SANITIZE_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present.")
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancelFn()

	// Create a new Sanitize client with Pangea token and domain
	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
		option.WithQueuedRetryEnabled(true),
		option.WithPollResultTimeout(120*time.Second),
	)
	if err != nil {
		log.Fatalf("expected no error got: %v", err)
	}
	client := sanitize.New(config)

	file, err := os.Open(SANITIZE_TESTFILE_PATH)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	// Create a SanitizeRequest.
	// In this case TransferMethod is set to TMmultipartURL
	input := &sanitize.SanitizeRequest{
		Content: &sanitize.SanitizeContent{
			URLIntel:            pangea.Bool(true),
			URLIntelProvider:    "crowdstrike",
			DomainIntel:         pangea.Bool(true),
			DomainIntelProvider: "crowdstrike",
			Defang:              pangea.Bool(true),
			DefangThreshold:     pangea.Int(20),
			Redact:              pangea.Bool(true),
		},
		// Enable Secure Share output and set output folder
		ShareOutput: &sanitize.SanitizeShareOutput{
			Enabled:      pangea.Bool(true),
			OutputFolder: "sdk_test/sanitize/",
		},
		File: &sanitize.SanitizeFile{
			ScanProvider: "crowdstrike",
		},
		UploadedFileName: "uploaded_file",
		TransferRequest: pangea.TransferRequest{
			TransferMethod: pangea.TMmultipart,
		},
	}

	fmt.Println("Sending Sanitize request as multipart...")
	resp, err := client.Sanitize(ctx, input, file)
	if err != nil {
		acceptedError, isAcceptedError := err.(*pangea.AcceptedError)
		if isAcceptedError {
			log.Printf("Result of request '%s' was not ready in time.\n", *acceptedError.RequestID)
			return err
		}

		log.Fatalf("Failed to process Sanitize request. Unexpected error: %v", err.Error())
	}

	fmt.Println("File Sanitized:")
	if resp.Result.DestShareID != nil {
		fmt.Printf("\tShare ID: %s\n", *resp.Result.DestShareID)
	}
	if resp.Result.DestURL != nil {
		fmt.Printf("\tDest URL: %s\n", *resp.Result.DestURL)
	}

	fmt.Printf("\tRedact data: %s\n", pangea.Stringify(resp.Result.Data.Redact))
	fmt.Printf("\tDefang data: %s\n", pangea.Stringify(resp.Result.Data.Defang))

	if resp.Result.Data.MaliciousFile {
		fmt.Println("File IS malicious")
	} else {
		fmt.Println("File is NOT malicious")
	}
	return nil
}

func sanitizeSplitUploadPostURL(cmd *cobra.Command) error {
	// Load Pangea token from environment variables
	token := os.Getenv("PANGEA_SANITIZE_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present.")
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancelFn()

	// Create a new Sanitize client with Pangea token and domain
	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
		option.WithQueuedRetryEnabled(true),
		option.WithPollResultTimeout(120*time.Second),
	)
	if err != nil {
		log.Fatalf("expected no error got: %v", err)
	}
	client := sanitize.New(config)

	file, err := os.Open(SANITIZE_TESTFILE_PATH)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	// Get file upload params
	params, err := pangea.GetUploadFileParams(file)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	// Create a SanitizeRequest to request an presigned upload url.
	// In this case TransferMethod is set to TMpostURL
	// Set CRC32C, SHA256 and Size
	input := &sanitize.SanitizeRequest{
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
		UploadedFileName: "uploaded_file",
		TransferRequest: pangea.TransferRequest{
			TransferMethod: pangea.TMpostURL,
		},
		SHA256: params.SHA256,
		CRC32C: params.CRC32C,
		Size:   pangea.Int(params.Size),
	}

	fmt.Println("Request upload URL...")
	resp, err := client.RequestUploadURL(ctx, input)
	if err != nil {
		log.Fatalf("Failed to get upload URL. Unexpected error: %v", err.Error())
	}

	// Get presigned url and data to post
	url := resp.AcceptedResult.PostURL
	data := resp.AcceptedResult.PostFormData
	fmt.Printf("Got URL: %s\n", url)

	fd := pangea.FileData{
		File:    file,
		Details: data,
		Name:    "someName",
	}

	// Create an uploader
	uploader := pangea.NewFileUploader()

	fmt.Println("Uploading file...")
	// Upload the file to the returned upload URL.
	// Need to set transfer method again to TMpostURL
	err = uploader.UploadFile(ctx, url, pangea.TMpostURL, fd)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	var sr *pangea.PangeaResponse[any]
	i := 0
	maxRetry := 24

	// Try to poll result
	for i < maxRetry {
		// Wait until result should be ready
		time.Sleep(time.Duration(10 * time.Second))
		fmt.Printf("Trying to poll result. Retry: %d\n", i)
		sr, err = client.PollResultByID(ctx, *resp.RequestID, &sanitize.SanitizeResult{})
		if err == nil {
			break
		}
		i++
		fmt.Println("Result is not ready yet. Retrying...")
	}

	if i >= maxRetry {
		log.Println("Result was not ready in time.")
		return nil
	}

	// After receiving the result, cast it before using it.
	rSanitize := (*sr.Result).(*sanitize.SanitizeResult)

	fmt.Println("File Sanitized:")
	if rSanitize.DestShareID != nil {
		fmt.Printf("\tShare ID: %s\n", *rSanitize.DestShareID)
	}
	if rSanitize.DestURL != nil {
		fmt.Printf("\tDest URL: %s\n", *rSanitize.DestURL)
	}

	fmt.Printf("\tRedact data: %s\n", pangea.Stringify(rSanitize.Data.Redact))
	fmt.Printf("\tDefang data: %s\n", pangea.Stringify(rSanitize.Data.Defang))

	if rSanitize.Data.MaliciousFile {
		fmt.Println("File IS malicious")
	} else {
		fmt.Println("File is NOT malicious")
	}
	return nil
}

func sanitizeSplitUploadPutURL(cmd *cobra.Command) error {
	// Load Pangea token from environment variables
	token := os.Getenv("PANGEA_SANITIZE_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present.")
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancelFn()

	// Create a new Sanitize client with Pangea token and domain
	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
		option.WithQueuedRetryEnabled(true),
		option.WithPollResultTimeout(120*time.Second),
	)
	if err != nil {
		log.Fatalf("expected no error got: %v", err)
	}
	client := sanitize.New(config)

	// Create a SanitizeRequest to request an presigned upload url.
	// In this case TransferMethod is set to TMputURL
	input := &sanitize.SanitizeRequest{
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
		UploadedFileName: "uploaded_file",
		TransferRequest: pangea.TransferRequest{
			TransferMethod: pangea.TMputURL,
		},
	}

	fmt.Println("Request upload URL...")
	resp, err := client.RequestUploadURL(ctx, input)
	if err != nil {
		log.Fatalf("Failed to get upload URL. Unexpected error: %v", err.Error())
	}

	// Get presigned url
	url := resp.AcceptedResult.PutURL
	fmt.Printf("Got URL: %s\n", url)

	file, err := os.Open(SANITIZE_TESTFILE_PATH)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	fd := pangea.FileData{
		File: file,
		Name: "someName",
	}

	// Create an uploader
	uploader := pangea.NewFileUploader()

	fmt.Println("Uploading file...")
	// Upload the file to the returned upload URL.
	// Need to set transfer method again to TMputURL
	err = uploader.UploadFile(ctx, url, pangea.TMputURL, fd)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	var sr *pangea.PangeaResponse[any]
	i := 0
	maxRetry := 24

	// Try to poll result
	for i < maxRetry {
		// Wait until result should be ready
		time.Sleep(time.Duration(10 * time.Second))

		fmt.Printf("Trying to poll result. Retry: %d\n", i)
		sr, err = client.PollResultByID(ctx, *resp.RequestID, &sanitize.SanitizeResult{})
		if err == nil {
			break
		}
		i++
		fmt.Println("Result is not ready yet. Retrying...")
	}

	if i >= maxRetry {
		log.Println("Result was not ready in time.")
		return nil
	}

	// After receiving the result, cast it before using it.
	rSanitize := (*sr.Result).(*sanitize.SanitizeResult)

	fmt.Println("File Sanitized:")
	if rSanitize.DestShareID != nil {
		fmt.Printf("\tShare ID: %s\n", *rSanitize.DestShareID)
	}
	if rSanitize.DestURL != nil {
		fmt.Printf("\tDest URL: %s\n", *rSanitize.DestURL)
	}

	fmt.Printf("\tRedact data: %s\n", pangea.Stringify(rSanitize.Data.Redact))
	fmt.Printf("\tDefang data: %s\n", pangea.Stringify(rSanitize.Data.Defang))

	if rSanitize.Data.MaliciousFile {
		fmt.Println("File IS malicious")
	} else {
		fmt.Println("File is NOT malicious")
	}
	return nil
}

func sanitizeNoShare(cmd *cobra.Command) error {
	// Load Pangea token from environment variables
	token := os.Getenv("PANGEA_SANITIZE_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present.")
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancelFn()

	// create a new sanitize client with pangea token and domain
	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
		option.WithQueuedRetryEnabled(true),
		option.WithPollResultTimeout(120*time.Second),
	)
	if err != nil {
		log.Fatalf("expected no error got: %v", err)
	}
	client := sanitize.New(config)

	file, err := os.Open(SANITIZE_TESTFILE_PATH)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	// Get file upload params
	params, err := pangea.GetUploadFileParams(file)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	// Create a SanitizeRequest.
	// In this case TransferMethod is set to TMpostURL
	// Set CRC32C, SHA256 and Size
	input := &sanitize.SanitizeRequest{
		Content: &sanitize.SanitizeContent{
			URLIntel:            pangea.Bool(true),
			URLIntelProvider:    "crowdstrike",
			DomainIntel:         pangea.Bool(true),
			DomainIntelProvider: "crowdstrike",
			Defang:              pangea.Bool(true),
			DefangThreshold:     pangea.Int(20),
			Redact:              pangea.Bool(true),
		},
		// Disable Secure Share output
		ShareOutput: &sanitize.SanitizeShareOutput{
			Enabled: pangea.Bool(false),
		},
		File: &sanitize.SanitizeFile{
			ScanProvider: "crowdstrike",
		},
		UploadedFileName: "uploaded_file",
		TransferRequest: pangea.TransferRequest{
			TransferMethod: pangea.TMpostURL,
		},
		SHA256: params.SHA256,
		CRC32C: params.CRC32C,
		Size:   pangea.Int(params.Size),
	}

	fmt.Println("Sending Sanitize request...")
	resp, err := client.Sanitize(ctx, input, file)
	if err != nil {
		acceptedError, isAcceptedError := err.(*pangea.AcceptedError)
		if isAcceptedError {
			log.Printf("Result of request '%s' was not ready in time.\n", *acceptedError.RequestID)
			return err
		}

		log.Fatalf("Failed to process Sanitize request. Unexpected error: %v", err.Error())
	}

	fmt.Println("File Sanitized:")
	if resp.Result.DestShareID != nil {
		fmt.Printf("\tShare ID: %s\n", *resp.Result.DestShareID)
	}
	if resp.Result.DestURL != nil {
		fmt.Printf("\tDest URL: %s\n", *resp.Result.DestURL)
	}

	fmt.Printf("\tRedact data: %s\n", pangea.Stringify(resp.Result.Data.Redact))
	fmt.Printf("\tDefang data: %s\n", pangea.Stringify(resp.Result.Data.Defang))

	if resp.Result.Data.MaliciousFile {
		fmt.Println("File IS malicious")
	} else {
		fmt.Println("File is NOT malicious")
	}
	return nil
}

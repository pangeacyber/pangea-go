package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/option"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/file_scan"
	"github.com/spf13/cobra"
)

var file string

func init() {
	fileScanCmd := &cobra.Command{
		Use:   "file_scan",
		Short: "File Scan examples",
	}

	fileScanSyncCrowdstrikeCmd := &cobra.Command{
		Use:  "file_scan_sync_crowdstrike",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return fileScanSyncCrowdstrike(cmd)
		},
	}
	fileScanSyncCrowdstrikeCmd.Flags().StringVarP(&file, "file", "f", "", "Path to a file to scan.")
	fileScanCmd.AddCommand(fileScanSyncCrowdstrikeCmd)

	fileScanSyncReversinglabsCmd := &cobra.Command{
		Use:  "file_scan_sync_reversinglabs",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return fileScanSyncReversinglabs(cmd)
		},
	}
	fileScanSyncReversinglabsCmd.Flags().StringVarP(&file, "file", "f", "", "Path to a file to scan.")
	fileScanCmd.AddCommand(fileScanSyncReversinglabsCmd)

	fileScanAsyncCrowdstrikeCmd := &cobra.Command{
		Use:  "file_scan_async_crowdstrike",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return fileScanAsyncCrowdstrike(cmd)
		},
	}
	fileScanAsyncCrowdstrikeCmd.Flags().StringVarP(&file, "file", "f", "", "Path to a file to scan.")
	fileScanCmd.AddCommand(fileScanAsyncCrowdstrikeCmd)

	fileScanAsyncReversinglabsCmd := &cobra.Command{
		Use:  "file_scan_async_reversinglabs",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return fileScanAsyncReversinglabs(cmd)
		},
	}
	fileScanAsyncReversinglabsCmd.Flags().StringVarP(&file, "file", "f", "", "Path to a file to scan.")
	fileScanCmd.AddCommand(fileScanAsyncReversinglabsCmd)

	fileScanRequestUploadURLPostCmd := &cobra.Command{
		Use:  "file_scan_request_upload_url_post",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return fileScanRequestUploadURLPost(cmd)
		},
	}
	fileScanRequestUploadURLPostCmd.Flags().StringVarP(&file, "file", "f", "", "Path to a file to scan.")
	fileScanCmd.AddCommand(fileScanRequestUploadURLPostCmd)

	fileScanRequestUploadURLPutCmd := &cobra.Command{
		Use:  "file_scan_request_upload_url_put",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return fileScanRequestUploadURLPut(cmd)
		},
	}
	fileScanRequestUploadURLPutCmd.Flags().StringVarP(&file, "file", "f", "", "Path to a file to scan.")
	fileScanCmd.AddCommand(fileScanRequestUploadURLPutCmd)

	ExamplesCmd.AddCommand(fileScanCmd)
}

func fileScanAsyncCrowdstrike(cmd *cobra.Command) error {
	fmt.Println("File Scan async start...")
	token := os.Getenv("PANGEA_FILE_SCAN_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancelFn()

	// To enable async mode, set QueuedRetryEnabled to false
	// When .Scan() is called it will return an AcceptedError immediately when server returns a 202 response
	config, err := pangea.NewConfig(option.WithToken(token), option.WithDomain(os.Getenv("PANGEA_DOMAIN")), option.WithQueuedRetryEnabled(false))
	if err != nil {
		log.Fatalf("expected no error got: %v", err)
	}
	client := file_scan.New(config)

	input := &file_scan.FileScanRequest{
		Raw:      true,
		Verbose:  true,
		Provider: "crowdstrike",
	}

	// This should be your own file to scan
	file, err := os.Open(file)
	if err != nil {
		log.Fatalf("expected no error got: %v", err)
	}

	fmt.Println("File Scan request...")
	sr, err := client.Scan(ctx, input, file)
	if sr != nil {
		// .Scan() will return a response immediately if the file has a known reputation
		fmt.Println("File Scan success on first attempt.")
		fmt.Println(pangea.Stringify(sr.Result))
		return nil
	}

	ae, ok := err.(*pangea.AcceptedError)
	if !ok {
		log.Fatalf("Unexpected error. This should be AcceptedError: %v", err.Error())
	}
	fmt.Println("Accepted error received (as expected).")

	var pr *pangea.PangeaResponse[any]
	i := 0
	maxRetry := 24

	fmt.Println("Let's try to poll result...")
	for i < maxRetry {
		// Wait for result
		time.Sleep(time.Duration(10 * time.Second))

		pr, err = client.PollResultByError(ctx, *ae)
		if err == nil {
			break
		}
		i++
		fmt.Printf("Result is not ready yet. Retry: %d\n", i)
	}

	if i == maxRetry {
		log.Fatal("Result still not ready")
	} else {
		r := (*pr.Result).(*file_scan.FileScanResult)
		fmt.Println("File Scan success.")
		fmt.Println(pangea.Stringify(r))
	}
	return nil
}

func fileScanAsyncReversinglabs(cmd *cobra.Command) error {
	fmt.Println("File Scan async start...")
	token := os.Getenv("PANGEA_FILE_SCAN_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancelFn()

	// To enable async mode, set QueuedRetryEnabled to false
	// When .Scan() is called it will return an AcceptedError immediately when server returns a 202 response
	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
		option.WithQueuedRetryEnabled(false),
	)
	if err != nil {
		log.Fatalf("expected no error got: %v", err)
	}
	client := file_scan.New(config)

	input := &file_scan.FileScanRequest{
		Raw:      true,
		Verbose:  true,
		Provider: "reversinglabs",
	}

	// This should be your own file to scan
	file, err := os.Open(file)
	if err != nil {
		log.Fatalf("expected no error got: %v", err)
	}

	fmt.Println("File Scan request...")
	sr, err := client.Scan(ctx, input, file)
	if sr != nil {
		// .Scan() will return a response immediately if the file has a known reputation
		fmt.Println("File Scan success on first attempt.")
		fmt.Println(pangea.Stringify(sr.Result))
		return nil
	}

	ae, ok := err.(*pangea.AcceptedError)
	if !ok {
		log.Fatalf("Unexpected error. This should be AcceptedError: %v", err.Error())
	}
	fmt.Println("Accepted error received (as expected).")

	var pr *pangea.PangeaResponse[any]
	i := 0
	maxRetry := 24

	fmt.Println("Let's try to poll result...")
	for i < maxRetry {
		// Wait for result
		time.Sleep(time.Duration(10 * time.Second))

		pr, err = client.PollResultByError(ctx, *ae)
		if err == nil {
			break
		}
		i++
		fmt.Printf("Result is not ready yet. Retry: %d\n", i)
	}

	if i == maxRetry {
		log.Fatal("Result still not ready")
	} else {
		r := (*pr.Result).(*file_scan.FileScanResult)
		fmt.Println("File Scan success.")
		fmt.Println(pangea.Stringify(r))
	}
	return nil
}

func fileScanRequestUploadURLPost(cmd *cobra.Command) error {
	fmt.Println("File Scan start...")
	token := os.Getenv("PANGEA_FILE_SCAN_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancelFn()

	// To enable sync mode, set QueuedRetryEnabled to true and set a timeout
	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
		option.WithQueuedRetryEnabled(true),
		option.WithPollResultTimeout(30*time.Second),
	)
	if err != nil {
		log.Fatalf("expected no error got: %v", err)
	}
	client := file_scan.New(config)

	file, err := os.Open(file)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	// calculate file info needed to request upload url
	params, err := file_scan.GetUploadFileParams(file)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	input := &file_scan.FileScanGetURLRequest{
		Raw:            true,
		Verbose:        true,
		Provider:       "reversinglabs",
		TransferMethod: pangea.TMpostURL,
		FileParams:     params,
	}

	// request an upload url
	resp, err := client.RequestUploadURL(ctx, input, file)
	if err != nil {
		log.Fatalf("unexpected error: %v", err.Error())
	}

	// extract upload url and upload details that should be posted with the file
	url := resp.AcceptedResult.PostURL
	ud := resp.AcceptedResult.PostFormData
	fmt.Printf("Got upload url: %s\n", url)

	fd := pangea.FileData{
		File:    file,
		Name:    "someName",
		Details: ud,
	}

	// Create an uploader and upload the file
	uploader := file_scan.NewFileUploader()
	err = uploader.UploadFile(ctx, url, pangea.TMpostURL, fd)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}
	fmt.Println("Upload file success")

	var pr *pangea.PangeaResponse[any]
	i := 0
	maxRetry := 24

	fmt.Println("Let's try to poll result...")
	for i < maxRetry {
		// Wait for result
		time.Sleep(time.Duration(10 * time.Second))

		pr, err = client.PollResultByID(ctx, *resp.RequestID, &file_scan.FileScanResult{})
		if err == nil {
			break
		}
		i++
		fmt.Printf("Result is not ready yet. Retry: %d\n", i)
	}

	if i == maxRetry {
		log.Fatal("Result still not ready")
	} else {
		r := (*pr.Result).(*file_scan.FileScanResult)
		fmt.Println("File Scan success.")
		fmt.Println(pangea.Stringify(r))
	}
	return nil
}

func fileScanRequestUploadURLPut(cmd *cobra.Command) error {
	fmt.Println("File Scan start...")
	token := os.Getenv("PANGEA_FILE_SCAN_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancelFn()

	// To enable sync mode, set QueuedRetryEnabled to true and set a timeout
	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
		option.WithQueuedRetryEnabled(true),
		option.WithPollResultTimeout(30*time.Second),
	)
	if err != nil {
		log.Fatalf("expected no error got: %v", err)
	}
	client := file_scan.New(config)

	file, err := os.Open(file)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	// Only the TransferMethod is needed when using TMputURL, in addition to the standard parameters
	input := &file_scan.FileScanGetURLRequest{
		Raw:            true,
		Verbose:        true,
		Provider:       "reversinglabs",
		TransferMethod: pangea.TMputURL,
	}

	// request an upload url
	resp, err := client.RequestUploadURL(ctx, input, file)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	// Extract the upload url
	// File details is not needed when using TransferMethod TMputURL
	url := resp.AcceptedResult.PutURL
	fmt.Printf("Got upload url: %s\n", url)

	fd := pangea.FileData{
		File: file,
		Name: "someName",
	}

	// Create an uploader and upload the file
	uploader := file_scan.NewFileUploader()
	err = uploader.UploadFile(ctx, url, pangea.TMputURL, fd)
	if err != nil {
		log.Fatalf("unexpected error: %v", err.Error())
	}
	fmt.Println("Upload file success")

	var pr *pangea.PangeaResponse[any]
	i := 0
	maxRetry := 24

	fmt.Println("Let's try to poll result...")
	for i < maxRetry {
		// Wait for result
		time.Sleep(time.Duration(10 * time.Second))

		pr, err = client.PollResultByID(ctx, *resp.RequestID, &file_scan.FileScanResult{})
		if err == nil {
			break
		}
		i++
		fmt.Printf("Result is not ready yet. Retry: %d\n", i)
	}

	if i == maxRetry {
		log.Fatal("Result still not ready")
	} else {
		r := (*pr.Result).(*file_scan.FileScanResult)
		fmt.Println("File Scan success.")
		fmt.Println(pangea.Stringify(r))
	}
	return nil
}

func fileScanSyncReversinglabs(cmd *cobra.Command) error {
	fmt.Println("File Scan sync start...")
	token := os.Getenv("PANGEA_FILE_SCAN_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancelFn()

	// To enable sync mode, set QueuedRetryEnabled to true and set a timeout
	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
		option.WithQueuedRetryEnabled(true),
		option.WithPollResultTimeout(60*time.Second),
	)
	if err != nil {
		log.Fatalf("expected no error got: %v", err)
	}
	client := file_scan.New(config)

	input := &file_scan.FileScanRequest{
		Raw:      true,
		Verbose:  true,
		Provider: "reversinglabs",
	}

	// This should be your own file to scan
	file, err := os.Open(file)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	fmt.Println("File Scan request...")
	resp, err := client.Scan(ctx, input, file)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("File Scan success.")
	fmt.Println(pangea.Stringify(resp.Result))
	return nil
}

func fileScanSyncCrowdstrike(cmd *cobra.Command) error {
	fmt.Println("File Scan sync start...")
	token := os.Getenv("PANGEA_FILE_SCAN_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancelFn()

	// To enable sync mode, set QueuedRetryEnabled to true and set a timeout
	config, err := pangea.NewConfig(
		option.WithToken(token),
		option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
		option.WithQueuedRetryEnabled(true),
		option.WithPollResultTimeout(60*time.Second),
	)
	if err != nil {
		log.Fatalf("expected no error got: %v", err)
	}
	client := file_scan.New(config)

	input := &file_scan.FileScanRequest{
		Raw:      true,
		Verbose:  true,
		Provider: "crowdstrike",
	}

	// This should be your own file to scan
	file, err := os.Open(file)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	fmt.Println("File Scan request...")
	resp, err := client.Scan(ctx, input, file)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("File Scan success.")
	fmt.Println(pangea.Stringify(resp.Result))
	return nil
}

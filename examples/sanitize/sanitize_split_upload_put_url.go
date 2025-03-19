package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/service/sanitize"
)

func main() {
	// Set filePath to your own file
	const filePath = "./testdata/test-sanitize.txt"

	// Load Pangea token from environment variables
	token := os.Getenv("PANGEA_SANITIZE_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present.")
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancelFn()

	// Create a new Sanitize client with Pangea token and domain
	client := sanitize.New(&pangea.Config{
		Token:              token,
		BaseURLTemplate:    os.Getenv("PANGEA_URL_TEMPLATE"),
		QueuedRetryEnabled: true,
		PollResultTimeout:  120 * time.Second,
		Retry:              true,
		RetryConfig: &pangea.RetryConfig{
			RetryMax: 4,
		},
	})

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

	file, err := os.Open(filePath)
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
		log.Fatal("Failed to poll result. Max retry reached.")
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

}

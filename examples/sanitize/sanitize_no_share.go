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

	// create a new sanitize client with pangea token and domain
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

	file, err := os.Open(filePath)
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

}

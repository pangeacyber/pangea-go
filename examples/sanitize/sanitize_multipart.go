package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/sanitize"
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
		Domain:             os.Getenv("PANGEA_DOMAIN"),
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
			return
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

}

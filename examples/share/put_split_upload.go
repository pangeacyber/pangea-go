package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/service/share"
)

func main() {
	var t = time.Now().Format("20060102_150405")
	var name = "file_name_" + t
	const filePath = "./testdata/testfile.pdf"

	// Load pangea token from environment variables
	token := os.Getenv("PANGEA_SHARE_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present.")
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancelFn()

	// create a new store client with pangea token and domain
	client := share.New(&pangea.Config{
		Token:              token,
		Domain:             os.Getenv("PANGEA_DOMAIN"),
		QueuedRetryEnabled: true,
		PollResultTimeout:  120 * time.Second,
		Retry:              true,
		RetryConfig: &pangea.RetryConfig{
			RetryMax: 4,
		},
	})

	// Create a PutRequest to request an presigned upload url.
	// In this case TransferMethod is set to TMputURL
	input := &share.PutRequest{
		Name: name,
		TransferRequest: pangea.TransferRequest{
			TransferMethod: pangea.TMputURL,
		},
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	resp, err := client.RequestUploadURL(ctx, input)
	if err != nil {
		log.Fatalf("unexpected error: %v", err.Error())
	}

	// Get presigned url
	url := resp.AcceptedResult.PutURL

	fd := pangea.FileData{
		File: file,
		Name: "someName",
	}

	// Create an upload
	uploader := share.NewFileUploader()

	// Upload the file to the url get previously
	// Need to set transfer method again to TMputURL
	err = uploader.UploadFile(ctx, url, pangea.TMputURL, fd)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	var pr *pangea.PangeaResponse[any]
	i := 0

	// Try to poll result
	for i < 24 {
		// Wait until result should be ready
		time.Sleep(time.Duration(10 * time.Second))

		pr, err = client.PollResultByID(ctx, *resp.RequestID, &share.PutResult{})
		if err == nil {
			break
		}
		i++
	}

	// Once got the result, cast it to use it
	rPut := (*pr.Result).(*share.PutResult)

	fmt.Println("File uploaded:")
	fmt.Printf("\tID: %s\n", rPut.Object.ID)
	fmt.Printf("\tName: %s\n", rPut.Object.Name)
}

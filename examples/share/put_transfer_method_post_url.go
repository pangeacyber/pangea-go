package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/share"
)

func main() {
	var t = time.Now().Format("20060102_150405")
	var name = "file_name_" + t
	const filePath = "./testdata/testfile.pdf"

	// Load Pangea token from environment variables
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

	// Create a PutRequest. In this case TransferMethod is set to TMpostURL
	// So SDK is going to request a post url, upload the file to that url and then request to pangea for the /put result
	input := &share.PutRequest{
		Name: name,
		TransferRequest: pangea.TransferRequest{
			TransferMethod: pangea.TMpostURL,
		},
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	rPut, err := client.Put(ctx, input, file)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	fmt.Println("File uploaded:")
	fmt.Printf("\tID: %s\n", rPut.Result.Object.ID)
	fmt.Printf("\tName: %s\n", rPut.Result.Object.Name)
}

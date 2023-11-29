package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/service/file_scan"
)

const TESTFILE_PATH = "./testdata/testfile.pdf"

func main() {
	fmt.Println("File Scan sync start...")
	token := os.Getenv("PANGEA_FILE_SCAN_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancelFn()

	// To enable sync mode, set QueuedRetryEnabled to true and set a timeout
	client := file_scan.New(&pangea.Config{
		Token:              token,
		Domain:             os.Getenv("PANGEA_DOMAIN"),
		QueuedRetryEnabled: true,
		PollResultTimeout:  60 * time.Second,
	})

	input := &file_scan.FileScanRequest{
		Raw:      true,
		Verbose:  true,
		Provider: "crowdstrike",
	}

	// This should be your own file to scan
	file, err := os.Open(TESTFILE_PATH)
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
}

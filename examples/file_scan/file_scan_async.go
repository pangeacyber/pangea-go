package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/file_scan"
)

const TESTFILE_PATH = "./testdata/testfile.pdf"

func main() {
	fmt.Println("File Scan async start...")
	token := os.Getenv("PANGEA_FILE_SCAN_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancelFn()

	// To work in async it's need to set up QueuedRetryEnabled to false
	// When we call .Scan() it will return an AcceptedError inmediatly if server return a 202 response
	client := file_scan.New(&pangea.Config{
		Token:              token,
		Domain:             os.Getenv("PANGEA_DOMAIN"),
		QueuedRetryEnabled: false,
	})

	input := &file_scan.FileScanRequest{
		Raw:      true,
		Verbose:  true,
		Provider: "crowdstrike",
	}

	// This should be your own file to scan
	file, err := os.Open(TESTFILE_PATH)
	if err != nil {
		log.Fatal("expected no error got: %v", err)
	}

	fmt.Println("File Scan request...")
	sr, err := client.Scan(ctx, input, file)
	if sr != nil {
		// this could happen if reputation call inside scan success
		fmt.Println("File Scan success on first attempt.")
		fmt.Println(pangea.Stringify(sr.Result))
		os.Exit(0)
	}

	ae, ok := err.(*pangea.AcceptedError)
	if ok == false {
		log.Fatal("Unexpected error. This should be AcceptedError: %v", err.Error())
	}
	fmt.Println("Accepted error received (as expected).")

	fmt.Println("Sleep some time until result should be ready.")
	// Wait until result should be ready
	time.Sleep(time.Duration(30 * time.Second))

	fmt.Println("File Scan poll result...")
	pr, err := client.PollResultByError(ctx, *ae)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("File Scan poll result success.")
	fmt.Println(pangea.Stringify(pr.Result))
}

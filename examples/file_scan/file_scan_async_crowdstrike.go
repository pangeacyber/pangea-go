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
	fmt.Println("File Scan async start...")
	token := os.Getenv("PANGEA_FILE_SCAN_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancelFn()

	// To enable async mode, set QueuedRetryEnabled to false
	// When .Scan() is called it will return an AcceptedError immediately when server returns a 202 response
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
		// .Scan() will return a response immediately if the file has a known reputation
		fmt.Println("File Scan success on first attempt.")
		fmt.Println(pangea.Stringify(sr.Result))
		os.Exit(0)
	}

	ae, ok := err.(*pangea.AcceptedError)
	if ok == false {
		log.Fatal("Unexpected error. This should be AcceptedError: %v", err.Error())
	}
	fmt.Println("Accepted error received (as expected).")

	fmt.Println("Sleep some time before polling.")
	// multiple polling attempts may be required
	time.Sleep(time.Duration(20 * time.Second))

	fmt.Println("File Scan poll result...")
	pr, err := client.PollResultByError(ctx, *ae)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("File Scan poll result success.")
	fmt.Println(pangea.Stringify(pr.Result))
}

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/service/file_scan"
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
		BaseURLTemplate:    os.Getenv("PANGEA_URL_TEMPLATE"),
		QueuedRetryEnabled: false,
	})

	input := &file_scan.FileScanRequest{
		Raw:      true,
		Verbose:  true,
		Provider: "reversinglabs",
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
}

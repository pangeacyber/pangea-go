// intel domain lookup is an example of how to use the lookup method
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/service/file_scan"
)

const TESTFILE_PATH = "./testdata/testfile.pdf"

func main() {
	fmt.Println("File Scan sync start...")
	token := os.Getenv("PANGEA_INTEL_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancelFn()

	// To work in sync it's need to set up QueuedRetryEnabled to true and set up a proper timeout
	// If timeout is so little service won't end up and will return an AcceptedError anyway
	client := file_scan.New(&pangea.Config{
		Token:              token,
		Domain:             os.Getenv("PANGEA_DOMAIN"),
		QueuedRetryEnabled: true,
		PollResultTimeout:  60 * time.Second,
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
	resp, err := client.Scan(ctx, input, file)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("File Scan success.")
	fmt.Println(pangea.Stringify(resp.Result))
}
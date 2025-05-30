package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/file_scan"
)

const TESTFILE_PATH = "./testdata/testfile.pdf"

func main() {
	fmt.Println("File Scan start...")
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
		PollResultTimeout:  30 * time.Second,
	})

	file, err := os.Open(TESTFILE_PATH)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	// calculate file info needed to request upload url
	params, err := file_scan.GetUploadFileParams(file)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	input := &file_scan.FileScanGetURLRequest{
		Raw:            true,
		Verbose:        true,
		Provider:       "reversinglabs",
		TransferMethod: pangea.TMpostURL,
		FileParams:     params,
	}

	// request an upload url
	resp, err := client.RequestUploadURL(ctx, input, file)
	if err != nil {
		log.Fatalf("unexpected error: %v", err.Error())
	}

	// extract upload url and upload details that should be posted with the file
	url := resp.AcceptedResult.PostURL
	ud := resp.AcceptedResult.PostFormData
	fmt.Printf("Got upload url: %s\n", url)

	fd := pangea.FileData{
		File:    file,
		Name:    "someName",
		Details: ud,
	}

	// Create an uploader and upload the file
	uploader := file_scan.NewFileUploader()
	err = uploader.UploadFile(ctx, url, pangea.TMpostURL, fd)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}
	fmt.Println("Upload file success")

	var pr *pangea.PangeaResponse[any]
	i := 0
	maxRetry := 24

	fmt.Println("Let's try to poll result...")
	for i < maxRetry {
		// Wait for result
		time.Sleep(time.Duration(10 * time.Second))

		pr, err = client.PollResultByID(ctx, *resp.RequestID, &file_scan.FileScanResult{})
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

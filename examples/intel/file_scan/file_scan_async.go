// intel domain lookup is an example of how to use the lookup method
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/service/file_scan"
)

const EICAR = "X5O!P%@AP[4\\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*\n"

func main() {
	fmt.Println("File Scan async start...")
	token := os.Getenv("PANGEA_INTEL_TOKEN")
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
		Provider: "reversinglabs",
	}

	// Here we create a file that will give us a malicious result as example,
	// This should be your own file to scan
	file := strings.NewReader(EICAR)

	fmt.Println("File Scan request...")
	_, err := client.Scan(ctx, input, file)
	if err == nil {
		log.Fatal("Should return AcceptedError")
	}

	ae := err.(*pangea.AcceptedError)
	fmt.Println("Accepted error received (as expected).")

	fmt.Println("Sleep some time until result should be ready.")
	// Wait until result should be ready
	time.Sleep(time.Duration(30 * time.Second))

	fmt.Println("File Scan poll result...")
	pr, err := client.PollResult(ctx, *ae)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("File Scan poll result success.")
	fmt.Println(pangea.Stringify(pr.Result))
}

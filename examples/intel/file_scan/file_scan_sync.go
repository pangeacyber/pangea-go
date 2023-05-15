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

	// Here we create a file that will give us a malicious result as example,
	// This should be your own file to scan
	file := strings.NewReader(EICAR)

	fmt.Println("File Scan request...")
	resp, err := client.Scan(ctx, input, file)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("File Scan success.")
	fmt.Println(pangea.Stringify(resp.Result))
}

# Pangea Go SDK

# Usage
```go
// embargo check is an example of how to use the check method
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/go-pangea/pangea"
	"github.com/pangeacyber/go-pangea/service/embargo"
)

func main() {
	token := os.Getenv("PANGEA_AUTH_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	configID := os.Getenv("EMBARGO_CONFIG_ID")
	if token == "" {
		log.Fatal("Configuration: No config ID present")
	}

	embargocli, err := embargo.New(&pangea.Config{
		Token: token,
		EndpointConfig: &pangea.EndpointConfig{
			Scheme: "https",
			CSP:    "aws",
		},
		CfgToken: configID,
	})
	if err != nil {
		log.Fatal("failed to create embargo client")
	}

	ctx := context.Background()
	input := &embargo.ISOCheckInput{
		ISOCode: pangea.String("CU"),
	}

	checkOutput, _, err := embargocli.ISOCheck(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(pangea.Stringify(checkOutput))
}
```

# Contributing

Currently, the setup scripts only have support for Mac/ZSH environments.
Future support is incoming.

To install our linters, simply run `./dev/setup_repo.sh`
These linters will run on every `git commit` operation.

# Generate SDK Documentation

## Overview

Throughout the SDK, there are go doc strings that serve as the source of our SDK docs.

The documentation pipeline here looks like:

1. Write doc strings throughout your go code. Please refer to existing doc strings as an example of what and how to document.
1. Make your pull request.
1. After the pull request is merged, go ahead and run the `autogen_docs.go` script to generate the JSON docs uses for rendering.
1. Copy the output from `autogen_docs.go` and overwrite the existing go_sdk.json file in the docs repo. File is located in platform/docs/sdk/go_sdk.json in the Pangea monorepo. Save this and make a merge request to update the Golang SDK docs in the Pangea monorepo.

## Running the autogen sdk doc script

From the root of the `go-pangea` repo run:
```sh
go run dev/autogen_docs.go
```
That will output the script in the terminal. If you're on a mac, you can do
```sh
go run dev/autogen_docs.go | pbcopy
```
to copy the output from the script into your clipboard.

<p>
  <br />
  <a href="https://pangea.cloud?utm_source=github&utm_medium=node-sdk" target="_blank" rel="noopener noreferrer">
    <img src="https://pangea-marketing.s3.us-west-2.amazonaws.com/pangea-color.svg" alt="Pangea Logo" height="40">
  </a>
  <br />
</p>

<p>
<br />

[![documentation](https://img.shields.io/badge/documentation-pangea-blue?style=for-the-badge&labelColor=551B76)](https://pangea.cloud/docs/sdk/go/)
[![Discord](https://img.shields.io/discord/1017567751818182786?color=%23551b76&label=Discord&logo=discord&logoColor=%23FFFFFF&style=for-the-badge)](https://discord.gg/z7yXhC7cQr)

<br />
</p>

# Pangea Go SDK

A Go SDK for integrating with Pangea Services.

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
	token := os.Getenv("PANGEA_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	configID := os.Getenv("EMBARGO_CONFIG_ID")
	if token == "" {
		log.Fatal("Configuration: No config ID present")
	}

	embargocli, err := embargo.New(&pangea.Config{
		Token: 		token,
		Domain: 	os.Getenv("PANGEA_DOMAIN"),
		Insecure: 	false,
		CfgToken: 	configID,
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

# Writing Docs

To maintain parity with documentation across all our SDKs, please follow this format when writing a doc comment for a *published* function or method. Published means the function ro method is listed as an endpoint in our API Reference docs.

Published Doc Example:
```
// Redact
//
// Redacts the content of a single text string.
//
// Example:
//
//  input := &redact.TextInput{
//  	Text: pangea.String("my phone number is 123-456-7890"),
//  }
//
//  redactOutput, _, err := redactcli.Redact(ctx, input)
//
```

Example breakdown:
```
// Redact <-- Displayed as the Summary/Heading field in docs
//
// Redacts the content of a single text string. <-- Displayed as the Description field in docs
//
// Example: <-- All lines below this are used as the code snippet field in docs
//
//  input := &redact.TextInput{
//  	Text: pangea.String("my phone number is 123-456-7890"),
//  }
//
//  redactOutput, _, err := redactcli.Redact(ctx, input)
//
```

# Generate SDK Documentation

## Overview

Throughout the SDK, there are go doc strings that serve as the source of our SDK docs.

The documentation pipeline here looks like:

1. Write doc strings throughout your go code. Please refer to existing doc strings as an example of what and how to document.
1. Make your pull request.
1. After the pull request is merged, go ahead and run the autogen docs script to generate the JSON docs uses for rendering.
1. Copy the output from autogen docs and overwrite the existing go_sdk.json file in the docs repo. File is located in platform/docs/sdk/go_sdk.json in the Pangea monorepo. Save this and make a merge request to update the Golang SDK docs in the Pangea monorepo.

## Running the autogen sdk doc script

From the root of the `go-pangea` repo run:
```sh
go run dev/autogendoc/main.go
```
That will output the script in the terminal. If you're on a mac, you can do
```sh
go run dev/autogendoc/main.go | pbcopy
```
to copy the output from the script into your clipboard.

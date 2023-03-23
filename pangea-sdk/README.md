<p>
  <br />
  <a href="https://pangea.cloud?utm_source=github&utm_medium=node-sdk" target="_blank" rel="noopener noreferrer">
    <img src="https://pangea-marketing.s3.us-west-2.amazonaws.com/pangea-color.svg" alt="Pangea Logo" height="40" />
  </a>
  <br />
</p>

<p>
<br />

[![documentation](https://img.shields.io/badge/documentation-pangea-blue?style=for-the-badge&labelColor=551B76)](https://pangea.cloud/docs/sdk/go/)
[![Slack](https://img.shields.io/badge/Slack-4A154B?style=for-the-badge&logo=slack&logoColor=white)](https://pangea.cloud/join-slack/)

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

	"github.com/pangeacyber/pangea-go/pangea-sdk/pangea pangea-sdk/v1.0.0"
	"github.com/pangeacyber/pangea-go/pangea-sdk/service/embargo pangea-sdk/v1.0.0"
)

func main() {
	token := os.Getenv("PANGEA_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	embargocli, err := embargo.New(&pangea.Config{
		Token: 		token,
		Domain: 	os.Getenv("PANGEA_DOMAIN"),
		Insecure: 	false,
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

To maintain parity with documentation across all our SDKs, please follow this format when writing a doc comment for a *published* function or method. Published means the function or method is listed as an endpoint in our API Reference docs.

Published Doc Example:
```
// @summary Redact
//
// @description Redacts the content of a single text string.
//
// @example
//
//	input := &redact.TextInput{
//  		Text: pangea.String("my phone number is 123-456-7890"),
//  }
//
//  redactOutput, _, err := redactcli.Redact(ctx, input)
//
```

Example breakdown:
```
// @summary Redact <-- Displayed as the Summary/Heading field in docs
//
// @description Redacts the content of a single text string. <-- Displayed as the Description field in docs
//
// @example <-- All lines below this are used as the code snippet field in docs
//
//  input := &redact.TextInput{
//  	Text: pangea.String("my phone number is 123-456-7890"),
//  }
//
//  redactOutput, _, err := redactcli.Redact(ctx, input)
//
```

Example with deprecation message:
```
// @summary Lookup a domain
//
// @description Lookup an internet domain to retrieve reputation data.
//
// @deprecated Use Reputation instead.
//
// @example
//
//	input := &domain_intel.DomainLookupInput{
//		Domain: "737updatesboeing.com",
//		Raw: true,
//		Verbose: true,
//		Provider: "domaintools",
//	}
//
//	checkResponse, err := domainintel.Lookup(ctx, input)
```

# autogendoc

A script for generating Pangea SDK Documentation.

Currently, this is an internal process handled by anyone at Pangea.

## Overview

Throughout the SDK, there are go doc strings that serve as the source of our SDK docs.

The documentation pipeline here looks like:

1. Write doc strings throughout your go code. Please refer to existing doc strings as an example of what and how to document.
1. Make your pull request.
1. After the pull request is merged, go ahead and run the autogen docs script to generate the JSON docs uses for rendering.
1. Copy the output from autogen docs and overwrite the existing go_sdk.json file in the Pangea documentation repo. The file is located in platform/docs/sdk/go_sdk.json in the Pangea monorepo. Save this and make a merge request to update the Golang SDK docs in the Pangea monorepo.

## Running the autogen sdk doc script

From the directory `pangea-go/pangea-sdk/v2` run:
```sh
go run ./dev/autogendoc/main.go
```
That will output the script in the terminal. If you're on a mac, you can do
```sh
go run ./dev/autogendoc/main.go | pbcopy
```
to copy the output from the script into your clipboard.

## Writing Docs

To maintain parity with documentation across all our SDKs, please follow this format when writing a doc comment for a *published* function or method. Published means the function or method is listed as an endpoint in our API Reference docs.

Published Doc Example:
```
// @summary Redact
//
// @description Redacts the content of a single text string.
//
// @operationId redact_post_v1_redact
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
// @operationId redact_post_v1_redact <-- The operationId from the OpenAPI spec
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

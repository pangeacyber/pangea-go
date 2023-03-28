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

From the root of the `go-pangea` repo run:
```sh
go run dev/autogendoc/main.go
```
That will output the script in the terminal. If you're on a mac, you can do
```sh
go run dev/autogendoc/main.go | pbcopy
```
to copy the output from the script into your clipboard.
# Pangea Go SDK examples

This directory contains full examples on how to use the Pangea Go SDK.

## Setup

Each example requires certain environment variables to be set. `PANGEA_URL_TEMPLATE`
must be set to a Pangea URL template (e.g. `https://{SERVICE_NAME}.aws.us.pangea.cloud`). Then, a token
variable must be set as well. This is typically in the format of
`PANGEA_{SERVICE_NAME}_TOKEN` (so: `PANGEA_REDACT_TOKEN` for the Redact service,
`PANGEA_VAULT_TOKEN` for the Vault service, etc.). Finally, note that some
examples require additional variables, so check out the example's source code
and look out for what environment variables it loads at the beginning.

## Run

To run an example, navigate to its directory and use `go run`:

```bash
$ cd examples/intel
$ go run url/reputation.go
```

# Pangea Go SDK examples

This is a quick example about how you use Pangea Go SDK, set up and run it.

## Setup

Set up environment variables ([Instructions](https://pangea.cloud/docs/getting-started/integrate/#set-environment-variables)) `PANGEA_AUDIT_TOKEN` and `PANGEA_DOMAIN` with your project token configured on Pangea User Console (token should have access to Audit service [Instructions](https://pangea.cloud/docs/getting-started/configure-services/#configure-a-pangea-service)) and with your pangea domain.

## Run

To run examples, move to service folder:
```
cd examples/intel
go mod tidy
```

and from service folder run:

```
go run url/reputation.go
```

<a href="https://pangea.cloud?utm_source=github&utm_medium=go-sdk" target="_blank" rel="noopener noreferrer">
  <img src="https://pangea-marketing.s3.us-west-2.amazonaws.com/pangea-color.svg" alt="Pangea Logo" height="40" />
</a>

<br />

[![Documentation](https://img.shields.io/badge/documentation-pangea-blue?style=for-the-badge&labelColor=551B76)][Documentation]

# Pangea Go SDK

A Go SDK for integrating with Pangea services. Supports Go v1.23 and above.

## Installation

#### GA releases

```bash
$ go get github.com/pangeacyber/pangea-go/pangea-sdk/v5
```

<a name="beta-releases"></a>

#### Beta releases

Pre-release versions may be available with the `beta` denotation in the version
number. These releases serve to preview Beta and Early Access services and APIs.
Per Semantic Versioning, they are considered unstable and do not carry the same
compatibility guarantees as stable releases.
[Beta changelog](https://github.com/pangeacyber/pangea-go/blob/beta/CHANGELOG.md).

```bash
$ go get github.com/pangeacyber/pangea-go/pangea-sdk/v5@v5.3.0-beta.2
```

## Usage

- [Documentation][]
- [GA Examples][]
- [Beta Examples][]

Set up the SDK in your project in 3 steps:

1. Pick your service. Full list of services available [here](https://pangea.cloud).
2. Initialize your client with your `Token` and `Domain`
3. Use your client to call the service's endpoints

Let's walk through an example using:

- Service: [Secure Audit Log](https://pangea.cloud/services/secure-audit-log/)
- Endpoint: [/v1/log](https://pangea.cloud/docs/api/audit#/v1/log-post)

We need two things to initialize your client: a `Token` and `Domain`. These can
be found on the service overview page. For the Secure Audit Log service, go to
<https://console.pangea.cloud/service/audit> and take a look at the
"Configuration Details" box where it has "Default Token" and "Domain" listed.

Go ahead and set the token and domain as environment variables in our terminal.

```bash
$ export PANGEA_AUDIT_TOKEN=pts_tokenvaluehere
```

```bash
$ export PANGEA_DOMAIN=aws.us.pangea.cloud
```

Now let's add the SDK to our code.

Import statements:

```go
import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v5/service/audit"
)
```

Initialize your client:

```go
// Client configuration.
config, err := pangea.NewConfig(
	option.WithToken(os.Getenv("PANGEA_AUDIT_TOKEN")),
	option.WithDomain(os.Getenv("PANGEA_DOMAIN")),
)
if err != nil {
	// Handle error.
}

// Initialize the Secure Audit Log client.
auditcli, err := audit.New(config)
if err != nil {
	// Handle error.
}
```

IMPORTANT! Never hardcode your token. [Use environment variables](https://gobyexample.com/environment-variables) to [avoid committing secrets in your codebase](https://www.thisdot.co/blog/a-guide-to-keeping-secrets-out-of-git-repositories/).

Make a call to the `/v1/log` endpoint using the SDK

```go
// Set up our parameters
ctx := context.Background()
event := &audit.StandardEvent{
	Message: "Hello, World!",
}

// Call the /v1/log endpoint
resp, err := auditcli.Log(ctx, event, true)
if err != nil {
	log.Fatal(err)
}

// Print the response
event := (resp.Result.EventEnvelope.Event).(*audit.StandardEvent)
fmt.Printf("Logged event: %s", pangea.Stringify(e))
```

Full code for the above example available in [the examples directory](https://github.com/pangeacyber/pangea-go/blob/main/examples/audit/log_standard_schema.go).

<a name="configuration"></a>

## Configuration

The SDK supports the following configuration options via `pangea.NewConfig`:

- `option.WithBaseURLTemplate()` — Template for constructing the base URL for
  API requests. The placeholder `{SERVICE_NAME}` will be replaced with the
  service name slug. This is a more powerful version of Domain that allows for
  setting more than just the host of the API server. Defaults to
  `https://{SERVICE_NAME}.aws.us.pangea.cloud`.
- `option.WithDomain()` — Base domain for API requests. This is a weaker version
  of `BaseURLTemplate` that only allows for setting the host of the API server.
  Use `BaseURLTemplate` for more control over the URL, such as setting
  service-specific paths. Defaults to `aws.us.pangea.cloud`.
- `option.WithHTTPClient()` - HTTP client to be used by the client. Defaults to
  `http.DefaultClient`.
- `option.WithLogger()` - Logger to be used by the client.
- `option.WithMaxRetries()` - Maximum number of retries to attempt. When set to
  0, the client only makes one request. By default, the client retries two
  times.
- `option.WithPollResultTimeout()` - Timeout used to poll results after
  HTTP/202.
- `option.WithQueuedRetryEnabled()` - Whether or not the client should retry
  queued requests.
- `option.WithToken()` - API token used to authenticate requests.

<a name="asynchronous-responses"></a>

## Asynchronous responses

Any response from Pangea may become [an asynchronous one][Asynchronous API Responses]
if the result takes too long to be ready. These HTTP/202 responses will manifest
themselves as `AcceptedError`s in the SDK. Use a type assertion to check for
this error type.

```go
resp, err := auditcli.Log(ctx, event, true)
if err != nil {
	acceptedError, isAcceptedError := err.(*pangea.AcceptedError)
	if isAcceptedError {
		// Result is not ready yet. One may poll for it by using
		// PollResultByError() or PollResultByID() + acceptedError.RequestID.
	} else {
		// Handle other error types.
	}
}
```

[Documentation]: https://pangea.cloud/docs/sdk/go/
[GA Examples]: https://github.com/pangeacyber/pangea-go/tree/main/examples
[Beta Examples]: https://github.com/pangeacyber/pangea-go/tree/beta/examples
[Pangea Console]: https://console.pangea.cloud/
[Secure Audit Log]: https://pangea.cloud/docs/audit
[Asynchronous API Responses]: https://pangea.cloud/docs/api/async

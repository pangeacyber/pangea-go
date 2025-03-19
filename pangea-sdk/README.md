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
$ go get github.com/pangeacyber/pangea-go/pangea-sdk/v4
```

<a name="beta-releases"></a>

#### Beta releases

Pre-release versions may be available with the `beta` denotation in the version
number. These releases serve to preview Beta and Early Access services and APIs.
Per Semantic Versioning, they are considered unstable and do not carry the same
compatibility guarantees as stable releases.
[Beta changelog](https://github.com/pangeacyber/pangea-go/blob/beta/CHANGELOG.md).

```bash
$ go get github.com/pangeacyber/pangea-go/pangea-sdk/v4@v4.2.0-beta.1
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
- Endpoint: [/v1/log](https://pangea.cloud/docs/api/audit#log-an-entry)

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

	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v4/service/audit"
)
```

Initialize your client:

```go
// Initialize the Secure Audit Log client.
auditcli, err := audit.New(&pangea.Config{
	Token: os.Getenv("PANGEA_AUDIT_TOKEN"), // NEVER hardcode your token here, always use env vars
	BaseURLTemplate: os.Getenv("PANGEA_URL_TEMPLATE"),
})
if err != nil {
	log.Fatal("failed to create Audit client")
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

   [Documentation]: https://pangea.cloud/docs/sdk/go/
   [GA Examples]: https://github.com/pangeacyber/pangea-go/tree/main/examples
   [Beta Examples]: https://github.com/pangeacyber/pangea-go/tree/beta/examples
   [Pangea Console]: https://console.pangea.cloud/
   [Secure Audit Log]: https://pangea.cloud/docs/audit

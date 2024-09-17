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
[![Discourse](https://img.shields.io/badge/Discourse-4A154B?style=for-the-badge&logo=discourse&logoColor=white)](https://l.pangea.cloud/Jd4wlGs)

<br />
</p>

# Pangea Go SDK

A Go SDK for integrating with Pangea Services.

## Installation

To add Pangea Go SDK to your project, will need to run next command on your project root directory where you should have `go.mod` file:

```bash
$ go get github.com/pangeacyber/pangea-go/pangea-sdk/v2
```

## Usage

### Examples

For the best example of how to set up and use the SDK in your code, look at the [/examples](https://github.com/pangeacyber/pangea-go/tree/main/examples) folder in this repository. There you will find basic samples apps for each services supported on this SDK. Each service folder has a README.md with instructions to install, setup and run.

## Getting started

Set up the SDK in your project in 3 steps:

1. Pick your service. Full list of services available [here](https://pangea.cloud).
1. Initialize your client with your `Token` and `Domain`
1. Use your client to call the service's endpoints

Let's walkthrough an example using:

- Service: [Secure Audit Log](https://pangea.cloud/services/secure-audit-log/)
- Endpoint: [/v1/log](https://pangea.cloud/docs/api/audit#log-an-entry)

We need two things to initialize your client -- a `Token` and `Domain`. These can be found on the service overview page. For the Secure Audit Log service, go to [https://console.pangea.cloud/service/audit](https://console.pangea.cloud/service/audit) and take a look at the "Configuration Details" box where it has "Default Token" and "Domain" listed.

Go ahead and set the token and domain as environment variables in our terminal.

```sh
$ export PANGEA_AUDIT_TOKEN=pts_mytokenvaluehere
```

```sh
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

  "github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
  "github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/audit"
)
```

Initialize your client:

```go
// Initialize the Secure Audit Log client
auditcli, err := audit.New(&pangea.Config{
  Token: os.Getenv("PANGEA_AUDIT_TOKEN"), // NEVER hardcode your token here, always use env vars
  Domain: os.Getenv("PANGEA_DOMAIN"),
})
if err != nil {
  log.Fatal("failed to create audit client")
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

## Reporting issues and new features

If faced some issue using or testing this SDK or a new feature request feel free to open an issue [clicking here](https://github.com/pangeacyber/pangea-go/issues).
We would need you to provide some basic information like what SDK's version you are using, stack trace if you got it, framework used, and steps to reproduce the issue.
Also feel free to contact [Pangea support](mailto:support@pangea.cloud) by email or send us a message on [Discourse](https://l.pangea.cloud/Jd4wlGs)

## Contributing

Currently, the setup scripts only have support for Mac/ZSH environments.
Future support is incoming.

To install our linters, simply run `./dev/setup_repo.sh`
These linters will run on every `git commit` operation.

### Send a PR

If you would like to [send a PR](https://github.com/pangeacyber/pangea-go/pulls) including a new feature or fixing a bug in code or an error in documents we will really appreciate it and after review and approval you will be included in our [contributors list](https://github.com/pangeacyber/pangea-go/blob/main/pangea-sdk/v2/CONTRIBUTING.md)

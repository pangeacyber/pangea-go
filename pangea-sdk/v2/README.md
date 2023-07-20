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

## Installation

To add Pangea Go SDK to your project, will need to run next command on your project root directory where you should have `go.mod` file:

```
go get github.com/pangeacyber/pangea-go/pangea-sdk
```


## Usage

For samples apps look at [/examples](https://github.com/pangeacyber/pangea-go/tree/main/examples) folder in this repository. There you will find basic samples apps for each services supported on this SDK. Each service folder has a README.md with instructions to install, setup and run.


## Reporting issues and new features

If faced some issue using or testing this SDK or a new feature request feel free to open an issue [clicking here](https://github.com/pangeacyber/pangea-go/issues).
We would need you to provide some basic information like what SDK's version you are using, stack trace if you got it, framework used, and steps to reproduce the issue.
Also feel free to contact [Pangea support](mailto:support@pangea.cloud) by email or send us a message on [Slack](https://pangea.cloud/join-slack/)


## Contributing

Currently, the setup scripts only have support for Mac/ZSH environments.
Future support is incoming.

To install our linters, simply run `./dev/setup_repo.sh`
These linters will run on every `git commit` operation.

### Send a PR

If you would like to [send a PR](https://github.com/pangeacyber/pangea-go/pulls) including a new feature or fixing a bug in code or an error in documents we will really appreciate it and after review and approval you will be included in our [contributors list](./CONTRIBUTING.md)

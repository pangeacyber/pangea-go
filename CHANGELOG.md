# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.2.0] - 2023-09-05

# Added

- Redact rulesets field support 
- FileScan service support


## [2.1.0] - 2023-07-14

# Added

- Vault /folder/create endpoint support


## [2.0.0] - 2023-07-06

# Added

- Logger support on each service
- Service methods to fetch async request's results

# Changed

- Audit service now allow user to setup CustomEvent format
- Audit.Log() now receive a IEvent interface instead of Event to log
- Make service structs private to force using its interface
- Rename Request/Result structs to user Request/Result postfix instead of Input/Output
- Update multiple fields to not be pointers unnecesarily
- Vault field RotationState is now a ItemVersionState

# Deleted

- All Intel deprecated methods


## [1.10.0] - 2023-06-26

# Added
- Multiconfig support
- Instructions to setup token and domain in examples


## [1.9.1] - 2023-06-09
# Added

- Defang examples
- Intel User breached password full example
- Intel IP /domain, /vpn and /proxy examples

# Changed

- UserBreachedPasswordResult now has maps instead of just interface


## [1.9.0] - 2023-05-25

# Added

- New algorithm support in Vault Service
- Algorithm field support in Audit Service
- Cymru IP Intel provider examples
- Support full url as domain in config for local use


## [1.8.0] - 2023-05-05

### Added

- AuthN support
- Docs to multiple services

### Fixed

- Redact Score field type
- Readme example links to SDK examples page


## [1.7.0] - 2023-04-10

### Added

- Audit-Vault signing integration support
- Intel User support
- Redact Service return_result field support
- Set custom user agent by config


## [1.6.0] - 2023-03-27

### Added

- Algorithm support in Vault Service

### Changed

- Algorithm name in Vault Service

## [1.5.0] - 2023-03-17

### Added

- Vault service support
- LICENSE

### Changed

- Update services examples
- Improve docs

## [1.4.0] - 2023-03-01

### Added

- IP service add /geolocate, /vpn, /domain and /proxy endpoints support

## [1.3.0] - 2023-02-28

### Added

- Tenant ID support in Audit Service

## [1.2.0] - 2023-02-03

### Added

- Rules parameter support in Redact service


## [1.1.2] - 2023-01-27

### Changed

- Intel Domain, URL and File add reputation endpoint that will replace lookup endpoint


## [1.1.1] - 2023-01-25

### Changed

- Intel IP add reputation endpoint that will replace lookup endpoint
- Change User-Agent format

### Added

- Count field in redact result


## [1.1.0] - 2023-01-05

### Added

- This CHANGELOG
- Intel add IP and URL services lookup endpoint

### Fixed

- Fix old references to deprecated SDK go repository

### Changed

- Unify token env var name on integration tests and sample apps


## [1.0.0] - 2022-11-29

### Added

- Audit client
- Embargo client
- File Intel client
- Domain Intel client
- Redact client


[unreleased]: https://github.com/pangeacyber/pangea-go/compare/pangea-sdk/v2.2.0...main
[2.2.0]: https://github.com/pangeacyber/pangea-go/compare/pangea-sdk/v2.1.0...v2.2.0
[2.1.0]: https://github.com/pangeacyber/pangea-go/compare/pangea-sdk/v2.0.0...v2.1.0
[2.0.0]: https://github.com/pangeacyber/pangea-go/compare/pangea-sdk/v1.10.0...v2.0.0
[1.10.0]: https://github.com/pangeacyber/pangea-go/compare/pangea-sdk/v1.9.1...v1.10.0
[1.9.1]: https://github.com/pangeacyber/pangea-go/compare/pangea-sdk/v1.9.0...v1.9.1
[1.9.0]: https://github.com/pangeacyber/pangea-go/compare/pangea-sdk/v1.8.0...v1.9.0
[1.8.0]: https://github.com/pangeacyber/pangea-go/compare/pangea-sdk/v1.7.0...v1.8.0
[1.7.0]: https://github.com/pangeacyber/pangea-go/compare/pangea-sdk/v1.6.0...v1.7.0
[1.6.0]: https://github.com/pangeacyber/pangea-go/compare/pangea-sdk/v1.5.0...v1.6.0
[1.5.0]: https://github.com/pangeacyber/pangea-go/compare/pangea-sdk/v1.4.0...v1.5.0
[1.4.0]: https://github.com/pangeacyber/pangea-go/compare/pangea-sdk/v1.3.0...v1.4.0
[1.3.0]: https://github.com/pangeacyber/pangea-go/compare/pangea-sdk/v1.2.0...v1.3.0
[1.2.0]: https://github.com/pangeacyber/pangea-go/compare/pangea-sdk/v1.1.2...v1.2.0
[1.1.2]: https://github.com/pangeacyber/pangea-go/compare/pangea-sdk/v1.1.1...v1.1.2
[1.1.1]: https://github.com/pangeacyber/pangea-go/compare/pangea-sdk/v1.1.0...v1.1.1
[1.1.0]: https://github.com/pangeacyber/pangea-go/compare/pangea-sdk/v1.0.0...v1.1.0
[1.0.0]: https://github.com/pangeacyber/pangea-go/releases/tag/pangea-sdk%2Fv1.0.0

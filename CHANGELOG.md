# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

### Added

- AI Guard: `GuardTextWithRelevantContent()` method which allows for sending
  only relevant messages to AI Guard.

### Changed

- Minimum supported Go version is now v1.24.

## 5.4.0 - 2025-08-27

### Changed

- Retries now include a request header `X-Pangea-Retried-Request-Ids` to track
  the request IDs of the retries.
- Configuring retries has been simplified to configuring only the number of
  retries via `option.WithMaxRetries(...)`.

## 5.3.0 - 2025-06-24

### Added

- Added `pangea.NewConfig()` as a constructor for the `Config` struct that will
  set defaults such as enabling failed request retrying (`Retry`).
- Added `option.WithAdditionalHeaders()` for setting additional request headers.

### Changed

- If a `Host` header is passed to `option.WithAdditionalHeaders()`, it will be
  promoted to net/http's `Request.Host`.

## 5.2.0 - 2025-04-25

### Added

- AuthZ: `ExpiresAt` to tuples.
- AuthN: groups.

### Fixed

- Secure Audit Log: "diabled" typo in `NewFilterUserList()`.

## 5.1.0 - 2025-04-15

### Added

- AI Guard: `ignore_recipe` in detector overrides.
- AI Guard: new detectors: Competitors, Gibberish, Negative Sentiment, and
  Self-Harm and Violence.

### Changed

- Disable logger by default
- AI Guard: updated the Topic detector's overrides.

## 5.0.0 - 2025-03-28

### Added

- `BaseURLTemplate` has been added to `Config` to allow for greater control over
  the complete API URL. This option may be a full URL with the optional
  `{SERVICE_NAME}` placeholder, which will be replaced by the slug of the
  respective service name. This supersedes `Environment` and `Insecure`.
- Redact: `Unredact()` method on service interface.
- Redact: `FPEContext` on `RedactResult` and `StructuredResult`.
- AI Guard: topic detector.

### Changed

- Minimum supported Go version is now v1.23

### Removed

- `Environment` and `Insecure` have been removed from `Config`, as the
  functionality they provided can now be accomplished via `BaseURLTemplate`
  (see above).
- AI Guard: `LlmInfo` and `LlmInput`.
- PangeaConfig: `ConfigID` and `Enviroment` field.

## 4.4.0 - 2025-02-16

### Added

- AI Guard and Prompt Guard services.

### Changed

- Deprecated `Config.Enviroment` (typo) in favor of `Config.Environment`.
- Clarified what `Config.Environment` affects.

### Removed

- CDR and PDF support in Sanitize.

## 4.3.0 - 2025-01-13

### Added

- `file_ttl` support to Secure Share.

## 4.2.0 - 2024-12-18

### Added

- More documentation to AuthZ and Secure Share filter lists.
- Support for `cursor` field on `v1/user/breached` of `user-intel` service.
- Millisecond time format support
- `vault_parameters` and `llm_request` fields support on Redact service.
- Support for `severity` field on `v1/user/breached` and `v2/user/breached` of `user-intel` service.
- `/v1/breach` endpoint support on `user-intel` service.

### Fixed

- Secure Share: share link list filter now uses "target_id" instead of "target".

## 4.1.0 - 2024-10-16

### Added

- Secure Share support.
- Multiple bucket ID support to Share.
- `attributes` field in `/list-resources` and `/list-subjects` endpoint
- Sanitize service support
- Secure Share support.
- Multiple bucket ID support to Share.
- `MetadataProtected` and `TagsProtected` support to Share `ItemData`
- `Password` and `PasswordAlgorithm` support to Share
- Filter fields to `FilterList` on Share service
- `Objects` field to Share `GetArchiveResult`
- `Title` and `Message` to Share `ShareCreateLinkItem` 

## 4.0.0 - 2024-10-15

### Added

- Vault KEM export support.

### Changed

- Vault v2 APIs support.
- Minimum supported Go version is now v1.22.

## 3.12.0 - 2024-10-15

### Added

- Detect-only Redact for Sanitize.
- Support for `domains` field in `v2/user/breached` endpoint in User Intel service.

## 3.11.0 - 2024-09-25

### Added 

- `attributes` field in `/list-resources` and `/list-subjects` endpoint.
- Sanitize service support.

## [3.10.0] - 2024-07-19

### Added

- Improvements in verification of Audit consistency proofs
- Doc example for unredact.
- Vault `/export` support
- `exportable` field support in Vault `/key/store` and `/key/generate`
- AuthN user password expiration support.
- `"state"` and other new properties to `Authenticator`.

### Changed

- `Enable` in `Authenticator` has been renamed to `Enabled`. The previous name
  did not match the name used in the API's response schema and JSON
  deserialization was not set up correctly, so `Enable` was unusable anyways.

## [3.9.0] - 2024-06-07

### Added

- `fpe_context` field in Audit search events
- `return_context` support in Audit `/search`, `/results` and `/download` endpoints
- Redact `/unredact` endpoint support
- `redaction_method_overrides` field support in `/redact` and `redact_structured` endpoints
- AuthN usernames support.
- Support for format-preserving encryption.

### Removed

- Beta tags from AuthZ.

## [3.8.0] - 2024-05-10

Note that Sanitize and Secure Share did not make it into this release.

### Added

- Support for Secure Audit Log's log stream API.
- Support for Secure Audit Log's export API.
- AuthZ service support.

### Changed

- Replaced usage of the deprecated io/ioutil package.
- Audit /download_results endpoint support

### Fixed

- All enums now have consistent types across their values.
- Put to presigned url. It should just put file in raw, not in form format.


## [3.7.0] - 2024-02-26

### Added 

- Vault service. Post quantum signing algorithms support

### Changed

- Rewrote `README.md`.

## [3.6.0] - 2024-01-16

### Added

- Vault encrypt structured support.

## [3.5.0] - 2023-12-18

### Added

- File Intel /v2/reputation support
- IP Intel /v2/reputation, /v2/domain, /v2/proxy, v2/vpn and /v2/geolocate support
- URL Intel /v2/reputation support
- Domain Intel /v2/reputation support
- User Intel /v2/user/breached and /v2/password/breached support


## [3.4.0] - 2023-12-07

### Changed

- 202 result format

### Removed

- accepted_status in 202 result

### Added

- put_url, post_url, post_form_data fields in 202 result


## [3.3.0] - 2023-11-28

### Added

- Authn unlock user support
- Redact multiconfig support
- File Scan post-url and put-url support


## [3.2.0] - 2023-11-15

### Added

- Support for audit /v2/log and /v2/log_async endpoints


## [3.1.0] - 2023-11-09

### Added

- Presigned URL upload support on FileScan service
- Folder settings support in Vault service

## [3.0.0] - 2023-10-23

### Added

- AuthN v2 support

### Removed

- AuthN v1 support


## [2.3.0] - 2023-09-26

### Added

- FileScan Reversinglabs provider example
- Domain WhoIs endpoint support
- AuthN Filters support

### Changed

- Deprecated config_id in PangeaConfig. Now is set in service initialization.

### Fixed

- HashType supported in File Intel

## [2.2.0] - 2023-09-05

### Added

- Redact rulesets field support
- FileScan service support


## [2.1.0] - 2023-07-14

### Added

- Vault /folder/create endpoint support


## [2.0.0] - 2023-07-06

### Added

- Logger support on each service
- Service methods to fetch async request's results

### Changed

- Audit service now allow user to setup CustomEvent format
- Audit.Log() now receive a IEvent interface instead of Event to log
- Make service structs private to force using its interface
- Rename Request/Result structs to user Request/Result postfix instead of Input/Output
- Update multiple fields to not be pointers unnecessarily
- Vault field RotationState is now a ItemVersionState

### Removed

- All Intel deprecated methods


## [1.10.0] - 2023-06-26

### Added
- Multiconfig support
- Instructions to setup token and domain in examples


## [1.9.1] - 2023-06-09
### Added

- Defang examples
- Intel User breached password full example
- Intel IP /domain, /vpn and /proxy examples

### Changed

- UserBreachedPasswordResult now has maps instead of just interface


## [1.9.0] - 2023-05-25

### Added

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


[unreleased]: https://github.com/pangeacyber/pangea-go/compare/pangea-sdk/v3.7.0...main
[3.7.0]: https://github.com/pangeacyber/pangea-go/compare/pangea-sdk/v3.6.0...v3.7.0
[3.6.0]: https://github.com/pangeacyber/pangea-go/compare/pangea-sdk/v3.5.0...v3.6.0
[3.5.0]: https://github.com/pangeacyber/pangea-go/compare/pangea-sdk/v3.4.0...v3.5.0
[3.4.0]: https://github.com/pangeacyber/pangea-go/compare/pangea-sdk/v3.3.0...v3.4.0
[3.3.0]: https://github.com/pangeacyber/pangea-go/compare/pangea-sdk/v3.2.0...v3.3.0
[3.2.0]: https://github.com/pangeacyber/pangea-go/compare/pangea-sdk/v3.1.0...v3.2.0
[3.1.0]: https://github.com/pangeacyber/pangea-go/compare/pangea-sdk/v3.0.0...v3.1.0
[3.0.0]: https://github.com/pangeacyber/pangea-go/compare/pangea-sdk/v2.3.0...v3.0.0
[2.3.0]: https://github.com/pangeacyber/pangea-go/compare/pangea-sdk/v2.2.0...v2.3.0
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

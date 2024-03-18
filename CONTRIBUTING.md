# Contributing

Currently, the setup scripts only have support for Mac/ZSH environments.
Future support is incoming.

To install our linters, simply run `./dev/setup_repo.sh`
These linters will run on every `git commit` operation.

## Publishing

Publishing pangea-sdk to pkg.go.dev is handled via a private GitLab CI pipeline.
This pipeline is triggered when a Git tag is pushed to the repository. Git tags
should be formatted as `pangea-sdk/vX.Y.Z`, where `vX.Y.Z` is the version number
to publish.

1. Update the `version` constant in `pangea-sdk/v3/pangea/pangea.go`.
2. Update the release notes in `CHANGELOG.md`.
3. Author a commit with these changes and land it on `main`.
4. `git tag -m pangea-sdk/vX.Y.Z pangea-sdk/vX.Y.Z 0000000`. Replace `vX.Y.Z`
  with the new version number and `0000000` with the commit SHA from the
  previous step.
5. `git push --tags origin main`.

From here the GitLab CI pipeline will pick up the pushed Git tag and publish
the package to pkg.go.dev.

## Contributors

- Andrés Tournour (andres.tournour@gmail.com). Code.
- Glenn Gallien (glenn.gallien@pangea.cloud). Code and docs.
- David Wayman (david.wayman@pangea.cloud). Code and docs.

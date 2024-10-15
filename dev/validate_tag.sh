#!/usr/bin/env bash

set -e

if [ $# -lt 1 ]; then
    echo "usage: validate_tag.sh <git_tag>"
    exit 1
fi

GIT_TAG=$1

if [[ ! $GIT_TAG == "pangea-sdk/v"* ]]; then
    echo "Git tag must begin with a 'pangea-sdk/v'."
    exit 1
fi

PACKAGE_NAME=$(echo "$GIT_TAG" | cut -d "/" -f 1)
VERSION=$(echo "$GIT_TAG" | cut -d "/" -f 2)

if [[ ! "$VERSION" == *"v"* ]]; then
    echo "Git tag must contain a version number that's prefixed with 'v'."
    exit 1
fi

# Trim the 'v'.
VERSION="${VERSION:1}"

# Move to repo root.
PARENT_PATH=$(cd "$(dirname "${BASH_SOURCE[0]}")"; pwd -P)
pushd "$PARENT_PATH/.."

GO_CONST_VERSION=$(grep -Eo "version.+=.+" pangea-sdk/pangea/pangea.go | head -1)

if [[ ! "$GO_CONST_VERSION" == *"$VERSION"* ]]; then
    echo "Git tag version '$VERSION' does not match Go constant version '$GO_CONST_VERSION'."
    exit 1
fi

popd

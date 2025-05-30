#!/usr/bin/env bash

set -e

cd -- "$(dirname -- "$0")/.."

npx start-server-and-test --expect 404 \
  "npx @stoplight/prism-cli mock -d --json-schema-faker-fillProperties=false dev/specs/ai-guard.openapi.json" \
  4010 \
  "go test ./pangea-sdk/service/ai_guard/mock_test.go"

#!/usr/bin/env bash

set -e

cd -- "$(dirname -- "$0")/.."

pnpm exec start-server-and-test --expect 404 \
  "pnpm exec prism mock -d --json-schema-faker-fillProperties=false dev/specs/ai-guard.openapi.json" \
  4010 \
  "go test ./pangea-sdk/service/ai_guard/mock_test.go"

pnpm exec start-server-and-test --expect 404 \
  "pnpm exec prism mock -d --json-schema-faker-fillProperties=false dev/specs/audit.openapi.json" \
  4010 \
  "go test ./pangea-sdk/service/audit/mock_test.go"

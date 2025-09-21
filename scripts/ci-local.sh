#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
IMAGE_NAME="hubless-ci-local"

docker build -f "$ROOT_DIR/.ci/Dockerfile" -t "$IMAGE_NAME" "$ROOT_DIR"

docker run --rm \
  -v "$ROOT_DIR:/workspace" \
  -w /workspace \
  "$IMAGE_NAME"

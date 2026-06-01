#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
IMAGE_NAME="hubless-release-test"
VERSION="${VERSION:-0.0.1}"

cd "$ROOT_DIR"

docker build --pull -f Dockerfile.release-test -t "$IMAGE_NAME" .

docker run --rm "$IMAGE_NAME" /bin/bash -lc "\
  export PATH=/usr/local/go/bin:/go/bin:\$PATH && \
  set -x && \
  cd /app && \
  git remote -v && \
  go run ./cmd/release --version $VERSION --dry-run --skip-checks && \
  go run ./cmd/release --version $VERSION --skip-checks && \
  git tag --list && \
  ls docs/reference | grep release-notes.md\
"

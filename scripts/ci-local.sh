#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
IMAGE_NAME="hubless-ci-local"
COMMON_GIT_DIR="$(git -C "$ROOT_DIR" rev-parse --path-format=absolute --git-common-dir)"
GIT_MOUNTS=()

case "$COMMON_GIT_DIR" in
  "$ROOT_DIR"/.git | "$ROOT_DIR"/.git/*) ;;
  *) GIT_MOUNTS=(-v "$COMMON_GIT_DIR:$COMMON_GIT_DIR") ;;
esac

docker build -f "$ROOT_DIR/.ci/Dockerfile" -t "$IMAGE_NAME" "$ROOT_DIR"

docker run --rm \
  -v "$ROOT_DIR:/workspace" \
  "${GIT_MOUNTS[@]}" \
  -w /workspace \
  "$IMAGE_NAME"

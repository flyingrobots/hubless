#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
HOOK_DIR="$ROOT_DIR/.githooks"

if [ ! -d "$HOOK_DIR" ]; then
  echo "No .githooks directory found" >&2
  exit 1
fi

git config core.hooksPath "$HOOK_DIR"
echo "Git hooks path set to $HOOK_DIR"

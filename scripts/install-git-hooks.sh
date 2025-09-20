#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
HOOK_DIR="$ROOT_DIR/.githooks"

if [ ! -d "$HOOK_DIR" ]; then
  echo "No .githooks directory found" >&2
  exit 1
fi

if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  echo "Not inside a git repository" >&2
  exit 1
fi
git config --local core.hooksPath "$HOOK_DIR"
echo "Git hooks path set to $HOOK_DIR"

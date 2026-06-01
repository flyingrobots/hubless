#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BIN="${MARKDOWN_TRANSCLUSION_BIN:-markdown-transclusion}"
SCRIPT_ENV="${MARKDOWN_TRANSCLUSION_SCRIPT:-}"
BASE_OVERRIDE="${MARKDOWN_TRANSCLUSION_BASE:-}"

cd "$ROOT_DIR"

if ! command -v "$BIN" >/dev/null 2>&1; then
    echo "ERROR: transclusion bin '$BIN' not found on PATH" >&2
    exit 127
fi

cmd=(go run ./cmd/docs-components --repo "$ROOT_DIR" --transclusion-bin "$BIN")

if [[ -n "$SCRIPT_ENV" ]]; then
    cmd+=(--transclusion-args "$SCRIPT_ENV")
fi

if [[ -n "$BASE_OVERRIDE" ]]; then
    cmd+=(--transclusion-base "$BASE_OVERRIDE")
fi

exec "${cmd[@]}"

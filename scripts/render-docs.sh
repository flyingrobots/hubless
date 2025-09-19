#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BIN="${MARKDOWN_TRANSCLUSION_BIN:-markdown-transclusion}"
BASE_OVERRIDE="${MARKDOWN_TRANSCLUSION_BASE:-}"
ARGS_ENV="${MARKDOWN_TRANSCLUSION_ARGS:-}"

cd "$ROOT_DIR"

cmd=(go run ./cmd/docs-components --repo "$ROOT_DIR" --transclusion-bin "$BIN")

if [[ -n "$BASE_OVERRIDE" ]]; then
    cmd+=(--transclusion-base "$BASE_OVERRIDE")
fi

if [[ -n "$ARGS_ENV" ]]; then
    # shellcheck disable=SC2206
    extra_args=($ARGS_ENV)
    for arg in "${extra_args[@]}"; do
        cmd+=(--transclusion-args "$arg")
    done
fi

exec "${cmd[@]}"

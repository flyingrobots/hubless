#!/usr/bin/env bash
set -euo pipefail

cd /workspace

# Ensure git allows operations inside mounted workspace
if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  git config --global --add safe.directory /workspace
fi

export PATH="/usr/local/go/bin:$PATH"
export GOTOOLCHAIN=local

make fmt-check
make lint
make vet
make test
make docs
make docs-verify

echo "All CI checks completed successfully inside the container."

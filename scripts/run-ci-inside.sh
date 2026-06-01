#!/usr/bin/env bash
set -euo pipefail

cd /workspace

# Ensure git allows operations inside mounted workspace.
git config --global --add safe.directory /workspace

export PATH="/usr/local/go/bin:$PATH"
export GOTOOLCHAIN=local

make fmt-check
make lint
make vet
make test
make docs
make docs-verify

echo "All CI checks completed successfully inside the container."

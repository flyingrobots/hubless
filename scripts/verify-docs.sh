#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TARGETS=(
  "@hubless/issues/generated"
  "@hubless/roadmap/generated"
  "docs/reference/release-notes.md"
  "CHANGELOG.md"
)

missing=()
failures=()

# contains_placeholder returns 0 if the specified file contains an unresolved placeholder of the form `![[...]]`; otherwise returns 1.
contains_placeholder() {
  local file="$1"
  if rg -n "!\\[\\[[^]]+\\]\\]" "$file" >/dev/null; then
    return 0
  fi
  return 1
}

for target in "${TARGETS[@]}"; do
  path="$ROOT_DIR/$target"
  if [ ! -e "$path" ]; then
    missing+=("$target")
    continue
  fi
  if [ -d "$path" ]; then
    while IFS= read -r -d '' file; do
      if contains_placeholder "$file"; then
        failures+=("${file#$ROOT_DIR/}")
      fi
    done < <(find "$path" -type f -name '*.md' -print0)
  else
    if contains_placeholder "$path"; then
      failures+=("${target}")
    fi
  fi
done

if ((${#missing[@]} > 0)); then
  printf 'verify-docs: missing generated targets:\n'
  printf '  %s\n' "${missing[@]}"
  exit 1
fi

if ((${#failures[@]} > 0)); then
  printf 'verify-docs: unresolved placeholders found in:\n'
  printf '  %s\n' "${failures[@]}"
  exit 1
fi

echo "Docs verification passed."

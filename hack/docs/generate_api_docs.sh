#!/bin/bash

set -euo pipefail

ROOT_PATH="${1:-.}"

echo "Generating API docs with gomarkdoc..."

find "$ROOT_PATH" -type d | while read -r dir; do
  if ls "$dir"/*.go >/dev/null 2>&1 && go list "$dir" >/dev/null 2>&1; then
    echo generating static package documentation at "${dir##*/}"...
    ~/go/bin/gomarkdoc -o docs/docs/api/packages/"${dir##*/}".md "${dir}"
  fi
done

#!/bin/bash

set -euo pipefail

ROOT_PATH="${1:-.}"

find "$ROOT_PATH" -type d | while read -r dir; do
  if ls "$dir"/*.go >/dev/null 2>&1; then
    if go list "$dir" >/dev/null 2>&1; then
      echo generating static package documentation at "${dir##*/}"...
      ~/go/bin/gomarkdoc -o docs/docs/api/packages/"${dir##*/}".md $dir
    fi
  fi
done

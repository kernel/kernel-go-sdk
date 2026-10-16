#!/usr/bin/env bash
# Reads NUL-delimited kernel/cli paths on stdin and prints, NUL-delimited,
# the ones the CLI coverage updater must not publish.
set -euo pipefail

while IFS= read -r -d '' path; do
  case "$path" in
    go.mod | go.sum | cmd/*.go | pkg/*.go) continue ;;
  esac
  printf '%s\0' "$path"
done

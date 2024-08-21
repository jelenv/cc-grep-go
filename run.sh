#!/bin/sh

set -e

# compile
(
  cd "$(dirname "$0")"
  go build -o build/mygrep cmd/mygrep/main.go
)

# test or run
if [ "$1" = "--test" ]; then
  shift
  go test cmd/mygrep/main_test.go "$@"
else
  # Run the mygrep program with provided arguments
  exec build/mygrep "$@"
fi

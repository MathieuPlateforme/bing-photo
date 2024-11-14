#!/bin/sh
set -e

if [ "$RUN_TESTS" = "true" ]; then
  if command -v go >/dev/null 2>&1; then
    echo "Running tests..."
    go test -v ./...
  else
    echo "Go environment not available, skipping tests."
  fi
fi

echo "Starting the application..."
exec "$@"
